package docs

const (
	LatestVersion = "latest"
	HeadVersion   = "head"
	LatestRelease = "1.5.5"
	V152Release   = "1.5.2"
)

type Page struct {
	Slug        string
	Title       string
	Description string
	Children    []Page
}

// WalkPages visits each page depth-first, parent before children. The parent
// pointer is only valid for the duration of fn.
func WalkPages(pages []Page, fn func(parent *Page, page Page)) {
	var walk func(parent *Page, pages []Page)
	walk = func(parent *Page, pages []Page) {
		for i := range pages {
			page := pages[i]
			fn(parent, page)
			if len(page.Children) > 0 {
				walk(&page, page.Children)
			}
		}
	}
	walk(nil, pages)
}

func findPage(pages []Page, slug string) (Page, bool) {
	for _, page := range pages {
		if page.Slug == slug {
			return page, true
		}
		if child, ok := findPage(page.Children, slug); ok {
			return child, true
		}
	}
	return Page{}, false
}

type Section struct {
	Title string
	Pages []Page
}

type Version struct {
	Name     string
	Content  string
	Sections []Section
}

func (v Version) ContentDir() string {
	if v.Content != "" {
		return v.Content
	}
	return v.Name
}

var v1Sections = []Section{
	{
		Title: "Getting Started",
		Pages: []Page{
			{
				Slug:        "introduction",
				Title:       "Introduction",
				Description: "Meet Andurel and learn when to use it.",
			},
			{
				Slug:        "installation",
				Title:       "Installation",
				Description: "Install Andurel and create your first application.",
			},
			{
				Slug:        "configuration",
				Title:       "Configuration",
				Description: "Configure an Andurel application through its environment.",
			},
			{
				Slug:        "directory-structure",
				Title:       "Directory Structure",
				Description: "Understand the files and layers in an Andurel project.",
			},
			{
				Slug:        "frontend-options",
				Title:       "Frontend Options",
				Description: "Choose Templ with Datastar or an Inertia frontend.",
			},
		},
	},
	{
		Title: "The Basics",
		Pages: []Page{
			{
				Slug:        "routing",
				Title:       "Routing",
				Description: "Declare typed routes and register HTTP handlers.",
			},
			{
				Slug:        "controllers",
				Title:       "Controllers",
				Description: "Handle requests and render responses.",
			},
			{
				Slug:        "views",
				Title:       "Views",
				Description: "Render server-side Templ views and interactive fragments.",
			},
			{
				Slug:        "validation",
				Title:       "Validation",
				Description: "Validate request data at the application boundary.",
			},
		},
	},
	{
		Title: "Data and Services",
		Pages: []Page{
			{
				Slug:        "database",
				Title:       "Database",
				Description: "Work with PostgreSQL, migrations, models, and factories.",
			},
			{
				Slug:        "authentication",
				Title:       "Authentication",
				Description: "Use the generated account and session flows.",
			},
			{
				Slug:        "email",
				Title:       "Email",
				Description: "Build and send transactional email.",
			},
			{
				Slug:        "queues",
				Title:       "Queues",
				Description: "Run PostgreSQL-backed background jobs.",
			},
			{
				Slug:        "telemetry",
				Title:       "Telemetry",
				Description: "Configure logs, metrics, and traces.",
			},
		},
	},
	{
		Title: "Command Line",
		Pages: []Page{
			{
				Slug:        "cli",
				Title:       "CLI Overview",
				Description: "Discover Andurel's project workflow and commands.",
			},
			{
				Slug:        "code-generation",
				Title:       "Code Generation",
				Description: "Generate models, controllers, views, and complete resources.",
			},
			{
				Slug:        "database-commands",
				Title:       "Database Commands",
				Description: "Create, migrate, rebuild, and inspect a database.",
			},
			{
				Slug:        "build-and-deploy",
				Title:       "Build and Deploy",
				Description: "Build a production binary and verify releases.",
			},
			{
				Slug:        "agent-workflows",
				Title:       "Agent Workflows",
				Description: "Use structured CLI output and dry-run mutations.",
			},
		},
	},
}

