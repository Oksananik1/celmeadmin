package helper

// Package httpflags provides a convenient way of filling struct fields from
// http request form values. Exposed struct fields should have special `flag`
// tag attached:
//
//	args := struct {
//		Name    string `flag:"name"`
//		Age     uint   `flag:"age"`
//		Married bool   // this won't be exposed
//	}{
// 		// default values
// 		Name: "John Doe",
// 		Age:  34,
//	}
//
// After declaring flags and their default values as above, call
// httpflags.Parse() inside http.Handler to fill struct fields from http request
// form values:
//
// 	func myHandler(w http.ResponseWriter, r *http.Request) {
//		args := struct {
// 			...
// 		}{}
// 		if err := httpflags.Parse(&args, r) ; err != nil {
//			http.Error(w, "Bad request", http.StatusBadRequest)
//			return
// 		}
//		// use args fields here
//
// Parse() calls ParseForm method of http.Request automatically, so it
// understands parsed values both from the URL field's query parameters and the
// POST or PUT form data.
//
// Package httpflags supports all basic types supported by xxxVar functions from
// standard library flag package: int, int64, uint, uint64, float64, bool,
// string, time.Duration as well as types implementing flag.Value interface.
// Parse panics on non-empty `flag` tag on unsupported type field.

import (
	"flag"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/artyom/autoflags"
)

// Parse fills dst struct with values extracted from r.Form. dst should be
// a non-nil pointer to struct having its exported attributes tagged with 'flag'
// tag — see autoflags package documentation. r.ParseForm is called
// automatically. Only the first value of each key from r.Form is used.
func Parse(dst interface{}, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	fs := new(flag.FlagSet)
	fs.SetOutput(ioutil.Discard)
	autoflags.DefineFlagSet(fs, dst)
	// this explicitly skips some sanity checks that are already done by
	// autoflags.DefineFlagSet
	st := reflect.Indirect(reflect.ValueOf(dst))
	args := make([]string, 0, st.NumField())
	for i := 0; i < st.NumField(); i++ {
		tag := st.Type().Field(i).Tag.Get("flag")
		if tag == "" {
			continue
		}
		key := tag
		if idx := strings.IndexRune(tag, ','); idx != -1 {
			key = tag[:idx]
		}
		for _, val := range r.Form[key] {
			if val != "" {
				args = append(args, "-"+key+"="+val)
			}
		}
	}
	return fs.Parse(args)
}
