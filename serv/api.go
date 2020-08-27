package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
	"github.com/nynicg/cake/lib/ahoy"
	"github.com/nynicg/cake/lib/log"
)

func runApiServ() {
	if !config.EnableAPI {
		return
	}
	router := httprouter.New()
	router.GET("/register/:uid/:cmd", wrap(Register))
	router.GET("/stat", wrap(Stat))

	log.Info("API service listen on ", config.LocalApiAddr)
	if e := http.ListenAndServe(config.LocalApiAddr, router); e != nil {
		log.Errorx("api service has crashed.", e)
	}
}

func wrap(h httprouter.Handle) httprouter.Handle {
	return BasicAuth(h, config.BAUserName, config.BAPassword)
}

func BasicAuth(h httprouter.Handle, usr, pwd string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == usr && password == pwd {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func Register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uid := params.ByName("uid")
	cmds, e := strconv.Atoi(params.ByName("cmd"))
	if e != nil || cmds < 0 || cmds > 255 {
		w.Write([]byte("error.param CMD only accept integer. range [0,255]"))
		return
	}
	RegisterUidCmd(ahoy.Command(cmds), uid)
	w.Write([]byte("ok"))
}

type ProxyStat struct {
	TotalUp   int64
	TotalDown int64
}

func (p *ProxyStat) Add(up, down int) {
	atomic.AddInt64(&p.TotalUp, int64(up))
	atomic.AddInt64(&p.TotalDown, int64(down))
}

func Stat(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	p := *proxyStat
	bts, e := json.Marshal(&p)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server internal error, check the lastest error log"))
		log.Error("api service.", e)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(bts)
}
