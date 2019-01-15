package render

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"
)

var TemplateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          TrimRightSpace,
	"trimTrailingWhitespaces": TrimRightSpace,
	"rpad":                    Rpad,
}

// AddTemplateFunc adds a template function that's available to Usage and Help
// template generation.
func AddTemplateFunc(name string, tmplFunc interface{}) {
	TemplateFuncs[name] = tmplFunc
}

// AddTemplateFuncs adds multiple template functions that are available to Usage and
// Help template generation.
func AddTemplateFuncs(tmplFuncs template.FuncMap) {
	for k, v := range tmplFuncs {
		TemplateFuncs[k] = v
	}
}

func TrimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// rpad adds padding to the right of a string.
func Rpad(s string, padding int) string {
	tmpl := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(tmpl, s)
}

// tmpl executes the given template text on data, writing the result to w.
func Tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(TemplateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

// ld compares two strings and returns the levenshtein distance between them.
func Ld(s, t string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
