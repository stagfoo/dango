package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func GetCurrentFolder() (string, error) {
	cmd := exec.Command("pwd")
	pwdBytes, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running pwd command: %v", err)
	}

	// Convert the byte slice to a string and trim any newline characters
	pwd := strings.TrimSpace(string(pwdBytes))

	return pwd, nil
}

func ListFilesInFolder() ([]string, error) {
	pwd, pwdErr := GetCurrentFolder()
	if pwdErr != nil {
		return nil, pwdErr
	}

	// Read the directory contents
	files, err := os.ReadDir(pwd)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, pwd+"/"+file.Name())
		}
	}
	return fileNames, nil
}

func RawListDirectoriesAndFiles() error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Use filepath.Glob to get all files and directories in the current directory
	files, err := filepath.Glob(filepath.Join(currentDir, "*"))
	if err != nil {
		return err
	}

	// Print each file/directory path to stdout
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}

func ListDirectoryContents(shouldPrint bool) ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %v", err)
	}

	files, err := filepath.Glob(filepath.Join(currentDir, "*"))
	if err != nil {
		return nil, fmt.Errorf("error using filepath.Glob: %v", err)
	}

	if shouldPrint {
		for _, file := range files {
			fmt.Println(file)
		}
		return nil, nil // Returning nil slice and nil error when printing
	} else {
		return files, nil
	}
}

func OutputSelectedItems() {
	db := ViewDB(Path)
	for _, item := range db.Items {
		fmt.Println(item)
	}
}

func RemoveAllItems() {
	err := os.Remove(Path)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error clearing dango file: %v\n", err)
		return
	}
	fmt.Println("Dango database has been removed")
}

func CopyToClipboard(text string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %s | pbcopy", text))
	return cmd.Run()
}

func ViewDBItems(program *tea.Program) []string {

	_, runError := program.Run()
	if runError != nil {
		fmt.Printf("Error running program: %s\n", runError)
	}

	return ViewDB(Path).Items
}

func HasInputFromPipe() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	result := (fileInfo.Mode() & os.ModeCharDevice) == 0
	if !result {
		os.Stdin.Close()
	}
	return result
}
