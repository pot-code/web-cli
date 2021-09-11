package constants

const (
	GoSuffix  = ".go"
	GoModFile = "go.mod"

	GoGithubPrefix = "github.com"

	GoApiServicePattern    = "%sUseCase"    // ModuleName+UseCase
	GoApiRepositoryPattern = "%sRepository" // ModuleName+Repository
	GoApiHandlerPattern    = "%sHandler"    // ModuleName+Handler
	GoModuleWireSePattern  = "%s.WireSet"   // packageName.WireSet

	GoConstructorPrefix = "New"
)