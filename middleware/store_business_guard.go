package middleware

import (
	"net/http"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/utils/database"
	httpx "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

// StoreBusinessGuard 拦截停业门店的业务写操作（记账/采购/库存/会员等）。
// 仅对 POST/PUT/PATCH/DELETE 生效；查询类 GET 不受影响。
func StoreBusinessGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			// continue
		default:
			c.Next()
			return
		}

		storeID := GetStoreID(c)
		if storeID == 0 || database.DB == nil {
			c.Next()
			return
		}

		var store model.Store
		if err := database.DB.Select("id", "status").First(&store, storeID).Error; err != nil {
			httpx.ErrorApp(c, apicode.NotFound)
			c.Abort()
			return
		}

		if store.Status == 2 {
			httpx.ErrorApp(c, apicode.StoreClosed)
			c.Abort()
			return
		}

		c.Next()
	}
}

