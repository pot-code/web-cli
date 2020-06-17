package backend

// ConfigTemplate template for <repo>/config/config.go
var ConfigTemplate = `package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// running environment
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

// AppConfig App options
type AppConfig struct {
	Host           string        ` + "`" + `mapstructure:"host" json:"host" yaml:"host"` + "`" + `                                               // bind host address
	Port           int           ` + "`" + `mapstructure:"port" json:"port" yaml:"port" validate:"required"` + "`" + `                           // bind listen port
	Env            string        ` + "`" + `mapstructure:"env" json:"env" yaml:"env" validate:"oneof=development production,required"` + "`" + ` // runtime enviroment, development or production
	{{- range .ConfigGlobal}}
	{{Title .Name}} {{.Type}} ` + "`" + `mapstructure:"{{.Name}}" json:"{{.Name}}" yaml:"{{.Name}}"{{if .Required}} validate:"required"{{end}}` + "`" + ` // {{.Usage}}
	{{- end}}
	Database       struct {
		Driver   string ` + "`" + `mapstructure:"driver" json:"driver" yaml:"driver" validate:"required"` + "`" + `                      // driver name
		Host     string ` + "`" + `mapstructure:"host" json:"host" yaml:"host" validate:"required"` + "`" + `                            // server host
		MaxConn  int32  ` + "`" + `mapstructure:"maxconn" json:"maxconn" yaml:"maxconn" validate:"required,min=100"` + "`" + `           // maximum opening connections number
		Password string ` + "`" + `mapstructure:"password" json:"password" yaml:"password" validate:"required"` + "`" + `                // db password
		Port     int    ` + "`" + `mapstructure:"port" json:"port" yaml:"port"` + "`" + `                                                // server port
		Protocol string ` + "`" + `mapstructure:"protocol" json:"protocol" yaml:"protocol" validate:"omitempty,oneof=tcp udp"` + "`" + ` // connection protocol, eg.tcp
		Query    string ` + "`" + `mapstructure:"query" json:"query" yaml:"query"` + "`" + `                                             // DSN query parameter
		Schema   string ` + "`" + `mapstructure:"schema" json:"schema" yaml:"schema" validate:"required"` + "`" + `                      // use schema
		User     string ` + "`" + `mapstructure:"username" json:"username" yaml:"username" validate:"required"` + "`" + `                // db username
	} ` + "`" + `mapstructure:"database" json:"database" yaml:"database"` + "`" + `
	Logging struct {
		FilePath string ` + "`" + `mapstructure:"filePath" json:"filePath" yaml:"filePath"` + "`" + `                               // log file path
		Level    string ` + "`" + `mapstructure:"level" json:"level" yaml:"level" validate:"oneof=debug info warn error"` + "`" + ` // global logging level
		{{- range .ConfigLogging}}
		{{Title .Name}} {{.Type}} ` + "`" + `mapstructure:"{{.Name}}" json:"{{.Name}}" yaml:"{{.Name}}"{{if .Required}} validate:"required"{{end}}` + "`" + ` // {{.Usage}}
		{{- end}}
	} ` + "`" + `mapstructure:"logging" json:"logging" yaml:"logging"` + "`" + `
	Security struct {
		IDLength  int      ` + "`" + `mapstructure:"idLength" json:"idLength" yaml:"idLength" validate:"required"` + "`" + ` // length of generated ID for entities
		{{- range .ConfigSecurity}}
		{{Title .Name}} {{.Type}} ` + "`" + `mapstructure:"{{.Name}}" json:"{{.Name}}" yaml:"{{.Name}}"{{if .Required}} validate:"required"{{end}}` + "`" + ` // {{.Usage}}
		{{- end}}
	} ` + "`" + `mapstructure:"security" json:"security" yaml:"security"` + "`" + `
	APM struct {
		Enabled bool ` + "`" + `mapstructure:"enabled" json:"enabled" yaml:"enabled"` + "`" + `
	}` + "`" + `mapstructure:"apm" json:"apm" yaml:"apm"` + "`" + `
}` + `
// InitCallback callback function, it will be called after config is fully initialized
type InitCallback func(cfg *AppConfig)

// InitConfig populate AppConfig with cobra(flags), env variables(if set) and config file(if exists)
// and merge them if necessary. The priority is env>config>cobra(flags)
//
// It receives a callback function that will be called after the setup, and pass the config to it
func InitConfig(cb InitCallback) *cobra.Command {
	config := new(AppConfig)
	rootCmd := &cobra.Command{
		Use:   "{{.Usage}}",
		Short: "{{.Short}}",
		Run: func(cmd *cobra.Command, args []string) {
			// read file config
			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					log.Println("config file not exists, skipped")
				} else {
					log.Fatal(err)
				}
			}
			// export viper output to config
			if err := viper.Unmarshal(config); err != nil {
				log.Fatal(err)
			}
			if err := validateConfig(config); err != nil {
				log.Fatal(err)
			}
			log.Println("finished parsing config")
			if config.Logging.Level == "debug" {
				log.Println("==============================[app config]=============================")
				log.Printf("%+v\n", config)
				log.Println("=================================[end]=================================")
			}
			cb(config)
		},
	}

	// global
	rootCmd.Flags().StringVar(&config.Host, "host", "", "binding address")
	rootCmd.Flags().StringVar(&config.Env, "env", "development", "runtime enviroment, can be 'development' or 'production'")
	rootCmd.Flags().IntVar(&config.Port, "port", 8081, "host port")
	{{- range .ConfigGlobal}}
	rootCmd.Flags().{{GoTypeToCobra .Type}}(&config.{{Title .Name}}, "{{ToKebabCase .Name}}", {{GetValueString .Type .DefaultValue}}, "{{.Usage}}")
	{{- end}}

	// database
	rootCmd.Flags().StringVar(&config.Database.Driver, "database.driver", "mysql", "database driver to use")
	rootCmd.Flags().StringVar(&config.Database.Host, "database.host", "127.0.0.1", "database host")
	rootCmd.Flags().IntVar(&config.Database.Port, "database.port", 3306, "database server port")
	rootCmd.Flags().StringVar(&config.Database.Protocol, "database.protocol", "", "connection protocol(if mysql is used, this flag must be set), eg.tcp")
	rootCmd.Flags().StringVar(&config.Database.User, "database.username", "", "database username (required)")
	rootCmd.Flags().StringVar(&config.Database.Password, "database.password", "", "database password (required)")
	rootCmd.Flags().StringVar(&config.Database.Schema, "database.schema", "", "database schema to use (required)")
	rootCmd.Flags().StringVar(&config.Database.Query, "database.query", "", ` + "`" + `additional DSN query parameters('?' is auto prefixed), if you work with mysql and wish to
work with time.Time, you may specify "parseTime=true"` + "`" + `)
	rootCmd.Flags().Int32Var(&config.Database.MaxConn, "database.maxconn", 200, ` + "`" + `max connection count, if you encounter a "too many connections" error, please consider
increasing the max_connection value of your db server, or lower this value` + "`" + `)

	// logging
	rootCmd.Flags().StringVar(&config.Logging.Level, "logging.level", "info", "logging level")
	rootCmd.Flags().StringVar(&config.Logging.FilePath, "logging.file-path", "", "log to file")
	{{- range .ConfigLogging}}
	rootCmd.Flags().{{GoTypeToCobra .Type}}(&config.Logging.{{Title .Name}}, "logging.{{ToKebabCase .Name}}", {{GetValueString .Type .DefaultValue}}, "{{.Usage}}")
	{{- end}}

	// security
	rootCmd.Flags().IntVar(&config.Security.IDLength, "security.id-length", 24, "set length of generated ID for entities")
	{{- range .ConfigSecurity}}
	rootCmd.Flags().{{GoTypeToCobra .Type}}(&config.Security.{{Title .Name}}, "security.{{ToKebabCase .Name}}", {{GetValueString .Type .DefaultValue}}, "{{.Usage}}")
	{{- end}}

	// register viper config
	registerEnvVariables("{{.EnvPrefix}}")
	registerConfigFile()
	return rootCmd
}

func validateConfig(config *AppConfig) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" || name == "" {
			name = fld.Tag.Get("yaml")
			if name == "-" || name == "" {
				return ""
			}
		}
		return name
	})
	err := validate.Struct(config)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		log.Fatalf("Failed to validate config: %s", err)
	}
	if err == nil {
		return nil
	}

	var msg []string
	for _, field := range err.(validator.ValidationErrors) {
		namespace := field.Namespace()
		fieldName := namespace[strings.IndexByte(namespace, '.')+1:] // trim top level namespace
		switch field.Tag() {
		case "required":
			msg = append(msg, fmt.Sprintf("%s is required", fieldName))
		case "oneof":
			msg = append(msg, fmt.Sprintf("%s must be one of (%s)", fieldName, field.Param()))
		}
	}
	if len(msg) > 0 {
		return fmt.Errorf("failed to validate config: \n%s", strings.Join(msg, "\n"))
	}
	return nil
}

func registerConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
}

func registerEnvVariables(envPrefix string) {
	replacer := strings.NewReplacer(".", "_")
	template := AppConfig{}
	rtype := reflect.TypeOf(template)

	var result []string
	for i := 0; i < rtype.NumField(); i++ {
		f := rtype.Field(i)
		extractFieldsDFS(nil, &f, &result)
	}
	for _, field := range result {
		viper.BindEnv(field, strings.ToUpper(replacer.Replace(envPrefix+"."+field)))
	}
}

func extractFieldsDFS(namespace []string, field *reflect.StructField, result *[]string) {
	tagVal := field.Tag.Get("mapstructure")
	if field.Type.Kind() == reflect.Struct {
		rtype := field.Type
		for i := 0; i < rtype.NumField(); i++ {
			f := rtype.Field(i)
			if tagVal != "" {
				extractFieldsDFS(append(namespace, tagVal), &f, result)
			} else {
				extractFieldsDFS(namespace, &f, result)
			}
		}
	} else if tagVal != "" {
		*result = append(*result, strings.Join(append(namespace, tagVal), "."))
	}
}`

