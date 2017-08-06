package main

import (
	"github.com/goadesign/goa"
)

// GoaSwaggerController implements the goa_swagger resource.
type GoaSwaggerController struct {
	*goa.Controller
}

// NewGoaSwaggerController creates a goa_swagger controller.
func NewGoaSwaggerController(service *goa.Service) *GoaSwaggerController {
	return &GoaSwaggerController{Controller: service.NewController("GoaSwaggerController")}
}
