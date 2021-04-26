# Deferred Encoder/Decoder

A decoder that supports lazy loading.

## Design:
For each composite value, encode as:
```
<type-tag>  <content-length>  <content>
```

### Decoding:
- Read the length and, load the raw-bytes. Do not try to decode.
- Create a new wrapper `CompositeDeferredValue`, that implements the same interfaces as`CompositeValue`.
  - Internally holds the raw bytes.
  - Decode the raw bytes on-demand, if the value actually is used (e.g: access a member, etc.). This again a
    shallow loading. Only loaded upto a one-level down.
  - Decoding the raw bytes will give a `CompositeValue`. Cache it in the wrapper.
- Once loaded, `CompositeDeferredValue` acts as a delegator for `CompositeValue`.
    
Example:

Array with 5 composite values, after decoding:

| Index | Value at Index | Wrapped Content |
| ----- | -------------- | --------------- |
| 0 | CompositeDeferredValue | raw bytes |
| 1 | CompositeDeferredValue | raw bytes |
| 2 | CompositeDeferredValue | raw bytes |
| 3 | CompositeDeferredValue | raw bytes |
| 4 | CompositeDeferredValue | raw bytes |

Getting element at index-3, and accessing its members (cause the element at 3 to be loaded/decoded):

| Index | Value at Index | Wrapped Content |
| ----- | -------------- | ------- |
| 0 | CompositeDeferredValue | raw bytes |
| 1 | CompositeDeferredValue | raw bytes |
| 2 | CompositeDeferredValue | raw bytes |
| **_3_** | **_CompositeDeferredValue_** | **_CompositeValue_** |
| 4 | CompositeDeferredValue | raw bytes |

### Encoding a decoded value:
- For fully built/loaded values, encode them in the normal way.
- For values that are not loaded (`CompositeDeferredValue` with raw-byte content), dump the raw bytes, as is.

## Benchmark results

_**NOTE:** Currently, uses a mocked byte reader/writer for writing bytes to the target.
Actual values may defer based on the low level API used for byte reading/writing._

```
Decoding:
---------
BenchmarkDecoding/decoding/normal-12           	      13	  90333817 ns/op     54317676 B/op	 2500003 allocs/op
BenchmarkDecoding/decoding/deferred-12         	     100	  12733971 ns/op      5005673 B/op	  200003 allocs/op

Encoding:
---------
BenchmarkDecoding/encoding/normal-12  	               1	2113130420 ns/op    566249104 B/op	 4400050 allocs/op
BenchmarkDecoding/encoding/deferred-12         	       2	1230219404 ns/op   1011646116 B/op	 5700027 allocs/op

Re-encoding (encoding back a decoded one):
-----------
BenchmarkDecoding/re-encoding/normal-12        	       3	 436506383 ns/op    575716240 B/op	 4400020 allocs/op
BenchmarkDecoding/re-encoding/deferred-12      	       6	 331096581 ns/op    526461602 B/op	  100011 allocs/op
```
