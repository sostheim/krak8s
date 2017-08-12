package main

import (
	"krak8s/app"

	"github.com/goadesign/goa"
)

// NamespaceController implements the namespace resource.
type NamespaceController struct {
	*goa.Controller
	ds *DataStore
}

// NewNamespaceController creates a namespace controller.
func NewNamespaceController(service *goa.Service, store *DataStore) *NamespaceController {
	return &NamespaceController{
		Controller: service.NewController("NamespaceController"),
		ds:         store,
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

	ns, ok := c.ds.Namespace(ctx.Namespaceid)
	if !ok {
		return ctx.NotFound()
	}
	c.ds.DeleteNamespace(ns)
	return ctx.NoContent()
	// NamespaceController_Delete: end_implement
}

// Get runs the get action.
func (c *NamespaceController) Get(ctx *app.GetNamespaceContext) error {
	// NamespaceController_Get: start_implement
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
