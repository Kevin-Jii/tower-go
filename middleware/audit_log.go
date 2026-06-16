package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/clientsource"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuditLogMiddleware(maxBody int) gin.HandlerFunc {
	if maxBody <= 0 {
		maxBody = 4096
	}
	return func(c *gin.Context) {
		if !shouldAuditRequest(c) {
			c.Next()
			return
		}

		start := time.Now()
		body := readAuditBody(c, maxBody)
		blw := &auditBodyLogWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
		c.Writer = blw
		c.Next()

		log := buildAuditLog(c, body, blw.body.String(), time.Since(start))
		if log == nil || database.DB == nil {
			return
		}
		if err := database.DB.Create(log).Error; err != nil {
			logging.LogWarn("写入操作日志失败", zap.Error(err))
		}
	}
}

func shouldAuditRequest(c *gin.Context) bool {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/api/v1/audit-logs") {
		return false
	}
	if path == "/api/v1/auth/login" {
		return true
	}
	switch c.Request.Method {
	case "POST", "PUT", "PATCH", "DELETE":
		return strings.HasPrefix(path, "/api/v1/")
	default:
		return false
	}
}

func readAuditBody(c *gin.Context, maxBody int) string {
	if c.Request.Body == nil {
		return ""
	}
	if strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "multipart/form-data") {
		return "[multipart body omitted]"
	}
	data, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	if len(data) == 0 {
		return ""
	}
	if len(data) > maxBody {
		data = append(data[:maxBody], []byte("...")...)
	}
	return maskAuditBody(string(data))
}

type auditBodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *auditBodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *auditBodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func buildAuditLog(c *gin.Context, body string, responseBody string, latency time.Duration) *model.AuditLog {
	moduleCode, moduleName, resourceType := auditModuleFromPath(c.Request.URL.Path)
	actionCode, actionName := auditActionFromMethod(c.Request.Method, c.Request.URL.Path)
	status := model.AuditStatusSuccess
	if c.Writer.Status() >= http.StatusBadRequest || responseBusinessFailed(responseBody) {
		status = model.AuditStatusFail
	}

	resourceID := resourceIDFromPath(c.Request.URL.Path)
	userID := GetUserID(c)
	storeID := GetStoreID(c)
	username := getStringCtx(c, "username")
	roleCode := GetRoleCode(c)
	var nickname, phone, roleName, storeName string

	if database.DB != nil && userID > 0 {
		var user model.User
		if err := database.DB.Preload("Role").Preload("Store").First(&user, userID).Error; err == nil {
			username = user.Username
			nickname = user.Nickname
			phone = user.Phone
			storeID = user.StoreID
			if user.Role != nil {
				roleName = user.Role.Name
				roleCode = user.Role.Code
			}
			if user.Store != nil {
				storeName = user.Store.Name
			}
		}
	}

	if c.Request.URL.Path == "/api/v1/auth/login" && userID == 0 {
		phone = phoneFromAuditBody(body)
	}
	deviceType, osName, browser := parseDeviceFromUA(c.Request.UserAgent(), clientsource.FromRequest(c.Request))

	return &model.AuditLog{
		UserID:       userID,
		Username:     username,
		Nickname:     nickname,
		Phone:        phone,
		RoleName:     roleName,
		RoleCode:     roleCode,
		StoreID:      storeID,
		StoreName:    storeName,
		Module:       moduleCode,
		ModuleName:   moduleName,
		Action:       actionCode,
		ActionName:   actionName,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Method:       c.Request.Method,
		Path:         c.Request.URL.Path,
		Query:        c.Request.URL.RawQuery,
		RequestBody:  body,
		Status:       status,
		StatusCode:   c.Writer.Status(),
		ErrorMessage: auditErrorMessage(c, responseBody),
		LatencyMs:    latency.Milliseconds(),
		ClientIP:     c.ClientIP(),
		ClientSource: clientsource.FromRequest(c.Request),
		DeviceType:   deviceType,
		OS:           osName,
		Browser:      browser,
		UserAgent:    c.Request.UserAgent(),
	}
}

