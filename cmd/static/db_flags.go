package static

import (
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func dbFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: consts.IDLPath, Usage: "Specify the IDL file path. (.thrift or .proto)"},
		&cli.StringFlag{Name: consts.Module, Aliases: []string{"mod"}, Usage: "Specify the Go module name to generate go.mod."},
		&cli.StringFlag{Name: consts.OutDir, Usage: "Specify output directory, default is current dir."},
		&cli.StringFlag{Name: consts.ModelDir, Usage: "Specify model output directory, default is biz/db/model."},
		&cli.StringFlag{Name: consts.DaoDir, Usage: "Specify dao output directory, default is biz/db/dao."},
		&cli.StringFlag{Name: consts.SQLDaoDir, Usage: "Specify SQL dao output directory, default is biz/db/dao/sql."},
		&cli.StringFlag{Name: consts.Name, Usage: "Specify specific db name, default is mysql."},
		&cli.StringSliceFlag{Name: consts.ProtoSearchPath, Aliases: []string{"I"}, Usage: "Add an IDL search path for includes."},
		&cli.StringSliceFlag{Name: consts.ThriftGo, Aliases: []string{"t"}, Usage: "Specify arguments for the thriftgo. ({flag}={value})"},
		&cli.StringSliceFlag{Name: consts.Protoc, Aliases: []string{"p"}, Usage: "Specify arguments for the protoc. ({flag}={value})"},
		&cli.BoolFlag{Name: consts.Verbose, Usage: "Turn on verbose mode, default is false."},
		&cli.BoolFlag{Name: consts.GenBase, Usage: "Generate base mysql code, default is false."},
	}
}
