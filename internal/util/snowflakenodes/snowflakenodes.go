package snowflakenodes

import (
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

var NodesReport []*snowflake.Node
var NodeBackup *snowflake.Node
var NodeTags *snowflake.Node
var NodeImages *snowflake.Node

var nodeMap map[int]string

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

func GetNodeName(nodeID int64) string {
	return nodeMap[int(nodeID)]
}
