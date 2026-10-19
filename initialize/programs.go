package initialize

const (
	// Programming languages

	Go         Language = "go"
	Rust       Language = "rust"
	C          Language = "c"
	Zig        Language = "zig"
	Python     Language = "python"
	Typescript Language = "ts"
	Typst      Language = "typ"
)

var initMap = map[Language]project{
	// Programming languages

	Go: {
		Tool:             "go",
		Arguments:        []string{"mod", "init"},
		Files:            []string{"main.go"},
		NeedsProjectName: true,
	},

	Rust: {
		Tool:      "cargo",
		Arguments: []string{"init"},
	},

	C: {
		Files:      []string{"main.c"},
		ManualInit: true,
	},

	Zig: {
		Tool:      "zig",
		Arguments: []string{"init"},
	},

	Python: {
		Tool:      "uv",
		Arguments: []string{"init"},
	},

	Typescript: {
		Tool:      "bun",
		Arguments: []string{"init", "--yes"},
	},

	Typst: {
		ManualInit: true,
		Files: []string{"main.typ"},
	},
}
