package options

// IsConvertible - Checks if type is convertible to another one.
func IsConvertible(from Type, to Type) bool {
	list, ok := Convertible[from]
	if !ok {
		return false
	}
	return ContainsType(list, to)
}

// Convertible - List of types every type can convert to.
var Convertible = map[Type][]Type{
	TypeInt: []Type{
		TypeInt64,
		TypeFloat,
	},
	TypeInt64: []Type{
		TypeInt,
		TypeFloat,
	},
	TypeFloat: []Type{
		TypeInt,
		TypeFloat,
		TypeInt64,
	},
	TypeString: []Type{
		TypeInt,
		TypeInt64,
		TypeFloat,
		TypeBool,
		TypeDuration,
		TypeKey,
		TypeDestination,
		TypeTemplate,
	},
	TypeTemplate: []Type{
		TypeString,
		TypeInt,
		TypeInt64,
		TypeFloat,
		TypeBool,
		TypeDuration,
		TypeKey,
		TypeDestination,
		TypeMap,
		TypeList,
	},
	TypeKey: []Type{
		TypeKey,
		TypeBool,
		TypeString,
		TypeDuration,
		TypeDestination,
		TypeFloat,
		TypeInt,
		TypeInt64,
		TypeMap,
		TypeList,
	},
	TypeList: []Type{
		TypeStringList,
	},
	TypeStringList: []Type{
		TypeList,
	},
}
