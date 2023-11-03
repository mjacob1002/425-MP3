package main;

import (
    "flag"
    "fmt"
    "strconv"
	"log"
    "time"
    "sort"
    "hash/fnv"
	"io"
    membership "github.com/mjacob1002/425-MP3/pkg/membership"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	fs "github.com/mjacob1002/425-MP3/pkg/filesystem"
)

var thisMachineName string
var thisMachineId string
var machineIds []string = []string{}
// stores the stubs used for gRPC methods
var machineStubs map[string] fs.FileSystemClient = make(map[string]fs.FileSystemClient)


func onAdd(machineId string, serverAddress string) {
    fmt.Println("Adding new node to membership list:", machineId)

    hasher := fnv.New32a()

    hasher.Write([]byte(machineId))
    machineIdHash := hasher.Sum32()
    hasher.Reset()

    index := sort.Search(len(machineIds), func(i int) bool {
        hasher.Write([]byte(machineIds[i]))
        machineIdsIHash := hasher.Sum32()
        hasher.Reset()
		return machineIdsIHash >= machineIdHash
	})

	machineIds = append(machineIds[:index], append([]string{machineId}, machineIds[index:]...)...)
	fmt.Println("Trying to connect to ", serverAddress);
	conn, err := grpc.Dial(serverAddress)
	if err != nil {
		fmt.Println(err);
	}
	client := fs.NewFileSystemClient(conn)
	// we should lock the machineStubs map
	machineStubs[machineId] = client;
	// unlock mutex
	fmt.Println("Just added the following's gRPC stuff: ", machineId, " with clientStub of ", machineStubs[machineId], " but client is ", client);
}

func onDelete(machineId string) {
    fmt.Println("Deleting node from membership list:", machineId)

    hasher := fnv.New32a()

    hasher.Write([]byte(machineId))
    machineIdHash := hasher.Sum32()
    hasher.Reset()

    index := sort.Search(len(machineIds), func(i int) bool {
        hasher.Write([]byte(machineIds[i]))
        machineIdsIHash := hasher.Sum32()
        hasher.Reset()
		return machineIdsIHash >= machineIdHash
	})

    machineIds = append(machineIds[:index], machineIds[index+1:]...)
}

func main() {
    // Collect arguments
    var hostname, port, introducer string
    flag.StringVar(&thisMachineName, "machine_name", "", "Machine Name")
    flag.StringVar(&hostname, "hostname", "", "Hostname")
    flag.StringVar(&port, "port", "", "Port")
    flag.StringVar(&introducer, "introducer", "", "Introducer Node Address")
	flag.Int64Var(&fs.Tcp_port, "tcp_port", 9999, "The port where gRPC will be listening for incoming requests")
    flag.Parse()

    thisMachineId = (thisMachineName + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10))
    machineIds = append(machineIds, thisMachineId)
	fmt.Printf("I am using tcp_port %d\n", fs.Tcp_port);
	go fs.InitializeGRPCServer()
	go membership.Join(
        thisMachineId,
        hostname,
        port,
        introducer,
        onAdd,
        onDelete,
    )

    // select {}
	time.Sleep(20 * time.Second)
	fmt.Println("Time to test RPC..");
	for key, value := range(machineStubs){
		fmt.Println("Working with key=", key);
		req := fs.GetRequest{SdfsName: "sdfsfile", LocalName: "localfname"} // create a request
		stream, err := value.Get(context.Background(), &req);
		if err != nil {
			fmt.Println(err);
		}
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err);
			}
			fmt.Printf("Response: %s\n", resp.String())
		}
	}

}
