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
	lengthStartIndex := enc.writeDummyLength()

	enc.w.WriteInt(len(array))

	for _, element := range array {
		enc.encodeValue(element)
	}

	// Update the size with the actual content size
	contentLength := enc.w.writeIndex - lengthStartIndex - intLength
	enc.encodeIntAt(contentLength, lengthStartIndex)
}

func (enc *DeferredEncoder2) encodeDeferredArray(deferredValue *DeferredArrayValue_V2) {
	// If the value is not built, then dump the content as is.
	if deferredValue.content != nil {
		enc.w.WriteInt(len(deferredValue.content))
		enc.encodeBytes(deferredValue.content)
	}

	return
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
	// Reserve 8-bits for length
	lengthStartIndex := enc.writeDummyLength()

	// Write meta content
	enc.encodeString(value.location)
	enc.encodeString(value.typeName)
	enc.encodeInt(value.kind)

	// Update the size with the actual content size
	contentLength := enc.w.writeIndex - lengthStartIndex - intLength
	enc.encodeIntAt(contentLength, lengthStartIndex)

	// Do the same for fields as well

	lengthStartIndex = enc.writeDummyLength()
	enc.encodeArray(value.fields)
	contentLength = enc.w.writeIndex - lengthStartIndex - intLength
	enc.encodeIntAt(contentLength, lengthStartIndex)
}

const intLength = 8

func (enc *DeferredEncoder2) writeDummyLength() int {
	lengthStartIndex := enc.w.writeIndex
	lengthBytes := make([]byte, intLength)
	enc.encodeBytes(lengthBytes)
	return lengthStartIndex
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
	case *DeferredArrayValue_V2:
		enc.encodeDeferredArray(val)

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

func (enc *DeferredEncoder2) encodeIntAt(value int, index int) {
	enc.w.WriteIntAt(value, index)
}

func (enc *DeferredEncoder2) encodeBytes(content []byte) {
	enc.w.WriteBytes(content)
}

func (enc *DeferredEncoder2) reset() {
	enc.w = NewDefaultReaderWriter()
}
