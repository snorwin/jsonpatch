# jsonpatch
[![GitHub Action](https://img.shields.io/badge/GitHub-Action-blue)](https://github.com/features/actions)
[![Documentation](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://pkg.go.dev/github.com/snorwin/jsonpatch)
[![Test](https://img.shields.io/github/workflow/status/snorwin/jsonpatch/Test?label=tests&logo=github)](https://github.com/snorwin/jsonpatch/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/snorwin/jsonpatch)](https://goreportcard.com/report/github.com/snorwin/jsonpatch)
[![Coverage Status](https://coveralls.io/repos/github/snorwin/jsonpatch/badge.svg?branch=main)](https://coveralls.io/github/snorwin/jsonpatch?branch=main)
[![Releases](https://img.shields.io/github/v/release/snorwin/jsonpatch)](https://github.com/snorwin/jsonpatch/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`jsonpatch` is a Go library to create JSON patches ([RFC6902](http://tools.ietf.org/html/rfc6902)) directly from arbitrary Go objects and facilitates the implementation of sophisticated custom (e.g. filtered, validated) patch creation.

## Basic Example

```go
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
```
```json
[{"op":"replace","path":"/name","value":"Jane Doe"},{"op":"replace","path":"/age","value":21}]
```

## Options
### Filter patches using Predicates
The option `WithPredicate` sets a patch `Predicate` which can be used to filter or validate the patch creation.
For each kind of patch (`add`, `remove` and `replace`) a dedicated filter function can be configured. The 
predicate will be checked before a patch is created, or the JSON object is processed further.

#### Example
```go
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
```

```json
[{"op":"replace","path":"/0/position","value":"Senior IT Trainer"}]
```


### Create partial patches
The option `WithPrefix` is used to specify a JSON pointer prefix if only a sub part of JSON structure needs to be patched,
but the patch still need to be applied on the entire JSON object.

#### Example
```go
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
```
```json
[{"op":"replace","path":"/jobs/0/position","value":"Senior IT Trainer"},{"op":"replace","path":"/jobs/0/volunteer","value":true},{"op":"add","path":"/jobs/1","value":{"position":"Software Engineer","company":"Github","volunteer":false}}]
```

### Ignore slice order
There are two options to ignore the slice order:
- `IgnoreSliceOrder` will ignore the order of all slices of built-in types (e.g. `int`, `string`) during the patch creation
  and will instead use the value itself in order to match and compare the current and modified JSON.
- `IgnoreSliceOrderWithPattern` allows to specify for which slices the order should be ignored using JSONPointer patterns (e.g. `/jobs`, `/jobs/*`).
  Furthermore, the slice order of structs (and pointer of structs) slices can be ignored by specifying a JSON field which should be used 
  to match the struct values. 

> NOTE: Ignoring the slice order only works if the elements (or the values used to match structs) are unique

#### Example
```go
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
```
```json
[{"op":"add","path":"/pseudonyms/2","value":"Jonny"},{"op":"remove","path":"/pseudonyms/1"},{"op":"replace","path":"/jobs/1/volunteer","value":true},{"op":"replace","path":"/jobs/0/position","value":"Senior Software Engineer"}]
```