# AI-Dock

`aidock` is a cross-platform terminal UI for discovering and launching AI tools
(CLI tools, IDEs, and desktop apps). It uses the Charm stack: Bubble Tea,
Bubbles, and Lip Gloss.

## Run

```sh
go run ./cmd/aidock
```

Tool configuration is stored at `~/.config/aidock/tools.json`. The first run
creates this file with an empty tools array and displays a couple of dummy
tools until configuration is added.
