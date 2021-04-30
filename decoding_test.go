package encoder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDecoding(b *testing.B) {

	b.Run("decoding", func(b *testing.B) {

		b.Run("normaldec", func(b *testing.B) {
			decoder := getDecoder()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decoder.Decode()
				decoder.reset()
			}
		})

		b.Run("deferred", func(b *testing.B) {
			decoder := getDeferredDecoder()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decoder.Decode()
				decoder.reset()
			}
		})

		b.Run("deferred_v2", func(b *testing.B) {
			decoder := getDeferredDecoder_V2()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decoder.Decode()
				decoder.reset()
			}
		})
	})

	b.Run("encoding", func(b *testing.B) {

		b.Run("normal", func(b *testing.B) {
			encoder := NewEncoder(NewDefaultReaderWriter())

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				encoder.Encode(valueArray)
				encoder.reset()
			}
		})

		b.Run("deferred", func(b *testing.B) {
			encoder := NewDeferredEncoder(NewDefaultReaderWriter())

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				encoder.Encode(valueArray)
				encoder.reset()
			}
		})

		b.Run("deferred_v2", func(b *testing.B) {
			encoder := NewDeferredEncoder_V2(NewDefaultReaderWriter())

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				encoder.Encode(valueArray)
				encoder.reset()
			}
		})

	})

	b.Run("re-encoding", func(b *testing.B) {

		b.Run("normal", func(b *testing.B) {
			decoder := getDecoder()
			decodedValue := decoder.Decode()

			r := NewDefaultReaderWriter()
			encoder := NewEncoder(r)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				encoder.Encode(decodedValue)
				encoder.reset()
			}

		})

		b.Run("deferred", func(b *testing.B) {

			decoder := getDeferredDecoder()
			decodedValue := decoder.Decode()

			encoder := NewDeferredEncoder(NewDefaultReaderWriter())

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				encoder.Encode(decodedValue)
				encoder.reset()
			}
		})

		b.Run("deferred_v2", func(b *testing.B) {

			decoder := getDeferredDecoder_V2()
			decodedValue := decoder.Decode()

			encoder := NewDeferredEncoder_V2(NewDefaultReaderWriter())

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				encoder.Encode(decodedValue)
				encoder.reset()
			}
		})
	})

}

func getDecoder() *Decoder {
	w := NewDefaultReaderWriter()
	encoder := NewEncoder(w)
	encoder.Encode(valueArray)

	return NewDecoder(w)
}

func getDeferredDecoder() *DeferredDecoder {
	w := NewDefaultReaderWriter()
	encoder := NewDeferredEncoder(w)
	encoder.Encode(valueArray)

	return NewDeferredDecoder(w)
}

func getDeferredDecoder_V2() *DeferredDecoder2 {
	w := NewDefaultReaderWriter()
	encoder := NewDeferredEncoder_V2(w)
	encoder.Encode(valueArray)

	return NewDeferredDecoder2(w)
}

func TestDecoding(t *testing.T) {
	decoder := getDecoder()

	decodedValue := decoder.Decode()

	// print the last value
	array, _ := decodedValue.([]interface{})
	lastValue := array[SIZE-1]
	fmt.Println(lastValue)
}

func TestDeferredDecoding(t *testing.T) {
	decoder := getDeferredDecoder()

	decodedValue := decoder.Decode()

	// print the last value
	array, _ := decodedValue.([]interface{})
	lastValue := array[SIZE-1].(*DeferredCompositeValue)
	lastValue.ensureLoaded()

	innerValue := lastValue.member(5).(*DeferredCompositeValue)

	// loading outer value does not load inner value
	assert.Nil(t, innerValue.value)

	innerValue.ensureLoaded()

	// now loaded
	assert.NotNil(t, innerValue.value)

	fmt.Println(lastValue.value)
	fmt.Println(innerValue.value)
}
