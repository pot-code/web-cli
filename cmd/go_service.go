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
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/constant"
	"github.com/pot-code/web-cli/internal/shell"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/pot-code/web-cli/internal/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type GoServiceConfig struct {
	Module string `arg:"0" alias:"module_name" validate:"required,var"` // go pkg name
}

var GoServiceCmd = command.NewCliCommand("service", "add a go service",
	&GoServiceConfig{},
	command.WithAlias([]string{"svc"}),
	command.WithArgUsage("module_name"),
).AddHandlers(
	&GenerateGoSimpleService{
		RegistryFile: path.Join("web", "server.go"),
	},
).BuildCommand()

type GenerateGoSimpleService struct {
	RegistryFile string
	PackageName  string
	ProjectName  string
	DomainName   string
	AuthorName   string
}

var _ command.CommandHandler = &GenerateGoSimpleService{}

func (gga *GenerateGoSimpleService) Handle(c *cli.Context, cfg interface{}) error {
	config := cfg.(*GoServiceConfig)
	pkgName := strings.ReplaceAll(config.Module, "-", "_")
	pkgName = strcase.ToSnake(pkgName)
	gga.PackageName = pkgName

	meta, err := util.ParseGoMod(constant.GoModFile)
	if err != nil {
		return errors.WithStack(err)
	}

	gga.DomainName = strcase.ToCamel(pkgName)
	gga.AuthorName = meta.GetAuthor()
	gga.ProjectName = meta.GetProject()

	dir := gga.PackageName
	if util.Exists(dir) {
		log.WithField("module", dir).Infof("module already exists")
		return nil
	}

	if err := gga.updateServerRegistry(); err != nil {
		return errors.Wrap(err, "failed to update registry")
	}

	return gga.generateFiles().Run()
}

func (gga *GenerateGoSimpleService) generateFiles() task.Task {
	modelName := gga.DomainName
	pkgName := gga.PackageName
	projectName := gga.ProjectName
	authorName := gga.AuthorName

	handlerDeclName := fmt.Sprintf(constant.GoApiHandlerPattern, modelName)
	svcDeclName := fmt.Sprintf(constant.GoApiServicePattern, modelName)
	repoDeclName := fmt.Sprintf(constant.GoApiRepositoryPattern, modelName)

	return task.NewSequentialExecutor(
		[]task.Task{
			task.NewParallelExecutor(
				task.BatchFileGenerationTask(
					task.NewFileGenerationTree(path.Join("internal", pkgName)).
						AddNodes( // root
							&task.FileGenerator{
								Name: "wire_set.go",
								Data: bytes.NewBufferString(templates.GoServiceWireSet(projectName, authorName, pkgName, handlerDeclName, svcDeclName, repoDeclName)),
							},
						).
						Branch("domain").
						AddNodes( // root/domain
							&task.FileGenerator{
								Name: fmt.Sprintf("%s.go", pkgName),
								Data: bytes.NewBufferString(templates.GoServiceDomainModel(modelName)),
							},
							&task.FileGenerator{
								Name: "type.go",
								Data: bytes.NewBufferString(templates.GoServiceDomainType(svcDeclName, repoDeclName)),
							},
						).Up().
						Branch("repository").
						AddNodes( // root/repository
							&task.FileGenerator{
								Name: "pgsql.go",
								Data: bytes.NewBufferString(templates.GoServiceRepo(projectName, authorName, pkgName, repoDeclName)),
							},
						).Up().
						Branch("port").
						AddNodes( // root/port
							&task.FileGenerator{
								Name: "http.go",
								Data: bytes.NewBufferString(templates.GoServiceWebHandler(projectName, authorName, pkgName, svcDeclName, handlerDeclName)),
							},
						).Up().
						Branch("service").
						AddNodes( // root/service
							&task.FileGenerator{
								Name: "service.go",
								Data: bytes.NewBufferString(templates.GoServiceService(projectName, authorName, pkgName, svcDeclName, repoDeclName)),
							},
						).Flatten(),
				),
			),
			shell.GoWire("./web"),
			shell.GoModTidy(),
		},
	)
}

func (gga *GenerateGoSimpleService) updateServerRegistry() error {
	pkg := gga.PackageName
	rf := gga.RegistryFile
	am := util.NewGoAstModifier(rf)

	am.AddModification(NewAddWireSetMod(pkg))
	am.AddImport(fmt.Sprintf("%s/%s/%s/internal/%s", constant.GoGithubPrefix, gga.AuthorName, gga.ProjectName, pkg))
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
