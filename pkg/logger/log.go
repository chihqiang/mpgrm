package logger

import (
	"github.com/fatih/color"
	"log"
)

// 彩色前缀
var (
	infoColor  = color.New(color.FgBlue, color.Bold)     // 蓝色：低严重程度
	warnColor  = color.New(color.FgHiYellow, color.Bold) // 黄色：中等严重程度
	errorColor = color.New(color.FgHiRed, color.Bold)    // 红色：高严重程度
)

// Info 输出普通信息
func Info(format string, a ...any) {
	log.Printf("%s", formatMessage(infoColor, "INFO", format, a...))
}

// Warning 输出警告信息
func Warning(format string, a ...any) {
	log.Printf("%s", formatMessage(warnColor, "WARN", format, a...))
}

// Error 输出错误信息
func Error(format string, a ...any) {
	log.Printf("%s", formatMessage(errorColor, "ERROR", format, a...))
}

// 格式化信息，将前缀和内容统一上色
func formatMessage(c *color.Color, level, format string, a ...any) string {
	return c.Sprintf("["+level+"] "+format, a...) // 整行颜色
}
