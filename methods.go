package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// MethodsController implements the methods resource.
type MethodsController struct {
	*goa.Controller
}

// NewMethodsController creates a methods controller.
func NewMethodsController(service *goa.Service) *MethodsController {
	return &MethodsController{Controller: service.NewController("MethodsController")}
}

// Deploy runs the deploy action.
func (c *MethodsController) Deploy(ctx *app.DeployMethodsContext) error {
	// MethodsController_Deploy: start_implement

	// Put your logic here
	// This is a stub only at the moment... 
	// TODO: add kubectl / k2 / helm stuff 

	// MethodsController_Deploy: end_implement
	return nil
}
