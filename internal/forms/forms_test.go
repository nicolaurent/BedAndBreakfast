package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should hanve been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if len(form.Errors) == 0 {
		t.Error("form shows no error when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if len(form.Errors) != 0 {
		t.Error("form shows error when required fields are exist")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	has := form.Has("a")

	if has {
		t.Error("form Has() returns true when field does not exist")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")

	if !has {
		t.Error("form Has() returns false when field exists")
	}
}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	postData.Add("length2", "aa")
	postData.Add("length3", "aaa")
	postData.Add("length4", "aaaa")

	requiredLength := 3

	form := New(postData)

	meetMinLength := form.MinLength("length2", requiredLength)
	if meetMinLength {
		t.Error("MinLength returns true when actual length is shorter")
	}

	isError := form.Errors.Get("length2")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	meetMinLength = form.MinLength("length3", requiredLength)
	if !meetMinLength {
		t.Error("MinLength returns false when actual length is met")
	}

	isError = form.Errors.Get("length3")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

	meetMinLength = form.MinLength("length4", requiredLength)
	if !meetMinLength {
		t.Error("MinLength returns false when actual length is longer")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	postData.Add("invalidEmail", "one@two")
	postData.Add("validEmail", "one@two.com")

	form := New(postData)

	form.IsEmail("non-existent-email")
	if _, exist := form.Errors["non-existent-email"]; !exist {
		t.Error("form shows no error for non existent email")
	}

	form.IsEmail("invalidEmail")
	if _, exist := form.Errors["invalidEmail"]; !exist {
		t.Error("form shows no error when invalid email exist")
	}

	form.IsEmail("validEmail")
	if _, exist := form.Errors["validEmail"]; exist {
		t.Error("form shows error when valid email exist")
	}
}
