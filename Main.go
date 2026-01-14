package main

import (
	pages "TestFyne-1119/Pages"
	"TestFyne-1119/Pages/tools"
	"image"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// === 全局变量声明 ===
var data []tools.ListItem
var preloadedPages [6]fyne.CanvasObject

// loadIcon 加载图标资源
func loadIcon(path string) fyne.Resource {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("加载图标失败 %s: %v", path, err)
		return nil
	}
	fileData, err := os.ReadFile(absPath)
	if err != nil {
		log.Printf("读取图标文件失败 %s: %v", path, err)
		return nil
	}
	return fyne.NewStaticResource(filepath.Base(path), fileData)
}

// getImageSize 获取图片尺寸
func getImageSize(path string) (fyne.Size, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fyne.Size{}, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fyne.Size{}, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return fyne.Size{}, err
	}

	return fyne.NewSize(float32(img.Width), float32(img.Height)), nil
}

// initData 初始化数据
func initData() {
	lang := tools.GetCurrentLanguage()
	err := tools.InitLocalization(lang)
	if err != nil {
		log.Printf("初始化本地化失败: %v", err)
	}

	iconPaths := []string{
		"images/message.png",
		"images/exam.png",
		"images/extract.png",
		"images/draw.png",
		"images/set.png",
		"images/download.png",
	}

	titles := []string{
		tools.GetLocalized("message"),
		tools.GetLocalized("exam"),
		tools.GetLocalized("extract"),
		tools.GetLocalized("detail"),
		tools.GetLocalized("settings"),
		tools.GetLocalized("download"),
	}

	data = make([]tools.ListItem, len(iconPaths))
	for i := 0; i < len(iconPaths); i++ {
		data[i] = tools.ListItem{
			Title: titles[i],
			Icon:  loadIcon(iconPaths[i]),
		}
	}
}

// cleanupTempDirs 清理临时目录
func cleanupTempDirs() {
	tempDirs := []string{
		"data/temp/tempImages",
		"data/temp/add",
		"data/output/outputTemp",
		"wrong",
	}

	for _, dir := range tempDirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			log.Printf("获取绝对路径失败 %s: %v", dir, err)
			continue
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			continue
		}

		files, err := os.ReadDir(absPath)
		if err != nil {
			log.Printf("读取目录失败 %s: %v", absPath, err)
			continue
		}

		for _, file := range files {
			filePath := filepath.Join(absPath, file.Name())
			err := os.RemoveAll(filePath)
			if err != nil {
				log.Printf("删除文件失败 %s: %v", filePath, err)
			} else {
				log.Printf("已删除: %s", filePath)
			}
		}

		log.Printf("清理完成: %s", absPath)
	}
}

// === 自定义布局（全局） ===
type MyCustomLayout struct {
	x, y, width, height float32
}

func (m *MyCustomLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Resize(fyne.NewSize(m.width, m.height))
	objects[0].Move(fyne.NewPos(m.x, m.y))

	contentWidth := containerSize.Width - m.width - m.x - 20
	contentX := m.x + m.width + 20
	objects[1].Resize(fyne.NewSize(contentWidth, m.height))
	objects[1].Move(fyne.NewPos(contentX, m.y))
}

