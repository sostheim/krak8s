package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// MongoController implements the mongo resource.
type MongoController struct {
	*goa.Controller
}

// NewMongoController creates a mongo controller.
func NewMongoController(service *goa.Service) *MongoController {
	return &MongoController{Controller: service.NewController("MongoController")}
}

// Create runs the create action.
func (c *MongoController) Create(ctx *app.CreateMongoContext) error {
	// MongoController_Create: start_implement

	// Put your logic here

	// MongoController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *MongoController) Delete(ctx *app.DeleteMongoContext) error {
	// MongoController_Delete: start_implement

	// Put your logic here

	// MongoController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *MongoController) Get(ctx *app.GetMongoContext) error {
	// MongoController_Get: start_implement

	// Put your logic here

	// MongoController_Get: end_implement
	res := &app.Mongo{}
	return ctx.OK(res)
}
