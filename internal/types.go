package internal

type IntroduceType struct {
	Status      string  `json:"status,omitempty"`
	Country     string  `json:"country,omitempty"`
	Countrycode string  `json:"countrycode,omitempty"`
	Region      string  `json:"region,omitempty"`
	Regionname  string  `json:"regionname,omitempty"`
	City        string  `json:"city,omitempty"`
	Zip         string  `json:"zip,omitempty"`
	Lat         float32 `json:"lat,omitempty"`
	Lon         float32 `json:"lon,omitempty"`
	Timezone    string  `json:"timezone,omitempty"`
	Isp         string  `json:"isp,omitempty"`
	Org         string  `json:"org,omitempty"`
	As          string  `json:"as,omitempty"`
	Query       string  `json:"query,omitempty"`
}

type AttackVector struct {
	Rate    float64 `json:"rate,omitempty" yaml:"rate"`
	Address string  `json:"address,omitempty" yaml:"address"`
	Method  string  `json:"method,omitempty" yaml:"method"`
}

type AttackConfig struct {
	Address string   `json:"address,omitempty" yaml:"address"`
	Methods []string `json:"methods,omitempty" yaml:"methods"`
}
