package main;

import (
	"flag"
	"fmt"
	"log"
	"os"
	pb "github.com/mjacob1002/425-MP3/pkg/gen_proto"
	"google.golang.org/protobuf/proto"
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
	responses := int(0);
	if introducer != "" {
		membership.SendHeartbeatMessage(membership.MakeHeartbeat(), introducer)
	}
	for responses < 1024 {
		buffer := make([]byte , 4096);
		log.Println("waiting for message")
		n, addr , err := membership.MembershipConnection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err);
			return
		}
		log.Println("received a message! has ", n , " bytes from ", addr);
		var received_heartbeat = pb.HeartbeatMessage{}
		err = proto.Unmarshal(buffer[0:n], &received_heartbeat);
		if err != nil {
			fmt.Println(err);
			return;
		}
		responses++;
		// merge the tables
		membership.MergeTables(*(received_heartbeat.Table));
		new_dst := (received_heartbeat.Host) + ":" + (received_heartbeat.Port)
		if (responses % 50) == 0{

			actual_string := received_heartbeat.String()
			log.Printf("Received %s from a node\n", actual_string);
			log.Println("New table: ", membership.MembershipList)
		}
		membership.SendHeartbeatMessage(membership.MakeHeartbeat(), new_dst);
	}
	log.Println("finished ping pong");
}
