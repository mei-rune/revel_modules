package controllers

import (
	"strings"

	p2 "github.com/runner-mei/pongo2"

	"github.com/revel/revel"
	"github.com/robfig/cron"
	"github.com/runner-mei/revel_modules/jobs/app/jobs"
)

type Jobs struct {
	*revel.Controller
}

func (c Jobs) Status() revel.Result {
	remoteAddress := c.Request.RemoteAddr
	if revel.Config.BoolDefault("jobs.acceptproxyaddress", false) {
		if proxiedAddress, isProxied := c.Request.Header["X-Forwarded-For"]; isProxied {
			remoteAddress = proxiedAddress[0]
		}
	}
	if !strings.HasPrefix(remoteAddress, "127.0.0.1") && !strings.HasPrefix(remoteAddress, "::1") {
		return c.Forbidden("%s is not local", remoteAddress)
	}
	entries := jobs.MainCron.Entries()
	return c.Render(entries)
}

func init() {
	revel.TemplateFuncs["castjob"] = func(job cron.Job) *jobs.Job {
		return job.(*jobs.Job)
	}
	p2.RegisterFilter("castjob", func(ctx *p2.ExecutionContext, in *p2.Value, param *p2.Value) (out *p2.Value, err *p2.Error) {
		job := in.Interface().(*jobs.Job)
		return p2.AsValue(job), nil
	})
}
