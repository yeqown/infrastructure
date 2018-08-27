package logger

import (
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

func Test_NewJSONLogger(t *testing.T) {
	l, _ := NewJSONLogger("./testdata", "new_logger_with_spiliter", "debug")
	timer := time.NewTimer(10 * time.Second)
	quit := false
	go func() {
		select {
		case <-timer.C:
			quit = true
		}
	}()

	for {
		if quit {
			break
		}

		l.Info("test msg")
		time.Sleep(1 * time.Second)
	}
	t.Log("done")
}
