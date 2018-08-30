package logger

import (
	"context"
	"testing"
	"time"
)

func Test_DefaultLogger(t *testing.T) {
	Log.Info("Info")
	Log.WithFields(map[string]interface{}{
		"fielda": "a",
		"fieldb": "b",
	}).Info("Info fields")

	type StructA struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	sa := StructA{
		A: 1,
		B: 2,
	}

	Log.Error(sa)
}

func Test_NewJSONLogger(t *testing.T) {
	logger, err := NewJSONLogger("./testdata", "test.log", "debug")
	if err != nil {
		t.Fail()
	}

	logger.Info("Info")
	logger.WithFields(map[string]interface{}{
		"fielda": "a",
		"fieldb": "b",
	}).Info("Info fields")

	type StructA struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	sa := StructA{
		A: 1,
		B: 2,
	}

	logger.Error(sa)
}

func Test_NewJSONLogger_multi(t *testing.T) {
	l, _ := NewJSONLogger("./testdata", "new_logger_multi", "debug")
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 10; i++ {
		go func(ctx context.Context, i int) {
			ticker := time.NewTicker(2 * time.Second)
			for {
				select {
				case <-ticker.C:
					l.Info("Test_NewJSONLogger_multi-", i)
				case <-ctx.Done():
					break
				}
			}
		}(ctx, i)
	}

	time.Sleep(10 * time.Second)
	cancel()
}
