package apicode

import (
	"errors"
	"testing"
)

func TestResolveWrappedError(t *testing.T) {
	err := Wrap(InventoryInsufficient, errors.New("quantity is not enough"))
	code, ok := Resolve(err)
	if !ok {
		t.Fatal("expected a recognized application error")
	}
	if code.Num != InventoryInsufficient.Num || code.Msg != InventoryInsufficient.Msg {
		t.Fatalf("unexpected code: %+v", code)
	}
	if !Is(err, InventoryInsufficient) {
		t.Fatal("expected Is to match the application code")
	}
}

func TestResolveDirectCode(t *testing.T) {
	code, ok := Resolve(DictTypeAlreadyExists)
	if !ok || code.Num != DictTypeAlreadyExists.Num {
		t.Fatalf("unexpected direct code resolution: %+v, %v", code, ok)
	}
}

func TestUnknownErrorFallsBackToInternal(t *testing.T) {
	code, ok := Resolve(errors.New("database details"))
	if ok || code.Num != InternalError.Num {
		t.Fatalf("unexpected fallback: %+v, %v", code, ok)
	}
}

func TestHTTPStatus(t *testing.T) {
	for code, expected := range map[int]int{
		AuthHeaderRequired.Num:    401,
		DictTypeAlreadyExists.Num: 409,
		ValidationFailed.Num:      422,
		InternalError.Num:         500,
		400:                       400,
	} {
		if actual := HTTPStatus(code); actual != expected {
			t.Fatalf("HTTPStatus(%d) = %d, want %d", code, actual, expected)
		}
	}
}
