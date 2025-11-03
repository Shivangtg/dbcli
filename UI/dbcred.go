package UI

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	components "github.com/pclubiitk/dbcli/Components"
	"github.com/pclubiitk/dbcli/DB"
)



func UpdateDBCred(m Model, msg tea.Msg) (Model,tea.Cmd) {
	var cmd tea.Cmd


	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Save input
			if m.IsSource {
				if m.SourceCred == nil {
					m.SourceCred = make(map[string]string)
				}
				m.SourceCred[m.CredKeys[m.CredIndex]] = m.CredInput.Value()
			} else {
				if m.DestCred == nil {
					m.DestCred = make(map[string]string)
				}
				m.DestCred[m.CredKeys[m.CredIndex]] = m.CredInput.Value()
			}

			// Reset input
			m.CredInput.SetValue("")
			m.CredIndex++
			if m.CredIndex >= len(m.CredKeys) {
				m.CredIndex = 0
				m.Step++
				var loweredVendor string
				if m.IsSource{
					loweredVendor = strings.ToLower(m.SourceCred["dbVendor"])
					m.SourceCred["dbVendor"] = loweredVendor
				}else{
					loweredVendor = strings.ToLower(m.DestCred["dbVendor"])
					m.DestCred["dbVendor"] = loweredVendor
				}
				
				if m.IsSource {
					switch loweredVendor {
					case "oracle":
						DB.ConnectOracle(m.SourceCred["host"], m.SourceCred["port"], m.SourceCred["dbname"], m.SourceCred["user"], m.SourceCred["password"])
						// Close any previously stored Source wrapper (if present)
						if m.Source != nil {
							_ = m.Source.Close()
						}
						m.Source = &DB.SQLWrapper{DB: DB.OracleDB}

						// Query tables owned by this user
						owner := strings.ToUpper(m.SourceCred["user"]) // Oracle stores owners uppercase
						query := "SELECT table_name FROM all_tables WHERE owner = :1 ORDER BY table_name"
						tables, err := DB.QueryStrings(m.Source, query, owner)
						if err != nil {
							m.ErrMsg = fmt.Sprintf("Failed to list tables: %v", err)
							log.Println("List tables error:", err)
							// keep SourceTables nil or empty
							m.SourceTables = []string{}
						} else {
							m.SourceTables = tables
							m.SourceTableList = components.CreateList(m.SourceTables, "📦 Source Tables",m.Width,m.Height)
						}
					case "mysql":
						DB.ConnectMySQL(m.SourceCred["host"], m.SourceCred["port"], m.SourceCred["dbname"], m.SourceCred["user"], m.SourceCred["password"])
						if m.Source != nil {
							_ = m.Source.Close()
						}
						m.Source = &DB.GormWrapper{DB: DB.MySQLDB}
						// MySQL: list tables in current database
						query := "SELECT table_name FROM information_schema.tables WHERE table_schema = ? ORDER BY table_name"
						tables, err := DB.QueryStrings(m.Source, query, m.SourceCred["dbname"])
						if err != nil {
							m.ErrMsg = fmt.Sprintf("Failed to list MySQL tables: %v", err)
							m.SourceTables = []string{}
						} else {
							m.SourceTables = tables
							m.SourceTableList = components.CreateList(m.SourceTables, "📦 Source Tables",m.Width,m.Height)
						}
					default:
						m.ErrMsg = "Unknown DB vendor for source"
					}
				} else {
					// Destination block (similar to above)
					switch loweredVendor {
					case "oracle":
						DB.ConnectOracle(m.DestCred["host"], m.DestCred["port"], m.DestCred["dbname"], m.DestCred["user"], m.DestCred["password"])
						if m.Dest != nil {
							_ = m.Dest.Close()
						}
						m.Dest = &DB.SQLWrapper{DB: DB.OracleDB}

						owner := strings.ToUpper(m.DestCred["user"])
						query := "SELECT table_name FROM all_tables WHERE owner = :1 ORDER BY table_name"
						tables, err := DB.QueryStrings(m.Dest, query, owner)
						if err != nil {
							m.ErrMsg = fmt.Sprintf("Failed to list dest tables: %v", err)
							m.DestTables = []string{}
						} else {
							m.DestTables = tables
							m.DestTableList = components.CreateList(m.DestTables, "📦 Destination Tables",m.Width,m.Height)
						}
					case "mysql":
						DB.ConnectMySQL(m.DestCred["host"], m.DestCred["port"], m.DestCred["dbname"], m.DestCred["user"], m.DestCred["password"])
						if m.Dest != nil {
							_ = m.Dest.Close()
						}
						m.Dest = &DB.GormWrapper{DB: DB.MySQLDB}
						query := "SELECT table_name FROM information_schema.tables WHERE table_schema = ? ORDER BY table_name"
						tables, err := DB.QueryStrings(m.Dest, query, m.DestCred["dbname"])
						if err != nil {
							m.ErrMsg = fmt.Sprintf("Failed to list dest tables: %v", err)
							m.DestTables = []string{}
						} else {
							m.DestTables = tables
							m.DestTableList = components.CreateList(m.DestTables, "📦 Destination Tables",m.Width,m.Height)
						}
					default:
						m.ErrMsg = "Unknown DB vendor for dest"
					}
				}
				m.IsSource = !m.IsSource
				m.CredInput, cmd = m.CredInput.Update(tea.WindowSize)
				return m , cmd

			}
		case "backspace":
			if len(m.CredInput.Value()) > 0 {
				m.CredInput.SetValue(m.CredInput.Value()[:len(m.CredInput.Value())-1])
			}
		default:
			m.CredInput.SetValue(m.CredInput.Value() + msg.String())
		}
	}
	m.CredInput, cmd = m.CredInput.Update(msg)
	return m , cmd
}
func ViewDBCred(m Model) string {
	dbType := "Source"
	if !m.IsSource {
		dbType = "Destination"
	}
	title    := titleStyle.Render(fmt.Sprintf("🔐 Enter %s Database Credentials", dbType))
	info_up  := infoStyle.Render(fmt.Sprintf("Step %d of %d", m.CredIndex+1, len(m.CredKeys)))
	key      := m.CredKeys[m.CredIndex]
	label    := labelStyle.Render(fmt.Sprintf("%s:", key))
	input    := inputStyle.Render(m.CredInput.Value())
	input_    := lipgloss.JoinHorizontal(lipgloss.Center,label,input)
	info_down :=infoStyle.Render("Type your input and press Enter ↵ to continue.")


	if m.ErrMsg != "" {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("Error: %s", m.ErrMsg)),
		)
	}

	if m.CredIndex < len(m.CredKeys) {
		
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			info_up,
			input_,
			info_down,
			),
		)
		
	} else if (m.CredIndex == len(m.CredKeys)){
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			msgStyle.Render(fmt.Sprintf("%s Database Credentials entered successfully", dbType)),
		)
	}
	return ""
}
