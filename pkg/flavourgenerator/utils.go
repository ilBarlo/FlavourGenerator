package flavourgenerator

import (
	"encoding/json"
	"log"
)

// marshallJson converts in json a NodeInfo struct
func marshallJson(node *NodeInfo) ([]byte, error) {
	jsonBody, err := json.Marshal(node)
	if err != nil {
		log.Fatalf("Errore durante la conversione in JSON: %v", err)
		return nil, err
	}
	return jsonBody, nil
}
