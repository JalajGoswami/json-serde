package main

import (
	"flag"
	"fmt"
	"json-serde/pkg/tokenizer"
	"json-serde/utils"
)

var fLong, fShort *string

func init() {
	fLong = flag.String("file", "", "path for json file")
	fShort = flag.String("f", "", "shorthand for -file or --file")
	*fLong = "input.json"
	flag.Parse()
}

func main() {
	utils.PrintLogo()
	file := utils.OpenFile(fLong, fShort)
	defer file.Close()
	// var buffer = make([]byte, 3)
	// n, err := file.Read(buffer)
	// fmt.Println(n, err, string(buffer), buffer[1] == '\\')

	// buffer = make([]byte, 3)
	// n, err = file.Read(buffer)
	// fmt.Println(n, err, string(buffer))
	tokenGenerator := tokenizer.NewTokenizer(file, tokenizer.TokenizerConfig{BufferLen: 4})
	token, err := tokenGenerator.Next()
	for err == nil {
		if token != nil {
			fmt.Println(
				fmt.Sprintf("%-8s", token.TokenType),
				fmt.Sprintf("(%v bytes)", len(token.Value)),
				string(token.Value),
			)
		} else {
			panic("No token found even when there is no Error !")
		}
		token, err = tokenGenerator.Next()
	}
	fmt.Println(err)
}
