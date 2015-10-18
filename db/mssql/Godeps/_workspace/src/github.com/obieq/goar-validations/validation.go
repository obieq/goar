package goar_validations

import (

	//"net/http"
	//"net/url"
	"fmt"
	"regexp"
)

// Simple struct to store the Message & Key of a validation error
type ValidationError struct {
	Message, Key string
}

// String returns the Message field of the ValidationError struct.
func (e *ValidationError) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// A Validation context manages data validation and error messages.
type Validation struct {
	Errors []*ValidationError `gorethink:"-" json:"-"`
	keep   bool
}

// Keep tells revel to set a flash cookie on the client to make the validation
// errors available for the next request.
// This is helpful  when redirecting the client after the validation failed.
// It is good practice to always redirect upon a HTTP POST request. Thus
// one should use this method when HTTP POST validation failed and redirect
// the user back to the form.
func (v *Validation) Keep() {
	v.keep = true
}

// Clear *all* ValidationErrors
func (v *Validation) Clear() {
	v.Errors = []*ValidationError{}
}

// HasErrors returns true if there are any (ie > 0) errors. False otherwise.
func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

// NumErrors returns the number of validation errors for a given instance
func (v *Validation) NumErrors() int {
	return len(v.Errors)
}

// ErrorMap returns the errors mapped by key.
// If there are multiple validation errors associated with a single key, the
// first one "wins".  (Typically the first validation will be the more basic).
func (v *Validation) ErrorMap() map[string]*ValidationError {
	m := map[string]*ValidationError{}
	for _, e := range v.Errors {
		if _, ok := m[e.Key]; !ok {
			m[e.Key] = e
		}
	}
	return m
}

// Error adds an error to the validation context.
func (v *Validation) Error(key string, message string, args ...interface{}) *ValidationResult {
	result := (&ValidationResult{
		Ok:    false,
		Error: &ValidationError{Key: key},
	}).Message(message, args...)
	v.Errors = append(v.Errors, result.Error)
	return result
}

// A ValidationResult is returned from every validation method.
// It provides an indication of success, and a pointer to the Error (if any).
type ValidationResult struct {
	Error *ValidationError
	Ok    bool
}

// Key sets the ValidationResult's Error "key" and returns itself for chaining
//func (r *ValidationResult) Key(key string) *ValidationResult {
//log.Fatal("Key: ", key)
//if r.Error != nil {
//r.Error.Key = key
//}
//return r
//}

// Message sets the error message for a ValidationResult. Returns itself to
// allow chaining.  Allows Sprintf() type calling with multiple parameters
func (r *ValidationResult) Message(message string, args ...interface{}) *ValidationResult {
	if r.Error != nil {
		if len(args) == 0 {
			r.Error.Message = message
		} else {
			r.Error.Message = fmt.Sprintf(message, args...)
		}
	}
	return r
}

// Required tests that the argument is non-nil and non-empty (if string or list)
func (v *Validation) Required(key string, obj interface{}) *ValidationResult {
	return v.apply(Required{}, obj, key)
}

func (v *Validation) Min(key string, n int, min int) *ValidationResult {
	return v.apply(Min{Min: min}, n, key)
}

func (v *Validation) Max(key string, n int, max int) *ValidationResult {
	return v.apply(Max{Max: max}, n, key)
}

func (v *Validation) Range(key string, n, min, max int) *ValidationResult {
	return v.apply(Range{Min: Min{Min: min}, Max: Max{Max: max}}, n, key)
}

func (v *Validation) MinSize(key string, obj interface{}, min int) *ValidationResult {
	return v.apply(MinSize{Min: min}, obj, key)
}

func (v *Validation) MaxSize(key string, obj interface{}, max int) *ValidationResult {
	return v.apply(MaxSize{Max: max}, obj, key)
}

func (v *Validation) Length(key string, obj interface{}, n int) *ValidationResult {
	return v.apply(Length{N: n}, obj, key)
}

func (v *Validation) Match(key string, str string, regex *regexp.Regexp) *ValidationResult {
	return v.apply(Match{Regexp: regex}, str, key)
}

func (v *Validation) Email(key string, str string) *ValidationResult {
	return v.apply(Email{Match: Match{Regexp: emailPattern}}, str, key)
}

func (v *Validation) apply(chk Validator, obj interface{}, key string) *ValidationResult {
	if chk.IsSatisfied(obj) {
		return &ValidationResult{Ok: true}
	}

	// Add the error to the validation context.
	err := &ValidationError{
		Message: chk.DefaultMessage(),
		Key:     key,
	}
	v.Errors = append(v.Errors, err)

	// Also return it in the result.
	return &ValidationResult{
		Ok:    false,
		Error: err,
	}
}

// Apply a group of validators to a field, in order, and return the
// ValidationResult from the first one that fails, or the last one that
// succeeds.
func (v *Validation) Check(obj interface{}, checks ...Validator) *ValidationResult {
	var result *ValidationResult
	for _, check := range checks {
		result = v.apply(check, obj, check.GetKey())
		if !result.Ok {
			return result
		}
	}
	return result
}

// Revel Filter function to be hooked into the filter chain.
//func ValidationFilter(c *Controller, fc []Filter) {
//errors, err := restoreValidationErrors(c.Request.Request)
//c.Validation = &Validation{
//Errors: errors,
//keep:   false,
//}
//hasCookie := (err != http.ErrNoCookie)

//fc[0](c, fc[1:])

//// Add Validation errors to RenderArgs.
//c.RenderArgs["errors"] = c.Validation.ErrorMap()

//// Store the Validation errors
//var errorsValue string
//if c.Validation.keep {
//for _, error := range c.Validation.Errors {
//if error.Message != "" {
//errorsValue += "\x00" + error.Key + ":" + error.Message + "\x00"
//}
//}
//}

//// When there are errors from Validation and Keep() has been called, store the
//// values in a cookie. If there previously was a cookie but no errors, remove
//// the cookie.
//if errorsValue != "" {
//c.SetCookie(&http.Cookie{
//Name:     CookiePrefix + "_ERRORS",
//Value:    url.QueryEscape(errorsValue),
//Path:     "/",
//HttpOnly: CookieHttpOnly,
//Secure:   CookieSecure,
//})
//} else if hasCookie {
//c.SetCookie(&http.Cookie{
//Name:     CookiePrefix + "_ERRORS",
//MaxAge:   -1,
//Path:     "/",
//HttpOnly: CookieHttpOnly,
//Secure:   CookieSecure,
//})
//}
//}

//// Restore Validation.Errors from a request.
//func restoreValidationErrors(req *http.Request) ([]*ValidationError, error) {
//var (
//err    error
//cookie *http.Cookie
//errors = make([]*ValidationError, 0, 5)
//)
//if cookie, err = req.Cookie(CookiePrefix + "_ERRORS"); err == nil {
//ParseKeyValueCookie(cookie.Value, func(key, val string) {
//errors = append(errors, &ValidationError{
//Key:     key,
//Message: val,
//})
//})
//}
//return errors, err
//}

// Register default validation keys for all calls to Controller.Validation.Func().
// Map from (package).func => (line => name of first arg to Validation func)
// E.g. "myapp/controllers.helper" or "myapp/controllers.(*Application).Action"
// This is set on initialization in the generated main.go file.
var DefaultValidationKeys map[string]map[int]string
