package module

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

type PrinterModule struct {
	db *gorm.DB
}

func NewPrinterModule(db *gorm.DB) *PrinterModule {
	return &PrinterModule{db: db}
}

// GetDB 返回底层数据库实例
func (m *PrinterModule) GetDB() *gorm.DB {
	return m.db
}

// Create 创建打印机
func (m *PrinterModule) Create(printer *model.Printer) error {
	return m.db.Create(printer).Error
}

// GetByID 根据ID获取打印机
func (m *PrinterModule) GetByID(id uint) (*model.Printer, error) {
	var printer model.Printer
	if err := m.db.First(&printer, id).Error; err != nil {
		return nil, err
	}
	return &printer, nil
}

// GetBySn 根据SN获取打印机
func (m *PrinterModule) GetBySn(sn string) (*model.Printer, error) {
	var printer model.Printer
	if err := m.db.Where("sn = ?", sn).First(&printer).Error; err != nil {
		return nil, err
	}
	return &printer, nil
}

// ListByStoreID 获取门店下的所有打印机
func (m *PrinterModule) ListByStoreID(storeID uint) ([]*model.Printer, error) {
	var printers []*model.Printer
	if err := m.db.Where("store_id = ?", storeID).Find(&printers).Error; err != nil {
		return nil, err
	}
	return printers, nil
}

// ListAll 获取所有打印机
func (m *PrinterModule) ListAll() ([]*model.Printer, int64, error) {
	var printers []*model.Printer
	var total int64

	if err := m.db.Model(&model.Printer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := m.db.Find(&printers).Error; err != nil {
		return nil, 0, err
	}

	return printers, total, nil
}

// GetDefaultByStoreID 获取门店默认打印机
func (m *PrinterModule) GetDefaultByStoreID(storeID uint) (*model.Printer, error) {
	var printer model.Printer
	if err := m.db.Where("store_id = ? AND is_default = 1", storeID).First(&printer).Error; err != nil {
		return nil, err
	}
	return &printer, nil
}

// Update 更新打印机信息
func (m *PrinterModule) Update(printer *model.Printer) error {
	return m.db.Save(printer).Error
}

// UpdateByID 根据ID更新打印机信息
func (m *PrinterModule) UpdateByID(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.Printer{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除打印机
func (m *PrinterModule) Delete(id uint) error {
	return m.db.Delete(&model.Printer{}, id).Error
}

// ClearDefault 重置门店下所有打印机为非默认
func (m *PrinterModule) ClearDefault(storeID uint) error {
	return m.db.Model(&model.Printer{}).Where("store_id = ?", storeID).Update("is_default", 0).Error
}

// BindStore 绑定打印机到门店（同时推送到芯烨云）
func (m *PrinterModule) BindStore(printer *model.Printer) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 如果设置为默认，先清除其他默认
		if printer.IsDefault == 1 {
			if err := tx.Model(&model.Printer{}).Where("store_id = ?", printer.StoreID).Update("is_default", 0).Error; err != nil {
				return err
			}
		}
		return tx.Create(printer).Error
	})
}