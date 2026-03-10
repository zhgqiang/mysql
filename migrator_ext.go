package mysql

import (
	"fmt"
	"gorm.io/gorm"
)

func (m Migrator) HasTable(value interface{}) bool {
	var count int64

	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		currentDatabase, table := m.CurrentSchema(stmt, stmt.Table)
		return m.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ? AND table_type = ?", currentDatabase, table, "BASE TABLE").Row().Scan(&count)
	})

	return count > 0
}

// HasIndex check has index `name` or not
func (m Migrator) HasIndex(value interface{}, name string) bool {
	var count int64
	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		//currentDatabase := m.DB.Migrator().CurrentDatabase()
		currentDatabase, table := m.CurrentSchema(stmt, stmt.Table)
		if stmt.Schema != nil {
			if idx := stmt.Schema.LookIndex(name); idx != nil {
				name = idx.Name
			}
		}

		return m.DB.Raw(
			"SELECT count(*) FROM information_schema.statistics WHERE table_schema = ? AND table_name = ? AND index_name = ?",
			currentDatabase, table, name,
		).Row().Scan(&count)
	})

	return count > 0
}

func GetTableName(schema, tableName string) string {
	return fmt.Sprintf("%s.%s", schema, tableName)
}
