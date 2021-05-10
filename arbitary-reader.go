package encoder

import (
	"encoding/binary"
)

type ArbitraryByteReaderWriter struct {
	bytes      []byte
	writeIndex int
}

func NewDefaultArbitraryReaderWriter() *ArbitraryByteReaderWriter {
	return &ArbitraryByteReaderWriter{}
}

func NewArbitraryReaderWriter(bytes []byte) *ArbitraryByteReaderWriter {
	return &ArbitraryByteReaderWriter{
		bytes: bytes,
	}
}

// Writer methods

func (rw *ArbitraryByteReaderWriter) WriteString(s string) {
	rw.WriteInt(len(s))
	rw.WriteBytes([]byte(s))
}

func (rw *ArbitraryByteReaderWriter) WriteInt(i int) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(i))
	rw.WriteBytes(bytes)
}

func (rw *ArbitraryByteReaderWriter) WriteIntAt(i, index int) {
	binary.BigEndian.PutUint64(rw.bytes[index:], uint64(i))
}

func (rw *ArbitraryByteReaderWriter) WriteByte(b byte) {
	rw.bytes = append(rw.bytes, b)
	rw.writeIndex++
}

func (rw *ArbitraryByteReaderWriter) WriteBytes(bytes []byte) {
	rw.bytes = append(rw.bytes, bytes...)
	rw.writeIndex += len(bytes)
}

func (rw *ArbitraryByteReaderWriter) ReadStringAt(index int) string {
	length := rw.ReadIntAt(index)
	startIndex := index + intLength
	endIndex := startIndex + length
	s := string(rw.bytes[startIndex:endIndex])
	return s
}

// Reader methods

func (rw *ArbitraryByteReaderWriter) ReadIntAt(index int) int {
	endIndex := index + 8
	value := binary.BigEndian.Uint64(rw.bytes[index:endIndex])
	return int(value)
}

func (rw *ArbitraryByteReaderWriter) ReadByteAt(index int) byte {
	b := rw.bytes[index]
	return b
}

func (rw *ArbitraryByteReaderWriter) ReadBytesAt(index, n int) []byte {
	endIndex := index + n
	b := rw.bytes[index:endIndex]

	return b
}

func (rw *ArbitraryByteReaderWriter) WriteReset() {
	rw.writeIndex = 0
	rw.bytes = nil
}
