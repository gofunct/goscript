package script

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"
	"io"
	"strings"
	"text/template"
	"unicode"
)

var TemplateFuncs = funcMap()

var legacyFuncs = template.FuncMap{
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"rpad":                    rpad,
}

// tmpl executes the given template text on data, writing the result to w.
func Compile(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(TemplateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

func funcMap() template.FuncMap {
	newMap := sprig.GenericFuncMap()
	for k, v := range legacyFuncs {
		newMap[k] = v
	}
	for k, v := range viper.AllSettings() {
		newMap[k] = v
	}
	return newMap
}

func addTemplateFunc(name string, tmplFunc interface{}) {
	TemplateFuncs[name] = tmplFunc
}

func addTemplateFuncs(tmplFuncs template.FuncMap) {
	for k, v := range tmplFuncs {
		TemplateFuncs[k] = v
	}
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func rpad(s string, padding int) string {
	tmpl := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(tmpl, s)
}

// ld compares two strings and returns the levenshtein distance between them.
func ld(s, t string, ignoreCase bool) int {
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
