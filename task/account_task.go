package task

import (
	"github.com/Gvinaxu/cli/handler"
)

type AccountTask struct {
	account *handler.Account
}

func (at *AccountTask) RefreshToken() {
	at.account.RefreshToken()
}
