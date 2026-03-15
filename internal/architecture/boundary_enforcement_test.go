package architecture

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestStorageBoundaryEnforcement(t *testing.T) {
	repoRoot := mustRepoRoot(t)
	violations := scanForPatternViolations(t, repoRoot, []patternRule{
		{re: regexp.MustCompile(`\bos\.ReadFile\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.WriteFile\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.MkdirAll\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.Rename\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.RemoveAll\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.Remove\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.CreateTemp\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.Chmod\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bos\.OpenFile\(`), allowPrefixes: []string{"storage/", "interactive/"}},
		{re: regexp.MustCompile(`\bos\.Open\(`), allowPrefixes: []string{"storage/", "interactive/"}},
		{re: regexp.MustCompile(`\bsql\.Open\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bdb\.Exec\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bdb\.QueryRow\(`), allowPrefixes: []string{"storage/"}},
		{re: regexp.MustCompile(`\bdb\.Query\(`), allowPrefixes: []string{"storage/"}},
	})

	if len(violations) > 0 {
		t.Fatalf("storage boundary violations found:\n%s", strings.Join(violations, "\n"))
	}
}

func TestNetworkBoundaryEnforcement(t *testing.T) {
	repoRoot := mustRepoRoot(t)
	violations := scanForPatternViolations(t, repoRoot, []patternRule{
		{re: regexp.MustCompile(`\bhttp\.NewRequest(?:WithContext)?\(`), allowPrefixes: []string{"network/"}},
		{re: regexp.MustCompile(`\bhttp\.DefaultClient\.Do\(`), allowPrefixes: []string{"network/"}},
	})

	if len(violations) > 0 {
		t.Fatalf("network boundary violations found:\n%s", strings.Join(violations, "\n"))
	}
}

type patternRule struct {
	re            *regexp.Regexp
	allowPrefixes []string
}

func scanForPatternViolations(t *testing.T, repoRoot string, rules []patternRule) []string {
	t.Helper()
	violations := make([]string, 0)
	walkErrors := make([]string, 0)

	if err := filepath.WalkDir(repoRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			walkErrors = append(walkErrors, path+": "+err.Error())
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		rel, relErr := filepath.Rel(repoRoot, path)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		text := string(content)

		for _, rule := range rules {
			if hasAllowedPrefix(rel, rule.allowPrefixes) {
				continue
			}
			if rule.re.FindStringIndex(text) != nil {
				violations = append(violations, rel+": matches "+rule.re.String())
			}
		}

		return nil
	}); err != nil {
		t.Fatalf("failed to walk repository tree: %v", err)
	}

	if len(walkErrors) > 0 {
		t.Fatalf("encountered file walk errors:\n%s", strings.Join(walkErrors, "\n"))
	}

	return violations
}

func hasAllowedPrefix(rel string, allowPrefixes []string) bool {
	for _, prefix := range allowPrefixes {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	return false
}

func mustRepoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}
