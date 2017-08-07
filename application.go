package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// ApplicationController implements the application resource.
type ApplicationController struct {
	*goa.Controller
}

// NewApplicationController creates a application controller.
func NewApplicationController(service *goa.Service) *ApplicationController {
	return &ApplicationController{Controller: service.NewController("ApplicationController")}
}

// Create runs the create action.
func (c *ApplicationController) Create(ctx *app.CreateApplicationContext) error {
	// ApplicationController_Create: start_implement

	// Put your logic here

	// ApplicationController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *ApplicationController) Delete(ctx *app.DeleteApplicationContext) error {
	// ApplicationController_Delete: start_implement

	// Put your logic here

	// ApplicationController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *ApplicationController) Get(ctx *app.GetApplicationContext) error {
	// ApplicationController_Get: start_implement

	// Put your logic here

	// ApplicationController_Get: end_implement
	res := &app.App{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *ApplicationController) List(ctx *app.ListApplicationContext) error {
	// ApplicationController_List: start_implement

	// Put your logic here

	// ApplicationController_List: end_implement
	res := app.AppCollection{}
	return ctx.OK(res)
}
