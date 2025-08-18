package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNumeric(t *testing.T) {
	type input struct {
		s string
	}
	type expected struct {
		valid bool
	}

	tester := func(t *testing.T, input input, expected expected) {
		t.Log(input.s)
		if expected.valid {
			require.True(t, IsNumeric(input.s))
		} else {
			require.False(t, IsNumeric(input.s))
		}
	}

	tester(t, input{"10"}, expected{true})
	tester(t, input{"-10"}, expected{true})
	tester(t, input{"0"}, expected{true})
	tester(t, input{"1"}, expected{true})
	tester(t, input{"-1"}, expected{true})
	tester(t, input{"10.0"}, expected{false})
	tester(t, input{"a"}, expected{false})
}

func TestIsBech32(t *testing.T) {
	type input struct {
		s string
	}
	type expected struct {
		valid bool
	}

	tester := func(t *testing.T, input input, expected expected) {
		t.Log(input.s)
		if expected.valid {
			require.True(t, IsBech32(input.s))
		} else {
			require.False(t, IsBech32(input.s))
		}
	}

	// valid
	tester(t, input{"g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"}, expected{true})
	tester(t, input{"g1ffzxha57dh0qgv9ma5v393ur0zexfvp6lsjpae"}, expected{true})
	tester(t, input{"g16a7etgm9z2r653ucl36rj0l2yqcxgrz2jyegzx"}, expected{true})
	tester(t, input{"g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d"}, expected{true})

	// invalid
	tester(t, input{""}, expected{false})
	tester(t, input{"a1asdfsdfsdfasdfwqefqwfwqd"}, expected{false})
	tester(t, input{"g1asdf"}, expected{false})
}
