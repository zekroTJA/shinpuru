package guildlog

type Logger interface {
	Debugf(guildID, message string, data ...interface{}) error
	Infof(guildID, message string, data ...interface{}) error
	Warnf(guildID, message string, data ...interface{}) error
	Errorf(guildID, message string, data ...interface{}) error
	Fatalf(guildID, message string, data ...interface{}) error

	Section(module string) Logger
}
