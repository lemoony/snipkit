package testutil

import (
	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/utils/system"
)

func NewTestSystem(options ...system.Option) *system.System {
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	options = append(options, system.WithFS(ufs))
	return system.NewSystem(options...)
}
