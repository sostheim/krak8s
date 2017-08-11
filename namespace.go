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
	return &NamespaceController{Controller: service.NewController("NamespaceController"), ds: store}
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
		return nil
		// return ctx.NotFound()
	}
	ns := c.ds.NewNamespace(ctx.Payload.Name)
	if ns == nil {
		return nil
		// return ctx.ServerError()
	}
	url := "/v1/projects/" + ctx.Projectid + "/namespaces/" + ns.OID
	proj.Namespaces = append(proj.Namespaces, &ObjectLink{OID: ns.OID, URL: url})
	// NamespaceController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *NamespaceController) Delete(ctx *app.DeleteNamespaceContext) error {
	// NamespaceController_Delete: start_implement

	ns, ok := c.ds.Namespace(ctx.Namespaceid)
	if !ok {
		return ctx.NotFound()
	}
	c.ds.DeleteNamespace(ns)

	// NamespaceController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *NamespaceController) Get(ctx *app.GetNamespaceContext) error {
	// NamespaceController_Get: start_implement

	ns, ok := c.ds.Namespace(ctx.Namespaceid)
	if !ok {
		// return ctx.NotFound()
	}
	res := MarshalNamespaceObject(ns)

	// NamespaceController_Get: end_implement
	return ctx.OK(res)
}

// List runs the list action.
func (c *NamespaceController) List(ctx *app.ListNamespaceContext) error {
	// NamespaceController_List: start_implement

	res := app.NamespaceCollection{}
	nses := c.ds.NamespacesCollection(ctx.Projectid)
	count := len(nses)
	if count > 0 {
		res = make(app.NamespaceCollection, count)
		i := 0
		for _, obj := range nses {
			res[i] = MarshalNamespaceObject(obj)
			i++
		}
	}

	// NamespaceController_List: end_implement
	return ctx.OK(res)
}
