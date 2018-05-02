package tests

import (
	"fmt"
	"github.com/davyxu/protoplus/codegen"
	"os/exec"
	"strings"
)

func compileLauncher(launcherFile, configFile, tableFile string) error {

	m := struct {
		ConfigFile string
	}{
		ConfigFile: strings.Replace(configFile, "\\", "\\\\", -1),
	}

	const textTemplate = `
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	data, err := ioutil.ReadFile("{{.ConfigFile}}")
	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config
	err = json.Unmarshal(data, &config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	outData, err := json.MarshalIndent(&config, "", "\t")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println(string(outData))
}
`
	var data []byte
	err := codegen.NewCodeGen("launcher").
		RegisterTemplateFunc(codegen.UsefulFunc).
		ParseTemplate(textTemplate, m).
		WriteOutputFile(launcherFile).Error()

	if err != nil {
		fmt.Println(string(data))
		return err
	}

	cmd := exec.Command("go", "run", launcherFile, tableFile)

	output, err := cmd.CombinedOutput()

	fmt.Println("go launcher:", string(output))

	if err != nil {

		return err
	}
	return nil
}
