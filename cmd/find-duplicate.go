package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func hashFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func findDuplicateFiles(rootDir string) map[string][]string {
	fileHashes := make(map[string][]string)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash, err := hashFile(path)
			if err == nil {
				fileHashes[hash] = append(fileHashes[hash], path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}

	duplicates := make(map[string][]string)
	for hash, paths := range fileHashes {
		if len(paths) > 1 {
			duplicates[hash] = paths
		}
	}
	return duplicates
}

func saveDuplicatesToFile(duplicates map[string][]string, outputFile string) {
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, paths := range duplicates {
		for _, path := range paths {
			file.WriteString(path + "\n")
		}
		file.WriteString("\n")
	}

	fmt.Println("Duplicate file list saved to", outputFile)
}

var findDupCmd = &cobra.Command{
	Use:   "find-dup",
	Short: "Find duplicated files in the folder and it's subfolder",
	Run: func(cmd *cobra.Command, args []string) {
		var folderToScan string
		fmt.Print("Enter the folder path to scan: ")
		fmt.Scanln(&folderToScan)
	
		if _, err := os.Stat(folderToScan); os.IsNotExist(err) {
			fmt.Println("Invalid directory.")
			return
		}
	
		duplicates := findDuplicateFiles(folderToScan)
		if len(duplicates) > 0 {
			saveDuplicatesToFile(duplicates, "duplicate_files.txt")
		} else {
			fmt.Println("No duplicate files found.")
		}
	},
}

func init() {
    // Add the greet command to the root command
    rootCmd.AddCommand(findDupCmd)
}
