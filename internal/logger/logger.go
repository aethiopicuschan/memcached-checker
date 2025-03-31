package logger

func OK(msg string) {
	println("\033[32m✓\033[0m " + msg)
}

func Fail(msg string) {
	println("\033[31m×\033[0m " + msg)
}
