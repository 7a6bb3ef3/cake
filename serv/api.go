package main

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nynicg/cake/lib/log"
)

func runApiServ() {
	if !globalConfig.ApiConfig.EnableApi {
		return
	}
	router := httprouter.New()
	router.GET("/stat", wrap(Stat))
	router.POST("/register/:uid", wrap(Register))

	log.Info("API service listen on ", globalConfig.ApiConfig.LocalApiAddr)
	if e := http.ListenAndServe(globalConfig.ApiConfig.LocalApiAddr, router); e != nil {
		log.Errorx("api service has crashed.", e)
	}
}

func wrap(h httprouter.Handle) httprouter.Handle {
	return BasicAuth(h, globalConfig.ApiConfig.BasicAuthUser, globalConfig.ApiConfig.BasicAuthPassword)
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

func Register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uid := params.ByName("uid")
	if len(uid) != 32 {
		w.Write([]byte("uid length must be 32"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DefUidManager().RegisterUid(uid ,UIDInfo{
		CreateTime: time.Now().Unix(),
		Addr:       r.RemoteAddr,
	})
	w.Write([]byte("ok"))
}
