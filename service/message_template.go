package service

import (
	"bytes"
	"text/template"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type MessageTemplateService struct {
	templateModule *module.MessageTemplateModule
}

func NewMessageTemplateService(templateModule *module.MessageTemplateModule) *MessageTemplateService {
	return &MessageTemplateService{
		templateModule: templateModule,
	}
}

// RenderTemplate 渲染模板
func (s *MessageTemplateService) RenderTemplate(code string, data map[string]interface{}) (title, content string, err error) {
	tpl, err := s.templateModule.GetByCode(code)
	if err != nil {
		return "", "", err
	}

	// 渲染标题
	if tpl.Title != "" {
		title, err = s.render(tpl.Title, data)
		if err != nil {
			return "", "", err
		}
	}

	// 渲染内容
	content, err = s.render(tpl.Content, data)
	if err != nil {
		return "", "", err
	}

	return title, content, nil
}

// render 使用 Go template 渲染
func (s *MessageTemplateService) render(tplStr string, data map[string]interface{}) (string, error) {
	tpl, err := template.New("msg").Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetByCode 获取模板
func (s *MessageTemplateService) GetByCode(code string) (*model.MessageTemplate, error) {
	return s.templateModule.GetByCode(code)
}

// List 获取模板列表
func (s *MessageTemplateService) List() ([]*model.MessageTemplate, error) {
	return s.templateModule.List()
}

// GetByID 获取模板详情
func (s *MessageTemplateService) GetByID(id uint) (*model.MessageTemplate, error) {
	return s.templateModule.GetByID(id)
}

// Create 创建模板
func (s *MessageTemplateService) Create(req *model.CreateMessageTemplateReq) (*model.MessageTemplate, error) {
	tpl := &model.MessageTemplate{
		Code:        req.Code,
		Name:        req.Name,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Variables:   req.Variables,
		IsEnabled:   true,
	}
	if req.IsEnabled != nil {
		tpl.IsEnabled = *req.IsEnabled
	}
	if err := s.templateModule.Create(tpl); err != nil {
		return nil, err
	}
	return tpl, nil
}

// Update 更新模板
func (s *MessageTemplateService) Update(id uint, req *model.UpdateMessageTemplateReq) error {
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Variables != nil {
		updates["variables"] = *req.Variables
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if len(updates) == 0 {
		return nil
	}
	return s.templateModule.Update(id, updates)
}

// Delete 删除模板
func (s *MessageTemplateService) Delete(id uint) error {
	return s.templateModule.Delete(id)
}

// InitDefaultTemplates 初始化默认模板
func (s *MessageTemplateService) InitDefaultTemplates() error {
	// 模板已迁移到 init_seed_data.sql，这里只做兼容处理
	// 如果数据库中没有模板，才初始化
	list, err := s.templateModule.List()
	if err != nil {
		// 表不存在通常发生在禁用自动迁移且未执行结构SQL的情况下；
		// 为避免启动阶段报错，这里直接跳过（由迁移/seed补齐表结构与数据）。
		if strings.Contains(err.Error(), "doesn't exist") || strings.Contains(err.Error(), "Error 1146") {
			return nil
		}
		return err
	}
	if len(list) > 0 {
		return nil // 已有模板，跳过
	}

	// 初始化基础模板（完整模板请执行 init_seed_data.sql）
	templates := []model.MessageTemplate{
		{Code: model.TemplateStoreAccountCreated, Name: "记账通知", Title: "📝 新记账通知 - {{.StoreName}}", Content: "记账编号: {{.AccountNo}}\n渠道: {{.ChannelName}}\n操作人: {{.OperatorName}}\n\n{{.ItemList}}\n\n合计: ¥{{.TotalAmount}}", IsEnabled: true},
		{Code: model.TemplateInventoryCreated, Name: "入库通知", Title: "📦 新入库通知 - {{.StoreName}}", Content: "入库单号: {{.OrderNo}}\n类型: {{.OrderType}}\n操作人: {{.OperatorName}}\n\n{{.ItemList}}\n\n总量: {{.TotalAmount}}", IsEnabled: true},
		{Code: model.TemplateBotHelp, Name: "机器人帮助", Title: "📋 功能菜单", Content: "## 📋 功能菜单\n\n**库存相关**\n- 库存查询\n- 查询库存 商品名\n\n**记账相关**\n- 今日记账\n\n**入库相关**\n- 今日入库\n\n---\n发送 **帮助** 查看此菜单", IsEnabled: true},
		{Code: model.TemplateBotInventoryQuery, Name: "库存查询回复", Title: "📦 库存查询", Content: "## 📦 库存查询\n\n**共{{.Total}}项**\n\n{{.ItemList}}\n\n---\n{{.CreateTime}}", IsEnabled: true},
		{Code: model.TemplateBotTodayAccount, Name: "今日记账回复", Title: "📝 今日记账", Content: "## 📝 今日记账汇总\n\n**日期：** {{.Date}}\n\n**笔数：** {{.Count}} 笔\n\n**总额：** ¥{{.TotalAmount}}\n\n---\n{{.CreateTime}}", IsEnabled: true},
		{Code: model.TemplateBotTodayInventory, Name: "今日入库回复", Title: "📦 今日入库", Content: "## 📦 今日入库汇总\n\n**日期：** {{.Date}}\n\n**单数：** {{.Count}} 单\n\n**总量：** {{.TotalQuantity}}\n\n{{.ItemList}}\n\n---\n{{.CreateTime}}", IsEnabled: true},
		{Code: model.TemplateBotSearchResult, Name: "搜索结果回复", Title: "🔍 库存搜索", Content: "## 🔍 库存搜索\n\n**关键词：** {{.Keyword}}\n\n**共{{.Total}}项**\n\n{{.ItemList}}\n\n---\n{{.CreateTime}}", IsEnabled: true},
		{Code: model.TemplateBotUnknown, Name: "未知命令回复", Title: "🤖 智能助手", Content: "## 🤖 智能助手\n\n您发送: {{.Content}}\n\n抱歉，暂时无法理解。\n\n发送 **帮助** 查看可用功能", IsEnabled: true},
	}

	for _, tpl := range templates {
		if err := s.templateModule.Upsert(&tpl); err != nil {
			return err
		}
	}
	return nil
}
