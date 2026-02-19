package parser

import (
	"testing"
)

func TestParseSelectStatement(t *testing.T) {
	input := "SELECT * FROM users;"
	
	l := NewLexer(input)
	p := NewParser(l)
	
	stmt := p.ParseStatement()
	
	if stmt == nil {
		t.Fatal("ParseStatement returned nil")
	}
	
	selectStmt, ok := stmt.(*SelectStatement)
	if !ok {
		t.Fatalf("stmt is not *SelectStatement. got=%T", stmt)
	}
	
	if selectStmt.TableName != "users" {
		t.Errorf("expected table name 'users', got %s", selectStmt.TableName)
	}
	
	if len(selectStmt.Columns) != 1 || selectStmt.Columns[0] != "*" {
		t.Errorf("expected columns ['*'], got %v", selectStmt.Columns)
	}
}

func TestParseSelectWithWhere(t *testing.T) {
	input := "SELECT name FROM users WHERE id = 123;"
	
	l := NewLexer(input)
	p := NewParser(l)
	
	stmt := p.ParseStatement()
	selectStmt := stmt.(*SelectStatement)
	
	if selectStmt.Where == nil {
		t.Fatal("WHERE clause is nil")
	}
	
	whereExpr, ok := selectStmt.Where.(*BinaryExpression)
	if !ok {
		t.Fatalf("WHERE is not BinaryExpression. got=%T", selectStmt.Where)
	}
	
	if whereExpr.Operator != EQ {
		t.Errorf("expected operator EQ, got %d", whereExpr.Operator)
	}
}

func TestParseInsertStatement(t *testing.T) {
	input := "INSERT INTO users VALUES (1, 'Alice');"
	
	l := NewLexer(input)
	p := NewParser(l)
	
	stmt := p.ParseStatement()
	insertStmt, ok := stmt.(*InsertStatement)
	if !ok {
		t.Fatalf("stmt is not *InsertStatement. got=%T", stmt)
	}
	
	if insertStmt.TableName != "users" {
		t.Errorf("expected table name 'users', got %s", insertStmt.TableName)
	}
	
	if len(insertStmt.Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(insertStmt.Values))
	}
}

func TestParseCreateTableStatement(t *testing.T) {
	input := "CREATE TABLE users (id INT, name VARCHAR);"
	
	l := NewLexer(input)
	p := NewParser(l)
	
	stmt := p.ParseStatement()
	createStmt, ok := stmt.(*CreateTableStatement)
	if !ok {
		t.Fatalf("stmt is not *CreateTableStatement. got=%T", stmt)
	}
	
	if createStmt.TableName != "users" {
		t.Errorf("expected table name 'users', got %s", createStmt.TableName)
	}
	
	if len(createStmt.Columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(createStmt.Columns))
	}
}
