package cli;

import (
    "bufio"
    "fmt"
    "os"
    "regexp"

	fs "github.com/mjacob1002/425-MP3/pkg/filesystem"
)

func ListenToCommands(){
    // Define command regular expressions
    putRe := regexp.MustCompile(`^put (\S+) (\S+)\n$`)
    getRe := regexp.MustCompile(`^get (\S+) (\S+)\n$`)
    deleteRe := regexp.MustCompile(`^delete (\S+)\n$`)
    lsRe := regexp.MustCompile(`^ls (\S+)\n$`)
    storeRe := regexp.MustCompile(`^store\n$`)

    listMemRe := regexp.MustCompile(`^list_mem\n$`)

    reader := bufio.NewReader(os.Stdin)
    for {
        // Read input
        input, _ := reader.ReadString('\n')

        switch {
        case putRe.MatchString(input):
            matches := putRe.FindStringSubmatch(input)
            if len(matches) == 3 {
                localFilename := matches[1]
                sdfsFilename := matches[2]
                index := fs.GetFileOwner(sdfsFilename)
                fs.Put(fs.MachineStubs[fs.MachineIds[index]], localFilename, sdfsFilename, false)
            }
        case getRe.MatchString(input):
            matches := getRe.FindStringSubmatch(input)
            if len(matches) == 3 {
                sdfsFilename := matches[1]
                localFilename := matches[2]
                index := fs.GetFileOwner(sdfsFilename)
                fs.Get(fs.MachineStubs[fs.MachineIds[index]], sdfsFilename, localFilename)
            }
        case deleteRe.MatchString(input):
            matches := deleteRe.FindStringSubmatch(input)
            if len(matches) == 2 {
                sdfsFilename := matches[1]
                index := fs.GetFileOwner(sdfsFilename)
                fs.Delete(fs.MachineStubs[fs.MachineIds[index]], sdfsFilename, false)
            }
        case lsRe.MatchString(input):
            matches := lsRe.FindStringSubmatch(input)
            if len(matches) == 2 {
                sdfsFilename := matches[1]
                index := fs.GetFileOwner(sdfsFilename)

                for i := 0; i < 4; i++ {
                    idx := (i + index) % len(fs.MachineIds)
                    if i != 0 && idx == index {
                        break
                    }

                    fmt.Printf("%v. %v\n", i + 1, fs.MachineIds[idx])
                }
            }
        case storeRe.MatchString(input):
            fmt.Printf("Files:\n")
            for i, file := range fs.Files {
                fmt.Printf("%v. %v\n", i + 1, file)
            }

            fmt.Printf("Replicated Files:\n")
            for i, file := range fs.ReplicaFiles {
                fmt.Printf("%v. %v\n", i + 1, file)
            }
        case listMemRe.MatchString(input):
            for i, machineId := range fs.MachineIds {
                fmt.Printf("%v. %v\n", i + 1, machineId)
            }
        default:
            fmt.Printf("Could not recoginize command\n")
        }
    }
}

