package util

import "github.com/bwmarrin/snowflake"

var ReportNodes []*snowflake.Node
var BackupNode *snowflake.Node

func SetupSnowflakeNodes() error {
	ReportNodes = make([]*snowflake.Node, len(ReportTypes))
	var err error
	for i := range ReportTypes {
		ReportNodes[i], err = snowflake.NewNode(int64(i))
		if err != nil {
			return err
		}
	}

	BackupNode, err = snowflake.NewNode(100)

	return err
}
