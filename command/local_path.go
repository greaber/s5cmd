package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/peak/s5cmd/v2/storage/url"
)

func joinLocalDestination(root *url.URL, relative string) (*url.URL, error) {
	if root.IsRemote() {
		return nil, fmt.Errorf("destination %q is not local", root.Absolute())
	}

	rootPath, err := filepath.Abs(root.Absolute())
	if err != nil {
		return nil, fmt.Errorf("resolve local destination %q: %v", root.Absolute(), err)
	}

	destination := root.Join(relative)
	destinationPath, err := filepath.Abs(destination.Absolute())
	if err != nil {
		return nil, fmt.Errorf("resolve local destination %q: %v", destination.Absolute(), err)
	}

	rel, err := filepath.Rel(rootPath, destinationPath)
	if err != nil {
		return nil, fmt.Errorf("compare local destination paths: %v", err)
	}
	if rel == ".." || filepath.IsAbs(rel) || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("refusing to write object path %q outside destination %q", relative, root.Absolute())
	}

	return destination, nil
}
