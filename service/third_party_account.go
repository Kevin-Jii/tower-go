package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

const tsBeerLoginURL = "https://tp-api.tsbeer.com/api/identity/v2/user/token"
const tsBeerOrderURL = "https://tp-api.tsbeer.com/api/icommerceb-trade/v1/trade/order/applet/page"

type ThirdPartyAccountService struct {
	module      *module.ThirdPartyAccountModule
	orderModule *module.ThirdPartyOrderModule
	client      *http.Client
	debugMode   bool
}

func NewThirdPartyAccountService(m *module.ThirdPartyAccountModule, orderModule *module.ThirdPartyOrderModule) *ThirdPartyAccountService {
	return &ThirdPartyAccountService{
		module:      m,
		orderModule: orderModule,
		client:      &http.Client{Timeout: 15 * time.Second},
		debugMode:   isThirdPartySyncDebugEnabled(),
	}
}

func isThirdPartySyncDebugEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("THIRD_PARTY_SYNC_DEBUG")))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func (s *ThirdPartyAccountService) List(keyword string) ([]*model.ThirdPartyAccount, error) {
	return s.module.List(strings.TrimSpace(keyword))
}

func (s *ThirdPartyAccountService) GetByID(id uint) (*model.ThirdPartyAccount, error) {
	return s.module.GetByID(id)
}

func (s *ThirdPartyAccountService) ListSyncedOrders(accountID uint, page, pageSize int) ([]*model.ThirdPartyOrder, int64, error) {
	if _, err := s.module.GetByID(accountID); err != nil {
		return nil, 0, err
	}
	return s.orderModule.ListByAccount(accountID, page, pageSize)
}

