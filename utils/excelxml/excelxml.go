package excelxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type Sheet struct {
	Name    string
	Headers []string
	Rows    [][]interface{}
}

func Build(sheets []Sheet) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<?mso-application progid="Excel.Sheet"?>` + "\n")
	b.WriteString(`<Workbook xmlns="urn:schemas-microsoft-com:office:spreadsheet" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns:ss="urn:schemas-microsoft-com:office:spreadsheet">`)
	b.WriteString(`<Styles><Style ss:ID="Header"><Font ss:Bold="1"/><Interior ss:Color="#D9EAF7" ss:Pattern="Solid"/></Style></Styles>`)
	for _, sheet := range sheets {
		name := sheet.Name
		if strings.TrimSpace(name) == "" {
			name = "Sheet1"
		}
		b.WriteString(`<Worksheet ss:Name="`)
		writeEscaped(&b, truncateSheetName(name))
		b.WriteString(`"><Table>`)
		if len(sheet.Headers) > 0 {
			b.WriteString(`<Row>`)
			for _, h := range sheet.Headers {
				writeCell(&b, h, true)
			}
			b.WriteString(`</Row>`)
		}
		for _, row := range sheet.Rows {
			b.WriteString(`<Row>`)
			for _, cell := range row {
				writeCell(&b, cell, false)
			}
			b.WriteString(`</Row>`)
		}
		b.WriteString(`</Table></Worksheet>`)
	}
	b.WriteString(`</Workbook>`)
	return []byte(b.String())
}

func Filename(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "export"
	}
	return fmt.Sprintf("%s-%s.xls", prefix, time.Now().Format("20060102150405"))
}

func writeCell(b *strings.Builder, v interface{}, header bool) {
	style := ""
	if header {
		style = ` ss:StyleID="Header"`
	}
	b.WriteString(`<Cell`)
	b.WriteString(style)
	b.WriteString(`><Data ss:Type="String">`)
	writeEscaped(b, formatValue(v))
	b.WriteString(`</Data></Cell>`)
}

func writeEscaped(b *strings.Builder, s string) {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	b.WriteString(buf.String())
}

func formatValue(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case time.Time:
		if x.IsZero() {
			return ""
		}
		if x.Hour() == 0 && x.Minute() == 0 && x.Second() == 0 {
			return x.Format("2006-01-02")
		}
		return x.Format("2006-01-02 15:04:05")
	case *time.Time:
		if x == nil {
			return ""
		}
		return formatValue(*x)
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprint(v)
	}
}

func truncateSheetName(name string) string {
	replacer := strings.NewReplacer("[", "(", "]", ")", "*", "", "?", "", "/", "-", "\\", "-", ":", "-")
	name = strings.TrimSpace(replacer.Replace(name))
	rs := []rune(name)
	if len(rs) > 31 {
		return string(rs[:31])
	}
	if name == "" {
		return "Sheet1"
	}
	return name
}
