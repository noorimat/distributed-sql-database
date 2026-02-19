package parser

import (
	"fmt"
)

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}
	
	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()
	
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %d, got %d instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// ParseStatement parses a single SQL statement
func (p *Parser) ParseStatement() Statement {
	switch p.curToken.Type {
	case SELECT:
		return p.parseSelectStatement()
	case INSERT:
		return p.parseInsertStatement()
	case CREATE:
		return p.parseCreateTableStatement()
	default:
		return nil
	}
}

func (p *Parser) parseSelectStatement() *SelectStatement {
	stmt := &SelectStatement{}
	
	// Move past SELECT
	p.nextToken()
	
	// Parse columns
	if p.curToken.Type == ASTERISK {
		stmt.Columns = []string{"*"}
		p.nextToken()
	} else {
		// Parse column list (comma-separated identifiers)
		for {
			if p.curToken.Type == IDENT {
				stmt.Columns = append(stmt.Columns, p.curToken.Literal)
				p.nextToken()
				
				if p.curToken.Type != COMMA {
					break
				}
				p.nextToken() // skip comma
			} else {
				break
			}
		}
	}
	
	// Expect FROM
	if p.curToken.Type != FROM {
		return nil
	}
	p.nextToken()
	
	// Expect table name
	if p.curToken.Type != IDENT {
		return nil
	}
	stmt.TableName = p.curToken.Literal
	
	p.nextToken()
	
	// Optional WHERE clause
	if p.curToken.Type == WHERE {
		p.nextToken()
		stmt.Where = p.parseExpression()
	}
	
	return stmt
}

func (p *Parser) parseInsertStatement() *InsertStatement {
	stmt := &InsertStatement{}
	
	// Move past INSERT
	if !p.expectPeek(INTO) {
		return nil
	}
	
	// Expect table name
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.TableName = p.curToken.Literal
	
	// Expect VALUES
	if !p.expectPeek(VALUES) {
		return nil
	}
	
	// Expect (
	if !p.expectPeek(LPAREN) {
		return nil
	}
	
	// Parse values
	p.nextToken()
	for p.curToken.Type != RPAREN && p.curToken.Type != EOF {
		expr := p.parseExpression()
		stmt.Values = append(stmt.Values, expr)
		
		if p.peekToken.Type == COMMA {
			p.nextToken() // move to comma
			p.nextToken() // move past comma
		} else {
			p.nextToken()
			break
		}
	}
	
	return stmt
}

func (p *Parser) parseCreateTableStatement() *CreateTableStatement {
	stmt := &CreateTableStatement{}
	
	// Move past CREATE
	if !p.expectPeek(TABLE) {
		return nil
	}
	
	// Expect table name
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.TableName = p.curToken.Literal
	
	// Expect (
	if !p.expectPeek(LPAREN) {
		return nil
	}
	
	// Parse column definitions
	p.nextToken()
	for p.curToken.Type != RPAREN && p.curToken.Type != EOF {
		if p.curToken.Type == IDENT {
			colName := p.curToken.Literal
			p.nextToken()
			
			if p.curToken.Type == IDENT {
				colType := p.curToken.Literal
				stmt.Columns = append(stmt.Columns, ColumnDefinition{
					Name: colName,
					Type: colType,
				})
			}
		}
		
		if p.peekToken.Type == COMMA {
			p.nextToken() // move to comma
			p.nextToken() // move past comma
		} else {
			p.nextToken()
			break
		}
	}
	
	return stmt
}

func (p *Parser) parseExpression() Expression {
	var left Expression
	
	// Parse left side
	switch p.curToken.Type {
	case IDENT:
		left = &Identifier{Value: p.curToken.Literal}
	case INT:
		left = &IntegerLiteral{Value: p.curToken.Literal}
	case STRING:
		left = &StringLiteral{Value: p.curToken.Literal}
	default:
		return nil
	}
	
	// Check for binary operator
	if p.peekToken.Type == EQ || p.peekToken.Type == NEQ ||
		p.peekToken.Type == LT || p.peekToken.Type == GT ||
		p.peekToken.Type == LTE || p.peekToken.Type == GTE {
		
		p.nextToken() // move to operator
		operator := p.curToken.Type
		p.nextToken() // move to right side
		
		var right Expression
		switch p.curToken.Type {
		case IDENT:
			right = &Identifier{Value: p.curToken.Literal}
		case INT:
			right = &IntegerLiteral{Value: p.curToken.Literal}
		case STRING:
			right = &StringLiteral{Value: p.curToken.Literal}
		}
		
		return &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	
	return left
}
