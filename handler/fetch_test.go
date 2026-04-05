package handler

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/AnTengye/NodeConverter/lib/network"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestFetchNodesLogsSubscriptionDownloadFailure(t *testing.T) {
	observed, logs := observer.New(zapcore.DebugLevel)
	previousLogger := zap.L()
	previousSugar := zap.S()
	logger := zap.New(observed)
	zap.ReplaceGlobals(logger)
	t.Cleanup(func() {
		_ = previousSugar
		zap.ReplaceGlobals(previousLogger)
	})

	previousCacheGet := network.CacheGET
	network.CacheGET = func(url string) ([]byte, error) {
		return nil, errors.New("dial tcp timeout")
	}
	t.Cleanup(func() {
		network.CacheGET = previousCacheGet
	})

	_, err := fetchNodes("http://example.com/sub")
	if err == nil {
		t.Fatal("expected fetchNodes to fail")
	}

	entry := logs.FilterMessage("fetch subscription failed").All()
	if len(entry) != 1 {
		t.Fatalf("expected exactly one download failure log, got %d", len(entry))
	}
	fields := entry[0].ContextMap()
	if fields["url"] != "http://example.com/sub" {
		t.Fatalf("expected url field to be recorded, got %#v", fields["url"])
	}
	if fields["stage"] != "download" {
		t.Fatalf("expected stage=download, got %#v", fields["stage"])
	}
	if fields["error"] != "dial tcp timeout" {
		t.Fatalf("expected error field to include root cause, got %#v", fields["error"])
	}
}

func TestHandlerSubResponseLogsShareConversionFailureDetails(t *testing.T) {
	observed, logs := observer.New(zapcore.DebugLevel)
	previousLogger := zap.L()
	previousSugar := zap.S()
	logger := zap.New(observed)
	zap.ReplaceGlobals(logger)
	t.Cleanup(func() {
		_ = previousSugar
		zap.ReplaceGlobals(previousLogger)
	})

	payload := base64.StdEncoding.EncodeToString([]byte("invalid-share\nss://bad"))
	nodes, err := handlerSubResponse([]byte(payload))
	if err != nil {
		t.Fatalf("expected handlerSubResponse to continue on invalid lines, got %v", err)
	}
	if len(nodes) != 0 {
		t.Fatalf("expected invalid shares to be skipped, got %d nodes", len(nodes))
	}

	entries := logs.FilterMessage("subscription node conversion failed").All()
	if len(entries) != 2 {
		t.Fatalf("expected two conversion failure logs, got %d", len(entries))
	}
	first := entries[0].ContextMap()
	if first["line"] != int64(1) {
		t.Fatalf("expected first failed line index to be 1, got %#v", first["line"])
	}
	if first["share_url"] != "invalid-share" {
		t.Fatalf("expected original share url in log, got %#v", first["share_url"])
	}
	if first["error"] == "" {
		t.Fatal("expected error detail in first log entry")
	}
	second := entries[1].ContextMap()
	if second["line"] != int64(2) {
		t.Fatalf("expected second failed line index to be 2, got %#v", second["line"])
	}
	if second["error"] == "" {
		t.Fatal("expected error detail in second log entry")
	}
}
