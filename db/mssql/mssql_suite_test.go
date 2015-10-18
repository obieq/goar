package mssql_test

import (
	"testing"

	. "github.com/obieq/goar"
	. "github.com/obieq/goar/db/mssql"
	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	AUTO_LIST_SP_NAME             = "spMsSqlAutomobileList"
	AUTO_LIST_WITH_PARAMS_SP_NAME = "spMsSqlAutomobileListWithParams"
)

type MsSqlAutomobile struct {
	ArMsSql
	Automobile
	SafetyRating int
}

// GORM requires that we override the model name via the TableName() method
func (m *MsSqlAutomobile) TableName() string {
	return "CustomMsSqlTableNameAutos"
}

func (m *MsSqlAutomobile) DBConnectionEnvironment() string {
	return "test" // NOTE: when using the goar package, this value should be pulled from ENV or config file
}

func (m *MsSqlAutomobile) DBConnectionName() string {
	return "aws"
}

func (m *MsSqlAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (model MsSqlAutomobile) ToActiveRecord() *MsSqlAutomobile {
	return ToAR(&model).(*MsSqlAutomobile)
}

func (dbModel MsSqlAutomobile) AssertDbPropertyMappings(model MsSqlAutomobile, isDbUpdate bool) {
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

func TestMsSql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MsSql Suite")
}

var _ = BeforeSuite(func() {
	auto := &MsSqlAutomobile{}
	client := auto.ToActiveRecord().Client()
	tblName := client.NewScope(auto).TableName()

	// clean up previous test data
	client.DropTable(auto)
	client.Exec("DROP PROCEDURE " + AUTO_LIST_SP_NAME + ";")
	client.Exec("DROP PROCEDURE " + AUTO_LIST_WITH_PARAMS_SP_NAME + ";")

	// prep for new test run
	client.CreateTable(auto)

	client.Exec("CREATE PROCEDURE " + AUTO_LIST_SP_NAME + " " +
		"AS " +
		"BEGIN SELECT * FROM dbo." + tblName + " " +
		"ORDER BY id ASC " +
		"END;")
	client.Exec("CREATE PROCEDURE " + AUTO_LIST_WITH_PARAMS_SP_NAME + " @Id int, @Model nvarchar(20) " +
		"AS " +
		"BEGIN SELECT * FROM dbo." + tblName + " " +
		"WHERE id > @Id AND model <> @Model " +
		"ORDER BY id ASC " +
		"END;")
})

var _ = AfterSuite(func() {
})
