package encoder

import (
	"fmt"
)

type DeferredEncoder3 struct {
	w *ByteReaderWriter
}

func NewDeferredEncoder_V3(w *ByteReaderWriter) *DeferredEncoder3 {
	return &DeferredEncoder3{
		w: w,
	}
}

func (enc *DeferredEncoder3) Encode(value interface{}) {
	enc.encodeValue(value)
}

func (enc *DeferredEncoder3) encodeComposite(value *CompositeValue) {
	enc.encodeCompositeContent(value)
}

func (enc *DeferredEncoder3) encodeDeferredComposite(deferredValue *DeferredCompositeValue_V2) {
	// If the value is not built, then dump the content as is.
	if deferredValue.metaContent != nil {
		enc.encodeBytes(deferredValue.metaContent)
	}

	if deferredValue.fieldsContent != nil {
		enc.encodeBytes(deferredValue.fieldsContent)
	}

	return
}

func (enc *DeferredEncoder3) encodeCompositeContent(value *CompositeValue) {
	enc.encodeTypeTag(TagComposite)

	indicesLen := len(value.fields) + 4
	offset := indicesLen*intLength + 1 // +1 for type tag

	contentWriter := NewDefaultReaderWriter()
	contentEnc := NewDeferredEncoder_V3(contentWriter)

	enc.encodeInt(contentWriter.writeIndex + offset)
	contentEnc.encodeString(value.location)

	enc.encodeInt(contentWriter.writeIndex + offset)
	contentEnc.encodeString(value.typeName)

	enc.encodeInt(contentWriter.writeIndex + offset)
	contentEnc.encodeInt(value.kind)

	enc.encodeInt(len(value.fields))
	for _, field := range value.fields {
		valueIndex := contentWriter.writeIndex + offset
		enc.encodeInt(valueIndex)
		contentEnc.encodeValue(field)
	}

	enc.encodeBytes(contentWriter.bytes)
}

func (enc *DeferredEncoder3) encodeArray(array []interface{}) {
	enc.encodeTypeTag(TagArray)

	indicesLen := len(array) + 1
	offset := indicesLen*intLength + 1 // +1 for type tag

	contentWriter := NewDefaultReaderWriter()
	contentEnc := NewDeferredEncoder_V3(contentWriter)

	enc.encodeInt(len(array))

	for _, element := range array {
		valueIndex := contentWriter.writeIndex + offset
		enc.encodeInt(valueIndex)
		contentEnc.encodeValue(element)
	}

	enc.encodeBytes(contentWriter.bytes)
}

func (enc *DeferredEncoder3) encodeValue(value interface{}) {
	switch val := value.(type) {
	case *CompositeValue:
		enc.encodeComposite(val)
	case string:
		enc.encodeTypeTag(TagString)
		enc.encodeString(val)
	case int:
		enc.encodeTypeTag(TagInt)
		enc.encodeInt(val)
	case []interface{}:
		enc.encodeArray(val)

	case *DeferredCompositeValue_V2:
		enc.encodeTypeTag(TagComposite)
		enc.encodeDeferredComposite(val)

	default:
		panic(fmt.Errorf("unknown type: %s", val))
	}
}

func (enc *DeferredEncoder3) encodeTypeTag(b byte) {
	enc.w.WriteByte(b)
}

func (enc *DeferredEncoder3) encodeString(value string) {
	enc.w.WriteString(value)
}

func (enc *DeferredEncoder3) encodeInt(value int) {
	enc.w.WriteInt(value)
}

func (enc *DeferredEncoder3) encodeIntAt(value int, index int) {
	enc.w.WriteIntAt(value, index)
}

func (enc *DeferredEncoder3) encodeBytes(content []byte) {
	enc.w.WriteBytes(content)
}

func (enc *DeferredEncoder3) reset() {
	enc.w = NewDefaultReaderWriter()
}
