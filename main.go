package main

import (
	"log"

	"github.com/mohammadne/gorillamq/cmd"
	"github.com/spf13/cobra"
)

func main() {
	const description = "Start GorillaMQ message broker"
	root := &cobra.Command{Short: description}

	root.AddCommand(
		cmd.Server{}.Command(),
	)

	if err := root.Execute(); err != nil {
		log.Fatalf("failed to execute root command: \n%v", err)
	}
}
