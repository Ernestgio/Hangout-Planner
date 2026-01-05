package main

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/domain"
)

func main() {
	stmts, err := gormschema.New("mysql").Load(
		&domain.MemoryFile{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(stmts)
}
