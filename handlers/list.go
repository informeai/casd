package handlers

import (
	"fmt"

	"github.com/informeai/casd/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List formulas by name",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := services.NewStorage()
		if err != nil {
			fmt.Println("Error initializing storage:", err)
			return
		}
		defer storage.Close()

		formulas, err := storage.ListFormulas()
		if err != nil {
			fmt.Println("Error retrieves  formulas:", err)
			return
		}

		if len(formulas) == 0 {
			fmt.Println("No formulas found with the given name.")
			return
		}
		fmt.Printf("FORMULAS:\n")
		for identifier, formula := range formulas {
			fmt.Printf("\nidentifier: %s\nname: %s\ncreated_at: %v\n", identifier, formula.Name, formula.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}
