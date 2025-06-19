package tool

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Operation string

const (
	IN                    Operation = "in"
	CONTAINS              Operation = "contains"
	SUBSET                Operation = "subset"
	SUPERSET              Operation = "superset"
	EQUALS                Operation = "eq"
	NOT_EQUALS            Operation = "neq"
	LESS_THAN             Operation = "lt"
	LESS_THAN_OR_EQUAL    Operation = "lte"
	GREATER_THAN          Operation = "gt"
	GREATER_THAN_OR_EQUAL Operation = "gte"

	EXCLUSIVE_OR Operation = "xor"
	OR           Operation = "or"
	AND          Operation = "and"
	NOT          Operation = "not"
)

type Filter struct {
	Field     string    `json:"field"`
	Operation Operation `json:"operation"`
	Items     []Filter  `json:"items"`
	Type      string    `json:"type"`
	Value     any       `json:"value"`
}

func (o Operation) IsValid() bool {
	switch o {
	case IN, CONTAINS, SUBSET, SUPERSET, EQUALS, NOT_EQUALS,
		LESS_THAN, LESS_THAN_OR_EQUAL, GREATER_THAN, GREATER_THAN_OR_EQUAL,
		EXCLUSIVE_OR, OR, AND, NOT:
		return true
	}
	return false
}

func (filter *Filter) UnmarshalJSON(data []byte) error {
	var tempFilter struct {
		Field     string          `json:"field"`
		Operation Operation       `json:"operation"`
		Items     []Filter        `json:"items"`
		Type      string          `json:"type"`
		Value     json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal(data, &tempFilter); err != nil {
		return err
	}

	typeHandlers := map[string]func() (interface{}, error){
		// BOOLEAN
		"bool":  func() (interface{}, error) { var v bool; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*bool": func() (interface{}, error) { var v *bool; err := json.Unmarshal(tempFilter.Value, &v); return v, err },

		// INT
		"int":    func() (interface{}, error) { var v int; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*int":   func() (interface{}, error) { var v *int; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"int8":   func() (interface{}, error) { var v int8; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*int8":  func() (interface{}, error) { var v *int8; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"int16":  func() (interface{}, error) { var v int16; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*int16": func() (interface{}, error) { var v *int16; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"int32":  func() (interface{}, error) { var v int32; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*int32": func() (interface{}, error) { var v *int32; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"int64":  func() (interface{}, error) { var v int64; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*int64": func() (interface{}, error) { var v *int64; err := json.Unmarshal(tempFilter.Value, &v); return v, err },

		// UINT
		"uint":    func() (interface{}, error) { var v uint; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*uint":   func() (interface{}, error) { var v *uint; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"uint8":   func() (interface{}, error) { var v uint8; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*uint8":  func() (interface{}, error) { var v *uint8; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"uint16":  func() (interface{}, error) { var v uint16; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*uint16": func() (interface{}, error) { var v *uint16; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"uint32":  func() (interface{}, error) { var v uint32; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*uint32": func() (interface{}, error) { var v *uint32; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"uint64":  func() (interface{}, error) { var v uint64; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*uint64": func() (interface{}, error) { var v *uint64; err := json.Unmarshal(tempFilter.Value, &v); return v, err },

		// FLOAT
		"float32": func() (interface{}, error) { var v float32; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*float32": func() (interface{}, error) {
			var v *float32
			err := json.Unmarshal(tempFilter.Value, &v)
			return v, err
		},
		"float64": func() (interface{}, error) { var v float64; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*float64": func() (interface{}, error) {
			var v *float64
			err := json.Unmarshal(tempFilter.Value, &v)
			return v, err
		},

		// STRING
		"string":  func() (interface{}, error) { var v string; err := json.Unmarshal(tempFilter.Value, &v); return v, err },
		"*string": func() (interface{}, error) { var v *string; err := json.Unmarshal(tempFilter.Value, &v); return v, err },

		// SPECIAL
		"uuid": func() (interface{}, error) {
			var v string
			if err := json.Unmarshal(tempFilter.Value, &v); err != nil {
				return nil, err
			}
			return uuid.Parse(v)
		},
		"*uuid": func() (interface{}, error) {
			var v *string
			if err := json.Unmarshal(tempFilter.Value, &v); err != nil {
				return nil, err
			}
			if v == nil {
				return NewPointerNil[uuid.UUID](), nil
			}
			parsed, err := uuid.Parse(*v)
			return &parsed, err
		},

		"time": func() (interface{}, error) {
			var v string
			if err := json.Unmarshal(tempFilter.Value, &v); err != nil {
				return nil, err
			}
			return time.Parse(time.RFC3339, v)
		},
		"*time": func() (interface{}, error) {
			var v *string
			if err := json.Unmarshal(tempFilter.Value, &v); err != nil {
				return nil, err
			}
			if v == nil {
				return NewPointerNil[time.Time](), nil
			}
			parsed, err := time.Parse(time.RFC3339, *v)
			return &parsed, err
		},

		"duration": func() (interface{}, error) {
			var v string
			if err := json.Unmarshal(tempFilter.Value, &v); err != nil {
				return nil, err
			}
			return time.ParseDuration(v)
		},
		"*duration": func() (interface{}, error) {
			var v *string
			if err := json.Unmarshal(tempFilter.Value, &v); err != nil {
				return nil, err
			}
			if v == nil {
				return NewPointerNil[time.Duration](), nil
			}
			parsed, err := time.ParseDuration(*v)
			return &parsed, err
		},
	}

	if handler, ok := typeHandlers[filter.Type]; ok {
		value, err := handler()
		if err != nil {
			return err
		}
		filter.Value = value
	} else {
		filter.Value = nil
	}

	return nil
}
