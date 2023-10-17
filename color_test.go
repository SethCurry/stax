package stax_test

import (
	"testing"

	"github.com/SethCurry/stax"
	"github.com/google/go-cmp/cmp"
)

func Test_Color_UnmarshalText(t *testing.T) {
	testCases := []struct {
		Value         string
		Expected      *stax.Color
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

	for _, tc := range testCases {
		t.Run(tc.Value, func(t *testing.T) {
			var color stax.Color
			err := color.UnmarshalText([]byte(tc.Value))
			errDiff := cmp.Diff(tc.ExpectedError, err)
			if errDiff != "" {
				t.Errorf("unexpected error: %s", errDiff)
			}

			if tc.Expected != nil {
				colorDiff := cmp.Diff(&color, tc.Expected)
				if colorDiff != "" {
					t.Errorf("unexpected color: %s", colorDiff)
				}
			}
		})
	}
}
