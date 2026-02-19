package sql

import (
	"testing"
	
	"github.com/noorimat/distributed-sql-database/pkg/sql/parser"
	"github.com/noorimat/distributed-sql-database/pkg/storage"
)

func TestExecutor_CreateAndInsert(t *testing.T) {
	store := storage.NewMemoryStorage()
	executor := NewExecutor(store)
	
	// CREATE TABLE
	createSQL := "CREATE TABLE users (id INT, name VARCHAR);"
	l := parser.NewLexer(createSQL)
	p := parser.NewParser(l)
	stmt := p.ParseStatement()
	
	result, err := executor.Execute(stmt)
	if err != nil {
		t.Fatalf("CREATE TABLE failed: %v", err)
	}
	if result.Message == "" {
		t.Error("Expected success message")
	}
	
	// INSERT
	insertSQL := "INSERT INTO users VALUES (1, 'Alice');"
	l = parser.NewLexer(insertSQL)
	p = parser.NewParser(l)
	stmt = p.ParseStatement()
	
	result, err = executor.Execute(stmt)
	if err != nil {
		t.Fatalf("INSERT failed: %v", err)
	}
	if result.RowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", result.RowsAffected)
	}
}

func TestExecutor_Select(t *testing.T) {
	store := storage.NewMemoryStorage()
	executor := NewExecutor(store)
	
	// Setup: Create table and insert data
	createSQL := "CREATE TABLE users (id INT, name VARCHAR);"
	l := parser.NewLexer(createSQL)
	p := parser.NewParser(l)
	executor.Execute(p.ParseStatement())
	
	insertSQL := "INSERT INTO users VALUES (1, 'Alice');"
	l = parser.NewLexer(insertSQL)
	p = parser.NewParser(l)
	executor.Execute(p.ParseStatement())
	
	// SELECT
	selectSQL := "SELECT * FROM users;"
	l = parser.NewLexer(selectSQL)
	p = parser.NewParser(l)
	stmt := p.ParseStatement()
	
	result, err := executor.Execute(stmt)
	if err != nil {
		t.Fatalf("SELECT failed: %v", err)
	}
	
	if len(result.Rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(result.Rows))
	}
	
	if len(result.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(result.Columns))
	}
}
