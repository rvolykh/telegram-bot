package scrap

import "strings"

func cleanHTML(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`
	for i, c := range text {
		if (i+1) == len(text) && end >= start {
			builder.WriteString(text[end:])
		}
		if c != '<' && c != '>' {
			continue
		}
		if c == '<' {
			if !in {
				start = i
				builder.WriteString(text[end:start])
			}
			in = true
			continue
		}
		end, in = i+1, false
	}

	return builder.String()
}
