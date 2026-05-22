package httputil

import (
	"bytes"
	"encoding/json"
)

type OptionalString struct {
	Set   bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if bytes.Equal(data, []byte("null")) {
		o.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	o.Value = &s
	return nil
}

type OptionalStringSlice struct {
	Set   bool
	Value []string
}

func (o *OptionalStringSlice) UnmarshalJSON(data []byte) error {
	o.Set = true
	if bytes.Equal(data, []byte("null")) {
		o.Value = nil
		return nil
	}
	var s []string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	o.Value = s
	return nil
}

