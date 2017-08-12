package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// API Resources

// The top level resource is /projects which contains a collection of the existing
// users in the system (users and projects are synonyms).
var _ = Resource("project", func() {
	Description("Manage {create, delete}, and get a project, or all projects")

	DefaultMedia(Project)
	BasePath("/projects")

	CanonicalActionName("get")

	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve all projects.")
		// Can always return success, including w/ empty list.
		Response(OK, CollectionOf(Project))
	})

	Action("get", func() {
		Routing(GET("/:projectid"))
		Description("Retrieve project with given id.")
		Response(OK, Project)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("create", func() {
		Routing(POST(""))
		Description("Create a new project entry with the provided name.")
		Payload(func() {
			Member("name")
			Required("name")
		})
		Response(Created, Project)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError)
	})

	Action("delete", func() {
		Routing(DELETE("/:projectid"))
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("namespace", func() {
	Description("Manage {create, delete}, and get project's namespace(s)")

	Parent("project")
	BasePath("namespaces")

	CanonicalActionName("get")

	Action("create", func() {
		Routing(POST(""))
		Description("Create a namespace in the specified project")
		Payload(func() {
			Member("name")
			Required("name")
		})
		Response(Created, Namespace)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError)
		Response(NotFound) // IFF projectid != valid project
	})

	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve all projects.")
		// Can always return success, including w/ empty list.
		Response(OK, CollectionOf(Namespace))
	})

	Action("get", func() {
		Routing(GET("/:namespaceid"))
		Description("Get the details of the specified namespace from the project")
		Response(NotFound)
		Response(OK, Namespace)
	})

	Action("delete", func() {
		Routing(DELETE("/:namespaceid"))
		Description("Delete the specified namespace from the project")
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("application", func() {
	Description("Manage {create, delete}, and get namespaces's Application(s)")

	Parent("project")
	BasePath("applications")

	CanonicalActionName("get")

	Action("create", func() {
		Routing(POST(""))
		Description("Create an application deployment")
		Payload(ApplicationPostBody)
		Response(Accepted, "^/projects/[a-f,0-9]+/applications/[a-f,0-9]+")
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError)
		Response(NotFound)
	})

	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve the collection of all applications in the namespace.")
		// Can always return success, including w/ empty list.
		Response(OK, CollectionOf(Application))
	})

	Action("get", func() {
		Routing(GET("/:appid"))
		Description("Get the status of the specified application from the project")
		Response(NotFound)
		Response(OK, Application)
	})

	Action("delete", func() {
		Routing(DELETE("/:appid"))
		Description("Delete the specified application from the project")
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("cluster", func() {
	Description("Manage {create, delete}, and get cluster resources")

	Parent("project")
	BasePath("cluster")

	CanonicalActionName("get")

	Action("create", func() {
		Routing(POST(""))
		Description("Create the cluster resources")
		Payload(ClusterPostBody)
		Response(Accepted, "^/projects/[a-f,0-9]+/cluster/[a-f,0-9]+")
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError)
		Response(Conflict)
		Response(NotFound)
	})

	Action("get", func() {
		Routing(GET("/:resourceid"))
		Description("Get the status of the cluster resources")
		Response(OK, Cluster)
		Response(NotFound)
	})

	Action("delete", func() {
		Routing(DELETE(""))
		Description("Delete the cluster resource")
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})
