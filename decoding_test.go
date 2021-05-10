package encoder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDecoding(b *testing.B) {

	b.Run("First field of last element", func(b *testing.B) {

		b.Run("normal", func(b *testing.B) {
			decoder := getDecoder()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.([]interface{})
				lastValue := array[SIZE-1].(*CompositeValue)
				lastValue.member(0)

				decoder.reset()
			}
		})

		b.Run("deferred", func(b *testing.B) {
			decoder := getDeferredDecoder()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.([]interface{})
				lastValue := array[SIZE-1].(*DeferredCompositeValue)
				innerValue := lastValue.member(5).(*DeferredCompositeValue)
				innerValue.member(0)

				decoder.reset()
			}
		})

		b.Run("deferred_v2", func(b *testing.B) {
			decoder := getDeferredDecoder_V2()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.(*DeferredArrayValue_V2)
				lastValue := array.member(SIZE - 1).(*DeferredCompositeValue_V2)
				lastValue.member(0)

				decoder.reset()
			}
		})

		b.Run("deferred_v3", func(b *testing.B) {
			decoder := getDeferredDecoder_V3()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.(*DeferredArrayValue_V3)
				lastValue := array.member(SIZE - 1).(*DeferredCompositeValue_V3)
				lastValue.member(0)

				decoder.reset()
			}
		})
	})

	b.Run("decoding", func(b *testing.B) {

		b.Run("normal", func(b *testing.B) {
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

		b.Run("deferred_v3", func(b *testing.B) {
			decoder := getDeferredDecoder_V3()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decoder.Decode()
				decoder.reset()
			}
		})
	})

	b.Run("decoding all", func(b *testing.B) {

		b.Run("normal", func(b *testing.B) {
			decoder := getDecoder()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.([]interface{})
				for i := 0; i < SIZE; i++ {
					lastValue := array[i].(*CompositeValue)
					for i := 0; i < 6; i++ {
						lastValue.member(i)
					}
				}

				decoder.reset()
			}
		})

		b.Run("deferred", func(b *testing.B) {
			decoder := getDeferredDecoder()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.([]interface{})
				for i := 0; i < SIZE; i++ {
					lastValue := array[i].(*DeferredCompositeValue)
					for i := 0; i < 6; i++ {
						lastValue.member(i)
					}
				}

				decoder.reset()
			}
		})

		b.Run("deferred_v2", func(b *testing.B) {
			decoder := getDeferredDecoder_V2()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.(*DeferredArrayValue_V2)

				for i := 0; i < SIZE; i++ {
					element := array.member(i).(*DeferredCompositeValue_V2)
					for i := 0; i < 6; i++ {
						element.member(i)
					}
				}

				decoder.reset()
			}
		})

		b.Run("deferred_v3", func(b *testing.B) {
			decoder := getDeferredDecoder_V3()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				decodedValue := decoder.Decode()
				array, _ := decodedValue.(*DeferredArrayValue_V3)
				for i := 0; i < SIZE; i++ {
					element := array.member(i).(*DeferredCompositeValue_V3)
					for i := 0; i < 6; i++ {
						element.member(i)
					}
				}

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

		b.Run("deferred_v3", func(b *testing.B) {
			encoder := NewDeferredEncoder_V3(NewDefaultReaderWriter())

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

		b.Run("deferred_v3", func(b *testing.B) {

			decoder := getDeferredDecoder_V3()
			decodedValue := decoder.Decode()

			encoder := NewDeferredEncoder_V3(NewDefaultReaderWriter())

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

func getDeferredDecoder_V3() *DeferredDecoder3 {
	w := NewDefaultReaderWriter()
	encoder := NewDeferredEncoder_V3(w)
	encoder.Encode(valueArray)

	return NewDeferredDecoder3(w)
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

func TestDeferredDecodingV3(t *testing.T) {
	decoder := getDeferredDecoder_V3()

	decodedValue := decoder.Decode()

	// print the last value
	array, _ := decodedValue.(*DeferredArrayValue_V3)
	lastValue := array.member(SIZE - 1).(*DeferredCompositeValue_V3)
	innerValue := lastValue.member(5).(*DeferredCompositeValue_V3)

	fmt.Println(lastValue.Name())
	fmt.Println(lastValue.member(0))
	fmt.Println(innerValue.member(0))
}
