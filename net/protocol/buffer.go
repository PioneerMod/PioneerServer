package protocol

import (
	"fmt"
	"math"
)

type PacketBuffer struct {
	Bytes []byte
}

func CreatePacketBuffer(id int) PacketBuffer {
	buffer := PacketBuffer{}
	buffer.WriteInt32(int32(id))
	return buffer
}

func ParsePacketBuffer(bytes []byte) (int, PacketBuffer) {
	buffer := PacketBuffer{bytes}
	return int(buffer.ReadInt32()), buffer
}

func (buffer *PacketBuffer) WriteInt16(value int16) {
	buffer.Bytes = append(buffer.Bytes, byte(value>>8), byte(value))
}

func (buffer *PacketBuffer) WriteInt32(value int32) {
	buffer.Bytes = append(buffer.Bytes, byte(value>>24), byte(value>>16), byte(value>>8), byte(value))
}

func (buffer *PacketBuffer) WriteInt64(value int64) {
	buffer.Bytes = append(buffer.Bytes, byte(value>>56), byte(value>>48), byte(value>>40), byte(value>>32), byte(value>>24), byte(value>>16), byte(value>>8), byte(value))
}

func (buffer *PacketBuffer) WriteDouble(value float64) {
	var buf [8]byte
	n := math.Float64bits(value)
	buf[0] = byte(n >> 56)
	buf[1] = byte(n >> 48)
	buf[2] = byte(n >> 40)
	buf[3] = byte(n >> 32)
	buf[4] = byte(n >> 24)
	buf[5] = byte(n >> 16)
	buf[6] = byte(n >> 8)
	buf[7] = byte(n)
}

func (buffer *PacketBuffer) WriteString(value string) {
	buffer.WriteInt16(int16(len(value)))
	buffer.Bytes = append(buffer.Bytes, value...)
}

func (buffer *PacketBuffer) WriteStringUTF16(value string) {
	buffer.WriteInt16(int16(len(value)))
	for _, char := range value {
		buffer.WriteInt16(int16(char))
	}
}

func (buffer *PacketBuffer) WriteBool(value bool) {
	if value {
		buffer.Bytes = append(buffer.Bytes, 1)
	} else {
		buffer.Bytes = append(buffer.Bytes, 0)
	}
}

func (buffer *PacketBuffer) WriteBytes(value []byte) {
	buffer.WriteInt32(int32(len(value)))
	buffer.Bytes = append(buffer.Bytes, value...)
}

func (buffer *PacketBuffer) WriteInt32Array(value []int) {
	buffer.WriteInt32(int32(len(value)))
	for _, val := range value {
		buffer.WriteInt32(int32(val))
	}
}

func (buffer *PacketBuffer) WriteInt64Array(value []int64) {
	buffer.WriteInt32(int32(len(value)))
	for _, val := range value {
		buffer.WriteInt64(val)
	}
}

func (buffer *PacketBuffer) WriteStringArray(value []string) {
	buffer.WriteInt32(int32(len(value)))
	for _, val := range value {
		buffer.WriteString(val)
	}
}

func (buffer *PacketBuffer) WriteBoolArray(value []bool) {
	buffer.WriteInt32(int32(len(value)))
	for _, val := range value {
		buffer.WriteBool(val)
	}
}

func (buffer *PacketBuffer) WriteDoubleArray(value []float64) {
	buffer.WriteInt32(int32(len(value)))
	for _, val := range value {
		buffer.WriteDouble(val)
	}
}

func (buffer *PacketBuffer) ReadInt16() int16 {
	return int16(buffer.Bytes[0])<<8 | int16(buffer.Bytes[1])
}

func (buffer *PacketBuffer) ReadInt32() int32 {
	return int32(buffer.Bytes[0])<<24 | int32(buffer.Bytes[1])<<16 | int32(buffer.Bytes[2])<<8 | int32(buffer.Bytes[3])
}

func (buffer *PacketBuffer) ReadInt64() int64 {
	return int64(buffer.Bytes[0])<<56 | int64(buffer.Bytes[1])<<48 | int64(buffer.Bytes[2])<<40 | int64(buffer.Bytes[3])<<32 | int64(buffer.Bytes[4])<<24 | int64(buffer.Bytes[5])<<16 | int64(buffer.Bytes[6])<<8 | int64(buffer.Bytes[7])
}

func (buffer *PacketBuffer) ReadString() string {
	length := buffer.ReadInt16()
	return string(buffer.Bytes[:length])
}

func (buffer *PacketBuffer) ReadStringUTF16() string {
	length := buffer.ReadInt16()
	var result string
	for i := 0; i < int(length); i++ {
		result += fmt.Sprintf("%c", buffer.ReadInt16())
	}
	return result
}

func (buffer *PacketBuffer) ReadBool() bool {
	return buffer.Bytes[0] == 1
}

func (buffer *PacketBuffer) ReadBytes() []byte {
	length := buffer.ReadInt32()
	return buffer.Bytes[:length]
}

func (buffer *PacketBuffer) ReadDouble() float64 {
	bits := uint64(buffer.Bytes[0])<<56 | uint64(buffer.Bytes[1])<<48 | uint64(buffer.Bytes[2])<<40 | uint64(buffer.Bytes[3])<<32 | uint64(buffer.Bytes[4])<<24 | uint64(buffer.Bytes[5])<<16 | uint64(buffer.Bytes[6])<<8 | uint64(buffer.Bytes[7])
	return math.Float64frombits(bits)
}

func (buffer *PacketBuffer) ReadInt16Array() []int16 {
	length := buffer.ReadInt32()
	var result []int16
	for i := 0; i < int(length); i++ {
		result = append(result, buffer.ReadInt16())
	}
	return result
}

func (buffer *PacketBuffer) ReadInt32Array() []int32 {
	length := buffer.ReadInt32()
	var result []int32
	for i := 0; i < int(length); i++ {
		result = append(result, buffer.ReadInt32())
	}
	return result
}

func (buffer *PacketBuffer) ReadInt64Array() []int64 {
	length := buffer.ReadInt32()
	var result []int64
	for i := 0; i < int(length); i++ {
		result = append(result, buffer.ReadInt64())
	}
	return result
}

func (buffer *PacketBuffer) ReadStringArray() []string {
	length := buffer.ReadInt32()
	var result []string
	for i := 0; i < int(length); i++ {
		result = append(result, buffer.ReadString())
	}
	return result
}

func (buffer *PacketBuffer) ReadBoolArray() []bool {
	length := buffer.ReadInt32()
	var result []bool
	for i := 0; i < int(length); i++ {
		result = append(result, buffer.ReadBool())
	}
	return result
}

func (buffer *PacketBuffer) ReadDoubleArray() []float64 {
	length := buffer.ReadInt32()
	var result []float64
	for i := 0; i < int(length); i++ {
		result = append(result, buffer.ReadDouble())
	}
	return result
}

func (buffer *PacketBuffer) Finalize() {
	buffer.Bytes = append(buffer.Bytes, '\n')
}
