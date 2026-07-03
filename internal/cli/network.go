package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
	"github.com/spf13/cobra"
)

// NetMessage is the newline-delimited JSON protocol shared by host and client.
// Only the host creates board states and decides whether moves are valid.
type NetMessage struct {
	Type         string    `json:"type"`
	Move         *NetMove  `json:"move,omitempty"`
	Board        string    `json:"board,omitempty"`
	Message      string    `json:"message,omitempty"`
	Player       game.Team `json:"player,omitempty"`
	Turn         game.Team `json:"turn,omitempty"`
	GameOver     bool      `json:"gameOver,omitempty"`
	FreeMove     bool      `json:"freeMove,omitempty"`
	NextBoardRow int       `json:"nextBoardRow,omitempty"`
	NextBoardCol int       `json:"nextBoardCol,omitempty"`
}

type NetMove struct {
	BoardRow int `json:"boardRow"`
	BoardCol int `json:"boardCol"`
	CellRow  int `json:"cellRow"`
	CellCol  int `json:"cellCol"`
}

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Host a LAN multiplayer game",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetString("port")
		return HostGame(port)
	},
}

var joinCmd = &cobra.Command{
	Use:   "join ADDRESS",
	Short: "Join a LAN multiplayer game",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return JoinGame(args[0])
	},
}

func init() {
	hostCmd.Flags().String("port", "8080", "TCP port to listen on")
	rootCmd.AddCommand(hostCmd, joinCmd)
}

// HostGame listens on every network interface and owns the official game state.
func HostGame(port string) error {
	port = strings.TrimPrefix(strings.TrimSpace(port), ":")
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("invalid port %q", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("could not host on port %s: %w", port, err)
	}
	defer listener.Close()

	fmt.Printf("Waiting for player 2 on port %s...\n", port)
	fmt.Println("Windows Firewall may need permission for this app on Private networks.")

	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("could not accept player: %w", err)
	}
	defer conn.Close()
	fmt.Println("Player 2 connected.")

	g := game.NewGame("Host", "Player 2")
	g.PickX(1)
	reader := bufio.NewReader(os.Stdin)
	remote := bufio.NewScanner(conn)

	if err := sendMessage(conn, NetMessage{Type: "welcome", Player: game.O, Message: "Connected to host. You are O."}); err != nil {
		return fmt.Errorf("could not welcome player: %w", err)
	}

	for g.Status == game.Playing {
		if err := sendState(conn, g); err != nil {
			fmt.Println("Opponent disconnected.")
			return nil
		}

		clearScreen()
		fmt.Println(renderBoard(g))

		if g.CurrentTurn == game.X {
			fmt.Println("Your turn. You are X.")
			move, err := promptMove(reader, g.FreeMove, g.NextBoardRow, g.NextBoardCol)
			if err != nil {
				return err
			}
			if err := applyNetMove(g, move); err != nil {
				fmt.Println("Invalid move:", err)
				continue
			}
			continue
		}

		fmt.Println("Waiting for opponent...")
		message, err := readMessage(remote)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Opponent disconnected.")
				return nil
			}
			if remote.Err() != nil {
				fmt.Println("Opponent disconnected:", remote.Err())
				return nil
			}
			_ = sendMessage(conn, NetMessage{Type: "error", Message: "Invalid message."})
			continue
		}
		if message.Type != "move" || message.Move == nil {
			_ = sendMessage(conn, NetMessage{Type: "error", Message: "Expected a move."})
			continue
		}
		if err := applyNetMove(g, *message.Move); err != nil {
			_ = sendMessage(conn, NetMessage{Type: "error", Message: "Invalid move: " + err.Error()})
			continue
		}
	}

	clearScreen()
	fmt.Println(renderBoard(g))
	result := gameResult(g)
	fmt.Println(result)
	if err := sendMessage(conn, stateMessage(g, "game_over", result)); err != nil {
		fmt.Println("Opponent disconnected before receiving the result.")
	}
	if err := remote.Err(); err != nil {
		return fmt.Errorf("network read failed: %w", err)
	}
	return nil
}

