package main

/*
I like go.mod approach for managing tools (and it can be done even better than described
on wiki, you can use submodule to not pollute you main go.mod), but it have two disadvantages:

1. It doesn't work with some tools (notably golangci-lint)
2. It doesn't work with non-go tools (notably protoc)

So, I decided to experiment with mage to solve this problem.
*/

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/cavaliergopher/grab/v3"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const toolsBinDir = ".tools"

type Tools mg.Namespace

// List of all tools that can be installed with Install.
//
// Initally I intended to have seperate manifest.json file to keep
// version seperated, but why to complicate things?
var managedTools = indexTools([]*Tool{
	{
		Name:    "gqlgen",
		Version: "0.17.36",
		Installer: &GoToolInstaller{
			URL: "github.com/99designs/gqlgen",
		},
	},
	{
		Name:    "golangci-lint",
		Version: "1.54.2",
		Installer: &ArchiveToolInstaller{
			URL:  "https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-linux-amd64.tar.gz",
			Path: "golangci-lint-{{.Version}}-linux-amd64/golangci-lint",
		},
	},
	{
		Name:    "mockery",
		Version: "2.27.1",
		Installer: &ArchiveToolInstaller{
			URL:  "https://github.com/vektra/mockery/releases/download/v{{.Version}}/mockery_{{.Version}}_Linux_x86_64.tar.gz",
			Path: "mockery",
		},
	},
	{
		Name:    "enumer",
		Version: "1.5.8",
		Installer: &GoToolInstaller{
			URL: "github.com/dmarkham/enumer",
		},
	},
})

func getToolPath(name string) string {
	tool, ok := managedTools[name]
	if !ok {
		panic(fmt.Sprintf("unknown tool %q", name))
	}
	return tool.Path()
}

// Install install specified tool (among defined, see tools.go) in .tools directory
func (Tools) Install(ctx context.Context, toolName string) error {
	tool, ok := managedTools[toolName]
	if !ok {
		return fmt.Errorf("unknown tool %q provided", toolName)
	}
	exists, err := tool.Check()
	if err != nil {
		return fmt.Errorf("check for tool %q failed: %w", toolName, err)
	}
	if exists {
		return nil
	}
	if err := ensureToolDir(); err != nil {
		return fmt.Errorf("unable to create directory for tools: %w", err)
	}
	return tool.Install(ctx)
}

// All installs all defined tools
func (t Tools) All(ctx context.Context) {
	var depFuncs []interface{}
	for name := range managedTools {
		depFuncs = append(
			depFuncs,
			mg.F(Tools.Install, name),
		)
	}
	mg.SerialCtxDeps(ctx, depFuncs...)
}

// Gqlgen installs gqlgen
func (t Tools) Gqlgen(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Tools.Install, "gqlgen"))
}

// Golangcilint installs golangci-lint
func (t Tools) Golangcilint(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Tools.Install, "golangci-lint"))
}

// Mockery installs mockery
func (t Tools) Mockery(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Tools.Install, "mockery"))
}

// Enumer installs enumer
func (t Tools) Enumer(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Tools.Install, "enumer"))
}

// Clear remove all tools
func (Tools) Clear() error {
	return os.RemoveAll(toolsBinDir)
}

type Tool struct {
	Name      string
	Version   string
	Installer ToolInstaller
}

func (t *Tool) Path() string {
	return filepath.Join(toolsBinDir, t.Name+"-"+t.Version)
}

func (t *Tool) Check() (bool, error) {
	stat, err := os.Stat(t.Path())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if stat.Mode().IsRegular() {
		return true, nil
	}
	return true, fmt.Errorf("path %q exists, but is not regular file", t.Path())
}

func (t *Tool) Install(ctx context.Context) error {
	return t.Installer.Install(ctx, t.Path(), t.Name, t.Version)
}

type ToolInstaller interface {
	Install(ctx context.Context, dst string, name string, version string) error
}

type GoToolInstaller struct {
	URL string
}

func (t *GoToolInstaller) binName() string {
	return path.Base(t.URL)
}

func (t *GoToolInstaller) Install(ctx context.Context, dst, name, version string) error {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tmpDir)
	err = sh.RunWithV(
		map[string]string{
			"GOBIN": tmpDir,
		},
		"go",
		"install",
		t.URL+"@v"+version,
	)
	if err != nil {
		return err
	}
	return sh.Copy(dst, filepath.Join(tmpDir, t.binName()))
}

type ArchiveToolInstaller struct {
	URL  string
	Path string
}

func (t *ArchiveToolInstaller) Install(ctx context.Context, dst, name, version string) error {
	url, err := interpolateWithVersionName(t.URL, name, version)
	if err != nil {
		return err
	}
	targetPath, err := interpolateWithVersionName(t.Path, name, version)
	if err != nil {
		return err
	}
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tmpDir)
	resp, err := grab.Get(tmpDir, url)
	if err != nil {
		return fmt.Errorf("archive download failed: %w", err)
	}
	file, err := os.Open(resp.Filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// hardcode for tar.gz archives, extend if needed
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)
	tmpDst, err := os.CreateTemp(tmpDir, "")
	if err != nil {
		return err
	}
	defer tmpDst.Close()
	for {
		header, err := tarReader.Next()
		if err != nil {
			return fmt.Errorf("file not found or error in tar: %w", err)
		}
		if header.Typeflag != tar.TypeReg || header.Name != targetPath {
			continue
		}
		_, err = io.Copy(tmpDst, tarReader)
		if err != nil {
			return err
		}
		if err := tmpDst.Chmod(fs.FileMode(header.Mode)); err != nil {
			return err
		}
		break
	}
	if err := tmpDst.Close(); err != nil {
		return err
	}
	return sh.Copy(dst, tmpDst.Name())
}

func ensureToolDir() error {
	return os.MkdirAll(toolsBinDir, 0o750)
}

func indexTools(tools []*Tool) map[string]*Tool {
	result := map[string]*Tool{}
	for _, tool := range tools {
		_, ok := result[tool.Name]
		if ok {
			panic(fmt.Sprintf("found duplicated tool definition for %q", tool.Name))
		}
		result[tool.Name] = tool
	}
	return result
}

func interpolateWithVersionName(src, name, version string) (string, error) {
	t := template.New("")
	t, err := t.Parse(src)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, map[string]string{
		"Name":    name,
		"Version": version,
	})
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
