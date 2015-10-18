package postgres_test

import (
	"testing"

	. "github.com/obieq/goar"
	. "github.com/obieq/goar/db/postgres"
	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	AUTO_LIST_SP_NAME             = "spPostgresAutomobileList"
	AUTO_LIST_WITH_PARAMS_SP_NAME = "spPostgresAutomobileListWithParams"
)

type PostgresAutomobile struct {
	ArPostgres
	Automobile
	SafetyRating int
}

func (m *PostgresAutomobile) DBConnectionEnvironment() string {
	return "test" // NOTE: when using the goar package, this value should be pulled from ENV or config file
}

func (m *PostgresAutomobile) DBConnectionName() string {
	return "aws"
}

// GORM requires that we override the model name via the TableName() method
func (m *PostgresAutomobile) TableName() string {
	return "CustomPostgresTableNameAutos"
}

func (m *PostgresAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (model PostgresAutomobile) ToActiveRecord() *PostgresAutomobile {
	return ToAR(&model).(*PostgresAutomobile)
}

func (dbModel PostgresAutomobile) AssertDbPropertyMappings(model PostgresAutomobile, isDbUpdate bool) {
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

func TestPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Suite")
}

var _ = BeforeSuite(func() {
	auto := &PostgresAutomobile{}
	client := auto.ToActiveRecord().Client()

	// clean up previous test data
	client.DropTable(auto)

	// prep for new test run
	client.CreateTable(auto)
})

var _ = AfterSuite(func() {
})
