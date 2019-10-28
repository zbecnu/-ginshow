package http

import (
"mygin/cache"
"net/http"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen(){
	http.Handle("/cache/",s.cacheHandler())
	http.Handle("/handler", s.cacheHandler())
	http.ListenAndServe(":1234",nil)
}


func (s *Server) cacheHandler() http.Handler{
	return &cacheHandler{s}
}

//new
func New(c cache.Cache)*Server{
	return &Server{c}
}


