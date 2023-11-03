package filesystem;

import (
	"golang.org/x/net/context"
	"fmt"
	"net"
	"log"
	"google.golang.org/grpc"
)

var Tcp_port int64;
type FileSystemServerStub struct {

}

type server struct {
	UnimplementedFileSystemServer;
}

func(s *FileSystemServerStub) Get(ctx context.Context, in *GetRequest) (*GetResponse, error) {
		fmt.Println("Got invoked Get with ", in.LocalName, " to be on sdfs as ", in.SdfsName);
		return &GetResponse{}, nil;
}

func (s *FileSystemServerStub) Put(ctx context.Context, in *PutRequest) (*PutResponse, error) {
		fmt.Println("Got invoked Put local_name=", in.LocalName, " sdfs=", in.SdfsName);
		return &PutResponse{}, nil;
}

func (s *FileSystemServerStub) Delete(ctx context.Context, in *DeleteRequest) (*DeleteResponse, error) {
		fmt.Println("Got invoked Delete by ", in.SdfsName);
		return &DeleteResponse{}, nil;
}


func InitializeGRPCServer(){
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", Tcp_port))
		if err != nil {
		  log.Fatalf("failed to listen: %v", err)
		}
		fmt.Println("Listening with ", lis);
		grpcServer := grpc.NewServer()
		fmt.Println("finished grpcserver line")
		serv := server{}
		fmt.Println("created serv");
		RegisterFileSystemServer(grpcServer, &serv)
		fmt.Println("Registered");
		grpcServer.Serve(lis)
		fmt.Println("Server up and running for gRPC...");
}
