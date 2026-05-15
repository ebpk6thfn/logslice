package output

import "fmt"

// ParseFormat converts a string into a Format constant.
// Returns an error if the string does not match a known format.
// An empty string defaults to FormatRaw.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatRaw, FormatNumbered:
		return Format(s), nil
	case "":
		return FormatRaw, nil
	default:
		return "", fmt.Errorf("unknown output format %q: must be one of [raw, numbered]", s)
	}
}

// String returns the string representation of the Format.
func (f Format) String() string {
	return string(f)
}

// IsValid reports whether the Format is a known, supported format.
func (f Format) IsValid() bool {
	switch f {
	case FormatRaw, FormatNumbered:
		return true
	default:
		return false
	}
}

// ValidFormats returns all supported format names as a slice of strings.
func ValidFormats() []string {
	return []string{
		string(FormatRaw),
		string(FormatNumbered),
	}
}
