package helpers

import "testing"

func TestWithDefaults(t *testing.T) {
	from := map[string]interface{}{
		"first": map[string]interface{}{
			"test": map[string]string{
				"key": "value",
			},
		},
	}
	fromDefaults := map[string]interface{}{
		"first": map[string]interface{}{
			"test": map[string]interface{}{
				"key2": "value2",
			},
		},
	}

	res := WithDefaults(from, fromDefaults).(map[string]interface{})

	first := res["first"].(map[string]interface{})
	test := first["test"].(map[string]string)

	if len(test) != 2 {
		t.Errorf("Invalid length: %d\n", len(test))
	}

	e := map[string]interface{}{
		"key":  "value",
		"key2": "value2",
	}

	for k, v := range e {
		if v != test[k] {
			t.Errorf("Invalid key %q value: %#v\n", k, v)
		}
	}
}
