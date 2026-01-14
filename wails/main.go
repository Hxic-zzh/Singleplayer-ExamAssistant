package main

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// 创建应用实例
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "题目编辑器 - Question Bank Creator",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 处理 /tempwails/ 路径的请求
				if len(r.URL.Path) > 11 && r.URL.Path[:11] == "/tempwails/" {
					// 从 tempwails 目录提供文件
					filePath := filepath.Join(".", r.URL.Path[1:]) // 去掉前导斜杠
					if _, err := os.Stat(filePath); err == nil {
						http.ServeFile(w, r, filePath)
						return
					}
					http.NotFound(w, r)
					return
				}
				// 其他请求返回 nil，让默认的 asset handler 处理
				w.WriteHeader(http.StatusNotFound)
			}),
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		panic(err)
	}
}
