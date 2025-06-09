package gateway

import (
	"encoding/json"
	"time"
)

type Login struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Handle null values
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}

	var ms float64
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}

	sec := int64(ms) / 1000
	nsec := int64((ms - float64(sec)*1000) * 1000000)

	t.Time = time.Unix(sec, nsec)
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.UnixMilli())
}
