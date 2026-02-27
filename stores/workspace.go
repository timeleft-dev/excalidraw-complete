package stores

import (
	"excalidraw-complete/core"
	"excalidraw-complete/stores/sqlite"
	"os"

	"github.com/sirupsen/logrus"
)

func GetWorkspaceStore() core.WorkspaceStore {
	dataSourceName := os.Getenv("DATA_SOURCE_NAME")
	if dataSourceName == "" {
		// Fall back to a shared workspace DB alongside the document store,
		// or /tmp for non-SQLite deployments.
		dataSourceName = "/tmp/workspace.db"
	}

	logrus.WithField("dataSourceName", dataSourceName).Info("Using SQLite workspace store")
	return sqlite.NewWorkspaceStore(dataSourceName)
}
