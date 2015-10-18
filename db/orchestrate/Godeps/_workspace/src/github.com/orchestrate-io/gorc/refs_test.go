// Copyright 2014, Orchestrate.IO, Inc.

package gorc

import (
	"testing"
	"testing/quick"
)

func TestRefsHasNext(t *testing.T) {
	f := func(results *RefResults) bool {
		return !(results.Next == "" && results.HasNext())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRefIsDeleted(t *testing.T) {
	f := func(result *RefResult) bool {
		return !(result.Path.Tombstone == false && result.IsDeleted()) &&
			!(result.Path.Tombstone == true && !result.IsDeleted())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
