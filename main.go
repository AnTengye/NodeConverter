package main

import (
	"flag"
	"os"

	"github.com/AnTengye/NodeConvertor/handler"
	"github.com/AnTengye/NodeConvertor/lib/network"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "config.yaml", "the config file")

// Read the example and its comments carefully.
func makeAccessLog() *accesslog.AccessLog {
	ac := accesslog.New(zap.CombineWriteSyncers(os.Stdout))

	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})

	return ac
}
func main() {
	flag.Parse()

	if *configFile != "" {
		viper.SetConfigFile(*configFile) // 指定配置文件（路径 + 配置文件名）
		viper.SetConfigType("yaml")      // 如果配置文件名中没有扩展名，则需要显式指定配置文件的格式
	} else {
		viper.AddConfigPath(".")             // 把当前目录加入到配置文件的搜索路径中
		viper.AddConfigPath("$HOME/.config") // 可以多次调用 AddConfigPath 来设置多个配置文件搜索路径
		viper.SetConfigName("config.yaml")   // 指定配置文件名（没有扩展名）
	}
	err := viper.ReadInConfig()
	if err != nil {
		viper.SetDefault("API.Listen", ":25500")
		_ = viper.WriteConfig()
	}
	var (
		logger *zap.Logger
		app    *iris.Application
	)
	if !viper.GetBool("API.Debug") {
		logger, _ = zap.NewProduction()
		app = iris.New()
	} else {
		logger, _ = zap.NewDevelopment()
		app = iris.Default()
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	ac := makeAccessLog()
	defer ac.Close()
	network.InitResty(viper.GetBool("API.Debug"))
	app.UseRouter(ac.Handler)
	app.Get("/sub", handler.Sub)
	app.Listen(viper.GetString("API.Listen"))
}
