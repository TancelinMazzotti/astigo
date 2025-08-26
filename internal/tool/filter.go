package tool

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Operation string

const (
	In       Operation = "in"
	Contains Operation = "cts"
	Subset   Operation = "subset"
	Superset Operation = "superset"

	Equals             Operation = "eq"
	NotEquals          Operation = "neq"
	LessThan           Operation = "lt"
	LessThanOrEqual    Operation = "lte"
	GreaterThan        Operation = "gt"
	GreaterThanOrEqual Operation = "gte"

	ExclusiveOr Operation = "xor"
	Or          Operation = "or"
	And         Operation = "and"
	Not         Operation = "not"
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
	case In, Contains, Subset, Superset, Equals, NotEquals,
		LessThan, LessThanOrEqual, GreaterThan, GreaterThanOrEqual,
		ExclusiveOr, Or, And, Not:
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
