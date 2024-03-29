// Code generated by protoc-gen-go. DO NOT EDIT.
// source: models.proto

package proto

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

// Block is proto message that represent Ethereum Chain Block
// structure with most important data fields
type Block struct {
	Number               int64    `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Hash                 []byte   `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Parent               []byte   `protobuf:"bytes,3,opt,name=parent,proto3" json:"parent,omitempty"`
	Timestamp            int64    `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TxCount              int64    `protobuf:"varint,5,opt,name=txCount,proto3" json:"txCount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_models_25eeb16b30aba51f, []int{0}
}
func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (dst *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(dst, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetNumber() int64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *Block) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Block) GetParent() []byte {
	if m != nil {
		return m.Parent
	}
	return nil
}

func (m *Block) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Block) GetTxCount() int64 {
	if m != nil {
		return m.TxCount
	}
	return 0
}

// Tx is proto message that represent Ethereum Transaction
// model structure with most important data fields
type Tx struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	From                 []byte   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To                   []byte   `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Value                int64    `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	Gas                  int64    `protobuf:"varint,5,opt,name=gas,proto3" json:"gas,omitempty"`
	GasPrice             int64    `protobuf:"varint,6,opt,name=gasPrice,proto3" json:"gasPrice,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tx) Reset()         { *m = Tx{} }
func (m *Tx) String() string { return proto.CompactTextString(m) }
func (*Tx) ProtoMessage()    {}
func (*Tx) Descriptor() ([]byte, []int) {
	return fileDescriptor_models_25eeb16b30aba51f, []int{1}
}
func (m *Tx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tx.Unmarshal(m, b)
}
func (m *Tx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tx.Marshal(b, m, deterministic)
}
func (dst *Tx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tx.Merge(dst, src)
}
func (m *Tx) XXX_Size() int {
	return xxx_messageInfo_Tx.Size(m)
}
func (m *Tx) XXX_DiscardUnknown() {
	xxx_messageInfo_Tx.DiscardUnknown(m)
}

var xxx_messageInfo_Tx proto.InternalMessageInfo

func (m *Tx) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Tx) GetFrom() []byte {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Tx) GetTo() []byte {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Tx) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Tx) GetGas() int64 {
	if m != nil {
		return m.Gas
	}
	return 0
}

func (m *Tx) GetGasPrice() int64 {
	if m != nil {
		return m.GasPrice
	}
	return 0
}

func init() {
	proto.RegisterType((*Block)(nil), "models.Block")
	proto.RegisterType((*Tx)(nil), "models.Tx")
}

func init() { proto.RegisterFile("models.proto", fileDescriptor_models_25eeb16b30aba51f) }

var fileDescriptor_models_25eeb16b30aba51f = []byte{
	// 218 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0x49, 0xda, 0x64, 0x75, 0x58, 0x44, 0x06, 0x91, 0x20, 0x1e, 0x96, 0x3d, 0xed, 0xc9,
	0x8b, 0x6f, 0xb0, 0xbe, 0x80, 0x2c, 0x9e, 0xbc, 0x65, 0xd7, 0xb8, 0x2d, 0x36, 0x9d, 0x92, 0x4c,
	0xa5, 0x47, 0xc1, 0x17, 0x97, 0xa6, 0x69, 0xdd, 0x53, 0xfe, 0xef, 0xcf, 0x30, 0x1f, 0x0c, 0xac,
	0x3d, 0x7d, 0xb8, 0x26, 0x3e, 0x75, 0x81, 0x98, 0x50, 0x4f, 0xb4, 0xfd, 0x15, 0xa0, 0xf6, 0x0d,
	0x9d, 0xbe, 0xf0, 0x1e, 0x74, 0xdb, 0xfb, 0xa3, 0x0b, 0x46, 0x6c, 0xc4, 0xae, 0x38, 0x64, 0x42,
	0x84, 0xb2, 0xb2, 0xb1, 0x32, 0x72, 0x23, 0x76, 0xeb, 0x43, 0xca, 0xe3, 0x6c, 0x67, 0x83, 0x6b,
	0xd9, 0x14, 0xa9, 0xcd, 0x84, 0x8f, 0x70, 0xcd, 0xb5, 0x77, 0x91, 0xad, 0xef, 0x4c, 0x99, 0xd6,
	0xfc, 0x17, 0x68, 0x60, 0xc5, 0xc3, 0x0b, 0xf5, 0x2d, 0x1b, 0x95, 0xfe, 0x66, 0xdc, 0xfe, 0x08,
	0x90, 0x6f, 0xc3, 0xa2, 0x12, 0x17, 0x2a, 0x84, 0xf2, 0x33, 0x90, 0x9f, 0xf5, 0x63, 0xc6, 0x1b,
	0x90, 0x4c, 0x59, 0x2d, 0x99, 0xf0, 0x0e, 0xd4, 0xb7, 0x6d, 0x7a, 0x97, 0x95, 0x13, 0xe0, 0x2d,
	0x14, 0x67, 0x1b, 0xb3, 0x6a, 0x8c, 0xf8, 0x00, 0x57, 0x67, 0x1b, 0x5f, 0x43, 0x7d, 0x72, 0x46,
	0xa7, 0x7a, 0xe1, 0xfd, 0xea, 0x5d, 0xa5, 0xcb, 0x1c, 0x75, 0x7a, 0x9e, 0xff, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x35, 0xa9, 0xfe, 0xb9, 0x30, 0x01, 0x00, 0x00,
}
