package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// GoaNamespaceController implements the goa_namespace resource.
type GoaNamespaceController struct {
	*goa.Controller
}

// NewGoaNamespaceController creates a goa_namespace controller.
func NewGoaNamespaceController(service *goa.Service) *GoaNamespaceController {
	return &GoaNamespaceController{Controller: service.NewController("GoaNamespaceController")}
}

// Create runs the create action.
func (c *GoaNamespaceController) Create(ctx *app.CreateGoaNamespaceContext) error {
	// GoaNamespaceController_Create: start_implement

	// Put your logic here

	// GoaNamespaceController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *GoaNamespaceController) Delete(ctx *app.DeleteGoaNamespaceContext) error {
	// GoaNamespaceController_Delete: start_implement

	// Put your logic here

	// GoaNamespaceController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *GoaNamespaceController) Get(ctx *app.GetGoaNamespaceContext) error {
	// GoaNamespaceController_Get: start_implement

	// Put your logic here

	// GoaNamespaceController_Get: end_implement
	res := &app.Namespace{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *GoaNamespaceController) List(ctx *app.ListGoaNamespaceContext) error {
	// GoaNamespaceController_List: start_implement

	// Put your logic here

	// GoaNamespaceController_List: end_implement
	res := app.NamespaceCollection{}
	return ctx.OK(res)
}
