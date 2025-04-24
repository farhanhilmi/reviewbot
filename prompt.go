package reviewbot

import "fmt"

func BuildPrompt(diff string) string {
	return fmt.Sprintf("You are a senior Go engineer. Review this code diff and suggest improvements:\n\n%s", diff)
}
