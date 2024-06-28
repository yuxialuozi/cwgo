package db

import (
	"errors"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"path/filepath"
)

func Db(c *config.DbArgument) error {
	if err := check(c); err != nil {
		return err
	}

	switch c.Name {
	case string(consts.MySQL):
		setLogVerbose(c.Verbose)

	default:
	}

	utils.ReplaceThriftVersion()

	return nil
}

func check(c *config.DbArgument) (err error) {
	//默认为 mysql
	if c.Name == "" {
		c.Name = string(consts.MySQL)
	}

	if c.Name != string(consts.MySQL) {
		return errors.New("db name not supported")
	}

	c.OutDir, err = filepath.Abs(c.OutDir)
	if err != nil {
		return err
	}

	return nil
}

func setLogVerbose(verbose bool) {
	if verbose {
		logs.SetLevel(logs.LevelDebug)
	} else {
		logs.SetLevel(logs.LevelWarn)
	}
}
