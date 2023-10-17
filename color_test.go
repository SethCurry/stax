package stax_test

import (
	"testing"

	"github.com/SethCurry/stax"
	"github.com/google/go-cmp/cmp"
)

func Test_Color_UnmarshalText(t *testing.T) {
	testCases := []struct {
		Value         string
		Expected      func() stax.Color
		ExpectedError error
	}{
		{
			Value:         "{R}",
			Expected:      stax.ColorRed,
			ExpectedError: nil,
		},
		{
			Value:         "{W}",
			Expected:      stax.ColorWhite,
			ExpectedError: nil,
		},
		{
			Value:         "{U}",
			Expected:      stax.ColorBlue,
			ExpectedError: nil,
		},
		{
			Value:         "{B}",
			Expected:      stax.ColorBlack,
			ExpectedError: nil,
		},
		{
			Value:         "{G}",
			Expected:      stax.ColorGreen,
			ExpectedError: nil,
		},
		{
			Value:         "G",
			Expected:      stax.ColorGreen,
			ExpectedError: nil,
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		testCase := tc
		t.Run(tc.Value, func(t *testing.T) {
			t.Parallel()
			var color stax.Color
			err := color.UnmarshalText([]byte(testCase.Value))
			errDiff := cmp.Diff(testCase.ExpectedError, err)
			if errDiff != "" {
				t.Errorf("unexpected error: %s", errDiff)
			}

			if testCase.Expected != nil {
				colorDiff := cmp.Diff(color, testCase.Expected())
				if colorDiff != "" {
					t.Errorf("unexpected color: %s", colorDiff)
				}
			}
		})
	}
}
