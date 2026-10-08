package module

import (
	"strings"
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestConversationOwnerScopeIncludesUserAndStore(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root@tcp(localhost:3306)/tower_test?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}

	var rows []model.AIAssistantConversation
	statement := conversationOwnerScope(db, 17, 23).
		Where("id = ?", 31).
		Find(&rows).
		Statement

	sql := statement.SQL.String()
	for _, clause := range []string{"user_id = ?", "store_id = ?", "id = ?"} {
		if !strings.Contains(sql, clause) {
			t.Fatalf("query %q does not contain %q", sql, clause)
		}
	}
	wantVars := []any{uint(17), uint(23), 31}
	if len(statement.Vars) != len(wantVars) {
		t.Fatalf("vars = %#v, want %#v", statement.Vars, wantVars)
	}
	for i := range wantVars {
		if statement.Vars[i] != wantVars[i] {
			t.Fatalf("vars[%d] = %#v, want %#v", i, statement.Vars[i], wantVars[i])
		}
	}
}
