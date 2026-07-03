# Ultimate Tic-Tac-Toe CLI

A command-line implementation of Ultimate Tic-Tac-Toe written in Go.

## Requirements

- [Go 1.26.4](https://go.dev/dl/) or a compatible newer version

## Run the project

From the project root, download the dependencies:

```bash
go mod download
```

Start the interactive CLI:

```bash
go run ./cmd/cli
```

At the `uttt>` prompt, use one of these commands:

- `play` starts a new game.
- `rules` displays the rules.
- `help` lists the available commands.
- `exit` closes the program.

You can also run a command directly:

```bash
go run ./cmd/cli play
go run ./cmd/cli rules
```

## Build an executable

On Windows:

```powershell
go build -o bin/uttt.exe ./cmd/cli
.\bin\uttt.exe
```

On macOS or Linux:

```bash
go build -o bin/uttt ./cmd/cli
./bin/uttt
```

## Playing

Rows and columns use values from `0` to `2`. On a free move, choose the small board first, then choose a cell. Otherwise, the previous move determines the required small board.
