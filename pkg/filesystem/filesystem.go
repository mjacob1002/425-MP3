package filesystem;

import (
    "golang.org/x/net/context"
    "path/filepath"
    "fmt"
    "io"
	"math/rand"
    "net"
    "os"
    "log"
    "google.golang.org/grpc"
    "sort"
	"strconv"
	"sync"
    "hash/fnv"
    "google.golang.org/protobuf/types/known/emptypb"
    "time"
)

var TempDirectory string
// stores the stubs used for gRPC methods
var MachineStubs map[string] FileSystemClient = make(map[string]FileSystemClient)
var MachineIdsLock sync.Mutex
var MachineIds []string = []string{}
var ThisMachineIdIdx int

// TODO: try and phase these out this iteration...
var Files []string = []string{}

var FilesMapMutex sync.Mutex // used for mutual exclusion in modifying this shit
var FilesMap map[string]*FileStruct = make(map[string]*FileStruct);

func SearchForFile(sdfsName string) (*FileStruct, int64) {
	FilesMapMutex.Lock()
	defer FilesMapMutex.Unlock();
	// have mutex to search the thingy
	fileStruct, ok := FilesMap[sdfsName];
	if !ok {
		return nil, 1;
	}
	return fileStruct, 0
}

func CreateFileStruct(sdfsName string) *FileStruct{
	FilesMapMutex.Lock(); // get mutex to the thingy
	defer FilesMapMutex.Unlock()
	FilesMap[sdfsName] = NewFileStruct(sdfsName, -1) // TODO: invoke constructor to add the replica to it
	return FilesMap[sdfsName]
}

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
        fmt.Printf(fmt.Errorf("grpc.Dial: %v\n", err).Error())
	}

	client := NewFileSystemClient(conn)
	MachineStubs[machineId] = client
}

func (s *Server) Get(in *GetRequest, stream FileSystem_GetServer) error {
	// do the mutex locking garbage
	fileStruct, status := SearchForFile(in.SdfsName)
	if status == 1 {
		// what do I do in the case of a Getting a non-existent file
		fmt.Printf("%s does not exist on the local filesystem...\n", in.SdfsName);
	}
	fileStruct.RLock();
	defer fileStruct.RUnlock();
    // Open file
    filename := filepath.Join(TempDirectory, in.SdfsName)
    file, err := os.Open(filename)
    if err != nil {
		fmt.Println("failure here");
        fmt.Printf(fmt.Errorf("os.Create: %v\n", err).Error())
        return err
    }

    // Loop over bytes in file and send over the network
    for {
        buffer := make([]byte, 4096)
        bytesRead, err := file.Read(buffer)

        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf(fmt.Errorf("file.Read: %v\n", err).Error())
            return err
        }

        resp := GetResponse{ Payload: string(buffer[:bytesRead]) }
        stream.Send(&resp)
    }

    err = stream.Send(&GetResponse{}) // do I need this to end the connection for RPC?
    if err != nil {
        fmt.Println(err)
        return err
    }

    return nil
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
            fmt.Printf(fmt.Errorf("stream.Recv: %v\n", err).Error())
            return err
        }

        if !initialized {
            // Initialize request variables
            sdfsFilename = req.SdfsName
            replica = req.Replica
            filename = filepath.Join(TempDirectory, sdfsFilename)
			// before, search the filelist map for the corresponding filestruct
			fileStruct, status := SearchForFile(sdfsFilename)
			if status == 1 {
				// have to create a block for this
				fileStruct = CreateFileStruct(sdfsFilename) // modify this function to set the field
			}
			// get writing access
			fileStruct.WLock();
			defer fileStruct.WUnlock();
			fileStruct.deleted = false; // this file is no longer deleted
            file, err = os.Create(filename)

            if err != nil {
                fmt.Printf(fmt.Errorf("os.Create: %v\n", err).Error())
                return err
            }

            initialized = true
        }

        file.Write([]byte(req.PayloadToWrite))
    }

    file.Close()

    // Add filename to local file list
    Files = append(Files, sdfsFilename)

    if !replica {
        // Send file to other machines as a replica
        for i := 1; i < 4; i++ {
            idx := (i + ThisMachineIdIdx) % len(MachineIds)
            if idx == ThisMachineIdIdx {
                break
            }
            if err := Put(MachineStubs[MachineIds[idx]], filename, sdfsFilename, true); err != nil {
                fmt.Printf(fmt.Errorf("Put: %v\n", err).Error())
            }
        }
    }

    return nil
}

func Put(targetStub FileSystemClient, localFilename string, sdfsFilename string, replica bool) error {
    // Open file
    file, err := os.Open(localFilename)
    if err != nil {
        fmt.Printf(fmt.Errorf("os.Open: %v\n", err).Error())
        return err
    }

    contextWithTimeout, cancel := context.WithTimeout(context.Background(), 300 * time.Second)
    defer cancel()

    stream, err := targetStub.Put(contextWithTimeout)
    if err != nil {
        fmt.Printf(fmt.Errorf("targetStub.Put: %v\n", err).Error())
        return err
    }

    // Loop over bytes in file and send over the network
    for {
        buffer := make([]byte, 4096)
        bytesRead, err := file.Read(buffer)

        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf(fmt.Errorf("file.Read: %v\n", err).Error())
        }

        req := PutRequest{ PayloadToWrite: string(buffer[:bytesRead]), SdfsName: sdfsFilename, Replica: replica }
        stream.Send(&req)
    }

    // Close stream and receive
    if _, err = stream.CloseAndRecv(); err != nil {
        fmt.Printf(fmt.Errorf("stream.CloseAndRecv: %v\n", err).Error())
        return err
    }

    return nil
}

