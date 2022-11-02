package app

import "github.com/lemoony/snipkit/internal/ui/uimsg"

func (a *appImpl) MigrateConfig() {
	newConfigStr := a.configService.Migrate(false)
	confirmed := a.tui.Confirmation(uimsg.ConfigFileMigrationConfirm(newConfigStr))

	if confirmed {
		a.configService.Migrate(true)
	}

	a.tui.Print(uimsg.ConfigFileMigrationResult(confirmed, a.configService.ConfigFilePath()))
}
