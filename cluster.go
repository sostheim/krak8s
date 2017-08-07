package main

import (
	"github.com/goadesign/goa"
	"krak8s/app"
)

// ClusterController implements the cluster resource.
type ClusterController struct {
	*goa.Controller
}

// NewClusterController creates a cluster controller.
func NewClusterController(service *goa.Service) *ClusterController {
	return &ClusterController{Controller: service.NewController("ClusterController")}
}

// Create runs the create action.
func (c *ClusterController) Create(ctx *app.CreateClusterContext) error {
	// ClusterController_Create: start_implement

	// Put your logic here

	// ClusterController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *ClusterController) Delete(ctx *app.DeleteClusterContext) error {
	// ClusterController_Delete: start_implement

	// Put your logic here

	// ClusterController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *ClusterController) Get(ctx *app.GetClusterContext) error {
	// ClusterController_Get: start_implement

	// Put your logic here

	// ClusterController_Get: end_implement
	res := &app.Cluster{}
	return ctx.OK(res)
}
