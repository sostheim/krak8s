package main

import (
	"errors"
	"krak8s/app"
	"time"

	"github.com/goadesign/goa"
)

// ApplicationController implements the application resource.
type ApplicationController struct {
	*goa.Controller
	ds      *DataStore
	backend *Runner
}

// NewApplicationController creates a application controller.
func NewApplicationController(service *goa.Service, store *DataStore, backend *Runner) *ApplicationController {
	return &ApplicationController{
		Controller: service.NewController("ApplicationController"),
		ds:         store,
		backend:    backend,
	}
}

// MarshalApplicationObject to project media type
func MarshalApplicationObject(obj *ApplicationObject) *app.Application {
	return &app.Application{
		ID:             obj.OID,
		Type:           obj.ObjType,
		NamespaceID:    obj.NamespaceID,
		DeploymentName: obj.Deployment,
		Server:         obj.Server,
		Registry:       obj.ChartRegistry,
		Name:           obj.ChartName,
		Version:        obj.ChartVersion,
		Channel:        obj.Channel,
		Username:       obj.Username,
		Config:         obj.Config,
		JSONValues:     obj.JSONValues,
		CreatedAt:      obj.CreatedAt,
		UpdatedAt:      obj.UpdatedAt,
		Status: &struct {
			DeployedAt time.Time `form:"deployed_at" json:"deployed_at" xml:"deployed_at"`
			Notes      *string   `form:"notes,omitempty" json:"notes,omitempty" xml:"notes,omitempty"`
			State      string    `form:"state" json:"state" xml:"state"`
		}{
			obj.Status.DeployedAt,
			nil,
			obj.Status.State,
		},
	}
}

// Create runs the create action.
func (c *ApplicationController) Create(ctx *app.CreateApplicationContext) error {
	// ApplicationController_Create: start_implement
	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	ns, ok := c.ds.Namespace(ctx.Payload.NamespaceID)
	if !ok {
		return ctx.NotFound()
	}

	found := false
	for _, val := range proj.Namespaces {
		if val.OID == ctx.Payload.NamespaceID {
			found = true
			break
		}
	}
	if !found {
		return ctx.BadRequest(errors.New("Inavlid Namespace Object ID specified in request"))
	}

	app := c.ds.NewApplication(
		ctx.Payload.NamespaceID,
		ctx.Payload.DeploymentName,
		ctx.Payload.Server,
		ctx.Payload.Registry,
		ctx.Payload.Name,
		ctx.Payload.Version,
		ctx.Payload.Channel,
		ctx.Payload.Username,
		ctx.Payload.Password,
		ctx.Payload.Set,
		ctx.Payload.JSONValues)
	if app == nil {
		return ctx.InternalServerError()
	}
	url := APIVersion + APIProjects + ctx.Projectid + APIApplications + app.OID
	ns.Applications = append(ns.Applications, &ObjectLink{OID: app.OID, URL: url})

	c.backend.ChartRequest(AddChart, c.ds, proj, ns, app)

	return ctx.Accepted(MarshalApplicationObject(app))
	// ApplicationController_Create: end_implement
}

// Delete runs the delete action.
func (c *ApplicationController) Delete(ctx *app.DeleteApplicationContext) error {
	// ApplicationController_Delete: start_implement
	app, ok := c.ds.Application(ctx.Appid)
	if !ok {
		return ctx.NotFound()
	}
	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	ns, ok := c.ds.Namespace(app.NamespaceID)
	if !ok {
		return ctx.NotFound()
	}

	index := 0
	found := false
	for i, val := range ns.Applications {
		if val.OID == ctx.Appid {
			index = i
			found = true
			break
		}
	}
	if !found {
		return ctx.BadRequest(errors.New("Inavlid Application Object ID specified in request"))
	}

	c.backend.ChartRequest(RemoveChart, c.ds, proj, ns, app)

	c.ds.DeleteApplication(app)

	copy(ns.Applications[index:], ns.Applications[index+1:])
	ns.Applications[len(ns.Applications)-1] = nil
	ns.Applications = ns.Applications[:len(ns.Applications)-1]

	return ctx.NoContent()

	// ApplicationController_Delete: end_implement
}

// Get runs the get action.
func (c *ApplicationController) Get(ctx *app.GetApplicationContext) error {
	// ApplicationController_Get: start_implement
	app, ok := c.ds.Application(ctx.Appid)
	if !ok {
		return ctx.NotFound()
	}
	res := MarshalApplicationObject(app)
	return ctx.OK(res)
	// ApplicationController_Get: end_implement
}

// List runs the list action.
func (c *ApplicationController) List(ctx *app.ListApplicationContext) error {
	// ApplicationController_List: start_implement
	_, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	// TODO: validation step for project oid + namespace oid
	_, ok = c.ds.Namespace(ctx.Payload.Namespaceid)
	if !ok {
		return ctx.NotFound()
	}
	collection := app.ApplicationCollection{}
	apps := c.ds.ApplicationsCollection(ctx.Payload.Namespaceid)
	count := len(apps)
	if count > 0 {
		collection = make(app.ApplicationCollection, count)
		i := 0
		for _, obj := range apps {
			collection[i] = MarshalApplicationObject(obj)
			i++
		}
	}
	return ctx.OK(collection)
	// ApplicationController_List: end_implement
}
