package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// MemberService 会员服务
type MemberService struct {
	module          *module.MemberModule
	storeModule     *module.StoreModule
	botModule       *module.DingTalkBotModule
	dictModule      *module.DictModule
	userModule      *module.UserModule
	dingTalkService *DingTalkService
}

// NewMemberService 创建会员服务
func NewMemberService(m *module.MemberModule) *MemberService {
	return &MemberService{module: m}
}

// SetDependencies 设置依赖（用于解耦初始化）
func (s *MemberService) SetDependencies(
	storeModule *module.StoreModule,
	botModule *module.DingTalkBotModule,
	dictModule *module.DictModule,
	userModule *module.UserModule,
	dingTalkService *DingTalkService,
) {
	s.storeModule = storeModule
	s.botModule = botModule
	s.dictModule = dictModule
	s.userModule = userModule
	s.dingTalkService = dingTalkService
}

// ========== Member 操作 ==========

// CreateMember 创建会员
func (s *MemberService) CreateMember(req *model.CreateMemberReq, storeID uint) (*model.Member, error) {
	return s.module.CreateMember(req, storeID)
}

// UpdateMember 更新会员
func (s *MemberService) UpdateMember(id uint, req *model.UpdateMemberReq, storeID uint, isAdmin bool) (*model.Member, error) {
	return s.module.UpdateMember(id, req, storeID, isAdmin)
}

// DeleteMember 删除会员
func (s *MemberService) DeleteMember(id uint, storeID uint, isAdmin bool) error {
	return s.module.DeleteMember(id, storeID, isAdmin)
}

// GetMember 获取会员
func (s *MemberService) GetMember(id uint, storeID uint, isAdmin bool) (*model.Member, error) {
	return s.module.GetMember(id, storeID, isAdmin)
}

// GetMemberByPhone 通过手机号获取会员
func (s *MemberService) GetMemberByPhone(phone string, storeID uint, isAdmin bool) (*model.Member, error) {
	return s.module.GetMemberByPhone(phone, storeID, isAdmin)
}

// GetMemberByUID 通过UID获取会员
func (s *MemberService) GetMemberByUID(uid string) (*model.Member, error) {
	return s.module.GetMemberByUID(uid)
}

// ListMembers 获取会员列表
func (s *MemberService) ListMembers(keyword string, page, pageSize int, storeID uint, isAdmin bool) ([]model.Member, int64, error) {
	return s.module.ListMembers(keyword, page, pageSize, storeID, isAdmin)
}

// AdjustBalance 调整余额
func (s *MemberService) AdjustBalance(id uint, amount model.DecimalType, changeType model.ChangeTypeEnum, remark string, version int, storeID, userID uint, isAdmin bool) (*model.Member, error) {
	member, err := s.module.AdjustBalanceWithLock(id, amount, changeType, remark, version, storeID, isAdmin)
	if err != nil {
		return nil, err
	}

	// 异步发送钉钉通知
	go s.sendAdjustBalanceDingTalkNotification(member, amount, changeType, remark, storeID, userID)

	return member, nil
}

// ========== WalletLog 操作 ==========

// CreateWalletLog 创建流水记录
func (s *MemberService) CreateWalletLog(log *model.WalletLog) (*model.WalletLog, error) {
	return s.module.CreateWalletLog(log)
}

// GetWalletLog 获取流水详情
func (s *MemberService) GetWalletLog(id uint) (*model.WalletLog, error) {
	return s.module.GetWalletLog(id)
}

// ListWalletLogs 查询流水列表
func (s *MemberService) ListWalletLogs(req *model.ListWalletLogReq, page, pageSize int, storeID uint, isAdmin bool) ([]model.WalletLog, int64, error) {
	return s.module.ListWalletLogs(req, page, pageSize, storeID, isAdmin)
}

// ========== RechargeOrder 操作 ==========

