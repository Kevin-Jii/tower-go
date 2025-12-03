package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type DictModule struct {
	db *gorm.DB
}

func NewDictModule(db *gorm.DB) *DictModule {
	return &DictModule{db: db}
}

// ========== 字典类型 ==========

// CreateType 创建字典类型
func (m *DictModule) CreateType(dictType *model.DictType) error {
	return m.db.Create(dictType).Error
}

// GetTypeByID 根据ID获取字典类型
func (m *DictModule) GetTypeByID(id uint) (*model.DictType, error) {
	var dictType model.DictType
	if err := m.db.First(&dictType, id).Error; err != nil {
		return nil, err
	}
	return &dictType, nil
}

// GetTypeByCode 根据编码获取字典类型
func (m *DictModule) GetTypeByCode(code string) (*model.DictType, error) {
	var dictType model.DictType
	if err := m.db.Where("code = ?", code).First(&dictType).Error; err != nil {
		return nil, err
	}
	return &dictType, nil
}

// ListTypes 获取字典类型列表
func (m *DictModule) ListTypes(keyword string, status *int8) ([]*model.DictType, error) {
	var types []*model.DictType
	query := m.db.Model(&model.DictType{})

	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("id ASC").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

// UpdateType 更新字典类型
func (m *DictModule) UpdateType(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.DictType{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteType 删除字典类型
func (m *DictModule) DeleteType(id uint) error {
	// 先删除关联的字典数据
	m.db.Where("type_id = ?", id).Delete(&model.DictData{})
	return m.db.Delete(&model.DictType{}, id).Error
}

// ========== 字典数据 ==========

// CreateData 创建字典数据
func (m *DictModule) CreateData(dictData *model.DictData) error {
	// 获取类型ID
	dictType, err := m.GetTypeByCode(dictData.TypeCode)
	if err != nil {
		return err
	}
	dictData.TypeID = dictType.ID
	return m.db.Create(dictData).Error
}

// GetDataByID 根据ID获取字典数据
func (m *DictModule) GetDataByID(id uint) (*model.DictData, error) {
	var dictData model.DictData
	if err := m.db.First(&dictData, id).Error; err != nil {
		return nil, err
	}
	return &dictData, nil
}

// ListDataByTypeCode 根据类型编码获取字典数据列表
func (m *DictModule) ListDataByTypeCode(typeCode string) ([]*model.DictData, error) {
	var dataList []*model.DictData
	if err := m.db.Where("type_code = ? AND status = 1", typeCode).Order("sort ASC, id ASC").Find(&dataList).Error; err != nil {
		return nil, err
	}
	return dataList, nil
}

// ListDataByTypeID 根据类型ID获取字典数据列表
func (m *DictModule) ListDataByTypeID(typeID uint, status *int8) ([]*model.DictData, error) {
	var dataList []*model.DictData
	query := m.db.Where("type_id = ?", typeID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Order("sort ASC, id ASC").Find(&dataList).Error; err != nil {
		return nil, err
	}
	return dataList, nil
}

// UpdateData 更新字典数据
func (m *DictModule) UpdateData(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.DictData{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteData 删除字典数据
func (m *DictModule) DeleteData(id uint) error {
	return m.db.Delete(&model.DictData{}, id).Error
}

// GetDataByTypeAndValue 根据类型编码和值获取字典数据
func (m *DictModule) GetDataByTypeAndValue(typeCode, value string) (*model.DictData, error) {
	var dictData model.DictData
	if err := m.db.Where("type_code = ? AND value = ?", typeCode, value).First(&dictData).Error; err != nil {
		return nil, err
	}
	return &dictData, nil
}
