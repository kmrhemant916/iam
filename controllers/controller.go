package controllers

import "net/http"


func Register(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello world"))
}
