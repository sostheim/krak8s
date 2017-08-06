package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// GoaMongoController implements the goa_mongo resource.
type GoaMongoController struct {
	*goa.Controller
}

// NewGoaMongoController creates a goa_mongo controller.
func NewGoaMongoController(service *goa.Service) *GoaMongoController {
	return &GoaMongoController{Controller: service.NewController("GoaMongoController")}
}

// Create runs the create action.
func (c *GoaMongoController) Create(ctx *app.CreateGoaMongoContext) error {
	// GoaMongoController_Create: start_implement

	// Put your logic here

	// GoaMongoController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *GoaMongoController) Delete(ctx *app.DeleteGoaMongoContext) error {
	// GoaMongoController_Delete: start_implement

	// Put your logic here

	// GoaMongoController_Delete: end_implement
	return nil
}

// Read runs the read action.
func (c *GoaMongoController) Read(ctx *app.ReadGoaMongoContext) error {
	// GoaMongoController_Read: start_implement

	// Put your logic here

	// GoaMongoController_Read: end_implement
	return nil
}