var headSections = []Section{
	{
		Title: "Prologue",
		Pages: []Page{
			{Slug: "whats-new", Title: "What's New in v2", Description: "Highlights from Andurel v2.0.0-alpha with a link to the full GitHub changelog."},
			{Slug: "upgrade", Title: "Upgrade Guide", Description: "Plan a manual move from an Andurel v1 application to v2."},
		},
	},
	{
		Title: "Getting Started",
		Pages: []Page{
			{Slug: "introduction", Title: "Introduction", Description: "Meet Andurel v2.0.0-alpha and its application architecture."},
			{Slug: "installation", Title: "Installation", Description: "Install the v2.0.0-alpha CLI and create an application."},
			{Slug: "configuration", Title: "Configuration", Description: "Compose validated application configuration with Fx."},
			{Slug: "directory-structure", Title: "Directory Structure", Description: "Understand the files and process boundaries in a v2 project."},
			{Slug: "frontend", Title: "Frontend", Description: "Choose Inertia v3 (default) or Templ with Datastar."},
			{Slug: "agentic-development", Title: "Agentic Development", Description: "Arm agents with the Andurel skill, AGENTS.md, and JSON CLI discovery."},
			{Slug: "deployment", Title: "Deployment", Description: "Build application assets and operate the web, queue, and SSR processes."},
		},
	},
	{
		Title: "Architecture",
		Pages: []Page{
			{Slug: "request-lifecycle", Title: "Request Lifecycle", Description: "Follow a request through the web, queue, and SSR processes."},
			{Slug: "dependency-injection", Title: "Dependency Injection", Description: "Wire constructors and lifecycles with Fx."},
			{Slug: "framework-packages", Title: "Framework Packages", Description: "Understand package boundaries, composition, and versioning."},
			{Slug: "project-lock", Title: "Project Lock", Description: "Record scaffold choices in andurel.toml and andurel.lock."},
		},
	},
	{
		Title: "The Basics",
		Pages: []Page{
			{Slug: "routing", Title: "Routing", Description: "Declare typed routes, Echo handlers, and Inertia route helpers."},
			{Slug: "controllers", Title: "Controllers", Description: "Handle requests with explicit injected dependencies."},
			{Slug: "cookies-sessions", Title: "Cookies & Sessions", Description: "Use kiks for bagged cookies, sessions, and flash messages."},
			{Slug: "validation", Title: "Validation", Description: "Build structured field errors and reusable validation rules."},
			{Slug: "views", Title: "Views", Description: "Render Templ, Datastar, and Inertia v3 responses."},
		},
	},
	{
		Title: "Database",
		Pages: []Page{
			{Slug: "database", Title: "Getting Started", Description: "Configure PostgreSQL through storage.Connection, transactions, and River."},
			{Slug: "queries", Title: "Queries", Description: "Author narsilc SQL and keep generated clients behind models."},
			{Slug: "migrations", Title: "Migrations & Seeding", Description: "Manage root migrations and seeds with andurel db."},
			{Slug: "models", Title: "Models", Description: "Construct model APIs with storage.Connection and inject them through Fx."},
			{Slug: "factories", Title: "Factories", Description: "Generate and sync model factories for tests and seeds."},
		},
	},
	{
		Title: "Inertia",
		Pages: []Page{
			{
				Slug:        "inertia",
				Title:       "Inertia",
				Description: "Understand the Inertia v3 protocol, props, Vite, SSR, and generators.",
				Children: []Page{
					{Slug: "inertia-renderer", Title: "Renderer", Description: "Construct the renderer, register middleware, and wire application lifecycle."},
					{Slug: "inertia-pages", Title: "Pages and Visits", Description: "Render initial documents and client visits from one Page call."},
					{Slug: "inertia-props", Title: "Props", Description: "Build JSON payloads and compose evaluation policies."},
					{Slug: "inertia-shared", Title: "Shared Data and Redirects", Description: "Share props, flash, validation errors, and protocol redirects."},
					{Slug: "inertia-vite", Title: "Root Document and Vite", Description: "Own the Templ root, Vite tags, and asset versioning."},
					{Slug: "inertia-ssr", Title: "SSR", Description: "Opt pages into SSR with a separate Node runtime."},
					{Slug: "inertia-typescript-sync", Title: "TypeScript Sync", Description: "Generate TypeScript route helpers and payload types with andurel sync."},
					{Slug: "inertia-diagnostics", Title: "Diagnostics", Description: "Classify protocol failures without leaking prop values."},
					{Slug: "inertia-generators", Title: "Generators", Description: "Scaffold Inertia apps and generate pages, types, and routes."},
				},
			},
		},
	},
	{
		Title: "Hypermedia",
		Pages: []Page{
			{Slug: "hypermedia", Title: "Templ & Datastar", Description: "Render Templ and build Datastar element, signal, and SSE flows."},
		},
	},
	{
		Title: "Digging Deeper",
		Pages: []Page{
			{Slug: "queues", Title: "Queues", Description: "Insert jobs from the web process and run dedicated workers."},
			{Slug: "email", Title: "Email", Description: "Compile and send typed transactional and marketing email."},
			{Slug: "telemetry", Title: "Telemetry", Description: "Configure logs, metrics, and traces."},
			{Slug: "http-server", Title: "HTTP Server", Description: "Configure HTTP bounds, lifecycle, and graceful shutdown."},
		},
	},
	{
		Title: "Security",
		Pages: []Page{
			{Slug: "authentication", Title: "Authentication", Description: "Use generated identity, session, and account flows."},
		},
	},
	{
		Title: "Testing",
		Pages: []Page{
			{Slug: "testing", Title: "Testing", Description: "Test Andurel applications with the v2 testing helpers."},
		},
	},
	{
		Title: "Command Line",
		Pages: []Page{
			{Slug: "cli", Title: "Overview", Description: "Discover the v2 project workflow and command surface."},
			{Slug: "generate", Title: "generate", Description: "Create application-owned models, queries, controllers, jobs, email, and resources."},
			{Slug: "sync", Title: "sync", Description: "Refresh derived views, queries, routes, payloads, email, and factories."},
			{Slug: "inspect", Title: "inspect", Description: "Read project shape without writing files."},
			{Slug: "db", Title: "db", Description: "Create, migrate, rebuild, seed, and inspect PostgreSQL."},
			{Slug: "packages", Title: "packages", Description: "List or update Andurel packages in go.mod."},
			{Slug: "skill", Title: "skill", Description: "Install or show the embedded Andurel agent skill."},
			{Slug: "agent-output", Title: "Agent Output", Description: "Use structured CLI output, discovery flags, and dry-run mutations."},
		},
	},
}

var Catalog = []Version{
	{
		Name:     LatestVersion,
		Content:  LatestRelease,
		Sections: v1Sections,
	},
	{
		Name:     LatestRelease,
		Sections: v1Sections,
	},
	{
		Name:     V152Release,
		Sections: v1Sections,
	},
	{
		Name:     HeadVersion,
		Sections: headSections,
	},
}
