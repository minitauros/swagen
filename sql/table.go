package sql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type (
	// TableService provides functionality related to database tables.
	TableService struct {
		db *sql.DB
	}
)

// NewTableService returns a new table service.
// DB is typically a *sql.DB database.
func NewTableService(db *sql.DB) TableService {
	return TableService{db: db}
}

// GetColumns returns the columns from the table with the given name.
func (s TableService) GetColumns(tableName string) ([]*sql.ColumnType, error) {
	rows, err := s.db.Query(`SELECT * FROM ` + tableName + ` LIMIT 1`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cols, nil
}
