[![Build Status](https://github.com/go-pkgz/auth/workflows/build/badge.svg)](https://github.com/butwhoare-you/gozdrofit-cli/actions)
[![Coverage Status](https://coveralls.io/repos/github/butwhoareyou/gozdrofit-cli/badge.svg?branch=master)](https://coveralls.io/github/butwhoareyou/gozdrofit-cli?branch=master)

Unofficial command line tool for Zdrofit Polska API.

---

## Usage

```
gozdrofit-cli book --url=https://example.com/ --username=username@mail.com --password=password --class.club-id=1 --class.name=Sztangi --class.date=2021-08-15 --class.hour-from=16 --class.hour-to=20 --debug
```
This command will book a class name `Sztangi` in club id `1` at 2021-08-15, and it will perform a lookup for a class with start time between 16 and 20 (4PM and 8PM).
