package main

import (
	"fmt"

	decoder "github.com/Pois-Noir/Mammad/decoder"
	encoder "github.com/Pois-Noir/Mammad/encoder"
)

func main() {
	testEncoder := encoder.NewEncoder()

	testMap := map[string]interface{}{
		"keyssssrewrwerewr1": "value 123",
		"key2":               2,
		"key3":               3,
		"key4":               4.123,
		"key5":               true,
		"key41": map[string]interface{}{
			"key45": 1,
			"key46": 2,
			"key47": 3,
		},
	}

	byteStream, _ := testEncoder.EncodeMap(testMap)
	decoder := decoder.NewDecoderBytes(byteStream)
	result, err := decoder.Decode(len(byteStream))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)

}
