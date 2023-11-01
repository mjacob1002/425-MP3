package main;

import (
    "flag"
    membership "github.com/mjacob1002/425-MP3/pkg/membership"
)

func main(){
    var machine_name, hostname, port, introducer string

    flag.StringVar(&machine_name, "machine_name", "", "Machine Name")
    flag.StringVar(&hostname, "hostname", "", "Hostname")
    flag.StringVar(&port, "port", "", "Port")
    flag.StringVar(&introducer, "introducer", "", "Introducer Node Address")

    flag.Parse()

    membership.Join(machine_name, hostname, port, introducer)
}
