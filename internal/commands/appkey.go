package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lemmego/api/app"
	"github.com/lemmego/api/utils"
	"github.com/spf13/cobra"
)

func init() {
	app.RegisterCommands(AppKeyCommand)
}

var AppKeyCommand = func(a app.App) *cobra.Command {
	return &cobra.Command{
		Use: "appkey",
		Run: func(cmd *cobra.Command, args []string) {
			replaceAppKey()
		},
	}
}

// readFile reads the content of a file into a byte slice
func readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// generateAppKey generates a new base64 encoded 256-bit key for AES encryption
func generateAppKey() string {
	key, err := utils.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}
	return utils.EncodeToBase64(key)
}

func replaceAppKey() {
	filePath := "./.env"
	searchString := "APP_KEY"
	replaceString := fmt.Sprintf("APP_KEY=%s", generateAppKey())

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "tempenv.txt")
	if err != nil {
		log.Printf("Error creating temporary file: %v", err)
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name()) // Clean up if something goes wrong
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)

	appKeyFound := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchString) {
			appKeyFound = true
			line = replaceString
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			log.Printf("Error writing to temporary file: %v", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error while reading file: %v", err)
		return
	}

	if err := writer.Flush(); err != nil {
		log.Printf("Error flushing writer: %v", err)
		return
	}

	if !appKeyFound {
		tempFile.Seek(0, 0)
		if _, err := tempFile.WriteString(replaceString + "\n"); err != nil {
			log.Printf("Error writing new APP_KEY: %v", err)
			return
		}
		envFile, err := readFile(filePath)
		if err != nil {
			log.Fatalf("Error reading ENV file: %v", err)
		}
		if _, err := tempFile.Write(envFile); err != nil {
			log.Printf("Error copying file content: %v", err)
			return
		}
	}

	if err := tempFile.Close(); err != nil {
		log.Printf("Error closing temporary file: %v", err)
		return
	}

	if err := os.Rename(tempFile.Name(), filePath); err != nil {
		log.Printf("Error renaming temp file: %v", err)
		return
	}

	fmt.Println("Application key [", replaceString, "] set successfully.")
}
