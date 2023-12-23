/*
	Venjector: Copyright (C) 2023 tizu69

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"github.com/ncruces/zenity"
	"golang.design/x/clipboard"
)

//go:embed core/*
var core embed.FS

var cli struct {
	LocalData  bool `help:"Do not use app data directory" default:"false"`
	AutoChoice int  `help:"Which user choice to make" default:"-1"`
	Visual     bool `help:"Visualize the progress" default:"false"`
	Tipless    bool `help:"No tips" default:"false"`
}
var progress zenity.ProgressDialog
var process = 0

func main() {
	log.SetReportCaller(true)

	kong.Parse(&cli)

	log.Info("Welcome to Venjector!")

	err := clipboard.Init()
	fatalIfError("Failed to initialize clipboard", err)

	newProgress(2)
	setVal(1, "Looking for PNPM", ensurePnpm)
	setVal(2, "Looking for Git", ensureGit)

	progress.Text("Welcome to Venjector!")
	time.Sleep(1 * time.Second) // This delay is unnecessary, but here to make the message readable

	if cli.LocalData {
		d, err := zenity.Entry("Please enter your local data directory", zenity.EntryText(getConfigPath()))
		fatalIfError("Failed to get local data directory", err)

		switch runtime.GOOS {
		case "linux":
			configLnxLocal = d
		case "darwin":
			configMacLocal = d
		case "windows":
			configWinLocal = d
		default:
			log.Fatal("Unsupported OS")
		}
	}

	for i := 0; true; i++ {
		if cli.AutoChoice != -1 && i > 0 {
			break
		}

		userChoice()

		switch process {
		case 0: // rebuild
			newProgress(9)
			setVal(1, "Downloading Vencord", pullRepo)
			setVal(2, "Installing dependencies", pnpmInstall)
			setVal(3, "Copying plugins", copyOverrides)
			setVal(3, "Downloading remote plugins", downloadPlugs)
			setVal(4, "Copying VenjectorCore", copyCore)
			setVal(5, "Changing reload-time variables", reloadVars)
			setVal(6, "Running tests", pnpmTest)
			setVal(7, "Building Vencord with plugins", pnpmBuild)
			setVal(8, "Adapting Vencord", replaceDev)

			extras := ""

			userpluginLocation := filepath.Join(getConfigPath(), "overrides", "src", "userplugins")
			if e, err := isEmpty(userpluginLocation); err != nil || e {
				extras += "\n\nWARN: Plugin directory was empty, so no custom plugins were injected."
			}

			pluginLocation := filepath.Join(getConfigPath(), "overrides", "src", "plugins")
			if e, err := isEmpty(pluginLocation); !os.IsNotExist(err) && !e {
				extras += "\n\nWARN: Explicit plugin override, prefer using 'userplugins' directory for custom plugins."
			}

			if !cli.Tipless {
				zenity.Info("Reloaded plugins!\n\n" +
					"If you haven't already, use the 'Install or uninstall Venjector' option to enable your custom plugins.\n" +
					"Updating through Vencord itself should work just fine. Please create an issue if you have any problems." + extras)
			} else if extras != "" {
				zenity.Warning("Reloaded plugins!\n\n" + extras)
			}

		case 2: // inject
			newProgress(1)

			if !cli.Tipless {
				err := zenity.Question("You're about to install or uninstall Venjector. Only use this if:\n" +
					"- You reloaded plugins at least once\n" +
					"- You don't have the Venjector patch installed\n" +
					"- You got it installed, but want to uninstall it\n" +
					"- You're using the vanilla client (see 'Install Vesktop' for Vesktop info)")
				if err != nil {
					progress.Close()
					continue
				}
			}

			setVal(1, "Injecting Discord with Venjector", injecc)
			zenity.Info("All done! Restart (not just hide!) your client to apply the changes.")
		case 4: // vesktop guide
			newProgress(1)

			if !cli.Tipless {
				err := zenity.Question("You're about to install Venjector for Vesktop. Only use this if:\n"+
					"- You reloaded plugins at least once\n"+
					"- You are using the Vesktop client!!\n"+
					"- You don't have the Vesktop patch installed\n"+
					"- VESKTOP CURRENTLY ISN'T RUNNING\n\n"+
					"To manually install Venjector, open Vesktop -> Settings -> Vesktop Settings -> Vencord Location and change"+
					" to the copied location (click 'Copy location')\n\n"+
					"To uninstall Venjector, open Vesktop -> Settings -> Vesktop Settings -> Vencord Location -> Reset",
					zenity.ExtraButton("Copy location"), zenity.OKLabel("Auto-install"))
				if err == zenity.ErrExtraButton {
					setVal(1, "Copying Vesktop path", func() {
						path, err := filepath.Abs(getConfigPath())
						fatalIfError("Failed to get Vesktop path", err)
						clipboard.Write(clipboard.FmtText, []byte(filepath.Join(path, "cord", "dist")))
						time.Sleep(1 * time.Second)
					})
					continue
				} else if err != nil {
					progress.Close()
					continue
				}
			}

			setVal(1, "Injecting Vesktop with Venjector", injeccVesktop)
			zenity.Info("All done! Restart your client to apply the changes.")
		case 1: // local plugins
			newProgress(1)
			setVal(1, "Opening plugin directory", func() {
				openByPath(filepath.Join(getConfigPath(), "overrides", "src", "userplugins"))
			})
		case 3: // remote plugins
			f, err := os.OpenFile(filepath.Join(getConfigPath(), "remote.json"), os.O_CREATE|os.O_RDONLY, 0644)
			fatalIfError("Failed to open remote.json", err)

			var data []string = []string{}
			json.NewDecoder(f).Decode(&data)
			// fatalIfError("Failed to decode remote.json", err)

			for {
				for i, v := range data {
					for j, w := range data {
						if v == w && i != j {
							data = remove(data, j)
							zenity.Warning("Removed duplicate plugin (" + v + ") from remote list")
							break
						}
					}
				}

				if len(data) != 0 {
					sel, err := zenity.List("Remote plugins (Venjector)", data,
						zenity.DisallowEmpty(), zenity.Width(512), zenity.Height(512),
						zenity.ExtraButton("Remove"), zenity.OKLabel("Add plugin"), zenity.CancelLabel("Done"))
					if err == zenity.ErrExtraButton {
						for i, plugin := range data {
							if plugin == sel {
								data = remove(data, i)
								break
							}
						}
					} else if err != nil {
						break
					}
				}

				inp, err := zenity.Entry("Enter plugin URL (use the raw URL!!)", zenity.Title("Venjector"))
				if err == zenity.ErrCanceled && len(data) == 0 {
					break
				} else if err == zenity.ErrCanceled {
					continue
				}

				// check if url is valid
				b, err := http.Get(inp)
				if err != nil || b.StatusCode != 200 {
					zenity.Error("Invalid plugin URL")
					b.Body.Close()
					continue
				}

				b.Body.Close()
				data = append(data, inp)
			}

			result, err := json.Marshal(data)
			fatalIfError("Failed to marshal remote.json", err)

			err = os.WriteFile(filepath.Join(getConfigPath(), "remote.json"), result, 0644)
			fatalIfError("Failed to write remote.json", err)
		}
	}
}
