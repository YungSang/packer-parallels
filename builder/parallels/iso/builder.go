package iso

import (
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

const BuilderId = "yungsang.parallels"

type Builder struct {
	//	config config
	runner multistep.Runner
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	warnings := make([]string, 0)
	return warnings, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	return nil, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
