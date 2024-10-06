# json-serde

`json-serde` is an efficient serializer & deserializer for JSON data interchange format completely written in Go. This allow go's native data types (structs, maps, slices, etc.) to be transformed into JSON encoded string and vice versa.

This project adheres to JSON encode/decode specifications provided by [json.org](https://www.json.org/)
which itself references formal specification in [STD 90](https://www.rfc-editor.org/std/std90.txt) / [RFC 8259](https://www.rfc-editor.org/rfc/rfc8259.txt) / [ECMA-404](https://ecma-international.org/wp-content/uploads/ECMA-404_2nd_edition_december_2017.pdf).

### Benefits of our own JSON implementation

- Knowledge gain in writing our own tokenizer/lexer and parser/deserializer
- Challenge to achieve high level of efficiency (memory footprint) and performance
- Instead of naive implementation (which can be achieved in exponentially less time), one which can rival standard lib's json package
- It will be a good starting point as JSON has limited tokens & syntax
- This can be a first step in understanding lexer, parser & AST generation for fully featured programming languages
