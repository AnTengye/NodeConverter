package main

import "github.com/kataras/iris/v12"

type ClashUrl struct {
	Url string `url:"url,required"`
}

// TODO: 未实现
func ClashToShare(ctx iris.Context) {
	var nu ClashUrl
	if err := ctx.ReadQuery(&nu); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	ctx.WriteString("Success")
}

type ShareReq struct {
	ShareUrl string `url:"share_url,required"`
}

func ShareToClash(ctx iris.Context) {
	var nu ShareReq
	if err := ctx.ReadQuery(&nu); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	node := NewVLESSNode()
	err := node.FromShare(nu.ShareUrl)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	ctx.WriteString(node.ToClash())
}
