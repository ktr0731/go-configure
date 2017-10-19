# go-configure

configure anything more effectively, easily, and speedy!  

## Description
`go-configure` is a library which manage config files.  

## Requirements
- [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure)

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

// You need to use mapstructure to convert map[string]interface{} to the defined type
func Get() *Config {
  var config Config
  mapstructure.Decode(meta.Get(), &config)
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
