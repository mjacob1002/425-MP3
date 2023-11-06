package cli;

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
	membership "github.com/mjacob1002/425-MP3/pkg/membership"
	fs "github.com/mjacob1002/425-MP3/pkg/filesystem"
)

func Multiread(sdfsName string) {
	counter := 0;
	machine_list := []string{}
	for machine, _ := range(membership.MembershipList) {
		if counter >= 6 {
			break;
		}
		machine_list = append(machine_list, machine);
		counter++;
	}
	fmt.Println("Nodes to ping: ", machine_list)
	for _, machine := range(machine_list){
		go fs.InvokeARead(machine,  sdfsName)
	}
}

func ListenToCommands(){
    // Define command regular expressions
    putRe := regexp.MustCompile(`^put (\S+) (\S+)\n$`)
    getRe := regexp.MustCompile(`^get (\S+) (\S+)\n$`)
    deleteRe := regexp.MustCompile(`^delete (\S+)\n$`)
    lsRe := regexp.MustCompile(`^ls (\S+)\n$`)
    storeRe := regexp.MustCompile(`^store\n$`)
	multiReadRe := regexp.MustCompile(`^multiread (\S+)\n$`)
    listMemRe := regexp.MustCompile(`^list_mem\n$`)

    reader := bufio.NewReader(os.Stdin)
    for {
        // Read input
        input, _ := reader.ReadString('\n')

        switch {
		case multiReadRe.MatchString(input):
			matches := multiReadRe.FindStringSubmatch(input)
			fmt.Println("Mutliread for file ", matches[1])
			Multiread(matches[1])
        case putRe.MatchString(input):
            matches := putRe.FindStringSubmatch(input)
            if len(matches) == 3 {
                localFilename := matches[1]
                sdfsFilename := matches[2]
                index := fs.GetFileOwner(sdfsFilename)
                fs.Put(fs.MachineStubs[fs.MachineIds[index]], localFilename, sdfsFilename, false)
                fmt.Printf("done")
            }
        case getRe.MatchString(input):
            matches := getRe.FindStringSubmatch(input)
            if len(matches) == 3 {
                sdfsFilename := matches[1]
                localFilename := matches[2]
                index := fs.GetFileOwner(sdfsFilename)
                fs.Get(fs.MachineStubs[fs.MachineIds[index]], sdfsFilename, localFilename)
                fmt.Printf("done")
            }
        case deleteRe.MatchString(input):
            matches := deleteRe.FindStringSubmatch(input)
            if len(matches) == 2 {
                sdfsFilename := matches[1]
                index := fs.GetFileOwner(sdfsFilename)
                fs.Delete(fs.MachineStubs[fs.MachineIds[index]], sdfsFilename, false)
                fmt.Printf("done")
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
            for i, file := range fs.Files {
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

