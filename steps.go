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
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ncruces/zenity"
	cp "github.com/otiai10/copy"
)

func ensurePnpm() {
	command := exec.Command("pnpm", "-v")

	buf := new(bytes.Buffer)
	command.Stdout = buf

	err := command.Run()
	fatalIfError("Failed to run PNPM, is it installed and in your PATH?", err)

	ver := strings.ReplaceAll(buf.String(), "\n", "")
	log.Info("Found PNPM", "version", ver)
}

func ensureGit() {
	command := exec.Command("git", "--version")

	buf := new(bytes.Buffer)
	command.Stdout = buf

	err := command.Run()
	fatalIfError("Failed to run Git, is it installed and in your PATH?", err)

	ver := strings.ReplaceAll(buf.String(), "\n", "")
	log.Info("Found Git", "version", ver)
}

func userChoice() {
	if cli.AutoChoice != -1 {
		log.Info("Auto choice", "choice", cli.AutoChoice)
		process = cli.AutoChoice
		return
	}

	rootLocation := getConfigPath()
	repoLocation := filepath.Join(rootLocation, "cord")
	if _, err := os.Stat(rootLocation); os.IsNotExist(err) {
		log.Info("No root directory, downloading Vencord repo")
		process = 0
		return
	} // I'd like to combine these, but go can't do that :(((((
	if _, err := os.Stat(repoLocation); os.IsNotExist(err) {
		log.Info("No cord directory, downloading Vencord repo")
		process = 0
		return
	}

	const (
		choiceRebuild = "Reload plugins"
		choiceOpen    = "Open plugin directory"
		choiceInject  = "Install or uninstall Venjector"
		choiceOpenWeb = "Manage downloaded plugins"
		choiceUpdate  = "Update Vencord"
		choiceVesktop = "Install Vesktop"
		choiceAbout   = "About Venjector"
	)

	result, err := zenity.List("Welcome to Venjector, the plugin loader for the cutest client mod :3\nWhat do you wish to do today?",
		[]string{choiceRebuild, choiceUpdate, choiceOpen, choiceOpenWeb, choiceInject, choiceVesktop, choiceAbout},
		zenity.Title("Venjector"), zenity.DisallowEmpty(), zenity.CancelLabel("Quit"))

	switch err {
	case zenity.ErrCanceled:
		newProgress(1)
		progress.Text("Have a nice day! :3")
		time.Sleep(1 * time.Second) // This delay is unnecessary, but here to make the message readable
		progress.Close()

		log.Fatal("Canceled by user")
	default:
		fatalIfError("User select error", err)
	}

	log.Info("Selected", "option", result)
	switch result {
	case choiceRebuild:
		process = 0
	case choiceOpen:
		process = 1
	case choiceInject:
		process = 2
	case choiceOpenWeb:
		process = 3
	case choiceUpdate: // same thing, just different name
		process = 0
	case choiceVesktop:
		process = 4
	case choiceAbout:
		zenity.Info(`Thanks for using Venjector!

Venjector is a plugin loader for the cutest client mod :3

This project is made possible thanks to the following awesome libraries: <3
	github.com/alecthomas/kong v0.8.1
	github.com/charmbracelet/log v0.3.1
	github.com/ncruces/zenity v0.10.10
	github.com/otiai10/copy v1.14.0
	golang.design/x/clipboard v0.7.0
	
Venjector: Copyright (C) 2023 tizu69
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain conditions.
Venjector is licensed under the following license: GNU GPL (version 3)`)
		process = -2
	}
}

func pullRepo() {
	log.Info("Pulling Vencord repo")
	repoLocation := filepath.Join(getConfigPath(), "cord")

	if _, err := os.Stat(repoLocation); err == nil {
		log.Info("Deleting old Vencord repo, YOLO", "location", repoLocation)
		os.RemoveAll(repoLocation)

		/* command := exec.Command("git", "pull")
		command.Dir = repoLocation

		buf := new(bytes.Buffer)
		command.Stderr = buf

		err := command.Run()
		log.Info("Ran Git pull", "output", buf.String())

		fatalIfError("Failed to run Git", err)
		log.Info("Successfully pulled Vencord repo")
		return */
	}

	abs, err := filepath.Abs(repoLocation)
	fatalIfError("Failed to get absolute path", err)

	command := exec.Command("git", "clone", repo, abs)

	buf := new(bytes.Buffer)
	command.Stderr = buf

	err = command.Run()
	log.Info("Ran Git clone", "output", buf.String())

	fatalIfError("Failed to run Git", err)
	log.Info("Successfully pulled Vencord repo")

}

