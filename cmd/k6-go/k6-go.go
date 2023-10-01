package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/bandorko/k6-go/internal/builder"
	"github.com/bandorko/k6-go/internal/fs"
)

func main() {
	filename := os.Args[len(os.Args)-1]
	tempDir, err := fs.CreateTempDir("")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(tempDir)
	customK6File := filepath.Join(tempDir, getCustomK6FileName())
	err = builder.Build(filename, customK6File)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(customK6File)
	cmd := exec.Command(customK6File, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func getCustomK6FileName() string {
	if runtime.GOOS == "windows" {
		return "custom-k6.exe"
	}
	return "custom-k6"
}
