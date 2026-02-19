package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	
	"github.com/noorimat/distributed-sql-database/pkg/sql"
	"github.com/noorimat/distributed-sql-database/pkg/sql/parser"
	"github.com/noorimat/distributed-sql-database/pkg/storage"
)

func main() {
	fmt.Println("NuvaDB Client v0.1.0")
	fmt.Println("Type 'exit' or 'quit' to exit")
	fmt.Println()
	
	// Initialize storage and executor
	store := storage.NewMemoryStorage()
	executor := sql.NewExecutor(store)
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("nuvadb> ")
		
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		
		if input == "" {
			continue
		}
		
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Goodbye!")
			break
		}
		
		// Parse and execute SQL
		l := parser.NewLexer(input)
		p := parser.NewParser(l)
		stmt := p.ParseStatement()
		
		if stmt == nil {
			fmt.Println("Error: Failed to parse SQL statement")
			if len(p.Errors()) > 0 {
				for _, err := range p.Errors() {
					fmt.Printf("  %s\n", err)
				}
			}
			continue
		}
		
		result, err := executor.Execute(stmt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		
		// Display results
		if result.Message != "" {
			fmt.Println(result.Message)
		}
		
		if len(result.Rows) > 0 {
			// Print column headers
			fmt.Println()
			for i, col := range result.Columns {
				fmt.Printf("%-20s", col)
				if i < len(result.Columns)-1 {
					fmt.Print(" | ")
				}
			}
			fmt.Println()
			fmt.Println(strings.Repeat("-", len(result.Columns)*23))
			
			// Print rows
			for _, row := range result.Rows {
				for i, val := range row {
					fmt.Printf("%-20v", val)
					if i < len(row)-1 {
						fmt.Print(" | ")
					}
				}
				fmt.Println()
			}
			fmt.Printf("\n%d row(s) returned\n", len(result.Rows))
		} else if result.RowsAffected > 0 {
			fmt.Printf("%d row(s) affected\n", result.RowsAffected)
		}
		
		fmt.Println()
	}
}
