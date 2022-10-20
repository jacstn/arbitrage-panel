package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadCharArr() []string {
	jsonFile, err := os.Open("./resources/charArr.json")
	if err != nil {

		fmt.Println(err)
		panic("cannot find charArr json file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var arr []string
	_ = json.Unmarshal([]byte(byteValue), &arr)

	defer jsonFile.Close()
	return arr
}
