package common

import (
	"songjiang/sign"
)

// Config 程序整体配置
type Config struct {
	Songjiang Songjiang  `yaml:"songjiang"`
	Hao4k     sign.Hao4k `yaml:"hao4k"`
	Apps      []App      `yaml:"apps"`
}

// Songjiang Songjiang的配置
type Songjiang struct {
	Debug          bool   `yaml:"debug"`
	LogLevel       string `yaml:"logLevel"`
	TimeFormat     string `yaml:"timeFormat"`
	BrowserWidth   int    `yaml:"browserWidth"`
	BrowserHeight  int    `yaml:"browserHeight"`
	BrowserTimeout string `yaml:"browserTimeout"`
	Redo           string `yaml:"redo"`
}

// App 描述一个可以被自动签到的应用
type App struct {
	Type   string `yaml:"type"`
	Cookie string `yaml:cookie`
}

// UnmarshalYAML 反序列化App对象时的默认值
func (app *App) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawType App
	raw := rawType{
		Type:   "hao4k",
		Cookie: "",
	}
	if err := unmarshal(&raw); nil != err {
		return err
	}

	*app = App(raw)

	return nil
}

// UnmarshalYAML 反序列化Songjiang对象时的默认值处理
func (songjiang *Songjiang) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawType Songjiang
	raw := rawType{
		Debug:          false,
		LogLevel:       "info",
		BrowserHeight:  1080,
		BrowserWidth:   1920,
		BrowserTimeout: "15s",
		Redo:           "5s",
	}
	if err := unmarshal(&raw); nil != err {
		return err
	}

	*songjiang = Songjiang(raw)

	return nil
}