// CreateRechargeOrder 创建充值单（自动完成支付）
func (s *MemberService) CreateRechargeOrder(req *model.CreateRechargeOrderReq, storeID, userID uint, isAdmin bool) (*model.RechargeOrder, error) {
	// 1. 创建充值单
	order, err := s.module.CreateRechargeOrder(req, storeID, isAdmin)
	if err != nil {
		return nil, err
	}

	// 2. 自动完成支付（更新会员余额、记录流水）
	order, err = s.module.PayRechargeOrder(order.OrderNo, storeID, isAdmin)
	if err != nil {
		return nil, err
	}

	// 3. 异步发送钉钉通知
	go s.sendRechargeDingTalkNotification(order, storeID, userID)

	return order, nil
}

// GetRechargeOrder 获取充值单
func (s *MemberService) GetRechargeOrder(id uint, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	return s.module.GetRechargeOrder(id, storeID, isAdmin)
}

// GetRechargeOrderByNo 通过单号获取充值单
func (s *MemberService) GetRechargeOrderByNo(orderNo string, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	return s.module.GetRechargeOrderByNo(orderNo, storeID, isAdmin)
}

// ListRechargeOrders 查询充值单列表
func (s *MemberService) ListRechargeOrders(memberID uint, status *model.PayStatusEnum, page, pageSize int, storeID uint, isAdmin bool) ([]model.RechargeOrder, int64, error) {
	return s.module.ListRechargeOrders(memberID, status, page, pageSize, storeID, isAdmin)
}

// PayRechargeOrder 支付充值单
func (s *MemberService) PayRechargeOrder(orderNo string, storeID, userID uint, isAdmin bool) (*model.RechargeOrder, error) {
	order, err := s.module.PayRechargeOrder(orderNo, storeID, isAdmin)
	if err != nil {
		return nil, err
	}

	// 异步发送钉钉通知
	go s.sendRechargeDingTalkNotification(order, storeID, userID)

	return order, nil
}

// sendRechargeDingTalkNotification 发送充值通知到门店钉钉群
func (s *MemberService) sendRechargeDingTalkNotification(order *model.RechargeOrder, storeID, userID uint) {
	if s.dingTalkService == nil || s.storeModule == nil || s.botModule == nil {
		return
	}

	// 获取门店信息
	store, err := s.storeModule.GetByID(storeID)
	if err != nil || store == nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to get store for recharge notification", "storeID", storeID, "error", err)
		}
		return
	}

	// 检查门店是否有联系电话
	if store.Phone == "" {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Store has no phone, skip recharge notification", "storeID", storeID)
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

	// 获取操作人名称
	operatorName := ""
	if s.userModule != nil && userID > 0 {
		if user, err := s.userModule.GetByID(userID); err == nil && user != nil {
			operatorName = user.Nickname
			if operatorName == "" {
				operatorName = user.Username
			}
		}
	}
	if operatorName == "" {
		operatorName = "未知"
	}

	// 获取会员名称
	memberName := order.MemberName
	if memberName == "" {
		memberName = fmt.Sprintf("会员%d", order.MemberID)
	}

	// 获取支付方式名称（使用字典 MDGL_ZFFS 翻译）
	payTypeName := getPayTypeName(int(order.PayType))
	if s.dictModule != nil {
		if dictData, err := s.dictModule.GetDataByTypeAndValue("MDGL_ZFFS", fmt.Sprintf("%d", order.PayType)); err == nil && dictData != nil {
			payTypeName = dictData.Label
		}
	}

	// 格式化金额
	payAmountStr := order.PayAmount.String()
	giftAmountStr := order.GiftAmount.String()

	// 构造通知消息
	title := "会员充值通知"
	text := fmt.Sprintf(`## 💰 会员充值通知

**门店**: %s

**会员**: %s 充值了 **%s** 元

**赠送金额**: %s 元

**支付方式**: %s

**操作人**: %s

**时间**: %s`,
		store.Name,
		memberName,
		payAmountStr,
		giftAmountStr,
		payTypeName,
		operatorName,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	// 发送钉钉消息
	if err := s.dingTalkService.SendStreamMarkdownToMobile(bot, title, text, store.Phone); err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to send recharge notification",
				"orderNo", order.OrderNo,
				"storeID", storeID,
				"error", err)
		}
		return
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Recharge notification sent successfully",
			"orderNo", order.OrderNo,
			"memberName", memberName,
			"storeName", store.Name)
	}
}

