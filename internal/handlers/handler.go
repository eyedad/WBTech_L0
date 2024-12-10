package handlers

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Register(oruter *httprouter.Router)
}
