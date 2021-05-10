package encoder

import (
	"encoding/binary"
)

type ByteReaderWriter struct {
	bytes      []byte
	readIndex  int
	writeIndex int
}

func NewDefaultReaderWriter() *ByteReaderWriter {
	return &ByteReaderWriter{}
}

func NewReaderWriter(bytes []byte) *ByteReaderWriter {
	return &ByteReaderWriter{
		bytes: bytes,
	}
}

// Writer methods

func (rw *ByteReaderWriter) WriteString(s string) {
	rw.WriteInt(len(s))
	rw.WriteBytes([]byte(s))
}

func (rw *ByteReaderWriter) WriteInt(i int) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(i))
	rw.WriteBytes(bytes)
}

func (rw *ByteReaderWriter) WriteIntAt(i, index int) {
	binary.BigEndian.PutUint64(rw.bytes[index:], uint64(i))
}

func (rw *ByteReaderWriter) WriteByte(b byte) {
	rw.bytes = append(rw.bytes, b)
	rw.writeIndex++
}

func (rw *ByteReaderWriter) WriteBytes(bytes []byte) {
	rw.bytes = append(rw.bytes, bytes...)
	rw.writeIndex += len(bytes)
}

func (rw *ByteReaderWriter) ReadString() string {
	length := rw.ReadInt()
	endIndex := rw.readIndex + length
	s := string(rw.bytes[rw.readIndex:endIndex])
	rw.readIndex = endIndex
	return s
}

func (rw *ByteReaderWriter) ReadStringAt(index int) string {
	length := rw.ReadInt()
	endIndex := rw.readIndex + length
	s := string(rw.bytes[rw.readIndex:endIndex])
	rw.readIndex = endIndex
	return s
}

// Reader methods

func (rw *ByteReaderWriter) ReadInt() int {
	endIndex := rw.readIndex + 8
	value := binary.BigEndian.Uint64(rw.bytes[rw.readIndex:endIndex])
	rw.readIndex = endIndex
	return int(value)
}

func (rw *ByteReaderWriter) ReadByte() byte {
	b := rw.bytes[rw.readIndex]
	rw.readIndex++
	return b
}

func (rw *ByteReaderWriter) ReadBytes(n int) []byte {
	endIndex := rw.readIndex + n
	b := rw.bytes[rw.readIndex:endIndex]

	rw.readIndex = endIndex
	return b
}

func (rw *ByteReaderWriter) ReadReset() {
	rw.readIndex = 0
}

func (rw *ByteReaderWriter) WriteReset() {
	rw.writeIndex = 0
	rw.bytes = nil
}
