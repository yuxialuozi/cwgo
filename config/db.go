package config

import (
	"fmt"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/urfave/cli/v2"
	"strings"
)

type DbArgument struct {
	GoMod           string
	PackagePrefix   string
	IdlPath         string
	IdlType         string
	OutDir          string
	SQLDaoDir       string
	Name            string
	ModelDir        string
	DaoDir          string
	Verbose         bool
	ProtoSearchPath []string
	ProtocOptions   []string // options to pass through to protoc
	ThriftOptions   []string // options to pass through to thriftgo for go flag
	GenBase         bool
}

func NewDbArgument() *DocArgument {
	return &DocArgument{}
}

func (d *DbArgument) ParseCli(ctx *cli.Context) error {
	d.IdlPath = ctx.String(consts.IDLPath)
	d.GoMod = ctx.String(consts.Module)
	d.OutDir = ctx.String(consts.OutDir)
	d.SQLDaoDir = ctx.String(consts.SQLDaoDir)
	d.ModelDir = ctx.String(consts.ModelDir)
	d.DaoDir = ctx.String(consts.DaoDir)
	d.Name = ctx.String(consts.Name)
	d.Verbose = ctx.Bool(consts.Verbose)
	d.ProtoSearchPath = ctx.StringSlice(consts.ProtoSearchPath)
	d.ProtocOptions = ctx.StringSlice(consts.Protoc)
	d.ThriftOptions = ctx.StringSlice(consts.ThriftGo)
	d.GenBase = ctx.Bool(consts.GenBase)
	return nil
}

func (d *DbArgument) Unpack(data []string) error {
	err := util.UnpackArgs(data, d)
	if err != nil {
		return fmt.Errorf("unpack argument failed: %s", err)
	}
	return nil
}

func (d *DbArgument) Pack() ([]string, error) {
	data, err := util.PackArgs(d)
	if err != nil {
		return nil, fmt.Errorf("pack argument failed: %s", err)
	}
	return data, nil
}

func (d *DbArgument) GetThriftgoOptions(prefix string) (string, error) {
	d.ThriftOptions = append(d.ThriftOptions, "package_prefix="+prefix)
	gas := "go:" + strings.Join(d.ThriftOptions, ",")
	return gas, nil
}
