package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strings"
)

// Cost bcrypt加密成本
const Cost = bcrypt.DefaultCost

// Config 密码配置
type Config struct {
	MinLength        int  `yaml:"min_length" json:"min_length"`
	RequireUppercase bool `yaml:"require_uppercase" json:"require_uppercase"`
	RequireLowercase bool `yaml:"require_lowercase" json:"require_lowercase"`
	RequireNumbers   bool `yaml:"require_numbers" json:"require_numbers"`
	RequireSpecial   bool `yaml:"require_special" json:"require_special"`
}

// DefaultConfig 返回默认密码配置
func DefaultConfig() *Config {
	return &Config{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   true,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.MinLength < 6 {
		return errors.New("password minimum length must be at least 6")
	}
	if c.MinLength > 128 {
		return errors.New("password minimum length cannot exceed 128")
	}
	return nil
}

// Strength 密码强度等级
type Strength int

const (
	StrengthWeak   Strength = iota // 弱密码
	StrengthFair                  // 一般密码
	StrengthGood                  // 良好密码
	StrengthStrong                // 强密码
)

// String 返回强度字符串
func (s Strength) String() string {
	switch s {
	case StrengthWeak:
		return "Weak"
	case StrengthFair:
		return "Fair"
	case StrengthGood:
		return "Good"
	case StrengthStrong:
		return "Strong"
	default:
		return "Unknown"
	}
}

// Password 密码管理器
type Password struct {
	config *Config
}

// New 创建密码管理器
func New(config *Config) (*Password, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid password config: %w", err)
	}
	return &Password{config: config}, nil
}

// NewWithMinLength 使用最小长度创建密码管理器
func NewWithMinLength(minLength int) (*Password, error) {
	config := DefaultConfig()
	config.MinLength = minLength
	return New(config)
}

// Hash 加密密码
func (p *Password) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// Verify 验证密码
func (p *Password) Verify(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword 验证密码强度
func (p *Password) ValidatePassword(password string) error {
	if len(password) < p.config.MinLength {
		return fmt.Errorf("password must be at least %d characters long", p.config.MinLength)
	}

	if p.config.RequireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if p.config.RequireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if p.config.RequireNumbers && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	if p.config.RequireSpecial && !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// CheckStrength 检查密码强度
func (p *Password) CheckStrength(password string) Strength {
	score := 0

	// 长度检查
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}

	// 字符类型检查
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		score++
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		score++
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		score++
	}
	if regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).MatchString(password) {
		score++
	}

	// 复杂性检查
	if !hasRepeatedChars(password) {
		score++
	}
	if !hasSequentialChars(password) {
		score++
	}

	// 根据得分返回强度等级
	switch {
	case score >= 7:
		return StrengthStrong
	case score >= 5:
		return StrengthGood
	case score >= 3:
		return StrengthFair
	default:
		return StrengthWeak
	}
}

// Generate 生成强密码
func (p *Password) Generate(length int) (string, error) {
	if length < p.config.MinLength {
		length = p.config.MinLength
	}

	// 定义字符集
	var charsets []string
	if p.config.RequireLowercase {
		charsets = append(charsets, "abcdefghijklmnopqrstuvwxyz")
	}
	if p.config.RequireUppercase {
		charsets = append(charsets, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	if p.config.RequireNumbers {
		charsets = append(charsets, "0123456789")
	}
	if p.config.RequireSpecial {
		charsets = append(charsets, "!@#$%^&*()_+-=[]{}|;:,.<>?")
	}

	// 如果没有特殊要求，使用默认字符集
	if len(charsets) == 0 {
		charsets = []string{
			"abcdefghijklmnopqrstuvwxyz",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"0123456789",
			"!@#$%^&*()_+-=[]{}|;:,.<>?",
		}
	}

	allChars := strings.Join(charsets, "")
	password := make([]byte, length)

	// 确保至少包含每种要求的字符类型
	for i, charset := range charsets {
		if i >= length {
			break
		}
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random character: %w", err)
		}
		password[i] = charset[n.Int64()]
	}

	// 填充剩余位置
	for i := len(charsets); i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random character: %w", err)
		}
		password[i] = allChars[n.Int64()]
	}

	// 随机打乱密码字符顺序
	for i := range password {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {
			return "", fmt.Errorf("failed to shuffle password: %w", err)
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

// GenerateWithStrength 生成指定强度的密码
func (p *Password) GenerateWithStrength(strength Strength) (string, error) {
	var length int
	switch strength {
	case StrengthWeak:
		length = max(p.config.MinLength, 8)
	case StrengthFair:
		length = max(p.config.MinLength, 10)
	case StrengthGood:
		length = max(p.config.MinLength, 12)
	case StrengthStrong:
		length = max(p.config.MinLength, 16)
	default:
		length = max(p.config.MinLength, 12)
	}

	// 生成密码直到满足强度要求
	for attempts := 0; attempts < 100; attempts++ {
		password, err := p.Generate(length)
		if err != nil {
			continue
		}
		if p.CheckStrength(password) >= strength {
			return password, nil
		}
	}

	// 如果无法生成满足要求的密码，生成最长的密码
	return p.Generate(max(length, 20))
}

// GetConfig 获取配置
func (p *Password) GetConfig() *Config {
	return p.config
}

// SetConfig 更新配置
func (p *Password) SetConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	p.config = config
	return nil
}

// 便利函数

// HashPassword 使用默认配置加密密码
func HashPassword(password string) (string, error) {
	p, _ := New(nil)
	return p.Hash(password)
}

// CheckPasswordHash 使用默认配置验证密码
func CheckPasswordHash(password, hash string) bool {
	p, _ := New(nil)
	return p.Verify(password, hash)
}

// GenerateStrongPassword 生成强密码
func GenerateStrongPassword(length int) (string, error) {
	p, _ := New(nil)
	return p.Generate(length)
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) error {
	p, _ := New(nil)
	return p.ValidatePassword(password)
}

// CheckPasswordStrength 检查密码强度
func CheckPasswordStrength(password string) Strength {
	p, _ := New(nil)
	return p.CheckStrength(password)
}

// 辅助函数

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func hasRepeatedChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}

func hasSequentialChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i+1] == password[i]+1 && password[i+2] == password[i]+2 {
			return true
		}
		if password[i+1] == password[i]-1 && password[i+2] == password[i]-2 {
			return true
		}
	}
	return false
}