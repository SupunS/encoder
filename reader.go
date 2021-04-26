package encoder

import "strconv"

type ByteReaderWriter struct {
	bytes     [][]byte
	readIndex int
}

func NewDefaultReaderWriter() *ByteReaderWriter {
	return &ByteReaderWriter{}
}

func NewReaderWriter(bytes [][]byte) *ByteReaderWriter {
	return &ByteReaderWriter{
		bytes: bytes,
	}
}

func (rw *ByteReaderWriter) WriteString(s string) {
	rw.bytes = append(rw.bytes, []byte(s))
}

func (rw *ByteReaderWriter) WriteInt(i int) {
	bytes := strconv.AppendInt(nil, int64(i), 10)
	rw.bytes = append(rw.bytes, bytes)
}

func (rw *ByteReaderWriter) WriteByte(b byte) {
	rw.bytes = append(rw.bytes, []byte{b})
}

func (rw *ByteReaderWriter) WriteBytes(b [][]byte) {
	rw.bytes = append(rw.bytes, b...)
}

func (rw *ByteReaderWriter) ReadString() string {
	s := string(rw.bytes[rw.readIndex])
	rw.readIndex++
	return s
}

func (rw *ByteReaderWriter) ReadInt() int {
	s := string(rw.bytes[rw.readIndex])
	value, _ := strconv.ParseInt(s, 10, 64)
	rw.readIndex++
	return int(value)
}

func (rw *ByteReaderWriter) ReadByte() byte {
	b := rw.bytes[rw.readIndex][0]
	rw.readIndex++
	return b
}

func (rw *ByteReaderWriter) Peek() byte {
	return rw.bytes[rw.readIndex][0]
}

func (rw *ByteReaderWriter) ReadBytes(n int) [][]byte {
	endIndex := rw.readIndex + n

	b := rw.bytes[rw.readIndex:endIndex]

	rw.readIndex = endIndex
	return b
}

func (rw *ByteReaderWriter) ReadReset() {
	rw.readIndex = 0
}

func (rw *ByteReaderWriter) WriteReset() {
	rw.bytes = nil
}
