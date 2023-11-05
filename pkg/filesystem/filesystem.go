package filesystem;

import (
    "golang.org/x/net/context"
    "path/filepath"
    "fmt"
    "io"
    "net"
    "os"
    "log"
    "google.golang.org/grpc"
    "sort"
    "hash/fnv"
)

var TempDirectory string
// stores the stubs used for gRPC methods
var MachineStubs map[string] FileSystemClient = make(map[string]FileSystemClient)
var MachineIds []string = []string{}
var ThisMachineIdIdx int

var Files []string = []string{}
var ReplicaFiles []string = []string{}

type Server struct {
    UnimplementedFileSystemServer
}

func GetFileOwner(filename string) int {
    hasher := fnv.New32a()

    // Calculate filename's hash
    hasher.Write([]byte(filename))
    filenameHash := hasher.Sum32()
    hasher.Reset()

    // Search for the first node who's hash is greater than or equal to file's hash
    return (sort.Search(len(MachineIds), func(i int) bool {
        hasher.Write([]byte(MachineIds[i]))
        machineIdsIHash := hasher.Sum32()
        hasher.Reset()
        return machineIdsIHash >= filenameHash
    }) % len(MachineIds))
}

func InitializeGRPCConnection(machineId string, serverAddress string) {
    // Establish TCP connection with new node
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
        fmt.Errorf("grpc.Dial: %v\n", err)
	}

	client := NewFileSystemClient(conn)
	MachineStubs[machineId] = client
}

func (s *Server) Get(in *GetRequest, stream FileSystem_GetServer) error {
    // Open file
    filename := filepath.Join(TempDirectory, in.SdfsName)
    file, err := os.Open(filename)
    if err != nil {
        fmt.Errorf("os.Open: %v\n", err)
    }

    // Loop over bytes in file and send over the network
    for {
        buffer := make([]byte, 4096)
        bytesRead, err := file.Read(buffer)

        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Errorf("file.Read: %v\n", err)
        }

        resp := GetResponse{ Payload: string(buffer[:bytesRead]) }
        stream.Send(&resp)
    }

    err = stream.Send(&GetResponse{}) // do I need this to end the connection for RPC?
    if err != nil {
        fmt.Println(err)
    }
    return err
}

func (s *Server) Put(stream FileSystem_PutServer) error {
    initialized := false
    var sdfsFilename, filename string
    var replica bool
    var file *os.File

    // Loop over bytes from the network and write to file
    for {
        req, err := stream.Recv()

        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Errorf("stream.Recv: %v\n", err)
        }

        if !initialized {
            // Initialize request variables
            sdfsFilename = req.SdfsName
            replica = req.Replica
            filename = filepath.Join(TempDirectory, sdfsFilename)
            file, err = os.Create(filename)

            if err != nil {
                fmt.Errorf("os.Create: %v\n", err)
                return err
            }

            initialized = true
        }

        file.Write([]byte(req.PayloadToWrite))
    }

    file.Close()

    // Add filename to local file list
    if !replica {
        Files = append(Files, sdfsFilename)
    } else {
        ReplicaFiles = append(ReplicaFiles, sdfsFilename)
    }

    if !replica {
        // Send file to other machines as a replica
        for i := 1; i < 4; i++ {
            idx := (i + ThisMachineIdIdx) % len(MachineIds)
            if idx == ThisMachineIdIdx {
                break
            }

            Put(MachineStubs[MachineIds[idx]], filename, sdfsFilename, true)
        }
    }

    return stream.SendAndClose(&PutResponse{ Err: 0 })
}

func Put(targetStub FileSystemClient, localFilename string, sdfsFilename string, replica bool) {
    // Open file
    file, err := os.Open(localFilename)
    if err != nil {
        fmt.Errorf("os.Open: %v\n", err)
        return
    }

    stream, err := targetStub.Put(context.Background())
    if err != nil {
        fmt.Errorf("targetStub.Put: %v\n", err)
        return
    }

    // Loop over bytes in file and send over the network
    for {
        buffer := make([]byte, 4096)
        bytesRead, err := file.Read(buffer)

        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Errorf("file.Read: %v\n", err)
        }

        req := PutRequest{ PayloadToWrite: string(buffer[:bytesRead]), SdfsName: sdfsFilename, Replica: replica }
        stream.Send(&req)
    }

    // Close stream and receive
    if _, err = stream.CloseAndRecv(); err != nil {
        fmt.Errorf("stream.CloseAndRecv: %v\n", err)
    }
}

