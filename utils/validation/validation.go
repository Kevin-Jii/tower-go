package validation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// InitValidator 初始化验证器
func InitValidator() *validator.Validate {
	validate := validator.New()

	// 注册自定义验证规则
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("employee_no", validateEmployeeNo)

	return validate
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 11 {
		return false
	}
	
	// 中国手机号正则
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	if len(password) < 6 {
		return false
	}
	
	// 至少包含字母和数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	return hasLetter && hasNumber
}

// validateEmployeeNo 验证工号
func validateEmployeeNo(fl validator.FieldLevel) bool {
	no := fl.Field().String()
	
	// 工号必须是6位数字
	if len(no) != 6 {
		return false
	}
	
	_, err := strconv.Atoi(no)
	return err == nil
}

// ValidationErrors 验证错误处理
type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}

// FormatValidationErrors 格式化验证错误
func FormatValidationErrors(err error) ValidationErrors {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()
			
			switch tag {
			case "required":
				errors[fieldName] = fieldName + "是必填的"
			case "min":
				errors[fieldName] = fieldName + "长度不能小于" + fieldError.Param()
			case "max":
				errors[fieldName] = fieldName + "长度不能大于" + fieldError.Param()
			case "len":
				errors[fieldName] = fieldName + "长度必须为" + fieldError.Param()
			case "email":
				errors[fieldName] = "请输入有效的邮箱地址"
			case "phone":
				errors[fieldName] = "请输入有效的手机号"
			case "oneof":
				errors[fieldName] = fieldName + "必须是以下值之一: " + fieldError.Param()
			default:
				errors[fieldName] = fieldName + "验证失败"
			}
		}
	}

	return ValidationErrors{Errors: errors}
}

// SanitizeInput 清理输入数据
func SanitizeInput(input string) string {
	// 移除首尾空格
	input = strings.TrimSpace(input)
	
	// 移除特殊字符（防止 SQL 注入和 XSS）
	input = regexp.MustCompile(`[<>""'\\]`).ReplaceAllString(input, "")
	
	return input
}

// ValidatePhone 验证手机号（独立函数）
func ValidatePhone(phone string) bool {
	if len(phone) != 11 {
		return false
	}
	
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) (bool, string) {
	if len(password) < 6 {
		return false, "密码长度不能小于6位"
	}
	
	if len(password) > 32 {
		return false, "密码长度不能大于32位"
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	
	strength := 0
	if hasUpper { strength++ }
	if hasLower { strength++ }
	if hasNumber { strength++ }
	if hasSpecial { strength++ }
	
	if strength < 2 {
		return false, "密码必须包含字母和数字，建议包含特殊字符"
	}
	
	return true, ""
}
