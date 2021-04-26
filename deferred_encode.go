package encoder

import (
	"fmt"
)

type DeferredEncoder struct {
	w *ByteReaderWriter
}

func NewDeferredEncoder(w *ByteReaderWriter) *DeferredEncoder {
	return &DeferredEncoder{
		w: w,
	}
}

func (enc *DeferredEncoder) Encode(value interface{}) {
	enc.encodeValue(value)
}

func (enc *DeferredEncoder) encodeArray(array []interface{}) {
	enc.w.WriteByte(TagArray)
	enc.encodeInt(len(array))

	for _, element := range array {
		enc.encodeValue(element)
	}
}

func (enc *DeferredEncoder) encodeComposite(value *CompositeValue) {
	enc.w.WriteByte(TagComposite)

	enc.encodeCompositeContent(value)
}

func (enc *DeferredEncoder) encodeDeferredComposite(deferredValue *DeferredCompositeValue) {
	enc.w.WriteByte(TagComposite)

	// If the value is not built, then dump the content as is.
	if deferredValue.content != nil {
		enc.encodeBytes(deferredValue.content)
		return
	}

	value := deferredValue.value
	enc.encodeCompositeContent(value)
}

func (enc *DeferredEncoder) encodeCompositeContent(value *CompositeValue) {
	w := NewDefaultReaderWriter()
	subEncoder := NewDeferredEncoder(w)


	subEncoder.encodeString(value.location)
	subEncoder.encodeString(value.typeName)
	subEncoder.encodeInt(value.kind)
	subEncoder.encodeValue(value.fields)

	enc.encodeInt(len(w.bytes))
	enc.encodeBytes(w.bytes)
}

func (enc *DeferredEncoder) encodeValue(value interface{}) {
	switch val := value.(type) {
	case *CompositeValue:
		enc.encodeComposite(val)
	case string:
		enc.encodeString(val)
	case int:
		enc.encodeInt(val)
	case []interface{}:
		enc.encodeArray(val)

	case *DeferredCompositeValue:
		enc.encodeDeferredComposite(val)

	default:
		panic(fmt.Errorf("unknown type: %s", val))
	}
}

func (enc *DeferredEncoder) encodeString(value string) {
	enc.w.WriteByte(TagString)
	enc.w.WriteString(value)
}

func (enc *DeferredEncoder) encodeInt(value int) {
	enc.w.WriteByte(TagInt)
	enc.w.WriteInt(value)
}

func (enc *DeferredEncoder) encodeBytes(content [][]byte) {
	enc.w.WriteBytes(content)
}

func (enc *DeferredEncoder) reset() {
	enc.w = NewDefaultReaderWriter()
}
