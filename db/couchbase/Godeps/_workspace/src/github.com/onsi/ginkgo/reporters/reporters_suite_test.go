package reporters_test

import (
	. "github.com/obieq/goar/db/couchbase/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/obieq/goar/db/couchbase/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestReporters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reporters Suite")
}
