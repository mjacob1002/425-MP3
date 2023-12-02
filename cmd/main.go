package main;

import (
    "flag"
    "fmt"
    "strconv"
    "time"
    "sort"
    "hash/fnv"
    "path/filepath"
    "os"
    membership "github.com/mjacob1002/425-MP3/pkg/membership"
    "github.com/mjacob1002/425-MP3/pkg/cli"
	fs "github.com/mjacob1002/425-MP3/pkg/filesystem"
)

var thisMachineName string
var thisMachineId string
var recentlyAdded bool = true

func deleteAllFilesInRange(startNode string, endNode string) {
    newFiles := []string{}

    hasher := fnv.New32a()

    hasher.Write([]byte(startNode))
    start := hasher.Sum32()
    hasher.Reset()

    hasher.Write([]byte(endNode))
    end := hasher.Sum32()
    hasher.Reset()

    fmt.Printf("start: %v, end: %v\n", startNode, endNode)

    for _, file := range fs.Files {
        hasher.Write([]byte(file))
        fileHash  := hasher.Sum32()
        hasher.Reset()

        if (start < end && start < fileHash && fileHash <= end) || (end <= start && (start < fileHash || fileHash <= end)) {
            filename := filepath.Join(fs.TempDirectory, file)
            fmt.Printf("deleting file %v\n", file)
            if err := os.Remove(filename); err != nil {
                fmt.Printf(fmt.Errorf("os.Remove: %v\n", err).Error())
            }
        } else {
            newFiles = append(newFiles, file)
        }
    }
    fs.Files = newFiles
    fmt.Printf("keeping files %v\n", newFiles)
}

func onAdd(machineId string, serverAddress string) {
    fmt.Println("Adding new node to membership list:", machineId)

    fs.MachineIdsLock.Lock()
    defer fs.MachineIdsLock.Unlock()

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

    fs.InitializeGRPCConnection(machineId, serverAddress)

    defer func() {
        // Append new machine id to list
        fs.MachineIds = append(fs.MachineIds[:index], append([]string{machineId}, fs.MachineIds[index:]...)...)
        if index <= fs.ThisMachineIdIdx {
            fs.ThisMachineIdIdx++
        }
    }()

    ahead := (index - fs.ThisMachineIdIdx + len(fs.MachineIds)) % len(fs.MachineIds)
    behind := (fs.ThisMachineIdIdx - index + len(fs.MachineIds)) % len(fs.MachineIds)

    if recentlyAdded {
        return
    }

    if len(fs.MachineIds) < 4 {
        fmt.Printf("line 47\n")
        // We need to copy files around to ensure we have 3 replicas of files
        sdfsFilenames := fs.FileRangeNodes(fs.MachineIds[(fs.ThisMachineIdIdx - 1 + len(fs.MachineIds)) % len(fs.MachineIds)], fs.MachineIds[fs.ThisMachineIdIdx])

        for _, sdfsFilename := range sdfsFilenames {
            fs.Put(fs.MachineStubs[machineId], filepath.Join(fs.TempDirectory, sdfsFilename), sdfsFilename, true)
        }

        return
    }

    if behind == 3 {
        // remove (x - 4, new]
        deleteAllFilesInRange(fs.MachineIds[(fs.ThisMachineIdIdx - 4 + len(fs.MachineIds)) % len(fs.MachineIds)], machineId)
    } else if behind <= 2 {
        // remove (x - 4, x - 3]
        deleteAllFilesInRange(fs.MachineIds[(fs.ThisMachineIdIdx - 4 + len(fs.MachineIds)) % len(fs.MachineIds)], fs.MachineIds[(fs.ThisMachineIdIdx - 3 + len(fs.MachineIds)) % len(fs.MachineIds)])
    }

    if fs.ThisMachineIdIdx == index {
        // copy (x - 1, new] to new
        sdfsFilenames := fs.FileRangeNodes(fs.MachineIds[(fs.ThisMachineIdIdx - 1 + len(fs.MachineIds)) % len(fs.MachineIds)], machineId)

        for _, sdfsFilename := range sdfsFilenames {
            fs.Put(fs.MachineStubs[machineId], filepath.Join(fs.TempDirectory, sdfsFilename), sdfsFilename, true)
        }
    }

    if ahead > 0 && ahead <= 3 {
        // copy (x - 1, x] to new
        sdfsFilenames := fs.FileRangeNodes(fs.MachineIds[(fs.ThisMachineIdIdx - 1 + len(fs.MachineIds)) % len(fs.MachineIds)], fs.MachineIds[fs.ThisMachineIdIdx])

        for _, sdfsFilename := range sdfsFilenames {
            fs.Put(fs.MachineStubs[machineId], filepath.Join(fs.TempDirectory, sdfsFilename), sdfsFilename, true)
        }
    }
}

func onDelete(machineId string) {
    fmt.Println("Deleting node from membership list:", machineId)

    fs.MachineIdsLock.Lock()
    defer fs.MachineIdsLock.Unlock()

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

    if len(fs.MachineIds) > 4 && (fs.ThisMachineIdIdx + len(fs.MachineIds) - index) % len(fs.MachineIds) <= 4  {
        // We need to copy files around to ensure we have 3 replicas of files

        // Check all 4 machines that occur previously in the ring
        for offset := 4; offset > 0; offset-- {
            hasher.Write([]byte(fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - 5) % len(fs.MachineIds)]))
            start := hasher.Sum32()
            hasher.Reset()

            hasher.Write([]byte(fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - 4) % len(fs.MachineIds)]))
            end := hasher.Sum32()
            hasher.Reset()

            newFiles, err := fs.FileRange(fs.MachineStubs[fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - offset) % len(fs.MachineIds)]], start, end)
            if err == nil {
                for _, newFile := range newFiles {
                    err := fs.Get(fs.MachineStubs[fs.MachineIds[(fs.ThisMachineIdIdx + len(fs.MachineIds) - offset) % len(fs.MachineIds)]], newFile, filepath.Join(fs.TempDirectory, newFile)) 
                    if err != nil {
                        fmt.Printf("fs.Get: %v\n", err)
                    } else {
                        fs.Files = append(fs.Files, newFile)
                    }
                }
                break
            }
        }
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

    go func() {
        time.Sleep(5 * time.Second)
        recentlyAdded = false
    }()
	fmt.Println(fs.TempDirectory)
    go cli.ListenToCommands()

	select {}
}

