package waffle

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type Validator func(string) error

func PromptStr(prompt string) string {
	return Prompt(prompt, func(str string) error {
		for _, r := range str {
			switch {
			case ' ' == r:
			case '\t' == r:
			case '\r' == r:
			case '\n' == r:
			default:
				continue
			}
			return fmt.Errorf("Only strings are allowed (no spaces)")
		}
		return nil
	})
}

func Prompt(prompt string, validator Validator) string {
	p := promptui.Prompt{
		Label:    prompt,
		Validate: promptui.ValidateFunc(validator),
	}

	result, _ := p.Run()
	return result
}
