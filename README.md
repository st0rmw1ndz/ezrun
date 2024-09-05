# ezrun

ezrun sits in between a stdin-based menu system like dmenu, and allows you to easily run command shortcuts by defining a shell function.

## Example Configuration

```sh
#!/bin/sh

# shellcheck disable=SC2034
CHOICE_COMMAND="dmenu -l 4 -g 8"

st() { command st; }
nvim() { command st -e nvim; }
pavucontrol() { command pavucontrol; }
signal() { signal-desktop; }
firefox() { command firefox; }
flycast() { flatpak run org.flycast.Flycast; }
steam() { flatpak run com.valvesoftware.Steam; }
```

*$XDG_CONFIG_HOME/ezrun/ezrunrc*
