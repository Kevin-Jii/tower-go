package service

import (
	"reflect"
	"testing"

	"github.com/Kevin-Jii/tower-go/pkg/apicode"
)

func TestOptionalStoreReturnClientReqID(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  *string
	}{
		{name: "missing id remains nullable", value: "", want: nil},
		{name: "whitespace id remains nullable", value: "  \t ", want: nil},
		{name: "id is trimmed", value: "  return-20260814-1  ", want: stringPointer("return-20260814-1")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := optionalStoreReturnClientReqID(tt.value)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("optionalStoreReturnClientReqID() = %q, want nil", *got)
				}
				return
			}
			if got == nil || *got != *tt.want {
				t.Fatalf("optionalStoreReturnClientReqID() = %v, want %q", got, *tt.want)
			}
		})
	}
}

func TestNormalizeStoreReturnPhotos(t *testing.T) {
	want := []string{"https://tower.example/images/one.jpg", "https://tower.example/images/two.jpg"}
	got, err := normalizeStoreReturnPhotos([]string{
		" https://tower.example/images/one.jpg ",
		"https://tower.example/images/one.jpg",
		"https://tower.example/images/two.jpg",
	})
	if err != nil {
		t.Fatalf("normalizeStoreReturnPhotos() error = %v", err)
	}
	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("normalizeStoreReturnPhotos() = %#v, want %#v", got, want)
	}
}

func TestNormalizeStoreReturnPhotosRejectsInvalidInput(t *testing.T) {
	if _, err := normalizeStoreReturnPhotos([]string{"one", "two", "three", "four"}); !apicode.Is(err, apicode.ValidationFailed) {
		t.Fatalf("too many photos error = %v", err)
	}
	if _, err := normalizeStoreReturnPhotos([]string{"javascript:alert(1)"}); !apicode.Is(err, apicode.ValidationFailed) {
		t.Fatalf("invalid URL error = %v", err)
	}
}

func stringPointer(value string) *string {
	return &value
}
