package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/businessdate"
)

type StoreExpenseService struct {
	expenseModule *module.StoreExpenseModule
	dictModule    *module.DictModule
	userModule    *module.UserModule
}

func NewStoreExpenseService(expenseModule *module.StoreExpenseModule, dictModule *module.DictModule, userModule *module.UserModule) *StoreExpenseService {
	return &StoreExpenseService{expenseModule: expenseModule, dictModule: dictModule, userModule: userModule}
}

func (s *StoreExpenseService) Create(storeID, operatorID uint, req *model.CreateStoreExpenseReq, hqUnbound bool) (*model.StoreExpense, error) {
	record, err := s.buildRecord(storeID, operatorID, hqUnbound, req.StoreID, req.CategoryCode, req.Amount, req.Remark)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 3; i++ {
		record.ExpenseNo = s.expenseModule.GenerateExpenseNo()
		if err := s.expenseModule.Create(record); err != nil {
			if module.IsDuplicateKeyError(err) {
				continue
			}
			return nil, err
		}
		return s.expenseModule.GetByIDScoped(record.ID, record.StoreID, true)
	}
	record.ExpenseNo = s.expenseModule.GenerateExpenseNo()
	if err := s.expenseModule.Create(record); err != nil {
		return nil, err
	}
	return s.expenseModule.GetByIDScoped(record.ID, record.StoreID, true)
}

func (s *StoreExpenseService) Update(id, storeID, operatorID uint, req *model.UpdateStoreExpenseReq, hqUnbound bool) (*model.StoreExpense, error) {
	existing, err := s.expenseModule.GetByIDScoped(id, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	updates := make(map[string]interface{})
	if strings.TrimSpace(req.CategoryCode) != "" {
		code, name, err := s.resolveCategory(req.CategoryCode)
		if err != nil {
			return nil, err
		}
		updates["category_code"] = code
		updates["category_name"] = name
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	updates["remark"] = strings.TrimSpace(req.Remark)
	if len(updates) == 0 {
		return existing, nil
	}
	if err := s.expenseModule.Update(id, storeID, hqUnbound, updates); err != nil {
		return nil, err
	}
	_ = operatorID
	return s.expenseModule.GetByIDScoped(id, storeID, hqUnbound)
}

func (s *StoreExpenseService) Delete(id, storeID uint, hqUnbound bool) error {
	return s.expenseModule.Delete(id, storeID, hqUnbound)
}

func (s *StoreExpenseService) Get(id, storeID uint, hqUnbound bool) (*model.StoreExpense, error) {
	return s.expenseModule.GetByIDScoped(id, storeID, hqUnbound)
}

func (s *StoreExpenseService) List(ctx context.Context, req *model.ListStoreExpenseReq) ([]*model.StoreExpense, int64, error) {
	_ = ctx
	return s.expenseModule.List(req)
}

func (s *StoreExpenseService) Stats(req *model.ListStoreExpenseReq) (*model.StoreExpenseStats, error) {
	return s.expenseModule.Stats(req)
}

func (s *StoreExpenseService) buildRecord(storeID, operatorID uint, hqUnbound bool, reqStoreID uint, categoryCode string, amount float64, remark string) (*model.StoreExpense, error) {
	realStoreID := storeID
	if hqUnbound && reqStoreID > 0 {
		realStoreID = reqStoreID
	}
	if realStoreID == 0 {
		return nil, fmt.Errorf("请选择门店")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("支出金额必须大于0")
	}
	categoryCode, categoryName, err := s.resolveCategory(categoryCode)
	if err != nil {
		return nil, err
	}
	operatorName := ""
	if s.userModule != nil {
		if user, err := s.userModule.GetByID(operatorID); err == nil && user != nil {
			operatorName = user.Nickname
			if operatorName == "" {
				operatorName = user.Username
			}
		}
	}
	return &model.StoreExpense{
		StoreID:      realStoreID,
		ExpenseDate:  businessdate.Date(time.Now()),
		CategoryCode: categoryCode,
		CategoryName: categoryName,
		Amount:       amount,
		Remark:       strings.TrimSpace(remark),
		OperatorID:   operatorID,
		OperatorName: operatorName,
	}, nil
}

func (s *StoreExpenseService) resolveCategory(code string) (string, string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", "", fmt.Errorf("请选择支出分类")
	}
	if s.dictModule == nil {
		return code, code, nil
	}
	data, err := s.dictModule.GetDataByTypeAndValue(model.StoreExpenseCategoryDictCode, code)
	if err != nil || data == nil || data.Status != 1 {
		return "", "", fmt.Errorf("支出分类不存在或已停用")
	}
	return data.Value, data.Label, nil
}
