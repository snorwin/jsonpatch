package main

import (
	"fmt"

	"github.com/snorwin/jsonpatch"
)

type Job struct {
	Position  string `json:"position"`
	Company   string `json:"company"`
	Volunteer bool   `json:"volunteer"`
}

func main() {
	original := []Job{
		{Position: "IT Trainer", Company: "Powercoders", Volunteer: true},
		{Position: "Software Engineer", Company: "Github"},
	}
	updated := []Job{
		{Position: "Senior IT Trainer", Company: "Powercoders", Volunteer: true},
		{Position: "Senior Software Engineer", Company: "Github"},
	}

	patch, _ := jsonpatch.CreateJSONPatch(updated, original, jsonpatch.WithPredicate(jsonpatch.Funcs{
		ReplaceFunc: func(pointer jsonpatch.JSONPointer, value, _ interface{}) bool {
			// only update volunteering jobs
			if job, ok := value.(Job); ok {
				return job.Volunteer
			}
			return true
		},
	}))
	fmt.Println(patch.String())
}
