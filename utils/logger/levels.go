package logger

import (
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
	return startColor(color) + str + endColor()
}

func ColorToInt(color color.RGBA) (int, int, int, int) {
	r, g, b, a := color.R, color.G, color.B, color.A

	red, green, blue, alpha := int(r), int(g), int(b), int(a)

	return red, green, blue, alpha
}

func startColor(color color.RGBA) string {
	red, green, blue, alpha := ColorToInt(color)

	mode := "38;2;"

	if alpha >= 255 {
		mode = "48;2;"
	}

	colorStr := strconv.Itoa(red) + ";" + strconv.Itoa(green) + ";" + strconv.Itoa(blue)

	return "\x1b[" + mode + colorStr + "m"
}

func endColor() string {
	return "\x1b[0m"
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
		return ColorCode("DEV  ", color.RGBA{
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
