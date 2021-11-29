package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if form.Has("anything") {
		t.Error("returns true when field does not exist")
	}

	postedData := url.Values{}
	postedData.Add("something", "something")

	r.PostForm = postedData
	form = New(r.PostForm)
	
	if !form.Has("something") {
		t.Error("returns false when field does exist")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	
	postedData := url.Values{}
	postedData.Add("something", "a")

	r.PostForm = postedData
	form := New(r.PostForm)

	form.MinLength("something", 3)

	if form.Valid() {
		t.Error("shows valid when string has less length than configured")
	}

	isError := form.Errors.Get("something")
	if isError == "" {
		t.Error("should have an error, but didn't have any")
	}

	postedData.Set("something", "aaaa")

	r.PostForm = postedData
	form = New(r.PostForm)

	form.MinLength("something", 3)

	if !form.Valid() {
		t.Error("returns false when string has more length than configured")
	}

	isError = form.Errors.Get("something")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

	form = New(r.PostForm)

	form.MinLength("some", 5)

	if form.Valid() {
		t.Error("returns true when field doesn't exist")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	
	postedData := url.Values{}
	postedData.Add("email", "ab@cd")

	r.PostForm = postedData
	form := New(r.PostForm)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("shows valid when email has wrong syntax")
	}

	postedData.Set("email", "ab@cd.org")

	r.PostForm = postedData
	form = New(r.PostForm)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("show invalid when email syntax is correct")
	}
}