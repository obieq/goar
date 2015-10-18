package orchestrate_test

import (
	"testing"

	. "github.com/obieq/goar"
	. "github.com/obieq/goar/db/orchestrate"
	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type IntegrationTestAutomobile struct {
	ArOrchestrate
	Automobile
	SafetyRating int         `json:"safety-rating,omitempty"` // NOTE: a dasherized json property name
	Junk         interface{} `json:"junk,omitempty"`
}

func (m *IntegrationTestAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (m *IntegrationTestAutomobile) DBConnectionEnvironment() string {
	return "test" // NOTE: when using the goar package, this value should be pulled from ENV or config file
}

func (m *IntegrationTestAutomobile) DBConnectionName() string {
	return "aws"
}

func (model IntegrationTestAutomobile) ToActiveRecord() *IntegrationTestAutomobile {
	return ToAR(&model).(*IntegrationTestAutomobile)
}

func TestOrchestrate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orchestrate Suite")
}

func (dbModel IntegrationTestAutomobile) AssertDbPropertyMappings(model IntegrationTestAutomobile, isDbUpdate bool) {
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
	_, err := IntegrationTestAutomobile{}.ToActiveRecord().Truncate()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
})
