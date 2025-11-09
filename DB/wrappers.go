package DB

import (
	"database/sql"
	"gorm.io/gorm"
)

// DBInterface is a common interface for any database type
type DBInterface interface {
    RawQuery(query string, args ...interface{}) (*sql.Rows, error)
    ExecQuery(query string, args ...interface{}) error
    Close() error
}

// --------------------------------------------------------

type SQLWrapper struct {
    DB *sql.DB
}

func (s *SQLWrapper) RawQuery(query string, args ...interface{}) (*sql.Rows, error) {
    return s.DB.Query(query, args...)
}

func (s *SQLWrapper) ExecQuery(query string, args ...interface{}) error {
    _, err := s.DB.Exec(query, args...)
    return err
}

func (s *SQLWrapper) Close() error {
    return s.DB.Close()
}

// --------------------------------------------------------

type GormWrapper struct {
    DB *gorm.DB
}

func (g *GormWrapper) RawQuery(query string, args ...interface{}) (*sql.Rows, error) {
    // GORM's Raw() returns *gorm.DB, which internally has *sql.Rows
    tx := g.DB.Raw(query, args...)
    return tx.Rows()
}

func (g *GormWrapper) ExecQuery(query string, args ...interface{}) error {
    tx := g.DB.Exec(query, args...)
    return tx.Error
}

func (g *GormWrapper) Close() error {
    sqlDB, err := g.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}
