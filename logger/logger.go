// Copyright 2024 wangxin.jeffry@gmail.com
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package nlog defines standardized ways for zap logging

package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

type Formatter string

const (
	JSON Formatter = "json"
	TEXT Formatter = "text"
)

type LogConfig struct {
	Formatter string `mapstructure:"formatter" json:"formatter" yaml:"formatter" validate:"oneof=json text"` // accepted values: "json", "text"
}

func (cfg *LogConfig) Levels() []zapcore.Level {
	levels := make([]zapcore.Level, 0)
	level, err := zapcore.ParseLevel("info")
	if err != nil {
		level = zapcore.DebugLevel
	}
	levels = append(levels, level)
	return levels
}

func (cfg *LogConfig) Encoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		NameKey:        "name",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if cfg.Formatter == string(JSON) {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)
}

func InitLogger(cfg *LogConfig) {
	level := cfg.Levels()
	encoder := cfg.Encoder()
	core := zapcore.NewCore(encoder, zapcore.AddSync(log.Writer()), level[0])

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(Logger)

}
