package bootstrap

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	seedMarkerFile    = ".seed_executed"
	dictMarkerFile    = ".dict_seed_executed"
	legacyMarkerValue = "executed"
)

// applicationRoot locates the directory that contains the migrations folder.
// This keeps runtime markers stable when the process is started from cmd/ or
// from the directory containing a built binary.
func applicationRoot() string {
	candidates := make([]string, 0, 2)
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, ancestors(cwd)...)
	}
	if executable, err := os.Executable(); err == nil {
		candidates = append(candidates, ancestors(filepath.Dir(executable))...)
	}

	for _, candidate := range candidates {
		if fileExists(filepath.Join(candidate, "migrations", "init_seed_data.sql")) ||
			fileExists(filepath.Join(candidate, "migrations", "init.sql")) {
			return candidate
		}
	}
	if len(candidates) > 0 {
		return candidates[0]
	}
	return "."
}

func ancestors(start string) []string {
	result := make([]string, 0, 4)
	current, err := filepath.Abs(start)
	if err != nil {
		return result
	}
	for {
		result = append(result, current)
		parent := filepath.Dir(current)
		if parent == current {
			return result
		}
		current = parent
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func applicationFile(relativePath string) string {
	return filepath.Join(applicationRoot(), relativePath)
}

func initializationMarkerPath(name string) string {
	return applicationFile(name)
}

func initializationComplete(path, version string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	value := strings.TrimSpace(string(data))
	return value == version || (version == "1" && value == legacyMarkerValue)
}

func markInitializationComplete(path, version string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".init-marker-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.WriteString(version + "\n"); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
