# Ultimate Tic-Tac-Toe CLI

A command-line implementation of Ultimate Tic-Tac-Toe written in Go.

## Requirements

- [Go 1.26.4](https://go.dev/dl/) or a compatible newer version

## Run the project

From the project root, download the dependencies:

```bash
go mod download
```

Start the Bubble Tea TUI:

```bash
go run ./cmd/cli
```

Use the arrow keys to highlight a small board and press **Enter** to select it.
Then use the arrow keys to highlight a cell and press **Enter** to play. Press
**Esc** to return to board selection during a free move, or **q** to quit.

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
go build -o bin\UltimateTicTacToe.exe .\cmd\cli
.\bin\UltimateTicTacToe.exe
```

On macOS or Linux:

```bash
go build -o bin/uttt ./cmd/cli
./bin/uttt
```

## Playing

On a free move, choose the highlighted small board first, then choose a cell.
Otherwise, the previous move determines and highlights the required small board.
