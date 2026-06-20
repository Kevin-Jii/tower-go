package module

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MemberModule 会员模块
type MemberModule struct {
	db *gorm.DB
}

// NewMemberModule 创建会员模块
func NewMemberModule(db *gorm.DB) *MemberModule {
	return &MemberModule{db: db}
}

// GetDB 返回底层数据库实例
func (m *MemberModule) GetDB() *gorm.DB {
	return m.db
}

func (m *MemberModule) scopedMemberQuery(storeID uint, isAdmin bool) *gorm.DB {
	q := m.db.Model(&model.Member{})
	if !isAdmin {
		q = q.Where("store_id = ?", storeID)
	}
	return q
}

func (m *MemberModule) scopedWalletLogQuery(storeID uint, isAdmin bool) *gorm.DB {
	q := m.db.Model(&model.WalletLog{}).Joins("JOIN t_member ON t_member.id = t_member_wallet_log.member_id")
	if !isAdmin {
		q = q.Where("t_member.store_id = ?", storeID)
	}
	return q
}

func (m *MemberModule) scopedRechargeOrderQuery(storeID uint, isAdmin bool) *gorm.DB {
	q := m.db.Model(&model.RechargeOrder{}).Joins("JOIN t_member ON t_member.id = t_recharge_order.member_id")
	if !isAdmin {
		q = q.Where("t_member.store_id = ?", storeID)
	}
	return q
}

func (m *MemberModule) scopedWineStorageQuery(storeID uint, isAdmin bool) *gorm.DB {
	q := m.db.Model(&model.MemberWineStorage{}).Preload("Member")
	if !isAdmin {
		q = q.Where("member_wine_storages.store_id = ?", storeID)
	}
	return q
}

func (m *MemberModule) scopedWineTransactionQuery(storeID uint, isAdmin bool) *gorm.DB {
	q := m.db.Model(&model.MemberWineTransaction{}).Preload("Member")
	if !isAdmin {
		q = q.Where("member_wine_transactions.store_id = ?", storeID)
	}
	return q
}

// ========== Member 操作 ==========

