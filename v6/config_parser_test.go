package configcat

import (
	"testing"
)

func TestConfigParser_Parse(t *testing.T) {
	jsonBody := "{ \"f\": { \"keyDouble\": { \"v\": 120.121238476, \"p\": [], \"r\": [], \"i\":\"\" }}}"
	parser := newParser(DefaultLogger(LogLevelWarn))

	val, err := parser.parse(jsonBody, "keyDouble", nil)

	if err != nil || val != 120.121238476 {
		t.Error("Expecting 120.121238476 as interface")
	}
}

func TestConfigParser_BadJson(t *testing.T) {
	jsonBody := ""
	parser := newParser(DefaultLogger(LogLevelWarn))

	_, err := parser.parse(jsonBody, "keyDouble", nil)

	if err == nil {
		t.Error("Expecting JSON error")
	}

	t.Log(err.Error())
}

func TestConfigParser_BadJson_String(t *testing.T) {
	jsonBody := ""
	parser := newParser(DefaultLogger(LogLevelWarn))

	_, err := parser.parse(jsonBody, "key", nil)

	if err == nil {
		t.Error("Expecting JSON error")
	}

	t.Log(err.Error())
}

func TestConfigParser_WrongKey(t *testing.T) {
	jsonBody := "{ \"keyDouble\": { \"Value\": 120.121238476, \"SettingType\": 0, \"RolloutPercentageItems\": [], \"RolloutRules\": [] }}"
	parser := newParser(DefaultLogger(LogLevelWarn))

	_, err := parser.parse(jsonBody, "wrongKey", nil)

	if err == nil {
		t.Error("Expecting key not found error")
	}

	t.Log(err.Error())
}

func TestConfigParser_EmptyNode(t *testing.T) {
	jsonBody := "{ \"keyDouble\": { }}"
	parser := newParser(DefaultLogger(LogLevelWarn))

	_, err := parser.parse(jsonBody, "keyDouble", nil)

	if err == nil {
		t.Error("Expecting invalid JSON error")
	}

	t.Log(err.Error())
}
