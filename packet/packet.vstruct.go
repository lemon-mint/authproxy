package packet

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

type _ = strings.Builder
type _ = unsafe.Pointer

var _ = math.Float32frombits
var _ = math.Float64frombits
var _ = strconv.FormatInt
var _ = strconv.FormatUint
var _ = strconv.FormatFloat
var _ = fmt.Sprint

type PacketType uint8

const (
	PacketType_Connect         PacketType = 0
	PacketType_ConnectResponse PacketType = 1
	PacketType_Upstream        PacketType = 2
	PacketType_Downstream      PacketType = 3
	PacketType_PublicKey       PacketType = 4
)

func (e PacketType) String() string {
	switch e {
	case PacketType_Connect:
		return "Connect"
	case PacketType_ConnectResponse:
		return "ConnectResponse"
	case PacketType_Upstream:
		return "Upstream"
	case PacketType_Downstream:
		return "Downstream"
	case PacketType_PublicKey:
		return "PublicKey"
	}
	return ""
}

func (e PacketType) Match(
	onConnect func(),
	onConnectResponse func(),
	onUpstream func(),
	onDownstream func(),
	onPublicKey func(),
) {
	switch e {
	case PacketType_Connect:
		onConnect()
	case PacketType_ConnectResponse:
		onConnectResponse()
	case PacketType_Upstream:
		onUpstream()
	case PacketType_Downstream:
		onDownstream()
	case PacketType_PublicKey:
		onPublicKey()
	}
}

type Packet []byte

func (s Packet) Type() PacketType {
	return PacketType(s[0])
}

func (s Packet) Payload() []byte {
	_ = s[8]
	var __off0 uint64 = 9
	var __off1 uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	return []byte(s[__off0:__off1])
}

func (s Packet) Vstruct_Validate() bool {
	if len(s) < 9 {
		return false
	}

	var __off0 uint64 = 9
	var __off1 uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	var __off2 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2
}

