package logger

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

const DeveloperLevel zapcore.Level = -2

func ParseLevel(s string) zapcore.Level {
	switch strings.ToLower(s) {
	case "dev":
		return DeveloperLevel
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func ColorCode(str string, color color.RGBA) string {
	r, g, b := color.R, color.G, color.B

	red, green, blue := int(r), int(g), int(b)

	colorStr := strconv.Itoa(red) + ";" + strconv.Itoa(green) + ";" + strconv.Itoa(blue)

	fmt.Println(colorStr)

	return "\x1b[38;2;" + colorStr + "m" + str + "\x1b[0m"
}

func LevelString(l zapcore.Level) string {
	switch l {
	case DeveloperLevel:
		return "dev"
	default:
		return l.CapitalString()
	}
}

func CapitalLevel(l zapcore.Level) string {
	switch l {
	case DeveloperLevel:
		return ColorCode("DEV", color.RGBA{
			R: 95, G: 175, B: 135,
		})
	default:
		return l.CapitalString()
	}
}

func CustomEncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case DeveloperLevel:
		enc.AppendString(CapitalLevel(l))
	default:
		zapcore.CapitalColorLevelEncoder(l, enc)
	}
}
