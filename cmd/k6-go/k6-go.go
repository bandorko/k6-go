package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

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
	customK6File := filepath.Join(tempDir, "custom-k6")
	err = builder.Build(filename, customK6File)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(customK6File)
	cmd := exec.Command(customK6File, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	/*	out, err := exec.RunCommand(context.Background(), 0, tempDir, customK6File, os.Args[1:]...)
		log.Println(out)
		log.Fatalln(err)
		if err != nil {
			log.Println(out)
			log.Fatalln(err)
		}*/
}
