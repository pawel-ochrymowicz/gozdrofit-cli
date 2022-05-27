package gozdrofitapi

import (
	"bytes"
	"encoding/json"
	log "github.com/go-pkgz/lgr"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Api interface {
	Authenticated() bool
	Authenticate(request LoginRequest) (*LoginResponse, error)
	DailyClasses(request DailyClassesRequest) (*DailyClassesResponse, error)
	BookClass(request BookClassRequest) error
	CancelClassBooking(request CancelBookingRequest) error
}

type LoginRequest struct {
	RememberMe bool   `json:"RememberMe"`
	Login      string `json:"Login"`
	Password   string `json:"Password"`
}

type DailyClassesRequest struct {
	ClubId int64 `json:"clubId"`
	Date   Date  `json:"date"`
}

type BookClassRequest struct {
	ClassId int64 `json:"classId"`
}

type CancelBookingRequest struct {
	ClassId int64 `json:"classId"`
}

type LoginResponse struct {
	User User `json:"User"`
}

type User struct {
	Member Member `json:"Member"`
}

type Member struct {
	Id            int64 `json:"Id"`
	HomeClubId    int64 `json:"HomeClubId"`
	DefaultClubId int64 `json:"DefaultClubId"`
}

type DailyClassesResponse struct {
	CalendarData []CalendarData `json:"CalendarData"`
}

type CalendarData struct {
	Classes []Class `json:"Classes"`
}

const (
	ClassStatusBookable  = "Bookable"
	ClassStatusAwaitable = "Awaitable"
)

type Class struct {
	Id               int64            `json:"Id"`
	Status           string           `json:"Status"`
	StatusReason     string           `json:"StatusReason,omitempty"`
	Name             string           `json:"Name"`
	StartTime        DateTime         `json:"StartTime"`
	BookingIndicator BookingIndicator `json:"BookingIndicator"`
	Users            []ClassUser      `json:"Users"`
}

type ClassUser struct {
	Id            int64 `json:"Id"`
	IsCurrentUser bool  `json:"IsCurrentUser"`
}

type BookingIndicator struct {
	Limit     int `json:"Limit"`
	Available int `json:"Available"`
}

const (
	defaultUserAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	authenticationCookie = "ClientPortal.Auth.bak"
)

func NewDefaultHttpClient() http.Client {
	jar, _ := cookiejar.New(nil)
	return http.Client{
		Timeout: time.Second * 10,
		Jar:     jar,
	}
}

type httpApi struct {
	baseUrl    url.URL
	httpClient http.Client
}

func NewHttpApi(baseUrl url.URL, httpClient http.Client, debug bool) Api {
	setupLog(debug)

	return &httpApi{baseUrl: baseUrl, httpClient: httpClient}
}

// Authenticated checks authentication cookies alivness
func (api *httpApi) Authenticated() bool {
	cookies := api.httpClient.Jar.Cookies(&api.baseUrl)

	var authCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == authenticationCookie {
			authCookie = cookie
		}
	}
	if authCookie == nil || cookieExpired(authCookie) {
		return false
	}

	return true
}

// Authenticate exchanges login and password for authentication cookies and creates a session inside Api
func (api *httpApi) Authenticate(request LoginRequest) (*LoginResponse, error) {
	resp, err := post(&api.httpClient, api.baseUrl.String()+"/ClientPortal2/Auth/Login", request)
	if err != nil {
		return nil, err
	}

	var loginResponse LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] received authentication cookies: %v", len(resp.Cookies()))

	return &loginResponse, nil
}

// DailyClasses Fetches daily classes using possible filters
func (api *httpApi) DailyClasses(request DailyClassesRequest) (*DailyClassesResponse, error) {
	resp, err := post(&api.httpClient, api.baseUrl.String()+"/ClientPortal2/Classes/ClassCalendar/DailyClasses", request)
	if err != nil {
		return nil, err
	}

	var dailyClassesResponse DailyClassesResponse
	err = json.NewDecoder(resp.Body).Decode(&dailyClassesResponse)
	return &dailyClassesResponse, err
}

// BookClass books specified class by id
func (api *httpApi) BookClass(request BookClassRequest) error {
	_, err := post(&api.httpClient, api.baseUrl.String()+"/ClientPortal2/Classes/ClassCalendar/BookClass", request)
	return err
}

// CancelClassBooking cancels previously ackuired booking
func (api *httpApi) CancelClassBooking(request CancelBookingRequest) error {
	_, err := post(&api.httpClient, api.baseUrl.String()+"/ClientPortal2/Classes/ClassCalendar/CancelBooking", request)
	return err
}

func post(client *http.Client, url string, payload interface{}) (*http.Response, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", defaultUserAgent)

	log.Printf("[DEBUG] POST %v %v ", url, len(jsonPayload))
	resp, err := client.Do(req)

	return resp, err
}

func cookieExpired(cookie *http.Cookie) bool {
	return time.Now().After(cookie.Expires.Add(time.Minute))
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
