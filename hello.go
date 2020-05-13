package hello

import (
	"fmt"
	"sync"

	"github.com/go-spirit/spirit/component"
	"github.com/go-spirit/spirit/mail"
	"github.com/go-spirit/spirit/worker"
)

type Body struct {
	Name string `json:"name"`
}

type Response struct {
	Reply string `json:"reply"`
}

type Hello struct {
	alias  string
	locker sync.RWMutex
}

func init() {
	component.RegisterComponent("hello-comp", NewHello)
}

func NewHello(alias string, opts ...component.Option) (comp component.Component, err error) {

	return &Hello{}, nil
}

func (h *Hello) Start() error {
	return nil
}

func (h *Hello) Stop() error {
	return nil
}

func (h *Hello) Alias() string {
	if h == nil {
		return ""
	}
	return h.alias
}

func (h *Hello) SayHi(session mail.Session) (err error) {

	h.locker.Lock()
	defer h.locker.Unlock()

	body := Body{}
	err = session.Payload().Content().ToObject(&body)
	if err != nil {
		return
	}

	if len(body.Name) == 0 {
		err = fmt.Errorf("name is empty")
		return
	}

	res := Response{"hello" + body.Name}
	err = session.Payload().Content().SetBody(res)
	if err != nil {
		return
	}

	return
}

func (h *Hello) Route(session mail.Session) worker.HandlerFunc {

	switch session.Query("action") {
	case "say_hi":
		{
			return h.SayHi
		}
	}

	return nil
}
