package membership;

import (
    "fmt"
    "time"
    "sync"
    "net"
    "os"
    "math/rand"

    "google.golang.org/protobuf/proto"
    pb "github.com/mjacob1002/425-MP3/pkg/gen_proto"
)

// Define node state variables
var MembershipList = make(map[string]pb.TableEntry)
var thisMachineId string
var thisHostname string
var thisPort string
var lock sync.Mutex
var deadMembers = make(map[string]int64)
var addressCache = make(map[string]*net.UDPAddr)

var thisAddCallback func(string, string)
var thisDeleteCallback func(string)

const (
    T_FAIL = 4000
    MAX_UDP_PACKET = 65535
    HEARTBEAT_FREQUENCY = 100
    CLEANUP_FREQUENCY = 500
    PEER_TARGET_COUNT = 4
)

func processHeartbeat(heartbeat *pb.HeartbeatMessage) {
    lock.Lock()
    defer lock.Unlock()

    for _, entry := range heartbeat.Table.Entries {
        // Skip entry referring to current node
        if entry.MachineId == thisMachineId {
            continue
        }

        // Skip entry if node has already been considered dead
        if _, ok := deadMembers[entry.MachineId]; ok {
            continue
        }

        if _, ok := MembershipList[entry.MachineId]; !ok {
            // Add new node to membership list
            newEntry := pb.TableEntry{}
            newEntry.MachineId = entry.MachineId
            newEntry.HeartbeatCounter = entry.HeartbeatCounter
            newEntry.Hostname = entry.Hostname
            newEntry.Port = entry.Port
            newEntry.LocalTime = time.Now().UnixMilli()
			newEntry.ApplicationPort= entry.ApplicationPort
            MembershipList[newEntry.MachineId] = newEntry
            serverAddress := newEntry.Hostname + ":" + newEntry.ApplicationPort
            go thisAddCallback(entry.MachineId, serverAddress)
        } else if entry.HeartbeatCounter > MembershipList[entry.MachineId].HeartbeatCounter {
            // Update pre-existing entry to higher heartbeat count
            updatedEntry := MembershipList[entry.MachineId]
            updatedEntry.HeartbeatCounter = entry.HeartbeatCounter
            updatedEntry.LocalTime = time.Now().UnixMilli()
            MembershipList[updatedEntry.MachineId] = updatedEntry
        }
    }
}

func listenInitializer() {
    // Initialize socket
    udpAddress, err := net.ResolveUDPAddr("udp", ":" + thisPort)
    if err != nil {
        fmt.Errorf("net.ResolveUDPAddr: %v\n", err)
        os.Exit(1)
    }

    conn, err := net.ListenUDP("udp", udpAddress)
    if err != nil {
        fmt.Errorf("net.ListenUDP: %v\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    buffer := make([]byte, MAX_UDP_PACKET)
    for {
        // Read packets
        n, _, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Errorf("conn.ReadFromUDP: %v\n", err)
        }

        // Parse protobuf messages
        message := &pb.HeartbeatMessage{}
        if err := proto.Unmarshal(buffer[0:n], message); err != nil {
            fmt.Errorf("proto.Unmarshal: %v\n", err)
            continue
        }

        // Pass message to helper function
        processHeartbeat(message)
    }
}

func incrementHeartbeat(){
    lock.Lock()
    defer lock.Unlock()

    // Update heartbeat counter for this node in membership list
    entry, ok := MembershipList[thisMachineId];
    if !ok {
        fmt.Errorf("Node does not exist in its own membership list\n")
        os.Exit(1)
    }
    entry.HeartbeatCounter = entry.HeartbeatCounter + 1
    MembershipList[thisMachineId] = entry
}

func sendPeriodicHeartbeats() {
    for {
        time.Sleep(HEARTBEAT_FREQUENCY * time.Millisecond)
        incrementHeartbeat()
        sendOutHeartbeats()
    }
}

func periodicCleanupTable() {
    for {
        time.Sleep(CLEANUP_FREQUENCY * time.Millisecond)
        cleanupTable()
    }
}

func sendOutHeartbeats() {
    lock.Lock()
    defer lock.Unlock()

    // Get a list of all the machine ids in the membership list
    machineIds := make([]string, 0, len(MembershipList) - 1)
    for machineId := range MembershipList {
        if machineId != thisMachineId {
            machineIds = append(machineIds, machineId)
        }
    }

    // Shuffle the array
    for i := len(machineIds) - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        machineIds[i], machineIds[j] = machineIds[j], machineIds[i]
    }

    // Select the machines to gossip to
    k := PEER_TARGET_COUNT
    if len(MembershipList) - 1 < k {
        k = len(MembershipList) - 1
    }
    selectedMachineIds := machineIds[:k]

    // Gossip and profit
    for _, machineId := range selectedMachineIds {
        sendHeartbeat(MembershipList[machineId].Hostname, MembershipList[machineId].Port)
    }
}