func (s Packet) String() string {
	if !s.Vstruct_Validate() {
		return "Packet (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("Packet {")
	__b.WriteString("Type: ")
	__b.WriteString(s.Type().String())
	__b.WriteString(", ")
	__b.WriteString("Payload: ")
	__b.WriteString(fmt.Sprint(s.Payload()))
	__b.WriteString("}")
	return __b.String()
}

type Connect []byte

func (s Connect) Host() string {
	_ = s[7]
	var __off0 uint64 = 8
	var __off1 uint64 = uint64(s[0]) |
		uint64(s[1])<<8 |
		uint64(s[2])<<16 |
		uint64(s[3])<<24 |
		uint64(s[4])<<32 |
		uint64(s[5])<<40 |
		uint64(s[6])<<48 |
		uint64(s[7])<<56
	var __v = s[__off0:__off1]

	return *(*string)(unsafe.Pointer(&__v))
}

func (s Connect) Vstruct_Validate() bool {
	if len(s) < 8 {
		return false
	}

	var __off0 uint64 = 8
	var __off1 uint64 = uint64(s[0]) |
		uint64(s[1])<<8 |
		uint64(s[2])<<16 |
		uint64(s[3])<<24 |
		uint64(s[4])<<32 |
		uint64(s[5])<<40 |
		uint64(s[6])<<48 |
		uint64(s[7])<<56
	var __off2 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2
}

func (s Connect) String() string {
	if !s.Vstruct_Validate() {
		return "Connect (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("Connect {")
	__b.WriteString("Host: ")
	__b.WriteString(strconv.Quote(s.Host()))
	__b.WriteString("}")
	return __b.String()
}

type ConnectResponse []byte

func (s ConnectResponse) Success() bool {
	return bool(s[0] != 0)
}

func (s ConnectResponse) Vstruct_Validate() bool {
	return len(s) >= 1
}

func (s ConnectResponse) String() string {
	if !s.Vstruct_Validate() {
		return "ConnectResponse (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("ConnectResponse {")
	__b.WriteString("Success: ")
	__b.WriteString(strconv.FormatBool(s.Success()))
	__b.WriteString("}")
	return __b.String()
}

type Stream []byte

func (s Stream) Error() bool {
	return bool(s[0] != 0)
}

func (s Stream) Data() []byte {
	_ = s[8]
	var __off0 uint64 = 9
	var __off1 uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	return []byte(s[__off0:__off1])
}

func (s Stream) Vstruct_Validate() bool {
	if len(s) < 9 {
		return false
	}

	var __off0 uint64 = 9
	var __off1 uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	var __off2 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2
}

func (s Stream) String() string {
	if !s.Vstruct_Validate() {
		return "Stream (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("Stream {")
	__b.WriteString("Error: ")
	__b.WriteString(strconv.FormatBool(s.Error()))
	__b.WriteString(", ")
	__b.WriteString("Data: ")
	__b.WriteString(fmt.Sprint(s.Data()))
	__b.WriteString("}")
	return __b.String()
}

func Serialize_Packet(dst Packet, Type PacketType, Payload []byte) Packet {
	_ = dst[8]
	dst[0] = byte(Type)

	var __index = uint64(9)
	__tmp_1 := uint64(len(Payload)) + __index
	dst[1] = byte(__tmp_1)
	dst[2] = byte(__tmp_1 >> 8)
	dst[3] = byte(__tmp_1 >> 16)
	dst[4] = byte(__tmp_1 >> 24)
	dst[5] = byte(__tmp_1 >> 32)
	dst[6] = byte(__tmp_1 >> 40)
	dst[7] = byte(__tmp_1 >> 48)
	dst[8] = byte(__tmp_1 >> 56)
	copy(dst[__index:__tmp_1], Payload)
	return dst
}

func New_Packet(Type PacketType, Payload []byte) Packet {
	var __vstruct__size = 9 + len(Payload)
	var __vstruct__buf = make(Packet, __vstruct__size)
	__vstruct__buf = Serialize_Packet(__vstruct__buf, Type, Payload)
	return __vstruct__buf
}

func Serialize_Connect(dst Connect, Host string) Connect {
	_ = dst[7]

	var __index = uint64(8)
	__tmp_0 := uint64(len(Host)) + __index
	dst[0] = byte(__tmp_0)
	dst[1] = byte(__tmp_0 >> 8)
	dst[2] = byte(__tmp_0 >> 16)
	dst[3] = byte(__tmp_0 >> 24)
	dst[4] = byte(__tmp_0 >> 32)
	dst[5] = byte(__tmp_0 >> 40)
	dst[6] = byte(__tmp_0 >> 48)
	dst[7] = byte(__tmp_0 >> 56)
	copy(dst[__index:__tmp_0], Host)
	return dst
}

func New_Connect(Host string) Connect {
	var __vstruct__size = 8 + len(Host)
	var __vstruct__buf = make(Connect, __vstruct__size)
	__vstruct__buf = Serialize_Connect(__vstruct__buf, Host)
	return __vstruct__buf
}

func Serialize_ConnectResponse(dst ConnectResponse, Success bool) ConnectResponse {
	_ = dst[0]
	dst[0] = *(*byte)(unsafe.Pointer(&Success))

	return dst
}

func New_ConnectResponse(Success bool) ConnectResponse {
	var __vstruct__size = 1
	var __vstruct__buf = make(ConnectResponse, __vstruct__size)
	__vstruct__buf = Serialize_ConnectResponse(__vstruct__buf, Success)
	return __vstruct__buf
}

func Serialize_Stream(dst Stream, Error bool, Data []byte) Stream {
	_ = dst[8]
	dst[0] = *(*byte)(unsafe.Pointer(&Error))

	var __index = uint64(9)
	__tmp_1 := uint64(len(Data)) + __index
	dst[1] = byte(__tmp_1)
	dst[2] = byte(__tmp_1 >> 8)
	dst[3] = byte(__tmp_1 >> 16)
	dst[4] = byte(__tmp_1 >> 24)
	dst[5] = byte(__tmp_1 >> 32)
	dst[6] = byte(__tmp_1 >> 40)
	dst[7] = byte(__tmp_1 >> 48)
	dst[8] = byte(__tmp_1 >> 56)
	copy(dst[__index:__tmp_1], Data)
	return dst
}

func New_Stream(Error bool, Data []byte) Stream {
	var __vstruct__size = 9 + len(Data)
	var __vstruct__buf = make(Stream, __vstruct__size)
	__vstruct__buf = Serialize_Stream(__vstruct__buf, Error, Data)
	return __vstruct__buf
}
