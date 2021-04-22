package database

// Migration describes a databse middleware that
// allows database version migration.
type Migration interface {
	// Migrate checks if the currently applied database
	// model is up to date and migrates the database to
	// the latest state.
	Migrate() error
}
