// Code generated by protoc-gen-go. DO NOT EDIT.
// source: annotated.proto

package testdata

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Foo struct {
	A                    string   `protobuf:"bytes,1,opt,name=a,proto3" json:"a,omitempty"`
	Bar                  *Bar     `protobuf:"bytes,2,opt,name=bar,proto3" json:"bar,omitempty"`
	Req                  string   `protobuf:"bytes,3,opt,name=req,proto3" json:"req,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Foo) Reset()         { *m = Foo{} }
func (m *Foo) String() string { return proto.CompactTextString(m) }
func (*Foo) ProtoMessage()    {}
func (*Foo) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{0}
}

func (m *Foo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Foo.Unmarshal(m, b)
}
func (m *Foo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Foo.Marshal(b, m, deterministic)
}
func (m *Foo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Foo.Merge(m, src)
}
func (m *Foo) XXX_Size() int {
	return xxx_messageInfo_Foo.Size(m)
}
func (m *Foo) XXX_DiscardUnknown() {
	xxx_messageInfo_Foo.DiscardUnknown(m)
}

var xxx_messageInfo_Foo proto.InternalMessageInfo

func (m *Foo) GetA() string {
	if m != nil {
		return m.A
	}
	return ""
}

func (m *Foo) GetBar() *Bar {
	if m != nil {
		return m.Bar
	}
	return nil
}

func (m *Foo) GetReq() string {
	if m != nil {
		return m.Req
	}
	return ""
}

type Bar struct {
	B                    string   `protobuf:"bytes,1,opt,name=b,proto3" json:"b,omitempty"`
	Baz                  *Baz     `protobuf:"bytes,2,opt,name=baz,proto3" json:"baz,omitempty"`
	A                    string   `protobuf:"bytes,3,opt,name=a,proto3" json:"a,omitempty"`
	Def                  string   `protobuf:"bytes,4,opt,name=def,proto3" json:"def,omitempty"`
	Remote               string   `protobuf:"bytes,5,opt,name=remote,proto3" json:"remote,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bar) Reset()         { *m = Bar{} }
func (m *Bar) String() string { return proto.CompactTextString(m) }
func (*Bar) ProtoMessage()    {}
func (*Bar) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{1}
}

func (m *Bar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bar.Unmarshal(m, b)
}
func (m *Bar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bar.Marshal(b, m, deterministic)
}
func (m *Bar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bar.Merge(m, src)
}
func (m *Bar) XXX_Size() int {
	return xxx_messageInfo_Bar.Size(m)
}
func (m *Bar) XXX_DiscardUnknown() {
	xxx_messageInfo_Bar.DiscardUnknown(m)
}

var xxx_messageInfo_Bar proto.InternalMessageInfo

func (m *Bar) GetB() string {
	if m != nil {
		return m.B
	}
	return ""
}

func (m *Bar) GetBaz() *Baz {
	if m != nil {
		return m.Baz
	}
	return nil
}

func (m *Bar) GetA() string {
	if m != nil {
		return m.A
	}
	return ""
}

func (m *Bar) GetDef() string {
	if m != nil {
		return m.Def
	}
	return ""
}

func (m *Bar) GetRemote() string {
	if m != nil {
		return m.Remote
	}
	return ""
}

type Baz struct {
	C                    string   `protobuf:"bytes,1,opt,name=c,proto3" json:"c,omitempty"`
	Biz                  []*Biz   `protobuf:"bytes,2,rep,name=biz,proto3" json:"biz,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Baz) Reset()         { *m = Baz{} }
func (m *Baz) String() string { return proto.CompactTextString(m) }
func (*Baz) ProtoMessage()    {}
func (*Baz) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{2}
}

func (m *Baz) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Baz.Unmarshal(m, b)
}
func (m *Baz) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Baz.Marshal(b, m, deterministic)
}
func (m *Baz) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Baz.Merge(m, src)
}
func (m *Baz) XXX_Size() int {
	return xxx_messageInfo_Baz.Size(m)
}
func (m *Baz) XXX_DiscardUnknown() {
	xxx_messageInfo_Baz.DiscardUnknown(m)
}

var xxx_messageInfo_Baz proto.InternalMessageInfo

func (m *Baz) GetC() string {
	if m != nil {
		return m.C
	}
	return ""
}

func (m *Baz) GetBiz() []*Biz {
	if m != nil {
		return m.Biz
	}
	return nil
}

type Biz struct {
	D                    string   `protobuf:"bytes,1,opt,name=d,proto3" json:"d,omitempty"`
	E                    string   `protobuf:"bytes,2,opt,name=e,proto3" json:"e,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Biz) Reset()         { *m = Biz{} }
func (m *Biz) String() string { return proto.CompactTextString(m) }
func (*Biz) ProtoMessage()    {}
func (*Biz) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{3}
}

