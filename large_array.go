package encoder

import "fmt"

const SIZE = 100_000

var valueArray interface{} = func() []interface{} {

	values := make([]interface{}, SIZE)

	for i := 0; i < SIZE; i++ {
		values[i] = &CompositeValue{
			location: "TestLocation",
			typeName: "Person",
			kind:     999,
			fields: []interface{}{
				"John",
				"Doe",
				30,
				"male",
				"single",
				&CompositeValue{
					location: "TestLocation",
					typeName: "Address",
					kind:     999,
					fields: []interface{}{
						fmt.Sprintf("No: %d", i),
						"Vancouver",
						"BC",
						"Canada",
					},
				},
			},
		}
	}

	return values
}()
