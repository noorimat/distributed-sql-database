package parser

// Statement represents a SQL statement
type Statement interface {
	statementNode()
}

// Expression represents a SQL expression
type Expression interface {
	expressionNode()
}

// SelectStatement represents a SELECT query
type SelectStatement struct {
	Columns   []string    // column names or ["*"]
	TableName string
	Where     Expression  // optional WHERE clause
}

func (s *SelectStatement) statementNode() {}

// InsertStatement represents an INSERT query
type InsertStatement struct {
	TableName string
	Values    []Expression
}

func (i *InsertStatement) statementNode() {}

// CreateTableStatement represents a CREATE TABLE query
type CreateTableStatement struct {
	TableName string
	Columns   []ColumnDefinition
}

func (c *CreateTableStatement) statementNode() {}

// ColumnDefinition represents a column in CREATE TABLE
type ColumnDefinition struct {
	Name string
	Type string
}

// BinaryExpression represents binary operations (e.g., id = 123)
type BinaryExpression struct {
	Left     Expression
	Operator TokenType
	Right    Expression
}

func (b *BinaryExpression) expressionNode() {}

// Identifier represents a column or table name
type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}

// IntegerLiteral represents an integer value
type IntegerLiteral struct {
	Value string
}

func (i *IntegerLiteral) expressionNode() {}

// StringLiteral represents a string value
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) expressionNode() {}