// sendAdjustBalanceDingTalkNotification 发送余额调整通知到门店钉钉群
func (s *MemberService) sendAdjustBalanceDingTalkNotification(member *model.Member, amount model.DecimalType, changeType model.ChangeTypeEnum, remark string, storeID, userID uint) {
	if s.dingTalkService == nil || s.storeModule == nil || s.botModule == nil {
		return
	}

	// 获取门店信息
	store, err := s.storeModule.GetByID(storeID)
	if err != nil || store == nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to get store for adjust balance notification", "storeID", storeID, "error", err)
		}
		return
	}

	// 检查门店是否有联系电话
	if store.Phone == "" {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Store has no phone, skip adjust balance notification", "storeID", storeID)
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

	// 获取操作人名称
	operatorName := ""
	if s.userModule != nil && userID > 0 {
		if user, err := s.userModule.GetByID(userID); err == nil && user != nil {
			operatorName = user.Nickname
			if operatorName == "" {
				operatorName = user.Username
			}
		}
	}
	if operatorName == "" {
		operatorName = "未知"
	}

	// 获取会员名称
	memberName := member.Name
	if memberName == "" {
		memberName = member.Phone
	}
	if memberName == "" {
		memberName = fmt.Sprintf("会员%d", member.ID)
	}

	// 判断是调增还是调减
	changeTypeName := "余额调整"
	changeSymbol := ""
	if changeType == model.ChangeTypeAdjustAdd {
		changeTypeName = "余额调增"
		changeSymbol = "+"
	} else if changeType == model.ChangeTypeAdjustLess {
		changeTypeName = "余额调减"
		changeSymbol = "-"
	}

	// 格式化金额
	amountStr := fmt.Sprintf("%s%s 元", changeSymbol, amount.String())

	// 构造通知消息
	title := "余额调整通知"
	text := fmt.Sprintf(`## 💰 余额调整通知

**门店**: %s

**会员**: %s

**调整类型**: %s

**调整金额**: %s

**当前余额**: %s 元

**备注**: %s

**操作人**: %s

**时间**: %s`,
		store.Name,
		memberName,
		changeTypeName,
		amountStr,
		member.Balance.String(),
		remark,
		operatorName,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	// 发送钉钉消息
	if err := s.dingTalkService.SendStreamMarkdownToMobile(bot, title, text, store.Phone); err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to send adjust balance notification",
				"memberID", member.ID,
				"storeID", storeID,
				"error", err)
		}
		return
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Adjust balance notification sent successfully",
			"memberID", member.ID,
			"memberName", memberName,
			"storeName", store.Name)
	}
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

// CancelRechargeOrder 取消充值单
func (s *MemberService) CancelRechargeOrder(orderNo string, storeID uint, isAdmin bool) (*model.RechargeOrder, error) {
	return s.module.CancelRechargeOrder(orderNo, storeID, isAdmin)
}

func (s *MemberService) ListMemberConsumptions(memberID uint, startDate, endDate string, page, pageSize int, storeID uint, isAdmin bool) ([]model.MemberConsumptionRecord, int64, *model.MemberConsumptionSummary, error) {
	member, err := s.module.GetMember(memberID, storeID, isAdmin)
	if err != nil || member == nil {
		return nil, 0, nil, errors.New("member not found")
	}

	records, total, summary, err := s.module.ListMemberConsumptions(memberID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, nil, err
	}
	if s.dictModule != nil {
		for i := range records {
			records[i].ChannelName = records[i].Channel
			if records[i].Channel != "" {
				if dictData, e := s.dictModule.GetDataByTypeAndValue("sales_channel", records[i].Channel); e == nil && dictData != nil {
					records[i].ChannelName = dictData.Label
				}
			}
		}
	}
	return records, total, summary, nil
}
