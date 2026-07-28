package middleware

import "testing"

func TestAuditBodyBufferRespectsLimit(t *testing.T) {
	buffer := &auditBodyBuffer{limit: 4}
	if n, err := buffer.Write([]byte("123456")); err != nil || n != 6 {
		t.Fatalf("Write() = (%d, %v), want (6, nil)", n, err)
	}
	if got := buffer.String(); got != "1234" {
		t.Fatalf("buffer = %q, want %q", got, "1234")
	}

	if n, err := buffer.Write([]byte("78")); err != nil || n != 2 {
		t.Fatalf("Write() after limit = (%d, %v), want (2, nil)", n, err)
	}
	if got := buffer.String(); got != "1234" {
		t.Fatalf("buffer grew after limit: %q", got)
	}
}
