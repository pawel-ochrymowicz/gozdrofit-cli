package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umputun/go-flags"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

const authenticationCookie = "ClientPortal.Auth.bak"

func TestBookCommand_DryRun(t *testing.T) {
	cmd := BookCommand{}
	cmd.SetCommon(CommonOpts{"", "any-user", "any-password", true, true, newHttpClient()})
	p := flags.NewParser(&cmd, flags.Default)
	_, err := p.ParseArgs([]string{"--class.club-id=-1", "--class.name=TBC", "--class.date=2021-08-12", "--class.hour-from=17", "--class.hour-to=17"})
	require.NoError(t, err)
	err = cmd.Execute(nil)

	assert.NoError(t, err)
}

func TestBookCommand_Execute(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ClientPortal2/Auth/Login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		http.SetCookie(w, &http.Cookie{
			Name:    authenticationCookie,
			Value:   "token",
			Expires: time.Now().Add(time.Hour * 24),
		})
		_, _ = fmt.Fprintf(w, "{\"User\":{\"Member\":{\"Id\":99,\"HomeClubId\":99,\"DefaultClubId\":99}}}")
	})
	mux.HandleFunc("/ClientPortal2/Classes/ClassCalendar/DailyClasses", func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, "{\"CalendarData\":["+
			"{\"Classes\":[{\"Id\":1,\"Status\":\"Bookable\",\"StatusReason\":null,\"Name\":\"TBC\",\"StartTime\":\"2021-08-12T17:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":7},\"Users\":[{\"Id\":1,\"IsCurrentUser\":true}]}]},"+
			"{\"Classes\":[{\"Id\":2,\"Status\":\"Awaitable\",\"StatusReason\":null,\"Name\":\"Pumba\",\"StartTime\":\"2021-08-12T18:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":1},\"Users\":[]}]}"+
			"]}")
	})
	mux.HandleFunc("ClientPortal2/Classes/ClassCalendar/BookClass", func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, "{\"Tickets\":[{\"TimeTableEventId\":1,\"Name\":\"Pumba\",\"StartTime\":\"2021-08-12T20:00:00\",\"ZoneName\":\"Zdrofit Centrum Krucza\",\"UserName\":\"Imie Nazwisko\",\"UserNumber\":\"5555555\",\"UserId\":99,\"Trainer\":\"IMIE NAZWISKO\"}],\"ClassId\":1,\"UserId\":99}")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := BookCommand{}
	cmd.SetCommon(CommonOpts{server.URL, "any-user", "any-password", false, true, newHttpClient()})
	p := flags.NewParser(&cmd, flags.Default)
	_, err := p.ParseArgs([]string{"--class.club-id=-1", "--class.name=Pumba", "--class.date=2021-08-12", "--class.hour-from=18", "--class.hour-to=20"})
	require.NoError(t, err)
	err = cmd.Execute(nil)

	assert.NoError(t, err)
}

func TestBookCommand_Execute_ClassAlreadyBooked(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ClientPortal2/Auth/Login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		http.SetCookie(w, &http.Cookie{
			Name:    authenticationCookie,
			Value:   "token",
			Expires: time.Now().Add(time.Hour * 24),
		})
		_, _ = fmt.Fprintf(w, "{\"User\":{\"Member\":{\"Id\":99,\"HomeClubId\":99,\"DefaultClubId\":99}}}")
	})
	mux.HandleFunc("/ClientPortal2/Classes/ClassCalendar/DailyClasses", func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, "{\"CalendarData\":["+
			"{\"Classes\":[{\"Id\":1,\"Status\":\"Bookable\",\"StatusReason\":null,\"Name\":\"TBC\",\"StartTime\":\"2021-08-12T17:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":7},\"Users\":[{\"Id\":1,\"IsCurrentUser\":true}]}]},"+
			"{\"Classes\":[{\"Id\":2,\"Status\":\"Awaitable\",\"StatusReason\":null,\"Name\":\"Pumba\",\"StartTime\":\"2021-08-12T18:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":1},\"Users\":[]}]}"+
			"]}")
	})
	mux.HandleFunc("ClientPortal2/Classes/ClassCalendar/BookClass", func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, "{\"Tickets\":[{\"TimeTableEventId\":1,\"Name\":\"TBC\",\"StartTime\":\"2021-08-12T20:00:00\",\"ZoneName\":\"Zdrofit Centrum Krucza\",\"UserName\":\"Imie Nazwisko\",\"UserNumber\":\"5555555\",\"UserId\":99,\"Trainer\":\"IMIE NAZWISKO\"}],\"ClassId\":1,\"UserId\":99}")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := BookCommand{}
	cmd.SetCommon(CommonOpts{server.URL, "any-user", "any-password", false, true, newHttpClient()})
	p := flags.NewParser(&cmd, flags.Default)
	_, err := p.ParseArgs([]string{"--class.club-id=-1", "--class.name=TBC", "--class.date=2021-08-12", "--class.hour-from=17", "--class.hour-to=17"})
	require.NoError(t, err)
	err = cmd.Execute(nil)

	assert.NoError(t, err)
}

