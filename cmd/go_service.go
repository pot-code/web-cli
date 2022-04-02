package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/commands"
	"github.com/pot-code/web-cli/internal/constants"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/transformation"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type GoServiceConfig struct {
	GenType string `flag:"type" alias:"t" usage:"api type" validate:"required,oneof=go"` // generation type
	ArgName string `arg:"0" alias:"module_name" validate:"required,var"`                 // go pkg name
}

var GoServiceCmd = util.NewCliCommand("service", "add a go service",
	&GoServiceConfig{
		GenType: "go",
	},
	util.WithAlias([]string{"svc"}),
	util.WithArgUsage("module_name"),
).AddFeature(
	&GenerateGoSimpleService{
		RegistryFile: path.Join("web", "server.go"),
	},
).ExportCommand()

type GenerateGoSimpleService struct {
	RegistryFile    string
	PackageName     string
	ProjectName     string
	CamelModuleName string
	AuthorName      string
	Config          *GoServiceConfig
}

var _ util.CommandFeature = &GenerateGoSimpleService{}

func (gga *GenerateGoSimpleService) Cond(c *cli.Context) bool {
	return true
}

func (gga *GenerateGoSimpleService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoServiceConfig)
	pkgName := strings.ReplaceAll(config.ArgName, "-", "_")
	pkgName = strcase.ToSnake(pkgName)
	gga.PackageName = pkgName

	meta, err := util.ParseGoMod(constants.GoModFile)
	if err != nil {
		return errors.WithStack(err)
	}

	gga.CamelModuleName = strcase.ToCamel(pkgName)
	gga.AuthorName = meta.Author
	gga.ProjectName = meta.ProjectName

	dir := gga.PackageName
	_, err = os.Stat(dir)
	if err == nil {
		log.WithField("module", dir).Infof("module already exists")
		return nil
	}

	if err := gga.updateServerRegistry(); err != nil {
		return errors.Wrap(err, "failed to update registry")
	}

	return gga.generateFiles().Run()
}

func (gga *GenerateGoSimpleService) generateFiles() task.Runner {
	modelName := gga.CamelModuleName
	pkgName := gga.PackageName
	projectName := gga.ProjectName
	authorName := gga.AuthorName

	handlerDeclName := fmt.Sprintf(constants.GoApiHandlerPattern, modelName)
	svcDeclName := fmt.Sprintf(constants.GoApiServicePattern, modelName)
	repoDeclName := fmt.Sprintf(constants.GoApiRepositoryPattern, modelName)

	return util.NewTaskComposer(path.Join("internal", pkgName)).AddFile(
		[]*task.FileDesc{
			{
				Path: path.Join("port", "http.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServiceWebHandler(buf, projectName, authorName, pkgName, svcDeclName, handlerDeclName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("service", "service.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServiceService(buf, projectName, authorName, pkgName, svcDeclName, repoDeclName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("repository", "pgsql.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServiceRepo(buf, projectName, authorName, pkgName, repoDeclName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("domain", fmt.Sprintf("%s.%s", pkgName, constants.GoSuffix)),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServiceDomainModel(buf, modelName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: path.Join("domain", "type.go"),
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServiceDomainType(buf, svcDeclName, repoDeclName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
			{
				Path: "wire_set.go",
				Source: func(buf *bytes.Buffer) error {
					templates.WriteGoServiceWireSet(buf, projectName, authorName, pkgName, handlerDeclName, svcDeclName, repoDeclName)
					return nil
				},
				Transforms: []task.Pipeline{transformation.GoFormatSource},
			},
		}...).AddCommand(
		// commands.GoEntInit(modelName),
		commands.GoWire("./web"),
		commands.GoModTidy(),
	)
}

func (gga *GenerateGoSimpleService) updateServerRegistry() error {
	pkg := gga.PackageName
	rf := gga.RegistryFile
	am := util.NewGoAstModifier(rf)

	am.AddModification(NewAddWireSetMod(pkg))
	am.AddImport(fmt.Sprintf("%s/%s/%s/internal/%s", constants.GoGithubPrefix, gga.AuthorName, gga.ProjectName, pkg))
	fset, at, err := am.ParseAndModify()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(rf, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to update '%s': %w", rf, err)
	}
	return format.Node(out, fset, at)
}