func pnpmInstall() {
	log.Info("Installing dependencies for Vencord")
	repoLocation := filepath.Join(getConfigPath(), "cord")

	command := exec.Command("pnpm", "install", "--frozen-lockfile")
	command.Dir = repoLocation

	buf := new(bytes.Buffer)
	command.Stderr = buf

	err := command.Run()
	log.Info("Ran PNPM install", "output", buf.String())

	fatalIfError("Failed to run PNPM", err)
	log.Info("Successfully installed dependencies for Vencord")
}

func copyOverrides() {
	log.Info("Copying overrides")
	repoLocation := filepath.Join(getConfigPath(), "cord")
	targetLocation := filepath.Join(repoLocation)
	overridesLocation := filepath.Join(getConfigPath(), "overrides")
	pluginLocation := filepath.Join(overridesLocation, "src", "userplugins")

	if _, err := os.Stat(pluginLocation); err != nil {
		log.Info("Creating plugin directory", "location", pluginLocation)
		os.MkdirAll(pluginLocation, 0755)
	}

	log.Info("Copying overrides recursively", "from", overridesLocation, "to", targetLocation)

	err := cp.Copy(overridesLocation, targetLocation)
	fatalIfError("Failed to copy overrides", err)

	log.Info("Successfully copied overrides")
}

func downloadPlugs() {
	log.Info("Downloading remote plugins")
	pluginLocation := filepath.Join(getConfigPath(), "cord", "src", "userplugins")

	f, err := os.OpenFile(filepath.Join(getConfigPath(), "remote.json"), os.O_CREATE|os.O_RDONLY, 0644)
	fatalIfError("Failed to open remote.json", err)

	var data []string = []string{}
	json.NewDecoder(f).Decode(&data)
	// fatalIfError("Failed to decode remote.json", err)

	for i, v := range data {
		for j, w := range data {
			if v == w && i != j {
				data = remove(data, j)
				zenity.Warning("Removed duplicate plugin (" + v + ") from remote list")
				break
			}
		}
	}

	for i, v := range data {
		log.Info("Downloading remote plugin", "plugin", v)
		downloadFile(filepath.Join(pluginLocation, "remotePlugin"+intToLetters(int32(i)), "index.tsx"), v)
	}

	log.Info("Successfully downloaded remote plugins")
}

func copyCore() {
	log.Info("Copying core")
	repoLocation := filepath.Join(getConfigPath(), "cord")
	targetLocation := filepath.Join(repoLocation, "src")
	pluginLocation := core

	log.Info("Copying core recursively", "from", pluginLocation, "to", targetLocation)

	os.MkdirAll(filepath.Join(targetLocation, "userplugins", "core"), 0755)

	files, err := getAllFilenames(&core)
	fatalIfError("Failed to get files in core", err)

	log.Info("Found some files in core", "files", files)

	for from, to := range map[string]string{
		"plugin.tsx":      "userplugins/core/index.tsx",
		"pluginNative.ts": "userplugins/core/native.ts",
		"tabUpdater.tsx":  "components/VencordSettings/UpdaterTab.tsx",
		"tabPlugins.tsx":  "components/PluginSettings/index.tsx",
	} {
		fileContent, err := core.ReadFile(filepath.Join("core", from))
		fatalIfError("Failed to read core file", err)

		if cli.Visual {
			progress.Text("Copying core: " + from)
		}

		err = os.WriteFile(filepath.Join(targetLocation, to), fileContent, 0644)
		fatalIfError("Failed to write core file", err)
	}

	log.Info("Successfully copied core plugins")
}

