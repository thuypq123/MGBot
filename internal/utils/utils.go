package utils

import (
	"strings"

	uuid "github.com/google/uuid"
)

func Uuid() string {
	return uuid.New().String()
}

func DeleteTextInString(text string, oldText string) string {
	return strings.Replace(text, oldText, "", -1)
}

// function to masking sensitive string
func MaskSensitiveString(text string) string {
	xString := "xxxxxx"
	return xString + text[len(text)-10:]
}
