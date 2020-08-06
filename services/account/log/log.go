package log

import (
	"go.uber.org/zap"
)

// Config はロガーの設定
type Config struct {
	Mode string
}

// NewLogger は与えられた設定を元にロガーを作成する
func NewLogger(conf Config) (*zap.Logger, error) {
	if conf.Mode == "development" {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}
