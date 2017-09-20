package main

import (
	"krak8s/app"
	"time"

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
	return ctx.OK([]byte("Health OK: " + time.Now().String() + "\n"))
	// HealthController_Health: end_implement
}
