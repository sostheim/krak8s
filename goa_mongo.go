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

// Get runs the get action.
func (c *GoaMongoController) Get(ctx *app.GetGoaMongoContext) error {
	// GoaMongoController_Get: start_implement

	// Put your logic here

	// GoaMongoController_Get: end_implement
	res := &app.Mongo{}
	return ctx.OK(res)
}
