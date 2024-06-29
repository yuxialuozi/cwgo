package plugin

import (
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/curd/extract"
	"github.com/cloudwego/cwgo/pkg/curd/parse"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/thriftgo/plugin"
	"io"
	"os"
)

type thriftGoPlugin struct {
	req    *plugin.Request
	dbArgs *config.DbArgument
}

func thriftPluginRun() int {
	plu := &thriftGoPlugin{}

	if err := plu.handleRequest(); err != nil {
		logs.Errorf("handle request failed: %s", err.Error())
		return meta.PluginError
	}

	if err := plu.parseArgs(); err != nil {
		logs.Errorf("parse args failed: %s", err.Error())
		return meta.PluginError
	}

	return 0
}

func (plu *thriftGoPlugin) handleRequest() int {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return meta.PluginError
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		return meta.PluginError
	}

	tfUsedInfo := &extract.ThriftUsedInfo{
		Req:    plu.req,
		DbArgs: plu.dbArgs,
	}
	rawStructs, err := tfUsedInfo.ParseThriftIdl()
	if err != nil {
		logs.Errorf("parse thrift idl failed: %s", err.Error())
		return meta.PluginError
	}

	operations, err := parse.HandleOperations(rawStructs)
	if err != nil {
		logs.Error(err.Error())
		return meta.PluginError
	}

	return 0
}

func (plu *thriftGoPlugin) parseArgs() error {
	if plu.req == nil {
		return fmt.Errorf("request is nil")
	}
	args := new(config.DbArgument)
	if err := args.Unpack(plu.req.PluginParameters); err != nil {
		logs.Errorf("unpack args failed: %s", err.Error())
		return err
	}
	plu.dbArgs = args
	return nil
}
