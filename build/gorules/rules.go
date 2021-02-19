// +build ignore

// WARNING Run `golangci-lint cache clean` after modifying this file!
package gorules

import "github.com/quasilyte/go-ruleguard/dsl/fluent"

func unexportedSensitive(m fluent.Matcher) {
	// TODO Won't match if exported sensitive field comes first:
	//   struct{ Exported sensitive.Type; unexported sensitive.Type }
	// TODO Add support for renamed package?
	// TODO Detect if struct with exported sensitive field is itself
	// inside unexported field in parent struct? Only if in print/log!
	// TODO Detect panic(sensitiveValue).
	m.Match(`struct{$*_; $field sensitive.$_; $*_}`,
		`struct{$*_; $field AccessToken; $*_}`,
		`struct{$*_; $field $_.AccessToken; $*_}`,
	).Where(m["field"].Text.Matches(`^[^A-Z]`)).
		Report("found sensitive value in unexported field $field")
}

// Forbid print/println because they'll output sensitive values as is.
func printfunc(m fluent.Matcher) {
	m.Match(`print($*args)`).
		Suggest("fmt.Print($args)")
	m.Match(`println($*args)`).
		Suggest("fmt.Println($args)")
}