func reloadVars() {
	log.Info("Inserting reload-time vars")
	pluginLocation := filepath.Join(getConfigPath(), "cord", "src", "userplugins")

	selfPath, err := os.Executable()
	fatalIfError("Failed to get self path", err)

	targets := map[string]string{
		"$VENJECTOR-SELFPATH": selfPath,
	}

	err = filepath.Walk(pluginLocation, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || err != nil {
			return nil
		}
		log.Info("Inserting reload-time vars", "file", path)

		if cli.Visual {
			progress.Text("Inserting reload-time vars: " + path)
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for k, v := range targets {
			contents = bytes.ReplaceAll(contents, []byte(k), []byte(v))
		}

		return os.WriteFile(path, contents, 0644)
	})
	fatalIfError("Failed to insert reload-time vars", err)

	log.Info("Successfully inserted reload-time vars")
}

func pnpmTest() {
	log.Info("Running tests")
	repoLocation := filepath.Join(getConfigPath(), "cord")

	command := exec.Command("pnpm", "test")
	command.Dir = repoLocation

	buf := new(bytes.Buffer)
	command.Stderr = buf

	/* err := */
	command.Run()
	log.Info("Ran PNPM test", "output", buf.String())

	/* TODO: if _, ok := err.(*exec.ExitError); ok {
		err := zenity.Error("Tests did not pass. If you feel experimental, click 'Continue anyway' to ignore test results.",
			zenity.ExtraButton("Continue anyway"))
		if err == zenity.ErrExtraButton {
			log.Info("Continuing without tests")
			return
		}

		log.Fatal("Tests did not pass", "err", err)
	} */

	log.Info("Successfully ran tests")
}

func pnpmBuild() {
	log.Info("Building Vencord with plugins")
	repoLocation := filepath.Join(getConfigPath(), "cord")

	command := exec.Command("pnpm", "build")
	command.Dir = repoLocation

	buf := new(bytes.Buffer)
	command.Stderr = buf

	err := command.Run()
	log.Info("Ran PNPM build", "output", buf.String())

	fatalIfError("Failed to run PNPM", err)
	log.Info("Successfully built Vencord with plugins")
}

func replaceDev() {
	log.Info("Turning Vencord into production")
	repoLocation := filepath.Join(getConfigPath(), "cord")

	f, err := os.OpenFile(filepath.Join(repoLocation, "scripts", "runInstaller.mjs"), os.O_RDWR, 0644)
	fatalIfError("Failed to open runInstaller.mjs", err)

	lines, err := io.ReadAll(f)
	fatalIfError("Failed to read runInstaller.mjs", err)

	// lines = []byte(strings.Replace(string(lines), `VENCORD_USER_DATA_DIR: BASE_DIR,`, "", 1))
	lines = []byte(strings.Replace(string(lines), `VENCORD_DEV_INSTALL: "1"`, "", 1))

	err = os.WriteFile(filepath.Join(repoLocation, "scripts", "runInstaller.mjs"), lines, 0644)
	fatalIfError("Failed to write runInstaller.mjs", err)
}

func injecc() {
	log.Info("Injecting Vencord with Venjector")
	repoLocation := filepath.Join(getConfigPath(), "cord")

	command := exec.Command("pnpm", "inject")
	command.Dir = repoLocation

	buf := new(bytes.Buffer)
	command.Stderr = buf

	err := command.Run()
	log.Info("Ran PNPM inject", "output", buf.String())

	fatalIfError("Failed to run PNPM", err)
	log.Info("Successfully injected Vencord with Venjector")
}

func injeccVesktop() {
	log.Info("Injecting Vencord with Venjector")
	repoLocation := filepath.Join(getConfigPath(), "cord", "dist")
	vesktopLocation := getVesktopPath()

	data, err := os.ReadFile(filepath.Join(vesktopLocation, "settings.json"))
	fatalIfError("Failed to read settings.json", err)

	var objmap map[string]interface{}
	err = json.Unmarshal(data, &objmap)
	fatalIfError("Failed to unmarshal settings.json", err)

	objmap["vencordDir"] = repoLocation

	data, err = json.Marshal(objmap)
	fatalIfError("Failed to marshal settings.json", err)

	err = os.WriteFile(filepath.Join(vesktopLocation, "settings.json"), data, 0644)
	fatalIfError("Failed to write settings.json", err)

	log.Info("Successfully injected Vencord with Venjector")
}
