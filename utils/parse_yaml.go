package main

import (
	flags "github.com/jessevdk/go-flags"
    "fmt"
	"gopkg.in/yaml.v2"
	"os"
    "reflect"
    "strings"
)

var opts struct {
	FileName string `short:"f" long:"file" description:"File name to read"`
	FieldName string `short:"k" long:"key" description:"Field name to match"`
	Value string `short:"v" long:"value" description:"Value to match"`
	Verbose bool `long:"verbose" description:"Verbose output"`
	Debug bool `long:"debug" description:"Debug output"`
}

func printVal(v interface{}, depth int, indent int, fieldName string, printIt bool, keyDepth int, prefix string) {
    typ := reflect.TypeOf(v).Kind()
    if typ == reflect.Int || typ == reflect.String {
	    if printIt && keyDepth >= 0 {
		    fmt.Printf("%s%s%v\n", strings.Repeat(" ", indent), prefix, v)
	    }
    } else if typ == reflect.Slice {
	    if printIt  && keyDepth >= 0 {
		    fmt.Printf("\n")
	    }
        printSlice(v.([]interface{}), depth+1, fieldName, printIt, keyDepth)
    } else if typ == reflect.Map {
	    if printIt && keyDepth >= 0  {
		    fmt.Printf("\n")
	    }
        printMap(v.(map[interface{}]interface{}), depth+1, fieldName, printIt, keyDepth)
    }

}

func printMap(m map[interface{}]interface{}, depth int, fieldName string, printIt bool, keyDepth int) {
	for k, v := range m {
		if k.(string) == fieldName {
			keyDepth = depth
			printIt = true
		}
		if k.(string) != fieldName && depth <= keyDepth {
			printIt = false
			keyDepth = -1
		}
	    if printIt {

		    fmt.Printf("%s%s:", strings.Repeat(" ", depth), k.(string))
		    printIt = true
	    }
		printVal(v, depth+1, 1, fieldName, printIt, keyDepth, "")
    }
	printIt = false
}

func printSlice(slc []interface{}, depth int, fieldName string, printIt bool, keyDepth int) {
    for _, v := range slc {
        printVal(v, depth+1, depth+1, fieldName, printIt, keyDepth, "- ")
    }
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flags.Parse(&opts)
	if opts.FileName == "" {
		fmt.Printf("Usage: parse_json filename\n")
		os.Exit(1)
	}
	byteValue, err := os.ReadFile(opts.FileName)
	check(err)

    yamlData := make(map[string]interface{})

    err = yaml.Unmarshal(byteValue, &yamlData)
    if err != nil {
        panic(err)
    }
	printIt := false
	keyDepth := -1
    for k, v := range yamlData {
	    if k == opts.FieldName {
		    keyDepth = 1
		    fmt.Printf("%s: ", k)
		    printIt = true
	    }
	    printVal(v, 1, 1, opts.FieldName, printIt, keyDepth, "")
    }
}
