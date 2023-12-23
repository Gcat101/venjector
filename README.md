# Venjector

Welcome to Venjector, the plugin loader for the cutest client mod :3

To install:

- [Experimental] A GitHub Action is available [here](https://github.com/tizu69/venjector/actions),
  select the newest successful build, then download for your OS.
- Use the go install command, [which requires go](https://go.dev/doc/install)

``` sh
go install github.com/tizu69/venjector@latest # requires go
# more install methods ... soon™️
```

**Make sure the following runtime dependencies are installed:** `pnpm`, `git`

Then, run it from the command line once:

```sh
venjector # requires the go bin (~/go/bin) to be in your path
```

Venjector will automatically initialize. Then, select 'Install' or 'Install Vesktop', depending
on if you're using Vesktop or not.

On reboot of your client, Venjector will take care of the rest. 4 neat buttons will be added to
the plugins page: Reload, open folder, open list of remote, and open Venjector.

**NOTE:** I only officially support Linux. If you're on Windows or Darwin (macOS), Venjector
may work, but if it doesn't, I don't care.

Enjoy!

## Credits

Obviously, thanks to everyone who's contributed to Vencord for making this project possible.
I can't tell you how much I appreciate them.

---

Venjector uses the following awesome libraries: <3

```sh
github.com/alecthomas/kong v0.8.1
github.com/charmbracelet/log v0.3.1
github.com/ncruces/zenity v0.10.10
github.com/otiai10/copy v1.14.0
golang.design/x/clipboard v0.7.0
```

(.. and their dependencies, of course <3)

## License

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

See [LICENSE.txt](LICENSE.txt) for more information.
