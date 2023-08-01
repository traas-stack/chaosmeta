package json

import (
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
)

func MargePatch(originalObj, updatedObj interface{}) ([]byte, error) {
	originalJSON, err := json.Marshal(originalObj)
	if err != nil {
		return nil, err
	}

	updatedJSON, err := json.Marshal(updatedObj)
	if err != nil {
		return nil, err
	}

	data, err := jsonpatch.CreateMergePatch(originalJSON, updatedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to marge patch data, error: %s", err)
	}

	return data, nil
}
