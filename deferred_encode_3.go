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

// TODO: this is incomplete
func (enc *DeferredEncoder3) encodeDeferredComposite(deferredValue *DeferredCompositeValue_V3) {
	enc.encodeTypeTag(TagComposite)

	rw := NewArbitraryReaderWriter(deferredValue.content)

	if deferredValue.locationLoaded {
	} else {
		enc.w.WriteInt(enc.w.writeIndex)                                       // write the new index
		length := rw.ReadIntAt(deferredValue.locationIndex)                    // get content length
		enc.encodeBytes(rw.ReadBytesAt(deferredValue.locationIndex+8, length)) // copy the bytes
	}

	if deferredValue.typeNameLoaded {

	} else {
		enc.w.WriteInt(enc.w.writeIndex)                                       // write the new index
		length := rw.ReadIntAt(deferredValue.locationIndex)                    // get content length
		enc.encodeBytes(rw.ReadBytesAt(deferredValue.locationIndex+8, length)) // copy the bytes
	}

	if deferredValue.kindLoaded {
	} else {
		enc.w.WriteInt(enc.w.writeIndex)                                        // write the new index
		enc.encodeBytes(rw.ReadBytesAt(deferredValue.locationIndex, intLength)) // copy the bytes
	}

	for i, fieldLoaded := range deferredValue.fieldsLoaded {
		if fieldLoaded {
		} else {
			enc.w.WriteInt(enc.w.writeIndex)
			oldIndex := deferredValue.fieldIndices[i]
			length := rw.ReadIntAt(oldIndex)                    // get content length
			enc.encodeBytes(rw.ReadBytesAt(oldIndex+8, length)) // copy the bytes
		}
	}

	return
}

func (enc *DeferredEncoder3) encodeCompositeContent(value *CompositeValue) {
	enc.encodeTypeTag(TagComposite)

	metaEncoder := NewDeferredEncoder_V3(NewDefaultReaderWriter())

	indexesLen := len(value.fields) + 4

	// +1 for type tag
	// +intLength for content length
	offset := indexesLen*intLength + 1 + intLength

	contentWriter := NewDefaultReaderWriter()
	contentEnc := NewDeferredEncoder_V3(contentWriter)

	metaEncoder.w.WriteInt(contentWriter.writeIndex + offset)
	contentEnc.encodeString(value.location)

	metaEncoder.w.WriteInt(contentWriter.writeIndex + offset)
	contentEnc.encodeString(value.typeName)

	metaEncoder.w.WriteInt(contentWriter.writeIndex + offset)
	contentEnc.encodeInt(value.kind)

	metaEncoder.w.WriteInt(len(value.fields))
	for _, field := range value.fields {
		valueIndex := contentWriter.writeIndex + offset
		metaEncoder.w.WriteInt(valueIndex)
		contentEnc.encodeValue(field)
	}

	enc.w.WriteInt(len(metaEncoder.w.bytes) + len(contentWriter.bytes))
	enc.encodeBytes(metaEncoder.w.bytes)
	enc.encodeBytes(contentWriter.bytes)
}

func (enc *DeferredEncoder3) encodeArray(array []interface{}) {
	enc.encodeTypeTag(TagArray)

	metaWriter := NewDefaultReaderWriter()
	metaEncoder := NewDeferredEncoder_V3(metaWriter)

	indicesLen := len(array) + 1

	// +1 for type tag
	// +intLength for content length
	offset := indicesLen*intLength + 1 + intLength

	elementWriter := NewDefaultReaderWriter()
	elementEnc := NewDeferredEncoder_V3(elementWriter)

	metaEncoder.w.WriteInt(len(array))

	for _, element := range array {
		valueIndex := elementWriter.writeIndex + offset
		metaEncoder.w.WriteInt(valueIndex)
		elementEnc.encodeValue(element)
	}

	enc.w.WriteInt(len(metaWriter.bytes) + len(elementWriter.bytes))
	enc.encodeBytes(metaWriter.bytes)
	enc.encodeBytes(elementWriter.bytes)
}

// TODO:
func (enc *DeferredEncoder3) encodeDeferredArray(deferredValue *DeferredArrayValue_V3) {
	// If the value is not built, then dump the content as is.
	if deferredValue.content != nil {
		enc.w.WriteInt(len(deferredValue.content))
		enc.encodeBytes(deferredValue.content)
	}

	return
}

func (enc *DeferredEncoder3) encodeValue(value interface{}) {
	switch val := value.(type) {
	case *CompositeValue:
		enc.encodeComposite(val)
	case string:
		enc.encodeString(val)
	case int:
		enc.encodeInt(val)
	case []interface{}:
		enc.encodeArray(val)

	case *DeferredCompositeValue_V3:
		enc.encodeDeferredComposite(val)

	case *DeferredArrayValue_V3:
		enc.encodeTypeTag(TagArray)
		enc.encodeDeferredArray(val)

	default:
		panic(fmt.Errorf("unknown type: %s", val))
	}
}

func (enc *DeferredEncoder3) encodeTypeTag(b byte) {
	enc.w.WriteByte(b)
}

func (enc *DeferredEncoder3) encodeString(value string) {
	enc.encodeTypeTag(TagString)
	enc.w.WriteString(value)
}

func (enc *DeferredEncoder3) encodeInt(value int) {
	enc.encodeTypeTag(TagInt)
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
