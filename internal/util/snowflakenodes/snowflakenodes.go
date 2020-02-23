package snowflakenodes

import (
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

var NodesReport []*snowflake.Node
var NodeBackup *snowflake.Node
var NodeLCHandler *snowflake.Node
var NodeTags *snowflake.Node
var NodeImages *snowflake.Node

func SetupSnowflakeNodes() error {
	NodesReport = make([]*snowflake.Node, len(static.ReportTypes))
	var err error
	for i := range static.ReportTypes {
		NodesReport[i], err = snowflake.NewNode(int64(i))
		if err != nil {
			return err
		}
	}

	NodeBackup, err = snowflake.NewNode(100)
	NodeLCHandler, err = snowflake.NewNode(110)
	NodeTags, err = snowflake.NewNode(120)
	NodeImages, err = snowflake.NewNode(130)

	return err
}
