// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: proto/filesystem.proto

package filesystem

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InvokeReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SdfsName string `protobuf:"bytes,1,opt,name=SdfsName,proto3" json:"SdfsName,omitempty"`
}

func (x *InvokeReadRequest) Reset() {
	*x = InvokeReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvokeReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvokeReadRequest) ProtoMessage() {}

func (x *InvokeReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvokeReadRequest.ProtoReflect.Descriptor instead.
func (*InvokeReadRequest) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{0}
}

func (x *InvokeReadRequest) GetSdfsName() string {
	if x != nil {
		return x.SdfsName
	}
	return ""
}

type InvokeReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LocalName string `protobuf:"bytes,1,opt,name=local_name,json=localName,proto3" json:"local_name,omitempty"`
}

func (x *InvokeReadResponse) Reset() {
	*x = InvokeReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvokeReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvokeReadResponse) ProtoMessage() {}

func (x *InvokeReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvokeReadResponse.ProtoReflect.Descriptor instead.
func (*InvokeReadResponse) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{1}
}

func (x *InvokeReadResponse) GetLocalName() string {
	if x != nil {
		return x.LocalName
	}
	return ""
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SdfsName string `protobuf:"bytes,2,opt,name=sdfs_name,json=sdfsName,proto3" json:"sdfs_name,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetSdfsName() string {
	if x != nil {
		return x.SdfsName
	}
	return ""
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload string `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Err     int64  `protobuf:"varint,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{3}
}

func (x *GetResponse) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *GetResponse) GetErr() int64 {
	if x != nil {
		return x.Err
	}
	return 0
}

type PutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PayloadToWrite string `protobuf:"bytes,1,opt,name=payload_to_write,json=payloadToWrite,proto3" json:"payload_to_write,omitempty"`
	SdfsName       string `protobuf:"bytes,2,opt,name=sdfs_name,json=sdfsName,proto3" json:"sdfs_name,omitempty"`
	Replica        bool   `protobuf:"varint,3,opt,name=replica,proto3" json:"replica,omitempty"`
}

func (x *PutRequest) Reset() {
	*x = PutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutRequest) ProtoMessage() {}

func (x *PutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutRequest.ProtoReflect.Descriptor instead.
func (*PutRequest) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{4}
}

func (x *PutRequest) GetPayloadToWrite() string {
	if x != nil {
		return x.PayloadToWrite
	}
	return ""
}

func (x *PutRequest) GetSdfsName() string {
	if x != nil {
		return x.SdfsName
	}
	return ""
}

func (x *PutRequest) GetReplica() bool {
	if x != nil {
		return x.Replica
	}
	return false
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SdfsName string `protobuf:"bytes,1,opt,name=sdfs_name,json=sdfsName,proto3" json:"sdfs_name,omitempty"`
	Replica  bool   `protobuf:"varint,2,opt,name=replica,proto3" json:"replica,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteRequest) GetSdfsName() string {
	if x != nil {
		return x.SdfsName
	}
	return ""
}

func (x *DeleteRequest) GetReplica() bool {
	if x != nil {
		return x.Replica
	}
	return false
}

type FileRangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start uint32 `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	End   uint32 `protobuf:"varint,2,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *FileRangeRequest) Reset() {
	*x = FileRangeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileRangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileRangeRequest) ProtoMessage() {}

func (x *FileRangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileRangeRequest.ProtoReflect.Descriptor instead.
func (*FileRangeRequest) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{6}
}

func (x *FileRangeRequest) GetStart() uint32 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *FileRangeRequest) GetEnd() uint32 {
	if x != nil {
		return x.End
	}
	return 0
}

type FileRangeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SdfsNames []string `protobuf:"bytes,1,rep,name=sdfs_names,json=sdfsNames,proto3" json:"sdfs_names,omitempty"`
}

func (x *FileRangeResponse) Reset() {
	*x = FileRangeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_filesystem_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileRangeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileRangeResponse) ProtoMessage() {}

