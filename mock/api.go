package mock

import (
	gozdrofitapi "github.com/butwhoareyou/gozdrofit-api"
	log "github.com/go-pkgz/lgr"
	"time"
)

type logOpApi struct {
	authenticatedAt time.Time
}

func NewLogOpApi() gozdrofitapi.Api {
	return &logOpApi{}
}

func (api logOpApi) Authenticated() bool {
	return api.authenticatedAt.Add(time.Hour).Before(time.Now())
}

func (api logOpApi) Authenticate(request gozdrofitapi.LoginRequest) (*gozdrofitapi.LoginResponse, error) {
	log.Printf("Authenticate %v", request)

	api.authenticatedAt = time.Now()

	return &gozdrofitapi.LoginResponse{
		User: gozdrofitapi.User{
			Member: gozdrofitapi.Member{
				Id:            99,
				HomeClubId:    99,
				DefaultClubId: 99,
			}}}, nil
}

func (api logOpApi) DailyClasses(request gozdrofitapi.DailyClassesRequest) (*gozdrofitapi.DailyClassesResponse, error) {
	log.Printf("DailyClasses %v", request)

	return &gozdrofitapi.DailyClassesResponse{
		CalendarData: []gozdrofitapi.CalendarData{
			{
				Classes: []gozdrofitapi.Class{
					{
						Id:               1,
						Status:           gozdrofitapi.ClassStatusBookable,
						Name:             "TBC",
						StartTime:        gozdrofitapi.DateTime{Time: time.Date(2021, 8, 12, 17, 0, 0, 0, time.UTC)},
						BookingIndicator: gozdrofitapi.BookingIndicator{Limit: 16, Available: 7},
					}}},
			{
				Classes: []gozdrofitapi.Class{
					{
						Id:               2,
						Status:           gozdrofitapi.ClassStatusAwaitable,
						Name:             "Tabata",
						StartTime:        gozdrofitapi.DateTime{Time: time.Date(2021, 8, 12, 18, 0, 0, 0, time.UTC)},
						BookingIndicator: gozdrofitapi.BookingIndicator{Limit: 16, Available: 1},
					}}},
		}}, nil
}

func (api logOpApi) BookClass(request gozdrofitapi.BookClassRequest) error {
	log.Printf("BookClass %v", request)

	return nil
}

func (api logOpApi) CancelClassBooking(request gozdrofitapi.CancelBookingRequest) error {
	log.Printf("CancelClassBooking %v", request)

	return nil
}
