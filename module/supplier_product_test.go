package module

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestIsForeignKeyReferenceError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "referenced row",
			err:  &mysql.MySQLError{Number: 1451, Message: "Cannot delete or update a parent row"},
			want: true,
		},
		{
			name: "wrapped referenced row",
			err:  errors.Join(errors.New("delete product"), &mysql.MySQLError{Number: 1451}),
			want: true,
		},
		{
			name: "duplicate key",
			err:  &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"},
			want: false,
		},
		{
			name: "generic database error",
			err:  errors.New("connection closed"),
			want: false,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isForeignKeyReferenceError(tt.err); got != tt.want {
				t.Fatalf("isForeignKeyReferenceError() = %v, want %v", got, tt.want)
			}
		})
	}
}