func (x *FileRangeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_filesystem_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileRangeResponse.ProtoReflect.Descriptor instead.
func (*FileRangeResponse) Descriptor() ([]byte, []int) {
	return file_proto_filesystem_proto_rawDescGZIP(), []int{7}
}

func (x *FileRangeResponse) GetSdfsNames() []string {
	if x != nil {
		return x.SdfsNames
	}
	return nil
}

var File_proto_filesystem_proto protoreflect.FileDescriptor

var file_proto_filesystem_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x11, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x64,
	0x66, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x64,
	0x66, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x33, 0x0a, 0x12, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65,
	0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x29, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x64, 0x66,
	0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x64,
	0x66, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x39, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x65, 0x72,
	0x72, 0x22, 0x6d, 0x0a, 0x0a, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x28, 0x0a, 0x10, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x74, 0x6f, 0x5f, 0x77, 0x72,
	0x69, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x54, 0x6f, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x64, 0x66,
	0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x64,
	0x66, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x22, 0x46, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x64, 0x66, 0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x64, 0x66, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x22, 0x3a, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x03, 0x65, 0x6e, 0x64, 0x22, 0x32, 0x0a, 0x11, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x64, 0x66,
	0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x64, 0x66, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x32, 0x85, 0x02, 0x0a, 0x0a, 0x46, 0x69, 0x6c,
	0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x24, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x0b,
	0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x2e, 0x0a,
	0x03, 0x50, 0x75, 0x74, 0x12, 0x0b, 0x2e, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x28, 0x01, 0x12, 0x32, 0x0a,
	0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0e, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x34, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x11,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x12, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0a, 0x49, 0x6e, 0x76, 0x6f, 0x6b,
	0x65, 0x52, 0x65, 0x61, 0x64, 0x12, 0x12, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x49, 0x6e, 0x76, 0x6f,
	0x6b, 0x65, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x13, 0x5a, 0x11, 0x2e, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_filesystem_proto_rawDescOnce sync.Once
	file_proto_filesystem_proto_rawDescData = file_proto_filesystem_proto_rawDesc
)

func file_proto_filesystem_proto_rawDescGZIP() []byte {
	file_proto_filesystem_proto_rawDescOnce.Do(func() {
		file_proto_filesystem_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_filesystem_proto_rawDescData)
	})
	return file_proto_filesystem_proto_rawDescData
}

var file_proto_filesystem_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_filesystem_proto_goTypes = []interface{}{
	(*InvokeReadRequest)(nil),  // 0: InvokeReadRequest
	(*InvokeReadResponse)(nil), // 1: InvokeReadResponse
	(*GetRequest)(nil),         // 2: GetRequest
	(*GetResponse)(nil),        // 3: GetResponse
	(*PutRequest)(nil),         // 4: PutRequest
	(*DeleteRequest)(nil),      // 5: DeleteRequest
	(*FileRangeRequest)(nil),   // 6: FileRangeRequest
	(*FileRangeResponse)(nil),  // 7: FileRangeResponse
	(*emptypb.Empty)(nil),      // 8: google.protobuf.Empty
}
var file_proto_filesystem_proto_depIdxs = []int32{
	2, // 0: FileSystem.Get:input_type -> GetRequest
	4, // 1: FileSystem.Put:input_type -> PutRequest
	5, // 2: FileSystem.Delete:input_type -> DeleteRequest
	6, // 3: FileSystem.FileRange:input_type -> FileRangeRequest
	0, // 4: FileSystem.InvokeRead:input_type -> InvokeReadRequest
	3, // 5: FileSystem.Get:output_type -> GetResponse
	8, // 6: FileSystem.Put:output_type -> google.protobuf.Empty
	8, // 7: FileSystem.Delete:output_type -> google.protobuf.Empty
	7, // 8: FileSystem.FileRange:output_type -> FileRangeResponse
	1, // 9: FileSystem.InvokeRead:output_type -> InvokeReadResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_filesystem_proto_init() }
func file_proto_filesystem_proto_init() {
	if File_proto_filesystem_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_filesystem_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvokeReadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_filesystem_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvokeReadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_filesystem_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_filesystem_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_filesystem_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_filesystem_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_filesystem_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileRangeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_filesystem_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileRangeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_filesystem_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_filesystem_proto_goTypes,
		DependencyIndexes: file_proto_filesystem_proto_depIdxs,
		MessageInfos:      file_proto_filesystem_proto_msgTypes,
	}.Build()
	File_proto_filesystem_proto = out.File
	file_proto_filesystem_proto_rawDesc = nil
	file_proto_filesystem_proto_goTypes = nil
	file_proto_filesystem_proto_depIdxs = nil
}
