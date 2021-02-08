package main

import (
	"fmt"

	"github.com/snorwin/jsonpatch"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Jobs []Job  `json:"jobs"`
}

type Job struct {
	Position  string `json:"position"`
	Company   string `json:"company"`
	Volunteer bool   `json:"volunteer"`
}

func main() {
	original := &Person{
		Name: "John Doe",
		Age:  42,
		Jobs: []Job{{Position: "IT Trainer", Company: "Powercoders"}},
	}
	updated := []Job{
		{Position: "Senior IT Trainer", Company: "Powercoders", Volunteer: true},
		{Position: "Software Engineer", Company: "Github"},
	}

	patch, _ := jsonpatch.CreateJSONPatch(updated, original.Jobs, jsonpatch.WithPrefix(jsonpatch.ParseJSONPointer("/jobs")))
	fmt.Println(patch.String())
}
