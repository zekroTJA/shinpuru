package snowflakenodes

import (
	"github.com/bwmarrin/snowflake"
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
	// NodeUnbanRequests is the snowflake node
	// for unban requests.
	NodeUnbanRequests *snowflake.Node
	// NodeUnbanRequests is the snowflake node
	// for karma rules.
	NodeKarmaRules *snowflake.Node
	// NodeGuildLog is the snowflake node
	// for guild logs.
	NodeGuildLog *snowflake.Node

	// nodeMap maps snowflake node IDs with
	// their identifier strings.
	nodeMap map[int]string
)

// Setup initializes the snowflake nodes and
// nodesMap.
func Setup() (err error) {
	nodeMap = make(map[int]string)

	NodeBackup, _ = RegisterNode(100, "backups")
	NodeTags, _ = RegisterNode(120, "tags")
	NodeImages, _ = RegisterNode(130, "images")
	NodeUnbanRequests, _ = RegisterNode(140, "unbanrequests")
	NodeKarmaRules, _ = RegisterNode(150, "karmarules")
	NodeGuildLog, _ = RegisterNode(160, "karmarules")

	return
}

func RegisterNode(id int, name string) (*snowflake.Node, error) {
	nodeMap[id] = name
	return snowflake.NewNode(int64(id))
}

// GetNodeName returns the identifier name of
// the snowflake node by nodeID.
func GetNodeName(nodeID int64) string {
	if ident, ok := nodeMap[int(nodeID)]; ok {
		return ident
	}
	return "undefined"
}
