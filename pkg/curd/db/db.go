package db

import (
	"errors"
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"os"
	"path/filepath"
	"strings"
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

	if c.ModelDir == "" {
		c.ModelDir = consts.DefaultDbModelOutDir
	}

	c.ModelDir, err = filepath.Abs(filepath.Join(c.OutDir, c.ModelDir))
	if err != nil {
		return err
	}

	if isExist, _ := utils.PathExist(c.ModelDir); !isExist {
		if err = os.MkdirAll(c.ModelDir, 0o755); err != nil {
			return err
		}
	}

	if c.DaoDir == "" {
		c.DaoDir = consts.DefaultDbDaoOutDir
	}
	c.DaoDir, err = filepath.Abs(filepath.Join(c.OutDir, c.DaoDir))
	if err != nil {
		return err
	}
	if isExist, _ := utils.PathExist(c.DaoDir); !isExist {
		if err = os.MkdirAll(c.DaoDir, 0o755); err != nil {
			return err
		}
	}

	if c.SQLDaoDir == "" {
		c.SQLDaoDir = consts.DefaultDbSQLDaoOutDir
	}
	if err != nil {
		return err
	}
	if isExist, _ := utils.PathExist(c.SQLDaoDir); !isExist {
		if err = os.MkdirAll(c.SQLDaoDir, 0o755); err != nil {
			return err
		}
	}

	gopath, err := utils.GetGOPATH()
	if err != nil {
		return fmt.Errorf("get gopath failed: %s", err)
	}
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	gosrc := filepath.Join(gopath, consts.Src)
	gosrc, err = filepath.Abs(gosrc)
	if err != nil {
		log.Warn("Get GOPATH/src path failed:", err.Error())
		os.Exit(1)
	}

	curpath, err := filepath.Abs(consts.CurrentDir)
	if err != nil {
		log.Warn("Get current path failed:", err.Error())
		os.Exit(1)
	}

	if strings.HasPrefix(curpath, gosrc) {
		goPkg := ""
		if goPkg, err = filepath.Rel(gosrc, curpath); err != nil {
			log.Warn("Get GOPATH/src relpath failed:", err.Error())
			os.Exit(1)
		}

		if c.GoMod == "" {
			if utils.IsWindows() {
				c.GoMod = strings.ReplaceAll(goPkg, consts.BackSlash, consts.Slash)
			} else {
				c.GoMod = goPkg
			}
		}

		if c.GoMod != "" {
			if utils.IsWindows() {
				goPkgSlash := strings.ReplaceAll(goPkg, consts.BackSlash, consts.Slash)
				if goPkgSlash != c.GoMod {
					return fmt.Errorf("module name: %s is not the same with GoPkg under GoPath: %s", c.GoMod, goPkgSlash)
				}
			} else {
				if c.GoMod != goPkg {
					return fmt.Errorf("module name: %s is not the same with GoPkg under GoPath: %s", c.GoMod, goPkg)
				}
			}
		}
	}

	if strings.HasPrefix(curpath, gosrc) {
		if c.PackagePrefix, err = filepath.Rel(gosrc, c.ModelDir); err != nil {
			log.Warn("Get GOPATH/src relpath failed:", err.Error())
			os.Exit(1)
		}
	} else {
		if c.GoMod == "" {
			log.Warn("Outside of $GOPATH. Please specify a module name with the '-module' flag.")
			os.Exit(1)
		}
	}

	if c.GoMod != "" {
		module, path, ok := utils.SearchGoMod(curpath, true)
		if ok {
			// go.mod exists
			if module != c.GoMod {
				log.Warnf("The module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)\n",
					c.GoMod, module, path)
				os.Exit(1)
			}
			if c.PackagePrefix, err = filepath.Rel(path, c.ModelDir); err != nil {
				log.Warn("Get package prefix failed:", err.Error())
				os.Exit(1)
			}
			c.PackagePrefix = filepath.Join(c.GoMod, c.PackagePrefix)
		} else {
			if err = utils.InitGoMod(c.GoMod); err != nil {
				log.Warn("Init go mod failed:", err.Error())
				os.Exit(1)
			}
			if c.PackagePrefix, err = filepath.Rel(curpath, c.ModelDir); err != nil {
				log.Warn("Get package prefix failed:", err.Error())
				os.Exit(1)
			}
			c.PackagePrefix = filepath.Join(c.GoMod, c.PackagePrefix)
		}
	}

	c.PackagePrefix = strings.ReplaceAll(c.PackagePrefix, consts.BackSlash, consts.Slash)

	return nil
}

func setLogVerbose(verbose bool) {
	if verbose {
		logs.SetLevel(logs.LevelDebug)
	} else {
		logs.SetLevel(logs.LevelWarn)
	}
}
