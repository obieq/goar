package migration

type Migrator interface {
	CreateDb(dbName string) error
	DropDb(dbName string) error
	CreateTable(tableName string) error
	DropTable(tableName string) error
	AddIndex(tableName string, fields []string, opts map[string]interface{}) error
}
