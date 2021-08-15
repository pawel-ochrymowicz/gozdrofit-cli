package main

import (
	gozdrofitapi "github.com/butwhoareyou/gozdrofit-api"
	"github.com/butwhoareyou/gozdrofit-cli/cmd"
	log "github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type Opts struct {
	BookCmd cmd.BookCommand `command:"book"`

	ZdrofitUrl      string `long:"url" env:"ZDROFIT_URL" required:"true" description:"url to zdrofit"`
	ZdrofitUsername string `long:"username" env:"ZDROFIT_USERNAME" required:"true" description:"registered zdrofit user name"`
	ZdrofitPassword string `long:"password" env:"ZDROFIT_PASSWORD" required:"true" description:"registered zdrofit user password"`

	DryRun bool `long:"dry-run" env:"DRY_RUN" description:"dry run mode"`
	Debug  bool `long:"debug" env:"DEBUG" description:"debug mode"`
}

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		c := command.(cmd.CommonCommander)
		c.SetCommon(cmd.CommonOpts{
			BaseUrl:    opts.ZdrofitUrl,
			Username:   opts.ZdrofitUsername,
			Password:   opts.ZdrofitPassword,
			DryRun:     opts.DryRun,
			Debug:      opts.Debug,
			HttpClient: gozdrofitapi.NewDefaultHttpClient(),
		})
		err := c.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed with %+v", err)
		}
		return err
	}
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

// getDump reads runtime stack and returns as a string
func getDump() string {
	maxSize := 5 * 1024 * 1024
	stacktrace := make([]byte, maxSize)
	length := runtime.Stack(stacktrace, true)
	if length > maxSize {
		length = maxSize
	}
	return string(stacktrace[:length])
}

// nolint:gochecknoinits // can't avoid it in this place
func init() {
	// catch SIGQUIT and print stack traces
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			log.Printf("[INFO] SIGQUIT detected, dump:\n%s", getDump())
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
}
