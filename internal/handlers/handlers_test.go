package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about/", "GET", []postData{}, http.StatusOK},
	{"presSuite", "/presidential-suites/", "GET", []postData{}, http.StatusOK},
	{"platSuite", "/platinum-rooms/", "GET", []postData{}, http.StatusOK},
	{"searchAv", "/search-availability/", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact/", "GET", []postData{}, http.StatusOK},
	{"makeReserv", "/make-reservation/", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability/", "POST", []postData{
		{key: "start", value: "2021-11-19"},
		{key: "end", value: "2021-11-30"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json/", "POST", []postData{
		{key: "start", value: "2021-11-19"},
		{key: "end", value: "2021-11-30"},
	}, http.StatusOK},
	{"make reservation post", "/make-reservation/", "POST", []postData{
		{key: "first_name", value: "Owakabi"},
		{key: "last_name", value: "Kapere"},
		{key: "email", value: "Owak@kapere.com"},
		{key: "phone", value: "0786220378"},
	}, http.StatusOK},
}

func TestNewHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes) //this creates a server for the test
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
