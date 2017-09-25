package main

import (
	"errors"
	"krak8s/app"

	"github.com/goadesign/goa"
)

// NamespaceController implements the namespace resource.
type NamespaceController struct {
	*goa.Controller
	ds      *DataStore
	backend *Runner
}

// NewNamespaceController creates a namespace controller.
func NewNamespaceController(service *goa.Service, store *DataStore, backend *Runner) *NamespaceController {
	return &NamespaceController{
		Controller: service.NewController("NamespaceController"),
		ds:         store,
		backend:    backend,
	}
}

// MarshaApplicationRef to project media type
func MarshaApplicationRef(obj *ObjectLink) *app.ApplicationRef {
	return &app.ApplicationRef{
		Oid: obj.OID,
		URL: obj.URL,
	}
}

// MarshalNamespaceObject to project media type
func MarshalNamespaceObject(obj *NamespaceObject) *app.Namespace {
	ns := &app.Namespace{
		ID:        obj.OID,
		Type:      obj.ObjType,
		Name:      obj.Name,
		CreatedAt: obj.CreatedAt,
	}

	if obj.Resources != nil {
		ns.Resources = &app.ClusterRef{Oid: obj.Resources.OID, URL: obj.Resources.URL}
	}
	count := len(obj.Applications)
	if count > 0 {
		ns.Applications = make(app.ApplicationRefCollection, count)
		i := 0
		for _, link := range obj.Applications {
			ns.Applications[i] = MarshaApplicationRef(link)
			i++
		}
	}
	return ns
}

// Create runs the create action.
func (c *NamespaceController) Create(ctx *app.CreateNamespaceContext) error {
	// NamespaceController_Create: start_implement
	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	ns := c.ds.NewNamespace(ctx.Payload.Name)
	if ns == nil {
		return ctx.InternalServerError()
	}
	url := APIVersion + APIProjects + ctx.Projectid + APINamespaces + ns.OID
	proj.Namespaces = append(proj.Namespaces, &ObjectLink{OID: ns.OID, URL: url})
	return ctx.Created(MarshalNamespaceObject(ns))
	// NamespaceController_Create: end_implement
}

// Delete runs the delete action.
func (c *NamespaceController) Delete(ctx *app.DeleteNamespaceContext) error {
	// NamespaceController_Delete: start_implement
	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	ns, ok := c.ds.Namespace(ctx.Namespaceid)
	if !ok {
		return ctx.NotFound()
	}

	// ensure that the namespace specified in the request is a member of this project
	index := 0
	found := false
	for i, val := range proj.Namespaces {
		if val.OID == ctx.Namespaceid {
			index = i
			found = true
			break
		}
	}
	if !found {
		return ctx.BadRequest(errors.New("Inavlid Namespace Object ID specified in request"))
	}

	for _, applink := range ns.Applications {
		if app, ok := c.ds.Application(applink.OID); ok {
			c.backend.ChartRequest(RemoveChart, c.ds, proj, ns, app)
		}
	}
	if res, ok := c.ds.Resource(ns.Resources.OID); ok {
		c.backend.ProjectRequest(RemoveProject, c.ds, proj, ns, res)
	}
	c.ds.DeleteNamespace(ns)

	copy(proj.Namespaces[index:], proj.Namespaces[index+1:])
	proj.Namespaces[len(proj.Namespaces)-1] = nil
	proj.Namespaces = proj.Namespaces[:len(proj.Namespaces)-1]

	return ctx.NoContent()
	// NamespaceController_Delete: end_implement
}

// Get runs the get action.
func (c *NamespaceController) Get(ctx *app.GetNamespaceContext) error {
	// NamespaceController_Get: start_implement
	if _, ok := c.ds.Project(ctx.Projectid); !ok {
		return ctx.NotFound()
	}
	ns, ok := c.ds.Namespace(ctx.Namespaceid)
	if !ok {
		return ctx.NotFound()
	}
	res := MarshalNamespaceObject(ns)
	return ctx.OK(res)
	// NamespaceController_Get: end_implement
}

// List runs the list action.
func (c *NamespaceController) List(ctx *app.ListNamespaceContext) error {
	// NamespaceController_List: start_implement
	if _, ok := c.ds.Project(ctx.Projectid); !ok {
		return ctx.NotFound()
	}
	collection := app.NamespaceCollection{}
	nses := c.ds.NamespacesCollection(ctx.Projectid)
	count := len(nses)
	if count > 0 {
		collection = make(app.NamespaceCollection, count)
		i := 0
		for _, obj := range nses {
			collection[i] = MarshalNamespaceObject(obj)
			i++
		}
	}
	return ctx.OK(collection)
	// NamespaceController_List: end_implement
}
