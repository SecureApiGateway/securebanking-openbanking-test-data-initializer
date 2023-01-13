package types

import "fmt"

func ToStr(config Configuration) string {
	return fmt.Sprintf("Config is %#v", config)
}

type Configuration struct {
	Environment environment `mapstructure:"ENVIRONMENT"`
	Identity    identity    `mapstructure:"IDENTITY"`
	Hosts       hosts       `mapstructure:"HOSTS"`
	Users       users       `mapstructure:"USERS"`
	Namespace   string      `mapstructure:"NAMESPACE"`
}

type hosts struct {
	RsFQDN               string `mapstructure:"RS_FQDN"`
	IdentityPlatformFQDN string `mapstructure:"IDENTITY_PLATFORM_FQDN"`
	Scheme               string `mapstructure:"SCHEME"`
}

type identity struct {
	AmRealm string `mapstructure:"AM_REALM"`
}

type environment struct {
	Verbose bool   `mapstructure:"VERBOSE"`
	Strict  bool   `mapstructure:"STRICT"`
	Type    string `mapstructure:"TYPE"`
	Paths   paths  `mapstructure:"PATHS"`
}

type users struct {
	FrPlatformAdminUsername string `mapstructure:"FR_PLATFORM_ADMIN_USERNAME"`
	FrPlatformAdminPassword string `mapstructure:"FR_PLATFORM_ADMIN_PASSWORD"`
	PsuUserId               string `mapstructure:"PSU_USER_ID"`
	PsuUsername             string `mapstructure:"PSU_USERNAME"`
	PsuPassword             string `mapstructure:"PSU_PASSWORD"`
}

type paths struct {
	ConfigAuthHelper string `mapstructure:"CONFIG_AUTH_HELPER"`
}
