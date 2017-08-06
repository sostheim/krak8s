package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// GoaProjectController implements the goa_project resource.
type GoaProjectController struct {
	*goa.Controller
}

// NewGoaProjectController creates a goa_project controller.
func NewGoaProjectController(service *goa.Service) *GoaProjectController {
	return &GoaProjectController{Controller: service.NewController("GoaProjectController")}
}

// Create runs the create action.
func (c *GoaProjectController) Create(ctx *app.CreateGoaProjectContext) error {
	// GoaProjectController_Create: start_implement

	// Put your logic here

	// GoaProjectController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *GoaProjectController) Delete(ctx *app.DeleteGoaProjectContext) error {
	// GoaProjectController_Delete: start_implement

	// Put your logic here

	// GoaProjectController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *GoaProjectController) Get(ctx *app.GetGoaProjectContext) error {
	// GoaProjectController_Get: start_implement

	// Put your logic here

	// GoaProjectController_Get: end_implement
	res := &app.Project{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *GoaProjectController) List(ctx *app.ListGoaProjectContext) error {
	// GoaProjectController_List: start_implement

	// Put your logic here

	// GoaProjectController_List: end_implement
	res := app.ProjectCollection{}
	return ctx.OK(res)
}
