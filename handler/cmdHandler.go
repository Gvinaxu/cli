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

	var err error
	switch cmd {
	case "tree", "t":
		path := ""
		if len(args) != 0 {
			path = args[0]
		}
		_, err = h.file.Tree(path)
	case "stat":
		path := ""
		if len(args) != 0 {
			path = args[0]
		} else {
			fmt.Println("file path is nil")
			return
		}
		_, err = h.file.Stat(path)
	case "move", "mv":
		old, new := "", ""
		if len(args) == 2 {
			old = args[0]
			new = args[1]
		} else {
			fmt.Println("parameter is missing")
			return
		}
		err = h.file.Mv(old, new)
	case "copy", "cp":
		old, new := "", ""
		if len(args) == 2 {
			old = args[0]
			new = args[1]
		} else {
			fmt.Println("parameter is missing")
			return
		}
		err = h.file.Cp(old, new)
	case "find":
		path := ""
		search := ""
		if len(args) >= 1 {
			search = args[0]
		} else {
			fmt.Println("search  is nil")
			return
		}
		if len(args) >= 2 {
			path = args[1]
		}
		_, err = h.file.Search(search, path)
	case "remove", "rm":
		path := ""
		if len(args) != 0 {
			path = args[0]
		} else {
			fmt.Println("file path is nil")
			return
		}
		err = h.file.Rm(path)
	case "mkdir":
		path := ""
		if len(args) != 0 {
			path = args[0]
		}
		if path == "" {
			fmt.Println("file path is nil")
			return
		}
		err = h.file.Mkdir(path)
	case "upload", "up":

	case "download", "down":

	case "help", "h":
		h.printHelp()
	case "":
		return
	default:
		fmt.Printf("command %s not exists \n", command)
		h.printHelp()
	}
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) printHelp() {
	fmt.Println("Usage: <command> [args]")
	fmt.Println()
	fmt.Println("commands:")
	fmt.Println("tree|t		[path]			--show Gavln account all files")
	fmt.Println("help|h					--show Gavln help")
}
