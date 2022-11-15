package round_test

import (
	"encoding/json"
	"fmt"

	"github.com/portfoliotree/round"
)

type Data struct {
	Number   float64
	Map      map[FloatKey]float64 `precision:"3"`
	Embedded struct {
		SomeRate float64 `precision:"2,percent"`
	}
	List []float64
}

func Example() {
	var data Data
	data.Number = 1.1111
	data.Map = map[FloatKey]float64{
		4.5555555: 4.5555555,
		4.4444444: 4.4444444,
	}
	data.Embedded.SomeRate = 0.998765
	data.List = []float64{7, 6.656}

	fmt.Printf("before: %s\n", data.String())
	_ = round.Recursive(&data, 2)
	fmt.Printf("after:  %s", data.String())

	// Output: before: {
	// 	"Number": 1.1111,
	// 	"Map": {
	// 		"4.4444444": 4.4444444,
	// 		"4.5555555": 4.5555555
	// 	},
	// 	"Embedded": {
	// 		"SomeRate": 0.998765
	// 	},
	// 	"List": [
	// 		7,
	// 		6.656
	// 	]
	// }
	// after:  {
	// 	"Number": 1.11,
	// 	"Map": {
	// 		"4.4444444": 4.444,
	// 		"4.5555555": 4.556
	// 	},
	// 	"Embedded": {
	// 		"SomeRate": 99.88
	// 	},
	// 	"List": [
	// 		7,
	// 		6.66
	// 	]
	//}
}

func (d Data) String() string {
	buf, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

type FloatKey float64

func (fk FloatKey) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%g", float64(fk))), nil
}
