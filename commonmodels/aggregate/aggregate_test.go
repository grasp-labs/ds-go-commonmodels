package aggregate_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/aggregate"
)

// The JSON names are the contract every service serves; changing one breaks
// every client.
func TestJSON(t *testing.T) {
	at := time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)
	bytes := int64(42)
	for name, tc := range map[string]struct {
		v    any
		want string
	}{
		"total":                {aggregate.Total{Count: 3, Bytes: &bytes, ModifiedAt: &at}, `{"count":3,"bytes":42,"modified_at":"2026-11-21T12:00:00Z"}`},
		"total without bytes":  {aggregate.Total{Count: 0}, `{"count":0,"modified_at":null}`},
		"group":                {aggregate.Group{Key: "reports/", Count: 2, Bytes: &bytes, ModifiedAt: &at}, `{"key":"reports/","count":2,"bytes":42,"modified_at":"2026-11-21T12:00:00Z"}`},
		"group without bytes":  {aggregate.Group{Key: "failed", Count: 1}, `{"key":"failed","count":1,"modified_at":null}`},
		"level":                {aggregate.Level{Name: "folder", Field: "key", Filter: "key_prefix", Bytes: true}, `{"name":"folder","field":"key","filter":"key_prefix","bytes":true}`},
		"level with intervals": {aggregate.Level{Name: "created-at", Field: "created_at", Filter: "created_at_gte", Intervals: []string{aggregate.IntervalHour, aggregate.IntervalDay}}, `{"name":"created-at","field":"created_at","filter":"created_at_gte","bytes":false,"intervals":["hour","day"]}`},
		"index":                {aggregate.Index{Levels: []aggregate.Level{}, OrderBy: aggregate.OrderBy}, `{"levels":[],"order_by":["key","count","bytes","modified_at"]}`},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("got  %s\nwant %s", got, tc.want)
			}
		})
	}
}
