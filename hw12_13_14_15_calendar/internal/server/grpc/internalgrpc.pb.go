// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: internal/server/grpc/internalgrpc.proto

package internalgrpc

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Email struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=Email,proto3" json:"Email,omitempty"`
}

func (x *Email) Reset() {
	*x = Email{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Email) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Email) ProtoMessage() {}

func (x *Email) ProtoReflect() protoreflect.Message {
	mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Email.ProtoReflect.Descriptor instead.
func (*Email) Descriptor() ([]byte, []int) {
	return file_internal_server_grpc_internalgrpc_proto_rawDescGZIP(), []int{0}
}

func (x *Email) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        int64  `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Email     *Email `protobuf:"bytes,2,opt,name=Email,proto3" json:"Email,omitempty"`
	FirstName string `protobuf:"bytes,3,opt,name=FirstName,proto3" json:"FirstName,omitempty"`
	LastName  string `protobuf:"bytes,4,opt,name=LastName,proto3" json:"LastName,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_internal_server_grpc_internalgrpc_proto_rawDescGZIP(), []int{1}
}

func (x *User) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *User) GetEmail() *Email {
	if x != nil {
		return x.Email
	}
	return nil
}

func (x *User) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *User) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID       int64                `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	UserID   int64                `protobuf:"varint,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Title    string               `protobuf:"bytes,3,opt,name=Title,proto3" json:"Title,omitempty"`
	Content  string               `protobuf:"bytes,4,opt,name=Content,proto3" json:"Content,omitempty"`
	DateFrom *timestamp.Timestamp `protobuf:"bytes,5,opt,name=DateFrom,proto3" json:"DateFrom,omitempty"`
	DateTo   *timestamp.Timestamp `protobuf:"bytes,6,opt,name=DateTo,proto3" json:"DateTo,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_internal_server_grpc_internalgrpc_proto_rawDescGZIP(), []int{2}
}

func (x *Event) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Event) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *Event) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Event) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Event) GetDateFrom() *timestamp.Timestamp {
	if x != nil {
		return x.DateFrom
	}
	return nil
}

func (x *Event) GetDateTo() *timestamp.Timestamp {
	if x != nil {
		return x.DateTo
	}
	return nil
}

type Events struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*Event `protobuf:"bytes,1,rep,name=Events,proto3" json:"Events,omitempty"`
}

func (x *Events) Reset() {
	*x = Events{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Events) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Events) ProtoMessage() {}

func (x *Events) ProtoReflect() protoreflect.Message {
	mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Events.ProtoReflect.Descriptor instead.
func (*Events) Descriptor() ([]byte, []int) {
	return file_internal_server_grpc_internalgrpc_proto_rawDescGZIP(), []int{3}
}

func (x *Events) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type DateEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User                `protobuf:"bytes,1,opt,name=User,proto3" json:"User,omitempty"`
	Date *timestamp.Timestamp `protobuf:"bytes,2,opt,name=Date,proto3" json:"Date,omitempty"`
}

func (x *DateEvent) Reset() {
	*x = DateEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DateEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DateEvent) ProtoMessage() {}

func (x *DateEvent) ProtoReflect() protoreflect.Message {
	mi := &file_internal_server_grpc_internalgrpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DateEvent.ProtoReflect.Descriptor instead.
func (*DateEvent) Descriptor() ([]byte, []int) {
	return file_internal_server_grpc_internalgrpc_proto_rawDescGZIP(), []int{4}
}

func (x *DateEvent) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *DateEvent) GetDate() *timestamp.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

var File_internal_server_grpc_internalgrpc_proto protoreflect.FileDescriptor

var file_internal_server_grpc_internalgrpc_proto_rawDesc = []byte{
	0x0a, 0x27, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1d, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14,
	0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x22, 0x7b, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x29, 0x0a, 0x05,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x52, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x46, 0x69, 0x72, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x46, 0x69, 0x72, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0xcb, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x08, 0x44, 0x61, 0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x32, 0x0a, 0x06, 0x44,
	0x61, 0x74, 0x65, 0x54, 0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x44, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x22,
	0x35, 0x0a, 0x06, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x63, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x65, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x04, 0x44,
	0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x44, 0x61, 0x74, 0x65, 0x32, 0x9f, 0x06, 0x0a, 0x08,
	0x43, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x12, 0x34, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x12, 0x36,
	0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x1a, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x12, 0x36,
	0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x1a, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x12, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x14, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x00, 0x12,
	0x39, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x13,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x1a, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x0b, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x13,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00,
	0x12, 0x3d, 0x0a, 0x0b, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12,
	0x17, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44,
	0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12,
	0x3e, 0x0a, 0x0c, 0x57, 0x65, 0x65, 0x6b, 0x6c, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12,
	0x17, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44,
	0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12,
	0x3f, 0x0a, 0x0d, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x6c, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x12, 0x17, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x44, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00,
	0x12, 0x45, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x61,
	0x64, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x13, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x13, 0x4d, 0x61, 0x72, 0x6b, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x41, 0x73, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x13,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x10, 0x5a,
	0x0e, 0x2e, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x67, 0x72, 0x70, 0x63, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_server_grpc_internalgrpc_proto_rawDescOnce sync.Once
	file_internal_server_grpc_internalgrpc_proto_rawDescData = file_internal_server_grpc_internalgrpc_proto_rawDesc
)

func file_internal_server_grpc_internalgrpc_proto_rawDescGZIP() []byte {
	file_internal_server_grpc_internalgrpc_proto_rawDescOnce.Do(func() {
		file_internal_server_grpc_internalgrpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_server_grpc_internalgrpc_proto_rawDescData)
	})
	return file_internal_server_grpc_internalgrpc_proto_rawDescData
}

var file_internal_server_grpc_internalgrpc_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_internal_server_grpc_internalgrpc_proto_goTypes = []interface{}{
	(*Email)(nil),               // 0: internalgrpc.Email
	(*User)(nil),                // 1: internalgrpc.User
	(*Event)(nil),               // 2: internalgrpc.Event
	(*Events)(nil),              // 3: internalgrpc.Events
	(*DateEvent)(nil),           // 4: internalgrpc.DateEvent
	(*timestamp.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*empty.Empty)(nil),         // 6: google.protobuf.Empty
}
var file_internal_server_grpc_internalgrpc_proto_depIdxs = []int32{
	0,  // 0: internalgrpc.User.Email:type_name -> internalgrpc.Email
	5,  // 1: internalgrpc.Event.DateFrom:type_name -> google.protobuf.Timestamp
	5,  // 2: internalgrpc.Event.DateTo:type_name -> google.protobuf.Timestamp
	2,  // 3: internalgrpc.Events.Events:type_name -> internalgrpc.Event
	1,  // 4: internalgrpc.DateEvent.User:type_name -> internalgrpc.User
	5,  // 5: internalgrpc.DateEvent.Date:type_name -> google.protobuf.Timestamp
	0,  // 6: internalgrpc.Calendar.GetUser:input_type -> internalgrpc.Email
	1,  // 7: internalgrpc.Calendar.CreateUser:input_type -> internalgrpc.User
	1,  // 8: internalgrpc.Calendar.UpdateUser:input_type -> internalgrpc.User
	1,  // 9: internalgrpc.Calendar.DeleteUser:input_type -> internalgrpc.User
	1,  // 10: internalgrpc.Calendar.GetEvents:input_type -> internalgrpc.User
	2,  // 11: internalgrpc.Calendar.CreateEvent:input_type -> internalgrpc.Event
	2,  // 12: internalgrpc.Calendar.UpdateEvent:input_type -> internalgrpc.Event
	2,  // 13: internalgrpc.Calendar.DeleteEvent:input_type -> internalgrpc.Event
	4,  // 14: internalgrpc.Calendar.DailyEvents:input_type -> internalgrpc.DateEvent
	4,  // 15: internalgrpc.Calendar.WeeklyEvents:input_type -> internalgrpc.DateEvent
	4,  // 16: internalgrpc.Calendar.MonthlyEvents:input_type -> internalgrpc.DateEvent
	6,  // 17: internalgrpc.Calendar.GetNotifyReadyEvents:input_type -> google.protobuf.Empty
	2,  // 18: internalgrpc.Calendar.MarkEventAsNotified:input_type -> internalgrpc.Event
	1,  // 19: internalgrpc.Calendar.GetUser:output_type -> internalgrpc.User
	1,  // 20: internalgrpc.Calendar.CreateUser:output_type -> internalgrpc.User
	1,  // 21: internalgrpc.Calendar.UpdateUser:output_type -> internalgrpc.User
	1,  // 22: internalgrpc.Calendar.DeleteUser:output_type -> internalgrpc.User
	3,  // 23: internalgrpc.Calendar.GetEvents:output_type -> internalgrpc.Events
	2,  // 24: internalgrpc.Calendar.CreateEvent:output_type -> internalgrpc.Event
	2,  // 25: internalgrpc.Calendar.UpdateEvent:output_type -> internalgrpc.Event
	2,  // 26: internalgrpc.Calendar.DeleteEvent:output_type -> internalgrpc.Event
	2,  // 27: internalgrpc.Calendar.DailyEvents:output_type -> internalgrpc.Event
	2,  // 28: internalgrpc.Calendar.WeeklyEvents:output_type -> internalgrpc.Event
	2,  // 29: internalgrpc.Calendar.MonthlyEvents:output_type -> internalgrpc.Event
	2,  // 30: internalgrpc.Calendar.GetNotifyReadyEvents:output_type -> internalgrpc.Event
	6,  // 31: internalgrpc.Calendar.MarkEventAsNotified:output_type -> google.protobuf.Empty
	19, // [19:32] is the sub-list for method output_type
	6,  // [6:19] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_internal_server_grpc_internalgrpc_proto_init() }
func file_internal_server_grpc_internalgrpc_proto_init() {
	if File_internal_server_grpc_internalgrpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_server_grpc_internalgrpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Email); i {
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
		file_internal_server_grpc_internalgrpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_internal_server_grpc_internalgrpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_internal_server_grpc_internalgrpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Events); i {
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
		file_internal_server_grpc_internalgrpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DateEvent); i {
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
			RawDescriptor: file_internal_server_grpc_internalgrpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_server_grpc_internalgrpc_proto_goTypes,
		DependencyIndexes: file_internal_server_grpc_internalgrpc_proto_depIdxs,
		MessageInfos:      file_internal_server_grpc_internalgrpc_proto_msgTypes,
	}.Build()
	File_internal_server_grpc_internalgrpc_proto = out.File
	file_internal_server_grpc_internalgrpc_proto_rawDesc = nil
	file_internal_server_grpc_internalgrpc_proto_goTypes = nil
	file_internal_server_grpc_internalgrpc_proto_depIdxs = nil
}