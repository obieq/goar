package couchbase_test

import (
	"testing"

	. "github.com/obieq/goar"
	. "github.com/obieq/goar/db/couchbase"
	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CouchbaseAutomobile struct {
	ArCouchbase
	Automobile
	SafetyRating int
}

func (m *CouchbaseAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (model CouchbaseAutomobile) ToActiveRecord() *CouchbaseAutomobile {
	return ToAR(&model).(*CouchbaseAutomobile)
}

func TestCouchbase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Couchbase Suite")
}

func (dbModel CouchbaseAutomobile) AssertDbPropertyMappings(model CouchbaseAutomobile, isDbUpdate bool) {
	Ω(dbModel.ID).Should(Equal(model.ID))
	Ω(dbModel.Year).Should(Equal(model.Year))
	Ω(dbModel.Make).Should(Equal(model.Make))
	Ω(dbModel.Model).Should(Equal(model.Model))
	Ω(dbModel.SafetyRating).Should(Equal(model.SafetyRating))

	Ω(dbModel.CreatedAt).ShouldNot(BeNil())
	if isDbUpdate {
		Ω(dbModel.UpdatedAt).ShouldNot(BeNil())
	} else {
		Ω(dbModel.UpdatedAt).Should(BeNil())
	}
}

var _ = BeforeSuite(func() {
	// drop collections from previous tests
	//_, err := CouchbaseAutomobile{}.ToActiveRecord().Truncate()
	//Expect(err).NotTo(HaveOccurred())

	// delete instances from prior test
	ids := []string{"id1", "id2", "id3"}
	for id := range ids {
		CouchbaseAutomobile{ArCouchbase: ArCouchbase{ID: ids[id]}}.ToActiveRecord().Delete()
	}
})

var _ = AfterSuite(func() {
})