func auditModuleFromPath(path string) (string, string, string) {
	p := strings.TrimPrefix(path, "/api/v1/")
	seg := strings.Split(strings.Trim(p, "/"), "/")
	key := ""
	if len(seg) > 0 {
		key = seg[0]
	}
	m := map[string][2]string{
		"auth":                  {"auth", "认证"},
		"users":                 {"user", "用户管理"},
		"roles":                 {"role", "角色管理"},
		"permission":            {"permission", "权限管理"},
		"stores":                {"store", "门店管理"},
		"menus":                 {"menu", "菜单管理"},
		"dict-types":            {"dict", "字典管理"},
		"dict-data":             {"dict", "字典管理"},
		"suppliers":             {"supplier", "供应商管理"},
		"supplier-products":     {"supplier_product", "供应商商品"},
		"supplier-categories":   {"supplier_category", "供应商分类"},
		"product-unit-specs":    {"product_unit_spec", "商品规格"},
		"store-suppliers":       {"store_supplier", "门店供应商"},
		"purchase-orders":       {"purchase_order", "采购管理"},
		"inventories":           {"inventory", "库存管理"},
		"inventory-orders":      {"inventory_order", "库存订单"},
		"inventory-loss-orders": {"inventory_loss", "库存损耗"},
		"store-accounts":        {"store_account", "门店记账"},
		"store-returns":         {"store_return", "门店退货"},
		"members":               {"member", "会员管理"},
		"wallet-logs":           {"wallet_log", "钱包流水"},
		"recharge-orders":       {"recharge_order", "充值订单"},
		"b2b":                   {"b2b", "B2B"},
		"price-lists":           {"price_list", "价格清单"},
		"statistics":            {"statistics", "统计分析"},
		"meituan-ai":            {"meituan_ai", "美团 AI"},
		"dingtalk":              {"dingtalk", "钉钉"},
		"message-templates":     {"message_template", "消息模板"},
		"printers":              {"printer", "打印机"},
		"third-party-accounts":  {"third_party_account", "第三方账号"},
		"third-party-routes":    {"third_party_route", "第三方路线"},
		"files":                 {"file", "文件管理"},
		"galleries":             {"gallery", "图库管理"},
	}
	if v, ok := m[key]; ok {
		return v[0], v[1], v[0]
	}
	if key == "" {
		return "system", "系统", "system"
	}
	return strings.ReplaceAll(key, "-", "_"), key, strings.ReplaceAll(key, "-", "_")
}

func auditActionFromMethod(method string, path string) (string, string) {
	if path == "/api/v1/auth/login" {
		return model.AuditActionLogin, "登录"
	}
	switch method {
	case "POST":
		if strings.Contains(path, "import") {
			return "import", "导入"
		}
		if strings.Contains(path, "export") {
			return "export", "导出"
		}
		if strings.Contains(path, "print") {
			return "print", "打印"
		}
		if strings.Contains(path, "approve") || strings.Contains(path, "audit") {
			return "approve", "审核"
		}
		if strings.Contains(path, "sync") {
			return "sync", "同步"
		}
		return model.AuditActionCreate, "新增"
	case "PUT", "PATCH":
		return model.AuditActionUpdate, "修改"
	case "DELETE":
		return model.AuditActionDelete, "删除"
	default:
		return model.AuditActionOther, "操作"
	}
}

func resourceIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if _, err := strconv.ParseUint(parts[i], 10, 64); err == nil {
			return parts[i]
		}
	}
	return ""
}

func maskAuditBody(body string) string {
	if body == "" {
		return ""
	}
	var v interface{}
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return body
	}
	maskJSONValue(v)
	out, err := json.Marshal(v)
	if err != nil {
		return body
	}
	return string(out)
}

