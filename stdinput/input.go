package stdinput

import (
	"fmt"
)

func AskQuestion(question string) string {
	fmt.Print(question)

	var answer string
	fmt.Scanln(&answer)
	return answer
}

func AskQuestionWithDefault(question string, defaultAnswer string) string {
	answer := AskQuestion(question)
	if len(answer) > 0 {
		return answer
	}

	return defaultAnswer
}
