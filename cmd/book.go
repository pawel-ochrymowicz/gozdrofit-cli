package cmd

import (
	"errors"
	gozdrofitapi "github.com/butwhoareyou/gozdrofit-api"
	mock "github.com/butwhoareyou/gozdrofit-cli/mock"
	log "github.com/go-pkgz/lgr"
	"net/url"
	"strings"
	"time"
)

type BookCommand struct {
	Class struct {
		ClubId   int64  `short:"c" long:"club-id" env:"CLUB_ID" description:"class club id" default:"-1"`
		Name     string `short:"n" long:"name" env:"NAME" description:"name of the class to book"`
		Date     string `short:"d" long:"date" env:"DATE" description:"class date like 2021-01-01"`
		HourFrom int    `short:"f" long:"hour-from" env:"HOUR_FROM" description:"from which hour book desired class"`
		HourTo   int    `short:"t" long:"hour-to" env:"HOUR_TO" description:"to which hour book desired class"`
	} `group:"class" namespace:"class" env-namespace:"CLASS"`
	CommonOpts
}

// gozdrofit-cli book --url=abc --username=user --password=pass --class.club-id=-1 --class.name=tumba --class.date=2021-01-01 --class.hour-from=18 --class.hour-to=20 --dry-run
func (cmd *BookCommand) Execute(_ []string) error {
	log.Print("Running booking command..")

	resetEnv("USERNAME", "PASSWORD")

	baseUrl, err := url.Parse(cmd.CommonOpts.BaseUrl)
	if err != nil {
		return err
	}

	api := gozdrofitapi.NewHttpApi(*baseUrl, cmd.HttpClient, cmd.Debug)
	if cmd.DryRun {
		api = mock.NewLogOpApi()
	}

	log.Print("Authenticating..")
	authenticated, err := api.Authenticate(gozdrofitapi.LoginRequest{
		RememberMe: true,
		Login:      cmd.Username,
		Password:   cmd.Password,
	})

	if err != nil {
		return err
	}

	log.Print("Authenticated.")

	clubId := cmd.Class.ClubId
	if clubId == -1 {
		clubId = authenticated.User.Member.DefaultClubId
	}

	date, err := parseISODate(cmd.Class.Date)

	if err != nil {
		return err
	}

	log.Print("Fetching daily classes..")
	dailyClasses, err := api.DailyClasses(gozdrofitapi.DailyClassesRequest{
		ClubId: clubId,
		Date:   gozdrofitapi.Date{Time: date},
	})

	if err != nil {
		return err
	}

	log.Print("Daily classes fetched.")

	classId := int64(-1)
	for _, calendarData := range dailyClasses.CalendarData {
		for _, class := range calendarData.Classes {
			if cmd.Class.HourFrom <= class.StartTime.Hour() &&
				cmd.Class.HourTo >= class.StartTime.Hour() &&
				strings.EqualFold(cmd.Class.Name, class.Name) {
				classId = class.Id

				log.Printf("Found %v class id %v.", cmd.Class.Name, class.Id)

				for _, user := range class.Users {
					if user.IsCurrentUser {
						log.Print("Class is already booked, exiting..")
						return nil
					}
				}

				if class.BookingIndicator.Available == 0 {
					return errors.New("the class is not available")
				}
			}
		}
	}

	if classId == -1 {
		return errors.New("specified class can not be found")
	}

	log.Printf("Booking class %v..", classId)

	err = api.BookClass(gozdrofitapi.BookClassRequest{
		ClassId: classId,
	})

	if err != nil {
		return err
	}

	log.Printf("Class %v booked.", classId)
	return nil
}

func parseISODate(dateString string) (time.Time, error) {
	return time.Parse(gozdrofitapi.DateFormat, dateString)
}
