package main

import (
	"github.com/goadesign/goa"
)

// GoaOpenapiController implements the goa_openapi resource.
type GoaOpenapiController struct {
	*goa.Controller
}

// NewGoaOpenapiController creates a goa_openapi controller.
func NewGoaOpenapiController(service *goa.Service) *GoaOpenapiController {
	return &GoaOpenapiController{Controller: service.NewController("GoaOpenapiController")}
}
