// Code generated by protoc-gen-go.
// source: manage.proto
// DO NOT EDIT!

/*
Package Report is a generated protocol buffer package.

It is generated from these files:
	manage.proto
	param.proto

It has these top-level messages:
	ManageCommand
	Param
*/
package Report

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ManageCommand_CommandType int32

const (
	ManageCommand_CMT_INVALID ManageCommand_CommandType = 0
	// das->dps
	ManageCommand_CMT_REQ_REGISTER ManageCommand_CommandType = 1
	ManageCommand_CMT_REQ_LOGIN    ManageCommand_CommandType = 2
	// dps->das
	ManageCommand_CMT_REP_REGISTER ManageCommand_CommandType = 32769
	ManageCommand_CMT_REP_LOGIN    ManageCommand_CommandType = 32770
)

var ManageCommand_CommandType_name = map[int32]string{
	0:     "CMT_INVALID",
	1:     "CMT_REQ_REGISTER",
	2:     "CMT_REQ_LOGIN",
	32769: "CMT_REP_REGISTER",
	32770: "CMT_REP_LOGIN",
}
var ManageCommand_CommandType_value = map[string]int32{
	"CMT_INVALID":      0,
	"CMT_REQ_REGISTER": 1,
	"CMT_REQ_LOGIN":    2,
	"CMT_REP_REGISTER": 32769,
	"CMT_REP_LOGIN":    32770,
}

func (x ManageCommand_CommandType) String() string {
	return proto.EnumName(ManageCommand_CommandType_name, int32(x))
}
func (ManageCommand_CommandType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type ManageCommand struct {
	Cpuid        []byte                    `protobuf:"bytes,1,opt,name=cpuid,proto3" json:"cpuid,omitempty"`
	Tid          uint64                    `protobuf:"varint,2,opt,name=tid" json:"tid,omitempty"`
	SerialNumber uint32                    `protobuf:"varint,3,opt,name=serial_number,json=serialNumber" json:"serial_number,omitempty"`
	Uuid         string                    `protobuf:"bytes,4,opt,name=uuid" json:"uuid,omitempty"`
	Type         ManageCommand_CommandType `protobuf:"varint,5,opt,name=type,enum=Report.ManageCommand_CommandType" json:"type,omitempty"`
	Paras        []*Param                  `protobuf:"bytes,6,rep,name=paras" json:"paras,omitempty"`
}

func (m *ManageCommand) Reset()                    { *m = ManageCommand{} }
func (m *ManageCommand) String() string            { return proto.CompactTextString(m) }
func (*ManageCommand) ProtoMessage()               {}
func (*ManageCommand) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ManageCommand) GetParas() []*Param {
	if m != nil {
		return m.Paras
	}
	return nil
}

func init() {
	proto.RegisterType((*ManageCommand)(nil), "Report.ManageCommand")
	proto.RegisterEnum("Report.ManageCommand_CommandType", ManageCommand_CommandType_name, ManageCommand_CommandType_value)
}

func init() { proto.RegisterFile("manage.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 272 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x54, 0x90, 0xbd, 0x4e, 0xc3, 0x30,
	0x10, 0x80, 0xc9, 0xaf, 0x84, 0x93, 0x40, 0x38, 0x10, 0x8a, 0x98, 0x4a, 0x59, 0x3a, 0x65, 0x28,
	0xe2, 0x01, 0x10, 0x54, 0x55, 0xa4, 0x36, 0x04, 0x13, 0xb1, 0x56, 0x2e, 0xb5, 0x50, 0x25, 0xd2,
	0x58, 0xae, 0x33, 0xb0, 0x15, 0x1e, 0x8d, 0x27, 0xe3, 0x6c, 0x37, 0x82, 0x0e, 0x51, 0xee, 0xbe,
	0xfb, 0xce, 0xbe, 0x33, 0x89, 0x1b, 0xb6, 0x61, 0xef, 0x3c, 0x17, 0xb2, 0x55, 0x2d, 0x84, 0x94,
	0x8b, 0x56, 0xaa, 0xab, 0x48, 0x30, 0xc9, 0x1a, 0x0b, 0x87, 0x3f, 0x2e, 0x49, 0xe6, 0xc6, 0x7a,
	0x68, 0x1b, 0xd4, 0x57, 0x70, 0x41, 0x82, 0x37, 0xd1, 0xad, 0x57, 0x99, 0x33, 0x70, 0x46, 0x31,
	0xb5, 0x09, 0xa4, 0xc4, 0x53, 0xc8, 0x5c, 0x64, 0x3e, 0xd5, 0x21, 0xdc, 0x90, 0x64, 0xcb, 0xe5,
	0x9a, 0x7d, 0x2c, 0x36, 0x5d, 0xb3, 0xe4, 0x32, 0xf3, 0xb0, 0x96, 0xd0, 0xd8, 0xc2, 0xd2, 0x30,
	0x00, 0xe2, 0x77, 0xfa, 0x2c, 0x1f, 0x6b, 0xc7, 0xd4, 0xc4, 0x70, 0x47, 0x7c, 0xf5, 0x29, 0x78,
	0x16, 0x20, 0x3b, 0x19, 0x5f, 0xe7, 0x76, 0xac, 0xfc, 0x60, 0x8a, 0x7c, 0xff, 0xaf, 0x51, 0xa4,
	0x46, 0xc7, 0xfb, 0x02, 0x3d, 0xf8, 0x36, 0x0b, 0x07, 0xde, 0x28, 0x1a, 0x27, 0x7d, 0x5f, 0xa5,
	0xb7, 0xa1, 0xb6, 0x36, 0x54, 0x24, 0xfa, 0xd7, 0x09, 0xa7, 0x98, 0xce, 0xeb, 0x45, 0x51, 0xbe,
	0xde, 0xcf, 0x8a, 0xc7, 0xf4, 0x08, 0x97, 0x4b, 0x35, 0xa0, 0x93, 0x67, 0xfc, 0xa6, 0xc5, 0x4b,
	0x3d, 0xa1, 0xa9, 0x03, 0x67, 0x24, 0xe9, 0xe9, 0xec, 0x69, 0x5a, 0x94, 0xa9, 0x0b, 0x97, 0xbd,
	0x58, 0xfd, 0x89, 0x5f, 0x3b, 0x17, 0xce, 0x7b, 0xb5, 0xda, 0xab, 0xdf, 0x3b, 0x77, 0x19, 0x9a,
	0xb7, 0xbc, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xa3, 0x18, 0x51, 0x5a, 0x70, 0x01, 0x00, 0x00,
}
