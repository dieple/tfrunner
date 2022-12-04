package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	opts "github.com/dieple/tfrunner/options"
	"github.com/ktr0731/go-fuzzyfinder"
)


// Walk thru ./infra/terraform<provider>/main ... subdirectory to
// find "main.tf" files for module name(s)
// then return full path of slice of module names
func FindModules(searchDirectory string) ([]string, error) {

	// fmt.Printf("searchDirectory: %v\n", searchDirectory)
	files, err := filepath.Glob(searchDirectory + "/**/**/" + opts.TerraformMainFile)
	var buf []string

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, file := range files {
		//fmt.Println(file)
		// ignoring these directories
		if strings.Contains(file, "/initial/") || strings.Contains(file, ".terraform") || strings.Contains(file, "/abstract/") {
			continue
		}
		moduleName := strings.TrimSuffix(file, "/" + opts.TerraformMainFile) + "\n"
		buf = append(buf, moduleName)
	}
	// fmt.Println(buf)
	return buf, err
}

func PromptToSelectModules(listModuleNames []string) []string {
	var modulesSelected []string

	idx, err := fuzzyfinder.FindMulti(
		listModuleNames,
		func(i int) string {
			return listModuleNames[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Select module to build:\n\nModule Name: %s\nModule Directory: %s", listModuleNames[i][strings.LastIndex(listModuleNames[i], "/")+1:], listModuleNames[i])
		}))
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range idx {
		modulesSelected = (append(modulesSelected, listModuleNames[i]))
	}

	return modulesSelected
}
