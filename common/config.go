package common

import (
    "songjiang/sign"
)

// Config 程序整体配置
type Config struct {
    Songjiang Songjiang
    Hao4k     sign.Hao4k
    Apps      []App
}

// Songjiang 程序整体配置
type Songjiang struct {
    Debug         bool   `default:"false"`
    LogLevel      string `default:"info"`
    Chans         []ServerChan
    BrowserWidth  int    `default:"1920"`
    BrowserHeight int    `default:"1080"`
    Redo          string `default:"5s"`
}

// App 应用配置
type App struct {
    Name      string `default:"应用1"`
    Chans     []ServerChan
    Type      string `default:"hao4k"`
    Cookies   string
    StartTime string `default:"8:00"`
    EndTime   string `default:"13:00"`
}
