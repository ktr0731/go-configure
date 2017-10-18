# go-configure

configure anything more effectively, easily, and speedy!  

## Description
`go-configure` is a library which manage config files.  

## Usage
First of all, define the structure for config.  
``` go
package config 

type Foo struct {
  Hoge string
  Fuga int
}

type Bar struct {
  Piyo string
}

type Config struct {
  Foo *Foo
  Bar *Bar
}

var meta *configure.Configure

func init() {
  // (example) init config if you need it
  conf := Config{
    Foo: &Foo{
      Hoge: "hoge",
      Fuga: "fuga",
    },
    Bar: &Bar{
      Piyo: "piyo",
    },
  }

  // init go-configure
  var err error
  meta, err = configure.NewConfigure("~/.example.config.toml", conf, nil)
}

func Get() *Config {
  config, _ := meta.Get().(Config)
  return &config
}

func Edit() error {
  return meta.Edit()
}
```

Use or edit config
``` go
// Get config
config.Get()

// Edit config by editor (Edit() opens the config file by editor)
config.Edit()
```
