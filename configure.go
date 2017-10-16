package configure

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	homedir "github.com/minio/go-homedir"
)

type NotationType int8

const (
	NotationTypeJSON NotationType = iota
	NotationTypeTOML
)

type Configure struct {
	path       string
	userConfig *interface{}

	opt *Option
}

type Option struct {
	NotationType NotationType
	SyncRealTime bool
	EnvVarPrefix string
}

func NewConfigure(fpath string, config interface{}, opt *Option) (*Configure, error) {
	var abs string
	if strings.HasPrefix(fpath, "~/") {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		abs = filepath.Join(home, strings.TrimPrefix(fpath, "~/"))
	} else {
		var err error
		abs, err = filepath.Abs(fpath)
		if err != nil {
			return nil, err
		}
	}

	if config == nil {
		var m *Option
		config = &m
	}

	return &Configure{
		path:       abs,
		userConfig: &config,
		opt:        opt,
	}, nil
}

func (c *Configure) Get() interface{} {
	return *c.userConfig
}

func (c *Configure) Init() error {
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = encode(f, c.userConfig, c.opt.NotationType)
	if err != nil {
		return err
	}
	return nil
}

func (c *Configure) Load() error {
	f, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return decode(f, c.userConfig, c.opt.NotationType)
}

func (c *Configure) Edit() error {
	_, err := os.Stat(c.path)
	if os.IsNotExist(err) {
		err = c.Init()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("nvim", c.path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if c.opt.SyncRealTime {
		if err := c.Load(); err != nil {
			return err
		}
	}

	return nil
}

func encode(w io.Writer, c *interface{}, t NotationType) error {
	switch t {
	case NotationTypeJSON:
		return json.NewEncoder(w).Encode(t)
	case NotationTypeTOML:
		return toml.NewEncoder(w).Encode(t)
	default:
		return fmt.Errorf("unknown notation type: %T", t)
	}
}

func decode(r io.Reader, c *interface{}, t NotationType) error {
	switch t {
	case NotationTypeJSON:
		return json.NewDecoder(r).Decode(c)
	case NotationTypeTOML:
		_, err := toml.DecodeReader(r, c)
		return err
	default:
		return fmt.Errorf("unknown notation type: %T", t)
	}
}
