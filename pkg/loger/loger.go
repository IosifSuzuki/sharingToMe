package loger

import "fmt"

type prefixTextInConsole string
const (
	info prefixTextInConsole = "Info"
	err prefixTextInConsole = "Error"
)

type textColorInConsole string
const (
	resetColor textColorInConsole = "\033[0m"
	redColor textColorInConsole   = "\033[31m"
	greenColor textColorInConsole = "\033[32m"
)

func PrintInfo(format string, args ...interface{}) {
	var colorInConsole, resetColor = string(greenColor), string(resetColor)
	fmt.Printf(colorInConsole + "\t%s: %s\n" + resetColor, info, fmt.Sprintf(format, args...))
}

func PrintError(err error) {
	var colorInConsole, resetColor = string(redColor), string(resetColor)
	fmt.Printf(colorInConsole + "\t%s: %v\n" + resetColor, info, err)
}