package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// ProjectController implements the project resource.
type ProjectController struct {
	*goa.Controller
}

// NewProjectController creates a project controller.
func NewProjectController(service *goa.Service) *ProjectController {
	return &ProjectController{Controller: service.NewController("ProjectController")}
}

// Create runs the create action.
func (c *ProjectController) Create(ctx *app.CreateProjectContext) error {
	// ProjectController_Create: start_implement

	// Put your logic here

	// ProjectController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *ProjectController) Delete(ctx *app.DeleteProjectContext) error {
	// ProjectController_Delete: start_implement

	// Put your logic here

	// ProjectController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *ProjectController) Get(ctx *app.GetProjectContext) error {
	// ProjectController_Get: start_implement

	// Put your logic here

	// ProjectController_Get: end_implement
	res := &app.Project{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *ProjectController) List(ctx *app.ListProjectContext) error {
	// ProjectController_List: start_implement

	// Put your logic here

	// ProjectController_List: end_implement
	res := app.ProjectCollection{}
	return ctx.OK(res)
}