// MainTemplate template for <repo>/cmd/main.go
var MainTemplate = `package main

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/go-injection"
	"{{.ModuleName}}/internal/controller"
	"{{.ModuleName}}/internal/db"
	"{{.ModuleName}}/internal/infra"
	"{{.ModuleName}}/internal/middleware"
	"{{.ModuleName}}/internal/service"
	"go.elastic.co/apm/module/apmechov4"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	if err := InitConfig(startServer).Execute(); err != nil {
		log.Fatal(err)
	}
}

func startServer(option *AppConfig) {
	env := option.Env
	app := echo.New()
	app.HTTPErrorHandler = func(err error, c echo.Context) {}

	logger, err := infra.NewLogger(&infra.LoggingConfig{
		FilePath: option.Logging.FilePath,
		Level:    option.Logging.Level,
	})
	if err != nil {
		log.Fatalf("Failed to create logger: %s\n", err)
	}

	conn, err := db.GetDBConnection(&db.Config{
		Logger:   logger,
		User:     option.Database.User,
		Password: option.Database.Password,
		MaxConn:  option.Database.MaxConn,
		Protocol: option.Database.Protocol,
		Driver:   option.Database.Driver,
		Host:     option.Database.Host,
		Port:     option.Database.Port,
		Query:    option.Database.Query,
		Schema:   option.Database.Schema,
	})
	if err != nil {
		log.Fatalf("Failed to create DB connection: %s\n", err)
	}

	dic := injection.NewDIContainer()
	dic.Register(app)
	dic.Register(conn)
	dic.Register(option)
	dic.Register(logger)
	if err := dic.Populate(); err != nil {
		log.Fatal(err)
	}

	app.Use(middleware.ErrorHandling(
		&middleware.ErrorHandlingOption{
			Logger: logger,
		},
	))
	app.Use(middleware.Logging(logger))
	if env == EnvProduction {
		app.Use(apmechov4.Middleware())
	}
	if env == EnvDevelopment {
		app.Use(middleware.CORSMiddleware)
	}
	app.Use(middleware.AbortRequest(&middleware.AbortRequestOption{
		Timeout: 30 * time.Second,
	}))

	if err := app.Start(fmt.Sprintf("%s:%d", option.Host, option.Port)); err != nil {
		log.Fatal(err)
	}
}`
