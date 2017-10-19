package configure

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

type NotationType int8

const (
	NotationTypeTOML NotationType = iota
	NotationTypeJSON
	NotationTypeYAML
)

type Configure struct {
	path       string
	userConfig interface{}

	opt *Option
}

type Option struct {
	NotationType NotationType
	SyncRealTime bool
	Editor       string
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

	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return nil, err
	}

	if opt == nil {
		var o Option
		opt = &o
	}

	conf := &Configure{
		path:       abs,
		userConfig: config,
		opt:        opt,
	}

	if err := conf.Load(); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Configure) Get() map[string]interface{} {
	return c.userConfig.(map[string]interface{})
}

func (c *Configure) Init() error {
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = encode(f, &c.userConfig, c.opt.NotationType)
	if err != nil {
		return err
	}
	return nil
}

func (c *Configure) Load() error {
	if !c.pathExist() {
		err := c.Init()
		if err != nil {
			return err
		}
	}

	f, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return decode(f, &c.userConfig, c.opt.NotationType)
}

func (c *Configure) Edit() error {
	if !c.pathExist() {
		err := c.Init()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(c.getEditor(), c.path)
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

func (c *Configure) pathExist() bool {
	_, err := os.Stat(c.path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Configure) getEditor() string {
	lookupEditor := func(editor string) string {
		p, err := exec.LookPath(editor)
		if err != nil && err.(*exec.Error).Err == exec.ErrNotFound {
			return "vim"
		}
		return p
	}

	if c.opt.Editor != "" {
		return lookupEditor(c.opt.Editor)
	}

	return lookupEditor(os.Getenv("EDITOR"))
}

func encode(w io.Writer, c *interface{}, t NotationType) error {
	switch t {
	case NotationTypeTOML:
		return toml.NewEncoder(w).Encode(c)
	case NotationTypeJSON:
		// For formatting
		d, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			return err
		}
		_, err = w.Write(d)
		return err
	case NotationTypeYAML:
		d, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		_, err = w.Write(d)
		return err
	default:
		return fmt.Errorf("unknown notation type: %T", t)
	}
}

func decode(r io.Reader, c *interface{}, t NotationType) error {
	switch t {
	case NotationTypeTOML:
		_, err := toml.DecodeReader(r, c)
		return err
	case NotationTypeJSON:
		err := json.NewDecoder(r).Decode(c)
		return err
	case NotationTypeYAML:
		d, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(d, c)
		return err
	default:
		return fmt.Errorf("unknown notation type: %T", t)
	}
}
