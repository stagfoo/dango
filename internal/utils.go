package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func CopyToClipboard(text string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %s | pbcopy", text))
	return cmd.Run()
}
