package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print(os.Getenv("THUY"))
	text := "1234567891011100"
	fmt.Print(MaskSensitiveString(text))
}

func MaskSensitiveString(text string) string {
	xString := "xxxxxx"
	return xString + text[len(text)-10:]
}
