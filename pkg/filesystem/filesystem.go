package filesystem;

import (
	"golang.org/x/net/context"
	"fmt"
)

var Tcp_port int64;
type FileSystemServerStub struct {

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

