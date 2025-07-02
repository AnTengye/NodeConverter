package network

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"resty.dev/v3"
)

var (
	RestyCli *resty.Client
	cacheCli *cache.Cache
	CacheGET func(url string) ([]byte, error)
)

func InitResty() {
	RestyCli = resty.New().SetLogger(zap.S()).
		SetDebug(viper.GetBool("API.Debug")).
		SetDebugLogFormatter(resty.DebugLogJSONFormatter).
		SetRetryCount(viper.GetInt("API.RetryCount")).
		SetTimeout(time.Duration(viper.GetInt("API.Timeout")) * time.Second)
	cacheCli = cache.New(5*time.Minute, 10*time.Minute)
	CacheGET = func(url string) ([]byte, error) {
		// 尝试读缓存
		zap.S().Debugw("try read cache", "url", url)
		if data, found := cacheCli.Get(url); found {
			zap.S().Debugw("read cache", "url", url)
			return data.([]byte), nil
		}
		zap.S().Debugw("cache not found", "url", url)
		// 缓存未命中，发起请求
		var result []byte
		resp, err := RestyCli.R().Get(url)
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return nil, fmt.Errorf("request error: %v", resp.Error())
		}
		result = resp.Bytes()
		// 缓存结果（根据需求调整过期时间）
		cacheCli.Set(url, result, 10*time.Minute) // 根据API特性设置合理过期时间
		return result, nil
	}
}

func DeleteCache(url string) {
	cacheCli.Delete(url)
}
