package templateutils

import (
	"strings"
	"text/template"
)

var funcMap = map[string]interface{}{
	// String
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,

	"toCamel":  ToCamel,
	"toKebab":  ToKebab,
	"toLower":  strings.ToLower,
	"toPascal": ToPascal,
	"toSnake":  ToSnack,
	"toUpper":  strings.ToUpper,

	"zeroValue": ZeroValue,

	// Case
	"toInt64": ToInt64,

	// Math
	"add": Add,

	// Utils
	"in":        In,
	"makeSlice": MakeSlice,
}

// FuncMap will return the template utils as FuncMap.
func FuncMap() template.FuncMap {
	return funcMap
}
