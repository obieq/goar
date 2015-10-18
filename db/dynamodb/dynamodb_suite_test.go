package dynamodb

import (
	"testing"

	. "github.com/obieq/goar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Automobile2 struct {
	Vehicle2 `json:"Vehicle2,omitempty"`
}

type Vehicle2 struct {
	Year  int    `json:"year,omitempty"`
	Make  string `json:"make,omitempty"`
	Model string `json:"model,omitempty"`
}

type DynamodbAutomobile struct {
	ArDynamodb
	Year         int         `json:"year,omitempty"`
	Make         string      `json:"make,omitempty"`
	Model        string      `json:"model,omitempty"`
	SafetyRating int         `json:"safety_rating,omitempty"`
	Junk         interface{} `json:"junk,omitempty"`
	// CreatedAt    *time.Time  `json:"-"` // TODO: resolve the adroll mapping issue
	// UpdatedAt    *time.Time  `json:"-"` // TODO: resolve the adroll mapping issue
}

func (m *DynamodbAutomobile) CustomModelName() string {
	return "DynamodbAutomobiles"
}

func (m *DynamodbAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (m *ArDynamodb) DBConnectionEnvironment() string {
	return "test" // NOTE: when using the goar package, this value should be pulled from ENV or config file
}

func (m *ArDynamodb) DBConnectionName() string {
	return "aws"
}

func (model DynamodbAutomobile) ToActiveRecord() *DynamodbAutomobile {
	return ToAR(&model).(*DynamodbAutomobile)
}

func TestDynamodb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dynamodb Suite")
}

func (dbModel DynamodbAutomobile) AssertDbPropertyMappings(model DynamodbAutomobile, isDbUpdate bool) {
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
	//_, err := DynamodbAutomobile{}.ToActiveRecord().Truncate()
	//Expect(err).NotTo(HaveOccurred())

	// delete instances from prior test
	//ids := []string{"id1", "id2", "id3"}
	//for id := range ids {
	//DynamodbAutomobile{ArDynamodb: ArDynamodb{ID: ids[id]}}.ToActiveRecord().Delete()
	//}
})

var _ = AfterSuite(func() {
})
