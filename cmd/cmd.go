package cmd

import (
	log "github.com/go-pkgz/lgr"
	"net/http"
	"os"
	"strings"
)

type CommonCommander interface {
	SetCommon(commonOpts CommonOpts)

	Execute(args []string) error
}

type CommonOpts struct {
	BaseUrl  string
	Username string
	Password string
	DryRun   bool
	Debug    bool

	HttpClient http.Client
}

func (c *CommonOpts) SetCommon(commonOpts CommonOpts) {
	c.BaseUrl = strings.TrimSuffix(commonOpts.BaseUrl, "/") // allow url with trailing /
	c.Username = commonOpts.Username
	c.Password = commonOpts.Password
	c.DryRun = commonOpts.DryRun
	c.Debug = commonOpts.Debug
	c.HttpClient = commonOpts.HttpClient
}

// resetEnv clears sensitive env vars
func resetEnv(envs ...string) {
	for _, env := range envs {
		if err := os.Unsetenv(env); err != nil {
			log.Printf("[WARN] can't unset env %s, %s", env, err)
		}
	}
}
