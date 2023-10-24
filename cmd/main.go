package main;

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	membership "github.com/mjacob1002/425-MP3/pkg/membership"
	"strconv"
)

func InitializeLogger(){
	logic_node_str := strconv.Itoa(int(membership.LogicalNode))
	fname := logic_node_str + ".log"
	log.Println(fname)
	file, err := os.OpenFile(fname, os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err);
		return;
	}
	log.SetOutput(file);
	log.Println("Logger initialized");
}


func main(){
	var introducer string;
	flag.StringVar(&membership.MEMBERSHIP_PORT, "listen_port", "8000", "the port to listen to incoming UDP messages")
	flag.StringVar(&introducer, "introducer", "", "the node to ping when you want to join the system")
	flag.Int64Var(&membership.LogicalNode, "logical_node", 0, "the number of the node");
	flag.Parse();
	host, err := os.Hostname();
	if err != nil {
		fmt.Println(err);
	}
	log.Println("My host address:", host);
	InitializeLogger();
	membership.InitializeMembership(); // setup the listening port	
	// establishing the writing connection
	go membership.ListenForHeartbeats();
	if introducer != "" {
		membership.SendHeartbeatMessage(membership.MakeHeartbeat(), introducer)
	}
	go membership.PingHeartbeats();
	go membership.PruneTable();
	for {
		 reader := bufio.NewReader(os.Stdin);
		 fmt.Printf("Type a command here\n");
		 text, _ := reader.ReadString('\n')
		 if text == "get members\n"{
			membership.MembershipMutex.Lock();
			for _, entry := range(membership.MembershipList){
					fmt.Printf("%s\n", entry.String());
			}
			membership.MembershipMutex.Unlock()
		 } else {
			fmt.Println("Didn't recognize the following command: ", text);
		 }

	}
}
