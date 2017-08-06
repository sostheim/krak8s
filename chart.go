package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// ChartController implements the chart resource.
type ChartController struct {
	*goa.Controller
}

// NewChartController creates a chart controller.
func NewChartController(service *goa.Service) *ChartController {
	return &ChartController{Controller: service.NewController("ChartController")}
}

// Create runs the create action.
func (c *ChartController) Create(ctx *app.CreateChartContext) error {
	// ChartController_Create: start_implement

	// Put your logic here

	// ChartController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *ChartController) Delete(ctx *app.DeleteChartContext) error {
	// ChartController_Delete: start_implement

	// Put your logic here

	// ChartController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *ChartController) Get(ctx *app.GetChartContext) error {
	// ChartController_Get: start_implement

	// Put your logic here

	// ChartController_Get: end_implement
	res := &app.Chart{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *ChartController) List(ctx *app.ListChartContext) error {
	// ChartController_List: start_implement

	// Put your logic here

	// ChartController_List: end_implement
	res := app.ChartCollection{}
	return ctx.OK(res)
}
