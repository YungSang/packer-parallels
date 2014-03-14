package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// These are the extensions of files that are unnecessary for the function
// of a Parallels virtual machine.
var UnnecessaryFileExtensions = []string{".log", ".backup", ".Backup"}

// This step removes unnecessary files from the final result.
//
// Uses:
//   dir    OutputDir
//   ui     packer.Ui
//
// Produces:
//   <nothing>
type StepCleanFiles struct {
	OutputDir string
}

func (s *StepCleanFiles) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Deleting unnecessary Parallels files...")
	files, err := s.ListFiles()
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	for _, path := range files {
		// If the file isn't critical to the function of the
		// virtual machine, we get rid of it.
		unnecessary := false
		ext := filepath.Ext(path)
		for _, unnecessaryExt := range UnnecessaryFileExtensions {
			if unnecessaryExt == ext {
				unnecessary = true
				break
			}
		}

		if unnecessary {
			ui.Message(fmt.Sprintf("Deleting: %s", path))
			if err = os.Remove(path); err != nil {
				if _, serr := os.Stat(path); serr == nil || !os.IsNotExist(serr) {
					state.Put("error", err)
					return multistep.ActionHalt
				}
			}
		}
	}

	return multistep.ActionContinue
}

func (StepCleanFiles) Cleanup(multistep.StateBag) {}

func (s *StepCleanFiles) ListFiles() ([]string, error) {
	files := make([]string, 0, 10)

	visit := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}

	return files, filepath.Walk(s.OutputDir, visit)
}
