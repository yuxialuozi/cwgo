package plugin

import (
	"fmt"
	"github.com/cloudwego/cwgo/config"
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

func (plu *thriftGoPlugin) handleRequest() error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("read request failed: %s", err.Error())
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		return fmt.Errorf("unmarshal request failed: %s", err.Error())
	}

	plu.req = req
	return nil
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
