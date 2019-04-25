// Code generated by protoc-gen-go. DO NOT EDIT.
// source: syncer.proto

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import proto2 "github.com/apache/servicecomb-service-center/server/core/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PullRequest struct {
	ServiceName string `protobuf:"bytes,1,opt,name=serviceName" json:"serviceName,omitempty"`
	Options     string `protobuf:"bytes,2,opt,name=options" json:"options,omitempty"`
	Time        string `protobuf:"bytes,3,opt,name=time" json:"time,omitempty"`
}

func (m *PullRequest) Reset()                    { *m = PullRequest{} }
func (m *PullRequest) String() string            { return proto1.CompactTextString(m) }
func (*PullRequest) ProtoMessage()               {}
func (*PullRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *PullRequest) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *PullRequest) GetOptions() string {
	if m != nil {
		return m.Options
	}
	return ""
}

func (m *PullRequest) GetTime() string {
	if m != nil {
		return m.Time
	}
	return ""
}

type SyncService struct {
	DomainProject string                         `protobuf:"bytes,1,opt,name=domainProject" json:"domainProject,omitempty"`
	Service       *proto2.MicroService           `protobuf:"bytes,2,opt,name=service" json:"service,omitempty"`
	Instances     []*proto2.MicroServiceInstance `protobuf:"bytes,3,rep,name=instances" json:"instances,omitempty"`
}

func (m *SyncService) Reset()                    { *m = SyncService{} }
func (m *SyncService) String() string            { return proto1.CompactTextString(m) }
func (*SyncService) ProtoMessage()               {}
func (*SyncService) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *SyncService) GetDomainProject() string {
	if m != nil {
		return m.DomainProject
	}
	return ""
}

func (m *SyncService) GetService() *proto2.MicroService {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *SyncService) GetInstances() []*proto2.MicroServiceInstance {
	if m != nil {
		return m.Instances
	}
	return nil
}

type SyncData struct {
	Services []*SyncService `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
}

func (m *SyncData) Reset()                    { *m = SyncData{} }
func (m *SyncData) String() string            { return proto1.CompactTextString(m) }
func (*SyncData) ProtoMessage()               {}
func (*SyncData) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *SyncData) GetServices() []*SyncService {
	if m != nil {
		return m.Services
	}
	return nil
}

func init() {
	proto1.RegisterType((*PullRequest)(nil), "proto.PullRequest")
	proto1.RegisterType((*SyncService)(nil), "proto.SyncService")
	proto1.RegisterType((*SyncData)(nil), "proto.SyncData")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Sync service

type SyncClient interface {
	Pull(ctx context.Context, in *PullRequest, opts ...grpc.CallOption) (*SyncData, error)
}

type syncClient struct {
	cc *grpc.ClientConn
}

func NewSyncClient(cc *grpc.ClientConn) SyncClient {
	return &syncClient{cc}
}

func (c *syncClient) Pull(ctx context.Context, in *PullRequest, opts ...grpc.CallOption) (*SyncData, error) {
	out := new(SyncData)
	err := grpc.Invoke(ctx, "/proto.Sync/Pull", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Sync service

type SyncServer interface {
	Pull(context.Context, *PullRequest) (*SyncData, error)
}

func RegisterSyncServer(s *grpc.Server, srv SyncServer) {
	s.RegisterService(&_Sync_serviceDesc, srv)
}

func _Sync_Pull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncServer).Pull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Sync/Pull",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncServer).Pull(ctx, req.(*PullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Sync_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Sync",
	HandlerType: (*SyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Pull",
			Handler:    _Sync_Pull_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syncer.proto",
}

func init() { proto1.RegisterFile("syncer.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 259 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x09, 0x09, 0xb4, 0x3d, 0xf3, 0x47, 0x3a, 0x16, 0xab, 0x2c, 0x51, 0xc4, 0xd0, 0xa5,
	0x19, 0x82, 0x18, 0x60, 0x66, 0x61, 0x00, 0x55, 0xee, 0xcc, 0x60, 0xcc, 0x0d, 0x46, 0x8d, 0x5d,
	0x6c, 0x17, 0xa9, 0x1f, 0x86, 0xef, 0x8a, 0xe2, 0x38, 0x10, 0xa4, 0x4e, 0x89, 0xdf, 0xfd, 0xfc,
	0xde, 0xf9, 0xc1, 0x99, 0xdf, 0x1b, 0x45, 0xae, 0xde, 0x3a, 0x1b, 0x2c, 0x9e, 0xc4, 0xcf, 0xfc,
	0xc2, 0x93, 0xfb, 0xd2, 0x8a, 0x7c, 0x2f, 0x57, 0xaf, 0xc0, 0x56, 0xbb, 0xcd, 0x46, 0xd0, 0xe7,
	0x8e, 0x7c, 0xc0, 0x12, 0x58, 0x02, 0x5e, 0x64, 0x4b, 0x3c, 0x2b, 0xb3, 0xc5, 0x4c, 0x8c, 0x25,
	0xe4, 0x30, 0xb1, 0xdb, 0xa0, 0xad, 0xf1, 0xfc, 0x38, 0x4e, 0x87, 0x23, 0x22, 0x14, 0x41, 0xb7,
	0xc4, 0xf3, 0x28, 0xc7, 0xff, 0xea, 0x3b, 0x03, 0xb6, 0xde, 0x1b, 0xb5, 0xee, 0x1d, 0xf0, 0x06,
	0xce, 0xdf, 0x6d, 0x2b, 0xb5, 0x59, 0x39, 0xfb, 0x41, 0x2a, 0xa4, 0x84, 0xff, 0x22, 0x2e, 0x61,
	0x92, 0x22, 0x63, 0x06, 0x6b, 0xae, 0xfa, 0x6d, 0xeb, 0x67, 0xad, 0x9c, 0x4d, 0x5e, 0x62, 0x60,
	0xf0, 0x1e, 0x66, 0xda, 0xf8, 0x20, 0x8d, 0x22, 0xcf, 0xf3, 0x32, 0x5f, 0xb0, 0xe6, 0xfa, 0xc0,
	0x85, 0xa7, 0xc4, 0x88, 0x3f, 0xba, 0x7a, 0x80, 0x69, 0xb7, 0xde, 0xa3, 0x0c, 0x12, 0x6b, 0x98,
	0x0e, 0xe5, 0xf0, 0x2c, 0xba, 0x60, 0x72, 0x19, 0xbd, 0x40, 0xfc, 0x32, 0xcd, 0x1d, 0x14, 0xdd,
	0x00, 0x97, 0x50, 0x74, 0x15, 0xe2, 0x40, 0x8f, 0xfa, 0x9c, 0x5f, 0x8e, 0x1c, 0xba, 0x90, 0xea,
	0xe8, 0xed, 0x34, 0x2a, 0xb7, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x48, 0xd4, 0x2e, 0x71, 0x9f,
	0x01, 0x00, 0x00,
}
