package UI

import (
	"fmt"
	"strings"
	"time"
)

// MigrateData performs the actual migration and reports progress.
func MigrateData(m Model, progressChan chan<- int, doneChan chan<- error) {
	defer func() {
		close(progressChan)
		close(doneChan)
	}()

	srcTable := m.SelectedSourceTbl
	destTable := m.SelectedDestTbl
	srcCols := append([]string{}, m.SelectedSourceCols...)
	mapping := m.ColumnMapping

	if srcTable == "" || destTable == "" {
		doneChan <- fmt.Errorf("source or destination table not selected")
		return
	}
	if len(srcCols) == 0 {
		doneChan <- fmt.Errorf("no source columns selected")
		return
	}

	// --- Ensure 'id' exists ---
	hasID := false
	for _, c := range m.SourceColumns {
		if strings.EqualFold(c, "id") {
			hasID = true
			break
		}
	}
	if !hasID {
		doneChan <- fmt.Errorf("source table must include an 'id' column for upsert logic")
		return
	}
	userSelectedID := false
	for _, c := range srcCols {
		if strings.EqualFold(c, "id") {
			userSelectedID = true
			break
		}
	}
	if !userSelectedID {
		srcCols = append(srcCols, "id")
	}

	selectQuery := fmt.Sprintf("SELECT %s FROM %s", strings.Join(srcCols, ", "), srcTable)
	rows, err := m.Source.RawQuery(selectQuery)
	if err != nil {
		doneChan <- fmt.Errorf("failed to fetch data from source: %v", err)
		return
	}
	defer rows.Close()

	// Count total rows
	totalRows := 0
	for rows.Next() {
		totalRows++
	}
	rows.Close()

	if totalRows == 0 {
		doneChan <- fmt.Errorf("no rows to migrate")
		return
	}

	// Re-fetch actual data
	rows, err = m.Source.RawQuery(selectQuery)
	if err != nil {
		doneChan <- fmt.Errorf("failed to refetch data: %v", err)
		return
	}
	defer rows.Close()

	finalDestCols := make([]string, len(srcCols))
	for i, srcCol := range srcCols {
		if mapped, ok := mapping[srcCol]; ok {
			finalDestCols[i] = mapped
		} else {
			finalDestCols[i] = srcCol
		}
	}

	values := make([]interface{}, len(srcCols))
	ptrs := make([]interface{}, len(srcCols))
	for i := range ptrs {
		ptrs[i] = &values[i]
	}

	progress := 0
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			doneChan <- fmt.Errorf("failed to scan source row: %v", err)
			return
		}

		var idVal interface{}
		for i, colName := range srcCols {
			if strings.EqualFold(colName, "id") {
				idVal = values[i]
				break
			}
		}

		if idVal == nil {
			doneChan <- fmt.Errorf("missing id value in source row")
			return
		}

		// Check if exists
		checkQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = $1", destTable)
		checkRows, err := m.Dest.RawQuery(checkQuery, idVal)
		if err != nil {
			doneChan <- fmt.Errorf("failed checking id existence: %v", err)
			return
		}

		var count int
		if checkRows.Next() {
			if err := checkRows.Scan(&count); err != nil {
				checkRows.Close()
				doneChan <- fmt.Errorf("failed scanning id existence: %v", err)
				return
			}
		}
		checkRows.Close()

		if count > 0 {
			// Update
			setParts := []string{}
			args := []interface{}{}
			argIndex := 1
			for i, col := range finalDestCols {
				if strings.EqualFold(col, "id") {
					continue
				}
				setParts = append(setParts, fmt.Sprintf("%s = $%d", col, argIndex))
				args = append(args, values[i])
				argIndex++
			}
			args = append(args, idVal)

			updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", destTable, strings.Join(setParts, ", "), argIndex)
			if err := m.Dest.ExecQuery(updateQuery, args...); err != nil {
				doneChan <- fmt.Errorf("update failed for id=%v: %v", idVal, err)
				return
			}
		} else {
			// Insert
			insertCols := []string{}
			placeholders := []string{}
			args := []interface{}{}
			argIndex := 1

			for i, col := range finalDestCols {
				insertCols = append(insertCols, col)
				placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
				args = append(args, values[i])
				argIndex++
			}

			insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
				destTable,
				strings.Join(insertCols, ", "),
				strings.Join(placeholders, ", "),
			)
			if err := m.Dest.ExecQuery(insertQuery, args...); err != nil {
				doneChan <- fmt.Errorf("insert failed for id=%v: %v", idVal, err)
				return
			}
		}

		progress++
		percent := progress * 100 / totalRows

		// Non-blocking send to progressChan
		select {
		case progressChan <- percent:
		default:
		}

		// Small delay for smooth UI updates (optional)
		time.Sleep(30 * time.Millisecond)
	}

	doneChan <- nil
}