func Delete(targetStub FileSystemClient, sdfsFilename string, replica bool) {
    request := DeleteRequest{ SdfsName: sdfsFilename, Replica: replica }
    _, err := targetStub.Delete(context.Background(), &request)
    if err != nil {
        fmt.Errorf("targetStub.Read: %v\n", err)
    }
}

func remove(slice []string, element string) []string {
    newSlice := []string{}
    for _, s := range slice {
        if s != element {
            newSlice = append(newSlice, s)
        }
    }
    return newSlice
}

func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*DeleteResponse, error) {
    sdfsFilename, replica := in.SdfsName, in.Replica
    filename := filepath.Join(TempDirectory, sdfsFilename) 
    if err := os.Remove(filename); err != nil {
        fmt.Errorf("os.Remove: %v\n", err)
        return &DeleteResponse{}, err
    }

    // Remove filename from local file list
    if !replica {
        Files = remove(Files, sdfsFilename)
    } else {
        ReplicaFiles = remove(ReplicaFiles, sdfsFilename)
    }

    if !replica {
        // Tell other machines to delete file as a replica
        for i := 1; i < 4; i++ {
            idx := (i + ThisMachineIdIdx) % len(MachineIds)
            if idx == ThisMachineIdIdx {
                break
            }

            Delete(MachineStubs[MachineIds[idx]], sdfsFilename, true)
        }
    }

    return &DeleteResponse{}, nil
}

// should be called instead of directly calling the RPC
func Get(targetStub FileSystemClient, sdfsFilename string, localFilename string) {
    request := GetRequest{ SdfsName: sdfsFilename }
    // implement location logic here; going to just go to the first entry in our map for noA
    file, err := os.Create(localFilename) // this is standard open - check if we need destructive write or something
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()
    stream, err := targetStub.Get(context.Background(), &request)
    if err != nil {
        log.Fatal(err)
    }
    for {
        response, err := stream.Recv()
        if err == io.EOF { // it's done sending data to client
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        // write to the file
        file.Write([]byte(response.Payload))
    }
}

func (s *Server) FileRange(ctx context.Context, in *FileRangeRequest) (*FileRangeResponse, error) {
    hasher := fnv.New32a()
    var sdfsNames []string

    for _, file := range Files {
        hasher.Write([]byte(file))
        fileHash := hasher.Sum32()
        hasher.Reset()

        if (in.Start < in.End && in.Start <= fileHash && fileHash < in.End) || (in.End < in.Start && (in.Start <= fileHash || fileHash < in.End)) {
            sdfsNames = append(sdfsNames, file)
        }
    }

    for _, file := range ReplicaFiles {
        hasher.Write([]byte(file))
        fileHash := hasher.Sum32()
        hasher.Reset()

        if (in.Start < in.End && in.Start <= fileHash && fileHash < in.End) || (in.End < in.Start && (in.Start <= fileHash || fileHash < in.End)) {
            sdfsNames = append(sdfsNames, file)
        }
    }

    return &(FileRangeResponse{ SdfsNames: sdfsNames }), nil
}

func FileRange(targetStub FileSystemClient, start uint32, end uint32) []string {
    request := &FileRangeRequest{ Start: start, End: end }

    response, err := targetStub.FileRange(context.Background(), request)
    if err != nil {
        fmt.Errorf("client.FileRange: %v", err)
        return []string{}
    }

    return response.SdfsNames
}

func InitializeFileSystem(port string) {
    TempDirectory, _ = os.MkdirTemp(".", "tmp")
    InitializeGRPCServer(port)
}

func InitializeGRPCServer(port string) {
    lis, err := net.Listen("tcp", ":" + port)
    if err != nil {
        fmt.Errorf("net.Listen: %v\n", err)
        os.Exit(1)
    }

    grpcServer := grpc.NewServer()
    serv := Server{}
    RegisterFileSystemServer(grpcServer, &serv)
    grpcServer.Serve(lis)
}

