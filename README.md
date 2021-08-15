[![Build Status](https://github.com/go-pkgz/auth/workflows/build/badge.svg)](https://github.com/butwhoare-you/gozdrofit-cli/actions)
[![Coverage Status](https://coveralls.io/repos/github/butwhoareyou/gozdrofit-cli/badge.svg?branch=master)](https://coveralls.io/github/butwhoareyou/gozdrofit-cli?branch=master)

# Unofficial command line tool for Zdrofit Polska API.

---

## Commands

### book

* Usage without Docker
```
gozdrofit-cli book --url=https://example.com/ --username=username@mail.com --password=password --class.club-id=1 --class.name=Sztangi --class.date=2021-08-15 --class.hour-from=16 --class.hour-to=20
```

* Usage with Docker
```
docker pull pochrymowicz/gozdrofit-cli:0.1.0
docker run pochrymowicz/gozdrofit-cli:0.1.0 book --url=https://example.com/ --username=username --password=password --class.club-id=1 --class.name=Sztangi --class.date=2021-08-15 --class.hour-from=16 --class.hour-to=20
```

This command will book a class named `Sztangi` in club id `1` at `2021-08-15`, and it will perform a lookup for a class
with start time between `16` and `20` (`4PM` and `8PM`).

* Arguments 

| Argument | Description | Example | Required | 
|---|---|---|---|
| url | Zdrofit API base url | https://example.com/ | true |
| username | Zdrofit username | login@mail.com | true |
| password | Zdrofit password | password | true |
| class.club-id | Zdrofit club id | 1 | default club id from Zdrofit profile will be used if not specified |
| class.name | Case-insensitive class name to book | Sztangi | true |
| class.date | Class date YYYY-MM-DD | 2021-01-01 | true |
| class.hour-from | Class start hour low bound (24h) (inclusive) | 16 | true |
| class.hour-to | Class start hour high bound (24h) (inclusive) | 16 | true |
| class.debug | Turn on debug logging |  |  |
| class.dry-run | Dry run |  |  |
