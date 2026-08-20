package module

import (
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
)

func TestGroupMembersWithUnsettledAccounts(t *testing.T) {
	memberOneID := uint(1)
	memberTwoID := uint(2)
	unknownMemberID := uint(3)
	members := []model.Member{
		{ID: memberOneID, UnsettledAmount: 999},
		{ID: memberTwoID, UnsettledAmount: 999},
	}
	accounts := []*model.StoreAccount{
		{ID: 11, MemberID: &memberOneID, TotalAmount: 12.5},
		nil,
		{ID: 12, MemberID: &memberOneID, TotalAmount: 7.5},
		{ID: 21, MemberID: &memberTwoID, TotalAmount: 30},
		{ID: 31, MemberID: &unknownMemberID, TotalAmount: 40},
		{ID: 32, TotalAmount: 50},
	}

	got := groupMembersWithUnsettledAccounts(members, accounts)
	if len(got) != 2 {
		t.Fatalf("group count = %d, want 2", len(got))
	}
	if len(got[0].UnsettledAccounts) != 2 || got[0].UnsettledAmount != 20 {
		t.Fatalf("member 1 result = %#v, want 2 accounts totaling 20", got[0])
	}
	if len(got[1].UnsettledAccounts) != 1 || got[1].UnsettledAmount != 30 {
		t.Fatalf("member 2 result = %#v, want 1 account totaling 30", got[1])
	}
}

func TestGroupMembersWithUnsettledAccountsReturnsEmptyAccountList(t *testing.T) {
	got := groupMembersWithUnsettledAccounts([]model.Member{{ID: 1}}, nil)
	if len(got) != 1 {
		t.Fatalf("group count = %d, want 1", len(got))
	}
	if got[0].UnsettledAccounts == nil || len(got[0].UnsettledAccounts) != 0 {
		t.Fatalf("unsettled accounts = %#v, want non-nil empty list", got[0].UnsettledAccounts)
	}
	if got[0].UnsettledAmount != 0 {
		t.Fatalf("unsettled amount = %v, want 0", got[0].UnsettledAmount)
	}
}
