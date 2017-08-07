// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "krak8s": Application Controllers
//
// Command:
// $ goagen
// --design=krak8s/design
// --out=$(GOPATH)/src/krak8s
// --version=v1.2.0-dirty

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(goa.NewGobDecoder, "application/gob", "application/x-gob")
	service.Decoder.Register(goa.NewXMLDecoder, "application/xml")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// ApplicationController is the controller interface for the Application actions.
type ApplicationController interface {
	goa.Muxer
	Create(*CreateApplicationContext) error
	Delete(*DeleteApplicationContext) error
	Get(*GetApplicationContext) error
	List(*ListApplicationContext) error
}

// MountApplicationController "mounts" a Application resource controller on the given service.
func MountApplicationController(service *goa.Service, ctrl ApplicationController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateApplicationContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*ApplicationPostBody)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/v1/projects/:project/ns/:ns/app", ctrl.MuxHandler("create", h, unmarshalCreateApplicationPayload))
	service.LogInfo("mount", "ctrl", "Application", "action", "Create", "route", "POST /v1/projects/:project/ns/:ns/app")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteApplicationContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	service.Mux.Handle("DELETE", "/v1/projects/:project/ns/:ns/app/:app", ctrl.MuxHandler("delete", h, nil))
	service.LogInfo("mount", "ctrl", "Application", "action", "Delete", "route", "DELETE /v1/projects/:project/ns/:ns/app/:app")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewGetApplicationContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Get(rctx)
	}
	service.Mux.Handle("GET", "/v1/projects/:project/ns/:ns/app/:app", ctrl.MuxHandler("get", h, nil))
	service.LogInfo("mount", "ctrl", "Application", "action", "Get", "route", "GET /v1/projects/:project/ns/:ns/app/:app")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListApplicationContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	service.Mux.Handle("GET", "/v1/projects/:project/ns/:ns/app", ctrl.MuxHandler("list", h, nil))
	service.LogInfo("mount", "ctrl", "Application", "action", "List", "route", "GET /v1/projects/:project/ns/:ns/app")
}

// unmarshalCreateApplicationPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateApplicationPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &applicationPostBody{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// ClusterController is the controller interface for the Cluster actions.
type ClusterController interface {
	goa.Muxer
	Create(*CreateClusterContext) error
	Delete(*DeleteClusterContext) error
	Get(*GetClusterContext) error
}

// MountClusterController "mounts" a Cluster resource controller on the given service.
func MountClusterController(service *goa.Service, ctrl ClusterController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateClusterContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CluterPostBody)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/v1/projects/:project/ns/:ns/cluster", ctrl.MuxHandler("create", h, unmarshalCreateClusterPayload))
	service.LogInfo("mount", "ctrl", "Cluster", "action", "Create", "route", "POST /v1/projects/:project/ns/:ns/cluster")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteClusterContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	service.Mux.Handle("DELETE", "/v1/projects/:project/ns/:ns/cluster", ctrl.MuxHandler("delete", h, nil))
	service.LogInfo("mount", "ctrl", "Cluster", "action", "Delete", "route", "DELETE /v1/projects/:project/ns/:ns/cluster")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewGetClusterContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Get(rctx)
	}
	service.Mux.Handle("GET", "/v1/projects/:project/ns/:ns/cluster", ctrl.MuxHandler("get", h, nil))
	service.LogInfo("mount", "ctrl", "Cluster", "action", "Get", "route", "GET /v1/projects/:project/ns/:ns/cluster")
}

// unmarshalCreateClusterPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateClusterPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &cluterPostBody{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	payload.Finalize()
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// NamespaceController is the controller interface for the Namespace actions.
type NamespaceController interface {
	goa.Muxer
	Create(*CreateNamespaceContext) error
	Delete(*DeleteNamespaceContext) error
	Get(*GetNamespaceContext) error
	List(*ListNamespaceContext) error
}

