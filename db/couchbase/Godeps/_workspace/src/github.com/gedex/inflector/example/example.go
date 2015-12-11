package main

import (
	"fmt"

	"github.com/obieq/goar/db/couchbase/Godeps/_workspace/src/github.com/gedex/inflector"
)

var (
	singulars = [...]string{
		"Person", "Hero",
	}
	plurals = [...]string{
		"Tooth", "child",
	}
)

func main() {
	for _, s := range singulars {
		fmt.Printf("Plural of %v = %v\n", s, inflector.Pluralize(s))
	}

	fmt.Println()

	for _, s := range plurals {
		fmt.Printf("Singular of %v = %v\n", s, inflector.Singularize(s))
	}
}
