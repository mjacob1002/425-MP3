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
	fs "github.com/mjacob1002/425-MP3/pkg/filesystem"
)

var thisMachineName string
var thisMachineId string

func onAdd(machineId string, serverAddress string) {
    fmt.Println("Adding new node to membership list:", machineId)

    hasher := fnv.New32a()

    // Calculate new machine's id's hash
    hasher.Write([]byte(machineId))
    machineIdHash := hasher.Sum32()
    hasher.Reset()

    // Search for new location of the new machine id
    index := sort.Search(len(fs.MachineIds), func(i int) bool {
        hasher.Write([]byte(fs.MachineIds[i]))
        machineIdsIHash := hasher.Sum32()
        hasher.Reset()
		return machineIdsIHash >= machineIdHash
	})

    // Append new machine id to list
	fs.MachineIds = append(fs.MachineIds[:index], append([]string{machineId}, fs.MachineIds[index:]...)...)
    if index <= fs.ThisMachineIdIdx {
        fs.ThisMachineIdIdx++
    }

    fs.InitializeGRPCConnection(machineId, serverAddress)
}

func onDelete(machineId string) {
    fmt.Println("Deleting node from membership list:", machineId)

    hasher := fnv.New32a()

    // Calculate old machine's id's hash
    hasher.Write([]byte(machineId))
    machineIdHash := hasher.Sum32()
    hasher.Reset()

    // Search for location of the old machine id
    index := sort.Search(len(fs.MachineIds), func(i int) bool {
        hasher.Write([]byte(fs.MachineIds[i]))
        machineIdsIHash := hasher.Sum32()
        hasher.Reset()
		return machineIdsIHash >= machineIdHash
	})

    if fs.ThisMachineIdIdx == (index + 1) % len(fs.MachineIds) {
        hasher.Write([]byte(fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - 5) % len(fs.MachineIds)]))
        start := hasher.Sum32()
        hasher.Reset()

        hasher.Write([]byte(fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - 4) % len(fs.MachineIds)]))
        end := hasher.Sum32()
        hasher.Reset()

        fmt.Printf("querying from %v:\n", fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - 4) % len(fs.MachineIds)])

        newReplicas := fs.FileRange(fs.MachineStubs[fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - 4) % len(fs.MachineIds)]], start, end)

        fmt.Printf("newReplicas: %v\n", newReplicas)
    }

    // Remove old machine id to list
    fs.MachineIds = append(fs.MachineIds[:index], fs.MachineIds[index+1:]...)
    if index <= fs.ThisMachineIdIdx {
        fs.ThisMachineIdIdx--
    }

    // Delete old connection from stubs map
    delete(fs.MachineStubs, machineId)
}

func main() {
    // Collect arguments
    var hostname, port, introducer, applicationPort string
    flag.StringVar(&thisMachineName, "machine_name", "", "Machine Name")
    flag.StringVar(&hostname, "hostname", "", "Hostname")
    flag.StringVar(&port, "port", "", "Port")
    flag.StringVar(&introducer, "introducer", "", "Introducer Node Address")
	flag.StringVar(&applicationPort, "application_port", "", "Application level port")
    flag.Parse()

    thisMachineId = (thisMachineName + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10))
    fs.MachineIds = append(fs.MachineIds, thisMachineId)
    fs.ThisMachineIdIdx = 0

    // Generate TCP connection to itself
    fs.InitializeGRPCConnection(thisMachineId, hostname + ":" + applicationPort)

	go fs.InitializeFileSystem(applicationPort)
	go membership.Join(
        thisMachineId,
        hostname,
        port,
        introducer,
        applicationPort,
        onAdd,
        onDelete,
    )
    go cli.ListenToCommands()

	select {}
}

