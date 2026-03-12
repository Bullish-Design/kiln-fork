// Build orchestrator that dispatches default or custom site generation. @feature:builder
package builder

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Build orchestrates the static site generation process.
func Build(log *slog.Logger) {
	CleanOutputDir(log)
	switch Mode {
	case "custom":
		log.Info("Building site in Custom mode")
		buildCustom(log)
	default:
		log.Info("Building site in Default mode")
		buildDefault(log)
	}
}

func IncrementalBuild(log *slog.Logger, rebuild []string, remove []string) {
	for _, relPath := range remove {
		outPath := filepath.Join(OutputDir, relPath)
		if strings.HasSuffix(relPath, ".md") {
			outPath = strings.TrimSuffix(outPath, ".md")
			if FlatUrls {
				outPath += ".html"
			} else {
				outPath = filepath.Join(outPath, "index.html")
			}
		}
		os.Remove(outPath)
	}

	RebuildFilter = make(map[string]struct{}, len(rebuild))
	for _, relPath := range rebuild {
		RebuildFilter[relPath] = struct{}{}
	}
	defer func() { RebuildFilter = nil }()

	switch Mode {
	case "custom":
		log.Info("Incremental build (custom mode)")
		buildCustom(log)
	default:
		log.Info("Incremental build (default mode)")
		buildDefault(log)
	}
}

var RebuildFilter map[string]struct{}

var (
	OutputDir         string // Destination directory
	InputDir          string // Source directory
	FlatUrls          bool   // Defines if the user opted in for flat urls
	ThemeName         string // Theme name
	FontName          string // Font name
	BaseURL           string // Base URL of the application
	SiteName          string // Sitename
	Mode              string // Mode, either default or custom
	LayoutName        string // Layout name
	DisableTOC        bool   // Disables table of contents
	DisableLocalGraph bool   // Disables local graph
	DisableBacklinks  bool   // Disables backlinks panel
	Lang              string // Language code for the site
	AccentColorName   string // Accent color override (palette color name)
)
