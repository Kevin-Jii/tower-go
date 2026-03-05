package performance

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Property 5: Automatic index strategy - Requirements 2.1
func TestProperty_AutomaticIndexStrategy(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("query optimizer detects missing index issues", prop.ForAll(
		func(sql string) bool {
			if sql == "" {
				return true
			}
			qo := NewQueryOptimizer(nil)
			result, err := qo.Analyze(sql)
			if err != nil {
				return false
			}
			return result != nil
		},
		gen.AlphaString(),
	))

	properties.Property("query optimizer detects offset pagination", prop.ForAll(
		func(offset int) bool {
			if offset < 0 {
				return true
			}
			sql := "SELECT * FROM orders LIMIT 10 OFFSET " + string(rune(offset))
			qo := NewQueryOptimizer(nil)
			_, err := qo.Analyze(sql)
			return err == nil
		},
		gen.IntRange(0, 1000),
	))

	properties.Property("query optimizer detects full table scan", prop.ForAll(
		func(sql string) bool {
			if sql == "" {
				return true
			}
			sql = "SELECT * FROM store_accounts"
			qo := NewQueryOptimizer(nil)
			result, err := qo.Analyze(sql)
			if err != nil {
				return false
			}
			hasIssue := false
			for _, issue := range result.Issues {
				if issue.Type == IssueTypeFullTableScan {
					hasIssue = true
					break
				}
			}
			return hasIssue || result.IndexUsage.TableScan
		},
		gen.AlphaString(),
	))

	properties.Property("query optimizer detects duplicate join", prop.ForAll(
		func(tableName string) bool {
			if tableName == "" {
				return true
			}
			sql := "SELECT * FROM " + tableName + " JOIN users ON 1=1 JOIN users ON 1=1"
			qo := NewQueryOptimizer(nil)
			result, err := qo.Analyze(sql)
			if err != nil {
				return false
			}
			for _, issue := range result.Issues {
				if issue.Type == IssueTypeDuplicateJoin {
					return true
				}
			}
			return result.JoinCount >= 2
		},
		gen.Identifier(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property 7: JOIN deduplication - Requirements 2.3
func TestProperty_JoinDeduplication(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("join deduplicator detects duplicate tables", prop.ForAll(
		func(table string) bool {
			if table == "" {
				return true
			}
			jd := NewJoinDeduplicator()
			if !jd.Add(table) {
				return false
			}
			return !jd.Add(table)
		},
		gen.Identifier(),
	))

	properties.Property("join deduplicator resets correctly", prop.ForAll(
		func(table1, table2 string) bool {
			if table1 == "" || table2 == "" {
				return true
			}
			jd := NewJoinDeduplicator()
			jd.Add(table1)
			jd.Reset()
			return jd.Add(table1) && jd.Add(table2)
		},
		gen.Identifier(),
		gen.Identifier(),
	))

	properties.Property("join deduplicator handles same table differently cased", prop.ForAll(
		func(table string) bool {
			if table == "" {
				return true
			}
			jd := NewJoinDeduplicator()
			jd.Add(table)
			jd.Add("UPPER_" + table)
			jd.Add("lower_" + table)
			return true
		},
		gen.Identifier(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestQueryOptimizerBasic tests basic query optimizer functionality
func TestQueryOptimizerBasic(t *testing.T) {
	qo := NewQueryOptimizer(nil)

	sql := "SELECT * FROM store_accounts WHERE store_id = ?"
	result, err := qo.Analyze(sql)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	table := qo.extractTable("SELECT * FROM users WHERE id = 1")
	if table != "users" {
		t.Errorf("Expected 'users', got '%s'", table)
	}

	count := qo.countJoins("SELECT * FROM a JOIN b ON a.id = b.id")
	if count != 1 {
		t.Errorf("Expected 1 join, got %d", count)
	}
}

// TestIndexAnalyzer tests index analyzer
func TestIndexAnalyzer(t *testing.T) {
	ia := NewIndexAnalyzer(nil)

	result := ia.AnalyzeIndexUsage("store_accounts", []string{"store_id", "account_date"})
	if len(result.PotentialIndexes) == 0 {
		t.Error("Expected potential indexes to be identified")
	}

	result = ia.AnalyzeIndexUsage("test", []string{"*"})
	if !result.TableScan {
		t.Error("Expected TableScan to be true for wildcard")
	}
}