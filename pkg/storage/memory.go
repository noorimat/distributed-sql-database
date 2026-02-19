package storage

import (
	"fmt"
	"sync"
)

// Table represents an in-memory table
type Table struct {
	Name    string
	Columns []Column
	Rows    []Row
	mu      sync.RWMutex
}

// Column represents a table column definition
type Column struct {
	Name string
	Type string
}

// Row represents a single row of data
type Row struct {
	Values []interface{}
}

// MemoryStorage manages in-memory tables
type MemoryStorage struct {
	tables map[string]*Table
	mu     sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tables: make(map[string]*Table),
	}
}

// CreateTable creates a new table
func (m *MemoryStorage) CreateTable(name string, columns []Column) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}
	
	m.tables[name] = &Table{
		Name:    name,
		Columns: columns,
		Rows:    []Row{},
	}
	
	return nil
}

// Insert adds a row to a table
func (m *MemoryStorage) Insert(tableName string, values []interface{}) error {
	m.mu.RLock()
	table, exists := m.tables[tableName]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	
	table.mu.Lock()
	defer table.mu.Unlock()
	
	if len(values) != len(table.Columns) {
		return fmt.Errorf("expected %d values, got %d", len(table.Columns), len(values))
	}
	
	table.Rows = append(table.Rows, Row{Values: values})
	return nil
}

// Select retrieves rows from a table
func (m *MemoryStorage) Select(tableName string, columns []string, filter func(Row) bool) ([]Row, []string, error) {
	m.mu.RLock()
	table, exists := m.tables[tableName]
	m.mu.RUnlock()
	
	if !exists {
		return nil, nil, fmt.Errorf("table %s does not exist", tableName)
	}
	
	table.mu.RLock()
	defer table.mu.RUnlock()
	
	// Determine which columns to return
	var columnIndices []int
	var columnNames []string
	
	if len(columns) == 1 && columns[0] == "*" {
		// Select all columns
		for i, col := range table.Columns {
			columnIndices = append(columnIndices, i)
			columnNames = append(columnNames, col.Name)
		}
	} else {
		// Select specific columns
		for _, colName := range columns {
			found := false
			for i, col := range table.Columns {
				if col.Name == colName {
					columnIndices = append(columnIndices, i)
					columnNames = append(columnNames, col.Name)
					found = true
					break
				}
			}
			if !found {
				return nil, nil, fmt.Errorf("column %s does not exist", colName)
			}
		}
	}
	
	// Filter and project rows
	var results []Row
	for _, row := range table.Rows {
		if filter == nil || filter(row) {
			// Project only requested columns
			projectedValues := make([]interface{}, len(columnIndices))
			for i, idx := range columnIndices {
				projectedValues[i] = row.Values[idx]
			}
			results = append(results, Row{Values: projectedValues})
		}
	}
	
	return results, columnNames, nil
}

// GetTable returns a table by name
func (m *MemoryStorage) GetTable(name string) (*Table, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	table, exists := m.tables[name]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", name)
	}
	
	return table, nil
}
