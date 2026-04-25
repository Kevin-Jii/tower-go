package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

type StoreAccountService struct {
	storeAccountModule    *module.StoreAccountModule
	inventoryModule       *module.InventoryModule
	productModule         *module.SupplierProductModule
	storeModule           *module.StoreModule
	userModule            *module.UserModule
	dictModule            *module.DictModule
	dingTalkService       *DingTalkService
	botModule             *module.DingTalkBotModule
	templateService       *MessageTemplateService
	imageGeneratorService *ImageGeneratorService
}

func NewStoreAccountService(
	storeAccountModule *module.StoreAccountModule,
	inventoryModule *module.InventoryModule,
	productModule *module.SupplierProductModule,
	storeModule *module.StoreModule,
	userModule *module.UserModule,
	dictModule *module.DictModule,
	dingTalkService *DingTalkService,
	botModule *module.DingTalkBotModule,
	templateService *MessageTemplateService,
	imageGeneratorService *ImageGeneratorService,
) *StoreAccountService {
	return &StoreAccountService{
		storeAccountModule:    storeAccountModule,
		inventoryModule:       inventoryModule,
		productModule:         productModule,
		storeModule:           storeModule,
		userModule:            userModule,
		dictModule:            dictModule,
		dingTalkService:       dingTalkService,
		botModule:             botModule,
		templateService:       templateService,
		imageGeneratorService: imageGeneratorService,
	}
}

// Create 创建记账
func (s *StoreAccountService) Create(storeID, operatorID uint, req *model.CreateStoreAccountReq) (*model.StoreAccount, error) {
	accountNo := s.storeAccountModule.GenerateAccountNo()

	// 解析记账日期
	accountDate := time.Now()
	if req.AccountDate != "" {
		if t, err := time.Parse("2006-01-02", req.AccountDate); err == nil {
			accountDate = t
		}
	}

	// 构建明细
	var items []model.StoreAccountItem
	var totalAmount float64

	for _, item := range req.Items {
		// 获取商品名称
		productName := ""
		unit := item.Unit
		if s.productModule != nil {
			if product, err := s.productModule.GetByID(item.ProductID); err == nil && product != nil {
				productName = product.Name
				if unit == "" {
					unit = product.Unit
				}
			}
		}

		// 计算金额
		amount := item.Amount
		if amount == 0 && item.Price > 0 && item.Quantity > 0 {
			amount = item.Price * item.Quantity
		}

		items = append(items, model.StoreAccountItem{
			ProductID:   item.ProductID,
			ProductName: productName,
			Spec:        item.Spec,
			Quantity:    item.Quantity,
			Unit:        unit,
			Price:       item.Price,
			Amount:      amount,
			Remark:      item.Remark,
		})

		totalAmount += amount
	}

	account := &model.StoreAccount{
		AccountNo:          accountNo,
		StoreID:            storeID,
		Channel:            req.Channel,
		OrderNo:            req.OrderNo,
		TotalAmount:        totalAmount,
		OtherExpenseAmount: req.OtherExpenseAmount,
		NetIncomeAmount:    totalAmount - req.OtherExpenseAmount,
		ItemCount:          len(items),
		TagCode:            req.TagCode,
		TagName:            req.TagName,
		Remark:             req.Remark,
		OperatorID:         operatorID,
		AccountDate:        accountDate,
		Items:              items,
	}

	inventoryOutOrder := &model.InventoryOrder{
		OrderNo:       s.inventoryModule.GenerateOrderNo(model.InventoryTypeOut),
		Type:          model.InventoryTypeOut,
		StoreID:       storeID,
		Reason:        model.ReasonSale,
		Remark:        fmt.Sprintf("记账自动出库，记账单号:%s", accountNo),
		TotalQuantity: 0,
		ItemCount:     len(items),
		OperatorID:    operatorID,
	}
	for _, item := range items {
		inventoryOutOrder.TotalQuantity += item.Quantity
		inventoryOutOrder.Items = append(inventoryOutOrder.Items, model.InventoryOrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			Remark:      "记账自动出库",
		})
	}
	if store, err := s.storeModule.GetByID(storeID); err == nil && store != nil {
		inventoryOutOrder.StoreName = store.Name
	}
	if user, err := s.userModule.GetByID(operatorID); err == nil && user != nil {
		inventoryOutOrder.OperatorName = user.Nickname
		if inventoryOutOrder.OperatorName == "" {
			inventoryOutOrder.OperatorName = user.Username
		}
		inventoryOutOrder.OperatorPhone = user.Phone
	}

	if err := s.storeAccountModule.CreateWithInventoryOut(account, inventoryOutOrder); err != nil {
		return nil, err
	}

	// 获取操作人名称
	operatorName := ""
	if s.userModule != nil {
		if user, err := s.userModule.GetByID(operatorID); err == nil && user != nil {
			operatorName = user.Nickname
			if operatorName == "" {
				operatorName = user.Username
			}
		}
	}

	// 获取渠道名称（字典转换）
	channelName := account.Channel
	if s.dictModule != nil && account.Channel != "" {
		if dictData, err := s.dictModule.GetDataByTypeAndValue("sales_channel", account.Channel); err == nil && dictData != nil {
			channelName = dictData.Label
		}
	}

	// 异步发送钉钉通知
	go s.sendDingTalkNotification(account, storeID, operatorName, channelName)

	return account, nil
}

