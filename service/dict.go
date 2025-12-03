package service

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type DictService struct {
	dictModule *module.DictModule
}

func NewDictService(dictModule *module.DictModule) *DictService {
	return &DictService{dictModule: dictModule}
}

// ========== 字典类型 ==========

// CreateType 创建字典类型
func (s *DictService) CreateType(req *model.CreateDictTypeReq) error {
	// 检查编码是否已存在
	if _, err := s.dictModule.GetTypeByCode(req.Code); err == nil {
		return errors.New("字典类型编码已存在")
	}

	dictType := &model.DictType{
		Code:   req.Code,
		Name:   req.Name,
		Remark: req.Remark,
		Status: req.Status,
	}
	if dictType.Status == 0 {
		dictType.Status = 1
	}
	return s.dictModule.CreateType(dictType)
}

// GetType 获取字典类型
func (s *DictService) GetType(id uint) (*model.DictType, error) {
	return s.dictModule.GetTypeByID(id)
}

// ListTypes 获取字典类型列表
func (s *DictService) ListTypes(keyword string, status *int8) ([]*model.DictType, error) {
	return s.dictModule.ListTypes(keyword, status)
}

// UpdateType 更新字典类型
func (s *DictService) UpdateType(id uint, req *model.UpdateDictTypeReq) error {
	updates := make(map[string]interface{})
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return nil
	}
	return s.dictModule.UpdateType(id, updates)
}

// DeleteType 删除字典类型
func (s *DictService) DeleteType(id uint) error {
	return s.dictModule.DeleteType(id)
}

// ========== 字典数据 ==========

// CreateData 创建字典数据
func (s *DictService) CreateData(req *model.CreateDictDataReq) error {
	dictData := &model.DictData{
		TypeCode:  req.TypeCode,
		Label:     req.Label,
		Value:     req.Value,
		Sort:      req.Sort,
		CssClass:  req.CssClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Remark:    req.Remark,
		Status:    req.Status,
	}
	if dictData.Status == 0 {
		dictData.Status = 1
	}
	return s.dictModule.CreateData(dictData)
}

// GetData 获取字典数据
func (s *DictService) GetData(id uint) (*model.DictData, error) {
	return s.dictModule.GetDataByID(id)
}

// ListDataByTypeCode 根据类型编码获取字典数据
func (s *DictService) ListDataByTypeCode(typeCode string) ([]*model.DictData, error) {
	return s.dictModule.ListDataByTypeCode(typeCode)
}

// ListDataByTypeID 根据类型ID获取字典数据
func (s *DictService) ListDataByTypeID(typeID uint, status *int8) ([]*model.DictData, error) {
	return s.dictModule.ListDataByTypeID(typeID, status)
}

// UpdateData 更新字典数据
func (s *DictService) UpdateData(id uint, req *model.UpdateDictDataReq) error {
	updates := make(map[string]interface{})
	if req.Label != "" {
		updates["label"] = req.Label
	}
	if req.Value != "" {
		updates["value"] = req.Value
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.CssClass != "" {
		updates["css_class"] = req.CssClass
	}
	if req.ListClass != "" {
		updates["list_class"] = req.ListClass
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return nil
	}
	return s.dictModule.UpdateData(id, updates)
}

// DeleteData 删除字典数据
func (s *DictService) DeleteData(id uint) error {
	return s.dictModule.DeleteData(id)
}

// GetAllDict 获取所有字典（用于前端缓存）
func (s *DictService) GetAllDict() (map[string][]*model.DictData, error) {
	types, err := s.dictModule.ListTypes("", nil)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*model.DictData)
	for _, t := range types {
		if t.Status != 1 {
			continue
		}
		dataList, err := s.dictModule.ListDataByTypeCode(t.Code)
		if err != nil {
			continue
		}
		result[t.Code] = dataList
	}
	return result, nil
}
