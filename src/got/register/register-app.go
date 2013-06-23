package register

import (
	"fmt"
)

// ________________________________________________________________________________
// Application Config
var Apps *AppConfigs = &AppConfigs{}

func RegisterApp(name string, displayName string, path string) {
	app := AppConfig{
		Name:        name,
		DisplayName: displayName,
		FilePath:    path,
	}
	Apps.Add(&app)
}

// ________________________________________________________________________________
// structs
//
type AppConfig struct {
	Name        string // application name, must be folder name of this app
	DisplayName string // Display Name
	FilePath    string // Absolute file path of this application, used to locate template.
}

type AppConfigs struct {
	Configs map[string]*AppConfig
}

func (c *AppConfigs) Add(config *AppConfig) {
	if c.Configs == nil {
		c.Configs = map[string]*AppConfig{}
	}
	c.Configs[config.Name] = config
}

func (c *AppConfigs) Get(appName string) *AppConfig {
	return c.Configs[appName]
}

func (c *AppConfig) String() string {
	return fmt.Sprintf("[APP:%v] '%v' - [Path:%v]",
		c.Name, c.DisplayName, c.FilePath)
}

func (c *AppConfigs) PrintALL() {
	fmt.Println("---- Apps ---------------------")
	for _, config := range c.Configs {
		fmt.Printf("  %v\n", config.String())
	}
}
