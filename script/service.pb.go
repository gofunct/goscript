// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: script/service.proto

package script

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Command struct {
	Dir  string   `protobuf:"bytes,1,opt,name=dir,proto3" json:"dir,omitempty"`
	Args []string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_4b3fcbd1c8b7cd18, []int{0}
}
func (m *Command) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Command.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(dst, src)
}
func (m *Command) XXX_Size() int {
	return m.Size()
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetDir() string {
	if m != nil {
		return m.Dir
	}
	return ""
}

func (m *Command) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

type Output struct {
	Retcode int32  `protobuf:"varint,1,opt,name=retcode,proto3" json:"retcode,omitempty"`
	Out     []byte `protobuf:"bytes,2,opt,name=out,proto3" json:"out,omitempty"`
}

func (m *Output) Reset()         { *m = Output{} }
func (m *Output) String() string { return proto.CompactTextString(m) }
func (*Output) ProtoMessage()    {}
func (*Output) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_4b3fcbd1c8b7cd18, []int{1}
}
func (m *Output) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Output) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Output.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Output) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Output.Merge(dst, src)
}
func (m *Output) XXX_Size() int {
	return m.Size()
}
func (m *Output) XXX_DiscardUnknown() {
	xxx_messageInfo_Output.DiscardUnknown(m)
}

var xxx_messageInfo_Output proto.InternalMessageInfo

func (m *Output) GetRetcode() int32 {
	if m != nil {
		return m.Retcode
	}
	return 0
}

func (m *Output) GetOut() []byte {
	if m != nil {
		return m.Out
	}
	return nil
}

func init() {
	proto.RegisterType((*Command)(nil), "script.Command")
	proto.RegisterType((*Output)(nil), "script.Output")
}
func (m *Command) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Command) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Dir) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintService(dAtA, i, uint64(len(m.Dir)))
		i += copy(dAtA[i:], m.Dir)
	}
	if len(m.Args) > 0 {
		for _, s := range m.Args {
			dAtA[i] = 0x12
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func (m *Output) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Output) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Retcode != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintService(dAtA, i, uint64(m.Retcode))
	}
	if len(m.Out) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintService(dAtA, i, uint64(len(m.Out)))
		i += copy(dAtA[i:], m.Out)
	}
	return i, nil
}

func encodeVarintService(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Command) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Dir)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	if len(m.Args) > 0 {
		for _, s := range m.Args {
			l = len(s)
			n += 1 + l + sovService(uint64(l))
		}
	}
	return n
}

func (m *Output) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Retcode != 0 {
		n += 1 + sovService(uint64(m.Retcode))
	}
	l = len(m.Out)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	return n
}

func sovService(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozService(x uint64) (n int) {
	return sovService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Command) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Command: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Command: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Dir", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Dir = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Args = append(m.Args, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Output) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Output: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Output: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Retcode", wireType)
			}
			m.Retcode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Retcode |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Out", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Out = append(m.Out[:0], dAtA[iNdEx:postIndex]...)
			if m.Out == nil {
				m.Out = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthService
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowService
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipService(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthService = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowService   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("script/service.proto", fileDescriptor_service_4b3fcbd1c8b7cd18) }

var fileDescriptor_service_4b3fcbd1c8b7cd18 = []byte{
	// 262 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xbf, 0x4e, 0x85, 0x30,
	0x14, 0xc6, 0x29, 0x28, 0x84, 0x13, 0xff, 0xa5, 0x71, 0x68, 0x1c, 0x1a, 0xc2, 0x22, 0x89, 0x86,
	0x9b, 0xa8, 0x4f, 0xa0, 0x26, 0x0c, 0x0e, 0x1a, 0xae, 0x2f, 0xc0, 0x2d, 0xd5, 0x10, 0x2f, 0x96,
	0x1c, 0x5a, 0x9f, 0xc3, 0xb7, 0xd2, 0xf1, 0x8e, 0x8e, 0x06, 0x5e, 0xc4, 0x14, 0x64, 0xef, 0xf6,
	0xf5, 0xe4, 0xd7, 0x7e, 0xbf, 0xf4, 0xc0, 0x69, 0x2f, 0xb0, 0xe9, 0xf4, 0xaa, 0x97, 0xf8, 0xd1,
	0x08, 0x99, 0x77, 0xa8, 0xb4, 0xa2, 0xe1, 0x3c, 0x4d, 0x57, 0x10, 0xdd, 0xa9, 0xb6, 0xad, 0xde,
	0x6b, 0x7a, 0x02, 0x41, 0xdd, 0x20, 0x23, 0x09, 0xc9, 0xe2, 0xd2, 0x46, 0x4a, 0x61, 0xaf, 0xc2,
	0xd7, 0x9e, 0xf9, 0x49, 0x90, 0xc5, 0xe5, 0x94, 0xd3, 0x1b, 0x08, 0x1f, 0x8d, 0xee, 0x8c, 0xa6,
	0x0c, 0x22, 0x94, 0x5a, 0xa8, 0x5a, 0x4e, 0x77, 0xf6, 0xcb, 0xe5, 0x68, 0x5f, 0x52, 0x46, 0x33,
	0x3f, 0x21, 0xd9, 0x41, 0x69, 0xe3, 0xd5, 0x97, 0x0f, 0x87, 0xeb, 0xa9, 0x71, 0x3d, 0x6b, 0xd0,
	0x0b, 0x08, 0xef, 0x95, 0x78, 0x93, 0x48, 0x8f, 0xf3, 0xd9, 0x25, 0xff, 0x17, 0x39, 0x3b, 0x5a,
	0x06, 0x73, 0x51, 0xea, 0x59, 0xb8, 0x10, 0x5b, 0x65, 0x6a, 0x17, 0x38, 0x83, 0xa0, 0x68, 0xb4,
	0x0b, 0x79, 0x0e, 0x7e, 0xa1, 0x5c, 0xc0, 0x4b, 0x88, 0x1e, 0xcc, 0x46, 0x0a, 0xbd, 0x75, 0xb4,
	0x7d, 0xb2, 0x9f, 0x2c, 0x5c, 0xe0, 0x1c, 0xe2, 0x67, 0x89, 0x58, 0xbd, 0x28, 0x6c, 0x1d, 0xf8,
	0x5b, 0xf6, 0x3d, 0x70, 0xb2, 0x1b, 0x38, 0xf9, 0x1d, 0x38, 0xf9, 0x1c, 0xb9, 0xb7, 0x1b, 0xb9,
	0xf7, 0x33, 0x72, 0x6f, 0x13, 0x4e, 0x9b, 0xbd, 0xfe, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x63, 0xd6,
	0x85, 0x53, 0xf1, 0x01, 0x00, 0x00,
}
