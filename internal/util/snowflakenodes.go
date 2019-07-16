package util

import "github.com/bwmarrin/snowflake"

var NodesReport []*snowflake.Node
var NodeBackup *snowflake.Node
var NodeLCHandler *snowflake.Node
var NodeTags *snowflake.Node

func SetupSnowflakeNodes() error {
	NodesReport = make([]*snowflake.Node, len(ReportTypes))
	var err error
	for i := range ReportTypes {
		NodesReport[i], err = snowflake.NewNode(int64(i))
		if err != nil {
			return err
		}
	}

	NodeBackup, err = snowflake.NewNode(100)
	NodeLCHandler, err = snowflake.NewNode(110)
	NodeTags, err = snowflake.NewNode(120)

	return err
}