func (m *MyCustomLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

// handlePageSelection 处理页面选择
func handlePageSelection(id int, myWindow fyne.Window, contentStack *fyne.Container, mainContent fyne.CanvasObject, finalContent fyne.CanvasObject, myApp fyne.App) {
	// 检查当前主题，如果不是金黄色主题，则设置
	if _, ok := myApp.Settings().Theme().(*tools.GoldTheme); !ok {
		myApp.Settings().SetTheme(tools.NewGoldTheme())
	}

	contentStack.Objects = []fyne.CanvasObject{preloadedPages[id]}
	contentStack.Refresh()
}

func main() {
	// 创建应用与主题
	myApp := app.NewWithID("io.Hxic.TestFyne-1119")
	myApp.Settings().SetTheme(tools.NewGoldTheme())
	// 显示闪屏（圆弧动画阶段）
	splash := tools.ShowStartupSplash(myApp)

	// 提前创建主窗口，但暂不显示
	myWindow := myApp.NewWindow("Hxic")

	// 后台执行：清理 + 初始化 + 预加载；完成后构建UI
	go func() {
		if err := tools.MachineToolClearExamTemp("data/Question/ExamTemp"); err != nil {
			log.Printf("启动时清理ExamTemp失败: %v", err)
		}
		cleanupTempDirs()
		initData()

		log.Println("开始构建UI...")

		fyne.Do(func() {
			// 创建带特效的自定义列表
			listConfig := tools.CustomListConfig{
				Items: data,
				OnSelected: func(id int) {
					// 处理选中事件
				},
				InitialSelected: 0, // 设置初始选中为第一项
			}

			// 加载背景图片并获取尺寸
			bgImage := loadIcon("images/listbackground.png")
			var bgSize fyne.Size
			if bgImage != nil {
				if size, err := getImageSize("images/listbackground.png"); err == nil {
					bgSize = size
				} else {
					bgSize = fyne.NewSize(220, 700) // 默认尺寸
				}
			}

			imageListConfig := tools.ImageListConfig{
				ListConfig: listConfig,
				Background: &tools.ListBackground{
					Image:     bgImage,
					ImageSize: bgSize,
					Opacity:   0.5,
				},
			}

			// 加载特效图片
			var effectImages [6]fyne.Resource
			var effectSizes [6]fyne.Size

			effectPaths := []string{
				"images/list1.png",
				"images/list2.png",
				"images/list3.png",
				"images/list4.png",
				"images/list5.png",
				"images/list6.png",
			}

			for i, path := range effectPaths {
				img := loadIcon(path)
				effectImages[i] = img

				if img != nil {
					if size, err := getImageSize(path); err == nil {
						effectSizes[i] = size
					} else {
						effectSizes[i] = fyne.NewSize(200, 50)
					}
				}
			}

			effectConfig := &tools.SelectionEffectConfig{
				Images:     effectImages,
				ImageSizes: effectSizes,
				DefaultPos: tools.EffectPosition{
					XOffset: 0,
					YOffset: 0,
					Scale:   1.0,
				},
				Positions: [6]*tools.EffectPosition{
					tools.NewEffectPosition(45, 2, 0.9),   // 第1项
					tools.NewEffectPosition(45, -3, 0.95), // 第2项
					tools.NewEffectPosition(0, -16, 0.85), // 第3项
					tools.NewEffectPosition(45, -25, 1.0), // 第4项
					tools.NewEffectPosition(20, -35, 1.1), // 第5项
					tools.NewEffectPosition(-45, -50, 5),  // 第6项
				},
			}

			// 创建带特效的列表
			customList := tools.NewImageBackgroundListWithEffect(imageListConfig, effectConfig)

			contentStack := container.NewStack()

			myLayout := &MyCustomLayout{
				x:      20,
				y:      20,
				width:  220,
				height: 700,
			}

			mainContent := container.New(myLayout, customList.CanvasObject, contentStack)

			// 加载主背景图片（按主题）
			var bgImageResource fyne.Resource
			if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantDark {
				bgImageResource = loadIcon("images/BlackBackGround.png")
			} else {
				bgImageResource = loadIcon("images/MainBackGround.png")
			}

			var mainBgCanvas *canvas.Image
			if bgImageResource != nil {
				mainBgCanvas = canvas.NewImageFromResource(bgImageResource)
				mainBgCanvas.FillMode = canvas.ImageFillStretch // 拉伸填充窗口
				mainBgCanvas.Resize(fyne.NewSize(1280, 720))    // 设置为窗口大小
			}

			// 创建包含背景的栈
			var finalContent fyne.CanvasObject
			if mainBgCanvas != nil {
				finalContent = container.NewStack(mainBgCanvas, mainContent)
			} else {
				finalContent = mainContent
			} // 预加载所有页面
			log.Println("开始预加载页面...")
			preloadedPages[0] = pages.MessagePage(myWindow)
			log.Println("MessagePage 加载完成")
			preloadedPages[1] = pages.MachinePage(myWindow, contentStack, mainContent, finalContent)
			log.Println("MachinePage 加载完成")
			preloadedPages[2] = pages.ExtractPage(myWindow)
			log.Println("ExtractPage 加载完成")
			preloadedPages[3] = pages.DrawPage(myWindow, mainContent, contentStack)
			log.Println("DrawPage 加载完成")
			preloadedPages[4] = pages.SettingsPage(myWindow)
			log.Println("SettingsPage 加载完成")
			preloadedPages[5] = pages.DownloadPage(myWindow)
			log.Println("DownloadPage 加载完成")

			// 设置列表选中事件
			customList.SetOnSelected(func(id int) {
				handlePageSelection(id, myWindow, contentStack, mainContent, finalContent, myApp)
			})

			// 手动触发初始页面显示
			handlePageSelection(0, myWindow, contentStack, mainContent, finalContent, myApp)

			// 设置主窗口内容与窗口属性
			myWindow.SetContent(finalContent)
			myWindow.Resize(fyne.NewSize(1280, 720))
			myWindow.SetFixedSize(true)
			myWindow.CenterOnScreen()
			myWindow.Content().Refresh()

			log.Println("UI构建完成，通知闪屏播放GIF...")

			// GIF播放完成后回调：进入主界面
			tools.SetSplashOnFinished(func() {
				myWindow.Content().Refresh()
				myWindow.Show()
				if splash != nil {
					splash.Close()
				}
				log.Println("主窗口已显示，闪屏GIF播放完成并关闭")
			})
			// 通知闪屏：主界面已准备好 -> 继续3秒圆弧 -> 播放GIF
			tools.NotifyReady()
		})
	}()

	// 进入事件循环
	myApp.Run()
}
