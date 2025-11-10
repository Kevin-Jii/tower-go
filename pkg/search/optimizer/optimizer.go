package optimizer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// SearchType 搜索类型
type SearchType int

const (
	SearchTypeExact SearchType = iota // 精确匹配
	SearchTypePrefix                  // 前缀匹配
	SearchTypeSuffix                  // 后缀匹配
	SearchTypeFull                    // 全文匹配
	SearchTypeRegex                   // 正则匹配
)

// String 返回搜索类型的字符串表示
func (st SearchType) String() string {
	switch st {
	case SearchTypeExact:
		return "exact"
	case SearchTypePrefix:
		return "prefix"
	case SearchTypeSuffix:
		return "suffix"
	case SearchTypeFull:
		return "full"
	case SearchTypeRegex:
		return "regex"
	default:
		return "unknown"
	}
}

// Condition 搜索条件
type Condition struct {
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
	SearchType SearchType  `json:"search_type"`
	Weight     float64     `json:"weight"`     // 权重，用于相关性排序
	Boost      float64     `json:"boost"`      // 提升因子
	Negated    bool        `json:"negated"`    // 是否取反
}

// Config 搜索优化配置
type Config struct {
	EnableFuzzy         bool    `yaml:"enable_fuzzy" json:"enable_fuzzy"`
	FuzzyThreshold      float64 `yaml:"fuzzy_threshold" json:"fuzzy_threshold"`
	MinKeywordLength    int     `yaml:"min_keyword_length" json:"min_keyword_length"`
	MaxKeywordLength    int     `yaml:"max_keyword_length" json:"max_keyword_length"`
	EnablePhoneOptimize bool    `yaml:"enable_phone_optimize" json:"enable_phone_optimize"`
	EnableEmailOptimize bool    `yaml:"enable_email_optimize" json:"enable_email_optimize"`
	EnableChineseOpt    bool    `yaml:"enable_chinese_opt" json:"enable_chinese_opt"`
	DefaultWeight       float64 `yaml:"default_weight" json:"default_weight"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		EnableFuzzy:         true,
		FuzzyThreshold:      0.6,
		MinKeywordLength:    1,
		MaxKeywordLength:    100,
		EnablePhoneOptimize: true,
		EnableEmailOptimize: true,
		EnableChineseOpt:    true,
		DefaultWeight:       1.0,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.FuzzyThreshold < 0 || c.FuzzyThreshold > 1 {
		return fmt.Errorf("fuzzy_threshold must be between 0 and 1")
	}
	if c.MinKeywordLength < 0 {
		return fmt.Errorf("min_keyword_length cannot be negative")
	}
	if c.MaxKeywordLength < c.MinKeywordLength {
		return fmt.Errorf("max_keyword_length must be greater than min_keyword_length")
	}
	if c.DefaultWeight < 0 {
		return fmt.Errorf("default_weight cannot be negative")
	}
	return nil
}

// Optimizer 搜索优化器
type Optimizer struct {
	config *Config
}

// New 创建搜索优化器
func New(config *Config) (*Optimizer, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &Optimizer{config: config}, nil
}

// NewWithDefaults 使用默认配置创建优化器
func NewWithDefaults() *Optimizer {
	opt, _ := New(DefaultConfig())
	return opt
}

// OptimizeKeyword 优化搜索关键字
func (o *Optimizer) OptimizeKeyword(keyword string, fields []string) []Condition {
	if keyword == "" {
		return []Condition{}
	}

	keyword = strings.TrimSpace(keyword)
	if len(keyword) < o.config.MinKeywordLength || len(keyword) > o.config.MaxKeywordLength {
		return []Condition{}
	}

	var conditions []Condition

	// 分析关键字特征
	keywordType := o.analyzeKeyword(keyword)

	switch keywordType {
	case "phone":
		conditions = append(conditions, o.optimizePhoneSearch(keyword, fields...)...)
	case "email":
		conditions = append(conditions, o.optimizeEmailSearch(keyword, fields...)...)
	case "id":
		conditions = append(conditions, o.optimizeIDSearch(keyword, fields...)...)
	case "chinese":
		conditions = append(conditions, o.optimizeChineseSearch(keyword, fields...)...)
	case "english":
		conditions = append(conditions, o.optimizeEnglishSearch(keyword, fields...)...)
	case "numeric":
		conditions = append(conditions, o.optimizeNumericSearch(keyword, fields...)...)
	default:
		conditions = append(conditions, o.optimizeGenericSearch(keyword, fields...)...)
	}

	// 启用模糊搜索时添加模糊匹配条件
	if o.config.EnableFuzzy && len(keyword) >= 3 {
		conditions = append(conditions, o.generateFuzzyConditions(keyword, fields...)...)
	}

	return conditions
}

// analyzeKeyword 分析关键字类型
func (o *Optimizer) analyzeKeyword(keyword string) string {
	// 手机号特征：11位数字，或4-6位数字（后缀搜索）
	if o.isPhone(keyword) {
		return "phone"
	}

	// 邮箱特征：包含@符号
	if strings.Contains(keyword, "@") {
		return "email"
	}

	// ID特征：短数字
	if o.isID(keyword) {
		return "id"
	}

	// 中文特征
	if o.containsChinese(keyword) {
		return "chinese"
	}

	// 英文特征
	if o.isEnglish(keyword) {
		return "english"
	}

	// 数字特征
	if o.isNumeric(keyword) {
		return "numeric"
	}

	return "unknown"
}

// optimizePhoneSearch 优化手机号搜索
func (o *Optimizer) optimizePhoneSearch(keyword string, fields ...string) []Condition {
	if !o.config.EnablePhoneOptimize {
		return []Condition{}
	}

	var conditions []Condition

	// 如果是完整手机号（11位），优先精确匹配
	if len(keyword) == 11 && o.isNumeric(keyword) {
		conditions = append(conditions, Condition{
			Field:      "phone",
			Value:      keyword,
			SearchType: SearchTypeExact,
			Weight:     2.0, // 高权重
			Boost:      1.5,
		})
	}

	// 后缀匹配（支持输入手机号后4位或后6位）
	if len(keyword) >= 4 && len(keyword) <= 6 && o.isNumeric(keyword) {
		conditions = append(conditions, Condition{
			Field:      "phone",
			Value:      "%" + keyword,
			SearchType: SearchTypeSuffix,
			Weight:     1.5, // 中等权重
		})
	}

	// 模糊匹配
	if len(keyword) >= 3 {
		conditions = append(conditions, Condition{
			Field:      "phone",
			Value:      "%" + keyword + "%",
			SearchType: SearchTypeFull,
			Weight:     0.8, // 低权重
		})
	}

	return conditions
}

// optimizeEmailSearch 优化邮箱搜索
func (o *Optimizer) optimizeEmailSearch(keyword string, fields ...string) []Condition {
	if !o.config.EnableEmailOptimize {
		return []Condition{}
	}

	return []Condition{
		{
			Field:      "email",
			Value:      keyword + "%",
			SearchType: SearchTypePrefix,
			Weight:     1.2,
		},
	}
}

// optimizeIDSearch 优化ID搜索
func (o *Optimizer) optimizeIDSearch(keyword string, fields ...string) []Condition {
	if id, err := strconv.ParseUint(keyword, 10, 32); err == nil {
		return []Condition{
			{
				Field:      "id",
				Value:      id,
				SearchType: SearchTypeExact,
				Weight:     3.0, // 最高权重
				Boost:      2.0,
			},
		}
	}
	return []Condition{}
}

// optimizeChineseSearch 优化中文搜索
func (o *Optimizer) optimizeChineseSearch(keyword string, fields ...string) []Condition {
	if !o.config.EnableChineseOpt {
		return []Condition{}
	}

	var conditions []Condition

	// 前缀匹配对于中文更有效
	conditions = append(conditions, Condition{
		Field:      "username",
		Value:      keyword + "%",
		SearchType: SearchTypePrefix,
		Weight:     1.3,
	})

	// 完整匹配（权重高）
	if len(keyword) <= 10 {
		conditions = append(conditions, Condition{
			Field:      "username",
			Value:      keyword,
			SearchType: SearchTypeExact,
			Weight:     2.0,
			Boost:      1.2,
		})
	}

	// 包含匹配
	conditions = append(conditions, Condition{
		Field:      "username",
		Value:      "%" + keyword + "%",
		SearchType: SearchTypeFull,
		Weight:     1.0,
	})

	return conditions
}

// optimizeEnglishSearch 优化英文搜索
func (o *Optimizer) optimizeEnglishSearch(keyword string, fields ...string) []Condition {
	var conditions []Condition

	// 英文单词优先前缀匹配
	conditions = append(conditions, Condition{
		Field:      "username",
		Value:      keyword + "%",
		SearchType: SearchTypePrefix,
		Weight:     1.4,
	})

	// 小写匹配
	if strings.ToLower(keyword) != keyword {
		conditions = append(conditions, Condition{
			Field:      "username",
			Value:      strings.ToLower(keyword) + "%",
			SearchType: SearchTypePrefix,
			Weight:     1.2,
		})
	}

	return conditions
}

// optimizeNumericSearch 优化数字搜索
func (o *Optimizer) optimizeNumericSearch(keyword string, fields ...string) []Condition {
	var conditions []Condition

	// 精确匹配
	if num, err := strconv.ParseInt(keyword, 10, 64); err == nil {
		// 尝试匹配可能的数字字段
		numericFields := []string{"id", "phone", "age", "score", "level"}
		for _, field := range numericFields {
			conditions = append(conditions, Condition{
				Field:      field,
				Value:      num,
				SearchType: SearchTypeExact,
				Weight:     1.8,
			})
		}
	}

	return conditions
}

// optimizeGenericSearch 通用搜索
func (o *Optimizer) optimizeGenericSearch(keyword string, fields ...string) []Condition {
	var conditions []Condition

	// 默认对用户名进行全文搜索
	conditions = append(conditions, Condition{
		Field:      "username",
		Value:      "%" + keyword + "%",
		SearchType: SearchTypeFull,
		Weight:     o.config.DefaultWeight,
	})

	// 如果提供了其他字段，也进行搜索
	for _, field := range fields {
		if field != "username" {
			conditions = append(conditions, Condition{
				Field:      field,
				Value:      "%" + keyword + "%",
				SearchType: SearchTypeFull,
				Weight:     o.config.DefaultWeight * 0.8,
			})
		}
	}

	return conditions
}

// generateFuzzyConditions 生成模糊搜索条件
func (o *Optimizer) generateFuzzyConditions(keyword string, fields ...string) []Condition {
	var conditions []Condition

	// 生成常见拼写错误的变体
	variations := o.generateVariations(keyword)
	for _, variation := range variations {
		if variation != keyword {
			conditions = append(conditions, Condition{
				Field:      "username",
				Value:      "%" + variation + "%",
				SearchType: SearchTypeFull,
				Weight:     o.config.DefaultWeight * 0.6,
				Boost:      0.8,
			})
		}
	}

	return conditions
}

// generateVariations 生成关键字变体
func (o *Optimizer) generateVariations(keyword string) []string {
	var variations []string

	// 转换大小写
	variations = append(variations, strings.ToLower(keyword))
	variations = append(variations, strings.ToUpper(keyword))

	// 移除常见错误字符
	cleaned := strings.ReplaceAll(keyword, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "_", "")
	if cleaned != keyword {
		variations = append(variations, cleaned)
	}

	return variations
}

// MultiTermOptimize 多词组搜索优化
func (o *Optimizer) MultiTermOptimize(keyword string, fields []string) []Condition {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []Condition{}
	}

	// 分割关键词
	terms := strings.Fields(keyword)
	if len(terms) == 1 {
		return o.OptimizeKeyword(terms[0], fields)
	}

	var conditions []Condition
	for _, term := range terms {
		termConditions := o.OptimizeKeyword(term, fields)
		conditions = append(conditions, termConditions...)
	}

	return conditions
}

// GetConfig 获取配置
func (o *Optimizer) GetConfig() *Config {
	return o.config
}

// SetConfig 更新配置
func (o *Optimizer) SetConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	o.config = config
	return nil
}

// 工具方法

func (o *Optimizer) isPhone(keyword string) bool {
	// 11位手机号
	if len(keyword) == 11 && o.isNumeric(keyword) {
		return true
	}
	// 4-6位数字（手机号后缀）
	if len(keyword) >= 4 && len(keyword) <= 6 && o.isNumeric(keyword) {
		return true
	}
	return false
}

func (o *Optimizer) isEmail(keyword string) bool {
	return strings.Contains(keyword, "@") && strings.Contains(keyword, ".")
}

func (o *Optimizer) isID(keyword string) bool {
	if len(keyword) > 10 {
		return false
	}
	return o.isNumeric(keyword)
}

func (o *Optimizer) containsChinese(keyword string) bool {
	for _, r := range keyword {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func (o *Optimizer) isEnglish(keyword string) bool {
	hasLetter := false
	for _, r := range keyword {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Han, r) {
			hasLetter = true
		}
		if unicode.Is(unicode.Han, r) || unicode.IsNumber(r) {
			return false
		}
	}
	return hasLetter
}

func (o *Optimizer) isNumeric(keyword string) bool {
	matched, _ := regexp.MatchString(`^\d+$`, keyword)
	return matched
}

// 全局默认优化器
var defaultOptimizer = NewWithDefaults()

// 便利函数（使用默认优化器）

// OptimizeKeyword 优化搜索关键字
func OptimizeKeyword(keyword string, fields []string) []Condition {
	return defaultOptimizer.OptimizeKeyword(keyword, fields)
}

// MultiTermOptimize 多词组搜索优化
func MultiTermOptimize(keyword string, fields []string) []Condition {
	return defaultOptimizer.MultiTermOptimize(keyword, fields)
}