// CreateMember 创建会员
func (m *MemberModule) CreateMember(req *model.CreateMemberReq, storeID uint) (*model.Member, error) {
	// 检查手机号是否已存在
	var count int64
	m.db.Model(&model.Member{}).Where("store_id = ? AND phone = ?", storeID, req.Phone).Count(&count)
	if count > 0 {
		return nil, errors.New("手机号已注册")
	}

	// 如果没有提供 UID，则自动生成
	uid := req.UID
	if uid == "" {
		uid = uuid.NewString()
	}

	// 设置默认等级
	level := 1
	if req.Level != nil {
		level = *req.Level
	}

	member := &model.Member{
		StoreID: storeID,
		UID:     uid,
		Name:    req.Name,
		Phone:   req.Phone,
		Balance: model.DecimalZero(),
		Points:  0,
		Level:   level,
		Version: 0,
	}
	if err := m.db.Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

// UpdateMember 更新会员
func (m *MemberModule) UpdateMember(id uint, req *model.UpdateMemberReq, storeID uint, isAdmin bool) (*model.Member, error) {
	member, err := m.GetMember(id, storeID, isAdmin)
	if err != nil {
		return nil, err
	}
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return member, nil
	}
	if err := m.db.Model(member).Updates(updateMap).Error; err != nil {
		return nil, err
	}
	updated, err := m.GetMember(id, storeID, isAdmin)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// DeleteMember 删除会员
func (m *MemberModule) DeleteMember(id uint, storeID uint, isAdmin bool) error {
	member, err := m.GetMember(id, storeID, isAdmin)
	if err != nil {
		return err
	}
	return m.db.Delete(&model.Member{}, member.ID).Error
}

// GetMember 获取会员
func (m *MemberModule) GetMember(id uint, storeID uint, isAdmin bool) (*model.Member, error) {
	var member model.Member
	query := m.scopedMemberQuery(storeID, isAdmin)
	if err := query.Where("id = ?", id).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMemberByPhone 通过手机号获取会员
func (m *MemberModule) GetMemberByPhone(phone string, storeID uint, isAdmin bool) (*model.Member, error) {
	var member model.Member
	query := m.scopedMemberQuery(storeID, isAdmin)
	if err := query.Where("phone = ?", phone).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMemberByUID 通过UID获取会员
func (m *MemberModule) GetMemberByUID(uid string) (*model.Member, error) {
	var member model.Member
	if err := m.db.Where("uid = ?", uid).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// ListMembers 获取会员列表
func (m *MemberModule) ListMembers(keyword string, page, pageSize int, storeID uint, isAdmin bool) ([]model.Member, int64, error) {
	var members []model.Member
	var total int64

	query := m.scopedMemberQuery(storeID, isAdmin)
	if keyword != "" {
		query = query.Where("phone LIKE ? OR uid LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	query = query.Order("id DESC")
	query = query.Offset((page - 1) * pageSize)
	query = query.Limit(pageSize)

	if err := query.Find(&members).Error; err != nil {
		return nil, 0, err
	}
	return members, total, nil
}

// ListWineStorages 查询会员存酒当前存量
func (m *MemberModule) ListWineStorages(req *model.ListMemberWineStorageReq, storeID uint, isAdmin bool) ([]model.MemberWineStorage, int64, error) {
	rows := make([]model.MemberWineStorage, 0)
	var total int64
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := m.scopedWineStorageQuery(storeID, isAdmin)
	if req.StoreID > 0 && isAdmin {
		query = query.Where("member_wine_storages.store_id = ?", req.StoreID)
	}
	if req.MemberID > 0 {
		query = query.Where("member_wine_storages.member_id = ?", req.MemberID)
	}
	if req.OnlyStock == 1 {
		query = query.Where("member_wine_storages.quantity > 0")
	}
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		like := "%" + kw + "%"
		query = query.Joins("LEFT JOIN t_member AS member_search ON member_search.id = member_wine_storages.member_id").
			Where("member_wine_storages.wine_name LIKE ? OR member_search.phone LIKE ? OR member_search.name LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("member_wine_storages.updated_at DESC, member_wine_storages.id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// AdjustWineStorage 存入或取出会员存酒
func (m *MemberModule) AdjustWineStorage(storeID, operatorID uint, operatorName string, isAdmin bool, txnType int, req *model.MemberWineAdjustReq) (*model.MemberWineStorage, error) {
	if req == nil {
		return nil, errors.New("请求不能为空")
	}
	wineName := strings.TrimSpace(req.WineName)
	if wineName == "" {
		return nil, errors.New("请填写酒品名称")
	}
	unit := strings.TrimSpace(req.Unit)
	if unit == "" {
		unit = "瓶"
	}
	if req.Quantity <= 0 {
		return nil, errors.New("数量必须大于0")
	}
	member, err := m.GetMember(req.MemberID, storeID, isAdmin)
	if err != nil {
		return nil, errors.New("会员不存在")
	}
	if !isAdmin && member.StoreID != storeID {
		return nil, errors.New("会员不属于当前门店")
	}
	realStoreID := member.StoreID
	if realStoreID == 0 {
		realStoreID = storeID
	}

	var storage model.MemberWineStorage
	if err := m.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("store_id = ? AND member_id = ? AND wine_name = ? AND unit = ?", realStoreID, req.MemberID, wineName, unit).
			First(&storage).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if txnType == model.MemberWineTxnWithdraw {
				return errors.New("该会员暂无此酒品存量")
			}
			storage = model.MemberWineStorage{
				StoreID:  realStoreID,
				MemberID: req.MemberID,
				WineName: wineName,
				Unit:     unit,
				Quantity: 0,
				Remark:   strings.TrimSpace(req.Remark),
			}
			if err := tx.Create(&storage).Error; err != nil {
				return err
			}
		}

		nextQty := storage.Quantity
		switch txnType {
		case model.MemberWineTxnDeposit:
			nextQty += req.Quantity
		case model.MemberWineTxnWithdraw:
			if storage.Quantity < req.Quantity {
				return fmt.Errorf("存酒数量不足，当前剩余 %.2f%s", storage.Quantity, storage.Unit)
			}
			nextQty -= req.Quantity
		default:
			return errors.New("不支持的存酒操作类型")
		}

		if err := tx.Model(&storage).Updates(map[string]interface{}{
			"quantity": nextQty,
			"remark":   strings.TrimSpace(req.Remark),
		}).Error; err != nil {
			return err
		}
		storage.Quantity = nextQty
		storage.Remark = strings.TrimSpace(req.Remark)
		txn := &model.MemberWineTransaction{
			StoreID:      realStoreID,
			StorageID:    storage.ID,
			MemberID:     req.MemberID,
			Type:         txnType,
			WineName:     wineName,
			Unit:         unit,
			Quantity:     req.Quantity,
			BalanceAfter: nextQty,
			Remark:       strings.TrimSpace(req.Remark),
			OperatorID:   operatorID,
			OperatorName: operatorName,
		}
		return tx.Create(txn).Error
	}); err != nil {
		return nil, err
	}

	return m.GetWineStorage(storage.ID, storeID, isAdmin)
}

func (m *MemberModule) GetWineStorage(id uint, storeID uint, isAdmin bool) (*model.MemberWineStorage, error) {
	var row model.MemberWineStorage
	query := m.scopedWineStorageQuery(storeID, isAdmin)
	if err := query.Where("member_wine_storages.id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ListWineTransactions 查询存取酒流水
func (m *MemberModule) ListWineTransactions(req *model.ListMemberWineTransactionReq, storeID uint, isAdmin bool) ([]model.MemberWineTransaction, int64, error) {
	rows := make([]model.MemberWineTransaction, 0)
	var total int64
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := m.scopedWineTransactionQuery(storeID, isAdmin)
	if req.StoreID > 0 && isAdmin {
		query = query.Where("member_wine_transactions.store_id = ?", req.StoreID)
	}
	if req.StorageID > 0 {
		query = query.Where("member_wine_transactions.storage_id = ?", req.StorageID)
	}
	if req.MemberID > 0 {
		query = query.Where("member_wine_transactions.member_id = ?", req.MemberID)
	}
	if req.Type > 0 {
		query = query.Where("member_wine_transactions.type = ?", req.Type)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(member_wine_transactions.created_at) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(member_wine_transactions.created_at) <= ?", req.EndDate)
	}
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		like := "%" + kw + "%"
		query = query.Joins("LEFT JOIN t_member AS member_search ON member_search.id = member_wine_transactions.member_id").
			Where("member_wine_transactions.wine_name LIKE ? OR member_search.phone LIKE ? OR member_search.name LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("member_wine_transactions.created_at DESC, member_wine_transactions.id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// AdjustBalanceWithLock 乐观锁调整余额
func (m *MemberModule) AdjustBalanceWithLock(id uint, amount model.DecimalType, changeType model.ChangeTypeEnum, remark string, version int, storeID uint, isAdmin bool) (*model.Member, error) {
	var member model.Member
	query := m.scopedMemberQuery(storeID, isAdmin).Clauses(clause.Locking{Strength: "UPDATE"})
	if err := query.Where("id = ?", id).First(&member).Error; err != nil {
		return nil, err
	}

	// 乐观锁校验
	if member.Version != version {
		return nil, errors.New("数据已被修改，请刷新后重试")
	}

	// 计算新余额
	var newBalance model.DecimalType
	if changeType == model.ChangeTypeAdjustAdd {
		newBalance = member.Balance.Add(amount)
	} else if changeType == model.ChangeTypeAdjustLess {
		newBalance = member.Balance.Sub(amount)
		if newBalance.LessThan(model.DecimalZero()) {
			return nil, errors.New("余额不足")
		}
	} else {
		return nil, errors.New("不支持的调整类型")
	}

	// 更新余额和版本号
	if err := m.db.Model(&member).Updates(map[string]interface{}{
		"balance": newBalance,
		"version": member.Version + 1,
	}).Error; err != nil {
		return nil, err
	}

	// 记录流水
	walletLog := &model.WalletLog{
		MemberID:       id,
		ChangeType:     changeType,
		ChangeAmount:   amount,
		BalanceAfter:   newBalance,
		RelatedOrderNo: "",
		Remark:         remark,
	}
	if err := m.db.Create(walletLog).Error; err != nil {
		return nil, err
	}

	// 重新查询
	if err := m.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// ========== WalletLog 操作 ==========

// CreateWalletLog 创建流水记录
func (m *MemberModule) CreateWalletLog(log *model.WalletLog) (*model.WalletLog, error) {
	if err := m.db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

// GetWalletLog 获取流水详情
func (m *MemberModule) GetWalletLog(id uint) (*model.WalletLog, error) {
	var log model.WalletLog
	if err := m.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// ListWalletLogs 查询流水列表
func (m *MemberModule) ListWalletLogs(req *model.ListWalletLogReq, page, pageSize int, storeID uint, isAdmin bool) ([]model.WalletLog, int64, error) {
	var logs []model.WalletLog
	var total int64

	query := m.scopedWalletLogQuery(storeID, isAdmin)
	if req.MemberID > 0 {
		query = query.Where("member_id = ?", req.MemberID)
	}
	if req.ChangeType != nil {
		query = query.Where("change_type = ?", *req.ChangeType)
	}
	if req.StartTime != nil {
		query = query.Where("create_time >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("create_time <= ?", *req.EndTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	query = query.Order("id DESC")
	query = query.Offset((page - 1) * pageSize)
	query = query.Limit(pageSize)

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ========== RechargeOrder 操作 ==========

// CreateRechargeOrder 创建充值单
func (m *MemberModule) CreateRechargeOrder(req *model.CreateRechargeOrderReq, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	if _, err := m.GetMember(req.MemberID, storeID, isAdmin); err != nil {
		return nil, err
	}
	orderNo := generateOrderNo()
	order := &model.RechargeOrder{
		OrderNo:     orderNo,
		MemberID:    req.MemberID,
		PayAmount:   req.PayAmount,
		GiftAmount:  req.GiftAmount,
		TotalAmount: req.PayAmount.Add(req.GiftAmount),
		PayStatus:   model.PayStatusPending,
		PayType:     req.PayType,
		Remark:      req.Remark,
	}
	if err := m.db.Create(order).Error; err != nil {
		return nil, err
	}

	// 关联查询会员信息
	var member model.Member
	if err := m.db.Where("id = ?", order.MemberID).First(&member).Error; err == nil {
		order.MemberName = member.Name
		order.MemberPhone = member.Phone
	}

	// 设置状态名称和支付方式名称
	order.StatusName = order.PayStatus.String()
	order.PayTypeName = getPayTypeName(req.PayType)

	return order, nil
}

// getPayTypeName 获取支付方式名称
func getPayTypeName(payType int) string {
	switch payType {
	case 1:
		return "微信支付"
	case 2:
		return "支付宝"
	case 3:
		return "现金"
	case 4:
		return "银行卡"
	default:
		return "其他"
	}
}

// GetRechargeOrder 获取充值单
func (m *MemberModule) GetRechargeOrder(id uint, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	query := m.scopedRechargeOrderQuery(storeID, isAdmin)
	if err := query.Where("t_recharge_order.id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetRechargeOrderByNo 通过单号获取充值单
func (m *MemberModule) GetRechargeOrderByNo(orderNo string, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	query := m.scopedRechargeOrderQuery(storeID, isAdmin)
	if err := query.Where("t_recharge_order.order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// ListRechargeOrders 查询充值单列表
func (m *MemberModule) ListRechargeOrders(memberID uint, status *model.PayStatusEnum, page, pageSize int, storeID uint, isAdmin bool) ([]model.RechargeOrder, int64, error) {
	var orders []model.RechargeOrder
	var total int64

	query := m.scopedRechargeOrderQuery(storeID, isAdmin)
	if memberID > 0 {
		query = query.Where("t_recharge_order.member_id = ?", memberID)
	}
	if status != nil {
		query = query.Where("t_recharge_order.pay_status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("id DESC")
	query = query.Offset((page - 1) * pageSize)
	query = query.Limit(pageSize)

	if err := query.Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// PayRechargeOrder 支付充值单（余额充值）
func (m *MemberModule) PayRechargeOrder(orderNo string, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	query := m.scopedRechargeOrderQuery(storeID, isAdmin)
	if err := query.Where("t_recharge_order.order_no = ? AND t_recharge_order.pay_status = ?", orderNo, model.PayStatusPending).First(&order).Error; err != nil {
		return nil, err
	}

	// 获取会员并加锁
	var member model.Member
	if err := m.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&member, order.MemberID).Error; err != nil {
		return nil, err
	}

	// 设置会员名称到订单中（用于通知）
	order.MemberName = member.Name

	// 计算新余额
	newBalance := member.Balance.Add(order.TotalAmount)

	// 更新会员余额
	if err := m.db.Model(&member).Updates(map[string]interface{}{
		"balance": newBalance,
		"version": member.Version + 1,
	}).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	// 更新订单状态
	if err := m.db.Model(&order).Updates(map[string]interface{}{
		"pay_status": model.PayStatusPaid,
		"pay_time":   &now,
	}).Error; err != nil {
		return nil, err
	}

	// 记录流水
	walletLog := &model.WalletLog{
		MemberID:       order.MemberID,
		ChangeType:     model.ChangeTypeRecharge,
		ChangeAmount:   order.TotalAmount,
		BalanceAfter:   newBalance,
		RelatedOrderNo: orderNo,
		Remark:         "充值",
	}
	if err := m.db.Create(walletLog).Error; err != nil {
		return nil, err
	}

	// 重新查询订单
	if err := m.db.First(&order, order.ID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// CancelRechargeOrder 取消充值单
func (m *MemberModule) CancelRechargeOrder(orderNo string, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	query := m.scopedRechargeOrderQuery(storeID, isAdmin)
	if err := query.Where("t_recharge_order.order_no = ? AND t_recharge_order.pay_status = ?", orderNo, model.PayStatusPending).First(&order).Error; err != nil {
		return nil, err
	}

	if err := m.db.Model(&order).Update("pay_status", model.PayStatusCancelled).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// ListMemberConsumptions 查询会员消费记录（来自门店记账）
func (m *MemberModule) ListMemberConsumptions(memberID uint, startDate, endDate string, page, pageSize int) ([]model.MemberConsumptionRecord, int64, *model.MemberConsumptionSummary, error) {
	records := make([]model.MemberConsumptionRecord, 0)
	var total int64
	summary := &model.MemberConsumptionSummary{}

	base := m.db.Table("store_accounts AS sa").
		Where("sa.member_id = ?", memberID)
	if startDate != "" {
		base = base.Where("sa.account_date >= ?", startDate)
	}
	if endDate != "" {
		base = base.Where("sa.account_date <= ?", endDate)
	}

	if err := base.Count(&total).Error; err != nil {
		return records, 0, nil, err
	}

	consSub := m.db.Table("store_account_consumables").
		Select("account_id, COALESCE(SUM(amount),0) AS consumable_amount").
		Group("account_id")
	costSub := m.db.Table("store_account_items AS sai").
		Select("sai.account_id, COALESCE(SUM(sai.quantity * COALESCE(ps.cost_price,0)),0) AS cost_amount").
		Joins("LEFT JOIN product_unit_specs AS ps ON ps.product_id = sai.product_id AND ps.is_enabled = 1 AND (ps.unit_code = sai.unit OR ps.unit_name = sai.unit)").
		Group("sai.account_id")

	offset := (page - 1) * pageSize
	listQuery := m.db.Table("store_accounts AS sa").
		Select(`sa.id AS account_id, sa.account_no, sa.account_date, sa.channel, sa.order_no,
			sa.total_amount, sa.other_expense_amount, sa.round_amount, sa.gift_wine_cost_amount,
			COALESCE(cons.consumable_amount,0) AS consumable_amount,
			(sa.total_amount - sa.other_expense_amount - sa.errand_fee - COALESCE(cons.consumable_amount,0) - COALESCE(costs.cost_amount,0) - sa.gift_wine_cost_amount - sa.round_amount) AS net_income_amount,
			sa.created_at`).
		Joins("LEFT JOIN (?) AS cons ON cons.account_id = sa.id", consSub).
		Joins("LEFT JOIN (?) AS costs ON costs.account_id = sa.id", costSub).
		Where("sa.member_id = ?", memberID)
	if startDate != "" {
		listQuery = listQuery.Where("sa.account_date >= ?", startDate)
	}
	if endDate != "" {
		listQuery = listQuery.Where("sa.account_date <= ?", endDate)
	}
	if err := listQuery.Order("sa.id DESC").Offset(offset).Limit(pageSize).Scan(&records).Error; err != nil {
		return records, 0, nil, err
	}

	summaryQuery := m.db.Table("store_accounts AS sa").
		Select(`COUNT(1) AS count,
			COALESCE(SUM(sa.total_amount),0) AS total_amount,
			COALESCE(SUM(sa.other_expense_amount),0) AS other_expense_amount,
			COALESCE(SUM(sa.round_amount),0) AS round_amount,
			COALESCE(SUM(sa.gift_wine_cost_amount),0) AS gift_wine_cost_amount,
			COALESCE(SUM(COALESCE(cons.consumable_amount,0)),0) AS consumable_amount,
			COALESCE(SUM(sa.total_amount - sa.other_expense_amount - sa.errand_fee - COALESCE(cons.consumable_amount,0) - COALESCE(costs.cost_amount,0) - sa.gift_wine_cost_amount - sa.round_amount),0) AS net_income_amount`).
		Joins("LEFT JOIN (?) AS cons ON cons.account_id = sa.id", consSub).
		Joins("LEFT JOIN (?) AS costs ON costs.account_id = sa.id", costSub).
		Where("sa.member_id = ?", memberID)
	if startDate != "" {
		summaryQuery = summaryQuery.Where("sa.account_date >= ?", startDate)
	}
	if endDate != "" {
		summaryQuery = summaryQuery.Where("sa.account_date <= ?", endDate)
	}
	if err := summaryQuery.Scan(summary).Error; err != nil {
		return records, 0, nil, err
	}

	return records, total, summary, nil
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
