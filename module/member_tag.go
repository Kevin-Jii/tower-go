package module

import (
	"errors"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemberTagModule struct{ db *gorm.DB }

func NewMemberTagModule(db *gorm.DB) *MemberTagModule { return &MemberTagModule{db: db} }

func (m *MemberTagModule) List(storeID uint, allStores bool) ([]model.MemberTag, error) {
	rows := make([]model.MemberTag, 0)
	q := m.db.Model(&model.MemberTag{})
	if !allStores || storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	if err := q.Order("store_id ASC, id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return rows, nil
	}
	ids := make([]uint, len(rows))
	for i := range rows {
		ids[i] = rows[i].ID
	}
	type countRow struct {
		TagID uint
		Count int64
	}
	var counts []countRow
	if err := m.db.Model(&model.MemberTagBinding{}).Select("tag_id, COUNT(*) AS count").Where("tag_id IN ?", ids).Group("tag_id").Scan(&counts).Error; err != nil {
		return nil, err
	}
	byID := make(map[uint]int64, len(counts))
	for _, c := range counts {
		byID[c.TagID] = c.Count
	}
	for i := range rows {
		rows[i].MemberCount = byID[rows[i].ID]
	}
	return rows, nil
}

func (m *MemberTagModule) Get(id, storeID uint, allStores bool) (*model.MemberTag, error) {
	var row model.MemberTag
	q := m.db.Model(&model.MemberTag{}).Where("id = ?", id)
	if !allStores || storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	if err := q.First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *MemberTagModule) Create(req *model.UpsertMemberTagReq, storeID uint) (*model.MemberTag, error) {
	if storeID == 0 {
		return nil, errors.New("请选择标签所属门店")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("请填写标签名称")
	}
	row := &model.MemberTag{StoreID: storeID, Name: name, Color: strings.TrimSpace(req.Color), Description: strings.TrimSpace(req.Description)}
	if err := m.db.Create(row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

func (m *MemberTagModule) Update(id, storeID uint, allStores bool, req *model.UpsertMemberTagReq) (*model.MemberTag, error) {
	row, err := m.Get(id, storeID, allStores)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("请填写标签名称")
	}
	if err := m.db.Model(row).Updates(map[string]interface{}{"name": name, "color": strings.TrimSpace(req.Color), "description": strings.TrimSpace(req.Description)}).Error; err != nil {
		return nil, err
	}
	return m.Get(id, row.StoreID, false)
}

func (m *MemberTagModule) Delete(id, storeID uint, allStores bool) error {
	row, err := m.Get(id, storeID, allStores)
	if err != nil {
		return err
	}
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tag_id = ?", row.ID).Delete(&model.MemberTagBinding{}).Error; err != nil {
			return err
		}
		return tx.Delete(row).Error
	})
}

func (m *MemberTagModule) ValidateTagIDs(tagIDs []uint, storeID uint) error {
	return m.ValidateTagIDsAccess(tagIDs, storeID, false)
}

func (m *MemberTagModule) ValidateTagIDsAccess(tagIDs []uint, storeID uint, allStores bool) error {
	if len(tagIDs) == 0 {
		return nil
	}
	unique := uniqueUint(tagIDs)
	q := m.db.Model(&model.MemberTag{}).Where("id IN ?", unique)
	if !allStores || storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(unique)) {
		return errors.New("标签不存在或不属于活动门店")
	}
	return nil
}

func (m *MemberTagModule) ListTagMembers(tagID, storeID uint, allStores bool) ([]model.SmsMemberSummary, error) {
	tag, err := m.Get(tagID, storeID, allStores)
	if err != nil {
		return nil, err
	}
	rows := make([]model.SmsMemberSummary, 0)
	err = m.db.Table("t_member").
		Select("t_member.id, t_member.store_id, t_member.uid, t_member.name, t_member.phone").
		Joins("JOIN member_tag_bindings b ON b.member_id = t_member.id").
		Where("b.tag_id = ? AND b.store_id = ?", tag.ID, tag.StoreID).
		Order("t_member.id ASC").Scan(&rows).Error
	return rows, err
}

