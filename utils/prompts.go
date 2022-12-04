package utils

import (
	"encoding/json"
	"fmt"
	"log"

	opts "github.com/dieple/tfrunner/options"
	"github.com/ktr0731/go-fuzzyfinder"
)

func GetProviderData() opts.Providers {
	obj := opts.Providers{}

	err := json.Unmarshal([]byte(opts.CloudProviders), &obj)

	if err != nil {
		log.Fatalln(err)
	}
	return obj
}

// user select provider, we return the src provider directory
func PromptToSelectProvider(p opts.Providers) (string, string) {
	idx, err := fuzzyfinder.Find(
		p.Providers,
		func(i int) string {
			return p.Providers[i].ID
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Cloud Provider: %s\nDescription: %s\nDirectory: %s", p.Providers[i].ID, p.Providers[i].Description, p.Providers[i].Directory)
		}))
	if err != nil {
		log.Fatal(err)
	}
	return p.Providers[idx].ID, p.Providers[idx].Directory
}

func PromptToSelectAccount(ws opts.Workspaces, provider string) int {
	idx, err := fuzzyfinder.Find(
		ws.Workspaces,
		func(i int) string {
			return ws.Workspaces[i].ID
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Cloud Provider: %s\nWorkspace: %s\nAccountId: %s\nAccountName: %s\nRegion: %s",
				provider,
				ws.Workspaces[i].ID,
				ws.Workspaces[i].AccountID,
				ws.Workspaces[i].AccountName,
				ws.Workspaces[i].Region)
		}))
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("selected: %v\n", ws.Workspaces[idx].AccountID)
	// return ws.Workspaces[idx].AccountID
	return idx
}

func PromptToSelectTFActions(tf opts.TFActions) string {
	idx, err := fuzzyfinder.Find(
		tf.TFActions,
		func(i int) string {
			return tf.TFActions[i].ID
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Terraform Action: %s\nDescription: %s",
				tf.TFActions[i].ID,
				tf.TFActions[i].Description)
		}))
	if err != nil {
		log.Fatal(err)
	}
	return tf.TFActions[idx].ID
}

func GetTFActions() opts.TFActions {
	obj := opts.TFActions{}

	err := json.Unmarshal([]byte(opts.ActionText), &obj)

	if err != nil {
		log.Fatalln(err)
	}
	return obj
}

