package PhoneNumber

import (
	"fmt"
	"strings"
)

type PhoneNumber struct {
	Id	   int
	Number string
}

func Normalize(dirtyString string) (string, error) {
	useString := strings.TrimSpace(dirtyString)
	if useString == ""{
		return "", fmt.Errorf("string cannot be empty")
	}	
	newString := ""

	for i := range dirtyString{
		char := dirtyString[i]
		if  char < 48 || char > 57 {
			continue
		}
		newString += string(dirtyString[i])
	}
	return newString, nil
}

