[![Build Status](https://github.com/go-pkgz/auth/workflows/build/badge.svg)](https://github.com/butwhoare-you/gozdrofit-api/actions)
[![Coverage Status](https://coveralls.io/repos/github/butwhoareyou/gozdrofit-api/badge.svg?branch=master)](https://coveralls.io/github/butwhoareyou/gozdrofit-api?branch=master)

# Unofficial implementation of Zdrofit Polska API in Golang.

---

## Usage

```go
api := gozdrofitapi.NewHttpApi(*baseUrl, gozdrofitapi.NewDefaultHttpClient(), true)

authenticated, err := api.Authenticate(gozdrofitapi.LoginRequest{
    RememberMe: true,
    Login:      "Username",
    Password:   "Password",
})

if err != nil {
    return err
}

dailyClasses, err := api.DailyClasses(gozdrofitapi.DailyClassesRequest{
    ClubId: 1,
    Date:   gozdrofitapi.Date{Time: time.Now()},
})
```

## API Documentation

### Login

`POST` `/ClientPortal2/Auth/Login`

```json
{
  "RememberMe": true,
  "Login": "some@mail.to",
  "Password": "password"
}
```

```json
{
  "User": {
    "Member": {
      "Id": 111111,
      "FirstName": "Imie",
      "LastName": "Nazwisko",
      "Email": "some@mail.to",
      "PhotoUrl": null,
      "HomeClubId": 99,
      "DefaultClubId": 99,
      "Type": "ClubMember",
      "IsGuest": false,
      "IsFamilyChild": false,
      "NotificationsData": {
        "CanUpdateFromTrial": false,
        "HasInvalidContractPaymentMethod": false,
        "HasFingerprintAssigned": false,
        "RemainingDeposit": 0,
        "ContractStatus": "Ended"
      }
    },
    "Employee": null,
    "Roles": [
      "Member"
    ]
  },
  "State": "Classes"
}
```

### Daily Classes

`POST` `/ClientPortal2/Classes/ClassCalendar/DailyClasses`

```json
{
  "clubId": 99,
  "date": "2021-08-12",
  "categoryId": null,
  "timeTableId": null,
  "trainerId": null,
  "zoneId": null
}
```

```json
{
  "CalendarData": [
    {
      "Hour": "1900-01-01T17:00:00",
      "Classes": [
        {
          "Id": 1,
          "Status": "Bookable",
          "StatusReason": null,
          "Name": "TBC",
          "StartTime": "2021-08-12T17:00:00",
          "Duration": "PT45M",
          "BookingIndicator": {
            "Indicator": 3,
            "Limit": 16,
            "Available": 7
          },
          "Trainer": "IMIE NAZWISKO",
          "Users": [],
          "HasRelatives": false,
          "AllowBookSeatNumber": false,
          "IsClassAvailableOnline": false,
          "ClassRatingSummaryInfo": {
            "TimeTableId": 1,
            "RatingsCount": 2734,
            "Rating": 4.819312362838332,
            "Ranking": 4.8
          }
        }
      ]
    },
    {
      "Hour": "1900-01-01T18:00:00",
      "Classes": [
        {
          "Id": 1,
          "Status": "Awaitable",
          "StatusReason": null,
          "Name": "Tabata",
          "StartTime": "2021-08-12T18:00:00",
          "Duration": "PT45M",
          "BookingIndicator": {
            "Indicator": 0,
            "Limit": 16,
            "Available": 0
          },
          "Trainer": "IMIE NAZWISKO",
          "Users": [],
          "HasRelatives": false,
          "AllowBookSeatNumber": false,
          "IsClassAvailableOnline": false,
          "ClassRatingSummaryInfo": {
            "TimeTableId": 1,
            "RatingsCount": 1453,
            "Rating": 4.90984170681349,
            "Ranking": 4.9
          }
        }
      ]
    },
    {
      "Hour": "1900-01-01T19:00:00",
      "Classes": [
        {
          "Id": 1,
          "Status": "Awaitable",
          "StatusReason": null,
          "Name": "Trening Funkcjonalny",
          "StartTime": "2021-08-12T19:00:00",
          "Duration": "PT45M",
          "BookingIndicator": {
            "Indicator": 0,
            "Limit": 16,
            "Available": 0
          },
          "Trainer": "IMIE NAZWISKO",
          "Users": [],
          "HasRelatives": false,
          "AllowBookSeatNumber": false,
          "IsClassAvailableOnline": false,
          "ClassRatingSummaryInfo": {
            "TimeTableId": 1,
            "RatingsCount": 268,
            "Rating": 4.951499937313433,
            "Ranking": 4.9
          }
        }
      ]
    },
    {
      "Hour": "1900-01-01T20:00:00",
      "Classes": [
        {
          "Id": 1,
          "Status": "Bookable",
          "StatusReason": null,
          "Name": "Pumba®",
          "StartTime": "2021-08-12T20:00:00",
          "Duration": "PT45M",
          "BookingIndicator": {
            "Indicator": 2,
            "Limit": 16,
            "Available": 6
          },
          "Trainer": "IMIE NAZWISKO",
          "Users": [],
          "HasRelatives": false,
          "AllowBookSeatNumber": false,
          "IsClassAvailableOnline": false,
          "ClassRatingSummaryInfo": {
            "TimeTableId": 1,
            "RatingsCount": 9969,
            "Rating": 4.9046321599885554,
            "Ranking": 4.9
          }
        }
      ]
    }
  ],
  "PagerData": {
    "Days": [
      "2021-08-11",
      "2021-08-12",
      "2021-08-13",
      "2021-08-14",
      "2021-08-15"
    ],
    "NextDate": "2021-08-17",
    "PreviousDate": "2021-08-07",
    "Date": "2021-08-12",
    "CanGoForward": true,
    "CanGoBack": false,
    "QueryStartDate": "2021-08-12",
    "QueryEndDate": "2021-08-12"
  }
}
```

### BookClass

`POST` `/ClientPortal2/Classes/ClassCalendar/BookClass`

```json
{
  "classId": 1
}
```

```json
{
  "Tickets": [
    {
      "TimeTableEventId": 1,
      "Name": "Pumba®",
      "StartTime": "2021-08-12T20:00:00",
      "ZoneName": "Zdrofit Centrum Krucza",
      "UserName": "Imie Nazwisko",
      "UserNumber": "5555555",
      "UserId": 111111,
      "Trainer": "IMIE NAZWISKO"
    }
  ],
  "ClassId": 1,
  "UserId": 111111
}
```

### CancelBooking

`POST` `/ClientPortal2/Classes/ClassCalendar/CancelBooking`

```json
{
  "classId": 1
}
```

```json
{
  "ClassId": 1,
  "UserId": 111111
}
```
