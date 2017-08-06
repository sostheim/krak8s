package main

import (
	"github.com/goadesign/goa"
)

// OpenapiController implements the openapi resource.
type OpenapiController struct {
	*goa.Controller
}

// NewOpenapiController creates a openapi controller.
func NewOpenapiController(service *goa.Service) *OpenapiController {
	return &OpenapiController{Controller: service.NewController("OpenapiController")}
}
