package snowflakenodes

import (
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

var (
	// NodesReport contains snowflake nodes for
	// each report type.
	NodesReport []*snowflake.Node

	// NodeBackup is the snowflake node for
	// backup IDs.
	NodeBackup *snowflake.Node
	// NodeTags is the snowflake node for
	// chat tag IDs.
	NodeTags *snowflake.Node
	// NodeImages is the snowflake node for
	// public images.
	NodeImages *snowflake.Node

	// nodeMap maps snowflake node IDs with
	// their identifier strings.
	nodeMap map[int]string
)

// Setup initializes the snowflake nodes and
// nodesMap.
func Setup() error {
	nodeMap = make(map[int]string)
	NodesReport = make([]*snowflake.Node, len(static.ReportTypes))
	var err error

	for i, t := range static.ReportTypes {
		NodesReport[i], err = snowflake.NewNode(int64(i))
		if err != nil {
			return err
		}
		nodeMap[i] = "report." + strings.ToLower(t)
	}

	NodeBackup, err = snowflake.NewNode(100)
	NodeTags, err = snowflake.NewNode(120)
	NodeImages, err = snowflake.NewNode(130)

	nodeMap[100] = "backups"
	nodeMap[110] = "lifecyclehandlers"
	nodeMap[120] = "tags"
	nodeMap[130] = "images"

	return err
}

// GetNodeName returns the identifier name of
// the snowflake node by nodeID.
func GetNodeName(nodeID int64) string {
	if ident, ok := nodeMap[int(nodeID)]; ok {
		return ident
	}
	return "undefined"
}
