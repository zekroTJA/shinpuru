package util

import "github.com/bwmarrin/snowflake"

var NodesReport []*snowflake.Node
var NodeBackup *snowflake.Node
var NodeLTCHandler *snowflake.Node

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
	NodeLTCHandler, err = snowflake.NewNode(110)

	return err
}
