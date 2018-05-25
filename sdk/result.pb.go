// Code generated by protoc-gen-go.
// source: result.proto
// DO NOT EDIT!

package sdk

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Result represents an execution Result
// Generate *.pb.go files with:
// 	protoc --go_out=plugins=grpc:. ./result.proto
// 	protoc-go-inject-tag -input=./log.pb.go
// 	=> github.com/favadi/protoc-go-inject-tag
type Result struct {
	Id         int64                      `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	BuildID    int64                      `protobuf:"varint,2,opt,name=buildID" json:"buildID,omitempty"`
	Status     string                     `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
	Version    int64                      `protobuf:"varint,4,opt,name=version" json:"version,omitempty"`
	Reason     string                     `protobuf:"bytes,5,opt,name=reason" json:"reason,omitempty"`
	RemoteTime *google_protobuf.Timestamp `protobuf:"bytes,6,opt,name=remoteTime" json:"remoteTime,omitempty"`
	Duration   string                     `protobuf:"bytes,7,opt,name=duration" json:"duration,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Result) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Result) GetBuildID() int64 {
	if m != nil {
		return m.BuildID
	}
	return 0
}

func (m *Result) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Result) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Result) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *Result) GetRemoteTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.RemoteTime
	}
	return nil
}

func (m *Result) GetDuration() string {
	if m != nil {
		return m.Duration
	}
	return ""
}

func init() {
	proto.RegisterType((*Result)(nil), "github.com.ovh.cds.sdk.Result")
}

func init() { proto.RegisterFile("result.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x4c, 0x8e, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0xe5, 0x84, 0xa6, 0x60, 0x10, 0x83, 0x87, 0xca, 0xca, 0x42, 0xc4, 0x94, 0xc9, 0x95,
	0x60, 0x63, 0x44, 0x2c, 0xac, 0x51, 0x27, 0xb6, 0xa4, 0x3e, 0x52, 0xab, 0x71, 0xaf, 0xf2, 0x9d,
	0xfb, 0x4b, 0xf9, 0x41, 0x28, 0x76, 0x83, 0x18, 0xbf, 0xd3, 0xfb, 0xde, 0x3d, 0xf9, 0x10, 0x80,
	0xe2, 0xc4, 0xe6, 0x1c, 0x90, 0x51, 0x6d, 0x46, 0xc7, 0x87, 0x38, 0x98, 0x3d, 0x7a, 0x83, 0x97,
	0x83, 0xd9, 0x5b, 0x32, 0x64, 0x8f, 0xf5, 0xd3, 0x88, 0x38, 0x4e, 0xb0, 0x4d, 0xa9, 0x21, 0x7e,
	0x6f, 0xd9, 0x79, 0x20, 0xee, 0xfd, 0x39, 0x8b, 0xcf, 0x3f, 0x42, 0x56, 0x5d, 0x6a, 0x52, 0x8f,
	0xb2, 0x70, 0x56, 0x8b, 0x46, 0xb4, 0x65, 0x57, 0x38, 0xab, 0xb4, 0x5c, 0x0f, 0xd1, 0x4d, 0xf6,
	0xf3, 0x43, 0x17, 0xe9, 0xb8, 0xa0, 0xda, 0xc8, 0x8a, 0xb8, 0xe7, 0x48, 0xba, 0x6c, 0x44, 0x7b,
	0xd7, 0x5d, 0x69, 0x36, 0x2e, 0x10, 0xc8, 0xe1, 0x49, 0xdf, 0x64, 0xe3, 0x8a, 0xb3, 0x11, 0xa0,
	0x27, 0x3c, 0xe9, 0x55, 0x36, 0x32, 0xa9, 0x37, 0x29, 0x03, 0x78, 0x64, 0xd8, 0x39, 0x0f, 0xba,
	0x6a, 0x44, 0x7b, 0xff, 0x52, 0x9b, 0x3c, 0xda, 0x2c, 0xa3, 0xcd, 0x6e, 0x19, 0xdd, 0xfd, 0x4b,
	0xab, 0x5a, 0xde, 0xda, 0x18, 0x7a, 0x9e, 0xdf, 0xad, 0x53, 0xeb, 0x1f, 0xbf, 0xaf, 0xbe, 0x4a,
	0xb2, 0xc7, 0xa1, 0x4a, 0x15, 0xaf, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6e, 0xb8, 0xaf, 0xc7,
	0x2d, 0x01, 0x00, 0x00,
}