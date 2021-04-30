package encoder

import (
	"fmt"
)

type DeferredEncoder2 struct {
	w *ByteReaderWriter
}

func NewDeferredEncoder_V2(w *ByteReaderWriter) *DeferredEncoder2 {
	return &DeferredEncoder2{
		w: w,
	}
}

func (enc *DeferredEncoder2) Encode(value interface{}) {
	enc.encodeValue(value)
}

func (enc *DeferredEncoder2) encodeArray(array []interface{}) {
	enc.w.WriteInt(len(array))

	for _, element := range array {
		enc.encodeValue(element)
	}
}

func (enc *DeferredEncoder2) encodeComposite(value *CompositeValue) {
	enc.encodeCompositeContent(value)
}

func (enc *DeferredEncoder2) encodeDeferredComposite(deferredValue *DeferredCompositeValue_V2) {
	// If the value is not built, then dump the content as is.
	if deferredValue.metaContent != nil {
		enc.encodeBytes(deferredValue.metaContent)
	}

	if deferredValue.fieldsContent != nil {
		enc.encodeBytes(deferredValue.fieldsContent)
	}

	return
}

func (enc *DeferredEncoder2) encodeCompositeContent(value *CompositeValue) {
	w := NewDefaultReaderWriter()
	subEncoder := NewDeferredEncoder_V2(w)

	subEncoder.encodeString(value.location)
	subEncoder.encodeString(value.typeName)
	subEncoder.encodeInt(value.kind)

	enc.encodeInt(len(w.bytes))
	enc.encodeBytes(w.bytes)

	subEncoder.reset()
	subEncoder.encodeArray(value.fields)

	enc.encodeInt(len(w.bytes))
	enc.encodeBytes(w.bytes)
}

func (enc *DeferredEncoder2) encodeValue(value interface{}) {
	switch val := value.(type) {
	case *CompositeValue:
		enc.w.WriteByte(TagComposite)
		enc.encodeComposite(val)
	case string:
		enc.w.WriteByte(TagString)
		enc.encodeString(val)
	case int:
		enc.w.WriteByte(TagInt)
		enc.encodeInt(val)
	case []interface{}:
		enc.w.WriteByte(TagArray)
		enc.encodeArray(val)

	case *DeferredCompositeValue_V2:
		enc.encodeDeferredComposite(val)

	default:
		panic(fmt.Errorf("unknown type: %s", val))
	}
}

func (enc *DeferredEncoder2) encodeString(value string) {
	enc.w.WriteString(value)
}

func (enc *DeferredEncoder2) encodeInt(value int) {
	enc.w.WriteInt(value)
}

func (enc *DeferredEncoder2) encodeBytes(content [][]byte) {
	enc.w.WriteBytes(content)
}

func (enc *DeferredEncoder2) reset() {
	enc.w = NewDefaultReaderWriter()
}
