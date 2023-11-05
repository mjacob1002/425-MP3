package main;

import (
    "flag"
    "fmt"
    "strconv"
    "time"
    "sort"
    "hash/fnv"
    membership "github.com/mjacob1002/425-MP3/pkg/membership"
    "github.com/mjacob1002/425-MP3/pkg/cli"
	"google.golang.org/grpc"
	//"golang.org/x/net/context"
	fs "github.com/mjacob1002/425-MP3/pkg/filesystem"
)

var thisMachineName string
var thisMachineId string


func onAdd(machineId string, serverAddress string) {
    fmt.Println("Adding new node to membership list:", machineId)

    hasher := fnv.New32a()

    hasher.Write([]byte(machineId))
    machineIdHash := hasher.Sum32()
    hasher.Reset()

    index := sort.Search(len(fs.MachineIds), func(i int) bool {
        hasher.Write([]byte(fs.MachineIds[i]))
        machineIdsIHash := hasher.Sum32()
        hasher.Reset()
		return machineIdsIHash >= machineIdHash
	})

	fs.MachineIds = append(fs.MachineIds[:index], append([]string{machineId}, fs.MachineIds[index:]...)...)
    if index <= fs.MachineIdsIdx {
        fs.MachineIdsIdx++
    }
	fmt.Println("Trying to connect to ", serverAddress);
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err);
	}
	if conn == nil {
		fmt.Println("What the fuck - why is this shit dead.");
	}
	client := fs.NewFileSystemClient(conn)
	// we should lock the machineStubs map
	fs.MachineStubs[machineId] = client;
	// unlock mutex
	fmt.Println("Just added the following's gRPC stuff: ", machineId, " with clientStub of ", fs.MachineStubs[machineId], " but client is ", client);
}

func onDelete(machineId string) {
    fmt.Println("Deleting node from membership list:", machineId)

    hasher := fnv.New32a()

    hasher.Write([]byte(machineId))
    machineIdHash := hasher.Sum32()
    hasher.Reset()

    index := sort.Search(len(fs.MachineIds), func(i int) bool {
        hasher.Write([]byte(fs.MachineIds[i]))
        machineIdsIHash := hasher.Sum32()
        hasher.Reset()
		return machineIdsIHash >= machineIdHash
	})

    fs.MachineIds = append(fs.MachineIds[:index], fs.MachineIds[index+1:]...)
    if index <= fs.MachineIdsIdx {
        fs.MachineIdsIdx--
    }
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
    fs.MachineIds = append(fs.MachineIds, thisMachineId)
    fs.MachineIdsIdx = 0
	fmt.Printf("I am using tcp_port %d\n", fs.Tcp_port);
	go fs.InitializeFileSystem()
	go membership.Join(
        thisMachineId,
        hostname,
        port,
        introducer,
        onAdd,
        onDelete,
    )
    go cli.ListenToCommands()

	select {}
}
