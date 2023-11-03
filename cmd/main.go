package main;

import (
    "flag"
    "fmt"
    "strconv"
    "time"
    membership "github.com/mjacob1002/425-MP3/pkg/membership"
)

var thisMachineName string

func onAdd(machineId string) {
    fmt.Println("Adding new node to membership list:", machineId)
}

func onDelete(machineId string) {
    fmt.Println("Deleting node from membership list:", machineId)
}

func main() {
    // Collect arguments
    var hostname, port, introducer string
    flag.StringVar(&thisMachineName, "machine_name", "", "Machine Name")
    flag.StringVar(&hostname, "hostname", "", "Hostname")
    flag.StringVar(&port, "port", "", "Port")
    flag.StringVar(&introducer, "introducer", "", "Introducer Node Address")
    flag.Parse()

    go membership.Join(
        (thisMachineName + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10)),
        hostname,
        port,
        introducer,
        onAdd,
        onDelete,
    )

    select {}
}
