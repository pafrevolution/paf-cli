package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var TARGET_PATH string

var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "Move a selected directory to Octobook",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		TARGET_PATH = os.Getenv("DEFAULT_TARGET_PATH")
		if TARGET_PATH == "" {
			log.Fatalf("TARGET_PATH is not set in .env file")
		}

		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current directory: %v", err)
		}

		folders, err := getFoldersInDirectory(currentDir)
		if err != nil {
			log.Fatalf("Error retrieving folders: %v", err)
		}

		if len(folders) == 0 {
			fmt.Println("No folders available to move.")
			return
		}

		var selectedFolders []string
		prompt := &survey.MultiSelect{
			Message: "Select the folder to move:",
			Options: folders,
		}
		err = survey.AskOne(prompt, &selectedFolders)
		if err != nil {
			log.Fatalf("Error selecting folder: %v", err)
		}

		for _, selectedFolder := range selectedFolders {
			sourcePath := filepath.Join(currentDir, selectedFolder)
			destPath, err := createDateDirectory(TARGET_PATH, sourcePath)
			if err != nil {
				log.Fatalf("Error creating destination directory: %v", err)
			}
			err = moveFolderWithProgress(sourcePath, destPath)
			if err != nil {
				log.Fatalf("Error moving folder: %v", err)
			}

			fmt.Printf("Folder %s successfully moved!", selectedFolder)
		}
	},
}

func getFoldersInDirectory(dir string) ([]string, error) {
	var folders []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}
	return folders, nil
}

func moveFolderWithProgress(sourcePath, destPath string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		return err
	}

	totalSize, err := calculateTotalSize(sourcePath)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions64(
		totalSize,
		progressbar.OptionSetDescription("Moving files..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("\nMove completed successfully!")
		}),
	)

	// Walk through the source directory
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path
		relativePath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		// Create destination path
		destFilePath := filepath.Join(destPath, relativePath)

		// If it's a directory, create it
		if info.IsDir() {
			return os.MkdirAll(destFilePath, os.ModePerm)
		}

		// Copy file
		if err := copyFileWithProgress(path, destFilePath, bar); err != nil {
			return err
		}

		// Verify file size and remove source file if sizes match
		sourceInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		destInfo, err := os.Stat(destFilePath)
		if err != nil {
			return err
		}

		if sourceInfo.Size() == destInfo.Size() {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove source file %s: %v", path, err)
			}
		} else {
			return fmt.Errorf("file size mismatch for %s", path)
		}

		return nil
	})
}

func copyFileWithProgress(sourcePath, destPath string, bar *progressbar.ProgressBar) error {
	// Open source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy file with progress
	_, err = io.Copy(destFile, io.TeeReader(sourceFile, bar))
	return err
}

func calculateTotalSize(path string) (int64, error) {
	var totalSize int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	return totalSize, err
}

func createDateDirectory(basePath string, srcPath string) (string, error) {
	// Get current date in 'YYYY-MM-DD' format
	currentDate := time.Now().Format("2006-01-02")
	splitSrcPath := strings.Split(srcPath, "/")
	srcFolder := splitSrcPath[len(splitSrcPath)-1]

	fullPath := filepath.Join(basePath, currentDate)
	fullPath = filepath.Join(fullPath, srcFolder)

	// Check if directory exists
	_, err := os.Stat(fullPath)

	// If directory does not exist, create it
	if os.IsNotExist(err) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
		fmt.Printf("Created directory: %s\n", fullPath)
	} else if err != nil {
		// Handle other potential errors
		return "", fmt.Errorf("error checking directory: %v", err)
	} else {
		fmt.Printf("Directory already exists: %s\n", fullPath)
	}

	return fullPath, nil
}

func init() {
	// Add the mv command to the root command
	rootCmd.AddCommand(mvCmd)
}
