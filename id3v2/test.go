package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"../utils/crypt"
)

func nmain() {
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["musicId"] = 1304665120
	json, _:= json.Marshal(data)
	fmt.Println(string(json))
}

func main() {
	data, err := base64.StdEncoding.DecodeString(os.Args[1])
	fmt.Println(err)
	decrypted := crypt.MarkerAesDecryptECB([]byte(data))
	fmt.Printf(string(decrypted))
}
