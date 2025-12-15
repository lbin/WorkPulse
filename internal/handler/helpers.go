package handler

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func parseUUID(input *string) *uuid.UUID {
	if input == nil || *input == "" {
		return nil
	}
	v := uuid.MustParse(*input)
	return &v
}

func mustJSON(v interface{}) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte(`[]`))
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte(`[]`))
	}
	return datatypes.JSON(b)
}