func TestBookCommand_Execute_Failed_ClassIdNotFound_NoClasses(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ClientPortal2/Auth/Login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		http.SetCookie(w, &http.Cookie{
			Name:    authenticationCookie,
			Value:   "token",
			Expires: time.Now().Add(time.Hour * 24),
		})
		_, _ = fmt.Fprintf(w, "{\"User\":{\"Member\":{\"Id\":99,\"HomeClubId\":99,\"DefaultClubId\":99}}}")
	})
	mux.HandleFunc("/ClientPortal2/Classes/ClassCalendar/DailyClasses", func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, "{\"CalendarData\":[]}")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := BookCommand{}
	cmd.SetCommon(CommonOpts{server.URL, "any-user", "any-password", false, true, newHttpClient()})
	p := flags.NewParser(&cmd, flags.Default)
	_, err := p.ParseArgs([]string{"--class.club-id=-1", "--class.name=Pumba", "--class.date=2021-08-12", "--class.hour-from=18", "--class.hour-to=20"})
	require.NoError(t, err)
	err = cmd.Execute(nil)

	assert.EqualError(t, err, "specified class can not be found")
}

func TestBookCommand_Execute_Failed_ClassIsNotAvailable_StartHourOutsideSpecifiedBrackets(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ClientPortal2/Auth/Login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		http.SetCookie(w, &http.Cookie{
			Name:    authenticationCookie,
			Value:   "token",
			Expires: time.Now().Add(time.Hour * 24),
		})
		_, _ = fmt.Fprintf(w, "{\"User\":{\"Member\":{\"Id\":99,\"HomeClubId\":99,\"DefaultClubId\":99}}}")
	})
	mux.HandleFunc("/ClientPortal2/Classes/ClassCalendar/DailyClasses", func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, "{\"CalendarData\":["+
			"{\"Classes\":[{\"Id\":1,\"Status\":\"Bookable\",\"StatusReason\":null,\"Name\":\"TBC\",\"StartTime\":\"2021-08-12T17:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":7},\"Users\":[{\"Id\":1,\"IsCurrentUser\":true}]}]},"+
			"{\"Classes\":[{\"Id\":2,\"Status\":\"Awaitable\",\"StatusReason\":null,\"Name\":\"Pumba\",\"StartTime\":\"2021-08-12T17:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":0},\"Users\":[]}]}"+
			"]}")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := BookCommand{}
	cmd.SetCommon(CommonOpts{server.URL, "any-user", "any-password", false, true, newHttpClient()})
	p := flags.NewParser(&cmd, flags.Default)
	_, err := p.ParseArgs([]string{"--class.club-id=-1", "--class.name=Pumba", "--class.date=2021-08-12", "--class.hour-from=18", "--class.hour-to=20"})
	require.NoError(t, err)
	err = cmd.Execute(nil)

	assert.EqualError(t, err, "specified class can not be found")
}

func TestBookCommand_Execute_Failed_ClassIsNotAvailable(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ClientPortal2/Auth/Login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		http.SetCookie(w, &http.Cookie{
			Name:    authenticationCookie,
			Value:   "token",
			Expires: time.Now().Add(time.Hour * 24),
		})
		_, _ = fmt.Fprintf(w, "{\"User\":{\"Member\":{\"Id\":99,\"HomeClubId\":99,\"DefaultClubId\":99}}}")
	})
	mux.HandleFunc("/ClientPortal2/Classes/ClassCalendar/DailyClasses", func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, "{\"CalendarData\":["+
			"{\"Classes\":[{\"Id\":1,\"Status\":\"Bookable\",\"StatusReason\":null,\"Name\":\"TBC\",\"StartTime\":\"2021-08-12T17:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":7},\"Users\":[{\"Id\":1,\"IsCurrentUser\":true}]}]},"+
			"{\"Classes\":[{\"Id\":2,\"Status\":\"Awaitable\",\"StatusReason\":null,\"Name\":\"Pumba\",\"StartTime\":\"2021-08-12T18:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":0},\"Users\":[]}]}"+
			"]}")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cmd := BookCommand{}
	cmd.SetCommon(CommonOpts{server.URL, "any-user", "any-password", false, true, newHttpClient()})
	p := flags.NewParser(&cmd, flags.Default)
	_, err := p.ParseArgs([]string{"--class.club-id=-1", "--class.name=Pumba", "--class.date=2021-08-12", "--class.hour-from=18", "--class.hour-to=20"})
	require.NoError(t, err)
	err = cmd.Execute(nil)

	assert.EqualError(t, err, "the class is not available")
}

func expectCookie(t *testing.T, name string, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			return
		}
	}
	t.Errorf("Cookie %v is expected but wasn't provided in the request", name)
}

type TestJar struct {
	m      sync.Mutex
	perURL map[string][]*http.Cookie
}

func (j *TestJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.m.Lock()
	defer j.m.Unlock()
	if j.perURL == nil {
		j.perURL = make(map[string][]*http.Cookie)
	}
	j.perURL[u.Host] = cookies
}

func (j *TestJar) Cookies(u *url.URL) []*http.Cookie {
	j.m.Lock()
	defer j.m.Unlock()
	return j.perURL[u.Host]
}

func newHttpClient() http.Client {
	return http.Client{Jar: &TestJar{}}
}
