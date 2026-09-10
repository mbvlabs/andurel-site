package docs

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
	Sections []Section
}

var Catalog = []Version{
	{
		Name: "latest",
		Sections: []Section{
			{
				Title: "Getting Started",
				Pages: []Page{
					{Slug: "introduction", Title: "Introduction", Description: "Meet Andurel v2 and its application architecture."},
					{Slug: "installation", Title: "Installation", Description: "Install the v2 development CLI and create an application."},
					{Slug: "configuration", Title: "Configuration", Description: "Compose validated application configuration with Fx."},
					{Slug: "directory-structure", Title: "Directory Structure", Description: "Understand the files and process boundaries in a v2 project."},
					{Slug: "frontend-options", Title: "Frontend Options", Description: "Choose Templ with Datastar or Inertia v3."},
					{Slug: "v2-migration", Title: "Moving to v2", Description: "Plan a manual move from an Andurel v1 application."},
				},
			},
			{
				Title: "The Basics",
				Pages: []Page{
					{Slug: "routing", Title: "Routing", Description: "Declare typed routes and register Echo handlers."},
					{Slug: "controllers", Title: "Controllers", Description: "Handle requests with explicit injected dependencies."},
					{Slug: "views", Title: "Views", Description: "Render Templ, Datastar, and Inertia v3 responses."},
				},
			},
			{
				Title: "Framework Packages",
				Pages: []Page{
					{Slug: "framework-packages", Title: "Framework Packages", Description: "Understand package boundaries, composition, and versioning."},
					{Slug: "storage", Title: "Storage", Description: "Configure PostgreSQL, transactions, sqlc, tests, and River clients."},
					{
						Slug:        "inertia",
						Title:       "Inertia",
						Description: "Understand the Inertia v3 protocol, props, Vite, SSR, and generators.",
						Children: []Page{
							{Slug: "inertia-renderer", Title: "Renderer", Description: "Construct the renderer, register middleware, and wire application lifecycle."},
							{Slug: "inertia-pages", Title: "Pages and Visits", Description: "Render initial documents and client visits from one Page call."},
							{Slug: "inertia-vite", Title: "Root Document and Vite", Description: "Own the Templ root, Vite tags, and asset versioning."},
							{Slug: "inertia-props", Title: "Props", Description: "Build JSON payloads and compose evaluation policies."},
							{Slug: "inertia-shared", Title: "Shared Data and Redirects", Description: "Share props, flash, validation errors, and protocol redirects."},
							{Slug: "inertia-ssr", Title: "SSR", Description: "Opt pages into SSR with a separate Node runtime."},
							{Slug: "inertia-diagnostics", Title: "Diagnostics", Description: "Classify protocol failures without leaking prop values."},
							{Slug: "inertia-generators", Title: "Generators", Description: "Scaffold Inertia apps and generate pages, types, and routes."},
						},
					},
					{Slug: "hypermedia", Title: "Hypermedia", Description: "Render Templ and build Datastar element, signal, and SSE flows."},
					{Slug: "routing-package", Title: "Routing Package", Description: "Construct typed Go URLs and generated Inertia route helpers."},
					{Slug: "server-package", Title: "Server", Description: "Configure HTTP bounds, lifecycle, and graceful shutdown."},
					{Slug: "email-package", Title: "Email Package", Description: "Use transport-neutral messages, providers, and retry classification."},
					{Slug: "validation", Title: "Validation", Description: "Build structured field errors and reusable validation rules."},
				},
			},
			{
				Title: "Data and Services",
				Pages: []Page{
					{Slug: "database", Title: "Database", Description: "Use the shared PostgreSQL connection, Bun models, and transactions."},
					{Slug: "sql-queries", Title: "SQL Queries", Description: "Use sqlc for complex queries inside the model layer."},
					{Slug: "authentication", Title: "Authentication", Description: "Use generated identity, session, and account flows."},
					{Slug: "email", Title: "Email", Description: "Compile and send typed transactional and marketing email."},
					{Slug: "queues", Title: "Queues", Description: "Insert jobs from the web process and run dedicated workers."},
					{Slug: "telemetry", Title: "Telemetry", Description: "Configure logs, metrics, and traces."},
				},
			},
			{
				Title: "Command Line",
				Pages: []Page{
					{Slug: "cli", Title: "CLI Overview", Description: "Discover the v2 project workflow and command surface."},
					{Slug: "code-generation", Title: "Code Generation", Description: "Generate models, queries, controllers, jobs, email, and resources."},
					{Slug: "database-commands", Title: "Database Commands", Description: "Create, migrate, rebuild, seed, and inspect PostgreSQL."},
					{Slug: "build-and-deploy", Title: "Build and Deploy", Description: "Build application assets and operate the web and queue processes."},
					{Slug: "agent-workflows", Title: "Agent Workflows", Description: "Use discovery, structured output, and dry-run mutations."},
				},
			},
		},
	},
	{
		Name: "1.5.2",
		Sections: []Section{
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
		},
	},
}
