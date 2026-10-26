package marshalutil

import (
	"encoding/json"
)

// RemarshalAny tries to unmarshal any value into another, by marshaling and unmarshaling the data. For example when
// needing to convert a map into a struct. May not be suitable for already byte-encoded data.
func RemarshalAny[T any](data any, target T) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
