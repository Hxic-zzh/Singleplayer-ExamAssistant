package main

import (
	pages "TestFyne-1119/Pages"
	"TestFyne-1119/Pages/tools"
	"os"
	"path/filepath"

	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// === 全局变量声明 ===
var data []struct { // 给菜单栏用的，结构体之后要加什么属性也方便点
	title string
	icon  fyne.Resource
}

func loadIcons() { // 先将数组存起来
	iconPaths := []string{
		"images/message.png",
		"images/exam.png",
		"images/extract.png",
		"images/draw.png",
		"images/set.png",
	}

	for i, path := range iconPaths { // 到数组里面去找
		absPath, err := filepath.Abs(path) // 先设置成工作目录的绝对路径，之后可能改，打包的时候再说
		if err != nil {
			continue
		}
		fileData, err := os.ReadFile(absPath) // 赋值给两个变量，但是实际上只会按照顺序赋值给第一个变量，所以检查第二个变量是不是nil，是nil说明前一个变量成功获取了数据
		if err != nil {
			continue
		}

		data[i].icon = fyne.NewStaticResource(filepath.Base(path), fileData) // 创建Fyne可用的资源对象，相当于就是保存了引用类的二进制文件，方便后面的读取，也是必须要做的
	}
}

// cleanupTempDirs 清理临时目录
func cleanupTempDirs() {
	tempDirs := []string{
		"data/temp/tempImages",
		"data/temp/add",
		"data/output/outputTemp", // 新增：清理输出临时目录
		"wrong",                  // 新增：清理wrong文件夹
	}

	for _, dir := range tempDirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			log.Printf("获取绝对路径失败 %s: %v", dir, err)
			continue
		}

		// 检查目录是否存在
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			continue // 目录不存在，跳过
		}

		// 读取目录内容
		files, err := os.ReadDir(absPath)
		if err != nil {
			log.Printf("读取目录失败 %s: %v", absPath, err)
			continue
		}

		// 删除目录内的所有文件
		for _, file := range files {
			filePath := filepath.Join(absPath, file.Name())
			err := os.RemoveAll(filePath) // 使用 RemoveAll 可以删除文件和子目录
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
type MyCustomLayout struct { // 自定义布局初始化
	x, y, width, height float32 // 搞dpi屏幕可能有单位缩放
}

// ㊀layout方法：1. 向函数传入自定义实例 2.然后用layout类型亦是方法，江所有需要用的控件包含（CanvasObject）3. 然后就是要修改的控件属性内容
func (m *MyCustomLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	// objects[0] 是 list，objects[1] 是 contentStack
	objects[0].Resize(fyne.NewSize(m.width, m.height))
	objects[0].Move(fyne.NewPos(m.x, m.y))

	contentWidth := containerSize.Width - m.width - m.x - 20 // 窗口宽度 - 列表宽度 - 列表x坐标 - 右边距
	contentX := m.x + m.width + 20                           // 列表右边 + 间距
	objects[1].Resize(fyne.NewSize(contentWidth, m.height))
	objects[1].Move(fyne.NewPos(contentX, m.y))
}

// ㊁MinSize方法：返回容器的最小大小
func (m *MyCustomLayout) MinSize(objects []fyne.CanvasObject) fyne.Size { // 会不会看不懂？→ func是"函数声明"；(xxx)是函数接收的外部数据表示；XXX是函数名；(yyy)是函数的参数列表，有这些参数（传入的外部数据的参数亦或是属性）要调用；最后是返回值的类型
	return fyne.NewSize(0, 0) // 搞飞机啊，我直接锁死窗口大小，搞什么动态布局，我又不是大厂员工
}

//	   |   |
//	   |   |
//	   |   |
//	----   ----
//	 \       /
//	  \     /
//	   \   /
//	    \ /
//	     ^
//
// === 主函数 ===
func main() {
	// 启动前清理考试临时目录
	err := tools.MachineToolClearExamTemp("data/Question/ExamTemp")
	if err != nil {
		log.Printf("启动时清理ExamTemp失败: %v", err)
	}
	cleanupTempDirs()

	// 初始化本地化
	lang := tools.GetCurrentLanguage()
	err = tools.InitLocalization(lang)
	if err != nil {
		log.Printf("初始化本地化失败: %v", err)
	}

	// 使用翻译更新菜单标题
	data = []struct { // 给菜单栏用的，结构体之后要加什么属性也方便点
		title string
		icon  fyne.Resource
	}{
		{tools.GetLocalized("message"), nil}, // nil占位
		{tools.GetLocalized("exam"), nil},
		{tools.GetLocalized("extract"), nil},
		{tools.GetLocalized("detail"), nil},
		{tools.GetLocalized("settings"), nil},
	}

	loadIcons()

	myApp := app.NewWithID("io.Hxic.TestFyne-1119")
	myWindow := myApp.NewWindow("Hxic")
	myWindow.Resize(fyne.NewSize(1280, 720))
	myWindow.SetFixedSize(true)

	list := widget.NewList( // 模板建议保存，我其实看不懂，官方文档有，这官方的基础上找的别人成功的直接用了；总之就是三函数
		func() int {
			return len(data) // ㈠提高运行效率：运算变常量
		},
		func() fyne.CanvasObject { // ㈡创建容器和布局模板
			// 创建包含图标和RichText的"水平容器"
			icon := widget.NewIcon(nil)
			richText := widget.NewRichTextWithText("template") // RichText最类似label的底层逻辑，用来给控件修改字体

			// 设置RichText的字体样式
			if len(richText.Segments) > 0 {
				if textSeg, ok := richText.Segments[0].(*widget.TextSegment); ok {
					textSeg.Style = widget.RichTextStyle{
						SizeName: theme.SizeNameSubHeadingText, // 使用大号字体来修改字体的问题，我还没找到直接修改字号的方法
					}
				}
			}

			// 创建水平布局，按照顺位渲染
			return container.NewHBox(
				icon,
				container.NewPadded(richText),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) { // ㈢设置模板内具体样式
			// 获取容器
			hbox := o.(*fyne.Container)
			icon := hbox.Objects[0].(*widget.Icon)
			paddedContainer := hbox.Objects[1].(*fyne.Container)
			richText := paddedContainer.Objects[0].(*widget.RichText)

			// 设置图标（如果没有图标就设为nil）
			if data[i].icon != nil {
				icon.SetResource(data[i].icon)
			} else {
				icon.SetResource(nil) // 没有图标就显示空白
			}

			// 设置RichText内容 - 这里要直接修改现有的segment，而不是创建新的
			if len(richText.Segments) > 0 {
				if textSeg, ok := richText.Segments[0].(*widget.TextSegment); ok {
					textSeg.Text = data[i].title
					textSeg.Style = widget.RichTextStyle{
						SizeName: theme.SizeNameSubHeadingText, // 使用大号字体来修改字体的问题，我还没找到直接修改字号的方法
					}
				}
			} else {
				// 如果没有segment，创建新的
				segment := &widget.TextSegment{
					Text: data[i].title,
					Style: widget.RichTextStyle{
						SizeName: theme.SizeNameSubHeadingText, // 使用大号字体来修改字体的问题，我还没找到直接修改字号的方法
					},
				}
				richText.Segments = []widget.RichTextSegment{segment}
			}
			richText.Refresh() // 重要：刷新显示
		})

	contentStack := container.NewStack() // 创建堆叠面板

	myLayout := &MyCustomLayout{ // 自定义布局是指针类型，所以这里是引用符号，不要忘
		x:      20,
		y:      20,
		width:  220,
		height: 700, // 最后一行依旧需要","这是Go简化的体现
	}

	mainContent := container.New(myLayout, list, contentStack) // 主容器

	// 列表点击事件
	list.OnSelected = func(id widget.ListItemID) {
		var page fyne.CanvasObject
		switch id {
		case 0:
			// 确保使用默认主题
			myApp.Settings().SetTheme(theme.DefaultTheme())
			page = pages.MessagePage(myWindow)
		case 1:
			// machine 页面使用默认主题，考试界面会在内部设置主题
			myApp.Settings().SetTheme(theme.DefaultTheme())
			page = pages.MachinePage(myWindow, contentStack, mainContent)
		case 2:
			// 恢复默认主题
			myApp.Settings().SetTheme(theme.DefaultTheme())
			page = pages.ExtractPage(myWindow)
		case 3:
			// 恢复默认主题
			myApp.Settings().SetTheme(theme.DefaultTheme())
			page = pages.DrawPage(myWindow, mainContent, contentStack)
		case 4:
			// 恢复默认主题
			myApp.Settings().SetTheme(theme.DefaultTheme())
			page = pages.SettingsPage(myWindow)
		}
		contentStack.Objects = []fyne.CanvasObject{page}
		contentStack.Refresh()
	}

	myWindow.SetContent(mainContent)

	// 默认选择第一项
	list.Select(0)

	myWindow.ShowAndRun()
}
