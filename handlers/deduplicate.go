package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/informeai/casd/dto"
	"github.com/informeai/casd/services"
	"github.com/informeai/casd/utils"
	"github.com/spf13/cobra"
)

var filePath string

func init() {
	rootCmd.AddCommand(deduplicatedCmd)
	deduplicatedCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file to deduplicate")
}

var deduplicatedCmd = &cobra.Command{
	Use:   "deduplicate",
	Short: "Execute deduplication process on files",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := services.NewStorage()
		if err != nil {
			fmt.Println("Error initializing storage:", err)
			return
		}
		defer storage.Close()

		if filePath == "" {
			fmt.Println("Please provide a file path using the --file or -f flag.")
			return
		}
		info, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("error stating file:", err)
			return
		}
		total := info.Size()

		f, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error open file:", err)
			return
		}

		const chunkSize = 32768
		var hashs []string

		buf := make([]byte, chunkSize)
		var processed int64 = 0
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err.Error() != "EOF" {
					fmt.Println("Error reading file chunk:", err)
				}
				break
			}
			if n > 0 {
				processed += int64(n)
			}
			fmt.Printf("\rprocessing: %.2f%%", (float64(processed)/float64(total))*100)
			chunk := buf[:n]
			hash := utils.HashData(chunk)
			err = storage.PutChunk(hash, chunk)
			if err != nil {
				fmt.Printf("Error processing chunk: %v\n", err)
				return
			}
			hashs = append(hashs, hash)
		}
		defer f.Close()

		baseName := filepath.Base(filePath)
		key := strings.Join(hashs, ":")
		hashKey := utils.HashData([]byte(key))

		if err := storage.SaveFormula(hashKey, &dto.Formula{Name: baseName, Hashs: hashs, CreatedAt: time.Now()}); err != nil {
			fmt.Println("Error saving formula:", err)
			return
		}

		fmt.Printf("\nidentifier: %s\n", hashKey)
	},
}
