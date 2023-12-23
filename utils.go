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
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ncruces/zenity"
)

const (
	repo = "https://github.com/Vendicated/Vencord"
)

func fatalIfError(task string, err error) {
	if err != nil {
		zenity.Error(fmt.Sprintf("ERROR: %s - %s", task, err.Error()))
		log.Fatal(task, "err", err)
	}
}

func setVal(val int, task string, run func()) {
	progress.Value(val)
	progress.Text(task + "..")
	run()
	time.Sleep(250 * time.Millisecond) // idk why, but without a delay this crashed Zenity on my end.
	if progress.MaxValue() == val+1 {
		time.Sleep(750 * time.Millisecond) // let the user read, duh :3
		progress.Close()
	}
}

var (
	configLnx = filepath.Join(os.Getenv("HOME"), ".config/Venjector")
	configMac = filepath.Join(os.Getenv("HOME"), "Library/Application Support/Venjector")
	configWin = filepath.Join(os.Getenv("LOCALAPPDATA"), "Venjector")
)

var (
	configLnxLocal = "./venjectorConfig"
	configMacLocal = "./venjectorConfig"
	configWinLocal = ".\\venjectorConfig"
)

func getConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		if cli.LocalData {
			return os.ExpandEnv(configLnxLocal)
		}
		return os.ExpandEnv(configLnx)
	case "darwin":
		if cli.LocalData {
			return os.ExpandEnv(configMacLocal)
		}
		return os.ExpandEnv(configMac)
	case "windows":
		if cli.LocalData {
			return os.ExpandEnv(configWinLocal)
		}
		return os.ExpandEnv(configWin)
	}

	log.Fatal("Unsupported OS")
	return ""
}

var (
	vesktopConfigLnx = filepath.Join(os.Getenv("HOME"), ".config/VencordDesktop/VencordDesktop")
	vesktopConfigMac = filepath.Join(os.Getenv("HOME"), "Library/Application Support/VencordDesktop/VencordDesktop")
	vesktopConfigWin = filepath.Join(os.Getenv("LOCALAPPDATA"), "VencordDesktop\\VencordDesktop")
)

func getVesktopPath() string {

	switch runtime.GOOS {
	case "linux":
		return os.ExpandEnv(vesktopConfigLnx)
	case "darwin":
		return os.ExpandEnv(vesktopConfigMac)
	case "windows":
		return os.ExpandEnv(vesktopConfigWin)
	}

	log.Fatal("Unsupported OS")
	return ""
}

// https://stackoverflow.com/a/30708914
func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func openByPath(path string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", path).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	case "darwin":
		err = exec.Command("open", path).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	fatalIfError("Failed to open "+path, err)
}

func newProgress(max int) {
	p, err := zenity.Progress(
		zenity.Title("Venjector"),
		zenity.AutoClose(),
		zenity.NoCancel(),
		zenity.MaxValue(max+1),
		zenity.TimeRemaining(),
	)
	fatalIfError("Failed to open the GUI, install one of 'zenity, matedialog, qarma' on Linux or 'osascript' on macOS, then try again", err)
	time.Sleep(1 * time.Second) // await a fade animation, if present
	progress = p
}

// https://gist.github.com/clarkmcc/1fdab4472283bb68464d066d6b4169bc?permalink_comment_id=4405804#gistcomment-4405804
func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

// https://stackoverflow.com/a/37335777
func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func downloadFile(path string, url string) {
	// Make the directory
	err := os.MkdirAll(filepath.Dir(path), 0755)
	fatalIfError("Failed to create "+filepath.Dir(path), err)

	// Create the file
	out, err := os.Create(path)
	fatalIfError("Failed to create "+path, err)
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	fatalIfError("Failed to download "+url, err)
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	fatalIfError("Failed to write "+path, err)
}

// https://stackoverflow.com/a/66172278
func intToLetters(number int32) (letters string) {
	number--
	if firstLetter := number / 26; firstLetter > 0 {
		letters += intToLetters(firstLetter)
		letters += string('A' + number%26)
	} else {
		letters += string('A' + number)
	}

	return
}
