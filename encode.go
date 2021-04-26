package encoder

import (
	"fmt"
)

const TagComposite byte = 1
const TagString = 2
const TagInt = 3
const TagArray = 4

const encodedCompositeValueLength = 4

type Encoder struct {
	w *ByteReaderWriter
}

func NewEncoder(w *ByteReaderWriter) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (enc *Encoder) Encode(value interface{}) {
	enc.encodeValue(value)
}

func (enc *Encoder) encodeArray(array []interface{}) {
	enc.w.WriteByte(TagArray)
	enc.encodeInt(len(array))

	for _, element := range array {
		enc.encodeValue(element)
	}
}

func (enc *Encoder) encodeComposite(value *CompositeValue) {
	enc.w.WriteByte(TagComposite)
	enc.encodeInt(encodedCompositeValueLength)

	enc.encodeString(value.location)
	enc.encodeString(value.typeName)
	enc.encodeInt(value.kind)
	enc.encodeValue(value.fields)
}

func (enc *Encoder) encodeValue(value interface{}) {
	switch val := value.(type) {
	case *CompositeValue:
		enc.encodeComposite(val)
	case string:
		enc.encodeString(val)
	case int:
		enc.encodeInt(val)
	case []interface{}:
		enc.encodeArray(val)
	default:
		panic(fmt.Errorf("unknown type: %s", val))
	}
}

func (enc *Encoder) encodeString(value string) {
	enc.w.WriteByte(TagString)
	enc.w.WriteString(value)
}

func (enc *Encoder) encodeInt(value int) {
	enc.w.WriteByte(TagInt)
	enc.w.WriteInt(value)
}

func (enc *Encoder) reset() {
	enc.w = NewDefaultReaderWriter()
}
