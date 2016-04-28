package options

import "time"

// Type - Option type identifier.
type Type string

const (
	// TypeKey - Option of type Context key.
	TypeKey Type = "key"

	// TypeBool - Option of type Bool.
	TypeBool Type = "bool"

	// TypeString - Option of type String.
	TypeString Type = "string"

	// TypeDuration - Option of type Duration.
	TypeDuration Type = "duration"

	// TypeDestination - Option of type Destination.
	TypeDestination Type = "destination"

	// TypeFloat - Option of type Float.
	TypeFloat Type = "float"

	// TypeInt - Option of type int.
	TypeInt Type = "int"

	// TypeInt64 - Option of type int64.
	TypeInt64 Type = "int64"

	// TypeMap - Option of type Map.
	TypeMap Type = "map"

	// TypeList - Option of type List.
	TypeList Type = "list"

	// TypeTemplate - Template option value type.
	TypeTemplate Type = "template"

	// TypeEmpty - Empty option value type.
	TypeEmpty Type = "empty"

	// TypeUnknown - Unknown option value type.
	TypeUnknown Type = "unknown"
)

// CheckType - Returns true if value type matches given type.
func CheckType(t Type, v interface{}) bool {
	if t == TypeKey || t == TypeTemplate || t == TypeDestination {
		t = TypeString
	}
	vt := ValueType(v)
	if vt == t {
		return true
	}
	return IsConvertible(vt, t)
}

// ValueType - Gets value type.
func ValueType(v interface{}) Type {
	if v == nil {
		return TypeEmpty
	}
	switch v.(type) {
	case string:
		return TypeString
	case bool:
		return TypeBool
	case float64:
		return TypeFloat
	case int:
		return TypeInt
	case int64:
		return TypeInt64
	case map[string]string:
		return TypeMap
	case map[string]interface{}:
		return TypeMap
	case map[interface{}]interface{}:
		return TypeMap
	case []string:
		return TypeList
	case []interface{}:
		return TypeList
	case time.Duration:
		return TypeDuration
	}
	return TypeUnknown
}

// ContainsType - Returns true if value type matches given type.
func ContainsType(list []Type, t Type) bool {
	for _, el := range list {
		if el == t {
			return true
		}
	}
	return false
}
