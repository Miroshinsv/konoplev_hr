package config

type ACL struct {
	ConfigPath string `env:"ACL_CONFIG_PATH" envDefault:"/opt/model.conf"`
	PolicyPath string `env:"ACL_POLICY_PATH" envDefault:"/opt/policy.csv"`
}
