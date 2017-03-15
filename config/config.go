package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	errUnsupportedOutputFormat = errors.New("unsupported output format")
)

type Configuration interface {
	Init(conf interface{}) error
}

type FlagConfiguration interface {
	Configuration
	SetCommandLineFlags(*flag.FlagSet)
}

// ConfigBase implements Configuration interface
type ConfigBase struct {
	// source of config, must be a filename or http URL
	SourceOfConfig string `json:"-" xml:"-" cli:"config-source" usage:"source of config, filename or http URL"`
	// output of config, exit process if non-empty
	OutputOfConfig string `json:"-" xml:"-" cli:"config-output" usage:"output of config, exit process if non-empty"`
}

func (c *ConfigBase) Init(conf interface{}) error {
	if c.SourceOfConfig != "" {
		var (
			reader io.Reader
			err    error
		)
		if strings.HasPrefix(c.SourceOfConfig, "http://") || strings.HasPrefix(c.SourceOfConfig, "https://") {
			var resp *http.Response
			resp, err = http.Get(c.SourceOfConfig)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					reader = resp.Body
				} else {
					err = errors.New(resp.Status)
				}
			}
		} else {
			var file *os.File
			file, err = os.Open(c.SourceOfConfig)
			if err == nil {
				defer file.Close()
				reader = file
			}
		}
		if reader != nil {
			err = json.NewDecoder(reader).Decode(conf)
		}
		if err != nil {
			return errors.New("read config from " + c.SourceOfConfig + ": " + err.Error())
		}
	}
	if c.OutputOfConfig != "" {
		var (
			format   = "json"
			filename = c.OutputOfConfig
			data     []byte
			err      error
		)
		out := strings.Split(c.OutputOfConfig, ":")
		if len(out) == 2 {
			// json:xxx
			if out[0] != "" {
				format = out[0]
			}
			filename = out[1]
		} else if len(out) == 1 {
			// yyy.json
			if i := strings.LastIndex(c.OutputOfConfig, "."); i > 0 && i+1 < len(c.OutputOfConfig) {
				format = c.OutputOfConfig[i+1:]
			}
		}
		switch format {
		case "json":
			data, err = json.MarshalIndent(conf, "", "    ")
		case "xml":
			data, err = xml.MarshalIndent(conf, "", "    ")
		default:
			err = errUnsupportedOutputFormat
		}
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filename, data, 0666); err != nil {
			return err
		}
		fmt.Printf("config output to file %s\n", filename)
		os.Exit(2)
	}
	return nil
}

// FlagConfig implements FlagConfiguration interface
type FlagConfig struct {
	ConfigBase
	sourceFlag string `json:"-" xml:"-" cli:"-"`
	outputFlag string `json:"-" xml:"-" cli:"-"`
}

func NewFlagConfig(sourceFlag, outputFlag string) *FlagConfig {
	return &FlagConfig{sourceFlag: sourceFlag, outputFlag: outputFlag}
}

func (c *FlagConfig) SetCommandLineFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&c.SourceOfConfig, c.sourceFlag, "", "source of config, filename or http URL")
	flagSet.StringVar(&c.OutputOfConfig, c.outputFlag, "", "output of config, exit process if non-empty")
}
