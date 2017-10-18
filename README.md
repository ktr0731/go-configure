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

  // it is used in bellow
  configure *configure
}

var meta *configure.Configure

func New() *Config {
  return 
}

// I recommend define it for convenience  
func (c *Config) Get() *Config {
  return c.configure.Get().(*Config)
}
```

Next, 
``` go
config := Config{
  Bar: &Bar{
    Hoge: "hoge",
    Fuga: "fuga",
  },
  Foo: &Foo{
    Piyo: "piyo",
  },
}

config.configure, _ = configure.NewConfigure("~/.config.toml", &config, nil)

// Edit config by editor (Edit() opens the config file by editor)
conf.Edit()

// 
```
