package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

func isLargePackUnit(unit string) bool {
	u := strings.ToLower(strings.TrimSpace(unit))
	if u == "" {
		return false
	}
	return strings.Contains(u, "箱") ||
		strings.Contains(u, "桶") ||
		strings.Contains(u, "case") ||
		strings.Contains(u, "barrel")
}

func resolveUnitPrice(unit string, product *model.SupplierProduct) float64 {
	if product == nil {
		return 0
	}
	if isLargePackUnit(unit) {
		if product.CasePrice > 0 {
			return product.CasePrice
		}
		if product.BottlePrice > 0 {
			return product.BottlePrice
		}
		if product.Price > 0 {
			return product.Price
		}
		return 0
	}

	if product.BottlePrice > 0 {
		return product.BottlePrice
	}
	if product.Price > 0 {
		return product.Price
	}
	return product.CasePrice
}

func resolveUnitPriceFromSpecs(unit string, specs []*model.ProductUnitSpec) float64 {
	if len(specs) == 0 {
		return 0
	}
	normalized := strings.ToLower(strings.TrimSpace(unit))

	// 1) 精确匹配（unit_code / unit_name）
	for _, spec := range specs {
		if spec == nil || !spec.IsEnabled || spec.SalePrice <= 0 {
			continue
		}
		if normalized == strings.ToLower(strings.TrimSpace(spec.UnitCode)) ||
			normalized == strings.ToLower(strings.TrimSpace(spec.UnitName)) {
			return spec.SalePrice
		}
	}

	// 2) 模糊包含匹配（兼容“L/瓶”“箱/桶”这类展示名）
	if normalized != "" {
		for _, spec := range specs {
			if spec == nil || !spec.IsEnabled || spec.SalePrice <= 0 {
				continue
			}
			code := strings.ToLower(strings.TrimSpace(spec.UnitCode))
			name := strings.ToLower(strings.TrimSpace(spec.UnitName))
			if strings.Contains(code, normalized) || strings.Contains(name, normalized) ||
				strings.Contains(normalized, code) || strings.Contains(normalized, name) {
				return spec.SalePrice
			}
		}
	}

	// 3) 兜底：按小/大规格选择
	needLarge := isLargePackUnit(unit)
	for _, spec := range specs {
		if spec == nil || !spec.IsEnabled || spec.SalePrice <= 0 {
			continue
		}
		if needLarge && spec.FactorToBase > 1 {
			return spec.SalePrice
		}
		if !needLarge && spec.FactorToBase <= 1 {
			return spec.SalePrice
		}
	}
	return 0
}

func tryResolveUnitSpecSalePrice(unit string, specs []*model.ProductUnitSpec) (float64, bool) {
	price := resolveUnitPriceFromSpecs(unit, specs)
	if price > 0 {
		return price, true
	}
	return 0, false
}

