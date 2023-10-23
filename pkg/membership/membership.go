package membership;

import (
	"fmt"
	"log"
	"net"
	"os"
	pb "github.com/mjacob1002/425-MP3/pkg/gen_proto"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

var MEMBERSHIP_ID string;
var MEMBERSHIP_PORT string;
var LogicalNode int64;
var MembershipConnection *net.UDPConn;
var MembershipList = make(map[string]pb.TableEntry);


func SetupMembershipPort(){
	s, err := net.ResolveUDPAddr("udp4", ":" + MEMBERSHIP_PORT);
	if err != nil {
		fmt.Println(err);
		return;
	}
	fmt.Println("My own socket", s);
	MembershipConnection, err = net.ListenUDP("udp4", s);
	fmt.Println("My own connection: ", MembershipConnection);
	if err != nil {
		fmt.Println(err);
		return;
	}
}

func InitializeMembership(){
	MEMBERSHIP_ID = strconv.Itoa(int(LogicalNode)) + strconv.Itoa(int(time.Now().Unix()))
	log.Printf("New Membership ID created: %s\n", MEMBERSHIP_ID);
	SetupMembershipPort();
	hostname, err := os.Hostname();
	if err != nil {
		fmt.Println(err);
		os.Exit(1)
	}
	MembershipList[MEMBERSHIP_ID] = pb.TableEntry{MachineId: MEMBERSHIP_ID, HeartbeatCounter: 0, Hostname: hostname, Port: MEMBERSHIP_PORT, LocalTime: time.Now().Unix(), CurrState: pb.NodeState_ALIVE}
	log.Println("Added myself to my membership list: ", MembershipList[MEMBERSHIP_ID]);
}

func MakeHeartbeat()(pb.HeartbeatMessage){
		// increment my own heartbeat
		entry, ok := MembershipList[MEMBERSHIP_ID]
		if ok {
			log.Printf("The heartbeat of %s: %s\n", MEMBERSHIP_ID, entry.String());
			entry.HeartbeatCounter = entry.HeartbeatCounter + 1;
			MembershipList[MEMBERSHIP_ID] = entry;
			log.Printf("The new heartbeat for myself: %s\n", entry.String());
		} else {
			log.Println("Something went wrong when trying to acess MEMBERSHIP_ID key in the map");
		}
		alive_list := []*pb.TableEntry{}
		for _, entry := range(MembershipList){
			if(entry.CurrState == pb.NodeState_ALIVE){
				copied_entry := entry
				alive_list = append(alive_list, &copied_entry);
			}
		}
		table := pb.Table{}
		table.Entries = alive_list;
		host, err := os.Hostname()
		if err != nil {
			log.Println(err);
			panic(err);
		}
		heartbeat_msg := pb.HeartbeatMessage{Table: &table, Host: host, Port: MEMBERSHIP_PORT}
		return heartbeat_msg;

}

func SendHeartbeatMessage(heartbeat pb.HeartbeatMessage, address string) {
		log.Println("inside the heartbeat function")
		log.Printf("Going to translate following address: %s\n", address)
		dst, err := net.ResolveUDPAddr("udp4", address);
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("resolved udp address, trying to write to it now: ", dst)
		payload_str, err := proto.Marshal(&heartbeat)// FINISH THIS LINE 
		if err != nil {
			log.Println(err);
			return;
		}
		payload := []byte(payload_str)
		num_bytes_written , err := MembershipConnection.WriteToUDP(payload, dst)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("wrote %d bytes\n", num_bytes_written);
}


func MergeTables(table pb.Table){
	for _, entry := range(table.Entries){
			log.Println("Currently looking at this entry that I received: %s\n", entry.String());
			if entry.CurrState == pb.NodeState_DEAD {
				// they marked the node as dead, I don't care
				log.Println("The other marked %s as dead, I don't care...", entry.MachineId);
				continue;
			}
			val, ok := MembershipList[entry.MachineId];
			if !ok { // 
				log.Println("I don't have the entry of %s. Will add %s now \n", entry.MachineId, entry.String());
				entry.LocalTime = time.Now().Unix(); // local time
				MembershipList[entry.MachineId] = *entry;
				list_obj, ok := MembershipList[entry.MachineId];
				if(!ok){
					log.Println("Didn't insert properly\n");
					os.Exit(1);
				}
				list_ptr := &list_obj
				log.Printf("Entry receieved: %s, Entry inserted: %s\n", entry.String(), list_ptr.String());
			} else {
				if (val.CurrState == pb.NodeState_DEAD){
					// I marked this node as dead, i don't care
					log.Println("I has marked %s as dead, I don't care...\n", val.MachineId);
					continue;
				}
				if(val.HeartbeatCounter > entry.HeartbeatCounter){
					log.Printf("Not updating entry for %s; I have heartbeat counter of %d while the other is %d\n", entry.MachineId, val.HeartbeatCounter, entry.HeartbeatCounter);
					continue; // no needto update, I have a more up to date entry
				} else {
					entry.LocalTime = time.Now().Unix();
					log.Println(MembershipList[entry.MachineId]);
					MembershipList[entry.MachineId] = *entry;
					log.Println(MembershipList[entry.MachineId]);
					log.Printf("Updated entry for %s to %s\n", entry.MachineId, entry.String());
				}
			}
	}
}
