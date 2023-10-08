package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/bandorko/k6-go/internal/builder"
	"github.com/bandorko/k6-go/internal/fs"
)

const (
	cmdHelp  = "help"
	cmdBuild = "build"
)

const usage = `
Usage:
k6-go build --output <custom_k6_binary_file> <go_script_file>
or
k6-go run <go_script_file>

Avaliable Commands:
  build			Only build custom k6 binary with the needed extensions
  help			Print help
  *			Any other command from k6

Flags:
  --output file		custom k6 binary output file
`

type Config struct {
	subcommand string
	output     string
	k6Args     []string
	script     string
}

func main() {
	config, err := readConfig(os.Args)
	if err != nil || config.subcommand == cmdHelp {
		printHelp(err)
	}

	tempDir, err := fs.CreateTempDir("")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(tempDir)
	customK6File := filepath.Join(tempDir, getCustomK6FileName())
	if config.output != "" {
		customK6File = config.output
	}
	err = builder.Build(config.script, customK6File, config.subcommand == cmdHelp)
	if err != nil {
		log.Fatalln(err)
	}
	if config.subcommand != cmdBuild { // not just build
		cmd := exec.Command(customK6File, config.k6Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func readConfig(args []string) (*Config, error) {
	config := &Config{
		subcommand: cmdHelp,
		k6Args:     make([]string, 0),
	}
	argn := len(args)

	if argn < 2 {
		return config, fmt.Errorf("Not enough parameter")
	}

	for i := 1; i < argn; i++ {
		arg := args[i]
		if arg == "--output" {
			i++
			if i >= argn {
				return config, fmt.Errorf("flag needs an argument: --output")
			}
			config.output = args[i]
			continue
		}
		config.k6Args = append(config.k6Args, arg)
	}
	config.subcommand = args[1]
	if argn > 2 {
		config.script = args[argn-1]
	}
	return config, nil
}

func printHelp(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(usage)
}

func getCustomK6FileName() string {
	if runtime.GOOS == "windows" {
		return "k6.exe"
	}
	return "k6"
}
