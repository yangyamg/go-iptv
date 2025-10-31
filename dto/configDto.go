package dto

type Build struct {
	Name    string `mapstructure:"name" json:"name" yaml:"name"`
	Package string `mapstructure:"package" json:"package" yaml:"package"`
	Sign    int64  `mapstructure:"sign" json:"sign" yaml:"sign"`
	Version string `mapstructure:"version" json:"version" yaml:"version"`
}

type AppUpdate struct {
	// Url  string `mapstructure:"url" json:"url" yaml:"url"`
	Set  int64  `mapstructure:"set" json:"set" yaml:"set"`
	Text string `mapstructure:"text" json:"text" yaml:"text"`
}

type App struct {
	NeedAuthor  int64 `mapstructure:"needauthor" json:"needauthor" yaml:"needauthor"`
	BuffTimeout int64 `mapstructure:"buff_time_out" json:"buff_time_out" yaml:"buff_time_out"`
	Decoder     int64 `mapstructure:"decoder" json:"decoder" yaml:"decoder"`
	// TrialDays   int64 `mapstructure:"trialdays" json:"trialdays" yaml:"trialdays"`
	// EPGApiChk     int64     `mapstructure:"epgapi_chk" json:"epgapi_chk" yaml:"epgapi_chk"`
	// MaxSameIPUser int64     `mapstructure:"max_sameip_user" json:"max_sameip_user" yaml:"max_sameip_user"`
	// IPCount       int64     `mapstructure:"ipcount" json:"ipcount" yaml:"ipcount"`
	Update AppUpdate `mapstructure:"update" json:"update" yaml:"update"`
}

type Tips struct {
	Loading       string `mapstructure:"loading" json:"loading" yaml:"loading"`
	UserExpired   string `mapstructure:"user_expired" json:"user_expired" yaml:"user_expired"`
	UserForbidden string `mapstructure:"user_forbidden" json:"user_forbidden" yaml:"user_forbidden"`
	UserNoReg     string `mapstructure:"user_noreg" json:"user_noreg" yaml:"user_noreg"`
}

type Ad struct {
	ShowTime     int64  `mapstructure:"showtime" json:"showtime" yaml:"showtime"`
	ShowInterval int64  `mapstructure:"showinterval" json:"showinterval" yaml:"showinterval"`
	AdText       string `mapstructure:"adtext" json:"adtext" yaml:"adtext"`
}

// type EPGErrors struct {
// 	Err1000 string `mapstructure:"err1000" json:"err1000" yaml:"err1000"`
// 	Err1001 string `mapstructure:"err1001" json:"err1001" yaml:"err1001"`
// 	Err1002 string `mapstructure:"err1002" json:"err1002" yaml:"err1002"`
// 	Err1003 string `mapstructure:"err1003" json:"err1003" yaml:"err1003"`
// 	Err1004 string `mapstructure:"err1004" json:"err1004" yaml:"err1004"`
// 	Err1005 string `mapstructure:"err1005" json:"err1005" yaml:"err1005"`
// }

// type Weather struct {
// 	ShowWea int64  `mapstructure:"showwea" json:"showwea" yaml:"showwea"`
// 	APIID   int64  `mapstructure:"api_id" json:"api_id" yaml:"api_id"`
// 	APIKey  string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`
// }

type ConfigChannel struct {
	Interval int64 `mapstructure:"interval" json:"interval" yaml:"interval"`
	Auto     int64 `mapstructure:"auto" json:"auto" yaml:"auto"`
}

// type Cache struct {
// 	Type  string `mapstructure:"type" json:"type" yaml:"type"`
// 	Redis Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
// }

type Redis struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Db       int    `mapstructure:"db" json:"db" yaml:"db"`
}

type Rss struct {
	Key string `mapstructure:"key" json:"key" yaml:"key"`
}

type Proxy struct {
	Status int64  `mapstructure:"status" json:"status" yaml:"status"`
	Port   int64  `mapstructure:"port" json:"port" yaml:"port"`
	PAddr  string `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type Config struct {
	ServerUrl string        `mapstructure:"server_url" json:"server_url" yaml:"server_url"`
	Build     Build         `mapstructure:"build" json:"build" yaml:"build"`
	App       App           `mapstructure:"app" json:"app" yaml:"app"`
	Tips      Tips          `mapstructure:"tips" json:"tips" yaml:"tips"`
	Ad        Ad            `mapstructure:"ad" json:"ad" yaml:"ad"`
	Channel   ConfigChannel `mapstructure:"channel" json:"channel" yaml:"channel"`
	Rss       Rss           `mapstructure:"rss" json:"rss" yaml:"rss"`
	Proxy     Proxy         `mapstructure:"proxy" json:"proxy" yaml:"proxy"`
	// Weather   Weather   `mapstructure:"weather" json:"weather" yaml:"weather"`
	// Cache     Cache     `mapstructure:"cache" json:"cache" yaml:"cache"`
	// EPGErrors EPGErrors `mapstructure:"epg_errors" json:"epg_errors" yaml:"epg_errors"`
}
