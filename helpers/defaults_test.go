package helpers

import (
	"testing"

	"tower.pro/renderer/template"
)

func TestWithDefaults(t *testing.T) {
	from := map[string]interface{}{
		"first": map[string]interface{}{
			"deep": map[string]interface{}{
				"test": map[string]string{
					"key": "value",
				},
			},
		},
	}
	fromDefaults := map[string]interface{}{
		"first": map[string]interface{}{
			"deep": map[string]interface{}{
				"test": map[string]string{
					"key2": "value2",
				},
			},
		},
	}

	res := WithDefaults(from, fromDefaults).(template.Context)

	first := res["first"].(template.Context)
	deep := first["deep"].(template.Context)
	test := deep["test"].(map[string]string)

	if len(test) != 2 {
		t.Errorf("Invalid length: %d\n", len(test))
	}

	e := template.Context{
		"key":  "value",
		"key2": "value2",
	}

	for k, v := range e {
		if v != test[k] {
			t.Errorf("Invalid key %q value: %#v\n", k, v)
		}
	}
}
