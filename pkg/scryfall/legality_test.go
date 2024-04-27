package scryfall_test

import (
	"testing"

	"github.com/SethCurry/stax/pkg/scryfall"
)

func Test_Legality_UnmarshalText(t *testing.T) {
	testCases := []struct {
		name    string
		txt     []byte
		want    scryfall.Legality
		wantErr bool
	}{
		{
			name:    "legal",
			txt:     []byte("legal"),
			want:    scryfall.LegalityLegal,
			wantErr: false,
		},
		{
			name:    "not legal",
			txt:     []byte("not_legal"),
			want:    scryfall.LegalityNotLegal,
			wantErr: false,
		},
		{
			name:    "unknown legality",
			txt:     []byte("unknown"),
			want:    scryfall.Legality(""),
			wantErr: true,
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			var l scryfall.Legality

			err := l.UnmarshalText(v.txt)
			if (err != nil) != v.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if l != v.want {
				t.Errorf("unexpected legality: got %v, want %v", l, v.want)
			}
		})
	}
}

func Test_Legality_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name    string
		txt     []byte
		want    scryfall.Legality
		wantErr bool
	}{
		{
			name:    "legal",
			txt:     []byte("\"legal\""),
			want:    scryfall.LegalityLegal,
			wantErr: false,
		},
		{
			name:    "not legal",
			txt:     []byte("\"not_legal\""),
			want:    scryfall.LegalityNotLegal,
			wantErr: false,
		},
		{
			name:    "unknown legality",
			txt:     []byte("\"unknown\""),
			want:    scryfall.Legality(""),
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			txt:     []byte("invalid"),
			want:    scryfall.Legality(""),
			wantErr: true,
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			var l scryfall.Legality

			err := l.UnmarshalJSON(v.txt)
			if (err != nil) != v.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if l != v.want {
				t.Errorf("unexpected legality: got %v, want %v", l, v.want)
			}
		})
	}
}
