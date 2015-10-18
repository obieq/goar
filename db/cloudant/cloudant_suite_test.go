package cloudant_test

import (
	"testing"

	c "github.com/obieq/go-cloudant"
	. "github.com/obieq/goar/db/cloudant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	cloudantTestDb *c.Database
)

func TestCloudant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloudant Suite")
}

var _ = BeforeSuite(func() {
	dbName := "golang_goar_suite_test"

	// drop database from prior test(s)
	Client().DeleteDatabase(dbName)

	// create database for test suite
	Client().CreateDatabase(dbName)

	// specify the database to use for test suite
	SetDbName(dbName)

	// drop databases from prior test(s)
	//err := Migration.DropDb(DbName())
	//Migration.DropDb(migrationDbName)

	//// prep for current test(s)
	//err = Migration.CreateDb(DbName())
	//Expect(err).NotTo(HaveOccurred())

	//err = Migration.CreateTable("rethink_db_automobiles")
	//Expect(err).NotTo(HaveOccurred())

	//err = Migration.CreateTable("callback_error_models")
	//Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
})
