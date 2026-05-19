package version

import (
	"bytes"
	"errors"
	"io/fs"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/projectdiscovery/pdtm/pkg/types"
)

var (
	RegexVersionNumber = regexp.MustCompile(`(?m)[v\s](\d+\.\d+\.\d+)`)
	versionCmd         = "--version"
)

func ExtractInstalledVersion(tool types.Tool, basePath string) (string, error) {
	toolPath := filepath.Join(basePath, tool.Name)
	version, err := tryVersionCommand(toolPath, versionCmd)
	if err != nil {
		return "", err
	} 
	return version, nil
}

func tryVersionCommand(toolPath, versionCmd string) (string, error) {
	cmd := exec.Command(toolPath, versionCmd)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb

	if err := cmd.Run(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", errors.New("not installed")
		}
		return "", err
	}

	output := outb.String()
	if output == "" {
		return "", errors.New("empty output")
	}

	installedVersion := RegexVersionNumber.FindString(strings.ToLower(output))
	if installedVersion == "" {
		return "", errors.New("no version found in output")
	}

	version := strings.TrimSpace(installedVersion)
	version = strings.TrimPrefix(version, "v")

	return version, nil
}