// sendDingTalkNotification 发送记账通知到门店负责人
func (s *StoreAccountService) sendDingTalkNotification(account *model.StoreAccount, storeID uint, operatorName, channelName string) {
	if s.dingTalkService == nil || s.storeModule == nil || s.botModule == nil {
		return
	}

	// 获取门店信息
	store, err := s.storeModule.GetByID(storeID)
	if err != nil || store == nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to get store for notification", "storeID", storeID, "error", err)
		}
		return
	}

	// 检查门店是否有联系电话
	if store.Phone == "" {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Store has no phone, skip notification", "storeID", storeID)
		}
		return
	}

	// 获取门店绑定的机器人
	bot, err := s.botModule.GetByStoreID(storeID)
	if err != nil || bot == nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("No bot found for store", "storeID", storeID, "error", err)
		}
		return
	}

	if !bot.IsEnabled || bot.BotType != "stream" {
		return
	}

	// 操作人显示
	operatorDisplay := operatorName
	if operatorDisplay == "" {
		operatorDisplay = "未知"
	}

	// 尝试生成通知图片
	var imageURL string
	if s.imageGeneratorService != nil {
		var items []AccountItemData
		for _, item := range account.Items {
			items = append(items, AccountItemData{
				Name:     item.ProductName,
				Quantity: item.Quantity,
				Unit:     item.Unit,
				Amount:   item.Amount,
			})
		}

		imgData := &AccountNotifyData{
			StoreName:    store.Name,
			AccountNo:    account.AccountNo,
			ChannelName:  channelName,
			AccountDate:  account.AccountDate.Format("2006-01-02"),
			OperatorName: operatorDisplay,
			Items:        items,
			TotalAmount:  account.TotalAmount,
			CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		}

		url, err := s.imageGeneratorService.GenerateAccountNotifyImage(imgData)
		if err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Warnw("Failed to generate notify image", "error", err)
			}
		} else {
			imageURL = url
		}
	}

	// 构建商品明细
	var itemLines []string
	for i, item := range account.Items {
		line := fmt.Sprintf("%d. %s x%.2f%s = ¥%.2f", i+1, item.ProductName, item.Quantity, item.Unit, item.Amount)
		itemLines = append(itemLines, line)
	}

	var title, text string

	// 构建文字消息内容
	// 尝试使用模板
	if s.templateService != nil {
		data := map[string]interface{}{
			"StoreName":    store.Name,
			"AccountNo":    account.AccountNo,
			"ChannelName":  channelName,
			"AccountDate":  account.AccountDate.Format("2006-01-02"),
			"OperatorName": operatorDisplay,
			"ItemList":     strings.Join(itemLines, "\n\n"),
			"TotalAmount":  fmt.Sprintf("%.2f", account.TotalAmount),
			"ItemCount":    account.ItemCount,
			"CreateTime":   time.Now().Format("2006-01-02 15:04:05"),
		}
		var err error
		title, text, err = s.templateService.RenderTemplate(model.TemplateStoreAccountCreated, data)
		if err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Warnw("Failed to render template, using default", "error", err)
			}
		}
	}

	// 如果模板渲染失败，使用默认格式
	if text == "" {
		title = fmt.Sprintf("📝 新记账通知 - %s", store.Name)
		text = fmt.Sprintf("## %s\n\n"+
			"**记账编号：** %s\n\n"+
			"**渠道来源：** %s\n\n"+
			"**记账日期：** %s\n\n"+
			"**操作人：** %s\n\n"+
			"### 商品明细\n\n"+
			"%s\n\n"+
			"**合计金额：** ¥%.2f\n\n"+
			"**商品数量：** %d 项\n\n"+
			"---\n\n"+
			"%s",
			title,
			account.AccountNo,
			channelName,
			account.AccountDate.Format("2006-01-02"),
			operatorDisplay,
			strings.Join(itemLines, "\n\n"),
			account.TotalAmount,
			account.ItemCount,
			time.Now().Format("2006-01-02 15:04:05"),
		)
	}

	// 发送通知：如果有图片，先发图片再发文字；否则只发文字
	var sendErr error
	if imageURL != "" {
		sendErr = s.dingTalkService.SendStreamMarkdownWithImageToMobile(bot, title, text, imageURL, store.Phone)
	} else {
		sendErr = s.dingTalkService.SendStreamMarkdownToMobile(bot, title, text, store.Phone)
	}

	if sendErr != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to send account notification",
				"storeID", storeID,
				"accountNo", account.AccountNo,
				"error", sendErr,
			)
		}
	} else {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("Account notification sent",
				"storeID", storeID,
				"accountNo", account.AccountNo,
				"mobile", store.Phone,
				"hasImage", imageURL != "",
			)
		}
	}
}

