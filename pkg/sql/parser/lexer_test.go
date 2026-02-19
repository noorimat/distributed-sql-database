package parser

import (
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	input := `SELECT * FROM users WHERE id = 123;`
	
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{SELECT, "SELECT"},
		{ASTERISK, "*"},
		{FROM, "FROM"},
		{IDENT, "users"},
		{WHERE, "WHERE"},
		{IDENT, "id"},
		{EQ, "="},
		{INT, "123"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}
	
	l := NewLexer(input)
	
	for i, tt := range tests {
		tok := l.NextToken()
		
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d",
				i, tt.expectedType, tok.Type)
		}
		
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_InsertStatement(t *testing.T) {
	input := `INSERT INTO users VALUES (1, 'Alice');`
	
	l := NewLexer(input)
	
	tok := l.NextToken()
	if tok.Type != INSERT {
		t.Fatalf("expected INSERT, got %d", tok.Type)
	}
	
	tok = l.NextToken()
	if tok.Type != INTO {
		t.Fatalf("expected INTO, got %d", tok.Type)
	}
}
