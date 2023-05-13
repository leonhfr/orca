package uci

import (
	"strconv"
	"strings"
	"time"
)

const fenFields = 6

// parse parses UCI commands and returns a Command object.
func parse(command []string) command {
	var index int
	if len(command) == 0 {
		return nil
	}

top:
	switch command[index] {
	case "uci":
		return commandUCI{}
	case "debug":
		if len(command) > 1 {
			return commandDebug{on: command[index+1] == "on"}
		}
	case "isready":
		return commandIsReady{}
	case "setoption":
		if len(command) > 1 {
			return parseCommandSetOption(command[index+1:])
		}
	case "ucinewgame":
		return commandUCINewGame{}
	case "position":
		if len(command) > 1 {
			return parseCommandPosition(command[index+1:])
		}
	case "go":
		if len(command) > 1 {
			return parseCommandGo(command[index+1:])
		}
	case "stop":
		return commandStop{}
	case "quit":
		return commandQuit{}
	default:
		if len(command) == index+1 {
			break
		}
		// unknown commands should be ignored
		index++
		goto top
	}
	return nil
}

// parseCommandSetOption parses setoption UCI commands.
func parseCommandSetOption(command []string) commandSetOption {
	var c commandSetOption
	if len(command) >= 4 && command[0] == "name" && command[2] == "value" {
		c.name = command[1]
		c.value = command[3]
	}
	return c
}

// parseCommandPosition parses position UCI commands.
func parseCommandPosition(command []string) commandPosition {
	var c commandPosition
	var index int

	if command[0] == "startpos" {
		c.startPos = true
		index = 1
	} else if command[0] == "fen" && len(command) >= fenFields+1 {
		c.fen = strings.Join(command[1:fenFields+1], " ")
		index = fenFields + 1
	}

	if len(command) > index && command[index] == "moves" {
		for index++; index < len(command); index++ {
			c.moves = append(c.moves, command[index])
		}
	}

	return c
}

// parseCommandGo parses go UCI commands.
func parseCommandGo(command []string) commandGo {
	var c commandGo

	for index := 0; index < len(command); index++ {
		switch command[index] {
		case "wtime":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.whiteTime = time.Duration(t) * time.Millisecond
				index++
			}
		case "btime":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.blackTime = time.Duration(t) * time.Millisecond
				index++
			}
		case "winc":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.whiteIncrement = time.Duration(t) * time.Millisecond
				index++
			}
		case "binc":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.blackIncrement = time.Duration(t) * time.Millisecond
				index++
			}
		case "movestogo":
			if len(command) >= index+1 {
				n, _ := strconv.Atoi(command[index+1])
				c.movesToGo = n
				index++
			}
		case "searchmoves":
			if len(command) >= index+1 {
				for index++; index < len(command); index++ {
					c.searchMoves = append(c.searchMoves, command[index])
				}
				return c
			}
		case "depth":
			if len(command) >= index+1 {
				n, _ := strconv.Atoi(command[index+1])
				c.depth = n
				index++
			}
		case "nodes":
			if len(command) >= index+1 {
				n, _ := strconv.Atoi(command[index+1])
				c.nodes = n
				index++
			}
		case "movetime":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.moveTime = time.Duration(t) * time.Millisecond
				index++
			}
		case "infinite":
			c.infinite = true
		}
	}

	return c
}
