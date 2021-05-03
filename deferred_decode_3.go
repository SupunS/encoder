package encoder

type DeferredDecoder3 struct {
	w *ByteReaderWriter
}

func NewDeferredDecoder3(w *ByteReaderWriter) *DeferredDecoder3 {
	return &DeferredDecoder3{
		w: w,
	}
}

func (dec *DeferredDecoder3) Decode() interface{} {
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

func (dec *DeferredDecoder3) decodeString() string {
	return dec.w.ReadString()
}

func (dec *DeferredDecoder3) decodeInt() int {
	return dec.w.ReadInt()
}

func (dec *DeferredDecoder3) decodeComposite() *DeferredCompositeValue_V3 {
	locationIndex := dec.decodeInt()
	typeNameIndex := dec.decodeInt()
	kindIndex := dec.decodeInt()

	fieldsLen := dec.decodeInt()
	fieldIndices := make([]int, fieldsLen)

	for i := 0; i < fieldsLen; i++ {
		fieldIndices[i] = dec.decodeInt()
	}

	fieldLoaded := make([]bool, fieldsLen)
	fields := make([]interface{}, fieldsLen)

	return &DeferredCompositeValue_V3{
		content:       dec.w.bytes,
		locationIndex: locationIndex,
		typeNameIndex: typeNameIndex,
		kindIndex:     kindIndex,
		fieldIndices:  fieldIndices,
		fieldsLoaded:  fieldLoaded,
		fields:        fields,
	}
}

func (dec *DeferredDecoder3) decodeArray() *DeferredArrayValue_V3 {
	arrayLen := dec.decodeInt()
	elementIndices := make([]int, arrayLen)

	for i := 0; i < arrayLen; i++ {
		elementIndices[i] = dec.decodeInt()
	}

	elementLoaded := make([]bool, arrayLen)
	elements := make([]interface{}, arrayLen)

	return &DeferredArrayValue_V3{
		content:        dec.w.bytes,
		elementIndices: elementIndices,
		elementLoaded:  elementLoaded,
		elements:       elements,
	}
}

func (dec *DeferredDecoder3) readByte() byte {
	return dec.w.ReadByte()
}

func (dec *DeferredDecoder3) readBytes(n int) []byte {
	return dec.w.ReadBytes(n)
}

func (dec *DeferredDecoder3) reset() {
	dec.w.ReadReset()
}
