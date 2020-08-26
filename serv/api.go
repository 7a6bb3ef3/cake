package main

import (
	"github.com/nynicg/cake/lib/ahoy"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/nynicg/cake/lib/log"
)

func runApiServ(){
	router := httprouter.New()
	router.GET("/register/:uid/:cmd", BasicAuth(Register ,"114514" ,"114514"))
	//router.POST("/register/:uid/:cmd", BasicAuth(Register ,"114514" ,"114514"))

	log.Info("API service listen on " ,config.LocalApi)
	if e := http.ListenAndServe(config.LocalApi, router);e != nil{
		log.Errorx("api service has crashed." ,e)
	}
}

func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func Register(w http.ResponseWriter,r *http.Request ,params httprouter.Params){
	uid := params.ByName("uid")
	cmds ,e := strconv.Atoi(params.ByName("cmd"))
	if e != nil{
		w.Write([]byte("error.param CMD only accept integer. range [0,255]"))
		return
	}
	RegisterUidCmd(ahoy.Command(cmds) ,uid)
	w.Write([]byte("ok"))
}
