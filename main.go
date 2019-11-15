package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"strings"

	"github.com/Gvinaxu/cli/handler"
	"github.com/Gvinaxu/cli/task"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	Version = "v1.0.0"
	Website = "https://gavln.com"
	Banner  = `
     ___          ___                      ___   ___     
    /  /\        /  /\        ___         /  /\ /  /\    
   /  /::\      /  /::\      /  /\       /  /://  /::|   
  /  /:/\:\    /  /:/\:\    /  /:/      /  /://  /:|:|   
 /  /:/  \:\  /  /::\ \:\  /  /:/      /  /://  /:/|:|__ 
/__/:/_\_ \:\/__/:/\:\_\:\/__/:/  ___ /__/://__/:/ |:| /\
\  \:\__/\_\/\__\/  \:\/:/|  |:| /  /\\  \:\\__\/  |:|/:/
 \  \:\ \:\       \__\::/ |  |:|/  /:/ \  \:\   |  |:/:/ 
  \  \:\/:/       /  /:/  |__|:|__/:/   \  \:\  |__|::/  
   \  \::/       /__/:/    \__\____/     \  \:\ /__/:/   
    \__\/        \__\/                    \__\/ \__\/        	%s
	
Know its white, keep its black
%s
______________________________
					
`
)

func main() {
	app := &App{}
	app.Start()
	app.Run()
}

type App struct {
	h       *handler.Handler
	account *handler.Account
	tm      *task.Maintainer
}

func (a *App) Start() {
	a.account = a.checkUser()

	f := &handler.FileReq{}
	a.h = handler.NewHandler(f)

	a.tm = task.NewTaskMaintainer(a.account)
	a.tm.Start()
}

func (a *App) Run() {

	fmt.Printf(Banner, Version, Website)
	history := make([]string, 0)
	f := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")

		line, err := f.ReadString('\n')
		if err != nil {
			panic(err)
		}
		history = append(history, line)
		command, args, err := a.getCommandAndArgs(line)
		if err != nil {
			continue
		}
		if command == "history" {
			for i, v := range history {
				fmt.Printf("%d		%s", i, v)
			}
			continue
		}
		if command == "quit" {
			fmt.Println("bye!")
			break
		}
		a.h.InvokeCmd(command, args)
	}
}

func (a *App) getCommandAndArgs(line string) (command string, args []string, err error) {
	line = strings.TrimSpace(line)
	all := strings.Split(line, " ")
	if len(all) == 0 {
		return "", nil, errors.New("input is nil")
	}
	args = make([]string, 0)
	for i, v := range all {
		if i == 0 {
			continue
		}
		if v == "" {
			continue
		}
		args = append(args, v)
	}
	return all[0], args, nil
}

func (a *App) checkUser() *handler.Account {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Gavln User Name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Print("Enter Password: ")
	bytes, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	password := string(bytes)
	name = strings.ReplaceAll(name, "\n", "")
	account := handler.NewAccount(name, password)
	_, err = account.Login()
	if err != nil {
		panic(err)
	}
	return account
}
