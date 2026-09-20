# AI-Dock

`aidock` is a cross-platform terminal UI for discovering and launching AI tools
(CLI tools, IDEs, and desktop apps). It uses the Charm stack: Bubble Tea,
Bubbles, and Lip Gloss.

## Run

```sh
go run ./cmd/aidock
```

Tool configuration is stored at `~/.config/aidock/tools.json`. The first run
creates this file with an empty tools array. There is no autodetection — tools
must be added manually via the UI.

Usage notes:
- Press `a` in the tools list to add a new tool. You'll be prompted for a
	tool name, then to add one or more variants. For each variant you'll enter a
	free-form label (e.g. "IDE", "CLI", "Agentic Mode"), a launch path or
	command, an optional launch command, and whether the variant should open in
	a terminal (y) or launch detached (n).
- Press `d` to delete the selected tool from the list. While viewing a tool's
	details, press `d` to delete the currently-selected variant. Changes are
	saved immediately to `~/.config/aidock/tools.json`.
