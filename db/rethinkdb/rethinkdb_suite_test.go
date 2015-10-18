package rethinkdb

import (
	"testing"

	r "github.com/dancannon/gorethink"
	. "github.com/obieq/goar"
	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRethinkDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RethinkDb Suite")
}

// Begin setting up ActiveRecord tests

type RethinkDbBaseModel struct {
	ArRethinkDb
}

type RethinkDbAutomobile struct {
	RethinkDbBaseModel
	Automobile
	SafetyRating int
	//Timestamps
	Junk interface{}
}

// Model for testing RethinkDb ActiveRecord error conditions
type ErrorTestingModel struct {
	ArRethinkDb
}

func (m *ArRethinkDb) DBConnectionEnvironment() string {
	return "test" // NOTE: when using the goar package, this value should be pulled from ENV or config file
}

func (m *ArRethinkDb) DBConnectionName() string {
	return "aws"
}

func (m RethinkDbAutomobile) ToActiveRecord() *RethinkDbAutomobile {
	return ToAR(&m).(*RethinkDbAutomobile)
}

func (m ErrorTestingModel) ToActiveRecord() *ErrorTestingModel {
	return ToAR(&m).(*ErrorTestingModel)
}

func (m *RethinkDbAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (m *ErrorTestingModel) Validate() {}

// End setting up ActiveRecord tests

// Begin setting up migration tests

var Migration = &RethinkDbMigration{}
var migrationDbName = "migration_test_db"
var migrationTestClient *r.Session
var rethinkTestDBName = "goar_test"

// End setting up migration tests

var _ = BeforeSuite(func() {

	// establish a db connection
	auto := RethinkDbAutomobile{}
	migrationTestClient = auto.ToActiveRecord().Client()

	// drop databases from prior test(s)
	err := Migration.DropDb(migrationTestClient, rethinkTestDBName)
	Migration.DropDb(migrationTestClient, migrationDbName)

	// prep for current test(s)
	err = Migration.CreateDb(migrationTestClient, rethinkTestDBName)
	Expect(err).NotTo(HaveOccurred())

	err = Migration.CreateTable(migrationTestClient, rethinkTestDBName, "rethink_db_automobiles")
	Expect(err).NotTo(HaveOccurred())

	err = Migration.CreateTable(migrationTestClient, rethinkTestDBName, "callback_error_models")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	migrationTestClient.Close()
})
