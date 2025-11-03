// DB/utils.go
package DB

import (
    "database/sql"
    "errors"
)

// RowsToStrings reads the first column from rows into a []string
func RowsToStrings(rows *sql.Rows) ([]string, error) {
    if rows == nil {
        return nil, errors.New("rows is nil")
    }
    defer rows.Close()
    var result []string
    for rows.Next() {
        var s sql.NullString
        if err := rows.Scan(&s); err != nil {
            return nil, err
        }
        if s.Valid {
            result = append(result, s.String)
        } else {
            result = append(result, "")
        }
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return result, nil
}

// QueryStrings executes a query using DBInterface and returns first-column strings
func QueryStrings(db DBInterface, query string, args ...interface{}) ([]string, error) {
    rows, err := db.RawQuery(query, args...)
    if err != nil {
        return nil, err
    }
    return RowsToStrings(rows)
}

// ExecSQL helper
func ExecSQL(db DBInterface, query string, args ...interface{}) error {
    return db.ExecQuery(query, args...)
}

// Convenience function: close DBInterface if non-nil
func CloseIfSet(db DBInterface) error {
    if db == nil {
        return nil
    }
    return db.Close()
}
