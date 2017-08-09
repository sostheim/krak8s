package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// NamespaceController implements the namespace resource.
type NamespaceController struct {
	*goa.Controller
}

// NewNamespaceController creates a namespace controller.
func NewNamespaceController(service *goa.Service) *NamespaceController {
	return &NamespaceController{Controller: service.NewController("NamespaceController")}
}

// Create runs the create action.
func (c *NamespaceController) Create(ctx *app.CreateNamespaceContext) error {
	// NamespaceController_Create: start_implement

	// Put your logic here

	// NamespaceController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *NamespaceController) Delete(ctx *app.DeleteNamespaceContext) error {
	// NamespaceController_Delete: start_implement

	// Put your logic here

	// NamespaceController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *NamespaceController) Get(ctx *app.GetNamespaceContext) error {
	// NamespaceController_Get: start_implement

	// Put your logic here

	// NamespaceController_Get: end_implement
	res := &app.Namespace{}
	return ctx.OK(res)
}
