package UI

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	components "github.com/pclubiitk/dbcli/Components"
	"github.com/pclubiitk/dbcli/DB"
)

const (
	StepSourceCred = iota
	StepDestCred
	StepSelectSourceTable
	StepSelectSourceColumns
	StepSelectDestTable
	StepSelectDestColumns
	StepMapping
	StepDumpOption
	StepMigrationConfirm
)

type Model struct {
	Step int

	// ---- DB CREDENTIAL INPUTS ----
	SourceCred  map[string]string
	DestCred    map[string]string
	Source	    DB.DBInterface  //these are the most imp fields
	Dest        DB.DBInterface  //they are direct connections to databases
	CredInput   components.Input
	CredKeys    []string
	CredIndex   int
	IsSource    bool

	// ---- SOURCE DATABASE ----
	SourceTables      []string
	SelectedSourceTbl string
	SourceTableList   list.Model
	SourceColumns     []string
	SelectedSourceCols []string
	SourceColumnList  list.Model

	// ---- DESTINATION DATABASE ----
	DestTables      []string
	SelectedDestTbl string
	DestTableList   list.Model
	DestColumns     []string
	SelectedDestCols []string
	DestColumnList  list.Model

	// ---- MAPPING ----
	ColumnMapping 		map[string]string
	CurrentMapIdx       int
	ColumnBeingMapped   string
	MappingDestColumnList list.Model

	// ---- DUMP OPTION ----
	WantDump    *bool
	DumpPath    string
	DumpPathInp components.Input

	// ---- Migration State ----
	MigrationInProgress bool
	MigrationDone       bool
	MigrationProgress   float64 // 0.0 → 1.0
	StatusMsg           string
	Progress            progress.Model
	ProgressChan chan int
	DoneChan     chan error



	// ---- MISC ----
	ErrMsg string

	// ---- WINDOW PROPERTIES ----
	Width  int
	Height int
}

func (m Model) Init() tea.Cmd {
	return m.CredInput.Blink
}