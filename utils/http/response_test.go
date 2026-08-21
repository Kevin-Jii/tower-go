package http

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type responseTimeFixture struct {
	AccountDate         time.Time              `json:"account_date"`
	CreatedAt           time.Time              `json:"created_at"`
	RecentConsumptionAt *time.Time             `json:"recent_consumption_at"`
	Nested              map[string]interface{} `json:"nested"`
}

func TestNormalizeResponseDataFormatsNestedTimes(t *testing.T) {
	createdAt := time.Date(2026, 8, 19, 15, 4, 5, 0, time.FixedZone("CST", 8*60*60))
	data := responseTimeFixture{
		AccountDate:         createdAt,
		CreatedAt:           createdAt,
		RecentConsumptionAt: &createdAt,
		Nested: map[string]interface{}{
			"updated_at": createdAt,
			"items":      []interface{}{map[string]interface{}{"order_date": createdAt}},
		},
	}

	normalized := normalizeResponseData(data)
	raw, err := json.Marshal(normalized)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	want := `{"account_date":"2026-08-19","created_at":"2026-08-19 15:04:05","recent_consumption_at":"2026-08-19 15:04:05","nested":{"updated_at":"2026-08-19 15:04:05","items":[{"order_date":"2026-08-19"}]}}`
	var gotValue interface{}
	var wantValue interface{}
	if err := json.Unmarshal(raw, &gotValue); err != nil {
		t.Fatalf("json.Unmarshal(got) error = %v", err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("json.Unmarshal(want) error = %v", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("normalized response = %s, want %s", raw, want)
	}
}

func TestNormalizeResponseDataKeepsNilTimeAndBytes(t *testing.T) {
	var at *time.Time
	normalized := normalizeResponseData(map[string]interface{}{
		"at":   at,
		"data": []byte("ok"),
	})
	raw, err := json.Marshal(normalized)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if string(raw) != `{"at":null,"data":"b2s="}` {
		t.Fatalf("normalized response = %s", raw)
	}
}

func TestSuccessFormatsResponseTimes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		t := time.Date(2026, 8, 19, 0, 0, 0, 0, time.Local)
		Success(c, map[string]interface{}{"created_at": t, "account_date": t})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), `"created_at":"2026-08-19 00:00:00"`) ||
		!strings.Contains(resp.Body.String(), `"account_date":"2026-08-19"`) {
		t.Fatalf("response = %s", resp.Body.String())
	}
}