func maskJSONValue(v interface{}) {
	switch x := v.(type) {
	case map[string]interface{}:
		for k, val := range x {
			lk := strings.ToLower(k)
			if strings.Contains(lk, "password") ||
				strings.Contains(lk, "token") ||
				strings.Contains(lk, "secret") ||
				strings.Contains(lk, "access_key") ||
				strings.Contains(lk, "user_key") ||
				strings.Contains(lk, "authorization") {
				x[k] = "***"
				continue
			}
			maskJSONValue(val)
		}
	case []interface{}:
		for _, item := range x {
			maskJSONValue(item)
		}
	}
}

func phoneFromAuditBody(body string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return ""
	}
	if v, ok := m["phone"]; ok {
		return strings.TrimSpace(toAuditString(v))
	}
	return ""
}

func toAuditString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	default:
		return ""
	}
}

func responseBusinessFailed(body string) bool {
	if body == "" {
		return false
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return false
	}
	if v, ok := m["code"]; ok {
		switch x := v.(type) {
		case float64:
			return int(x) != http.StatusOK
		case int:
			return x != http.StatusOK
		}
	}
	return false
}

func auditErrorMessage(c *gin.Context, responseBody string) string {
	if len(c.Errors) > 0 {
		return c.Errors.String()
	}
	if responseBody != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(responseBody), &m); err == nil {
			if code, ok := m["code"].(float64); ok && int(code) != http.StatusOK {
				if msg, ok := m["message"].(string); ok {
					return msg
				}
			}
		}
	}
	if c.Writer.Status() >= 400 {
		return "请求失败"
	}
	return ""
}

func getStringCtx(c *gin.Context, key string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func parseDeviceFromUA(ua string, source string) (string, string, string) {
	lower := strings.ToLower(ua)
	deviceType := "desktop"
	if ua == "" {
		deviceType = "unknown"
	}
	if strings.Contains(lower, "bot") || strings.Contains(lower, "spider") || strings.Contains(lower, "crawler") {
		deviceType = "bot"
	} else if strings.Contains(lower, "ipad") || strings.Contains(lower, "tablet") {
		deviceType = "tablet"
	} else if strings.Contains(lower, "mobile") || strings.Contains(lower, "iphone") || strings.Contains(lower, "android") {
		deviceType = "mobile"
	}
	if strings.Contains(strings.ToLower(source), "weapp") || strings.Contains(strings.ToLower(source), "mini") {
		deviceType = "mobile"
	}

	osName := "Unknown"
	switch {
	case strings.Contains(lower, "windows nt 10"):
		osName = "Windows 10/11"
	case strings.Contains(lower, "windows"):
		osName = "Windows"
	case strings.Contains(lower, "mac os x") || strings.Contains(lower, "macintosh"):
		osName = "macOS"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") || strings.Contains(lower, "ios"):
		osName = "iOS"
	case strings.Contains(lower, "android"):
		osName = "Android"
	case strings.Contains(lower, "linux"):
		osName = "Linux"
	}

	browser := "Unknown"
	switch {
	case strings.Contains(lower, "edg/"):
		browser = "Edge"
	case strings.Contains(lower, "micromessenger"):
		browser = "WeChat"
	case strings.Contains(lower, "dingtalk"):
		browser = "DingTalk"
	case strings.Contains(lower, "chrome/") && !strings.Contains(lower, "chromium"):
		browser = "Chrome"
	case strings.Contains(lower, "firefox/"):
		browser = "Firefox"
	case strings.Contains(lower, "safari/") && strings.Contains(lower, "version/"):
		browser = "Safari"
	case strings.Contains(lower, "curl/"):
		browser = "curl"
	case strings.Contains(lower, "postman"):
		browser = "Postman"
	}
	if ua == "" && source != "" {
		browser = source
	}
	return deviceType, osName, browser
}
