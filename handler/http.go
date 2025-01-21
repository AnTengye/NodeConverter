package handler

import "github.com/go-resty/resty/v2"

var restyCli *resty.Client

func InitResty() {
	restyCli = resty.New()
}
