package main

import (
	"fmt"
	"os"

	"github.com/rickcrawford/autocomplete.es/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
