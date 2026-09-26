package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type MemberTagService struct{ module *module.MemberTagModule }

func NewMemberTagService(m *module.MemberTagModule) *MemberTagService {
	return &MemberTagService{module: m}
}
func (s *MemberTagService) List(storeID uint, allStores bool) ([]model.MemberTag, error) {
	return s.module.List(storeID, allStores)
}
func (s *MemberTagService) Get(id, storeID uint, allStores bool) (*model.MemberTag, error) {
	return s.module.Get(id, storeID, allStores)
}
func (s *MemberTagService) Create(req *model.UpsertMemberTagReq, storeID uint) (*model.MemberTag, error) {
	return s.module.Create(req, storeID)
}
func (s *MemberTagService) Update(id, storeID uint, allStores bool, req *model.UpsertMemberTagReq) (*model.MemberTag, error) {
	return s.module.Update(id, storeID, allStores, req)
}
func (s *MemberTagService) Delete(id, storeID uint, allStores bool) error {
	return s.module.Delete(id, storeID, allStores)
}
func (s *MemberTagService) ListTagMembers(tagID, storeID uint, allStores bool) ([]model.SmsMemberSummary, error) {
	return s.module.ListTagMembers(tagID, storeID, allStores)
}
func (s *MemberTagService) SearchMembers(keyword string, storeID uint, allStores bool, limit int) ([]model.SmsMemberSummary, error) {
	return s.module.SearchMembers(keyword, storeID, allStores, limit)
}
func (s *MemberTagService) BindMember(tagID, memberID, storeID uint, allStores bool) (*model.MemberTagBinding, error) {
	return s.module.BindMember(tagID, memberID, storeID, allStores)
}
func (s *MemberTagService) UnbindMember(tagID, memberID, storeID uint, allStores bool) error {
	return s.module.UnbindMember(tagID, memberID, storeID, allStores)
}
func (s *MemberTagService) ListMemberTags(memberID, storeID uint, allStores bool) ([]model.MemberTag, error) {
	return s.module.ListMemberTags(memberID, storeID, allStores)
}
func (s *MemberTagService) AssignMemberTags(memberID, storeID uint, allStores bool, tagIDs []uint) error {
	return s.module.AssignMemberTags(memberID, storeID, allStores, tagIDs)
}
