package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
)

func InitCodeExec(container di.Container) codeexec.Factory {
	return codeexec.NewJdoodleFactory(container)
}
