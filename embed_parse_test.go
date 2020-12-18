package main

import "testing"


type ep_test_case struct {
	input string
	expect_pointer bool
	expect_path string
	expect_name string
}

var ep_samples = []ep_test_case {
	ep_test_case {
		input: "Simple",
		expect_pointer: false,
		expect_path: "",
		expect_name: "Simple",
	},
	ep_test_case {
		input: "*StarSimple",
		expect_pointer: true,
		expect_path: "",
		expect_name: "StarSimple",
	},
	ep_test_case {
		input: "*foo/bar/Ident",
		expect_pointer: true,
		expect_path: "foo/bar",
		expect_name: "Ident",
	},
}


func TestParseEmbed(t *testing.T) {
	for _, testcase := range ep_samples {
		t.Logf("input: %q", testcase.input)
		parsed := ParseEmbedImpl(testcase.input)
		if parsed.err != nil {
			t.Errorf("%s", parsed.err)
			continue
		}
		if parsed.pointer != testcase.expect_pointer {
			t.Errorf("Expected pointer to be %v, got %v",
				testcase.expect_pointer,
				parsed.pointer)
		}
		if got := parsed.Path(); got != testcase.expect_path {
			t.Errorf("Expected path to be %q, got %q",
				testcase.expect_path,
				got)
		}
		if got := parsed.Name(); got != testcase.expect_name {
			t.Errorf("Expected name to be %q, got %q",
				testcase.expect_name,
				got)
		}
	}
}