// MountNamespaceController "mounts" a Namespace resource controller on the given service.
func MountNamespaceController(service *goa.Service, ctrl NamespaceController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateNamespaceContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateNamespacePayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/v1/projects/:project/ns", ctrl.MuxHandler("create", h, unmarshalCreateNamespacePayload))
	service.LogInfo("mount", "ctrl", "Namespace", "action", "Create", "route", "POST /v1/projects/:project/ns")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteNamespaceContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	service.Mux.Handle("DELETE", "/v1/projects/:project/ns/:ns", ctrl.MuxHandler("delete", h, nil))
	service.LogInfo("mount", "ctrl", "Namespace", "action", "Delete", "route", "DELETE /v1/projects/:project/ns/:ns")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewGetNamespaceContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Get(rctx)
	}
	service.Mux.Handle("GET", "/v1/projects/:project/ns/:ns", ctrl.MuxHandler("get", h, nil))
	service.LogInfo("mount", "ctrl", "Namespace", "action", "Get", "route", "GET /v1/projects/:project/ns/:ns")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListNamespaceContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	service.Mux.Handle("GET", "/v1/projects/:project/ns", ctrl.MuxHandler("list", h, nil))
	service.LogInfo("mount", "ctrl", "Namespace", "action", "List", "route", "GET /v1/projects/:project/ns")
}

// unmarshalCreateNamespacePayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateNamespacePayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createNamespacePayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// OpenapiController is the controller interface for the Openapi actions.
type OpenapiController interface {
	goa.Muxer
	goa.FileServer
}

// MountOpenapiController "mounts" a Openapi resource controller on the given service.
func MountOpenapiController(service *goa.Service, ctrl OpenapiController) {
	initService(service)
	var h goa.Handler

	h = ctrl.FileHandler("/openapi", "swagger/swagger.json")
	service.Mux.Handle("GET", "/openapi", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Openapi", "files", "swagger/swagger.json", "route", "GET /openapi")

	h = ctrl.FileHandler("/openapi.json", "swagger/swagger.json")
	service.Mux.Handle("GET", "/openapi.json", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Openapi", "files", "swagger/swagger.json", "route", "GET /openapi.json")

	h = ctrl.FileHandler("/openapi.yaml", "swagger/swagger.yaml")
	service.Mux.Handle("GET", "/openapi.yaml", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Openapi", "files", "swagger/swagger.yaml", "route", "GET /openapi.yaml")
}

// ProjectController is the controller interface for the Project actions.
type ProjectController interface {
	goa.Muxer
	Create(*CreateProjectContext) error
	Delete(*DeleteProjectContext) error
	Get(*GetProjectContext) error
	List(*ListProjectContext) error
}

// MountProjectController "mounts" a Project resource controller on the given service.
func MountProjectController(service *goa.Service, ctrl ProjectController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateProjectPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/v1/projects", ctrl.MuxHandler("create", h, unmarshalCreateProjectPayload))
	service.LogInfo("mount", "ctrl", "Project", "action", "Create", "route", "POST /v1/projects")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	service.Mux.Handle("DELETE", "/v1/projects/:project", ctrl.MuxHandler("delete", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "Delete", "route", "DELETE /v1/projects/:project")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewGetProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Get(rctx)
	}
	service.Mux.Handle("GET", "/v1/projects/:project", ctrl.MuxHandler("get", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "Get", "route", "GET /v1/projects/:project")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListProjectContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	service.Mux.Handle("GET", "/v1/projects", ctrl.MuxHandler("list", h, nil))
	service.LogInfo("mount", "ctrl", "Project", "action", "List", "route", "GET /v1/projects")
}

// unmarshalCreateProjectPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateProjectPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createProjectPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// SwaggerController is the controller interface for the Swagger actions.
type SwaggerController interface {
	goa.Muxer
	goa.FileServer
}

// MountSwaggerController "mounts" a Swagger resource controller on the given service.
func MountSwaggerController(service *goa.Service, ctrl SwaggerController) {
	initService(service)
	var h goa.Handler

	h = ctrl.FileHandler("/swagger", "swagger/swagger.json")
	service.Mux.Handle("GET", "/swagger", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Swagger", "files", "swagger/swagger.json", "route", "GET /swagger")

	h = ctrl.FileHandler("/swagger.json", "swagger/swagger.json")
	service.Mux.Handle("GET", "/swagger.json", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Swagger", "files", "swagger/swagger.json", "route", "GET /swagger.json")

	h = ctrl.FileHandler("/swagger.yaml", "swagger/swagger.yaml")
	service.Mux.Handle("GET", "/swagger.yaml", ctrl.MuxHandler("serve", h, nil))
	service.LogInfo("mount", "ctrl", "Swagger", "files", "swagger/swagger.yaml", "route", "GET /swagger.yaml")
}
