package utils

import "fmt"

func ShouldProceedInput(label string) bool {
	fmt.Print(label)
	var input string
	fmt.Scanln(&input)
	return input == "y"
}

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

func LogInfo(message string) {
	fmt.Println(Green + "[INFO] " + message + Reset)
}

func LogError(message string) {
	fmt.Println(Red + "[ERROR] " + message + Reset)
}
