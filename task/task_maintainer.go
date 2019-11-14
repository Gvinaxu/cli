package task

import (
	"github.com/Gvinaxu/cli/handler"
)

type Maintainer struct {
	scheduler *Scheduler
}

func NewTaskMaintainer(account *handler.Account) *Maintainer {
	scheduler := NewScheduler()
	tm := &Maintainer{
		scheduler: scheduler,
	}

	tm.taskInit(account)
	return tm
}

func (tm *Maintainer) taskInit(account *handler.Account) {
	// add task
	a := &AccountTask{account: account}
	tm.scheduler.Every(1).Hours().Do(a.RefreshToken)

}

func (tm *Maintainer) Start() {
	go tm.scheduler.Start()
}

func (tm *Maintainer) Stop() {
	tm.scheduler.Clear()
}
