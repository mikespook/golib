package signal

import (
	"os"
	S "os/signal"

	"github.com/mikespook/golib/idgen"
)

// return true if need to break
type Callback func() bool

type cbHandler struct {
	Id       interface{}
	Callback Callback
}

type Handler struct {
	schan chan os.Signal
	cb    map[os.Signal][]cbHandler
	id    idgen.IdGen
}

func NewHandler(id idgen.IdGen) (sh *Handler) {
	if id == nil {
		id = idgen.NewAutoIncId()
	}
	sh = &Handler{
		schan: make(chan os.Signal),
		cb:    make(map[os.Signal][]cbHandler, 5),
		id:    id,
	}
	return
}

func (sh *Handler) Bind(s os.Signal, cb Callback) (id interface{}) {
	S.Notify(sh.schan, s)
	id = sh.id.Id()
	sh.cb[s] = append(sh.cb[s], cbHandler{id, cb})
	return
}

func (sh *Handler) Unbind(s os.Signal, id interface{}) bool {
	for k, v := range sh.cb[s] {
		if v.Id == id {
			sh.cb[s] = append(sh.cb[s][:k], sh.cb[s][k+1:]...)
			return true
		}
	}
	return false
}

func (sh *Handler) Loop() os.Signal {
	for s := range sh.schan {
		if cbs, ok := sh.cb[s]; ok && cbs != nil {
			for _, v := range cbs {
				if v.Callback() {
					return s
				}
			}
		}
	}
	return nil
}

func (sh *Handler) Close() {
	S.Stop(sh.schan)
	close(sh.schan)
}

var (
	DefaultHandler = NewHandler(nil)
)

func Bind(s os.Signal, cb Callback) interface{} {
	return DefaultHandler.Bind(s, cb)
}

func Unbind(s os.Signal, id interface{}) bool {
	return DefaultHandler.Unbind(s, id)
}

func Loop() os.Signal {
	return DefaultHandler.Loop()
}

func Send(pid int, signal os.Signal) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Signal(signal); err != nil {
		return err
	}
	return nil
}

func Close() {
	DefaultHandler.Close()
}
