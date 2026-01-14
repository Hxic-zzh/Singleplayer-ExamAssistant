/*
$$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\
$$  __$$\  \____$$\ $$  __$$\ $$  __$$\
$$ /  $$ | $$$$$$$ |$$ /  $$ |$$$$$$$$ |
$$ |  $$ |$$  __$$ |$$ |  $$ |$$   ____|
$$$$$$$  |\$$$$$$$ |\$$$$$$$ |\$$$$$$$\
$$  ____/  \_______| \____$$ | \_______|
$$ |                $$\   $$ |
$$ |                \$$$$$$  |
\__|                 \______/
*/
package pages

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"TestFyne-1119/Pages/tools"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// DownloadItem 表示 JSON 中的每个下载项
type DownloadItem struct {
	Index       int    `json:"序号"`
	Enabled     bool   `json:"是否启用"`
	Description string `json:"题库描述"`
	URL         string `json:"下载链接"`
}

// DownloadPage 下载管理界面
func DownloadPage(window fyne.Window) fyne.CanvasObject {
	log.Println("DownloadPage: 初始化下载管理界面")

	// 创建进度条
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	// 创建状态标签
	statusLabel := widget.NewLabel(tools.GetLocalized("download_status_init"))

	// 创建滚动容器用于显示卡片
	scrollContainer := container.NewScroll(container.NewVBox())
	scrollContainer.SetMinSize(fyne.NewSize(800, 500))

	// 获取数据按钮
	var fetchBtn *widget.Button
	fetchBtn = widget.NewButton(tools.GetLocalized("download_fetch_btn"), func() {
		log.Println("DownloadPage: 点击获取数据按钮")
		progress.Show()
		statusLabel.SetText(tools.GetLocalized("download_status_downloading"))
		fetchBtn.Disable()

		go func() {
			log.Println("DownloadPage: 开始下载JSON文件")

			// 下载 JSON 文件
			data, err := downloadJSON()
			if err != nil {
				log.Printf("DownloadPage: 下载失败 - %v", err)
				// 在主线程中更新UI
				updateUI := func() {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   tools.GetLocalized("error"),
						Content: fmt.Sprintf(tools.GetLocalized("download_failed"), err),
					})
					statusLabel.SetText(fmt.Sprintf(tools.GetLocalized("download_failed"), err))
					progress.Hide()
					fetchBtn.Enable()
				}

				// 在主线程中更新UI
				fyne.DoAndWait(updateUI)
				return
			}

			// 解析 JSON
			log.Printf("DownloadPage: JSON下载成功，开始解析，数据大小: %d bytes", len(data))
			downloadItems, err := parseJSON(data)
			if err != nil {
				log.Printf("DownloadPage: JSON解析失败 - %v", err)
				// 在主线程中更新UI
				updateUI := func() {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   tools.GetLocalized("error"),
						Content: fmt.Sprintf(tools.GetLocalized("parse_failed"), err),
					})
					statusLabel.SetText(fmt.Sprintf(tools.GetLocalized("parse_failed"), err))
					progress.Hide()
					fetchBtn.Enable()
				}

				// 在主线程中更新UI
				fyne.DoAndWait(updateUI)
				return
			}

			// 在主线程中更新UI
			updateUI := func() {
				// 创建卡片
				log.Printf("DownloadPage: 开始创建卡片，共 %d 个项目", len(downloadItems))
				createDownloadCards(downloadItems, scrollContainer, window)

				statusLabel.SetText(fmt.Sprintf(tools.GetLocalized("download_success"), len(downloadItems)))
				log.Printf("DownloadPage: 界面更新完成，显示 %d 个项目", len(downloadItems))
				progress.Hide()
				fetchBtn.Enable()
			}

			// 在主线程中更新UI
			fyne.DoAndWait(updateUI)
		}()
	})

	// 顶部控制面板
	controlPanel := container.NewHBox(
		fetchBtn,
		progress,
		container.NewCenter(statusLabel),
	)

	// 主布局
	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(tools.GetLocalized("download_title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			controlPanel,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		scrollContainer,
	)
}

