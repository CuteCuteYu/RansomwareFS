package new_notepad

import (
	"os"
	"os/exec"
)

const Content = "hello world"

func NewNotepad() {
	// Create file in current directory
	filename := "notepad.txt"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Write "hello world" to file
	if _, err := file.WriteString(Content); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}

	// Open file with notepad
	cmd := exec.Command("notepad.exe", filename)
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}
