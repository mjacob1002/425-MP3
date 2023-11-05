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

var TempDirectory string;
// stores the stubs used for gRPC methods
var MachineStubs map[string] FileSystemClient = make(map[string]FileSystemClient)
var MachineIds []string = []string{}
var ThisMachineIdIdx int

var Files []string = []string{}
var ReplicaFiles []string = []string{}

type Server struct {
    UnimplementedFileSystemServer;
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
        fmt.Errorf("grpc.Dial: %v\n", err);
	}

	client := NewFileSystemClient(conn)
	MachineStubs[machineId] = client
}

func (s *Server) Get(in *GetRequest, stream FileSystem_GetServer) error {
    filename := filepath.Join(TempDirectory, in.SdfsName)
    file, err := os.Open(filename)
    if err != nil {
        fmt.Errorf("os.Open: %v\n", err)
    }

    for {
        buffer := make([]byte, 4096)
        bytesRead, err := file.Read(buffer)

        if err == io.EOF {
            break;
        } else if err != nil {
            fmt.Errorf("file.Read: %v\n", err)
        }

        resp := GetResponse{ Payload: string(buffer[:bytesRead]) }
        stream.Send(&resp);
    }

    err = stream.Send(&GetResponse{}); // do I need this to end the connection for RPC?
    if err != nil {
        fmt.Println(err);
    }
    return err;
}

func (s *Server) Put(stream FileSystem_PutServer) error {
    // acquire a lock here with the information - also incorporate the filestruct idea here
    // go through all of their messages and write them
    fname, sdfsname := "", ""
    replica := false
    var file *os.File;
    for {
        req, err := stream.Recv();
        if err == io.EOF {
            break;
        }
        if err != nil {
            log.Fatal(err);
        }
        if fname == "" {
            fname = filepath.Join(TempDirectory, req.SdfsName)
            tempvar, err := os.Create(fname);
            if err != nil {
                log.Fatal(err)
            }
            file = tempvar
            replica = req.Replica
            sdfsname = req.SdfsName
        }
        file.Write([]byte(req.PayloadToWrite));
    }
    file.Close()
    fmt.Printf("just wrote to the file in %v\n", TempDirectory)

    if !replica {
        Files = append(Files, sdfsname)
    } else {
        ReplicaFiles = append(ReplicaFiles, sdfsname)
    }

    if !replica {
        // send out to other machines
        for i := 1; i < 4; i++ {
            idx := (i + ThisMachineIdIdx) % len(MachineIds)
            if idx == ThisMachineIdIdx {
                break
            }
            // send to machine at idx
            fmt.Printf("should be sending out %v to index %v: %v\n", sdfsname, idx, MachineIds[idx])
            Put(MachineStubs[MachineIds[idx]], fname, sdfsname, true)
        }
    }
    return stream.SendAndClose(&PutResponse{Err: 0});
}

func Put(targetStub FileSystemClient, localFilename string, sdfsFilename string, replica bool) {
    file, err := os.Open(localFilename)
    if err != nil {
        fmt.Errorf("os.Open: %v\n", err)
        return
    }

    stream, err := targetStub.Put(context.Background())
    if err != nil {
        fmt.Errorf("FileSystemClient.Put: %v\n", err)
        return
    }

    // Loop over bytes in file and send over the network
    for {
        buffer := make([]byte, 4096)
        bytesRead, err := file.Read(buffer)

        if err == io.EOF {
            break;
        } else if err != nil {
            fmt.Errorf("file.Read: %v\n", err)
        }

        req := PutRequest{ PayloadToWrite: string(buffer[:bytesRead]), LocalName: localFilename, SdfsName: sdfsFilename, Replica: replica }
        stream.Send(&req);
    }

    // Close stream and receive
    if _, err = stream.CloseAndRecv(); err != nil {
        fmt.Errorf("file.Read: %v\n", err)
    }
}

func Delete(targetStub FileSystemClient, remotefname string){
    deletemsg := DeleteRequest{SdfsName: remotefname}
    resp, err := targetStub.Delete(context.Background(), &deletemsg);
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Delete response: %s\n", resp.String());
}

func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*DeleteResponse, error) {
    fmt.Println("Got invoked Delete by ", in.SdfsName);
    fname := filepath.Join(TempDirectory, in.SdfsName);
    err := os.Remove(fname);
    if err != nil {
        log.Fatal(err);
    }
    return &DeleteResponse{}, nil;
}

// should be called instead of directly calling the RPC
func Get(targetStub FileSystemClient, sdfsFilename string, localFilename string) {
    request := GetRequest{ SdfsName: sdfsFilename, LocalName: localFilename }
    // implement location logic here; going to just go to the first entry in our map for noA
    file, err := os.Create(localFilename) // this is standard open - check if we need destructive write or something
    if err != nil {
        fmt.Println(err);
    }
    defer file.Close()
    stream, err := targetStub.Get(context.Background(), &request);
    if err != nil {
        log.Fatal(err);
    }
    for {
        response, err := stream.Recv()
        if err == io.EOF { // it's done sending data to client
            break;
        }
        if err != nil {
            log.Fatal(err);
        }
        // write to the file
        file.Write([]byte(response.Payload));
    }
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

