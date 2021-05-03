package encoder

type DeferredValue interface {
	ensureLoaded()
}

type CadenceCompositeValue interface {
	member(int) interface{}
}

var _ CadenceCompositeValue = &CompositeValue{}
var _ CadenceCompositeValue = &DeferredCompositeValue{}

type CompositeValue struct {
	location string
	typeName string
	kind     int
	fields   []interface{}
}

func (c *CompositeValue) member(i int) interface{} {
	return c.fields[i]
}

type DeferredCompositeValue struct {
	content []byte
	value   *CompositeValue
}

func (val *DeferredCompositeValue) member(i int) interface{} {
	val.ensureLoaded() // Make sure the content is built before doing any operation
	return val.value.member(i)
}

// Perform a shallow-build of the content.
func (val *DeferredCompositeValue) ensureLoaded() {
	if val.value != nil {
		return
	}

	rw := NewReaderWriter(val.content)
	decoder := NewDeferredDecoder(rw)

	location := decoder.decodeString()
	name := decoder.decodeString()
	kind := decoder.decodeInt()
	fields := decoder.decodeArray()

	val.value = &CompositeValue{
		location: location,
		typeName: name,
		kind:     kind,
		fields:   fields,
	}

	// clear the content
	val.content = nil
}

type DeferredCompositeValue_V2 struct {
	metaContent   []byte
	fieldsContent []byte
	value         *CompositeValue
}

func (val *DeferredCompositeValue_V2) member(i int) interface{} {
	val.ensureLoaded() // Make sure the content is built before doing any operation
	return val.value.member(i)
}

// Perform a shallow-build of the content.
func (val *DeferredCompositeValue_V2) ensureLoaded() {
	if val.value != nil {
		return
	}

	rw := NewReaderWriter(val.metaContent)
	decoder := NewDeferredDecoder(rw)

	location := decoder.decodeString()
	name := decoder.decodeString()
	kind := decoder.decodeInt()

	rw = NewReaderWriter(val.fieldsContent)
	decoder = NewDeferredDecoder(rw)
	fields := decoder.decodeArray()

	val.value = &CompositeValue{
		location: location,
		typeName: name,
		kind:     kind,
		fields:   fields,
	}

	// clear the content
	val.metaContent = nil
	val.fieldsContent = nil
}

type DeferredCompositeValue_V3 struct {
	content []byte

	location       string
	locationIndex  int
	locationLoaded bool

	typeName       string
	typeNameIndex  int
	typeNameLoaded bool

	kind       int
	kindIndex  int
	kindLoaded bool

	fieldIndices []int
	fieldsLoaded []bool
	fields       []interface{}
}

func (val *DeferredCompositeValue_V3) member(i int) interface{} {
	if val.fieldsLoaded[i] {
		return val.fields[i]
	}

	fieldIndex := val.fieldIndices[i]
	rw := NewReaderWriter(val.content[fieldIndex:])
	decoder := NewDeferredDecoder3(rw)
	element := decoder.Decode()

	val.fields[i] = element
	return element
}

func (val *DeferredCompositeValue_V3) Name() string {
	if val.typeNameLoaded {
		return val.typeName
	}

	rw := NewReaderWriter(val.content[val.typeNameIndex:])
	decoder := NewDeferredDecoder3(rw)
	val.typeName = decoder.decodeString()
	return val.typeName
}

// Perform a shallow-build of the content.
func (val *DeferredCompositeValue_V3) ensureLoaded() {
	// clear the content
	val.content = nil
}

type DeferredArrayValue_V3 struct {
	content []byte

	elementIndices []int
	elementLoaded  []bool
	elements       []interface{}
}

func (val *DeferredArrayValue_V3) member(i int) interface{} {
	if val.elementLoaded[i] {
		return val.elements[i]
	}

	elementIndex := val.elementIndices[i]
	rw := NewReaderWriter(val.content[elementIndex:])
	decoder := NewDeferredDecoder3(rw)
	element := decoder.Decode()

	val.elements[i] = element
	return element
}
