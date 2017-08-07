package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// API Resources

// The top level resource is /projects which contains a collection of the existing
// users in the system (users and projects are synonyms).
var _ = Resource("project", func() {
	Description("Manage {create, delete} individual projects, read the list of all projects, read a specific project")

	DefaultMedia(Project)
	BasePath("/projects")

	CanonicalActionName("get")

	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve all projects.")
		Response(OK, CollectionOf(Project))
	})

	Action("get", func() {
		Routing(GET("/:project"))
		Description("Retrieve project with given id.")
		Params(func() {
			Param("project", String, "project name")
			Required("project")
		})
		Response(OK)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("create", func() {
		Routing(POST(""))
		Description("Create a new user/project")
		Payload(func() {
			Member("identity")
			Required("identity")
		})
		Response(Created, "/projects/[a-z,A-Z,0-9]+")
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Routing(DELETE("/:project"))
		Params(func() {
			Param("project", String, "project name")
			Required("project")
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("namespace", func() {
	Description("Manage {create, delete}, and get project's namespace(s)")

	Parent("project")
	BasePath("ns")

	CanonicalActionName("get")

	Action("create", func() {
		Routing(POST(""))
		Description("Create a namespace in the specified project")
		Payload(func() {
			Member("name")
			Required("name")
		})
		Response(Created, "^/projects/[a-z,A-Z,0-9]+/ns/[a-z,A-Z,0-9]+")
		Response(BadRequest, ErrorMedia)
	})

	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve the collection of all namespaces in the project.")
		Response(OK, CollectionOf(Namespace))
	})

	Action("get", func() {
		Routing(GET("/:ns"))
		Description("Get the details of the specified namespace (ns) in the project")
		Params(func() {
			Param("ns", String, "namespace identifier")
			Required("ns")
		})
		Response(OK, Namespace)
	})

	Action("delete", func() {
		Routing(DELETE("/:ns"))
		Description("Delete the specified namespace (ns)")
		Params(func() {
			Param("ns", String, "namespace identifier")
			Required("ns")
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("application", func() {
	Description("Manage {create, delete}, and get namespaces's Application(s)")

	Parent("namespace")
	BasePath("app")

	CanonicalActionName("get")

	Action("create", func() {
		Routing(POST(""))
		Description("Create an application deployment")
		Payload(ApplicationPostBody)
		Response(Accepted, "^/projects/[a-z,A-Z,0-9]+/ns/[a-z,A-Z,0-9]+/app/[a-z,A-Z,0-9]+")
		Response(BadRequest, ErrorMedia)
	})

	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve the collection of all applications in the namespace.")
		Response(OK, CollectionOf(Application))
	})

	Action("get", func() {
		Routing(GET("/:app"))
		Description("Get the status of the specified application")
		Response(OK, Application)
	})

	Action("delete", func() {
		Routing(DELETE("/:app"))
		Description("Delete the specified application")
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("cluster", func() {
	Description("Manage {create, delete}, and get cluster resources")

	Parent("namespace")
	BasePath("cluster")

	CanonicalActionName("get")

	Action("create", func() {
		Routing(POST(""))
		Description("Create the specified cluster resources")
		Payload(ClusterPostBody)
		Response(Accepted, "^/projects/[a-z,A-Z,0-9]+/ns/[a-z,A-Z,0-9]+/cluster")
		Response(BadRequest, ErrorMedia)
	})

	Action("get", func() {
		Routing(GET(""))
		Description("Get the status of the cluster resources")
		Response(OK, Cluster)
	})

	Action("delete", func() {
		Routing(DELETE(""))
		Description("Delete the cluster resource")
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})
