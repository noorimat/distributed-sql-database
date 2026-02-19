package parser

// TokenType represents the type of token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	
	// Literals
	IDENT  // table names, column names
	INT    // 123
	STRING // 'hello'
	
	// Keywords
	SELECT
	FROM
	WHERE
	INSERT
	INTO
	VALUES
	CREATE
	TABLE
	UPDATE
	DELETE
	SET
	AND
	OR
	
	// Operators
	EQ     // =
	NEQ    // !=
	LT     // 
	GT     // >
	LTE    // <=
	GTE    // >=
	
	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	ASTERISK  // *
)

var keywords = map[string]TokenType{
	"SELECT": SELECT,
	"FROM":   FROM,
	"WHERE":  WHERE,
	"INSERT": INSERT,
	"INTO":   INTO,
	"VALUES": VALUES,
	"CREATE": CREATE,
	"TABLE":  TABLE,
	"UPDATE": UPDATE,
	"DELETE": DELETE,
	"SET":    SET,
	"AND":    AND,
	"OR":     OR,
}

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
