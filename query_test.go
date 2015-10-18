package goar

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query", func() {

	It("should have a New method", func() {
		q := NewQuery()
		Ω(q.Aggregations).ShouldNot(BeNil())
		Ω(q.Distinct).Should(BeFalse())
	})

})
