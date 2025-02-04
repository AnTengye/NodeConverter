package network

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

var (
	RestyCli *resty.Client
	cacheCli *cache.Cache
	CacheGET func(url string) ([]byte, error)
)

func InitResty(debug bool) {
	RestyCli = resty.New().SetLogger(zap.S()).SetDebug(debug)
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
		result = resp.Body()
		// 缓存结果（根据需求调整过期时间）
		cacheCli.Set(url, result, 10*time.Minute) // 根据API特性设置合理过期时间
		return result, nil
	}
}

func DeleteCache(url string) {
	cacheCli.Delete(url)
}
