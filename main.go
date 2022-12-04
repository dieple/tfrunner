package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	opts "github.com/dieple/tfrunner/options"
	"github.com/dieple/tfrunner/utils"
	"github.com/jessevdk/go-flags"
)

// func isValidCloudProvider(provider string) bool {
// 	switch provider {
// 	case
// 		"aws",
// 		"azure",
// 		"cgp":
// 		return true
// 	}
// 	return false
// }

func processArguments() opts.Options {
	var opts opts.Options

	p := flags.NewParser(&opts, flags.Default&^flags.HelpFlag)
	_, err := p.Parse()

	if err != nil {
		fmt.Printf("fail to parse args: %v", err)
		os.Exit(1)
	}
	if opts.Help {
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	}
	if opts.Verbose != nil {
		fmt.Printf("Argument values: %v\n", opts)
	}

	return opts
}

func fetchWorkspaceDataFromJson() opts.Workspaces {
	wsFile, err := os.OpenFile("workspaces.json", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer wsFile.Close()
	byteValue, _ := ioutil.ReadAll(wsFile)
	var ws opts.Workspaces
	json.Unmarshal(byteValue, &ws)
	return ws
}

func main() {
	args := processArguments()

	if args.Interactive {
		// run terraform interactively
		wsData := fetchWorkspaceDataFromJson()

		providerSelected, providerDirectorySelected := utils.PromptToSelectProvider(utils.GetProviderData())
		//fmt.Printf("providerSelected: %v, providerDirectorySelected %v\n", providerSelected, providerDirectorySelected)

		idxSelected := utils.PromptToSelectAccount(wsData, providerSelected)
		// fmt.Printf("accountSelected: %v\n", idxSelected)
		// fmt.Println("Account details:")
		// fmt.Println(wsData.Workspaces[idxSelected])

		buf, err := utils.FindModules(fmt.Sprintf("%s/%s", providerDirectorySelected, args.ComponentDir))
		if err != nil {
            log.Fatalln(err)
		}

		// fmt.Println(buf)
		//provider string, terraformMainFile string, searchDirectory string) ([]string, error
		modulesSelected := utils.PromptToSelectModules(buf)
		// fmt.Printf("Module(s) selected %v\n", modulesSelected)

		tfData := utils.GetTFActions()
		tfActionSelected := utils.PromptToSelectTFActions(tfData)

		buildData := &opts.BuildData{
			Workspace:    wsData.Workspaces[idxSelected].ID,
			Interactive:  true,
			Bucket:       wsData.Workspaces[idxSelected].Bucket,
			BucketRegion: wsData.Workspaces[idxSelected].BucketRegion,
			Dynamodb:     wsData.Workspaces[idxSelected].DynamoDb,
			Modules:      modulesSelected,
			TFaction:     tfActionSelected,
			AutoApprove:  args.AutoApprove,
			ProviderPath: providerDirectorySelected,
		}

		// fmt.Printf("buildData: %v\n", buildData)
		utils.TFRun(args, buildData)
	} else {
		// run terraform in CI/CD or non interactive mode
        // TODO: implement this to work with your own specifics CI/CD tool
	}
}