// JoinGame is intentionally display/input-only; all rule decisions stay on the host.
func JoinGame(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("could not connect to %s: %w", address, err)
	}
	defer conn.Close()

	fmt.Println("Connected to host.")
	remote := bufio.NewScanner(conn)
	input := bufio.NewReader(os.Stdin)

	for {
		message, err := readMessage(remote)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Host disconnected.")
				return nil
			}
			if remote.Err() != nil {
				fmt.Println("Host disconnected:", remote.Err())
				return nil
			}
			fmt.Println("Could not parse host message:", err)
			continue
		}

		switch message.Type {
		case "welcome", "info", "error":
			fmt.Println(message.Message)
		case "board", "game_over":
			clearScreen()
			fmt.Println(message.Board)
			if message.Message != "" {
				fmt.Println(message.Message)
			}
			if message.GameOver {
				return nil
			}
			if message.Turn != game.O {
				fmt.Println("Waiting for opponent...")
				continue
			}

			fmt.Println("Your turn. You are O.")
			move, err := promptMove(input, message.FreeMove, message.NextBoardRow, message.NextBoardCol)
			if err != nil {
				return err
			}
			if err := sendMessage(conn, NetMessage{Type: "move", Move: &move}); err != nil {
				fmt.Println("Host disconnected.")
				return nil
			}
		}
	}
}

func sendState(conn net.Conn, g *game.Game) error {
	message := "Waiting for opponent..."
	if g.CurrentTurn == game.O {
		message = "Your turn."
	}
	return sendMessage(conn, stateMessage(g, "board", message))
}

func stateMessage(g *game.Game, messageType, message string) NetMessage {
	return NetMessage{
		Type:         messageType,
		Board:        renderBoard(g),
		Message:      message,
		Turn:         g.CurrentTurn,
		GameOver:     g.Status == game.Finished,
		FreeMove:     g.FreeMove,
		NextBoardRow: g.NextBoardRow,
		NextBoardCol: g.NextBoardCol,
	}
}

func sendMessage(conn net.Conn, message NetMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = conn.Write(data)
	return err
}

func readMessage(scanner *bufio.Scanner) (NetMessage, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return NetMessage{}, err
		}
		return NetMessage{}, io.EOF
	}

	var message NetMessage
	if err := json.Unmarshal(scanner.Bytes(), &message); err != nil {
		return NetMessage{}, err
	}
	return message, nil
}

func applyNetMove(g *game.Game, move NetMove) error {
	return g.MakeMove(move.BoardRow, move.BoardCol, move.CellRow, move.CellCol)
}

func promptMove(reader *bufio.Reader, freeMove bool, nextBoardRow, nextBoardCol int) (NetMove, error) {
	move := NetMove{BoardRow: nextBoardRow, BoardCol: nextBoardCol}
	var err error

	if freeMove {
		fmt.Println("Free move: choose any board.")
		if move.BoardRow, err = promptCoordinate(reader, "Choose board row (0-2): "); err != nil {
			return NetMove{}, err
		}
		if move.BoardCol, err = promptCoordinate(reader, "Choose board col (0-2): "); err != nil {
			return NetMove{}, err
		}
	} else {
		fmt.Printf("Required board: (%d, %d)\n", nextBoardRow, nextBoardCol)
	}

	if move.CellRow, err = promptCoordinate(reader, "Choose cell row (0-2): "); err != nil {
		return NetMove{}, err
	}
	if move.CellCol, err = promptCoordinate(reader, "Choose cell col (0-2): "); err != nil {
		return NetMove{}, err
	}
	return move, nil
}

func promptCoordinate(reader *bufio.Reader, prompt string) (int, error) {
	for {
		fmt.Print(prompt)
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return 0, err
		}
		value, parseErr := strconv.Atoi(strings.TrimSpace(line))
		if parseErr == nil && value >= 0 && value <= 2 {
			return value, nil
		}
		if errors.Is(err, io.EOF) {
			return 0, io.EOF
		}
		fmt.Println("Enter a number from 0 to 2.")
	}
}

func gameResult(g *game.Game) string {
	if g.Winner != game.Empty {
		return "Winner: " + string(g.Winner)
	}
	return "Draw!"
}
