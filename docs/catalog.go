package docs

type Page struct {
	Slug        string
	Title       string
	Description string
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
