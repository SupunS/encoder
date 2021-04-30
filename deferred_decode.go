package encoder

type DeferredDecoder struct {
	w *ByteReaderWriter
}

func NewDeferredDecoder(w *ByteReaderWriter) *DeferredDecoder {
	return &DeferredDecoder{
		w: w,
	}
}

func (dec *DeferredDecoder) Decode() interface{} {
	b := dec.readByte()

	switch b {
	case TagComposite:
		return dec.decodeComposite()
	case TagArray:
		return dec.decodeArray()
	case TagString:
		return dec.decodeString()
	case TagInt:
		return dec.decodeInt()
	}

	return nil
}

func (dec *DeferredDecoder) decodeString() string {
	return dec.w.ReadString()
}

func (dec *DeferredDecoder) decodeInt() int {
	return dec.w.ReadInt()
}

func (dec *DeferredDecoder) decodeComposite() *DeferredCompositeValue {
	len := dec.decodeInt() // length of content
	content := dec.readBytes(len)

	return &DeferredCompositeValue{
		content: content,
		value:   nil,
	}
}

func (dec *DeferredDecoder) decodeArray() []interface{} {
	len := dec.decodeInt()
	values := make([]interface{}, len)

	for i := 0; i < len; i++ {
		values[i] = dec.Decode()
	}

	return values
}

func (dec *DeferredDecoder) readByte() byte {
	return dec.w.ReadByte()
}

func (dec *DeferredDecoder) readBytes(n int) []byte {
	return dec.w.ReadBytes(n)
}

func (dec *DeferredDecoder) reset() {
	dec.w.ReadReset()
}
