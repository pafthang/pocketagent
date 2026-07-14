package tools

import "testing"

func TestCodeExecDisabled(t *testing.T) {
	cfg := Config{CodeExecEnabled: false}
	_, err := codeExec(cfg, map[string]interface{}{"code": "print(1)"})
	if err == nil {
		t.Fatal("expected disabled error")
	}
}