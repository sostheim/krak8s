package main

import (
	"krak8s/app"

	"github.com/goadesign/goa"
)

// ProjectController implements the project resource.
type ProjectController struct {
	*goa.Controller
	ds      *DataStore
	backend *Runner
}

// NewProjectController creates a project controller.
func NewProjectController(service *goa.Service, store *DataStore, backend *Runner) *ProjectController {
	return &ProjectController{
		Controller: service.NewController("ProjectController"),
		ds:         store,
		backend:    backend,
	}
}

// MarshaNamespaceRef to project media type
func MarshaNamespaceRef(obj *ObjectLink) *app.NamespaceRef {
	return &app.NamespaceRef{
		Oid: obj.OID,
		URL: obj.URL,
	}
}

// MarshalProjectObject to project media type
func MarshalProjectObject(obj *ProjectObject) *app.Project {
	proj := &app.Project{
		ID:        obj.OID,
		Type:      obj.ObjType,
		Name:      obj.Name,
		CreatedAt: obj.CreatedAt,
	}

	count := len(obj.Namespaces)
	if count > 0 {
		proj.Namespaces = make(app.NamespaceRefCollection, count)
		i := 0
		for _, link := range obj.Namespaces {
			proj.Namespaces[i] = MarshaNamespaceRef(link)
			i++
		}
	}
	return proj
}

// Create runs the create action.
func (c *ProjectController) Create(ctx *app.CreateProjectContext) error {
	// ProjectController_Create: start_implement
	proj := c.ds.NewProject(ctx.Payload.Name)
	if proj == nil {
		return ctx.InternalServerError()
	}
	return ctx.Created(MarshalProjectObject(proj))
	// ProjectController_Create: end_implement
}

// Delete runs the delete action.
func (c *ProjectController) Delete(ctx *app.DeleteProjectContext) error {
	// ProjectController_Delete: start_implement
	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}

	for _, nslink := range proj.Namespaces {
		if ns, ok := c.ds.Namespace(nslink.OID); ok {
			for _, applink := range ns.Applications {
				if app, ok := c.ds.Application(applink.OID); ok {
					if app.Status.State == ApplicationDeployed {
						c.backend.ChartRequest(RemoveChart, c.ds, proj, ns, app)
					}
				}
			}
			if ns.Resources != nil {
				if res, ok := c.ds.Resource(ns.Resources.OID); ok {
					if res.State == ResourceErrorStarting || res.State == ResourceActive || res.State == ResourceErrorDeleting {
						c.backend.ProjectRequest(RemoveProject, c.ds, proj, ns, res)
					}
				}
			}
		}
	}

	c.ds.DeleteProject(proj)
	return ctx.NoContent()
	// ProjectController_Delete: end_implement
}

// Get runs the get action.
func (c *ProjectController) Get(ctx *app.GetProjectContext) error {
	// ProjectController_Get: start_implement
	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	res := MarshalProjectObject(proj)
	return ctx.OK(res)
	// ProjectController_Get: end_implement
}

// List runs the list action.
func (c *ProjectController) List(ctx *app.ListProjectContext) error {
	// ProjectController_List: start_implement
	collection := app.ProjectCollection{}
	projects := c.ds.ProjectsCollection()
	count := len(projects)
	if count > 0 {
		collection = make(app.ProjectCollection, count)
		i := 0
		for _, obj := range projects {
			collection[i] = MarshalProjectObject(obj)
			i++
		}
	}
	return ctx.OK(collection)
	// ProjectController_List: end_implement
}
