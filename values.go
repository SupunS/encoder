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
	content [][]byte
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
