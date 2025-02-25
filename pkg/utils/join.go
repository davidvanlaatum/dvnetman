package utils

import "bytes"

func JoinPtr(sep string, parts ...*string) *string {
	var valid bool
	var buf bytes.Buffer
	for _, p := range parts {
		if p != nil {
			valid = true
			if buf.Len() > 0 {
				buf.WriteString(sep)
			}
			buf.WriteString(*p)
		}
	}
	if !valid {
		return nil
	}
	return ToPtr(buf.String())
}
