package main

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
)

func main() {
	_, err := gormschema.New("mysql").Load(&domain.User{}, &domain.Hangout{}, &domain.Activity{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
}
