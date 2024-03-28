package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/msaufi2325/06_bookings/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

// go test -coverprofile=coverage.out && go tool cover -html=coverage.out

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
	// {"post-search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"make-reservation-post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "js@email.com"},
	// 	{key: "phone", value: "555-555-5555"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			res, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusOK)
	}

	// Test case where reservation is not in session (reset everything)
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-exixtent room
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody += "&end_date=2050-01-02"
	reqBody += "&first_name=John"
	reqBody += "&last_name=Smith"
	reqBody += "&email=john@smith.com"
	reqBody += "&phone=555-555-5555"
	reqBody += "&room_id=1"

	req, err := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, err = http.NewRequest("POST", "/make-reservation", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date
	reqBody = "start_date=invalid"
	reqBody += "&end_date=2050-01-02"
	reqBody += "&first_name=John"
	reqBody += "&last_name=Smith"
	reqBody += "&email=john@smith.com"
	reqBody += "&phone=555-555-5555"
	reqBody += "&room_id=1"

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date
	reqBody = "start_date=2050-01-01"
	reqBody += "&end_date=invalid"
	reqBody += "&first_name=John"
	reqBody += "&last_name=Smith"
	reqBody += "&email=john@smith.com"
	reqBody += "&phone=555-555-5555"
	reqBody += "&room_id=1"

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid room id
	reqBody = "start_date=2050-01-01"
	reqBody += "&end_date=2050-01-02"
	reqBody += "&first_name=John"
	reqBody += "&last_name=Smith"
	reqBody += "&email=john@smith.com"
	reqBody += "&phone=555-555-5555"
	reqBody += "&room_id=invalid"

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid room id: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid data
	reqBody = "start_date=2050-01-01"
	reqBody += "&end_date=2050-01-02"
	reqBody += "&first_name=J"
	reqBody += "&last_name=Smith"
	reqBody += "&email=john@smith.com"
	reqBody += "&phone=555-555-5555"
	reqBody += "&room_id=1"

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, want %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert reservation into database
	reqBody = "start_date=2050-01-01"
	reqBody += "&end_date=2050-01-02"
	reqBody += "&first_name=John"
	reqBody += "&last_name=Smith"
	reqBody += "&email=john@smith.com"
	reqBody += "&phone=555-555-5555"
	reqBody += "&room_id=2"

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for failure to insert reservation into database: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure to insert room restriction into database
	reqBody = "start_date=2050-01-01"
	reqBody += "&end_date=2050-01-02"
	reqBody += "&first_name=John"
	reqBody += "&last_name=Smith"
	reqBody += "&email=john@smith.com"
	reqBody += "&phone=555-555-5555"
	reqBody += "&room_id=1000"

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for failure to insert room restriction into database: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRespository_AvailabilityJSON(t *testing.T) {
	// first case: rooms are not available
	reqBody := "start=2050-01-01"
	reqBody += "&end=2050-01-02"
	reqBody += "&room_id=1"

	// create request
	req, err := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// make handler handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// get response recorder
	rr := httptest.NewRecorder()

	// make request to our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
