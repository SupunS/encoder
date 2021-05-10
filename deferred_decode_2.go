package encoder

type DeferredDecoder2 struct {
	w *ByteReaderWriter
}

func NewDeferredDecoder2(w *ByteReaderWriter) *DeferredDecoder2 {
	return &DeferredDecoder2{
		w: w,
	}
}

func (dec *DeferredDecoder2) Decode() interface{} {
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

func (dec *DeferredDecoder2) decodeString() string {
	return dec.w.ReadString()
}

func (dec *DeferredDecoder2) decodeInt() int {
	return dec.w.ReadInt()
}

func (dec *DeferredDecoder2) decodeComposite() *DeferredCompositeValue_V2 {
	metaLen := dec.decodeInt() // length of content
	metaContent := dec.readBytes(metaLen)

	fieldsLen := dec.decodeInt() // length of content
	fieldsContent := dec.readBytes(fieldsLen)

	return &DeferredCompositeValue_V2{
		metaContent:   metaContent,
		fieldsContent: fieldsContent,
	}
}

func (dec *DeferredDecoder2) decodeArray() *DeferredArrayValue_V2 {
	contentLen := dec.decodeInt() // length of content
	content := dec.readBytes(contentLen)

	return &DeferredArrayValue_V2{
		content: content,
	}
}

func (dec *DeferredDecoder2) readByte() byte {
	return dec.w.ReadByte()
}

func (dec *DeferredDecoder2) readBytes(n int) []byte {
	return dec.w.ReadBytes(n)
}

func (dec *DeferredDecoder2) reset() {
	dec.w.ReadReset()
}
