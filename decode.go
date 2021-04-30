package encoder

type Decoder struct {
	w *ByteReaderWriter
}

func NewDecoder(w *ByteReaderWriter) *Decoder {
	return &Decoder{
		w: w,
	}
}

func (dec *Decoder) Decode() interface{} {
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

func (dec *Decoder) decodeString() string {
	return dec.w.ReadString()
}

func (dec *Decoder) decodeInt() int {
	return dec.w.ReadInt()
}

func (dec *Decoder) decodeComposite() *CompositeValue {
	dec.decodeInt() // ignore fields length

	location := dec.decodeString()
	name := dec.decodeString()
	kind := dec.decodeInt()

	fields := dec.decodeArray()

	return &CompositeValue{
		location: location,
		typeName: name,
		kind:     kind,
		fields:   fields,
	}
}

func (dec *Decoder) decodeArray() []interface{} {
	len := dec.decodeInt()
	values := make([]interface{}, len)

	for i := 0; i < len; i++ {
		values[i] = dec.Decode()
	}

	return values
}

func (dec *Decoder) readByte() byte {
	return dec.w.ReadByte()
}

func (dec *Decoder) readBytes(n int) []byte {
	return dec.w.ReadBytes(n)
}

func (dec *Decoder) reset() {
	dec.w.ReadReset()
}
