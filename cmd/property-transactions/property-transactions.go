package main

import (
	"github.com/spf13/cobra"
	"log"
	"property_transactions/cmd/property-transactions/app"
)

func main() {
	mainApp := &cobra.Command{
		Use:   "app",
		Short: "My CLI application",
		Long:  "This application demonstrates a Cobra CLI structure without a rootCmd.",
	}

	mainApp.AddCommand(app.ServCmd)

	if err := mainApp.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
