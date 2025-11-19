package logx

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger1(t *testing.T) {
	lo := NewLogger()
	lo.Info("test log info")
}

func TestNewLogger2(t *testing.T) {
	lo := NewFromCustom()
	lo.Info("test log info")
}

func TestNewByRoom(t *testing.T) {
	lo := NewByRoom(101, "../../logs/game")
	lo.Info("test log info")
	lo.Printf("test log info uid: %d", 1)
	lo.Debug("test log info", zap.Any("uid", 1))
	lo.Error("test log info", zap.Any("uid", 1))
	lo.Error("test log info", zap.Any("uid", 1))
}
