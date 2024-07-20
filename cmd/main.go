package main

import (
	"flag"
	"json-serde/pkg/deserializer"
	"json-serde/utils"
)

var fLong, fShort *string

func init() {
	fLong = flag.String("file", "", "path for json file")
	fShort = flag.String("f", "", "shorthand for -file or --file")
	flag.Parse()
}

func main() {
	utils.PrintLogo()
	file := utils.OpenFile(fLong, fShort)
	defer file.Close()

	var data any
	deserializer.Deserialize(file, &data)
	// dt, _ := io.ReadAll(file)
	// json.Unmarshal(dt, &data)
}
