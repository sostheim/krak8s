package main

import (
	"krak8s/app"

	"github.com/goadesign/goa"
)

// ClusterController implements the cluster resource.
type ClusterController struct {
	*goa.Controller
	ds *DataStore
}

// NewClusterController creates a cluster controller.
func NewClusterController(service *goa.Service, store *DataStore) *ClusterController {
	return &ClusterController{Controller: service.NewController("ClusterController"), ds: store}
}

// MarshalResourcesObject to project media type
func MarshalResourcesObject(obj *ResourceObject) *app.Cluster {
	return &app.Cluster{
		ID:           obj.OID,
		Type:         obj.ObjType,
		NodePoolSize: obj.NodePoolSize,
		NamespaceID:  obj.NamespaceID,
	}
}

// Create runs the create action.
func (c *ClusterController) Create(ctx *app.CreateClusterContext) error {
	// ClusterController_Create: start_implement
	_, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	// TODO: validation step for project oid + namespace oid
	ns, ok := c.ds.Namespace(ctx.Payload.NamespaceID)
	if !ok {
		return ctx.NotFound()
	} else if ns.Resources != nil {
		return ctx.Conflict()
	}
	res := c.ds.NewResource(ctx.Payload.NamespaceID, ctx.Payload.NodePoolSize)
	if res == nil {
		return ctx.InternalServerError()
	}
	url := "/v1/projects/" + ctx.Projectid + "/cluster/" + res.OID
	ns.Resources = &ObjectLink{OID: res.OID, URL: url}

	// ClusterController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *ClusterController) Delete(ctx *app.DeleteClusterContext) error {
	// ClusterController_Delete: start_implement

	_, ok := c.ds.Resource(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	c.ds.DeleteResource(ctx.Projectid)

	// ClusterController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *ClusterController) Get(ctx *app.GetClusterContext) error {
	// ClusterController_Get: start_implement
	resource, ok := c.ds.Resource(ctx.Resourceid)
	if !ok {
		return ctx.NotFound()
	}
	res := MarshalResourcesObject(resource)

	// ClusterController_Get: end_implement
	return ctx.OK(res)
}
