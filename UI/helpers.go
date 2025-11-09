package UI

import (
	"fmt"
	"os"
	"path/filepath"
)

func Contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func Remove(slice []string, val string) []string {
	result := []string{}
	for _, s := range slice {
		if s != val {
			result = append(result, s)
		}
	}
	return result
}


func createDumpFile(path string, m Model) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	// Create or overwrite dump file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer f.Close()

	// --- write a simple dump summary ---
	fmt.Fprintf(f, "-- DBCLI Dump\n")
	fmt.Fprintf(f, "-- Source DB: %s\n", m.SourceCred["dbname"])
	fmt.Fprintf(f, "-- Destination DB: %s\n", m.DestCred["dbname"])
	fmt.Fprintf(f, "-- Tables: %s -> %s\n\n", m.SelectedSourceTbl, m.SelectedDestTbl)

	fmt.Fprintf(f, "-- Column Mapping:\n")
	for src, dst := range m.ColumnMapping {
		fmt.Fprintf(f, "%s -> %s\n", src, dst)
	}

	fmt.Fprintf(f, "\n-- Selected Source Columns: %v\n", m.SelectedSourceCols)
	fmt.Fprintf(f, "-- Selected Destination Columns: %v\n", m.SelectedDestCols)

	return f.Sync()
}
