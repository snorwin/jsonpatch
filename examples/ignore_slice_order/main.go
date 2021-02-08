package main

import (
	"fmt"

	"github.com/snorwin/jsonpatch"
)

type Person struct {
	Name       string   `json:"name"`
	Pseudonyms []string `json:"pseudonyms"`
	Jobs       []Job    `json:"jobs"`
}

type Job struct {
	Position  string `json:"position"`
	Company   string `json:"company"`
	Volunteer bool   `json:"volunteer"`
}

func main() {
	original := Person{
		Name:       "John Doe",
		Pseudonyms: []string{"Jo", "JayD"},
		Jobs: []Job{
			{Position: "Software Engineer", Company: "Github"},
			{Position: "IT Trainer", Company: "Powercoders"},
		},
	}
	updated := Person{
		Name:       "John Doe",
		Pseudonyms: []string{"Jonny", "Jo"},
		Jobs: []Job{
			{Position: "IT Trainer", Company: "Powercoders", Volunteer: true},
			{Position: "Senior Software Engineer", Company: "Github"},
		},
	}

	patch, _ := jsonpatch.CreateJSONPatch(updated, original,
		jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{Pattern: "/*", JSONField: "company"}}),
		jsonpatch.IgnoreSliceOrder(),
	)
	fmt.Println(patch.String())
}