// downloadJSON 使用 wget_x64.exe 下载 JSON 文件
func downloadJSON() ([]byte, error) {
	// 获取当前文件目录，然后找到main.go所在目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("无法获取当前文件路径")
	}

	// 获取main.go目录（pages目录的父目录）
	mainDir := filepath.Dir(filepath.Dir(filename))
	log.Printf("downloadJSON: main.go目录为: %s", mainDir)

	// 构建 wget 路径
	wgetPath := filepath.Join(mainDir, "wget_x64.exe")
	log.Printf("downloadJSON: 查找wget路径: %s", wgetPath)

	// 检查 wget 是否存在
	if _, err := os.Stat(wgetPath); os.IsNotExist(err) {
		log.Printf("downloadJSON: wget_x64.exe不存在")
		return nil, fmt.Errorf("wget_x64.exe 不存在于: %s", wgetPath)
	}

	// 创建临时文件保存 JSON
	tempFile := filepath.Join(mainDir, "users_test_temp.json")
	log.Printf("downloadJSON: 临时文件路径: %s", tempFile)

	// 构建命令
	cmd := exec.Command(wgetPath, "https://codelearnsiteofzzhandpzy.top/users_test.json", "-O", tempFile)
	log.Printf("downloadJSON: 执行命令: %s", cmd.String())

	// 设置超时
	timeout := time.Second * 30
	done := make(chan error, 1)

	go func() {
		log.Println("downloadJSON: 开始执行wget命令")
		err := cmd.Run()
		if err != nil {
			// 尝试获取输出
			if output, outErr := cmd.CombinedOutput(); outErr == nil {
				log.Printf("downloadJSON: wget输出: %s", string(output))
			}
		}
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("downloadJSON: wget执行失败 - %v", err)
			return nil, fmt.Errorf("wget 执行失败: %v", err)
		}
		log.Println("downloadJSON: wget执行成功")
	case <-time.After(timeout):
		log.Println("downloadJSON: 下载超时")
		return nil, fmt.Errorf("下载超时 (超过 %v)", timeout)
	}

	// 读取下载的文件
	log.Println("downloadJSON: 读取临时文件")
	data, err := os.ReadFile(tempFile)
	if err != nil {
		log.Printf("downloadJSON: 读取临时文件失败 - %v", err)
		return nil, fmt.Errorf("读取下载文件失败: %v", err)
	}
	log.Printf("downloadJSON: 成功读取文件，大小: %d bytes", len(data))

	// 打印文件内容用于调试
	log.Printf("downloadJSON: 文件内容:\n%s", string(data))

	// 清理临时文件
	log.Println("downloadJSON: 清理临时文件")
	if err := os.Remove(tempFile); err != nil {
		log.Printf("downloadJSON: 删除临时文件失败 - %v", err)
	}
	log.Println("downloadJSON: 临时文件已清理")

	return data, nil
}

// parseJSON 解析 JSON 数据
func parseJSON(data []byte) ([]DownloadItem, error) {
	log.Println("parseJSON: 开始解析JSON")

	// 先尝试解析为数组
	var downloadItems []DownloadItem
	err := json.Unmarshal(data, &downloadItems)
	if err != nil {
		log.Printf("parseJSON: JSON解析错误 - %v", err)

		// 尝试调试：打印原始数据
		log.Printf("parseJSON: 原始JSON数据:\n%s", string(data))

		// 尝试另一种方式：先解析为interface{}看看结构
		var rawData interface{}
		if err2 := json.Unmarshal(data, &rawData); err2 == nil {
			log.Printf("parseJSON: JSON结构类型: %T", rawData)
		}

		return nil, fmt.Errorf("JSON 解析错误: %v", err)
	}

	log.Printf("parseJSON: 解析完成，找到 %d 个项目", len(downloadItems))

	// 打印每个项目的信息用于调试
	for i, item := range downloadItems {
		log.Printf("parseJSON: 项目 %d - 序号: %d, 启用: %v, 描述: %s, URL: %s",
			i+1, item.Index, item.Enabled, item.Description, item.URL)
	}

	return downloadItems, nil
}

// createDownloadCards 创建下载卡片
func createDownloadCards(items []DownloadItem, scrollContainer *container.Scroll, window fyne.Window) {
	createDownloadCardsImpl(items, scrollContainer, window)
}

