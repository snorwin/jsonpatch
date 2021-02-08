package main

import (
	"fmt"

	"github.com/snorwin/jsonpatch"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	original := &Person{
		Name: "John Doe",
		Age:  42,
	}
	updated := &Person{
		Name: "Jane Doe",
		Age:  21,
	}

	patch, _ := jsonpatch.CreateJSONPatch(updated, original)
	fmt.Println(patch.String())
}
