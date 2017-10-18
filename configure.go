package configure

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

type Configure struct {
	path         string
	userConfig   *interface{}
	syncRealTime bool
}

func NewConfigure(fpath string, config interface{}, syncRealTime bool) (*Configure, error) {
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

	return &Configure{
		path:         abs,
		userConfig:   &config,
		syncRealTime: syncRealTime,
	}, nil
}

func (c *Configure) Init() error {
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = toml.NewEncoder(f).Encode(c.userConfig)
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

	_, err = toml.DecodeReader(f, c.userConfig)
	return err
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

	if c.syncRealTime {
		if err := c.Load(); err != nil {
			return err
		}
	}

	return nil
}
