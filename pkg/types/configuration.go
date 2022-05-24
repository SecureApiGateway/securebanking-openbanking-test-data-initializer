package types

type Configuration struct {
	Environment environment `mapstructure:"ENVIRONMENT"`
	Hosts       hosts       `mapstructure:"HOSTS"`
	Users       users       `mapstructure:"USERS"`
	Namespace   string      `mapstructure:"NAMESPACE"`
}

type hosts struct {
	RsFQDN               string `mapstructure:"RS_FQDN"`
	IdentityPlatformFQDN string `mapstructure:"IDENTITY_PLATFORM_FQDN"`
	Scheme               string `mapstructure:"SCHEME"`
}
type environment struct {
	Verbose bool   `mapstructure:"VERBOSE"`
	Strict  bool   `mapstructure:"STRICT"`
	Type    string `mapstructure:"TYPE"`
	Paths   paths  `mapstructure:"PATHS"`
}

type paths struct {
	ConfigBaseDirectory    string `mapstructure:"CONFIG_BASE_DIRECTORY"`
	ConfigSecureBanking    string `mapstructure:"CONFIG_SECURE_BANKING"`
	ConfigIdentityPlatform string `mapstructure:"CONFIG_IDENTITY_PLATFORM"`
}

type users struct {
	FrPlatformAdminUsername string `mapstructure:"FR_PLATFORM_ADMIN_USERNAME"`
	FrPlatformAdminPassword string `mapstructure:"FR_PLATFORM_ADMIN_PASSWORD"`
	PsuUsername             string `mapstructure:"PSU_USERNAME"`
	PsuPassword             string `mapstructure:"PSU_PASSWORD"`
}
