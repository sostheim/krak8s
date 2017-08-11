package main

import (
	"krak8s/app"

	"github.com/goadesign/goa"
)

// ProjectController implements the project resource.
type ProjectController struct {
	*goa.Controller
	ds *DataStore
}

// NewProjectController creates a project controller.
func NewProjectController(service *goa.Service, store *DataStore) *ProjectController {
	return &ProjectController{Controller: service.NewController("ProjectController"), ds: store}
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

	p := c.ds.NewProject(ctx.Payload.Name)
	if p == nil {
		return nil
		// return ctx.ServerError()
	}

	// ProjectController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *ProjectController) Delete(ctx *app.DeleteProjectContext) error {
	// ProjectController_Delete: start_implement

	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	c.ds.DeleteProject(proj)

	// ProjectController_Delete: end_implement
	return nil
}

// Get runs the get action.
func (c *ProjectController) Get(ctx *app.GetProjectContext) error {
	// ProjectController_Get: start_implement

	proj, ok := c.ds.Project(ctx.Projectid)
	if !ok {
		return ctx.NotFound()
	}
	res := MarshalProjectObject(proj)

	// ProjectController_Get: end_implement
	return ctx.OK(res)
}

// List runs the list action.
func (c *ProjectController) List(ctx *app.ListProjectContext) error {
	// ProjectController_List: start_implement

	res := app.ProjectCollection{}
	projects := c.ds.ProjectsCollection()
	count := len(projects)
	if count > 0 {
		res = make(app.ProjectCollection, count)
		i := 0
		for _, obj := range projects {
			res[i] = MarshalProjectObject(obj)
			i++
		}
	}

	// ProjectController_List: end_implement
	return ctx.OK(res)
}
