package service

import "testing"

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

func stringPointer(value string) *string {
	return &value
}