func resolveUDPAddress(address string) (*net.UDPAddr, error) {
    // Check cache map
    if cachedAddress, ok := addressCache[address]; ok {
        return cachedAddress, nil
    }

    // Manually call the underlying resolve function
    udpAddr, err := net.ResolveUDPAddr("udp", address)
    if err != nil {
        return nil, err
    }

    // Store the result in the cache
    addressCache[address] = udpAddr

    return udpAddr, nil
}

func makeHeartbeatFromMembershipList() (pb.HeartbeatMessage) {
    // Generate array of pointers to each membership list entry
    membershipArray := make([]*pb.TableEntry, 0, len(MembershipList))
    for _, value := range MembershipList {
        copiedValue := value
        membershipArray = append(membershipArray, &copiedValue)
    }

    // Generate membership list table
    membershipTable := pb.Table{ Entries : membershipArray }

    // Generate heartbeat message
    heartbeatMessage := pb.HeartbeatMessage{
        Table: &membershipTable,
    }

    return heartbeatMessage
}

func sendHeartbeat(hostname string, port string){
    sendHeartbeatAddress(hostname + ":" + port)
}

func sendHeartbeatAddress(address string){
    // Setup connection
    udpAddr, err := resolveUDPAddress(address)
    if err != nil {
        fmt.Errorf("resolveUDPAddress: \n", err)
        os.Exit(1)
    }

    conn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        fmt.Errorf("net.DialUDP: %v\n", err)
        return
    }
    defer conn.Close()

    // Create the hearbeat
    heartbeat := makeHeartbeatFromMembershipList()
    serializedHeartbeat, err := proto.Marshal(&heartbeat)
    if err != nil {
        fmt.Errorf("proto.Marshal: %v\n", err)
        return
    }

    // Write serialized heartbeat to connection
    _, err = conn.Write(serializedHeartbeat)
    if err != nil {
        fmt.Errorf("conn.Write: %v\n", err)
    }
}

func cleanupTable() {
    lock.Lock()
    defer lock.Unlock()

    for key, value := range MembershipList {
        // Skip entry referring to current node
        if key == thisMachineId {
            continue
        }

        if time.Now().UnixMilli() - value.LocalTime >= T_FAIL {
            // Delete the node from membership list and add to dead members
            deadMembers[key] = value.LocalTime
            delete(MembershipList, key)
            go thisDeleteCallback(key)
        }
    }
}

func Join(machineId string, hostname string, port string, introducer string, applicationPort string, addCallback func(string, string), deleteCallback func(string)) {
    // Initialize node state variables
    thisHostname = hostname
    thisPort = port
    thisMachineId = machineId
    thisAddCallback = addCallback
    thisDeleteCallback = deleteCallback

    fmt.Println("Joining Node Info:", thisMachineId, thisHostname, thisPort, introducer)

    // Add node to membership list
    MembershipList[thisMachineId] = pb.TableEntry {
        MachineId: thisMachineId,
        HeartbeatCounter: 0,
        Hostname: thisHostname,
        Port: thisPort,
        LocalTime: time.Now().UnixMilli(),
		ApplicationPort: applicationPort,
    }

    // Introduce node to known node
    if introducer != "" {
        sendHeartbeatAddress(introducer)
    }

    // Start all go routines
    go listenInitializer()
    go sendPeriodicHeartbeats()
    go periodicCleanupTable()

    // Wait
    select {}
}