// 真正的卡片创建逻辑
func createDownloadCardsImpl(items []DownloadItem, scrollContainer *container.Scroll, window fyne.Window) {
	log.Println("createDownloadCards: 开始创建卡片")

	// 清空现有内容
	scrollContainer.Content.(*fyne.Container).Objects = nil

	// 为每个项目创建卡片
	for _, item := range items {
		item := item // 创建局部变量副本

		log.Printf("createDownloadCards: 创建卡片 - 序号: %d, 描述: %s", item.Index, item.Description)

		// 检查是否启用，如果不启用则跳过
		if !item.Enabled {
			log.Printf("createDownloadCards: 卡片 %d 未启用，跳过", item.Index)
			continue
		}

		// 格式化描述（确保不为空）
		description := item.Description
		if description == "" {
			description = fmt.Sprintf(tools.GetLocalized("download_item_title"), item.Index)
		}

		// 格式化 URL（智能省略中间部分）
		urlText := item.URL
		if len(urlText) > 50 {
			urlText = urlText[:20] + "..." + urlText[len(urlText)-20:]
		}
		log.Printf("createDownloadCards: URL原始长度 %d, 格式化后: %s", len(item.URL), urlText)

		// 创建状态标签
		statusText := tools.GetLocalized("download_enabled")
		if !item.Enabled {
			statusText = tools.GetLocalized("download_disabled")
		}
		statusLabel := widget.NewLabel(fmt.Sprintf("%s: %s", tools.GetLocalized("download_status_label"), statusText))

		// 创建 URL 标签（不可点击）
		urlLabel := widget.NewLabel(urlText)
		urlLabel.Wrapping = fyne.TextWrapWord

		// 创建打开按钮
		openBtn := widget.NewButton(tools.GetLocalized("download_open_link"), func() {
			log.Printf("createDownloadCards: 点击打开链接 - %s", item.URL)
			// 在新标签页/浏览器中打开链接
			err := openBrowser(item.URL)
			if err != nil {
				log.Printf("createDownloadCards: 打开链接失败 - %v", err)
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   tools.GetLocalized("error"),
					Content: fmt.Sprintf(tools.GetLocalized("open_link_failed"), err),
				})
			} else {
				log.Println("createDownloadCards: 链接已打开")
			}
		})

		// 创建卡片容器
		card := widget.NewCard(
			fmt.Sprintf(tools.GetLocalized("download_item_title"), item.Index),
			description,
			container.NewVBox(
				statusLabel,
				container.NewBorder(
					nil, nil,
					widget.NewLabel(tools.GetLocalized("download_link_label")), nil,
					urlLabel,
				),
				container.NewHBox(
					widget.NewLabel(""), // 占位符
					container.NewCenter(openBtn),
				),
			),
		)

		// 添加到滚动容器
		scrollContainer.Content.(*fyne.Container).Add(card)
	}

	// 如果没有启用的项目
	enabledCount := 0
	for _, item := range items {
		if item.Enabled {
			enabledCount++
		}
	}

	if enabledCount == 0 {
		log.Println("createDownloadCards: 没有启用的项目，显示无数据提示")
		noDataLabel := widget.NewLabel(tools.GetLocalized("download_no_enabled"))
		noDataLabel.Alignment = fyne.TextAlignCenter
		scrollContainer.Content.(*fyne.Container).Add(noDataLabel)
	}

	scrollContainer.Refresh()
	log.Printf("createDownloadCards: 卡片创建完成，显示 %d 个启用的项目", enabledCount)
}

// openBrowser 在浏览器中打开链接
func openBrowser(url string) error {
	log.Printf("openBrowser: 尝试打开链接 - %s", url)

	// 检查URL是否有效
	if url == "" {
		log.Println("openBrowser: URL为空")
		return fmt.Errorf("URL为空")
	}

	// 使用系统命令打开浏览器
	var cmd *exec.Cmd

	switch {
	case strings.Contains(strings.ToLower(url), "http://") || strings.Contains(strings.ToLower(url), "https://"):
		// 对于 Windows
		cmd = exec.Command("cmd", "/c", "start", "", url)
		log.Printf("openBrowser: 执行命令 - %s", cmd.String())
	default:
		log.Printf("openBrowser: 不支持的URL格式 - %s", url)
		return fmt.Errorf("不支持的URL格式")
	}

	err := cmd.Start()
	if err != nil {
		log.Printf("openBrowser: 启动浏览器失败 - %v", err)
		return err
	}

	log.Println("openBrowser: 浏览器命令已启动")
	return nil
}
