package helpers

import (
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
)

var ProjectRoot string

func init() {
	var err error
	ProjectRoot, err = projectRoot()
	if err != nil {
		panic(err)
	}
}

func projectRoot() (string, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	gomod := strings.TrimSpace(string(out))
	if gomod == "" {
		return "", nil // hors module
	}

	u := url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(filepath.Dir(gomod)), // C:\x\y -> C:/x/y
	}

	return u.Path, nil
}
