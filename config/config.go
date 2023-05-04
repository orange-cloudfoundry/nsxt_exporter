package config

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type logConfig struct {
	Level   string `yaml:"level"`
	InJSON  bool   `yaml:"in_json"`
	NoColor bool   `yaml:"no_color"`
}

func (c *logConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain logConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.InJSON {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: c.NoColor,
		})
	}
	switch strings.ToUpper(c.Level) {
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	return nil
}

type Regexp struct {
	*regexp.Regexp
}

func (re *Regexp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	regex, err := regexp.Compile("^(?:" + s + ")$")
	if err != nil {
		return err
	}
	re.Regexp = regex
	return nil
}

type exporterConfig struct {
	Async                 bool          `yaml:"async"`
	IntervalDuration      time.Duration `yaml:"interval_duration"`
	ErrorIntervalDuration time.Duration `yaml:"error_interval_duration"`
	Port                  int           `yaml:"port"`
	Path                  string        `yaml:"path"`
	Namespace             string        `yaml:"namespace"`
}

func (c *exporterConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain exporterConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.Path == "" {
		c.Path = "/metrics"
	}
	if c.Port <= 0 {
		c.Port = 8080
	}
	if c.IntervalDuration == 0 {
		return fmt.Errorf("missing or zero key 'exporter.interval_duration'")
	}
	if c.ErrorIntervalDuration == 0 {
		return fmt.Errorf("missing or zero key 'exporter.error_interval_duration'")
	}
	return nil
}

type NSXConfig struct {
	URL            string   `yaml:"url"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
	ClientCertPath string   `yaml:"client_cert_path"`
	ClientKeyPath  string   `yaml:"client_key_path"`
	SkipSslVerify  bool     `yaml:"skip_ssl_verify"`
	CaCertPath     string   `yaml:"ca_cert_path"`
	MaxRetries     int      `yaml:"max_retries"`
	T0Filters      []string `yaml:"t0_filters"`
	T1Filters      []string `yaml:"t1_filters"`
	LBFilters      []string `yaml:"lb_filters"`
	VSFilters      []string `yaml:"vs_filters"`
}

func (n *NSXConfig) NeedPasswordLogin() bool {
	return len(n.ClientCertPath) == 0 || len(n.ClientKeyPath) == 0
}

func (n *NSXConfig) NSXHost() (string, error) {
	url, err := url.Parse(n.URL)
	if err != nil {
		return "", err
	}
	return url.Host, nil
}

func (n *NSXConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain NSXConfig
	if err := unmarshal((*plain)(n)); err != nil {
		return err
	}
	if n.URL == "" {
		return fmt.Errorf("missing mandatory key url")
	}
	if _, err := n.NSXHost(); err != nil {
		return fmt.Errorf("invalid url '%s' in key url", n.URL)
	}
	if (len(n.ClientCertPath) > 0 && len(n.ClientKeyPath) == 0) || (len(n.ClientCertPath) == 0 && len(n.ClientKeyPath) > 0) {
		return fmt.Errorf("one of {client_cert,client_key} keys are missing")
	}
	if (len(n.Username) > 0 && len(n.Password) == 0) || (len(n.Username) == 0 && len(n.Password) > 0) {
		return fmt.Errorf("one of {username,password} keys are missing")
	}
	if (len(n.Username) > 0 && len(n.ClientCertPath) > 0) || (len(n.Username) == 0 && len(n.ClientCertPath) == 0) {
		return fmt.Errorf("one and only one of {username,password} or {client_cert,client_key} should be given")
	}
	if n.MaxRetries == 0 {
		n.MaxRetries = 3
	}
	return nil
}

type Config struct {
	Log      *logConfig      `yaml:"log"`
	Exporter *exporterConfig `yaml:"exporter"`
	Nsxt     *NSXConfig      `yaml:"nsxt"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.Nsxt == nil {
		return fmt.Errorf("missing mandatory key nsxt.url, nsxt.username and nsxt.password")
	}
	if c.Exporter == nil {
		return fmt.Errorf("missing mandatory key exporter.interval_duration and exporter.error_interval_duration")
	}
	return nil
}

// NewConfig - Creates and validates config from given reader
func NewConfig(file io.Reader) *Config {
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("unable to read configuration file : %s", err)
	}
	config := Config{}
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Fatalf("Error when loading yaml config: %s", err)
	}
	return &config
}
