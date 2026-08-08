package middleware

import (
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/gin-gonic/gin"
)

func TestIsStoreManager(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		roleCode string
		want     bool
	}{
		{roleCode: model.RoleCodeSuperAdmin, want: true},
		{roleCode: model.RoleCodeAdmin, want: true},
		{roleCode: model.RoleCodeStoreAdmin, want: true},
		{roleCode: model.RoleCodeStaff, want: false},
		{roleCode: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.roleCode, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			ctx.Set("roleCode", tt.roleCode)
			if got := IsStoreManager(ctx); got != tt.want {
				t.Fatalf("IsStoreManager() = %v, want %v", got, tt.want)
			}
		})
	}
}
