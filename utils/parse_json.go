package main

// Read an arbitrary JSON file
// Find a field/property name and/or a value
import (
	"encoding/json"
	flags "github.com/jessevdk/go-flags"
	"fmt"
	"math"
	"os"
	"strconv"
	// "strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var opts struct {
	FileName string `short:"f" long:"file" description:"File name to read"`
	FieldName string `short:"k" long:"key" description:"Field name to match"`
	Value string `short:"v" long:"value" description:"Value to match"`
	Verbose bool `long:"verbose" description:"Verbose output"`
	Debug bool `long:"debug" description:"Debug output"`
}

func main() {
	flags.Parse(&opts)
	if opts.FileName == "" {
		fmt.Printf("Usage: parse_json filename\n")
		os.Exit(1)
	}
	fileName := opts.FileName
	findValue := opts.Value
	fieldMatch := opts.FieldName

	byteValue, err := os.ReadFile(fileName)
	check(err)

	var f interface{}
	err = json.Unmarshal(byteValue, &f)
	check(err)
	// fmt.Printf("findValue: '%s' fieldMatch: '%s'\n", findValue, fieldMatch)

	findJSON("", "", findValue, fieldMatch, f, false)
}

func findJSON(fieldName string, indent string, findValue string, fieldMatch string, v interface{}, isMatchObj bool) {
	orig_indent := indent
	isValue := findValue != ""
	isField := fieldMatch != ""

	if opts.Debug {
		fmt.Printf("Debug fieldMatch: %s findValue: %s fieldName: %s isValue: %t isField: %t Type: %T\n", fieldMatch, findValue, fieldName, isValue, isField, v)
	}

	switch vv := v.(type) {
	case string:
		if MatchString(vv, fieldMatch, findValue, fieldName, isValue, isField) || isMatchObj {
			fmt.Printf("%s\"%s\": \"%s\",\n", indent, fieldName, vv)
		}
	case float64:

		if vv != math.Round(vv) {

			if MatchFloat(vv, fieldMatch, findValue, fieldName, isValue, isField) || isMatchObj {
				fmt.Printf("%s\"%s\": %.3f,\n", indent, fieldName, vv)
			}
		} else {
			if MatchInt(int(vv), fieldMatch, findValue, fieldName, isValue, isField) || isMatchObj {
				fmt.Printf("%s\"%s\": %d,\n", indent, fieldName, int(vv))
			}
		}
	case []interface{}:
		if isMatchObj || (isField && fieldName == fieldMatch) {
			fmt.Printf("%s\"%s\": [\n", indent, fieldName)
			isMatchObj = true
		}
		for _, u := range vv {
			indent += "  "
			findJSON("", indent, findValue, fieldMatch, u, isMatchObj)
			indent = orig_indent
		}
		if isMatchObj {
			fmt.Printf("%s],\n", indent)
		}
		isMatchObj = false
	case map[string]interface{}:
		if isMatchObj || (isField && fieldName == fieldMatch) {
			if fieldName == fieldMatch {
				fmt.Printf("%s\"%s\": {\n", indent, fieldName)
			} else {
				fmt.Printf("%s{\n", indent)
			}
			isMatchObj = true
		}
		for i, u := range vv {
			indent += "  "
			findJSON(i, indent, findValue, fieldMatch, u, isMatchObj)
			indent = orig_indent
		}
		if isMatchObj {
			fmt.Printf("%s},\n", indent)
		}
		isMatchObj = false
	case nil:
		if findValue == "null" {
			fmt.Printf("\"%s\": null\n", fieldName)
		}
	case bool:
		fvBool, err := strconv.ParseBool(findValue)
		if err == nil && fvBool == vv {
			fmt.Printf("\"%s\": %t\n", fieldName, vv)
		}
	default:
		fmt.Printf("\"%s\": Unknown type: %T\n", fieldName, vv)
	}
}

func MatchString(value string, fieldMatch string, findValue string, fieldName string, isValue bool, isField bool) bool {
	isMatch := false

	if isValue && value == findValue && isField && fieldName == fieldMatch {
		isMatch = true
	} else if !isField && isValue && value == findValue {
		isMatch = true
	} else if !isValue && isField && fieldName == fieldMatch {
		isMatch = true
	}

	if isMatch && opts.Verbose {
		fmt.Printf("MatchString: value: %s fieldMatch: %s findValue: %s fieldName: %s isValue: %t isField: %t\n", value, fieldMatch, findValue, fieldName, isValue, isField)
	}

	return isMatch
}

func MatchFloat(value float64, fieldMatch string, findValue string, fieldName string, isValue bool, isField bool) bool {
	isMatch := false

	fvFloat, err := strconv.ParseFloat(findValue, 64)

	if isValue && err != nil {
		return false
	}

	if isValue && value == fvFloat && isField && fieldName == fieldMatch {
		isMatch = true
	} else if !isField && value == fvFloat {
		isMatch = true
	} else if !isValue && fieldName == fieldMatch {
		isMatch = true
	}

	if isMatch && opts.Verbose {
		fmt.Printf("MatchFloat: value: %.3f fieldMatch: %s findValue: %s fieldName: %s isValue: %t isField: %t\n", value, fieldMatch, findValue, fieldName, isValue, isField)
	}

	return isMatch
}

func MatchInt(value int, fieldMatch string, findValue string, fieldName string, isValue bool, isField bool) bool {
	isMatch := false
	fvInt, err := strconv.Atoi(findValue)
	if isValue && err != nil {
		return false
	}

	if isValue && value == fvInt && isField && fieldName == fieldMatch {
		isMatch = true
	} else if !isField && value == fvInt {
		isMatch = true
	} else if !isValue && fieldName == fieldMatch {
		isMatch = true
	}
	if isMatch && opts.Verbose {
		fmt.Printf("MatchInt: value: %d fieldMatch: %s findValue: %s fieldName: %s isValue: %t isField: %t\n", value, fieldMatch, findValue, fieldName, isValue, isField)
	}
	return isMatch
}

func printJSON(indent string, v interface{}) {
	orig_indent := indent
	switch vv := v.(type) {
	case string:
		fmt.Printf("\"%s\",\n", vv)
	case float64:
		if vv == math.Round(vv) {
			fmt.Printf("%d,\n", int(vv))
		} else {
			fmt.Printf("%f,\n", vv)
		}
	case []interface{}:
		// fmt.Println("is an array:")
		for _, u := range vv {
			fmt.Printf(" [\n")
			indent += "  "
			printJSON(indent, u)
			fmt.Printf(" ]\n")
			indent = orig_indent
		}
	case map[string]interface{}:
		// fmt.Println("is an object:")
		for i, u := range vv {
			fmt.Printf("%s\"%s\": {\n", indent, i)
			indent += "  "
			printJSON(indent, u)
			fmt.Printf("%s}\n", indent)
			indent = orig_indent
		}
	case nil:
		fmt.Println("null")
	default:
		fmt.Println("Unknown type")
	}
}
