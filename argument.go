package wetgear

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Ugly as hell
var ArgumentRegex = regexp.MustCompile(`[^" \n\r\t` + "`" + `]+|"[^"\n\r\t` + "`" + `]*"`)

type Argument struct {
	raw    string
	quoted bool
}

func (a Argument) Empty() bool {
	return a.raw == ""
}

func (a Argument) Raw() string {
	return a.raw
}

func (a Argument) Quoted() bool {
	return a.quoted
}

func (a Argument) SurroundQuotes() string {
	return fmt.Sprintf(`"%s"`, a.raw)
}

// GetInt attempts to parse the argument as a int64. Upon failure, returns nil
func (a Argument) GetInt() *int64 {
	if i, err := strconv.ParseInt(a.raw, 10, 64); err == nil {
		return &i
	} else {
		return nil
	}
}

// GetFloat attempts to parse the argument as a float64. Upon failure, returns nil
func (a Argument) GetFloat() *float64 {
	if i, err := strconv.ParseFloat(a.raw, 64); err == nil {
		return &i
	} else {
		return nil
	}
}

// GetUint attempts to parse the argument as a uint64. Upon failure, returns nil
func (a Argument) GetUint() *uint64 {
	if i, err := strconv.ParseUint(a.raw, 10, 64); err == nil {
		return &i
	} else {
		return nil
	}
}

// GetBool attempts to parse the argument as a bool. Upon failure, returns nil
func (a Argument) GetBool() *bool {
	if i, err := strconv.ParseBool(a.raw); err == nil {
		return &i
	} else {
		return nil
	}
}

// ParseArguments takes in an input string and parses out valid arguments for use in commands
func ParseArguments(content string) []Argument {
	if content == "" {
		return nil
	}

	arguments := make([]Argument, 0)
	for _, match := range ArgumentRegex.FindAllString(content, -1) {
		quoted := strings.HasPrefix(match, `"`) && strings.HasSuffix(match, `"`)
		if quoted {
			runes := []rune(match)
			match = string(runes[1 : len(runes)-1])
		}
		arguments = append(arguments, Argument{match, quoted})
	}

	return arguments
}
