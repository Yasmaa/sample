package templateutils

import (
	"regexp"
	"strings"
)

var (
	intTypeRegex     = regexp.MustCompile(`^[u]?int`)
	floatTypeRegex   = regexp.MustCompile(`^float`)
	complexTypeRegex = regexp.MustCompile(`^complex`)
)

// ZeroValue will return input type name's default zero value.
func ZeroValue(in string) string {
	switch in {
	case "bool":
		return "false"
	case "string":
		return `""`
	case "interface{}":
		return "nil"
	}

	// Handle all integer, float and complex via return 0.
	if intTypeRegex.MatchString(in) || floatTypeRegex.MatchString(in) || complexTypeRegex.MatchString(in) {
		return "0"
	}

	// If type starts with a "*", it must be a pointer, return nil directly.
	if strings.HasPrefix(in, "*") {
		return "nil"
	}

	// If type starts with "[]", it must a slice, return nil directly.
	if strings.HasPrefix(in, "[]") {
		return "nil"
	}

	return in + "{}"
}

const (
	stateUpper = "upper"
	stateLower = "lower"
	stateDigit = "digit"
	stateSplit = "split"
)

func getState(r rune) string {
	switch {
	case r >= '0' && r <= '9':
		return stateDigit
	case r >= 'a' && r <= 'z':
		return stateLower
	case r >= 'A' && r <= 'Z':
		return stateUpper
	default:
		return stateSplit
	}
}

// SplitStringViaSpecialChars will split strings via special chars.
//
// Examples:
//   a_bc => [a, bc]
//   a-bc => [a, bc]
//   a bc => [a, bc]
//   a--b => [a, b]
func SplitStringViaSpecialChars(in string) []string {
	idx := 0
	ans := make([]string, 1)
	for _, v := range in {
		s := getState(v)
		// If not a split char, append to current string.
		if s != stateSplit {
			ans[idx] += string(v)
			continue
		}
		// If current string is empty, just skip.
		if len(ans[idx]) == 0 {
			continue
		}
		// If current string is not empty, append a new string.
		idx++
		ans = append(ans, "")
	}
	return ans
}

// SplitStringViaUpperChars will split string via upper case chars.
//
// Examples:
//   abc => [abc]
//   Abc => [Abc]
//   AbC => [Ab, C]
//   AbCd => [Ab, Cd]
//   AAA => [AAA]
//   AAABb => [AAA, Bb]
//
// All state: empty, upper, digit+lower
// All action: add, append
//
// prev, cur, next, action
// empty, upper, empty, append
// empty, upper, upper, append
// empty, upper, lower, append
// upper, upper, empty, append
// upper, upper, upper, append
// upper, upper, lower, add
// lower, upper, empty, add
// lower, upper, upper, add
// lower, upper, lower, add
// empty, lower, empty, append
// empty, lower, upper, append
// empty, lower, lower, append
// upper, lower, empty, append
// upper, lower, upper, append
// upper, lower, lower, append
// lower, lower, empty, append
// lower, lower, upper, append
// lower, lower, lower, append
func SplitStringViaUpperChars(in string) []string {
	if len(in) == 0 {
		return nil
	}

	idx := 0
	ans := make([]string, 1)
	for i := 0; i < len(in); i++ {
		// If prev is empty, we always add current into buf.
		// This case handles all prev == empty cases.
		if i == 0 {
			ans[idx] += string(in[i])
			continue
		}
		// Handle all add cases.
		curs := getState(rune(in[i]))
		if curs == stateUpper {
			prevs := getState(rune(in[i-1]))
			// lower, upper, empty, add
			// lower, upper, upper, add
			// lower, upper, lower, add
			if prevs == stateLower || prevs == stateDigit {
				ans = append(ans, string(in[i]))
				idx++
				continue
			}
			// upper, upper, lower, add
			if i != len(in)-1 {
				nexts := getState(rune(in[i+1]))
				if nexts == stateLower || nexts == stateDigit {
					ans = append(ans, string(in[i]))
					idx++
					continue
				}
			}
		}
		// all other cases
		ans[idx] += string(in[i])
	}
	return ans
}

// splitStringInParts will split string into different parts.
//
// Notes:
//   - Only be used in function/class name converts
//   - Only be used for ASCII chars
// Examples:
//   abc => [abc]
//   ABC => [abc]
//   a_b_c => [a, b, c]
//   A_B_C => [a, b, c]
//   A_b_C => [a, b, c]
//   AB_C => [ab, c]
//   a-b-c => [a, b, c]
//   a b c => [a, b, c]
//   AAbAACcccc => [a, ab, aa, ccccc]
func splitStringInParts(in string) []string {
	ans := make([]string, 0)
	for _, v := range SplitStringViaSpecialChars(in) {
		x := SplitStringViaUpperChars(v)
		for i := range x {
			// Convert all string to lower case.
			x[i] = strings.ToLower(x[i])
		}
		ans = append(ans, x...)
	}
	return ans
}

// ToCamel combines words by capitalizing all words following the first word and
// removing the space
func ToCamel(in string) string {
	x := splitStringInParts(in)
	for i := range x {
		if i == 0 {
			continue
		}
		if commonInitialisms[x[i]] {
			x[i] = strings.ToUpper(x[i])
			continue
		}
		x[i] = ToUpperFirst(x[i])
	}
	return strings.Join(x, "")
}

// ToPascal combines words by capitalizing all words (even the first word) and
// removing the space
func ToPascal(in string) string {
	x := splitStringInParts(in)
	for i := range x {
		if commonInitialisms[x[i]] {
			x[i] = strings.ToUpper(x[i])
			continue
		}
		x[i] = ToUpperFirst(x[i])
	}
	return strings.Join(x, "")
}

// ToSnack combines words by replacing each space with an underscore (_)
func ToSnack(in string) string {
	x := splitStringInParts(in)
	return strings.Join(x, "_")
}

// ToKebab combines words by replacing each space with a dash (-)
func ToKebab(in string) string {
	x := splitStringInParts(in)
	return strings.Join(x, "-")
}

// ToUpperFirst will convert first letter to upper case.
func ToUpperFirst(in string) string {
	return strings.ToUpper(in[:1]) + in[1:]
}

// borrowed from https://github.com/golang/lint/blob/master/lint.go#L770
var commonInitialisms = map[string]bool{
	"acl":   true, // ACL
	"api":   true, // API
	"ascii": true, // ASCII
	"cpu":   true, // CPU
	"css":   true, // CSS
	"dns":   true, // DNS
	"eof":   true, // EOF
	"guid":  true, // GUID
	"html":  true, // HTML
	"http":  true, // HTTP
	"https": true, // HTTPS
	"id":    true, // ID
	"ip":    true, // IP
	"json":  true, // JSON
	"lhs":   true, // LHS
	"qps":   true, // QPS
	"ram":   true, // RAM
	"rhs":   true, // RHS
	"rpc":   true, // RPC
	"sla":   true, // SLA
	"smtp":  true, // SMTP
	"sql":   true, // SQL
	"ssh":   true, // SSH
	"tcp":   true, // TCP
	"tls":   true, // TLS
	"ttl":   true, // TTL
	"udp":   true, // UDP
	"ui":    true, // UI
	"uid":   true, // UID
	"uuid":  true, // UUID
	"uri":   true, // URI
	"url":   true, // URL
	"utf8":  true, // UTF8
	"vm":    true, // VM
	"xml":   true, // XML
	"xmpp":  true, // XMPP
	"xsrf":  true, // XSRF
	"xss":   true, // XSS
}
