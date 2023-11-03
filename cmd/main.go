package main;

import (
    "flag"
    "fmt"
    "strconv"
    "time"
    "sort"
    "hash/fnv"
    membership "github.com/mjacob1002/425-MP3/pkg/membership"
)

var thisMachineName string
var thisMachineId string
var machineIds []string = []string{}

func onAdd(machineId string) {
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

    machineIds := append(machineIds[:index], machineIds[index+1:]...)
}

func main() {
    // Collect arguments
    var hostname, port, introducer string
    flag.StringVar(&thisMachineName, "machine_name", "", "Machine Name")
    flag.StringVar(&hostname, "hostname", "", "Hostname")
    flag.StringVar(&port, "port", "", "Port")
    flag.StringVar(&introducer, "introducer", "", "Introducer Node Address")
    flag.Parse()

    thisMachineId = (thisMachineName + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10))
    machineIds = append(machineIds, thisMachineId)

    go membership.Join(
        thisMachineId,
        hostname,
        port,
        introducer,
        onAdd,
        onDelete,
    )

    select {}
}
