package service

import (
	"bytes"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type MeituanAIService struct {
	module *module.MeituanAIModule
	client *http.Client
}

func NewMeituanAIService(m *module.MeituanAIModule) *MeituanAIService {
	return &MeituanAIService{
		module: m,
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

func (s *MeituanAIService) ListAccounts(storeID uint, hqUnbound bool) ([]*model.MeituanAIOperatorAccount, error) {
	return s.module.ListAccounts(storeID, hqUnbound)
}

func (s *MeituanAIService) CreateAccount(storeID uint, hqUnbound bool, req *model.CreateMeituanAIAccountReq) (*model.MeituanAIOperatorAccount, error) {
	realStoreID := storeID
	if hqUnbound && req.StoreID > 0 {
		realStoreID = req.StoreID
	}
	if realStoreID == 0 {
		return nil, fmt.Errorf("请选择门店")
	}
	row := &model.MeituanAIOperatorAccount{
		StoreID:        realStoreID,
		ShopName:       strings.TrimSpace(req.ShopName),
		ShopID:         strings.TrimSpace(req.ShopID),
		LoginName:      strings.TrimSpace(req.LoginName),
		DeveloperID:    strings.TrimSpace(req.DeveloperID),
		SignKey:        strings.TrimSpace(req.SignKey),
		AppAuthToken:   strings.TrimSpace(req.AppAuthToken),
		BusinessID:     req.BusinessID,
		APIVersion:     strings.TrimSpace(req.APIVersion),
		APIBaseURL:     strings.TrimSpace(req.APIBaseURL),
		QueryOrderPath: strings.TrimSpace(req.QueryOrderPath),
		AuthStatus:     "manual",
		IsEnabled:      true,
		Remark:         strings.TrimSpace(req.Remark),
	}
	normalizeMeituanAccountConfig(row)
	if req.IsEnabled != nil {
		row.IsEnabled = *req.IsEnabled
	}
	if err := s.module.CreateAccount(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *MeituanAIService) UpdateAccount(id, storeID uint, hqUnbound bool, req *model.UpdateMeituanAIAccountReq) error {
	updates := map[string]interface{}{}
	if strings.TrimSpace(req.ShopName) != "" {
		updates["shop_name"] = strings.TrimSpace(req.ShopName)
	}
	if strings.TrimSpace(req.ShopID) != "" {
		updates["shop_id"] = strings.TrimSpace(req.ShopID)
	}
	if strings.TrimSpace(req.LoginName) != "" {
		updates["login_name"] = strings.TrimSpace(req.LoginName)
	}
	if strings.TrimSpace(req.DeveloperID) != "" {
		updates["developer_id"] = strings.TrimSpace(req.DeveloperID)
	}
	if strings.TrimSpace(req.SignKey) != "" {
		updates["sign_key"] = strings.TrimSpace(req.SignKey)
	}
	if strings.TrimSpace(req.AppAuthToken) != "" {
		updates["app_auth_token"] = strings.TrimSpace(req.AppAuthToken)
	}
	if req.BusinessID > 0 {
		updates["business_id"] = req.BusinessID
	}
	if strings.TrimSpace(req.APIVersion) != "" {
		updates["api_version"] = strings.TrimSpace(req.APIVersion)
	}
	if strings.TrimSpace(req.APIBaseURL) != "" {
		updates["api_base_url"] = strings.TrimSpace(req.APIBaseURL)
	}
	if strings.TrimSpace(req.QueryOrderPath) != "" {
		updates["query_order_path"] = strings.TrimSpace(req.QueryOrderPath)
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if strings.TrimSpace(req.Remark) != "" {
		updates["remark"] = strings.TrimSpace(req.Remark)
	}
	if len(updates) == 0 {
		return nil
	}
	return s.module.UpdateAccount(id, storeID, hqUnbound, updates)
}

func (s *MeituanAIService) ImportOrders(accountID, storeID uint, hqUnbound bool, req *model.ImportMeituanAIOrdersReq) (map[string]interface{}, error) {
	account, err := s.module.GetAccount(accountID, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	orders, err := buildMeituanAIOrders(req.Orders)
	if err != nil {
		return nil, err
	}
	n, err := s.module.UpsertOrders(account, orders)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"imported": n}, nil
}

func (s *MeituanAIService) SyncOrdersFromFile(accountID, storeID uint, hqUnbound bool, filename string, reader io.Reader) (*model.SyncMeituanAIOrdersResp, error) {
	account, err := s.module.GetAccount(accountID, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	items, skipped, err := parseMeituanAIOrderFile(filename, reader)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("文件里没有识别到有效订单")
	}
	orders, err := buildMeituanAIOrders(items)
	if err != nil {
		return nil, err
	}
	n, err := s.module.UpsertOrders(account, orders)
	if err != nil {
		return nil, err
	}
	return &model.SyncMeituanAIOrdersResp{Imported: n, Skipped: skipped}, nil
}

func (s *MeituanAIService) SyncOrdersFromOpenAPI(accountID, storeID uint, hqUnbound bool, req *model.SyncMeituanAIOpenAPIOrdersReq) (*model.SyncMeituanAIOrdersResp, error) {
	account, err := s.module.GetAccount(accountID, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	normalizeMeituanAccountConfig(account)
	orderIDs := cleanMeituanOrderIDs(req.OrderIDs)
	if strings.TrimSpace(req.OrderID) != "" {
		orderIDs = append(orderIDs, strings.TrimSpace(req.OrderID))
	}
	orderIDs = dedupeStrings(orderIDs)
	if len(orderIDs) == 0 {
		return nil, fmt.Errorf("请输入需要同步的美团订单号")
	}
	orders := make([]model.MeituanAIOrder, 0, len(orderIDs))
	skipped := 0
	for _, orderID := range orderIDs {
		row, err := s.fetchMeituanOrderByID(account, orderID)
		if err != nil {
			skipped++
			continue
		}
		orders = append(orders, row)
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("没有同步到有效订单，请检查开放平台凭证、订单号和接口权限")
	}
	n, err := s.module.UpsertOrders(account, orders)
	if err != nil {
		return nil, err
	}
	return &model.SyncMeituanAIOrdersResp{Imported: n, Skipped: skipped}, nil
}

func (s *MeituanAIService) ImportReviews(accountID, storeID uint, hqUnbound bool, req *model.ImportMeituanAIReviewsReq) (map[string]interface{}, error) {
	account, err := s.module.GetAccount(accountID, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	reviews := make([]model.MeituanAIReview, 0, len(req.Reviews))
	for _, item := range req.Reviews {
		t, err := parseFlexibleTime(item.ReviewTime)
		if err != nil {
			return nil, fmt.Errorf("评价%s时间格式错误", item.ReviewID)
		}
		sentiment, tags := classifyReview(item.Rating, item.Content)
		reviewID := strings.TrimSpace(item.ReviewID)
		if reviewID == "" {
			reviewID = fmt.Sprintf("%s-%d", strings.TrimSpace(item.OrderNo), t.Unix())
		}
		reviews = append(reviews, model.MeituanAIReview{
			ReviewID:       reviewID,
			OrderNo:        strings.TrimSpace(item.OrderNo),
			Rating:         item.Rating,
			Content:        strings.TrimSpace(item.Content),
			Sentiment:      sentiment,
			Tags:           strings.Join(tags, ","),
			SuggestedReply: buildReviewReply(item.Rating, item.Content, tags),
			ReviewTime:     t,
			ReplyStatus:    "pending",
		})
	}
	n, err := s.module.UpsertReviews(account, reviews)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"imported": n}, nil
}

func (s *MeituanAIService) ListOrders(storeID uint, hqUnbound bool, req *model.ListMeituanAIReq) ([]*model.MeituanAIOrder, int64, error) {
	if !hqUnbound {
		req.StoreID = storeID
	}
	if req.StoreID == 0 {
		return nil, 0, fmt.Errorf("请选择门店")
	}
	return s.module.ListOrders(req)
}

func (s *MeituanAIService) ListReviews(storeID uint, hqUnbound bool, req *model.ListMeituanAIReq) ([]*model.MeituanAIReview, int64, error) {
	if !hqUnbound {
		req.StoreID = storeID
	}
	if req.StoreID == 0 {
		return nil, 0, fmt.Errorf("请选择门店")
	}
	return s.module.ListReviews(req)
}

func (s *MeituanAIService) ListSuggestions(storeID uint, hqUnbound bool, req *model.ListMeituanAIReq) ([]*model.MeituanAISuggestion, int64, error) {
	if !hqUnbound {
		req.StoreID = storeID
	}
	if req.StoreID == 0 {
		return nil, 0, fmt.Errorf("请选择门店")
	}
	return s.module.ListSuggestions(req)
}

func (s *MeituanAIService) UpdateSuggestionStatus(id, storeID uint, req *model.UpdateMeituanAISuggestionStatusReq) error {
	if storeID == 0 {
		return fmt.Errorf("请选择门店")
	}
	return s.module.UpdateSuggestionStatus(id, storeID, req.Status)
}

func (s *MeituanAIService) Dashboard(storeID uint, hqUnbound bool, req *model.ListMeituanAIReq) (*model.MeituanAIDashboard, error) {
	if !hqUnbound {
		req.StoreID = storeID
	}
	if req.StoreID == 0 {
		return nil, fmt.Errorf("请选择门店")
	}
	return s.module.Dashboard(req.StoreID, req.AccountID, req.StartDate, req.EndDate)
}

func (s *MeituanAIService) GenerateSuggestions(accountID, storeID uint, hqUnbound bool, req *model.ListMeituanAIReq) (map[string]interface{}, error) {
	account, err := s.module.GetAccount(accountID, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	req.StoreID = account.StoreID
	req.AccountID = account.ID
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 100
	}
	dash, err := s.module.Dashboard(req.StoreID, req.AccountID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	orders, _, _ := s.module.ListOrders(req)
	reviews, _, _ := s.module.ListReviews(req)

	now := time.Now()
	suggestions, aiErr := s.generateDeepSeekSuggestions(req.StoreID, req.AccountID, dash, orders, reviews, now)
	usedAI := aiErr == nil && len(suggestions) > 0
	source := "deepseek"
	if !usedAI {
		suggestions = s.generateRuleSuggestions(req.StoreID, req.AccountID, dash, orders, reviews, now)
		source = "rules"
	}

	if err := s.module.ClearPendingSuggestions(req.StoreID, req.AccountID); err != nil {
		return nil, err
	}
	if err := s.module.CreateSuggestions(suggestions); err != nil {
		return nil, err
	}
	return map[string]interface{}{"generated": len(suggestions), "ai_enabled": usedAI, "source": source}, nil
}

func (s *MeituanAIService) generateRuleSuggestions(storeID, accountID uint, dash *model.MeituanAIDashboard, orders []*model.MeituanAIOrder, reviews []*model.MeituanAIReview, now time.Time) []model.MeituanAISuggestion {
	suggestions := make([]model.MeituanAISuggestion, 0)
	if dash.OrderCount == 0 {
		suggestions = append(suggestions, makeSuggestion(storeID, accountID, "data", "先导入美团订单数据", "当前周期没有订单数据，AI无法判断爆品、客单价和活动效果。", "从美团商家后台导出订单后，在本模块导入。", 90, now))
	}
	if dash.NegativeRate >= 15 {
		suggestions = append(suggestions, makeSuggestion(storeID, accountID, "review", "差评率偏高，优先处理差评原因", fmt.Sprintf("当前差评率 %.1f%%，已经影响店铺转化。", dash.NegativeRate), "优先回复低分评价，并把高频问题拆成包装、配送、口味、缺货四类处理。", 88, now))
	}
	if dash.AvgOrderAmount > 0 && dash.AvgOrderAmount < 60 {
		suggestions = append(suggestions, makeSuggestion(storeID, accountID, "bundle", "设计高客单套餐", fmt.Sprintf("当前客单价 %.2f，适合通过组合套餐提升单均收入。", dash.AvgOrderAmount), "选择销量最高商品，搭配杯具、小食或第二规格，做成高毛利套餐。", 78, now))
	}
	if dash.PlatformFee > 0 && dash.SalesAmount > 0 && dash.PlatformFee/dash.SalesAmount >= 0.08 {
		suggestions = append(suggestions, makeSuggestion(storeID, accountID, "profit", "复核平台费用和活动力度", fmt.Sprintf("平台费用占销售额 %.1f%%。", dash.PlatformFee/dash.SalesAmount*100), "检查满减、配送费补贴和平台服务费，避免活动后利润被吃掉。", 82, now))
	}
	if product := topProductName(orders); product != "" {
		suggestions = append(suggestions, makeSuggestion(storeID, accountID, "product", "强化爆品展示", "当前周期已有明显高频商品。", fmt.Sprintf("把「%s」放到美团店铺靠前位置，标题加容量/场景词，并搭配套餐入口。", product), 75, now))
	}
	for _, r := range reviews {
		if r.Rating <= 3 {
			title := "生成差评回复草稿"
			content := strings.TrimSpace(r.SuggestedReply)
			if content == "" {
				content = buildReviewReply(r.Rating, r.Content, strings.Split(r.Tags, ","))
			}
			suggestions = append(suggestions, makeSuggestion(storeID, accountID, "reply", title, fmt.Sprintf("订单%s存在%d星评价：%s", r.OrderNo, r.Rating, truncateText(r.Content, 80)), content, 72, now))
			if len(suggestions) >= 8 {
				break
			}
		}
	}
	if len(suggestions) == 0 {
		suggestions = append(suggestions, makeSuggestion(storeID, accountID, "routine", "保持当前运营节奏", "当前周期未发现明显风险。", "继续保持商品库存、评价回复和活动利润监控；建议每日导入订单评价刷新一次。", 55, now))
	}
	return suggestions
}

type deepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepSeekRequest struct {
	Model          string            `json:"model"`
	Messages       []deepSeekMessage `json:"messages"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
	Thinking       map[string]string `json:"thinking,omitempty"`
	Temperature    float64           `json:"temperature"`
	MaxTokens      int               `json:"max_tokens,omitempty"`
	Stream         bool              `json:"stream"`
}

type deepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

type deepSeekSuggestionPayload struct {
	Suggestions []struct {
		Type        string `json:"type"`
		Title       string `json:"title"`
		Reason      string `json:"reason"`
		Content     string `json:"content"`
		ImpactScore int    `json:"impact_score"`
	} `json:"suggestions"`
}

func (s *MeituanAIService) generateDeepSeekSuggestions(storeID, accountID uint, dash *model.MeituanAIDashboard, orders []*model.MeituanAIOrder, reviews []*model.MeituanAIReview, now time.Time) ([]model.MeituanAISuggestion, error) {
	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("DEEPSEEK_API_KEY not configured")
	}
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("DEEPSEEK_BASE_URL")), "/")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	modelName := strings.TrimSpace(os.Getenv("DEEPSEEK_MODEL"))
	if modelName == "" {
		modelName = "deepseek-v4-flash"
	}

	prompt := buildDeepSeekMeituanPrompt(dash, orders, reviews)
	body, _ := json.Marshal(deepSeekRequest{
		Model: modelName,
		Messages: []deepSeekMessage{
			{Role: "system", Content: "你是精酿酒门店的美团外卖运营顾问。只输出严格 JSON，不要 Markdown。所有建议必须是半自动执行：给出建议和话术，由商家确认后执行。"},
			{Role: "user", Content: prompt},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
		Thinking:       map[string]string{"type": "disabled"},
		Temperature:    0.3,
		MaxTokens:      2200,
		Stream:         false,
	})
	req, err := http.NewRequest(http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result deepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if result.Error != nil && result.Error.Message != "" {
			return nil, fmt.Errorf("%s", result.Error.Message)
		}
		return nil, fmt.Errorf("deepseek request failed: %d", resp.StatusCode)
	}
	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return nil, fmt.Errorf("deepseek empty response")
	}

	var payload deepSeekSuggestionPayload
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &payload); err != nil {
		return nil, err
	}
	suggestions := make([]model.MeituanAISuggestion, 0, len(payload.Suggestions))
	for _, item := range payload.Suggestions {
		title := strings.TrimSpace(item.Title)
		content := strings.TrimSpace(item.Content)
		if title == "" || content == "" {
			continue
		}
		score := item.ImpactScore
		if score <= 0 {
			score = 60
		}
		if score > 100 {
			score = 100
		}
		suggestions = append(suggestions, makeSuggestion(storeID, accountID, ifEmptyString(item.Type, "ai"), title, strings.TrimSpace(item.Reason), content, score, now))
		if len(suggestions) >= 8 {
			break
		}
	}
	return suggestions, nil
}

func buildDeepSeekMeituanPrompt(dash *model.MeituanAIDashboard, orders []*model.MeituanAIOrder, reviews []*model.MeituanAIReview) string {
	type orderLite struct {
		OrderNo        string  `json:"order_no"`
		ProductSummary string  `json:"product_summary"`
		ActualAmount   float64 `json:"actual_amount"`
		PlatformFee    float64 `json:"platform_fee"`
		RefundAmount   float64 `json:"refund_amount"`
		Status         string  `json:"status"`
	}
	type reviewLite struct {
		Rating  int    `json:"rating"`
		Content string `json:"content"`
		Tags    string `json:"tags"`
	}
	os := make([]orderLite, 0, minInt(len(orders), 30))
	for i, order := range orders {
		if i >= 30 {
			break
		}
		os = append(os, orderLite{OrderNo: order.OrderNo, ProductSummary: order.ProductSummary, ActualAmount: order.ActualAmount, PlatformFee: order.PlatformFee, RefundAmount: order.RefundAmount, Status: order.Status})
	}
	rs := make([]reviewLite, 0, minInt(len(reviews), 30))
	for i, review := range reviews {
		if i >= 30 {
			break
		}
		rs = append(rs, reviewLite{Rating: review.Rating, Content: review.Content, Tags: review.Tags})
	}
	data := map[string]interface{}{
		"dashboard": dash,
		"orders":    os,
		"reviews":   rs,
	}
	b, _ := json.Marshal(data)
	return `请基于以下美团外卖经营数据，生成 3-8 条可人工确认执行的运营建议。
输出 JSON 格式必须为：
{"suggestions":[{"type":"product|review|reply|bundle|profit|activity|data|routine","title":"短标题","reason":"为什么建议这样做","content":"具体怎么做，包含可复制话术或执行步骤","impact_score":1-100}]}
要求：
1. 不要建议自动登录或绕过美团规则。
2. 评价回复必须礼貌、具体、可复制。
3. 活动建议必须考虑平台费用、退款和客单价。
4. 商品建议要围绕精酿酒/桶装/规格/套餐。
数据：` + string(b)
}

func ifEmptyString(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return strings.TrimSpace(v)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildMeituanAIOrders(items []model.ImportMeituanAIOrderItem) ([]model.MeituanAIOrder, error) {
	orders := make([]model.MeituanAIOrder, 0, len(items))
	for _, item := range items {
		orderNo := strings.TrimSpace(item.OrderNo)
		if orderNo == "" {
			continue
		}
		t, err := parseFlexibleTime(item.OrderTime)
		if err != nil {
			return nil, fmt.Errorf("订单%s时间格式错误", item.OrderNo)
		}
		raw := strings.TrimSpace(item.RawJSON)
		if raw == "" {
			if b, err := json.Marshal(item); err == nil {
				raw = string(b)
			}
		}
		orders = append(orders, model.MeituanAIOrder{
			OrderNo:        orderNo,
			OrderTime:      t,
			CustomerName:   strings.TrimSpace(item.CustomerName),
			ProductSummary: strings.TrimSpace(item.ProductSummary),
			OriginalAmount: item.OriginalAmount,
			ActualAmount:   item.ActualAmount,
			DeliveryFee:    item.DeliveryFee,
			PlatformFee:    item.PlatformFee,
			RefundAmount:   item.RefundAmount,
			Status:         strings.TrimSpace(item.Status),
			RawJSON:        raw,
		})
	}
	return orders, nil
}

func parseMeituanAIOrderFile(filename string, reader io.Reader) ([]model.ImportMeituanAIOrderItem, int, error) {
	data, err := io.ReadAll(io.LimitReader(reader, 10<<20))
	if err != nil {
		return nil, 0, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, 0, fmt.Errorf("文件为空")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".json" || strings.HasPrefix(strings.TrimSpace(string(data)), "[") {
		var items []model.ImportMeituanAIOrderItem
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, 0, fmt.Errorf("JSON格式错误: %w", err)
		}
		return items, 0, nil
	}
	return parseMeituanAIOrderCSV(data)
}

func parseMeituanAIOrderCSV(data []byte) ([]model.ImportMeituanAIOrderItem, int, error) {
	text := strings.TrimPrefix(string(data), "\ufeff")
	r := csv.NewReader(strings.NewReader(text))
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil || !meituanOrderCSVHeaderRecognized(records) {
		decoded, decodeErr := io.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder()))
		if decodeErr == nil {
			text = strings.TrimPrefix(string(decoded), "\ufeff")
			r = csv.NewReader(strings.NewReader(text))
			r.FieldsPerRecord = -1
			r.TrimLeadingSpace = true
			records, err = r.ReadAll()
		}
	}
	if err != nil {
		return nil, 0, fmt.Errorf("CSV格式错误: %w", err)
	}
	if len(records) < 2 {
		return nil, 0, fmt.Errorf("CSV至少需要表头和一行数据")
	}
	header := map[string]int{}
	for i, name := range records[0] {
		key := normalizeMeituanHeader(name)
		if key != "" {
			header[key] = i
		}
	}
	items := make([]model.ImportMeituanAIOrderItem, 0, len(records)-1)
	skipped := 0
	for _, row := range records[1:] {
		item := model.ImportMeituanAIOrderItem{
			OrderNo:        csvCell(row, header, "order_no"),
			OrderTime:      csvCell(row, header, "order_time"),
			CustomerName:   csvCell(row, header, "customer_name"),
			ProductSummary: csvCell(row, header, "product_summary"),
			OriginalAmount: parseAmount(csvCell(row, header, "original_amount")),
			ActualAmount:   parseAmount(csvCell(row, header, "actual_amount")),
			DeliveryFee:    parseAmount(csvCell(row, header, "delivery_fee")),
			PlatformFee:    parseAmount(csvCell(row, header, "platform_fee")),
			RefundAmount:   parseAmount(csvCell(row, header, "refund_amount")),
			Status:         csvCell(row, header, "status"),
		}
		if item.OrderNo == "" || item.OrderTime == "" {
			skipped++
			continue
		}
		if item.OriginalAmount == 0 {
			item.OriginalAmount = item.ActualAmount + item.RefundAmount
		}
		if b, err := json.Marshal(row); err == nil {
			item.RawJSON = string(b)
		}
		items = append(items, item)
	}
	return items, skipped, nil
}

func meituanOrderCSVHeaderRecognized(records [][]string) bool {
	if len(records) == 0 {
		return false
	}
	for _, name := range records[0] {
		if normalizeMeituanHeader(name) == "order_no" {
			return true
		}
	}
	return false
}

func normalizeMeituanHeader(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	switch s {
	case "orderno", "orderid", "订单号", "订单编号", "美团订单号":
		return "order_no"
	case "ordertime", "orderdate", "createtime", "下单时间", "订单时间", "创建时间", "支付时间":
		return "order_time"
	case "customername", "username", "buyername", "顾客", "顾客昵称", "客户", "客户名称":
		return "customer_name"
	case "productsummary", "products", "goods", "商品", "商品信息", "商品名称", "订单商品", "菜品":
		return "product_summary"
	case "originalamount", "totalamount", "原价", "订单原价", "订单金额", "商品总价":
		return "original_amount"
	case "actualamount", "income", "paidamount", "实收", "实收金额", "商家实收", "预计收入", "到账金额":
		return "actual_amount"
	case "deliveryfee", "配送费", "配送金额":
		return "delivery_fee"
	case "platformfee", "servicefee", "commission", "平台费", "平台服务费", "佣金":
		return "platform_fee"
	case "refundamount", "refund", "退款", "退款金额":
		return "refund_amount"
	case "status", "订单状态", "状态":
		return "status"
	default:
		return ""
	}
}

func csvCell(row []string, header map[string]int, key string) string {
	idx, ok := header[key]
	if !ok || idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func parseAmount(v string) float64 {
	s := strings.TrimSpace(v)
	s = strings.ReplaceAll(s, "¥", "")
	s = strings.ReplaceAll(s, "￥", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "元", "")
	if s == "" {
		return 0
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return n
}

func normalizeMeituanAccountConfig(account *model.MeituanAIOperatorAccount) {
	if account.BusinessID <= 0 {
		account.BusinessID = 2
	}
	if strings.TrimSpace(account.APIVersion) == "" {
		account.APIVersion = "2"
	}
	if strings.TrimSpace(account.APIBaseURL) == "" {
		account.APIBaseURL = "https://api-open-cater.meituan.com"
	}
	if strings.TrimSpace(account.QueryOrderPath) == "" {
		account.QueryOrderPath = "/api/order/queryById"
	}
}

func (s *MeituanAIService) fetchMeituanOrderByID(account *model.MeituanAIOperatorAccount, orderID string) (model.MeituanAIOrder, error) {
	if strings.TrimSpace(account.DeveloperID) == "" || strings.TrimSpace(account.SignKey) == "" || strings.TrimSpace(account.AppAuthToken) == "" {
		return model.MeituanAIOrder{}, fmt.Errorf("请先配置美团开放平台DeveloperId、SignKey和appAuthToken")
	}
	biz := map[string]interface{}{"orderId": orderID}
	raw, err := s.callMeituanOpenAPI(account, account.QueryOrderPath, biz)
	if err != nil {
		return model.MeituanAIOrder{}, err
	}
	row, err := parseMeituanOpenAPIOrder(orderID, raw)
	if err != nil {
		return model.MeituanAIOrder{}, err
	}
	row.RawJSON = string(raw)
	return row, nil
}

func (s *MeituanAIService) callMeituanOpenAPI(account *model.MeituanAIOperatorAccount, path string, biz map[string]interface{}) ([]byte, error) {
	bizBytes, _ := json.Marshal(biz)
	values := map[string]string{
		"developerId":  strings.TrimSpace(account.DeveloperID),
		"businessId":   strconv.Itoa(account.BusinessID),
		"appAuthToken": strings.TrimSpace(account.AppAuthToken),
		"charset":      "UTF-8",
		"timestamp":    strconv.FormatInt(time.Now().UnixMilli(), 10),
		"version":      strings.TrimSpace(account.APIVersion),
		"biz":          string(bizBytes),
	}
	values["sign"] = signMeituanOpenAPI(values, account.SignKey)

	form := url.Values{}
	for k, v := range values {
		form.Set(k, v)
	}
	endpoint := strings.TrimRight(account.APIBaseURL, "/") + "/" + strings.TrimLeft(path, "/")
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("美团接口请求失败: %d %s", resp.StatusCode, truncateText(string(raw), 200))
	}
	if ok, msg := meituanOpenAPISuccess(raw); !ok {
		return nil, fmt.Errorf("美团接口返回失败: %s", msg)
	}
	return raw, nil
}

func signMeituanOpenAPI(values map[string]string, signKey string) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString(values[k])
	}
	b.WriteString(strings.TrimSpace(signKey))
	sum := sha1.Sum([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func meituanOpenAPISuccess(raw []byte) (bool, string) {
	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return false, "返回不是JSON: " + err.Error()
	}
	if code, ok := body["code"]; ok {
		switch v := code.(type) {
		case float64:
			if v == 0 || v == 200 {
				return true, ""
			}
		case string:
			if v == "0" || v == "200" || strings.EqualFold(v, "success") {
				return true, ""
			}
		}
	}
	if success, ok := body["success"].(bool); ok && success {
		return true, ""
	}
	msg := firstString(body, "msg", "message", "errorMsg", "error", "error_message")
	if msg == "" {
		msg = truncateText(string(raw), 200)
	}
	return false, msg
}

func parseMeituanOpenAPIOrder(orderID string, raw []byte) (model.MeituanAIOrder, error) {
	var body interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return model.MeituanAIOrder{}, err
	}
	data := pickMeituanData(body)
	m, ok := data.(map[string]interface{})
	if !ok {
		return model.MeituanAIOrder{}, fmt.Errorf("美团订单返回结构无法识别")
	}
	orderNo := firstString(m, "orderId", "order_id", "orderNo", "order_no", "wmOrderIdView", "wm_order_id_view")
	if orderNo == "" {
		orderNo = orderID
	}
	orderTime := firstString(m, "ctime", "createTime", "create_time", "orderTime", "order_time", "payTime", "pay_time")
	t, err := parseFlexibleTime(orderTime)
	if err != nil {
		if n := firstFloat(m, "ctime", "createTime", "orderTime", "payTime"); n > 0 {
			t = unixMaybeMilli(n)
		} else {
			t = time.Now()
		}
	}
	return model.MeituanAIOrder{
		OrderNo:        orderNo,
		OrderTime:      t,
		CustomerName:   firstString(m, "recipientName", "recipient_name", "userName", "user_name", "customerName", "customer_name"),
		ProductSummary: buildMeituanProductSummary(m),
		OriginalAmount: firstFloat(m, "total", "totalPrice", "originalPrice", "originalAmount", "boxTotalPrice"),
		ActualAmount:   firstFloat(m, "actualPrice", "actualAmount", "wmPoiReceiveCent", "poiReceive", "income"),
		DeliveryFee:    firstFloat(m, "shippingFee", "deliveryFee", "logisticsFee"),
		PlatformFee:    firstFloat(m, "platformFee", "commission", "serviceFee"),
		RefundAmount:   firstFloat(m, "refundAmount", "refundPrice"),
		Status:         firstString(m, "status", "statusDesc", "orderStatus", "order_status"),
	}, nil
}

func pickMeituanData(v interface{}) interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return v
	}
	for _, key := range []string{"data", "result", "body"} {
		if next, exists := m[key]; exists {
			return pickMeituanData(next)
		}
	}
	return m
}

func buildMeituanProductSummary(m map[string]interface{}) string {
	if s := firstString(m, "productSummary", "product_summary", "detail", "goods", "foodNames"); s != "" {
		return s
	}
	for _, key := range []string{"detail", "details", "cartDetailVos", "foodList", "foodlist"} {
		if arr, ok := m[key].([]interface{}); ok {
			names := make([]string, 0, len(arr))
			for _, item := range arr {
				if im, ok := item.(map[string]interface{}); ok {
					name := firstString(im, "name", "foodName", "food_name", "skuName", "sku_name")
					count := firstFloat(im, "count", "quantity", "num", "boxNum")
					if name != "" {
						if count > 0 {
							name = fmt.Sprintf("%s x%.0f", name, count)
						}
						names = append(names, name)
					}
				}
			}
			if len(names) > 0 {
				return strings.Join(names, "、")
			}
		}
	}
	return ""
}

func firstString(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			switch t := v.(type) {
			case string:
				if strings.TrimSpace(t) != "" {
					return strings.TrimSpace(t)
				}
			case float64:
				if t > 0 {
					return strconv.FormatInt(int64(t), 10)
				}
			case json.Number:
				return t.String()
			}
		}
	}
	return ""
}

func firstFloat(m map[string]interface{}, keys ...string) float64 {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			switch t := v.(type) {
			case float64:
				return normalizeMoneyAmount(t)
			case int:
				return normalizeMoneyAmount(float64(t))
			case json.Number:
				n, _ := t.Float64()
				return normalizeMoneyAmount(n)
			case string:
				return parseAmount(t)
			}
		}
	}
	return 0
}

func normalizeMoneyAmount(v float64) float64 {
	if v > 10000 {
		return v / 100
	}
	return v
}

func unixMaybeMilli(v float64) time.Time {
	if v > 100000000000 {
		return time.UnixMilli(int64(v))
	}
	return time.Unix(int64(v), 0)
}

func cleanMeituanOrderIDs(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		for _, part := range strings.FieldsFunc(item, func(r rune) bool {
			return r == ',' || r == '，' || r == '\n' || r == '\r' || r == '\t' || r == ' '
		}) {
			if s := strings.TrimSpace(part); s != "" {
				out = append(out, s)
			}
		}
	}
	return out
}

func dedupeStrings(items []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func makeSuggestion(storeID, accountID uint, typ, title, reason, content string, score int, t time.Time) model.MeituanAISuggestion {
	return model.MeituanAISuggestion{
		StoreID:     storeID,
		AccountID:   accountID,
		Type:        typ,
		Title:       title,
		Reason:      reason,
		Content:     content,
		ImpactScore: score,
		Status:      model.MeituanSuggestionStatusPending,
		GeneratedAt: t,
	}
}

func parseFlexibleTime(v string) (time.Time, error) {
	s := strings.TrimSpace(v)
	layouts := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"}
	var last error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			return t, nil
		}
		last = err
	}
	return time.Time{}, last
}

func classifyReview(rating int, content string) (string, []string) {
	text := strings.ToLower(content)
	tags := make([]string, 0)
	rules := map[string][]string{
		"包装": {"漏", "洒", "破", "包装", "袋子"},
		"配送": {"慢", "骑手", "配送", "超时", "凉"},
		"口味": {"难喝", "味道", "酸", "淡", "苦"},
		"缺货": {"没货", "缺货", "少", "漏发"},
		"价格": {"贵", "价格", "不值"},
	}
	for tag, words := range rules {
		for _, word := range words {
			if strings.Contains(text, word) {
				tags = append(tags, tag)
				break
			}
		}
	}
	if len(tags) == 0 {
		tags = append(tags, "体验")
	}
	if rating <= 3 {
		return "negative", tags
	}
	if rating == 4 {
		return "neutral", tags
	}
	return "positive", tags
}

func buildReviewReply(rating int, content string, tags []string) string {
	if rating <= 3 {
		return fmt.Sprintf("非常抱歉这次体验没有达到预期，我们已经记录您反馈的%s问题，会马上复盘并优化。也欢迎您再次联系门店，我们会尽力给您一个更好的处理方案。", strings.Join(cleanTags(tags), "、"))
	}
	return "感谢您的认可和反馈，我们会继续保持出品和配送体验，期待再次为您服务。"
}

func cleanTags(tags []string) []string {
	out := make([]string, 0)
	seen := map[string]bool{}
	for _, tag := range tags {
		t := strings.TrimSpace(tag)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

func topProductName(orders []*model.MeituanAIOrder) string {
	counts := map[string]int{}
	re := regexp.MustCompile(`[、,，/]+`)
	for _, order := range orders {
		parts := re.Split(order.ProductSummary, -1)
		for _, p := range parts {
			name := strings.TrimSpace(p)
			if name == "" {
				continue
			}
			counts[name]++
		}
	}
	type pair struct {
		Name  string
		Count int
	}
	list := make([]pair, 0, len(counts))
	for name, count := range counts {
		list = append(list, pair{Name: name, Count: count})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Count > list[j].Count })
	if len(list) == 0 {
		return ""
	}
	return list[0].Name
}

func truncateText(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n]) + "..."
}
