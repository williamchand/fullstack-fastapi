package util

import (
	"bytes"
	"encoding/json"
	"html/template"
)

func StructToJSONMap(data interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func FillTextTemplate(tplStr string, data interface{}) (string, error) {
	m, err := StructToJSONMap(data)
	if err != nil {
		return "", err
	}

	tpl, err := template.New("").Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, m); err != nil {
		return "", err
	}

	return buf.String(), nil
}
