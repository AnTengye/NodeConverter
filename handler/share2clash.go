package handler

import (
	"fmt"
	"strings"

	"github.com/AnTengye/NodeConvertor/core"
	"github.com/kataras/iris/v12"
)

type ShareReq struct {
	Share string `url:"share,required"`
}

func ShareToClash(ctx iris.Context) {
	var nu ShareReq
	if err := ctx.ReadQuery(&nu); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	if nu.Share == "" {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("share_url is empty"))
		return
	}
	split := strings.Split(nu.Share, "://")
	if len(split) != 2 {
		ctx.StopWithError(
			iris.StatusBadRequest,
			fmt.Errorf("share_url is invalid, it should start with xxx://"))
		return
	}
	var node core.Node
	switch split[0] {
	case "vless":
		node = core.NewVLESSNode()
		err := node.FromShare(split[1])
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}
	case "vmess":
		ctx.StopWithError(
			iris.StatusBadRequest,
			fmt.Errorf("todo"),
		)
		return
	case "ss":
		ctx.StopWithError(
			iris.StatusBadRequest,
			fmt.Errorf("todo"),
		)
		return
	case "trojan":
		ctx.StopWithError(
			iris.StatusBadRequest,
			fmt.Errorf("todo"),
		)
		return
	default:
		ctx.StopWithError(
			iris.StatusBadRequest,
			fmt.Errorf("not support protocol: %s", split[0]),
		)
		return
	}
	ctx.WriteString(node.ToClash())
}
