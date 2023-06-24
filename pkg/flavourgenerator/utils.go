package flavourgenerator

import (
	"encoding/json"
	"fmt"
)

// marshallJson converts a NodeInfo struct to JSON
func marshallJson(flavour *Flavour) ([]byte, error) {
	jsonBody, err := json.Marshal(flavour)
	if err != nil {
		return nil, fmt.Errorf("error converting to JSON: %v", err)
	}
	return jsonBody, nil
}
