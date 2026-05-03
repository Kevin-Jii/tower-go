package service

import (
	"errors"
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

func resolveUnitCostFromSpecs(unit string, specs []*model.ProductUnitSpec) float64 {
	if len(specs) == 0 {
		return 0
	}
	normalized := strings.ToLower(strings.TrimSpace(unit))
	for _, spec := range specs {
		if spec == nil || !spec.IsEnabled || spec.CostPrice < 0 {
			continue
		}
		if normalized == strings.ToLower(strings.TrimSpace(spec.UnitCode)) ||
			normalized == strings.ToLower(strings.TrimSpace(spec.UnitName)) {
			return spec.CostPrice
		}
	}
	if normalized != "" {
		for _, spec := range specs {
			if spec == nil || !spec.IsEnabled || spec.CostPrice < 0 {
				continue
			}
			code := strings.ToLower(strings.TrimSpace(spec.UnitCode))
			name := strings.ToLower(strings.TrimSpace(spec.UnitName))
			if strings.Contains(code, normalized) || strings.Contains(name, normalized) ||
				strings.Contains(normalized, code) || strings.Contains(normalized, name) {
				return spec.CostPrice
			}
		}
	}
	return 0
}

type StoreAccountService struct {
	storeAccountModule    *module.StoreAccountModule
	inventoryModule       *module.InventoryModule
	productModule         *module.SupplierProductModule
	unitSpecModule        *module.ProductUnitSpecModule
	storeModule           *module.StoreModule
	memberModule          *module.MemberModule
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
	memberModule *module.MemberModule,
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
		memberModule:          memberModule,
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
	if req.MemberID != nil && *req.MemberID > 0 && s.memberModule != nil {
		if _, err := s.memberModule.GetMember(*req.MemberID, storeID, false); err != nil {
			return nil, fmt.Errorf("会员不存在")
		}
	}

	accountNo := s.storeAccountModule.GenerateAccountNo()
	orderNo := fmt.Sprintf("DD%s%03d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000)

	// 记账日期由后端默认当前时间
	accountDate := time.Now()

	// 构建明细
	var items []model.StoreAccountItem
	var consumables []model.StoreAccountConsumable
	var totalAmount float64
	var consumableAmount float64
	var itemCostAmount float64
	productMap := make(map[uint]*model.SupplierProduct)

	for _, item := range req.Items {
		if item.ProductID == model.StoreAccountItemCustomProductID {
			name := strings.TrimSpace(item.ProductName)
			if name == "" {
				return nil, fmt.Errorf("自定义明细描述不能为空")
			}
			unit := strings.TrimSpace(item.Unit)
			if unit == "" {
				return nil, fmt.Errorf("自定义明细「%s」请填写单位", name)
			}
			price := item.Price
			if price <= 0 {
				return nil, fmt.Errorf("自定义明细「%s」请填写单价", name)
			}
			amount := item.Amount
			if amount <= 0 && item.Quantity > 0 {
				amount = price * item.Quantity
			}
			if amount <= 0 {
				return nil, fmt.Errorf("自定义明细「%s」金额无效", name)
			}
			items = append(items, model.StoreAccountItem{
				ProductID:   model.StoreAccountItemCustomProductID,
				ProductName: name,
				Spec:        strings.TrimSpace(item.Spec),
				Quantity:    item.Quantity,
				Unit:        unit,
				Price:       price,
				Amount:      amount,
				Remark:      strings.TrimSpace(item.Remark),
			})
			totalAmount += amount
			continue
		}

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
		if item.Quantity > 0 {
			itemCostAmount += item.Quantity * resolveUnitCostFromSpecs(unit, productUnitSpecs)
		}
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
		MemberID:           req.MemberID,
		PaymentStatus:      resolvePaymentStatus(req.PaymentStatus),
		Channel:            req.Channel,
		OrderNo:            orderNo,
		TotalAmount:        totalAmount,
		OtherExpenseAmount: req.OtherExpenseAmount,
		NetIncomeAmount:    totalAmount - req.OtherExpenseAmount - consumableAmount - itemCostAmount,
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
		ItemCount:     0,
		OperatorID:    operatorID,
	}
	for _, item := range items {
		if item.ProductID == model.StoreAccountItemCustomProductID {
			continue
		}
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
	inventoryOutOrder.ItemCount = len(inventoryOutOrder.Items)
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

	var outForTx *model.InventoryOrder = inventoryOutOrder
	if len(inventoryOutOrder.Items) == 0 {
		outForTx = nil
	}
	if err := s.storeAccountModule.CreateWithInventoryOut(account, outForTx); err != nil {
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
			name := strings.TrimSpace(item.ProductName)
			if name == "" {
				name = fmt.Sprintf("商品#%d", item.ProductID)
			}
			items = append(items, AccountItemData{
				Name:     name,
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
			OtherExpense: account.OtherExpenseAmount,
			NetIncome:    account.NetIncomeAmount,
			ItemCount:    account.ItemCount,
			Remark:       account.Remark,
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
	type templateItem struct {
		ProductName string
		Quantity    string
		Amount      string
	}
	templateItems := make([]templateItem, 0, len(account.Items))
	for i, item := range account.Items {
		line := fmt.Sprintf("%d. %s x%.2f%s = ¥%.2f", i+1, item.ProductName, item.Quantity, item.Unit, item.Amount)
		itemLines = append(itemLines, line)
		templateItems = append(templateItems, templateItem{
			ProductName: strings.TrimSpace(item.ProductName),
			Quantity:    fmt.Sprintf("%.2f%s", item.Quantity, strings.TrimSpace(item.Unit)),
			Amount:      fmt.Sprintf("¥%.2f", item.Amount),
		})
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
			"Items":        templateItems,
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

	cardTitle, cardText, cardButtonTitle, cardButtonURL := s.buildAccountActionCard(account, store.Name, operatorDisplay, channelName, itemLines, imageURL)

	// 发送通知：记账通知优先使用钉钉卡片，失败后回退到 Markdown，避免通知丢失
	var sendErr error
	if strings.TrimSpace(bot.CardMsgKey) != "" {
		itemListForCard := make([]map[string]interface{}, 0, len(account.Items))
		for _, it := range account.Items {
			name := strings.TrimSpace(it.ProductName)
			if name == "" {
				name = fmt.Sprintf("商品#%d", it.ProductID)
			}
			itemListForCard = append(itemListForCard, map[string]interface{}{
				"name":     name,
				"quantity": fmt.Sprintf("%.2f", it.Quantity),
				"unit":     strings.TrimSpace(it.Unit),
				"amount":   fmt.Sprintf("%.2f", it.Amount),
			})
		}
		accountBlock := map[string]interface{}{
			"account_no":    account.AccountNo,
			"channel":       channelName,
			"account_date":  account.AccountDate.Format("2006-01-02"),
			"other_expense": fmt.Sprintf("%.2f", account.OtherExpenseAmount),
			"net_income":    fmt.Sprintf("%.2f", account.NetIncomeAmount),
		}
		cardParam := map[string]interface{}{
			"title":        title,
			"storeName":    store.Name,
			"storename":    store.Name,
			"accountNo":    account.AccountNo,
			"channelName":  channelName,
			"accountDate":  account.AccountDate.Format("2006-01-02"),
			"operatorName": operatorDisplay,
			"content":      text,
			"item_list":    itemListForCard,
			"itemList":     strings.Join(itemLines, "\n"),
			"shangpinls":   strings.Join(itemLines, "\n"),
			"shangpinimg":  imageURL,
			"itemCount":    account.ItemCount,
			"totalAmount":  fmt.Sprintf("%.2f", account.TotalAmount),
			"createTime":   time.Now().Format("2006-01-02 15:04:05"),
			"imageUrl":     imageURL,
			"account":      accountBlock,
			// 兼容钉钉模板使用扁平点路径变量名的场景
			"account.account_no":    accountBlock["account_no"],
			"account.channel":       accountBlock["channel"],
			"account.account_date":  accountBlock["account_date"],
			"account.other_expense": accountBlock["other_expense"],
			"account.net_income":    accountBlock["net_income"],
			"account.total_amount":  fmt.Sprintf("%.2f", account.TotalAmount),
			"ccount.total_amount":   fmt.Sprintf("%.2f", account.TotalAmount),
		}
		sendErr = s.dingTalkService.SendStreamCardToMobile(bot, bot.CardMsgKey, store.Phone, cardParam)
		if sendErr != nil && logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to send custom account card, fallback to action card",
				"storeID", storeID,
				"accountNo", account.AccountNo,
				"cardMsgKey", bot.CardMsgKey,
				"error", sendErr,
			)
		}
	}
	if sendErr != nil || strings.TrimSpace(bot.CardMsgKey) == "" {
		sendErr = s.dingTalkService.SendStreamActionCardToMobile(bot, cardTitle, cardText, cardButtonTitle, cardButtonURL, store.Phone)
		if sendErr != nil && logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to send account action card, fallback to markdown",
				"storeID", storeID,
				"accountNo", account.AccountNo,
				"error", sendErr,
			)
		}
	}
	if sendErr != nil {
		if imageURL != "" {
			sendErr = s.dingTalkService.SendStreamMarkdownWithImageToMobile(bot, title, text, imageURL, store.Phone)
		} else {
			sendErr = s.dingTalkService.SendStreamMarkdownToMobile(bot, title, text, store.Phone)
		}
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
				"messageType", "card",
			)
		}
	}
}

func (s *StoreAccountService) buildAccountActionCard(account *model.StoreAccount, storeName, operatorName, channelName string, itemLines []string, imageURL string) (string, string, string, string) {
	title := fmt.Sprintf("新记账通知 - %s", storeName)
	buttonTitle := "查看记账回单"
	buttonURL := strings.TrimSpace(imageURL)
	if buttonURL == "" {
		buttonTitle = "查看详情"
		buttonURL = "https://www.dingtalk.com/"
	}

	detailLines := itemLines
	if len(detailLines) == 0 {
		detailLines = []string{"暂无商品明细"}
	}
	if len(detailLines) > 8 {
		detailLines = append(detailLines[:8], fmt.Sprintf("...等共 %d 项", len(itemLines)))
	}

	var b strings.Builder
	b.WriteString("### ")
	b.WriteString(title)
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("- 记账编号：%s\n", account.AccountNo))
	b.WriteString(fmt.Sprintf("- 渠道来源：%s\n", channelName))
	b.WriteString(fmt.Sprintf("- 记账日期：%s\n", account.AccountDate.Format("2006-01-02")))
	b.WriteString(fmt.Sprintf("- 操作人：%s\n\n", operatorName))
	b.WriteString("#### 记账明细\n")
	for _, line := range detailLines {
		b.WriteString("- ")
		b.WriteString(strings.TrimSpace(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("**合计：¥%.2f**\n\n", account.TotalAmount))
	if account.OtherExpenseAmount > 0 {
		b.WriteString(fmt.Sprintf("- 其他支出：¥%.2f\n", account.OtherExpenseAmount))
	}
	b.WriteString(fmt.Sprintf("- 净收入：¥%.2f\n", account.NetIncomeAmount))
	if strings.TrimSpace(account.Remark) != "" {
		b.WriteString(fmt.Sprintf("\n备注：%s\n", strings.TrimSpace(account.Remark)))
	}
	if imageURL != "" {
		b.WriteString("\n![记账回单](")
		b.WriteString(imageURL)
		b.WriteString(")\n")
	}
	b.WriteString("\n> 本消息由系统自动发送")

	return title, b.String(), buttonTitle, buttonURL
}

// Get 获取记账详情
func (s *StoreAccountService) Get(id uint) (*model.StoreAccount, error) {
	account, err := s.storeAccountModule.GetByID(id)
	if err != nil {
		return nil, err
	}
	account.NetIncomeAmount = s.calculateAccountNetIncome(account)
	return account, nil
}

// List 记账列表
func (s *StoreAccountService) List(req *model.ListStoreAccountReq) ([]*model.StoreAccount, int64, error) {
	list, total, err := s.storeAccountModule.List(req)
	if err != nil {
		return nil, 0, err
	}
	for _, account := range list {
		account.NetIncomeAmount = s.calculateAccountNetIncome(account)
	}
	return list, total, nil
}

// Update 更新记账
func (s *StoreAccountService) Update(id uint, req *model.UpdateStoreAccountReq) error {
	account, err := s.storeAccountModule.GetByID(id)
	if err != nil {
		return err
	}
	if !s.IsAccountEditable(account) {
		return errors.New("记账已超过可编辑时间，仅支持当天编辑（23:00后创建可延长至次日03:00）")
	}

	updates := make(map[string]interface{})

	if req.Channel != "" {
		updates["channel"] = req.Channel
	}
	if req.PaymentStatus != nil {
		updates["payment_status"] = resolvePaymentStatus(*req.PaymentStatus)
	}
	if req.MemberID != nil {
		if *req.MemberID > 0 {
			if s.memberModule != nil {
				if _, err := s.memberModule.GetMember(*req.MemberID, account.StoreID, false); err != nil {
					return fmt.Errorf("会员不存在")
				}
			}
			updates["member_id"] = *req.MemberID
		} else {
			updates["member_id"] = nil
		}
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
		var consumableTotal float64
		for _, c := range account.Consumables {
			consumableTotal += c.Amount
		}
		itemCostTotal := s.calculateAccountItemCost(account.Items)
		updates["net_income_amount"] = account.TotalAmount - *req.OtherExpenseAmount - consumableTotal - itemCostTotal
	}

	if len(updates) == 0 {
		return nil
	}

	return s.storeAccountModule.Update(id, updates)
}

// Delete 删除记账
func (s *StoreAccountService) Delete(id uint) error {
	return errors.New("记账记录不允许删除")
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

func (s *StoreAccountService) calculateAccountNetIncome(account *model.StoreAccount) float64 {
	if account == nil {
		return 0
	}
	var consumableTotal float64
	for _, c := range account.Consumables {
		consumableTotal += c.Amount
	}
	itemCostTotal := s.calculateAccountItemCost(account.Items)
	return account.TotalAmount - account.OtherExpenseAmount - consumableTotal - itemCostTotal
}

func (s *StoreAccountService) calculateAccountItemCost(items []model.StoreAccountItem) float64 {
	if len(items) == 0 || s.unitSpecModule == nil {
		return 0
	}
	specCache := make(map[uint][]*model.ProductUnitSpec)
	var total float64
	for _, it := range items {
		if it.ProductID == 0 || it.Quantity <= 0 {
			continue
		}
		specs, ok := specCache[it.ProductID]
		if !ok {
			rows, err := s.unitSpecModule.ListByProductID(it.ProductID)
			if err == nil {
				specs = rows
			}
			specCache[it.ProductID] = specs
		}
		costPrice := resolveUnitCostFromSpecs(it.Unit, specs)
		total += it.Quantity * costPrice
	}
	return total
}

func resolvePaymentStatus(v int) int {
	if v == model.StoreAccountPaymentUnpaid {
		return model.StoreAccountPaymentUnpaid
	}
	return model.StoreAccountPaymentPaid
}

func (s *StoreAccountService) BindConsumables(accountID uint, req *model.BindStoreAccountConsumablesReq) error {
	account, err := s.storeAccountModule.GetByID(accountID)
	if err != nil {
		return err
	}
	if !s.IsAccountEditable(account) {
		return errors.New("记账已超过可编辑时间，仅支持当天编辑（23:00后创建可延长至次日03:00）")
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

// IsAccountEditable 判断记账记录是否允许编辑：
// 1. 默认仅允许在记账创建当日（自然日）编辑
// 2. 若创建时间在当日23点（含）后，延长至次日03:00
func (s *StoreAccountService) IsAccountEditable(account *model.StoreAccount) bool {
	if account == nil {
		return false
	}
	now := time.Now()
	created := account.CreatedAt
	loc := now.Location()
	if !created.IsZero() {
		created = created.In(loc)
	}

	dayStart := time.Date(created.Year(), created.Month(), created.Day(), 0, 0, 0, 0, loc)
	nextDayStart := dayStart.Add(24 * time.Hour)
	cutoff := nextDayStart
	if created.Hour() >= 23 {
		cutoff = dayStart.Add(27 * time.Hour) // 次日03:00
	}
	return now.Before(cutoff)
}
