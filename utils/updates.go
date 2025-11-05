package utils

import (
	"reflect"
	"unicode"
)

// BuildUpdatesFromReq 根据请求结构体构建更新字段 map：
// 规则：
// 1. 指针字段：非 nil 则加入（值为指针解引用）。
// 2. 字符串：非空加入。
// 3. 数值/布尔：非零值加入；如需允许零值更新请在字段 tag 中增加 `patch:"allowZero"`。
// 4. 结构体内嵌或切片不自动展开（需手工处理）。
// 5. 可用 tag `patch:"always"` 强制加入，无视零值；`patch:"allowZero"` 允许 0 / false。
// 注意：仅处理导出字段。
func BuildUpdatesFromReq(req interface{}) map[string]interface{} {
	updates := make(map[string]interface{})
	if req == nil {
		return updates
	}
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return updates
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return updates
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		f := typ.Field(i)
		// 跳过未导出字段
		if f.PkgPath != "" {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		// 字段名取 json tag(逗号前) 或原字段名的 snake 风格（这里简单用原名小写）
		name := jsonTag
		if idx := findComma(jsonTag); idx != -1 {
			name = jsonTag[:idx]
		}
		if name == "" {
			name = toDBName(f.Name)
		}

		patchTag := f.Tag.Get("patch")
		allowZero := patchTag == "allowZero" || patchTag == "always"
		always := patchTag == "always"

		fv := val.Field(i)
		kind := fv.Kind()

		// 指针字段
		if kind == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			elem := fv.Elem()
			updates[name] = elem.Interface()
			continue
		}

		// 根据类型决定
		switch kind {
		case reflect.String:
			if fv.String() != "" || always {
				updates[name] = fv.String()
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			if fv.Interface() != reflect.Zero(fv.Type()).Interface() || allowZero || always {
				updates[name] = fv.Interface()
			}
		case reflect.Bool:
			if fv.Bool() || allowZero || always { // false 仅在 allowZero/always 时允许
				updates[name] = fv.Bool()
			}
		default:
			// 其它类型默认跳过（避免误更新复杂结构）
			if always {
				updates[name] = fv.Interface()
			}
		}
	}
	return updates
}

func findComma(s string) int {
	for i, ch := range s {
		if ch == ',' {
			return i
		}
	}
	return -1
}

// 简单字段名转 db 名：全部小写。可扩展为 snake_case。
// 将驼峰转为 snake_case，例如 StoreID -> store_id, RoleCode -> role_code
// 保留原有全部小写逻辑作为基础
func toDBName(s string) string {
	if s == "" {
		return s
	}
	var out []rune
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			// 前一位存在且不是下划线且：前一位是小写 或 后一位是小写（处理连续大写拆分）
			if i > 0 && runes[i-1] != '_' && (unicode.IsLower(runes[i-1]) || (i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
				out = append(out, '_')
			}
			out = append(out, unicode.ToLower(r))
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}
