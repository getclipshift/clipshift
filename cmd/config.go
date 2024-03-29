package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getclipshift/clipshift/backends"
	"github.com/nbd-wtf/go-nostr"
	"github.com/oleiade/reflections"
	"github.com/pterm/pterm"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	app *tview.Application
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Manage configuration
Use 'get' and 'set' subcommands to get/set properties
Use 'add-backend' to add a new backend
User 'edit-backend' to edit an existing backend`,
}

func writeConfig() error {
	out, err := yaml.Marshal(config)
	if err != nil {
		log.WithError(err).Println("Error converting config to yaml")
		return err
	}
	err = os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0755)
	if err != nil {
		log.WithError(err).Error("Error creating config directory")
		return err
	}
	err = os.WriteFile(viper.ConfigFileUsed(), out, 0755)
	if err != nil {
		log.WithError(err).Println("Error writing config file")
		return err
	}
	return nil
}

func addEditBackendForm(configIndex int) {
	var b *backends.BackendConfig
	add := false
	if configIndex > -1 {
		b = &config.Backends[configIndex]
	} else {
		b = &config.Backends[len(config.Backends)-1]
		add = true
	}

	// Create form
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(b.Type + " Options").SetTitleAlign(tview.AlignLeft)

	// Host address
	if add {
		switch b.Type {
		case backends.Hosts.Nostr:
			b.Host = "wss://"
		default:
			b.Host = "https://"
		}
	}
	form.AddInputField("Address:", b.Host, 40, nil, func(text string) {
		b.Host = text
	})

	// User and Pass
	switch b.Type {
	case backends.Hosts.Nostr:
		// With nostr, a private key is all we need
		pkLabel := "Private key:"
		form.AddInputField(pkLabel, b.Pass, 65, nil, func(text string) {
			b.Pass = text
		})
		form.AddButton("Generate Private Key", func() {
			generated := nostr.GeneratePrivateKey()
			pkIndex := form.GetFormItemIndex(pkLabel)
			form.GetFormItem(pkIndex).(*tview.InputField).SetText(generated)
			form.SetFocus(pkIndex)
		})
	default:
		form.AddInputField("Username:", b.User, 40, nil, func(text string) {
			b.User = text
		})
		form.AddInputField("Password:", b.Pass, 65, nil, func(text string) {
			b.Pass = text
		})
	}

	// platform-specific options
	switch b.Type {
	case backends.Hosts.Ntfy:
		form.AddInputField("Topic:", b.Topic, 65, nil, func(text string) {
			b.Topic = text
		})
		form.AddInputField("Encryption Key (optional):", b.EncryptionKey, 40, nil, func(text string) {
			b.EncryptionKey = text
		})
		form.AddButton("Generate Topic", func() {
			generated := nostr.GeneratePrivateKey()
			topicIndex := form.GetFormItemIndex("Topic:")
			form.GetFormItem(topicIndex).(*tview.InputField).SetText(generated)
			form.SetFocus(topicIndex)
		})
	}

	// Action
	var availableActions []string
	selectedAction := 0
	index := 0
	actions, _ := reflections.Items(&backends.SyncActions)
	for _, v := range actions {
		action := fmt.Sprintf("%v", v)
		availableActions = append(availableActions, action)
		if add && action == backends.SyncActions.Sync {
			selectedAction = index
		} else if action == b.Action {
			selectedAction = index
		}
		index += 1
	}
	form.AddDropDown("Action (dropdown):", availableActions, selectedAction, func(text string, i int) {
		b.Action = text
	})

	// Form buttons
	form.AddButton("Save", func() {
		app.Stop()
		err := writeConfig()
		if err != nil {
			pterm.Error.Printfln("Error writing config file: %v", err)
		} else {
			pterm.Success.Println("Config updated")
		}
	})
	if !add {
		form.AddButton("Delete", func() {
			app.Stop()
			var keepBackends []backends.BackendConfig
			for i, existing := range config.Backends {
				if i != configIndex {
					keepBackends = append(keepBackends, existing)
				}
			}
			config.Backends = keepBackends
			err := writeConfig()
			if err != nil {
				pterm.Error.Printfln("Error writing config file: %v", err)
				os.Exit(1)
			}
			pterm.Success.Println("Config updated")
		})
	}
	form.AddButton("Quit", func() {
		app.Stop()
	})

	if app.GetFocus() == nil {
		if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
			log.WithError(err).Error("Error in TUI")
			os.Exit(1)
		}
	} else {
		app.SetRoot(form, true)
	}
}

func indexToRune(i int) rune {
	var shortcut rune
	if i < 9 {
		shortcut = []rune(fmt.Sprintf("%d", i+1))[0]
	} else {
		shortcut = []rune("abcdefghijklmnopqrstuvwxyz")[i-9]
	}
	return shortcut
}
