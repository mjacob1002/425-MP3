package filesystem;

import (
	"golang.org/x/net/context"
	"path/filepath"
	"fmt"
	"io"
	"net"
	"os"
	"log"
	"google.golang.org/grpc"
)

var Tcp_port int64;
var TempDirectory string;
// stores the stubs used for gRPC methods
var MachineStubs map[string] FileSystemClient = make(map[string]FileSystemClient)

type Server struct {
	UnimplementedFileSystemServer;
}

// the get function
func(s *Server) Get(in *GetRequest, stream FileSystem_GetServer) error {
		fmt.Println("Got invoked Get with ", in.LocalName, " to be on sdfs as ", in.SdfsName);
		fname := filepath.Join(TempDirectory, in.SdfsName);
		file, err := os.Open(fname);
		if err != nil {
			log.Fatal(err);
		}
		for {
			buffer := make([]byte, 4096)
			bytesRead, err := file.Read(buffer);
			if err == io.EOF{
				break;
			}
			if err != nil {
				log.Fatal(err)
			}
			resp := GetResponse{Payload: string(buffer[:bytesRead])}
			stream.Send(&resp); // send stuff over the network
		}
		err = stream.Send(&GetResponse{}); // do I need this to end the connection for RPC?
		if err != nil {
				fmt.Println(err);
		}
		return err;
}

func (s *Server) Put(stream FileSystem_PutServer) error {
		// acquire a lock here with the information - also incorporate the filestruct idea here
		// go through all of their messages and write them
		fname := ""
		var file *os.File;
		for {
			req, err := stream.Recv();
			if err == io.EOF {
				break;
			}
			if err != nil {
				log.Fatal(err);
			}
			if fname == "" {
				fname = filepath.Join(TempDirectory, req.SdfsName)
				tempvar, err := os.Create(fname);
				if err != nil {
					log.Fatal(err)
				}
				file = tempvar
			}
			file.Write([]byte(req.PayloadToWrite));
		}
		file.Close()
		fmt.Println("just wrote to the file in %s\n", TempDirectory);
		return stream.SendAndClose(&PutResponse{Err: 0});
}

func Put(targetStub FileSystemClient, localfname string, sdfsname string){
	// do we have to do mutex shit here? Not sure. We don't need to do any struct stuff here
	file, err := os.Open(localfname)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := targetStub.Put(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for {
			buffer := make([]byte, 4096)
			bytesRead, err := file.Read(buffer);
			if err == io.EOF{
				break;
			}
			if err != nil {
				log.Fatal(err)
			}
			req := PutRequest{PayloadToWrite: string(buffer[:bytesRead]), LocalName: localfname, SdfsName: sdfsname}
			stream.Send(&req); // send stuff over the network
	}
	response, err := stream.CloseAndRecv() // close and receive
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response to PUT: %s\n", response.String());
}

func Delete(targetStub FileSystemClient, remotefname string){
	deletemsg := DeleteRequest{SdfsName: remotefname}
	resp, err := targetStub.Delete(context.Background(), &deletemsg);
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Delete response: %s\n", resp.String());
}

func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*DeleteResponse, error) {
		fmt.Println("Got invoked Delete by ", in.SdfsName);
		fname := filepath.Join(TempDirectory, in.SdfsName);
		err := os.Remove(fname);
		if err != nil {
			log.Fatal(err);
		}
		return &DeleteResponse{}, nil;
}

// should be called instead of directly calling the RPC
func Get(targetStub FileSystemClient, sdfsname string, localfname string){
		request := GetRequest{SdfsName: sdfsname, LocalName: localfname};
		// implement location logic here; going to just go to the first entry in our map for noA
		file, err := os.Create(localfname) // this is standard open - check if we need destructive write or something
		if err != nil {
			fmt.Println(err);
		}
		defer file.Close()
		stream, err := targetStub.Get(context.Background(), &request);
		if err != nil {
			log.Fatal(err);
		}
		for {
			response, err := stream.Recv()
			if err == io.EOF { // it's done sending data to client
				break;
			}
			if err != nil {
				log.Fatal(err);
			}
			// write to the file
			file.Write([]byte(response.Payload));
		}
}
func InitializeFileSystem(){
		TempDirectory, _ = os.MkdirTemp(".", "LOL");
		InitializeGRPCServer();
}

func InitializeGRPCServer(){
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Tcp_port))
		if err != nil {
		  log.Fatalf("failed to listen: %v", err)
		}
		fmt.Println("Listening with ", lis);
		grpcServer := grpc.NewServer()
		fmt.Println("finished grpcserver line")
		serv := Server{}
		fmt.Println("created serv");
		RegisterFileSystemServer(grpcServer, &serv)
		fmt.Println("Registered");
		grpcServer.Serve(lis)
		fmt.Println("Server up and running for gRPC...");
}
