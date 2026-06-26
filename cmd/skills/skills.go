package skills

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

type Skill struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "List and install bundled agent skills",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List bundled EchoTik skills",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print bundled skills directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := bundledDir()
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), map[string]any{"path": dir}, nil)
		},
	})
	cmd.AddCommand(newInstall())
	return cmd
}

func newInstall() *cobra.Command {
	var targetDir string
	cmd := &cobra.Command{
		Use:   "install codex",
		Short: "Install bundled skills into Codex",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || args[0] != "codex" {
				return fmt.Errorf("expected: echotik skills install codex")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if targetDir == "" {
				targetDir = defaultCodexSkillsDir()
			}
			skills, err := bundledSkills()
			if err != nil {
				return err
			}
			if len(skills) == 0 {
				return output.NewError("not_found", "no bundled skills found", "reinstall echotik-cli or check: echotik skills path")
			}
			if err := os.MkdirAll(targetDir, 0o755); err != nil {
				return output.NewError("filesystem_error", err.Error(), "check target directory permissions")
			}
			installed := []Skill{}
			for _, skill := range skills {
				dst := filepath.Join(targetDir, skill.Name)
				if err := os.RemoveAll(dst); err != nil {
					return output.NewError("filesystem_error", err.Error(), "check target directory permissions")
				}
				if err := copyDir(skill.Path, dst); err != nil {
					return output.NewError("filesystem_error", err.Error(), "check target directory permissions")
				}
				installed = append(installed, Skill{Name: skill.Name, Path: dst})
			}
			return output.Success(cmd.OutOrStdout(), map[string]any{
				"target":    targetDir,
				"installed": installed,
				"count":     len(installed),
				"hint":      "Restart Codex or open a new Codex session to reload skills.",
			}, nil)
		},
	}
	cmd.Flags().StringVar(&targetDir, "target", "", "skills target directory; default ~/.codex/skills")
	return cmd
}

func list(cmd *cobra.Command) error {
	skills, err := bundledSkills()
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), map[string]any{
		"skills": skills,
		"count":  len(skills),
	}, nil)
}

func bundledSkills() ([]Skill, error) {
	dir, err := bundledDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, output.NewError("not_found", "bundled skills directory not found", "reinstall echotik-cli or check package contents")
	}
	skills := []Skill{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillPath := filepath.Join(dir, name)
		if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err == nil {
			skills = append(skills, Skill{Name: name, Path: skillPath})
		}
	}
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})
	return skills, nil
}

func bundledDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", output.NewError("runtime_error", "cannot locate bundled skills", "reinstall echotik-cli")
	}
	candidates := []string{
		filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "skills")),
		filepath.Clean(filepath.Join(filepath.Dir(os.Args[0]), "..", "skills")),
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", output.NewError("not_found", "bundled skills directory not found", "reinstall echotik-cli or check package contents")
}

func defaultCodexSkillsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".codex/skills"
	}
	return filepath.Join(home, ".codex", "skills")
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
