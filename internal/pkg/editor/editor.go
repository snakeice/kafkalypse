package editor

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Modified struct {
	Original string
	Edited   string
	HadEdits bool
}

func EditText(in, ext string) (*Modified, error) {
	tmpFile, err := os.CreateTemp(fmt.Sprintf("*.%s", ext), "kafkalypse")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(in)); err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// handle signals (e.g. ctrl+c) in the editor
	// https://stackoverflow.com/a/32566546/145400

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run editor: %w", err)
	}

	edited, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read temporary file: %w", err)
	}

	return &Modified{
		Original: in,
		Edited:   string(edited),
		HadEdits: string(edited) != in,
	}, nil
}