func (m *Biz) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Biz.Unmarshal(m, b)
}
func (m *Biz) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Biz.Marshal(b, m, deterministic)
}
func (m *Biz) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Biz.Merge(m, src)
}
func (m *Biz) XXX_Size() int {
	return xxx_messageInfo_Biz.Size(m)
}
func (m *Biz) XXX_DiscardUnknown() {
	xxx_messageInfo_Biz.DiscardUnknown(m)
}

var xxx_messageInfo_Biz proto.InternalMessageInfo

func (m *Biz) GetD() string {
	if m != nil {
		return m.D
	}
	return ""
}

func (m *Biz) GetE() string {
	if m != nil {
		return m.E
	}
	return ""
}

type Qux struct {
	Req                  string   `protobuf:"bytes,1,opt,name=req,proto3" json:"req,omitempty"`
	E                    string   `protobuf:"bytes,2,opt,name=e,proto3" json:"e,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Qux) Reset()         { *m = Qux{} }
func (m *Qux) String() string { return proto.CompactTextString(m) }
func (*Qux) ProtoMessage()    {}
func (*Qux) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{4}
}

func (m *Qux) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Qux.Unmarshal(m, b)
}
func (m *Qux) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Qux.Marshal(b, m, deterministic)
}
func (m *Qux) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Qux.Merge(m, src)
}
func (m *Qux) XXX_Size() int {
	return xxx_messageInfo_Qux.Size(m)
}
func (m *Qux) XXX_DiscardUnknown() {
	xxx_messageInfo_Qux.DiscardUnknown(m)
}

var xxx_messageInfo_Qux proto.InternalMessageInfo

func (m *Qux) GetReq() string {
	if m != nil {
		return m.Req
	}
	return ""
}

func (m *Qux) GetE() string {
	if m != nil {
		return m.E
	}
	return ""
}

type Wibble struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Wibble) Reset()         { *m = Wibble{} }
func (m *Wibble) String() string { return proto.CompactTextString(m) }
func (*Wibble) ProtoMessage()    {}
func (*Wibble) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{5}
}

func (m *Wibble) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Wibble.Unmarshal(m, b)
}
func (m *Wibble) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Wibble.Marshal(b, m, deterministic)
}
func (m *Wibble) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Wibble.Merge(m, src)
}
func (m *Wibble) XXX_Size() int {
	return xxx_messageInfo_Wibble.Size(m)
}
func (m *Wibble) XXX_DiscardUnknown() {
	xxx_messageInfo_Wibble.DiscardUnknown(m)
}

var xxx_messageInfo_Wibble proto.InternalMessageInfo

type Wobble struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Wobble) Reset()         { *m = Wobble{} }
func (m *Wobble) String() string { return proto.CompactTextString(m) }
func (*Wobble) ProtoMessage()    {}
func (*Wobble) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{6}
}

func (m *Wobble) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Wobble.Unmarshal(m, b)
}
func (m *Wobble) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Wobble.Marshal(b, m, deterministic)
}
func (m *Wobble) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Wobble.Merge(m, src)
}
func (m *Wobble) XXX_Size() int {
	return xxx_messageInfo_Wobble.Size(m)
}
func (m *Wobble) XXX_DiscardUnknown() {
	xxx_messageInfo_Wobble.DiscardUnknown(m)
}

var xxx_messageInfo_Wobble proto.InternalMessageInfo

func (m *Wobble) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Wubble struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Wubble) Reset()         { *m = Wubble{} }
func (m *Wubble) String() string { return proto.CompactTextString(m) }
func (*Wubble) ProtoMessage()    {}
func (*Wubble) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{7}
}

func (m *Wubble) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Wubble.Unmarshal(m, b)
}
func (m *Wubble) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Wubble.Marshal(b, m, deterministic)
}
func (m *Wubble) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Wubble.Merge(m, src)
}
func (m *Wubble) XXX_Size() int {
	return xxx_messageInfo_Wubble.Size(m)
}
func (m *Wubble) XXX_DiscardUnknown() {
	xxx_messageInfo_Wubble.DiscardUnknown(m)
}

var xxx_messageInfo_Wubble proto.InternalMessageInfo

func (m *Wubble) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Flob struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Flob) Reset()         { *m = Flob{} }
func (m *Flob) String() string { return proto.CompactTextString(m) }
func (*Flob) ProtoMessage()    {}
func (*Flob) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9e975ae609bc8c7, []int{8}
}

func (m *Flob) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Flob.Unmarshal(m, b)
}
func (m *Flob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Flob.Marshal(b, m, deterministic)
}
func (m *Flob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Flob.Merge(m, src)
}
func (m *Flob) XXX_Size() int {
	return xxx_messageInfo_Flob.Size(m)
}
func (m *Flob) XXX_DiscardUnknown() {
	xxx_messageInfo_Flob.DiscardUnknown(m)
}

var xxx_messageInfo_Flob proto.InternalMessageInfo

func (m *Flob) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Foo)(nil), "annotated.Foo")
	proto.RegisterType((*Bar)(nil), "annotated.Bar")
	proto.RegisterType((*Baz)(nil), "annotated.Baz")
	proto.RegisterType((*Biz)(nil), "annotated.Biz")
	proto.RegisterType((*Qux)(nil), "annotated.Qux")
	proto.RegisterType((*Wibble)(nil), "annotated.Wibble")
	proto.RegisterType((*Wobble)(nil), "annotated.Wobble")
	proto.RegisterType((*Wubble)(nil), "annotated.Wubble")
	proto.RegisterType((*Flob)(nil), "annotated.Flob")
}

func init() { proto.RegisterFile("annotated.proto", fileDescriptor_e9e975ae609bc8c7) }

var fileDescriptor_e9e975ae609bc8c7 = []byte{
	// 570 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0xbf, 0x6e, 0xd4, 0x4e,
	0x10, 0xc7, 0x7f, 0x1b, 0x27, 0x27, 0x65, 0x7f, 0x28, 0x7f, 0x36, 0x89, 0x72, 0x9c, 0x40, 0x09,
	0xae, 0xd2, 0xe4, 0x16, 0x12, 0x68, 0x02, 0x8d, 0x5d, 0x44, 0x42, 0x02, 0xa1, 0x98, 0x22, 0x12,
	0xcd, 0x69, 0xd6, 0xbb, 0x76, 0x16, 0x7c, 0x1e, 0xb3, 0x5e, 0x47, 0xc1, 0xa7, 0x13, 0x05, 0x8f,
	0x41, 0x41, 0xc3, 0x73, 0xf0, 0x2c, 0xd4, 0xf7, 0x08, 0x57, 0xa1, 0xb5, 0x1d, 0x82, 0xc8, 0xa5,
	0xa4, 0x9b, 0x9b, 0x99, 0xef, 0xcc, 0x67, 0xee, 0xbb, 0xa6, 0xeb, 0x90, 0xe7, 0x68, 0xc1, 0x2a,
	0x39, 0x2c, 0x0c, 0x5a, 0x64, 0xab, 0xbf, 0x13, 0x83, 0xbd, 0x14, 0x31, 0xcd, 0x14, 0x87, 0x42,
	0xf3, 0x44, 0xab, 0x4c, 0x8e, 0x84, 0xba, 0x80, 0x4b, 0x8d, 0xa6, 0xed, 0x1d, 0xdc, 0xff, 0xa3,
	0xc1, 0xa8, 0x12, 0x2b, 0x13, 0xab, 0xb6, 0xe4, 0x7f, 0x21, 0xd4, 0x3b, 0x45, 0x64, 0xf7, 0x28,
	0x81, 0x3e, 0xd9, 0x27, 0x07, 0xab, 0x11, 0x01, 0xb6, 0x4f, 0x3d, 0x01, 0xa6, 0xbf, 0xb4, 0x4f,
	0x0e, 0xfe, 0x3f, 0x5a, 0x1b, 0xde, 0xec, 0x0e, 0xc1, 0x44, 0xae, 0xc4, 0x76, 0xa8, 0x67, 0xd4,
	0xc7, 0xbe, 0xe7, 0x14, 0xa1, 0xf7, 0x33, 0x58, 0x8a, 0xdc, 0xef, 0x93, 0x67, 0xb3, 0xe0, 0x88,
	0xae, 0x27, 0x88, 0x43, 0x01, 0x66, 0x18, 0xe3, 0x98, 0xbb, 0xe1, 0x0f, 0x0a, 0x83, 0xef, 0x55,
	0x6c, 0x4b, 0x3e, 0xe9, 0xa2, 0x29, 0x4f, 0x10, 0xf9, 0x24, 0x41, 0x9c, 0x0e, 0x08, 0xf8, 0xdf,
	0x09, 0xf5, 0x42, 0x30, 0x6c, 0x93, 0x12, 0xd1, 0x52, 0xb8, 0x99, 0x24, 0x22, 0xa2, 0x45, 0xa9,
	0x17, 0xa2, 0xd4, 0x0e, 0xa5, 0x66, 0xbe, 0x43, 0x6f, 0x41, 0xb6, 0xe7, 0xc1, 0xe6, 0xad, 0xf5,
	0xee, 0xa0, 0x0d, 0xea, 0x49, 0x95, 0xf4, 0x97, 0x9b, 0x03, 0x5d, 0xc8, 0x8e, 0x68, 0xcf, 0xa8,
	0x31, 0x5a, 0xd5, 0x5f, 0x69, 0xa4, 0x83, 0x79, 0xb0, 0xcb, 0x76, 0x12, 0x44, 0xa7, 0x84, 0x42,
	0x5f, 0xab, 0xdd, 0xc5, 0x5d, 0xa7, 0xff, 0xc2, 0x51, 0xd6, 0x6c, 0x8b, 0x92, 0xb8, 0xa3, 0x5c,
	0x99, 0x07, 0x4b, 0xf4, 0xbf, 0x88, 0xc4, 0x0d, 0xa7, 0x76, 0x9c, 0xde, 0xdf, 0x9c, 0xda, 0x71,
	0xea, 0xda, 0x7f, 0x4d, 0xbd, 0x50, 0x37, 0xb8, 0xb2, 0x53, 0xdf, 0xc6, 0x0d, 0xab, 0x3a, 0x22,
	0xd2, 0xf5, 0xa8, 0xe6, 0xe4, 0x45, 0x3d, 0x67, 0xd5, 0x55, 0x44, 0x94, 0xff, 0x92, 0x7a, 0x67,
	0xd5, 0x15, 0x7b, 0xda, 0x1a, 0xd1, 0x0e, 0xf4, 0xe7, 0xc1, 0x1e, 0x7d, 0x18, 0x67, 0x58, 0xc9,
	0x61, 0xeb, 0x39, 0x14, 0xba, 0x6c, 0x54, 0xaf, 0x30, 0x06, 0xab, 0x31, 0x6f, 0x7c, 0x72, 0x76,
	0x77, 0x0b, 0xdc, 0xa8, 0x47, 0xb4, 0x77, 0xae, 0x85, 0xc8, 0xd4, 0xc9, 0xee, 0x2c, 0xd8, 0xa6,
	0x0c, 0x44, 0x3c, 0x94, 0x2a, 0x69, 0x74, 0x6d, 0xc1, 0x7f, 0x42, 0x7b, 0xe7, 0xe8, 0x22, 0xc6,
	0xe8, 0x72, 0x0e, 0x63, 0xd5, 0x3d, 0x96, 0x26, 0x3e, 0xd9, 0x9a, 0x05, 0x1b, 0x6c, 0x6d, 0xac,
	0xcb, 0x52, 0xe7, 0x29, 0x9f, 0xd8, 0x4f, 0x85, 0x9a, 0xfa, 0x3f, 0x08, 0xed, 0x9d, 0x57, 0x77,
	0x6a, 0xbe, 0x91, 0x59, 0xf0, 0x95, 0xd0, 0xcf, 0x02, 0xe4, 0xf5, 0x9b, 0x1c, 0x39, 0xe1, 0xe8,
	0x83, 0xce, 0x25, 0x5f, 0x98, 0xfd, 0x17, 0x49, 0xd6, 0x17, 0x20, 0xf9, 0xe4, 0x76, 0x61, 0xea,
	0xbf, 0xa1, 0xcb, 0xa7, 0x19, 0x8a, 0x85, 0xf4, 0xc7, 0xb3, 0xe0, 0x31, 0xdd, 0xee, 0x2e, 0x1e,
	0x95, 0xca, 0x5c, 0xea, 0x58, 0x8d, 0x5c, 0xe9, 0xee, 0x81, 0xe1, 0xdb, 0x77, 0x67, 0xa9, 0xb6,
	0x17, 0x95, 0x68, 0xfe, 0xd9, 0x1b, 0x83, 0x78, 0x0a, 0x85, 0x8e, 0x0f, 0x63, 0xcc, 0x13, 0x9d,
	0x1e, 0x5e, 0x42, 0xa6, 0x25, 0x58, 0x34, 0x5c, 0xe7, 0x56, 0x99, 0x1c, 0x32, 0x7e, 0x93, 0xb2,
	0xaa, 0xb4, 0x12, 0x2c, 0x3c, 0xbf, 0x0e, 0x44, 0xaf, 0xf9, 0x90, 0x8f, 0x7f, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x92, 0xb5, 0x60, 0x15, 0x22, 0x04, 0x00, 0x00,
}
