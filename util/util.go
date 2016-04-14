package util

import (
	"encoding/json"
	"fmt"
)

func PrintJson(v interface{}) error {
	jsonMsg, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(jsonMsg))
	return nil
}
