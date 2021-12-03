package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adrialopezbou/bookings-go/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"mr", "/make-reservation", "GET", http.StatusOK},

	/* {"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "01-01-2020"},
		{key: "end", value: "02-01-2020"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "01-01-2020"},
		{key: "end", value: "02-01-2020"},
	}, http.StatusOK},
	{"make-reservation-post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "at@co.org"},
		{key: "phone", value: "555-666-7777"},
	}, http.StatusOK}, */
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} 
	}

}


func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID: 1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case with non-existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	// testing correct post body
	reqBody := "start_date=01-01-2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02-01-2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=adria")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=lopez")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=adria@lopez.es")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=66582")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	layout := "02-1-2006"
	sd, _ := time.Parse(layout, "01-01-2050")
	ed, _ := time.Parse(layout, "02-01-2050")
	sessionalRes := models.Reservation{
		StartDate: sd,
		EndDate: ed,
		RoomID: 1,
	}

	session.Put(ctx, "reservation", sessionalRes)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", sessionalRes)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for error getting reservation from session
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "foo", nil)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}


	// testing for invalid data
	reqBody = "start_date=01-01-2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02-01-2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=a")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=lopez")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=adria@lopez.es")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=66582")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _= http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	layout = "02-1-2006"
	sd, _ = time.Parse(layout, "01-01-2050")
	ed, _ = time.Parse(layout, "02-01-2050")
	sessionalRes = models.Reservation{
		StartDate: sd,
		EndDate: ed,
	}

	session.Put(ctx, "reservation", sessionalRes)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// testing for failure to insert reservation into database
	reqBody = "start_date=01-01-2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02-01-2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=adria")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=lopez")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=adria@lopez.es")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=66582")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _= http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	layout = "02-1-2006"
	sd, _ = time.Parse(layout, "01-01-2050")
	ed, _ = time.Parse(layout, "02-01-2050")
	sessionalRes = models.Reservation{
		StartDate: sd,
		EndDate: ed,
	}

	sessionalRes.RoomID = 2
	session.Put(ctx, "reservation", sessionalRes)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("inserting reservation to database didnt fail as it should : got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for failure to insert room restriction into database
	reqBody = "start_date=01-01-2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02-01-2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=adria")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=lopez")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=adria@lopez.es")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=66582")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _= http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	layout = "02-1-2006"
	sd, _ = time.Parse(layout, "01-01-2050")
	ed, _ = time.Parse(layout, "02-01-2050")
	sessionalRes = models.Reservation{
		StartDate: sd,
		EndDate: ed,
	}

	sessionalRes.RoomID = 1000
	session.Put(ctx, "reservation", sessionalRes)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("inserting reservation to database didnt fail as it should : got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// tests when rooms are not available
	reqBody := "start=01-01-2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=01-02-2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	// test for error parsing form
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.Ok {
		t.Error("empty body on post didn't return an error when parsing form")
	}

	// test for error parsing form
	reqBody = "start=01-01-2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=01-02-2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.Ok {
		t.Error("didn't return an error when failing inserting into database")
	}

}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
