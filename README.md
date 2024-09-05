# ezrun

ezrun sits in between a stdin-based menu system like dmenu, and allows you to easily run command shortcuts by defining a shell function.

## Example Configuration

The following is an ezrun configuration with the accompanying output.

```sh
#!/bin/sh

# shellcheck disable=SC2034
CHOICE_COMMAND="dmenu -l 4 -g 8"
REPLACEMENTS="_:-
0: "

nvim() { command st -e nvim; }
signal() { signal-desktop; }
firefox() { command firefox; }
obs_studio() { obs-studio; }
steam() { flatpak run com.valvesoftware.Steam; }
prism0launcher() { flatpak run org.prismlauncher.PrismLauncher; }
```

*Contents of `$XDG_CONFIG_HOME/ezrun/ezrunrc`.*

```
nvim
signal
firefox
obs-studio
steam
prism launcher
```

*ezrun output.*
