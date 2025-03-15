package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type editorFinishedMsg struct {
	tempJsonPath string
	err          error
}

func editFileExecCmd(filepath string) *exec.Cmd {
	return editorFileExecCmd(filepath, false)
}

func editorFileExecCmd(filepath string, readOnly bool) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	var c *exec.Cmd
	if readOnly && (editor == "vim" || editor == "nvim") {
		c = exec.Command(editor, "-R", filepath)
	} else {
		c = exec.Command(editor, filepath)
	}
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	return c
}

func getTempJsonPath() (string, error) {
	f, err := os.CreateTemp("", fmt.Sprintf("tmp_jqc%s_*.json", time.Now().Format("2006-01-02T15:04:05")))
	f.Close() // Just get path
	return f.Name(), err
}

func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create or truncate the destination file
	// The os.O_TRUNC flag ensures any existing file is truncated (overwritten)
	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
