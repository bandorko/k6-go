package builder

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bandorko/k6-go/internal/exec"
	"github.com/bandorko/k6-go/internal/fs"
	"go.k6.io/xk6"
	"golang.org/x/mod/modfile"
)

const exportsFileTemplate = `
package project
import (
	_ "github.com/szkiba/xk6-g0"
	"github.com/szkiba/xk6-g0/g0"
	"github.com/traefik/yaegi/interp"
	"go.k6.io/k6/js/modules"
)

var Symbols = interp.Exports{}

func exports(vu modules.VU) interp.Exports {
	return Symbols
}

func init() {
	g0.RegisterExports(exports)
}	
`

// Build builds the custom k6 binary including all the needed libraries for running the given go script.
func Build(filename string, out string, silent bool) error {

	origStrdOut := os.Stdout
	origStrdErr := os.Stderr
	origLogOutput := log.Writer()
	defer func() {
		os.Stdout = origStrdOut
		os.Stderr = origStrdErr
		log.SetOutput(origLogOutput)
	}()
	if silent {
		if null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0); err == nil {
			os.Stderr = null
			os.Stdout = null
			log.SetOutput(os.Stderr)
		}
	}
	externalImports, err := collectExternalImports(filename)
	if err != nil {
		return err
	}
	fmt.Println("external imports in the script file:", externalImports)
	tempDir := ""
	moduleName := ""
	if len(externalImports) > 0 {
		tempDir, err = fs.CreateTempDir("")
		if err != nil {
			return err
		}

		moduleDir, err := getModuleDir(filename)
		if err != nil {
			return err
		}

		moduleName, err = getModuleName(moduleDir)
		if err != nil {
			return err
		}
		err = fs.CopyDirectory(moduleDir, tempDir)
		if err != nil {
			return err
		}
		rootTempDir, err := fs.CreateTempDir(tempDir)
		if err != nil {
			return err
		}
		err = fs.MoveFilesByExtension(tempDir, rootTempDir, ".go")
		if err != nil {
			return err
		}
		err = createExportsFile(tempDir, externalImports)
		if err != nil {
			return err
		}
		err = generate(tempDir)
		if err != nil {
			return err
		}
	}
	err = build(moduleName, tempDir, out)
	if err != nil {
		return err
	}
	os.RemoveAll(tempDir)
	return nil
}

func collectExternalImports(filename string) ([]string, error) {
	fset := &token.FileSet{}
	f, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
	if err != nil {
		fmt.Println("Can not parse script file: ", filename)
		return []string{}, nil
	} else {
		externalImports := make([]string, 0)
		for _, s := range f.Imports {
			imp := s.Path.Value
			if strings.Contains(imp, ".") {
				externalImports = append(externalImports, imp)
			}
		}
		return externalImports, nil
	}
}

func getModuleDir(filename string) (string, error) {
	// script file directory
	scriptDir, scriptName := filepath.Split(filename)
	log.Println("scriptDir:", scriptDir, " scriptName:", scriptName)

	// find go.mod file
	goModFile, err := exec.RunCommand(context.Background(), 0, scriptDir, "go", "env", "GOMOD")
	if err != nil {
		return "", err
	}
	moduleDir, _ := filepath.Split(goModFile)
	return moduleDir, nil
}

func createExportsFile(dir string, externalImports []string) error {
	exportsFile, err := os.Create(filepath.Join(dir, "exports.go"))
	if err != nil {
		return err
	}
	defer exportsFile.Close()

	exportsFile.WriteString(exportsFileTemplate)
	exportsFile.WriteString("//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract -name project ")
	for _, v := range externalImports {
		exportsFile.WriteString(v + " ")
	}
	exportsFile.WriteString("\n")
	return nil
}

func generate(dir string) error {
	_, err := exec.RunCommand(context.Background(), 0, dir, "go", "get", "github.com/szkiba/xk6-g0")
	if err != nil {
		return err
	}
	_, err = exec.RunCommand(context.Background(), 0, dir, "go", "mod", "tidy")
	if err != nil {
		return err
	}
	_, err = exec.RunCommand(context.Background(), 0, dir, "go", "generate")
	if err != nil {
		return err
	}
	return nil
}

func getModuleName(dir string) (string, error) {
	modfileBytes, err := fs.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return "", err
	}
	mf, err := modfile.Parse("go.mod", modfileBytes, nil)
	if err != nil {
		return "", err
	}
	return mf.Module.Mod.Path, nil
}

func build(modulename string, moduledir string, out string) error {
	k6Version := os.Getenv("K6_VERSION")
	k6Repo := os.Getenv("XK6_K6_REPO")
	raceDetector := os.Getenv("XK6_RACE_DETECTOR") == "1"
	skipCleanup := os.Getenv("XK6_SKIP_CLEANUP") == "1"

	extensions := []xk6.Dependency{
		{
			PackagePath: "github.com/szkiba/xk6-g0",
			Version:     "latest",
		},
	}
	replacements := []xk6.Replace{}
	if modulename != "" {
		extensions = append(extensions, xk6.Dependency{
			PackagePath: modulename,
			Version:     "v0.0.0",
		})
		replacements = append(replacements, xk6.NewReplace(modulename, moduledir))
	}
	builder := xk6.Builder{
		K6Version:    k6Version,
		K6Repo:       k6Repo,
		RaceDetector: raceDetector,
		SkipCleanup:  skipCleanup,
		Extensions:   extensions,
		Replacements: replacements,
	}

	err := builder.Build(context.Background(), out)
	if err != nil {
		return err
	}
	return err
}
