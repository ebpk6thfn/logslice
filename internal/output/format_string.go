package output

import "fmt"

// ParseFormat converts a string into a Format constant.
// Returns an error if the string does not match a known format.
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

// ValidFormats returns all supported format names as a slice of strings.
func ValidFormats() []string {
	return []string{
		string(FormatRaw),
		string(FormatNumbered),
	}
}