func (s *ThirdPartyAccountService) Create(req *model.CreateThirdPartyAccountReq) (*model.ThirdPartyAccount, error) {
	row := &model.ThirdPartyAccount{
		PlatformName:   strings.TrimSpace(req.PlatformName),
		Name:           strings.TrimSpace(req.Name),
		LoginName:      strings.TrimSpace(req.LoginName),
		Phone:          strings.TrimSpace(req.Phone),
		Password:       strings.TrimSpace(req.Password),
		ApplicationKey: strings.TrimSpace(req.ApplicationKey),
		LoginType:      strings.TrimSpace(req.LoginType),
		Channel:        strings.TrimSpace(req.Channel),
		ShopID:         strings.TrimSpace(req.ShopID),
		CustomerID:     strings.TrimSpace(req.CustomerID),
		Remark:         strings.TrimSpace(req.Remark),
		IsEnabled:      true,
	}
	if row.LoginType == "" {
		row.LoginType = "2"
	}
	if row.Channel == "" {
		row.Channel = "WEB"
	}
	if req.IsEnabled != nil {
		row.IsEnabled = *req.IsEnabled
	}
	if err := s.module.Create(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *ThirdPartyAccountService) Update(id uint, req *model.UpdateThirdPartyAccountReq) error {
	updates := map[string]interface{}{}
	if req.PlatformName != nil {
		updates["platform_name"] = strings.TrimSpace(*req.PlatformName)
	}
	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.LoginName != nil {
		updates["login_name"] = strings.TrimSpace(*req.LoginName)
	}
	if req.Phone != nil {
		updates["phone"] = strings.TrimSpace(*req.Phone)
	}
	if req.Password != nil {
		updates["password"] = strings.TrimSpace(*req.Password)
	}
	if req.ApplicationKey != nil {
		updates["application_key"] = strings.TrimSpace(*req.ApplicationKey)
	}
	if req.LoginType != nil {
		updates["login_type"] = strings.TrimSpace(*req.LoginType)
	}
	if req.Channel != nil {
		updates["channel"] = strings.TrimSpace(*req.Channel)
	}
	if req.ShopID != nil {
		updates["shop_id"] = strings.TrimSpace(*req.ShopID)
	}
	if req.CustomerID != nil {
		updates["customer_id"] = strings.TrimSpace(*req.CustomerID)
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.Remark != nil {
		updates["remark"] = strings.TrimSpace(*req.Remark)
	}
	if len(updates) == 0 {
		return nil
	}
	return s.module.Update(id, updates)
}

func (s *ThirdPartyAccountService) Delete(id uint) error {
	return s.module.Delete(id)
}

func (s *ThirdPartyAccountService) TestLogin(id uint) (map[string]interface{}, error) {
	row, err := s.module.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !row.IsEnabled {
		return nil, errors.New("账号已禁用")
	}

	payload := map[string]string{
		"loginName": row.LoginName,
		"loginType": ifEmpty(row.LoginType, "2"),
		"phone":     ifEmpty(row.Phone, row.LoginName),
		"password":  row.Password,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, tsBeerLoginURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("application-key", row.ApplicationKey)
	req.Header.Set("channel", ifEmpty(row.Channel, "WEB"))
	req.Header.Set("isMock", "false")

	resp, err := s.client.Do(req)
	if err != nil {
		_ = s.module.Update(id, map[string]interface{}{
			"last_test_ok":  false,
			"last_test_msg": err.Error(),
			"last_test_at":  time.Now(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		_ = s.module.Update(id, map[string]interface{}{
			"last_test_ok":  false,
			"last_test_msg": "返回解析失败",
			"last_test_at":  time.Now(),
		})
		return nil, err
	}

	resultCode, _ := result["resultCode"].(string)
	resultMsg, _ := result["resultMsg"].(string)
	ok := resp.StatusCode >= 200 && resp.StatusCode < 300 && resultCode == "0"

	updates := map[string]interface{}{
		"last_test_ok":  ok,
		"last_test_msg": resultMsg,
		"last_test_at":  time.Now(),
	}

	if data, okData := result["data"].(map[string]interface{}); okData {
		if token, okToken := data["token"].(string); okToken {
			updates["last_token"] = token
		}
		switch v := data["tokenValidTime"].(type) {
		case float64:
			updates["token_valid_time"] = int64(v)
		case int64:
			updates["token_valid_time"] = v
		}
	}
	_ = s.module.Update(id, updates)
	if !ok {
		return result, errors.New(ifEmpty(resultMsg, "登录测试失败"))
	}
	return result, nil
}

func (s *ThirdPartyAccountService) SyncLatestOrders(id uint) (map[string]interface{}, error) {
	row, err := s.module.GetByID(id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(row.ShopID) == "" {
		return nil, errors.New("请先维护 shop_id")
	}
	customerID := ifEmpty(strings.TrimSpace(row.CustomerID), strings.TrimSpace(row.ShopID))
	token := strings.TrimSpace(row.LastToken)
	if token == "" {
		if _, err := s.TestLogin(id); err != nil {
			return nil, fmt.Errorf("登录失败，无法同步订单: %w", err)
		}
		row, err = s.module.GetByID(id)
		if err != nil {
			return nil, err
		}
		token = strings.TrimSpace(row.LastToken)
		if token == "" {
			return nil, errors.New("未获取到 access-token")
		}
	}

	now := time.Now()
	placeTimeEnd := now.Format("2006-01-02 15:04:05")
	placeTimeUp := now.AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
	if latestSynced, err := s.orderModule.GetLatestPlaceTimeByAccount(row.ID); err == nil && latestSynced != nil {
		// 留 2 分钟重叠，避免边界时钟误差漏单
		placeTimeUp = latestSynced.Add(-2 * time.Minute).Format("2006-01-02 15:04:05")
	}

	pageNum := 1
	pageSize := 10
	allList := make([]interface{}, 0)
	maxPages := 0

	for {
		u, _ := url.Parse(tsBeerOrderURL)
		q := u.Query()
		q.Set("pageSize", fmt.Sprintf("%d", pageSize))
		q.Set("pageNum", fmt.Sprintf("%d", pageNum))
		q.Set("shopId", row.ShopID)
		q.Set("itemName", "")
		q.Set("placeTimeUp", placeTimeUp)
		q.Set("placeTimeEnd", placeTimeEnd)
		u.RawQuery = q.Encode()
		if s.debugMode {
			fmt.Printf("[TP_SYNC_DEBUG] account=%d request pageNum=%d pageSize=%d shopId=%s placeTimeUp=%s placeTimeEnd=%s url=%s\n",
				row.ID, pageNum, pageSize, row.ShopID, placeTimeUp, placeTimeEnd, u.String())
		}

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "*/*")
		req.Header.Set("content-type", "application/json")
		req.Header.Set("application-key", row.ApplicationKey)
		req.Header.Set("channel", ifEmpty(row.Channel, "WEB"))
		req.Header.Set("isMock", "false")
		req.Header.Set("access-token", token)
		req.Header.Set("customerId", customerID)

		resp, err := s.client.Do(req)
		if err != nil {
			_ = s.module.Update(id, map[string]interface{}{
				"last_sync_at":    time.Now(),
				"last_sync_msg":   err.Error(),
				"last_sync_count": 0,
			})
			return nil, err
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		payload := result
		if d, ok := result["data"].(map[string]interface{}); ok && d != nil {
			payload = d
		}

		listAny, _ := payload["list"].([]interface{})
		allList = append(allList, listAny...)
		total := int(asFloat(payload["total"]))
		if total > 0 && pageSize > 0 {
			maxPages = (total + pageSize - 1) / pageSize
		}

		hasNext := false
		if v, ok := payload["hasNextPage"].(bool); ok {
			hasNext = v
		}
		if s.debugMode {
			fmt.Printf("[TP_SYNC_DEBUG] account=%d response pageNum=%d list=%d total=%d hasNextPage=%v maxPages=%d\n",
				row.ID, pageNum, len(listAny), total, hasNext, maxPages)
		}
		// 第三方 hasNextPage 字段偶尔不准确：优先用 total/pageSize 判断是否还有下一页
		if len(listAny) == 0 {
			break
		}
		if maxPages > 0 && pageNum >= maxPages {
			break
		}
		if maxPages == 0 && !hasNext {
			break
		}
		pageNum++
	}

	if len(allList) == 0 {
		now := time.Now()
		_ = s.module.Update(id, map[string]interface{}{
			"last_sync_at":    now,
			"last_sync_msg":   "未获取到订单数据",
			"last_sync_count": 0,
		})
		return map[string]interface{}{
			"synced_count": 0,
			"latest_date":  "",
			"message":      "未获取到订单数据",
		}, nil
	}

	rows := make([]model.ThirdPartyOrder, 0)
	latestDate := ""
	for _, item := range allList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		placeTimeStr, _ := m["placeTime"].(string)
		orderNo, _ := m["orderNo"].(string)
		if strings.TrimSpace(orderNo) == "" {
			continue
		}
		var placeTimePtr *time.Time
		placeDate := ""
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", placeTimeStr, time.Local); err == nil {
			placeTimePtr = &t
			placeDate = t.Format("2006-01-02")
		}
		if placeDate == "" && len(placeTimeStr) >= 10 {
			placeDate = placeTimeStr[:10]
		}
		if placeDate != "" && placeDate > latestDate {
			latestDate = placeDate
		}
		raw, _ := json.Marshal(m)
		rows = append(rows, model.ThirdPartyOrder{
			AccountID:        row.ID,
			PlatformName:     row.PlatformName,
			OrderNo:          orderNo,
			PlaceTime:        placeTimePtr,
			PlaceDate:        placeDate,
			OrderTradeStatus: asString(m["orderTradeStatus"]),
			StatusName:       asString(m["orderTradeStatusName"]),
			PayAmount:        asFloat(m["payAmount"]),
			TotalAmount:      asFloat(m["totalAmount"]),
			TotalItemNum:     asFloat(m["totalItemNum"]),
			RawJSON:          string(raw),
			SyncedAt:         now,
		})
	}
	if len(rows) == 0 {
		return nil, errors.New("未解析到有效订单")
	}

	if err := s.orderModule.UpsertBatch(rows); err != nil {
		_ = s.module.Update(id, map[string]interface{}{
			"last_sync_at":    now,
			"last_sync_msg":   err.Error(),
			"last_sync_count": 0,
		})
		return nil, err
	}

	msg := fmt.Sprintf("同步成功，时间范围 %s ~ %s，共 %d 单", placeTimeUp, placeTimeEnd, len(rows))
	_ = s.module.Update(id, map[string]interface{}{
		"last_sync_at":    now,
		"last_sync_msg":   msg,
		"last_sync_count": len(rows),
	})
	return map[string]interface{}{
		"synced_count": len(rows),
		"latest_date":  latestDate,
		"place_time_up": placeTimeUp,
		"place_time_end": placeTimeEnd,
		"message":      msg,
	}, nil
}

func ifEmpty(v string, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return v
}

func asString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func asFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case json.Number:
		f, _ := x.Float64()
		return f
	default:
		return 0
	}
}
