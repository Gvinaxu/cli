package handler

import (
	"fmt"
	"strings"
)

type Handler struct {
	file *FileReq
}

func NewHandler(file *FileReq) *Handler {
	return &Handler{
		file: file,
	}
}

func (h *Handler) InvokeCmd(command string, args []string) {
	cmd := strings.ToLower(command)
	switch cmd {
	case "tree", "t":
		path := ""
		if len(args) != 0 {
			path = args[0]
		}
		h.file.Tree(path)
	case "help", "h":
		h.printHelp()
	case "":
		return
	default:
		fmt.Printf("command %s not exists \n", command)
		h.printHelp()
	}
}

func (h *Handler) printHelp() {
	fmt.Println("Usage: <command> [args]")
	fmt.Println()
	fmt.Println("commands:")
	fmt.Println("tree|t		[path]			--show Gavln account all files")
	fmt.Println("help|h					--show Gavln help")
}
