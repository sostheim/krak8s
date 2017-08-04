package main

import (
	"krak8s/app"

	"github.com/goadesign/goa"
	"github.com/golang/glog"
)

// MongodbController implements the mongodb resource.
type MongodbController struct {
	*goa.Controller
}

// NewMongodbController creates a mongodb controller.
func NewMongodbController(service *goa.Service) *MongodbController {
	return &MongodbController{Controller: service.NewController("MongodbController")}
}

// Deploy runs the deploy action.
func (c *MongodbController) Deploy(ctx *app.DeployMongodbContext) error {
	// MongodbController_Deploy: start_implement

	glog.V(3).Infof("Deploy(): deployging mongodb-repset to client: %s, ns: %s",
		ctx.Client, ctx.Ns)

	// MongodbController_Deploy: end_implement
	return nil
}