func Delete(targetStub FileSystemClient, sdfsFilename string, replica bool) error {
    contextWithTimeout, cancel := context.WithTimeout(context.Background(), 300 * time.Second)
    defer cancel()

    request := DeleteRequest{ SdfsName: sdfsFilename, Replica: replica }
    fmt.Printf("%s\n", request.String())
    _, err := targetStub.Delete(contextWithTimeout, &request)
    if err != nil {
        fmt.Printf(fmt.Errorf("targetStub.Read: %v\n", err).Error())
        return err
    }

    return nil
}

func Remove(slice []string, element string) []string {
    newSlice := []string{}
    for _, s := range slice {
        if s != element {
            newSlice = append(newSlice, s)
        }
    }
    return newSlice
}

func InvokeARead(machine_id string, sdfsname string){
		res, err := MachineStubs[machine_id].InvokeRead(context.Background(), &InvokeReadRequest{SdfsName: sdfsname})
		if err != nil {
			fmt.Println(err);
		}
		fmt.Printf("%s read %s to create %s\n", machine_id, sdfsname, res.LocalName)
}

func (s *Server) InvokeRead(ctx context.Context, in *InvokeReadRequest) (*InvokeReadResponse, error) {
	sdfsName := in.SdfsName;
	randNum := rand.Intn(1000)
	randNumString := strconv.Itoa(randNum)
	localFilename := sdfsName + "_" + randNumString
	index  := GetFileOwner(sdfsName)
	fmt.Printf("about to ping %s for get...\n", MachineIds[index])
	err := Get(MachineStubs[MachineIds[index]], sdfsName, localFilename)
	if err != nil {
		fmt.Println(err);
	}
	return &InvokeReadResponse{LocalName: localFilename}, nil
}

func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*emptypb.Empty, error) {
    sdfsFilename, replica := in.SdfsName, in.Replica
	// get write access to this shit
	fileStruct, status := SearchForFile(sdfsFilename)
	if status == 1 {
			fmt.Printf("trying to delete %s, but the file isn't here...");
	}
	// get writing access
	fileStruct.WLock();
	defer fileStruct.WUnlock();
	fileStruct.deleted = true; // set the deleted field to true;
    filename := filepath.Join(TempDirectory, sdfsFilename)
    if err := os.Remove(filename); err != nil {
        fmt.Printf(fmt.Errorf("os.Remove: %v\n", err).Error())
        return nil, err
    }

    // Remove filename from local file list
    Files = Remove(Files, sdfsFilename)

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

    return &emptypb.Empty{}, nil
}

// should be called instead of directly calling the RPC
func Get(targetStub FileSystemClient, sdfsFilename string, localFilename string) error {
    request := GetRequest{ SdfsName: sdfsFilename }
    // implement location logic here; going to just go to the first entry in our map for noA
    file, err := os.Create(localFilename) // this is standard open - check if we need destructive write or something
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

    contextWithTimeout, cancel := context.WithTimeout(context.Background(), 300 * time.Second)
    defer cancel()

    stream, err := targetStub.Get(contextWithTimeout, &request)
    if err != nil {
        log.Fatal(err)
        return err
    }

    for {
        response, err := stream.Recv()
        if err == io.EOF { // it's done sending data to client
            break
        }
        if err != nil {
            log.Fatal(err)
            return err
        }
		fmt.Println("about to write payload...\n")
        // write to the file
        file.Write([]byte(response.Payload))
    }
    
    return nil
}

func (s *Server) FileRange(ctx context.Context, in *FileRangeRequest) (*FileRangeResponse, error) {
    sdfsNames := FileRangeHash(in.Start, in.End)

    return &(FileRangeResponse{ SdfsNames: sdfsNames }), nil
}

func FileRangeNodes(start string, end string) []string {
    hasher := fnv.New32a()

    hasher.Write([]byte(start))
    startHash := hasher.Sum32()
    hasher.Reset()
    
    hasher.Write([]byte(end))
    endHash := hasher.Sum32()
    hasher.Reset()

    return FileRangeHash(startHash, endHash)
}

func FileRangeHash(start uint32, end uint32) []string {
    hasher := fnv.New32a()
    var sdfsNames []string

    for _, file := range Files {
        hasher.Write([]byte(file))
        fileHash := hasher.Sum32()
        hasher.Reset()

        if (start < end && start <= fileHash && fileHash < end) || (end <= start && (start <= fileHash || fileHash < end)) {
            sdfsNames = append(sdfsNames, file)
        }
    }

    return sdfsNames
}

func FileRange(targetStub FileSystemClient, start uint32, end uint32) ([]string, error) {
    request := &FileRangeRequest{ Start: start, End: end }

    contextWithTimeout, cancel := context.WithTimeout(context.Background(), 300 * time.Second)
    defer cancel()

    response, err := targetStub.FileRange(contextWithTimeout, request)
    if err != nil {
        fmt.Printf(fmt.Errorf("client.FileRange: %v", err).Error())
        return nil, err
    }

    return response.SdfsNames, nil
}

func InitializeFileSystem(port string) {
    TempDirectory, _ = os.MkdirTemp(".", "tmp")
    InitializeGRPCServer(port)
}

func InitializeGRPCServer(port string) {
    lis, err := net.Listen("tcp", ":" + port)
    if err != nil {
        fmt.Printf(fmt.Errorf("net.Listen: %v\n", err).Error())
        os.Exit(1)
    }

    grpcServer := grpc.NewServer()
    serv := Server{}
    RegisterFileSystemServer(grpcServer, &serv)
    grpcServer.Serve(lis)
}

