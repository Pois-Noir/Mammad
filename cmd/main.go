package main

import (
	decoder "botzila/parser/decoder"
	encoder "botzila/parser/encoder"
	"fmt"
)

func main() {
	testEncoder := encoder.NewEncoder()

	testMap := map[string]interface{}{
		"keyssssrewrwerewr1": "value 123",
		"key2":               2,
		"key3":               3,
		"key4":               4.123,
		"key5":               true,
		// "key4": map[string]int{
		// 	"key45": 1,
		// 	"key46": 2,
		// 	"key47": 3,
		// },
	}

	byteStream, err := testEncoder.EncodeMap(testMap)
	decoder := decoder.NewDecoderBytes(byteStream)
	result, err := decoder.Decode()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

}