func (m *MemberTagModule) SearchMembers(keyword string, storeID uint, allStores bool, limit int) ([]model.SmsMemberSummary, error) {
	rows := make([]model.SmsMemberSummary, 0)
	q := m.db.Table("t_member").Select("id, store_id, uid, name, phone")
	if !allStores || storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("phone LIKE ? OR uid LIKE ? OR name LIKE ?", like, like, like)
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if err := q.Order("id DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *MemberTagModule) BindMember(tagID, memberID, storeID uint, allStores bool) (*model.MemberTagBinding, error) {
	tag, err := m.Get(tagID, storeID, allStores)
	if err != nil {
		return nil, err
	}
	var member model.Member
	if err := m.db.Where("id = ? AND store_id = ?", memberID, tag.StoreID).First(&member).Error; err != nil {
		return nil, errors.New("会员不存在或不属于标签门店")
	}
	row := &model.MemberTagBinding{StoreID: tag.StoreID, MemberID: member.ID, TagID: tag.ID}
	if err := m.db.Clauses(clause.OnConflict{DoNothing: true}).Create(row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		if err := m.db.Where("store_id = ? AND member_id = ? AND tag_id = ?", tag.StoreID, member.ID, tag.ID).First(row).Error; err != nil {
			return nil, err
		}
	}
	return row, nil
}

func (m *MemberTagModule) UnbindMember(tagID, memberID, storeID uint, allStores bool) error {
	tag, err := m.Get(tagID, storeID, allStores)
	if err != nil {
		return err
	}
	result := m.db.Where("store_id = ? AND member_id = ? AND tag_id = ?", tag.StoreID, memberID, tag.ID).Delete(&model.MemberTagBinding{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *MemberTagModule) ListMemberTags(memberID, storeID uint, allStores bool) ([]model.MemberTag, error) {
	var member model.Member
	q := m.db.Model(&model.Member{}).Where("id = ?", memberID)
	if !allStores || storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	if err := q.First(&member).Error; err != nil {
		return nil, err
	}
	rows := make([]model.MemberTag, 0)
	err := m.db.Model(&model.MemberTag{}).Joins("JOIN member_tag_bindings b ON b.tag_id = member_tags.id").Where("b.member_id = ? AND b.store_id = ?", member.ID, member.StoreID).Order("member_tags.id DESC").Find(&rows).Error
	return rows, err
}

func (m *MemberTagModule) AssignMemberTags(memberID, storeID uint, allStores bool, tagIDs []uint) error {
	var member model.Member
	q := m.db.Model(&model.Member{}).Where("id = ?", memberID)
	if !allStores || storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	if err := q.First(&member).Error; err != nil {
		return err
	}
	ids := uniqueUint(tagIDs)
	if err := m.ValidateTagIDs(ids, member.StoreID); err != nil {
		return err
	}
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("member_id = ? AND store_id = ?", member.ID, member.StoreID).Delete(&model.MemberTagBinding{}).Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}
		rows := make([]model.MemberTagBinding, 0, len(ids))
		for _, id := range ids {
			rows = append(rows, model.MemberTagBinding{StoreID: member.StoreID, MemberID: member.ID, TagID: id})
		}
		return tx.Create(&rows).Error
	})
}

func (m *MemberTagModule) MemberTagMap(memberIDs []uint, storeID uint) (map[uint]map[uint]struct{}, error) {
	out := make(map[uint]map[uint]struct{}, len(memberIDs))
	if len(memberIDs) == 0 {
		return out, nil
	}
	var rows []model.MemberTagBinding
	q := m.db.Where("member_id IN ?", memberIDs)
	if storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		if out[row.MemberID] == nil {
			out[row.MemberID] = map[uint]struct{}{}
		}
		out[row.MemberID][row.TagID] = struct{}{}
	}
	return out, nil
}

func uniqueUint(in []uint) []uint {
	out := make([]uint, 0, len(in))
	seen := map[uint]struct{}{}
	for _, id := range in {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
