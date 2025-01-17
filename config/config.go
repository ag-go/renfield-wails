package config

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/jsawo/renfield/editor"
	"github.com/spf13/viper"
)

var mtx = &sync.Mutex{}

type Config struct {
	Currentproject string `mapstructure:"currentproject"`
	Projects       map[string]ProjectConfig
	Tags           []Tag
	TinkerTimeout  int // in seconds
}

type ProjectConfig struct {
	Id        string
	Name      string
	Path      string
	Tag       string
	Command   string
	Tinker    TinkerConfig
	JSONTools JSONToolsConfig
}

type TinkerConfig struct {
	Tabs      []editor.Tab
	ActiveTab string
}

type JSONToolsConfig struct {
	Tabs       []editor.Tab
	ActiveTab  string
	ActiveTool string
	Filter     string
}

type Tag struct {
	Label string
	Color string
}

var AppConfig Config

const (
	appDir      = "renfield"
	projectsDir = "projects"
	configFile  = "config"
	tempFileExt = ".tmp"
)

// Initialize creates the file and sets the default config values
func Initialize() {
	dirPath := GetAppConfigDir()
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: Error creating config directory:", err.Error())
	}

	setDefaults()

	err = viper.SafeWriteConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: Failed write to config file: \n", err.Error())
	}
}

// UpdateProject saves the project config to the config file
func UpdateProject(project ProjectConfig) {
	mtx.Lock()
	AppConfig.Projects[project.Id] = project
	Save(AppConfig)
	mtx.Unlock()
}

// Save saves the given config to the config file
func Save(config Config) {
	AppConfig = config
	viper.Set("currentproject", AppConfig.Currentproject)
	viper.Set("projects", AppConfig.Projects)
	viper.Set("tags", AppConfig.Tags)

	err := viper.WriteConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: Failed write to config file:", err.Error())
	}
}

// GetFreshConfig returns the current config from the config file
func GetFreshConfig() Config {
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("ERROR: unable to decode into struct: %w", err))
	}

	return AppConfig
}

func GetAppConfigDir() string {
	userConfigDir, _ := os.UserConfigDir()
	return path.Join(userConfigDir, appDir)
}

func GetWASMCachePath() string {
	return path.Join(GetAppConfigDir(), "cache_wasm")
}

func GetTempFilePath(tempName string) string {
	return path.Join(GetAppConfigDir(), projectsDir, GetCurrentProject().Id, tempName+tempFileExt)
}

func Load() {
	viper.SetConfigName(configFile)
	viper.SetConfigType("json")
	viper.AddConfigPath(GetAppConfigDir())

	AppConfig = Config{}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Initialize()
		} else {
			fmt.Fprint(os.Stderr, "ERROR:", err.Error())
			os.Exit(2)
		}
	}

	AppConfig = GetFreshConfig()

	if AppConfig.TinkerTimeout == 0 {
		AppConfig.TinkerTimeout = defaultTinkerTimeout
	}
}

func GetCurrentProject() ProjectConfig {
	return AppConfig.Projects[AppConfig.Currentproject]
}
