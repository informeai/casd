package handlers

import (
	"fmt"
	"os"

	"github.com/informeai/casd/services"
	"github.com/spf13/cobra"
)

var hashFile string

func init() {
	rootCmd.AddCommand(rehydrateCmd)
	rehydrateCmd.Flags().StringVarP(&hashFile, "identifier", "i", "", "Identifier of the file to rehydrate")
}

var rehydrateCmd = &cobra.Command{
	Use:   "rehydrate",
	Short: "Execute rehydration process on files",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := services.NewStorage()
		if err != nil {
			fmt.Println("Error initializing storage:", err)
			return
		}
		defer storage.Close()

		if hashFile == "" {
			fmt.Println("Please provide a file identifier using the --identifier or -i flag.")
			return
		}

		formula, err := storage.GetFormula(hashFile)
		if err != nil {
			fmt.Println("Error retrieving formula:", err)
			return
		}

		if formula == nil {
			fmt.Println("No formula found with the given identifier.")
			return
		}

		f, err := os.Create(formula.Name)
		if err != nil {
			fmt.Println("Error creating rehydrated file:", err)
			return
		}
		defer f.Close()
		var total int64 = 0
		for _, hash := range formula.Hashs {
			chunk, err := storage.GetChunk(hash)
			if err != nil {
				fmt.Printf("Error retrieving chunk %s: %v\n", hash, err)
				return
			}
			total += int64(len(chunk))
			fmt.Printf("\r%d", total)

		}

		var processed int64 = 0

		for _, hash := range formula.Hashs {
			chunk, err := storage.GetChunk(hash)
			if err != nil {
				fmt.Printf("Error retrieving chunk %s: %v\n", hash, err)
				return
			}
			n, err := f.Write(chunk)
			if err != nil {
				fmt.Printf("Error writing chunk %s to file: %v\n", hash, err)
				return
			}
			if n > 0 {
				processed += int64(n)
			}
			fmt.Printf("\rprocessing: %.2f%%", (float64(processed)/float64(total))*100)
		}

		outputPath := formula.Name
		fmt.Printf("\nFile rehydrated successfully: %s\n", outputPath)
	},
}
