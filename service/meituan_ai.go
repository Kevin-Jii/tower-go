package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type MeituanAIService struct {
	module *module.MeituanAIModule
}

func NewMeituanAIService(m *module.MeituanAIModule) *MeituanAIService {
	return &MeituanAIService{module: m}
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
		StoreID:    realStoreID,
		ShopName:   strings.TrimSpace(req.ShopName),
		ShopID:     strings.TrimSpace(req.ShopID),
		LoginName:  strings.TrimSpace(req.LoginName),
		AuthStatus: "manual",
		IsEnabled:  true,
		Remark:     strings.TrimSpace(req.Remark),
	}
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
	orders := make([]model.MeituanAIOrder, 0, len(req.Orders))
	for _, item := range req.Orders {
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
			OrderNo:        strings.TrimSpace(item.OrderNo),
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
	n, err := s.module.UpsertOrders(account, orders)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"imported": n}, nil
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
	suggestions := make([]model.MeituanAISuggestion, 0)
	if dash.OrderCount == 0 {
		suggestions = append(suggestions, makeSuggestion(req.StoreID, req.AccountID, "data", "先导入美团订单数据", "当前周期没有订单数据，AI无法判断爆品、客单价和活动效果。", "从美团商家后台导出订单后，在本模块导入。", 90, now))
	}
	if dash.NegativeRate >= 15 {
		suggestions = append(suggestions, makeSuggestion(req.StoreID, req.AccountID, "review", "差评率偏高，优先处理差评原因", fmt.Sprintf("当前差评率 %.1f%%，已经影响店铺转化。", dash.NegativeRate), "优先回复低分评价，并把高频问题拆成包装、配送、口味、缺货四类处理。", 88, now))
	}
	if dash.AvgOrderAmount > 0 && dash.AvgOrderAmount < 60 {
		suggestions = append(suggestions, makeSuggestion(req.StoreID, req.AccountID, "bundle", "设计高客单套餐", fmt.Sprintf("当前客单价 %.2f，适合通过组合套餐提升单均收入。", dash.AvgOrderAmount), "选择销量最高商品，搭配杯具、小食或第二规格，做成高毛利套餐。", 78, now))
	}
	if dash.PlatformFee > 0 && dash.SalesAmount > 0 && dash.PlatformFee/dash.SalesAmount >= 0.08 {
		suggestions = append(suggestions, makeSuggestion(req.StoreID, req.AccountID, "profit", "复核平台费用和活动力度", fmt.Sprintf("平台费用占销售额 %.1f%%。", dash.PlatformFee/dash.SalesAmount*100), "检查满减、配送费补贴和平台服务费，避免活动后利润被吃掉。", 82, now))
	}
	if product := topProductName(orders); product != "" {
		suggestions = append(suggestions, makeSuggestion(req.StoreID, req.AccountID, "product", "强化爆品展示", "当前周期已有明显高频商品。", fmt.Sprintf("把「%s」放到美团店铺靠前位置，标题加容量/场景词，并搭配套餐入口。", product), 75, now))
	}
	for _, r := range reviews {
		if r.Rating <= 3 {
			title := "生成差评回复草稿"
			content := strings.TrimSpace(r.SuggestedReply)
			if content == "" {
				content = buildReviewReply(r.Rating, r.Content, strings.Split(r.Tags, ","))
			}
			suggestions = append(suggestions, makeSuggestion(req.StoreID, req.AccountID, "reply", title, fmt.Sprintf("订单%s存在%d星评价：%s", r.OrderNo, r.Rating, truncateText(r.Content, 80)), content, 72, now))
			if len(suggestions) >= 8 {
				break
			}
		}
	}
	if len(suggestions) == 0 {
		suggestions = append(suggestions, makeSuggestion(req.StoreID, req.AccountID, "routine", "保持当前运营节奏", "当前周期未发现明显风险。", "继续保持商品库存、评价回复和活动利润监控；建议每日导入订单评价刷新一次。", 55, now))
	}

	if err := s.module.ClearPendingSuggestions(req.StoreID, req.AccountID); err != nil {
		return nil, err
	}
	if err := s.module.CreateSuggestions(suggestions); err != nil {
		return nil, err
	}
	return map[string]interface{}{"generated": len(suggestions)}, nil
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
