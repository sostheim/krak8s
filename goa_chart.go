package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// GoaChartController implements the goa_chart resource.
type GoaChartController struct {
	*goa.Controller
}

// NewGoaChartController creates a goa_chart controller.
func NewGoaChartController(service *goa.Service) *GoaChartController {
	return &GoaChartController{Controller: service.NewController("GoaChartController")}
}

// Create runs the create action.
func (c *GoaChartController) Create(ctx *app.CreateGoaChartContext) error {
	// GoaChartController_Create: start_implement

	// Put your logic here

	// GoaChartController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *GoaChartController) Delete(ctx *app.DeleteGoaChartContext) error {
	// GoaChartController_Delete: start_implement

	// Put your logic here

	// GoaChartController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *GoaChartController) Get(ctx *app.GetGoaChartContext) error {
	// GoaChartController_Get: start_implement

	// Put your logic here

	// GoaChartController_Get: end_implement
	res := &app.Chart{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *GoaChartController) List(ctx *app.ListGoaChartContext) error {
	// GoaChartController_List: start_implement

	// Put your logic here

	// GoaChartController_List: end_implement
	res := app.ChartCollection{}
	return ctx.OK(res)
}