type StoreAccountService struct {
	storeAccountModule    *module.StoreAccountModule
	inventoryModule       *module.InventoryModule
	productModule         *module.SupplierProductModule
	unitSpecModule        *module.ProductUnitSpecModule
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
	unitSpecModule *module.ProductUnitSpecModule,
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
		unitSpecModule:        unitSpecModule,
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
	orderNo := fmt.Sprintf("DD%s%03d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000)

	// 记账日期由后端默认当前时间
	accountDate := time.Now()

	// 构建明细
	var items []model.StoreAccountItem
	var consumables []model.StoreAccountConsumable
	var totalAmount float64
	var consumableAmount float64
	productMap := make(map[uint]*model.SupplierProduct)

	for _, item := range req.Items {
		// 获取商品名称
		productName := ""
		unit := item.Unit
		price := item.Price
		var productUnitSpecs []*model.ProductUnitSpec
		var product *model.SupplierProduct
		if s.productModule != nil {
			if p, err := s.productModule.GetByID(item.ProductID); err == nil && p != nil {
				product = p
				productMap[item.ProductID] = p
				productName = p.Name
				if unit == "" {
					unit = p.Unit
				}
			}
		}
		if s.unitSpecModule != nil {
			if specs, err := s.unitSpecModule.ListByProductID(item.ProductID); err == nil {
				productUnitSpecs = specs
			}
		}

		// 严格模式：启用规格表后，单位必须匹配到规格售价
		if s.unitSpecModule != nil {
			specPrice, matched := tryResolveUnitSpecSalePrice(unit, productUnitSpecs)
			if !matched {
				name := productName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return nil, fmt.Errorf("商品【%s】单位【%s】未配置售价，请先在商品单位配置中维护该单位售价", name, unit)
			}
			price = specPrice
		} else if price <= 0 && product != nil {
			// 兼容兜底：未启用规格模块时，沿用旧逻辑
			price = resolveUnitPrice(unit, product)
		}

		// 计算金额
		amount := item.Amount
		// 只要能确定单价，就由后端统一重算金额，避免前端旧金额污染
		if price > 0 && item.Quantity > 0 {
			amount = price * item.Quantity
		}

		items = append(items, model.StoreAccountItem{
			ProductID:   item.ProductID,
			ProductName: productName,
			Spec:        item.Spec,
			Quantity:    item.Quantity,
			Unit:        unit,
			Price:       price,
			Amount:      amount,
			Remark:      item.Remark,
		})

		totalAmount += amount
	}

	for _, item := range req.Consumables {
		productName := ""
		unit := item.Unit
		price := item.Price
		var productUnitSpecs []*model.ProductUnitSpec
		var product *model.SupplierProduct
		if s.productModule != nil {
			if p, err := s.productModule.GetByID(item.ProductID); err == nil && p != nil {
				product = p
				productName = p.Name
				if unit == "" {
					unit = p.Unit
				}
			}
		}
		if s.unitSpecModule != nil {
			if specs, err := s.unitSpecModule.ListByProductID(item.ProductID); err == nil {
				productUnitSpecs = specs
			}
		}
		if s.unitSpecModule != nil {
			specPrice, matched := tryResolveUnitSpecSalePrice(unit, productUnitSpecs)
			if !matched {
				name := productName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return nil, fmt.Errorf("消耗品【%s】单位【%s】未配置售价，请先维护该单位售价", name, unit)
			}
			price = specPrice
		} else if price <= 0 && product != nil {
			price = resolveUnitPrice(unit, product)
		}
		amount := item.Amount
		if price > 0 && item.Quantity > 0 {
			amount = price * item.Quantity
		}
		consumables = append(consumables, model.StoreAccountConsumable{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Unit:        unit,
			Price:       price,
			Amount:      amount,
			Remark:      item.Remark,
		})
		consumableAmount += amount
	}

	account := &model.StoreAccount{
		AccountNo:          accountNo,
		StoreID:            storeID,
		Channel:            req.Channel,
		OrderNo:            orderNo,
		TotalAmount:        totalAmount,
		OtherExpenseAmount: req.OtherExpenseAmount,
		NetIncomeAmount:    totalAmount - req.OtherExpenseAmount - consumableAmount,
		ItemCount:          len(items),
		TagCode:            req.TagCode,
		TagName:            req.TagName,
		Remark:             req.Remark,
		OperatorID:         operatorID,
		AccountDate:        accountDate,
		Items:              items,
		Consumables:        consumables,
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
		var product *model.SupplierProduct
		if s.productModule != nil {
			if p, err := s.productModule.GetByID(item.ProductID); err == nil && p != nil {
				product = p
			}
		}
		baseQuantity, baseUnit := convertToBaseQuantity(s.unitSpecModule, product, item.ProductID, item.Quantity, item.Unit)
		inventoryOutOrder.TotalQuantity += baseQuantity
		inventoryOutOrder.Items = append(inventoryOutOrder.Items, model.InventoryOrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    baseQuantity,
			Unit:        baseUnit,
			Remark:      fmt.Sprintf("记账自动出库(原始: %.2f%s)", item.Quantity, item.Unit),
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

	// 兼容历史库存：若库存仍是“箱/桶”等旧单位，先换算到基础单位再扣减
	for _, outItem := range inventoryOutOrder.Items {
		inv, err := s.inventoryModule.GetByStoreAndProduct(storeID, outItem.ProductID)
		if err != nil || inv == nil {
			continue
		}
		product := productMap[outItem.ProductID]
		invUnit := strings.TrimSpace(inv.Unit)
		outUnit := strings.TrimSpace(outItem.Unit)

		// 场景1：单位不同，直接按库存当前单位换算为基础单位
		if invUnit != outUnit {
			convertedQty, baseUnit := convertToBaseQuantity(s.unitSpecModule, product, outItem.ProductID, inv.Quantity, inv.Unit)
			if convertedQty <= 0 {
				continue
			}
			if err := s.inventoryModule.UpdateQuantityAndUnit(inv.ID, convertedQty, baseUnit); err != nil {
				return nil, err
			}
			continue
		}

		// 场景2：单位相同但库存明显偏小（历史把“箱数”写进了“瓶单位”）
		if inv.Quantity < outItem.Quantity && s.unitSpecModule != nil {
			specs, err := s.unitSpecModule.ListByProductID(outItem.ProductID)
			if err != nil {
				continue
			}
			legacyFixed := false
			for _, spec := range specs {
				if spec == nil || !spec.IsEnabled || spec.FactorToBase <= 1 {
					continue
				}
				// 若“库存数量 * 大规格系数”能覆盖本次出库，按历史箱数纠偏一次
				candidate := inv.Quantity * spec.FactorToBase
				if candidate >= outItem.Quantity {
					if err := s.inventoryModule.UpdateQuantityAndUnit(inv.ID, candidate, outUnit); err != nil {
						return nil, err
					}
					legacyFixed = true
					break
				}
			}
			if legacyFixed {
				continue
			}
		}
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
	totalAmount, netIncomeAmount, count, err := s.storeAccountModule.GetStatsByDateRange(storeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_amount":      totalAmount,
		"net_income_amount": netIncomeAmount,
		"count":             count,
	}, nil
}

func (s *StoreAccountService) BindConsumables(accountID uint, req *model.BindStoreAccountConsumablesReq) error {
	if _, err := s.storeAccountModule.GetByID(accountID); err != nil {
		return err
	}
	consumables := make([]model.StoreAccountConsumable, 0, len(req.Consumables))
	for _, item := range req.Consumables {
		productName := ""
		unit := item.Unit
		price := item.Price
		var productUnitSpecs []*model.ProductUnitSpec
		var product *model.SupplierProduct
		if s.productModule != nil {
			if p, err := s.productModule.GetByID(item.ProductID); err == nil && p != nil {
				product = p
				productName = p.Name
				if unit == "" {
					unit = p.Unit
				}
			}
		}
		if s.unitSpecModule != nil {
			if specs, err := s.unitSpecModule.ListByProductID(item.ProductID); err == nil {
				productUnitSpecs = specs
			}
		}
		if s.unitSpecModule != nil {
			specPrice, matched := tryResolveUnitSpecSalePrice(unit, productUnitSpecs)
			if !matched {
				name := productName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return fmt.Errorf("消耗品【%s】单位【%s】未配置售价，请先维护该单位售价", name, unit)
			}
			price = specPrice
		} else if price <= 0 && product != nil {
			price = resolveUnitPrice(unit, product)
		}
		amount := item.Amount
		if price > 0 && item.Quantity > 0 {
			amount = price * item.Quantity
		}
		consumables = append(consumables, model.StoreAccountConsumable{
			AccountID:   accountID,
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Unit:        unit,
			Price:       price,
			Amount:      amount,
			Remark:      item.Remark,
		})
	}
	return s.storeAccountModule.ReplaceConsumables(accountID, consumables)
}
