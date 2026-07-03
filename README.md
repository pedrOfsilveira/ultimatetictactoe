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
go run ./cmd/cli local
go run ./cmd/cli rules
```

## LAN multiplayer

The host is player **X** and owns the official game state. The joining player is
**O**. All moves are validated by the host.

Start a host that listens on every network interface:

```powershell
.\bin\UltimateTicTacToe.exe host --port 8080
```

Join it from another terminal using the host's LAN or Radmin VPN address:

```powershell
.\bin\UltimateTicTacToe.exe join 192.168.0.25:8080
```

For a test on one computer, join `127.0.0.1:8080`. Windows Firewall may ask you
to allow the executable; allow it on Private networks for LAN play.

## Build an executable

On Windows:

```powershell
go build -ldflags="-H windowsgui" -o bin\UltimateTicTacToe.exe .\cmd\launcher
.\bin\UltimateTicTacToe.exe
```

On macOS or Linux:

```bash
go build -o bin/uttt ./cmd/cli
./bin/uttt
```

## Playing

Rows and columns use values from `0` to `2`. On a free move, choose the small board first, then choose a cell. Otherwise, the previous move determines the required small board.
