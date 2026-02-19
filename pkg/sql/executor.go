package sql

import (
	"fmt"
	"strconv"
	
	"github.com/noorimat/distributed-sql-database/pkg/sql/parser"
	"github.com/noorimat/distributed-sql-database/pkg/storage"
)

// Executor executes parsed SQL statements
type Executor struct {
	storage *storage.MemoryStorage
}

// NewExecutor creates a new SQL executor
func NewExecutor(storage *storage.MemoryStorage) *Executor {
	return &Executor{
		storage: storage,
	}
}

// Execute runs a SQL statement and returns results
func (e *Executor) Execute(stmt parser.Statement) (*Result, error) {
	switch s := stmt.(type) {
	case *parser.SelectStatement:
		return e.executeSelect(s)
	case *parser.InsertStatement:
		return e.executeInsert(s)
	case *parser.CreateTableStatement:
		return e.executeCreateTable(s)
	default:
		return nil, fmt.Errorf("unsupported statement type")
	}
}

// Result represents query execution results
type Result struct {
	Columns      []string
	Rows         [][]interface{}
	RowsAffected int
	Message      string
}

func (e *Executor) executeSelect(stmt *parser.SelectStatement) (*Result, error) {
	// Build filter function from WHERE clause
	var filter func(storage.Row) bool
	if stmt.Where != nil {
		filter = e.buildFilter(stmt.Where)
	}
	
	rows, columns, err := e.storage.Select(stmt.TableName, stmt.Columns, filter)
	if err != nil {
		return nil, err
	}
	
	// Convert to result format
	resultRows := make([][]interface{}, len(rows))
	for i, row := range rows {
		resultRows[i] = row.Values
	}
	
	return &Result{
		Columns:      columns,
		Rows:         resultRows,
		RowsAffected: len(rows),
	}, nil
}

func (e *Executor) executeInsert(stmt *parser.InsertStatement) (*Result, error) {
	// Convert expressions to values
	values := make([]interface{}, len(stmt.Values))
	for i, expr := range stmt.Values {
		values[i] = e.evaluateExpression(expr)
	}
	
	err := e.storage.Insert(stmt.TableName, values)
	if err != nil {
		return nil, err
	}
	
	return &Result{
		RowsAffected: 1,
		Message:      "INSERT successful",
	}, nil
}

func (e *Executor) executeCreateTable(stmt *parser.CreateTableStatement) (*Result, error) {
	columns := make([]storage.Column, len(stmt.Columns))
	for i, col := range stmt.Columns {
		columns[i] = storage.Column{
			Name: col.Name,
			Type: col.Type,
		}
	}
	
	err := e.storage.CreateTable(stmt.TableName, columns)
	if err != nil {
		return nil, err
	}
	
	return &Result{
		Message: fmt.Sprintf("Table %s created", stmt.TableName),
	}, nil
}

func (e *Executor) buildFilter(expr parser.Expression) func(storage.Row) bool {
	return func(row storage.Row) bool {
		binExpr, ok := expr.(*parser.BinaryExpression)
		if !ok {
			return true
		}
		
		// For now, simplified WHERE clause evaluation
		// Assumes first column is the one being filtered
		// In a real implementation, you'd look up column indices
		
		// Get right value
		rightValue := e.evaluateExpression(binExpr.Right)
		
		// For simplicity, assume column 0 for now
		leftValue := row.Values[0]
		
		// Compare based on operator
		switch binExpr.Operator {
		case parser.EQ:
			return fmt.Sprintf("%v", leftValue) == fmt.Sprintf("%v", rightValue)
		case parser.NEQ:
			return fmt.Sprintf("%v", leftValue) != fmt.Sprintf("%v", rightValue)
		case parser.LT:
			return compareValues(leftValue, rightValue) < 0
		case parser.GT:
			return compareValues(leftValue, rightValue) > 0
		case parser.LTE:
			return compareValues(leftValue, rightValue) <= 0
		case parser.GTE:
			return compareValues(leftValue, rightValue) >= 0
		}
		
		return true
	}
}

func (e *Executor) evaluateExpression(expr parser.Expression) interface{} {
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		val, _ := strconv.Atoi(e.Value)
		return val
	case *parser.StringLiteral:
		return e.Value
	case *parser.Identifier:
		return e.Value
	default:
		return nil
	}
}

func compareValues(left, right interface{}) int {
	// Simple numeric comparison
	leftInt, lok := left.(int)
	rightInt, rok := right.(int)
	
	if lok && rok {
		if leftInt < rightInt {
			return -1
		} else if leftInt > rightInt {
			return 1
		}
		return 0
	}
	
	// String comparison
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)
	
	if leftStr < rightStr {
		return -1
	} else if leftStr > rightStr {
		return 1
	}
	return 0
}
