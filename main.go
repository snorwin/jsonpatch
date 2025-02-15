package jsonpatch

type Complex struct {
	String string            `json:"string,omitempty"`
	Bolean bool              `json:"boolean"`
	Float  float64           `json:"float"`
	Uint   uint              `json:"uint"`
	Int    int               `json:"int"`
	Slice  []string          `json:"slice"`
	Map    map[string]string `json:"map"`
}

type Basic struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Complex Complex `json:"complex"`
}

func main() {
	// base := Basic{Name: "Matheus", Complex: Complex{
	// 	String: "a",
	// 	Bolean: true,
	// 	Float:  float64(1),
	// 	Uint:   uint(1),
	// 	Int:    int(1),
	// 	Slice:  []string{"a"},
	// 	Map:    map[string]string{"a": "a"},
	// }}
	// modified := Basic{Name: "Matheus", Complex: Complex{}}

	// p, err := jsonpatch.CreateJSONPatch(modified, base)
	// fmt.Println(p, err)
}
