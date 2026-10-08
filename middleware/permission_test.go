package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	apphttp "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

func TestPermissionOrRolesAllowsExplicitRoleWithoutPermissionLookup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/assistant",
		func(c *gin.Context) {
			c.Set("roleCode", model.RoleCodeStoreAdmin)
		},
		PermissionOrRoles("ai:assistant:use", model.RoleCodeStoreAdmin),
		func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		},
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/assistant", nil))

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestPermissionOrRolesStillRejectsOtherRolesWithoutIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/assistant",
		func(c *gin.Context) {
			c.Set("roleCode", model.RoleCodeStaff)
		},
		PermissionOrRoles("ai:assistant:use", model.RoleCodeStoreAdmin),
		func(c *gin.Context) {
			t.Fatal("protected handler should not run")
		},
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/assistant", nil))

	var body apphttp.Response
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != 40301 {
		t.Fatalf("code = %d, want 40301", body.Code)
	}
}
