package plugin

import (
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/parser"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/pkg/curd/doc/mongo/codegen"
	"github.com/cloudwego/cwgo/pkg/curd/extract"
	"github.com/cloudwego/cwgo/pkg/curd/parse"
	"github.com/cloudwego/cwgo/pkg/curd/template"
	"github.com/cloudwego/hertz/cmd/hz/meta"
	"os"
	"os/exec"
	"strings"
)

func MysqlTriggerPlugin(c *config.DbArgument) error {
	cmd, err := buildPluginCmd(c)
	if err != nil {
		return fmt.Errorf("build plugin command failed: %v", err)
	}

	buf, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("plugin cwgo-db returns error: %v, cause:\n%v", err, string(buf))
	}

	// If len(buf) != 0, the plugin returned the log.
	if len(buf) != 0 {
		fmt.Println(string(buf))
	}

	if c.IdlType == meta.IdlProto {
		info := &extract.PbUsedInfo{
			DbArgs: c,
		}
		rawStructs, err := info.ParsePbIdl()
		if err != nil {
			return err
		}
		operations, err := parse.HandleOperations(rawStructs)
		if err != nil {
			return err
		}
		methodRenders := codegen.HandleCodegen(operations)

		if c.GenBase {
			if err = generateBaseMongoFile(info.DocArgs.DaoDir, info.ImportPaths, codegen.HandleBaseCodegen()); err != nil {
				return err
			}
		}

		if err = info.GeneratePbFile(); err != nil {
			return err
		}
		if err = generatePbMongoFile(rawStructs, methodRenders, info); err != nil {
			return err
		}
	}

	return nil
}

func buildPluginCmd(args *config.DbArgument) (*exec.Cmd, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to detect current executable, err: %v", err)
	}

	argPacks, err := args.Pack()
	if err != nil {
		return nil, err
	}
	kas := strings.Join(argPacks, ",")

	path, err := utils.LookupTool(args.IdlType)
	if err != nil {
		return nil, err
	}
	cmd := &exec.Cmd{
		Path: path,
	}
	if args.IdlType == meta.IdlThrift {
		os.Setenv(consts.CwgoDbPluginMode, consts.ThriftCwgoDbPluginName)

		thriftOpt, err := args.GetThriftgoOptions(args.PackagePrefix)
		if err != nil {
			return nil, err
		}
		cmd.Args = append(cmd.Args, meta.TpCompilerThrift)
		if args.Verbose {
			cmd.Args = append(cmd.Args, "-v")
		}
		cmd.Args = append(cmd.Args,
			"-o", args.ModelDir,
			"-p", "cwgo-doc="+exe+":"+kas,
			"-g", thriftOpt,
			"-r",
			args.IdlPath,
		)
	} else {
		cmd.Args = append(cmd.Args, meta.TpCompilerProto)

		var isFindIdl bool

		var importPaths []string

		for _, inc := range args.ProtoSearchPath {
			idlParser := parser.NewProtoParser()

			if !isFindIdl {
				_, importPaths, err = idlParser.GetDependentFilePaths(inc, args.IdlPath)
				if err == nil {
					isFindIdl = true
				}

			}

			cmd.Args = append(cmd.Args, "-I", inc)

			cmd.Args = append(cmd.Args, "--go_out="+args.ModelDir)
			for _, kv := range args.ProtocOptions {
				cmd.Args = append(cmd.Args, "--"+kv)
			}

			cmd.Args = append(cmd.Args, importPaths...)
			cmd.Args = append(cmd.Args, args.IdlPath)
		}

	}
	return cmd, err
}

func MysqlPluginMode() {
	mode := os.Getenv(consts.CwgoDbPluginMode)
	if len(os.Args) <= 1 && mode != "" {
		switch mode {
		case consts.ThriftCwgoDbPluginName:
			os.Exit(thriftPluginRun())
		}
	}
}

func generatePbMongoFile(structs []*extract.IdlExtractStruct, methodRenders [][]*template.MethodRender, info *extract.PbUsedInfo) error {

	return nil
}
