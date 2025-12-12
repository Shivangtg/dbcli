package UI

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/pclubiitk/dbcli/DB"
)

// --- Utility: log executed queries ---
func logQuery(vendor string, query string, args []interface{}, start time.Time) {
	f, err := os.OpenFile("migration_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("⚠️ Could not open migration log file: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	duration := time.Since(start).Milliseconds()

	formattedArgs := []string{}
	for _, a := range args {
		formattedArgs = append(formattedArgs, fmt.Sprintf("%v", a))
	}

	entry := fmt.Sprintf(
		"[%s] [DB: %s] [Exec: %dms]\nQuery: %s\nArgs: [%s]\n\n",
		timestamp,
		strings.ToUpper(vendor),
		duration,
		query,
		strings.Join(formattedArgs, ", "),
	)
	f.WriteString(entry)
}

// --- Utility: compare two rows ---
func rowDiffers(srcVals, destVals []interface{}) bool {
	if len(srcVals) != len(destVals) {
		return true
	}
	for i := range srcVals {
		if !reflect.DeepEqual(srcVals[i], destVals[i]) {
			return true
		}
	}
	return false
}

// --- Main migration function ---
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

	// Ensure 'id' exists
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

	// Add id if not selected
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

	// Query source data
	selectQuery := fmt.Sprintf("SELECT %s FROM %s", strings.Join(srcCols, ", "), srcTable)
	rows, err := m.Source.RawQuery(selectQuery)
	if err != nil {
		doneChan <- fmt.Errorf("failed to fetch data from source: %v", err)
		return
	}
	defer rows.Close()

	// --- Count total rows ---
	totalRows := 0
	for rows.Next() {
		totalRows++
	}
	rows.Close()

	if totalRows == 0 {
		doneChan <- fmt.Errorf("no rows to migrate")
		return
	}

	// --- Re-fetch actual data for migration ---
	rows, err = m.Source.RawQuery(selectQuery)
	if err != nil {
		doneChan <- fmt.Errorf("failed to refetch data: %v", err)
		return
	}
	defer rows.Close()

	destVendor := m.DestCred["dbVendor"]

	// Normalize table case
	if destVendor == "oracle" {
		destTable = strings.ToUpper(destTable)
	} else if destVendor == "mysql" {
		destTable = strings.ToLower(destTable)
	}

	// Map source→destination column names
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

		// 1️⃣ Check if row exists
		checkQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = %s",
			destTable, DB.Placeholder(destVendor, 1))
		start := time.Now()
		logQuery(destVendor, checkQuery, []interface{}{idVal}, start)

		checkRows, err := m.Dest.RawQuery(checkQuery, idVal)
		if err != nil {
			doneChan <- fmt.Errorf("failed checking id existence: %v", err)
			return
		}
		var count int
		if checkRows.Next() {
			checkRows.Scan(&count)
		}
		checkRows.Close()

		// 2️⃣ If exists → compare row
		if count > 0 {
			destQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id = %s",
				strings.Join(finalDestCols, ", "), destTable, DB.Placeholder(destVendor, 1))
			destRows, err := m.Dest.RawQuery(destQuery, idVal)
			if err != nil {
				doneChan <- fmt.Errorf("failed reading destination row: %v", err)
				return
			}

			destVals := make([]interface{}, len(finalDestCols))
			destPtrs := make([]interface{}, len(finalDestCols))
			for i := range destPtrs {
				destPtrs[i] = &destVals[i]
			}

			if destRows.Next() {
				destRows.Scan(destPtrs...)
			}
			destRows.Close()

			// Compare — skip update if identical
			if !rowDiffers(values, destVals) {
				progress++
				continue
			}

			// 3️⃣ Update only changed
			setParts := []string{}
			args := []interface{}{}
			argIndex := 1
			for i, col := range finalDestCols {
				if strings.EqualFold(col, "id") {
					continue
				}
				setParts = append(setParts, fmt.Sprintf("%s = %s", col, DB.Placeholder(destVendor, argIndex)))
				args = append(args, values[i])
				argIndex++
			}
			args = append(args, idVal)

			updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = %s",
				destTable, strings.Join(setParts, ", "), DB.Placeholder(destVendor, argIndex))

			start = time.Now()
			logQuery(destVendor, updateQuery, args, start)
			if err := m.Dest.ExecQuery(updateQuery, args...); err != nil {
				doneChan <- fmt.Errorf("update failed for id=%v: %v", idVal, err)
				return
			}
			if destVendor == "oracle" {
				_ = m.Dest.ExecQuery("COMMIT")
			}

		} else {
			// 4️⃣ Insert new row
			insertCols := []string{}
			placeholders := []string{}
			args := []interface{}{}
			argIndex := 1

			for i, col := range finalDestCols {
				insertCols = append(insertCols, col)
				placeholders = append(placeholders, DB.Placeholder(destVendor, argIndex))
				args = append(args, values[i])
				argIndex++
			}

			insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
				destTable, strings.Join(insertCols, ", "), strings.Join(placeholders, ", "))

			start = time.Now()
			logQuery(destVendor, insertQuery, args, start)
			if err := m.Dest.ExecQuery(insertQuery, args...); err != nil {
				doneChan <- fmt.Errorf("insert failed for id=%v: %v", idVal, err)
				return
			}
			if destVendor == "oracle" {
				_ = m.Dest.ExecQuery("COMMIT")
			}
		}

		// 5️⃣ Update progress
		progress++
		percent := progress * 100 / totalRows
		select {
		case progressChan <- percent:
		default:
		}
		time.Sleep(20 * time.Millisecond)
	}

	doneChan <- nil
}