// Get 获取记账详情
func (s *StoreAccountService) Get(id uint) (*model.StoreAccount, error) {
	return s.storeAccountModule.GetByID(id)
}

// List 记账列表
func (s *StoreAccountService) List(req *model.ListStoreAccountReq) ([]*model.StoreAccount, int64, error) {
	return s.storeAccountModule.List(req)
}

// Update 更新记账
func (s *StoreAccountService) Update(id uint, req *model.UpdateStoreAccountReq) error {
	updates := make(map[string]interface{})

	if req.Channel != "" {
		updates["channel"] = req.Channel
	}
	if req.OrderNo != "" {
		updates["order_no"] = req.OrderNo
	}
	if req.TagCode != "" {
		updates["tag_code"] = req.TagCode
	}
	if req.TagName != "" {
		updates["tag_name"] = req.TagName
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.AccountDate != "" {
		if t, err := time.Parse("2006-01-02", req.AccountDate); err == nil {
			updates["account_date"] = t
		}
	}
	if req.OtherExpenseAmount != nil {
		updates["other_expense_amount"] = *req.OtherExpenseAmount
		if account, err := s.storeAccountModule.GetByID(id); err == nil && account != nil {
			updates["net_income_amount"] = account.TotalAmount - *req.OtherExpenseAmount
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return s.storeAccountModule.Update(id, updates)
}

// Delete 删除记账
func (s *StoreAccountService) Delete(id uint) error {
	return s.storeAccountModule.Delete(id)
}

// GetStats 获取统计
func (s *StoreAccountService) GetStats(storeID uint, startDate, endDate string) (map[string]interface{}, error) {
	totalAmount, count, err := s.storeAccountModule.GetStatsByDateRange(storeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_amount": totalAmount,
		"count":        count,
	}, nil
}
