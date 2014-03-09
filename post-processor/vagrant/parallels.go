package vagrant

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/packer/packer"
)

type ParallelsProvider struct{}

func (p *ParallelsProvider) KeepInputArtifact() bool {
	return false
}

func (p *ParallelsProvider) Process(ui packer.Ui, artifact packer.Artifact, dir string) (vagrantfile string, metadata map[string]interface{}, err error) {
	// Create the metadata
	metadata = map[string]interface{}{"provider": "parallels"}

	// Copy all of the original contents into the temporary directory
	for _, path := range artifact.Files() {
		if extension := filepath.Ext(path); extension == ".log" {
			continue
		}
		if extension := filepath.Ext(path); extension == ".backup" {
			continue
		}
		if extension := filepath.Ext(path); extension == ".Backup" {
			continue
		}

		var pvmPath string

		tmpPath := filepath.ToSlash(path)
		pathRe := regexp.MustCompile(`^(.+?)([^/]+\.pvm/.+?)$`)
		matches := pathRe.FindStringSubmatch(tmpPath)
		if matches != nil {
			pvmPath = filepath.FromSlash(matches[2])
		} else {
			continue // Just copy a pvm
		}
		dstPath := filepath.Join(dir, pvmPath)

		ui.Message(fmt.Sprintf("Copying: %s", path))
		if err = CopyContents(dstPath, path); err != nil {
			return
		}
	}

	return
}
