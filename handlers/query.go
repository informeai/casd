package handlers

import (
	"fmt"

	"github.com/informeai/casd/dto"
	"github.com/informeai/casd/services"
	"github.com/spf13/cobra"
)

var queryName string

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&queryName, "name", "n", "", "Name the file for querying")
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query file information by name",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := services.NewStorage()
		if err != nil {
			fmt.Println("Error initializing storage:", err)
			return
		}
		defer storage.Close()

		if queryName == "" {
			fmt.Println("Please provide a file name using the --name or -n flag.")
			return
		}

		formulas, err := storage.FilterFormulas(func(f *dto.Formula) bool {
			return f.Name == queryName
		})
		if err != nil {
			fmt.Println("Error querying formulas:", err)
			return
		}

		if len(formulas) == 0 {
			fmt.Println("No formulas found with the given name.")
			return
		}

		for identifier, formula := range formulas {
			fmt.Printf("FORMULAS:\n\nidentifier: %s\nname: %s\ncreated_at: %v\n", identifier, formula.Name, formula.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}
