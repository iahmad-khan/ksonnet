package utils

import (
	"testing"

	jsonnet "github.com/strickyak/jsonnet_cgo"
)

// check there is no err, and a == b.
func check(t *testing.T, err error, actual, expected string) {
	if err != nil {
		t.Errorf("Expected %q, got error: %q", expected, err.Error())
	} else if actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestParseJson(t *testing.T) {
	vm := jsonnet.Make()
	defer vm.Destroy()
	RegisterNativeFuncs(vm, NewIdentityResolver())

	_, err := vm.EvaluateSnippet("failtest", `std.native("parseJson")("barf{")`)
	if err == nil {
		t.Errorf("parseJson succeeded on invalid json")
	}

	x, err := vm.EvaluateSnippet("test", `std.native("parseJson")("null")`)
	check(t, err, x, "null\n")

	x, err = vm.EvaluateSnippet("test", `
    local a = std.native("parseJson")('{"foo": 3, "bar": 4}');
    a.foo + a.bar`)
	check(t, err, x, "7\n")
}

func TestParseYaml(t *testing.T) {
	vm := jsonnet.Make()
	defer vm.Destroy()
	RegisterNativeFuncs(vm, NewIdentityResolver())

	_, err := vm.EvaluateSnippet("failtest", `std.native("parseYaml")("[barf")`)
	if err == nil {
		t.Errorf("parseYaml succeeded on invalid yaml")
	}

	x, err := vm.EvaluateSnippet("test", `std.native("parseYaml")("")`)
	check(t, err, x, "[ ]\n")

	x, err = vm.EvaluateSnippet("test", `
    local a = std.native("parseYaml")("foo:\n- 3\n- 4\n")[0];
    a.foo[0] + a.foo[1]`)
	check(t, err, x, "7\n")

	x, err = vm.EvaluateSnippet("test", `
    local a = std.native("parseYaml")("---\nhello\n---\nworld");
    a[0] + a[1]`)
	check(t, err, x, "\"helloworld\"\n")
}
