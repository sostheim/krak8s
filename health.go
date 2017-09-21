package main

import (
	"krak8s/app"
	"time"

	"github.com/blang/semver"
	"github.com/goadesign/goa"
)

// HealthController implements the health resource.
type HealthController struct {
	*goa.Controller
}

// NewHealthController creates a health controller.
func NewHealthController(service *goa.Service) *HealthController {
	return &HealthController{Controller: service.NewController("HealthController")}
}

// Health runs the health action.
func (c *HealthController) Health(ctx *app.HealthHealthContext) error {
	// HealthController_Health: start_implement
	ver := "unknown"
	semVer, err := semver.Make(MajorMinorPatch + "-" + ReleaseType + "+git.sha." + GitCommit)
	if err == nil {
		ver = semVer.String()
	}
	return ctx.OK([]byte("Health OK: " + time.Now().String() + ", semVer: " + ver + "\n"))
	// HealthController_Health: end_implement
}
