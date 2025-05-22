package main

import (
	encoder "botzila/parser/encoder"
	"fmt"
)

func main() {
	testEncoder := encoder.NewEncoder()

	testMap := map[string]interface{}{
		"key1": 1,
		"key2": 2,
		"key3": 3,
	}

	byteStream, err := testEncoder.EncodeMap(testMap)

	if err == nil {
		fmt.Println(byteStream)
	}

}
