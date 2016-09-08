// Various utilities in support on sparkcli
package util

import (
	"encoding/json"
	"fmt"
)

// PrintJson prints v with proper indents
func PrintJson(v interface{}) error {
	jsonMsg, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(jsonMsg))
	return nil
}
