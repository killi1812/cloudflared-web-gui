package app

const (
	BuildDev  = "dev"
	BuildProd = "prod"
)

var (
	// Build describes app build type
	//
	// Could be: dev, prod
	Build          = BuildDev
	Version        = "0.0.0" // Version is a semver version of the app
	CommitHash     = "n/a"   // CommitHash is latest build commit hash
	BuildTimestamp = "n/a"   // BuildTimestamp stores when the app was build
)

// Envirment variables

var (
	Port       int    // Port is app port
	AccessKey  string // AccessKey is secrete for jwt access key
	RefreshKey string // RefreshKey is secrete for jwt refresh key
)
