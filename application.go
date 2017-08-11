package main

import (
	"krak8s/app"

	"github.com/goadesign/goa"
)

// ApplicationController implements the application resource.
type ApplicationController struct {
	*goa.Controller
	ds *DataStore
}

// NewApplicationController creates a application controller.
func NewApplicationController(service *goa.Service, store *DataStore) *ApplicationController {
	return &ApplicationController{
		Controller: service.NewController("ApplicationController"),
		ds:         store,
	}
}

// MarshalApplicationObject to project media type
func MarshalApplicationObject(obj *ApplicationObject) *app.Application {
	return &app.Application{
		ID:          obj.OID,
		Type:        obj.ObjType,
		Name:        obj.Name,
		Version:     obj.Version,
		NamespaceID: obj.NamespaceID,
	}
}

// Create runs the create action.
func (c *ApplicationController) Create(ctx *app.CreateApplicationContext) error {
	// ApplicationController_Create: start_implement
	_, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return nil
		// return ctx.NotFound()
	}
	// TODO: validation step for project oid + namespace oid
	ns, ok := c.ds.Namespace(ctx.Payload.NamespaceID)
	if !ok {
		return nil
		// return ctx.NotFound()
	}
	app := c.ds.NewApplication(
		ctx.Payload.NamespaceID,
		ctx.Payload.Name,
		ctx.Payload.Version,
		*ctx.Payload.Set,
		*ctx.Payload.Registry)
	if app == nil {
		return nil
		// return ctx.ServerError()
	}
	url := APIVersion + APIProjects + ctx.Projectid + APIApplications + app.OID
	ns.Applications = append(ns.Applications, &ObjectLink{OID: app.OID, URL: url})

	// ApplicationController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *ApplicationController) Delete(ctx *app.DeleteApplicationContext) error {
	// ApplicationController_Delete: start_implement

	app, ok := c.ds.Application(ctx.Appid)
	if !ok {
		return ctx.NotFound()
	}
	c.ds.DeleteApplication(app)

	// ApplicationController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *ApplicationController) Get(ctx *app.GetApplicationContext) error {
	// ApplicationController_Get: start_implement

	app, ok := c.ds.Application(ctx.Appid)
	if !ok {
		// return ctx.NotFound()
	}
	res := MarshalApplicationObject(app)

	// ApplicationController_Get: end_implement
	return ctx.OK(res)
}

// List runs the list action.
func (c *ApplicationController) List(ctx *app.ListApplicationContext) error {
	// ApplicationController_List: start_implement

	res := app.ApplicationCollection{}
	apps := c.ds.ApplicationsCollection(ctx.Projectid)
	count := len(apps)
	if count > 0 {
		res = make(app.ApplicationCollection, count)
		i := 0
		for _, obj := range apps {
			res[i] = MarshalApplicationObject(obj)
			i++
		}
	}

	// ApplicationController_List: end_implement
	return ctx.OK(res)
}
