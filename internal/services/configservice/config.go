package configservice

type Config struct {
	Forwards []*Forward `mapstructure:"forwards,omitempty" json:"forwards,omitempty"`
}

type Forward struct {
	Src     string `mapstructure:"src,omitempty" json:"src,omitempty"`
	Target  string `mapstructure:"target,omitempty" json:"target,omitempty"`
	Network string `mapstructure:"network,omitempty" json:"network,omitempty"`
}
