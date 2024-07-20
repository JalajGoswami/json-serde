package deserializer

import (
	"errors"
	"fmt"
	"io"
	"json-serde/utils"
)

var typeNestingStack utils.Stack[utils.DataType]
var buffer = make([]byte, 4*1024)

type _deserializer struct {
	valueStartIndx int
	carryOverValue []byte
	value          []byte
	key            string
	result         any
}

func (s *_deserializer) handleDataType(dataType utils.DataType, currentIndx int) {
	if dataType == utils.StringType {
		last, isEmpty := typeNestingStack.Top()
		if !isEmpty && last == dataType {
			s.value = append(s.carryOverValue, buffer[s.valueStartIndx+1:currentIndx]...)
			s.valueStartIndx = -1
		} else {
			s.valueStartIndx = currentIndx
			typeNestingStack.Push(utils.StringType)
		}
	}
}

func (s *_deserializer) saveValue() {
	last, isEmpty := typeNestingStack.Top()
	if isEmpty {
		utils.LogError("no datatype found")
	}
	var value any
	switch last {
	case utils.StringType:
		value = string(s.value)
	}

	typeNestingStack.Pop()

}

func Deserialize(reader io.Reader, data any) {
	des := _deserializer{valueStartIndx: -1}
	n, err := reader.Read(buffer)
	for err == nil {
		for i := 0; i < len(buffer); i++ {
			ch := buffer[i]

			if isWhiteSpace(ch) {
				continue
			}

			dataType := scan(ch)
			des.handleDataType(dataType, i)

			des.saveValue()
		}
		n, err = reader.Read(buffer)
	}

	fmt.Println(n, err, err == io.EOF, errors.Is(err, io.EOF))
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' ||
		ch == '\n' ||
		ch == '\r' ||
		ch == '\t'
}

func scan(ch byte) utils.DataType {
	var dt = utils.NoType
	switch ch {
	case '"':
		dt = utils.StringType
	case '[', ']':
		dt = utils.ArrayType
	case '{', '}':
		dt = utils.ObjectType
	}
	return dt
}

func setValue(value []byte) {

}