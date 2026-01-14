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
	"TestFyne-1119/Pages/tools"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CustomTheme 自定义主题
type CustomTheme struct {
	fyne.Theme
}

func (c *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{0x2E, 0x8B, 0x57, 0xFF} // 海绿色
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x69, G: 0x69, B: 0x69, A: 0x60} // 半透明灰色
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x69, G: 0x69, B: 0x69, A: 0x30} // 更透明的hover状态
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0x69, G: 0x69, B: 0x69, A: 0x40} // 按下状态
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x2E, G: 0x8B, B: 0x57, A: 0x40} // 焦点状态半透明
	case theme.ColorNameForeground:
		return color.Black // 设置前景色（文本颜色）为黑色
	default:
		return c.Theme.Color(name, variant)
	}
}

// 全局 Lightbox 查看器
var lightboxViewer *tools.LightboxViewer

// 全局题目映射表 - 用于广播定位
var questionMap = make(map[int]fyne.CanvasObject)

// 注册题目卡片到映射表
func registerQuestionCard(questionNumber int, card fyne.CanvasObject) {
	questionMap[questionNumber] = card
}

// 显示指定题目 - 修复版：使用正确的位置计算
func showQuestion(questionNumber int, scrollContainer *container.Scroll) {
	if card, exists := questionMap[questionNumber]; exists {
		// 正确的位置计算方法：
		// 获取卡片在内容容器中的相对位置，而不是相对于滚动容器的位置
		cardPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(card)
		contentPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(scrollContainer.Content)

		// 计算卡片相对于内容容器顶部的距离
		relativeY := cardPos.Y - contentPos.Y

		// 滚动到该位置
		scrollContainer.Offset = fyne.NewPos(0, relativeY)
		scrollContainer.Refresh()

		fmt.Printf("正确定位到第%d题，卡片位置:%.0f, 内容位置:%.0f, 滚动位置:%.0f\n",
			questionNumber, cardPos.Y, contentPos.Y, relativeY)
	} else {
		fmt.Printf("未找到第%d题的卡片\n", questionNumber)
	}
}

// 清理题目映射表
func clearQuestionMap() {
	questionMap = make(map[int]fyne.CanvasObject)
}

// createThemeDialog 创建自适应主题的对话框
func createThemeDialog(title, message string, window fyne.Window) dialog.Dialog {
	// 创建标签并设置颜色
	label := widget.NewLabel(message)

	// 检查当前主题，如果是深色主题则设置白色文字
	if app := fyne.CurrentApp(); app != nil {
		if app.Settings().ThemeVariant() == theme.VariantDark {
			label.TextStyle = fyne.TextStyle{Bold: true}
		} else {
			label.TextStyle = fyne.TextStyle{}
		}
	}

	return dialog.NewCustom(title, "确定", label, window)
}

// createConfirmDialog 创建自适应主题的确认对话框
func createConfirmDialog(title, message, confirmText, cancelText string, callback func(bool), window fyne.Window) dialog.Dialog {
	// 创建对话框内容
	var label fyne.CanvasObject
	// 检查当前主题，如果是深色主题则设置白色文字
	if app := fyne.CurrentApp(); app != nil && app.Settings().ThemeVariant() == theme.VariantDark {
		text := canvas.NewText(message, color.White)
		text.TextStyle = fyne.TextStyle{}
		label = text
	} else {
		text := canvas.NewText(message, color.Black)
		text.TextStyle = fyne.TextStyle{}
		label = text
	}

	content := container.NewVBox(label)
	return dialog.NewCustomConfirm(title, confirmText, cancelText, content, callback, window)
}

// MachinePage 生成题库选择页面
func MachinePage(window fyne.Window, contentStack *fyne.Container, mainContent fyne.CanvasObject, finalContent fyne.CanvasObject) fyne.CanvasObject {
	zipList, _ := tools.MachineToolListZips("data/Question")
	zipSelect := widget.NewSelect(zipList, nil)
	zipSelect.PlaceHolder = tools.GetLocalized("machine_select_zip_placeholder")

	progressLabelBinding := binding.NewString()
	progressLabel := widget.NewLabelWithData(progressLabelBinding)

	progressBarBinding := binding.NewFloat()
	progressBar := tools.NewPigProgressBar()
	progressBar.SetPosition(-400, 300) // 手动放置到指定坐标
	progressBar.SetAutoCenter(false)   // 禁用自动居中
	progressBar.Hide()
	// 确保初始状态干净（避免从上一次残留状态继续）
	progressBar.Reset()

	// 绑定进度变化到猪猪进度条
	progressBarBinding.AddListener(binding.NewDataListener(func() {
		v, _ := progressBarBinding.Get()
		progressBar.SetProgress(v)
	}))

	startBtn := widget.NewButton(tools.GetLocalized("machine_start_exam"), func() {
		if zipSelect.Selected == "" {
			fmt.Println(tools.GetLocalized("machine_select_zip_first"))
			widget.ShowPopUp(widget.NewLabel(tools.GetLocalized("machine_select_zip_first")), window.Canvas())
			return
		}

		// 每次开始都重置，防止上一轮动画/进度影响本轮
		progressBar.Reset()
		progressLabelBinding.Set(tools.GetLocalized("machine_preparing_exam"))
		progressBarBinding.Set(0)
		progressBar.Show()

		// 创建信号桥
		signalBridge := tools.NewMachineExamSignalBridge()

		// 监听进度
		go func() {
			for progress := range signalBridge.ProgressChan {
				fyne.DoAndWait(func() {
					progressLabelBinding.Set(tools.GetLocalized(progress.Stage))
					progressBarBinding.Set(float64(progress.Percent) / 100.0)
					if progress.Error != "" {
						fmt.Println("进度错误:", progress.Error)
						progressLabelBinding.Set(tools.GetLocalized("machine_error_prefix") + progress.Error)
					}
				})
			}
		}()

		// 监听完成
		go func() {
			<-signalBridge.DoneChan
			fyne.DoAndWait(func() {
				progressLabelBinding.Set(tools.GetLocalized("machine_ready"))
				progressBarBinding.Set(1)
				// 进度到1后把控件隐藏（否则后续 SetContent 切换时可能还在渲染树里可见）
				progressBar.Hide()
				progressBar.Reset()

				window.SetFullScreen(true)
				// 显示"加载中"界面
				loadingLabel := widget.NewLabel("正在加载考试界面，请稍候……")
				window.SetContent(container.NewVBox(
					container.NewCenter(loadingLabel),
				))
			})
			// 后台生成主界面
			examMain := NewExamMainPage(window, zipSelect.Selected, contentStack, mainContent, finalContent)
			// 平滑切换到主界面
			fyne.DoAndWait(func() {
				window.SetContent(examMain)
			})
		}()

		// 启动准备过程
		go func() {
			// 1. 清理ExamTemp
			err := tools.MachineToolClearExamTemp("data/Question/ExamTemp")
			if err != nil {
				signalBridge.SendProgress("machine_clear_examtemp_failed", 0, err.Error())
				return
			}
			signalBridge.SendProgress("machine_examtemp_cleared", 25, "")

			// 2. 解压zip到ExamTemp
			zipPath := "data/Question/" + zipSelect.Selected + ".zip"
			err = tools.MachineToolUnzip(zipPath, "data/Question/ExamTemp")
			if err != nil {
				signalBridge.SendProgress("machine_unzip_failed", 25, err.Error())
				return
			}
			signalBridge.SendProgress("machine_unzip_success", 50, "")

			// 3. 加载题库数据并生成正确答案
			examData, err := tools.MachineToolLoadExamData("data/Question/ExamTemp")
			if err != nil {
				signalBridge.SendProgress("machine_load_exam_failed", 50, err.Error())
				return
			}
			correctAnswers := tools.MachineToolGenerateCorrectAnswers(examData)
			err = tools.MachineToolSaveCorrectAnswers("data/Question/ExamTemp", correctAnswers)
			if err != nil {
				signalBridge.SendProgress("machine_save_correct_failed", 50, err.Error())
				return
			}
			signalBridge.SendProgress("machine_generate_correct", 75, "")

			time.Sleep(time.Second)
			signalBridge.SendProgress("machine_all_ready", 100, "")
			signalBridge.SendDone()
		}()
	})

	return container.NewVBox(
		widget.NewLabel(tools.GetLocalized("machine_exam_title")),
		zipSelect,
		startBtn,
		progressLabel,
		container.NewCenter(progressBar), // 居中显示猪猪进度条
	)
}

// NewExamMainPage 生成考试主界面（手动布局版本）
func NewExamMainPage(window fyne.Window, examName string, contentStack *fyne.Container, mainContent fyne.CanvasObject, finalContent fyne.CanvasObject) fyne.CanvasObject {
	// 设置自定义主题
	fyne.CurrentApp().Settings().SetTheme(&CustomTheme{theme.DefaultTheme()})

	// 获取当前分辨率配置
	resConfig := tools.GetCurrentResolutionConfig()

	examTempPath := "data/Question/ExamTemp"

	// 加载题库数据
	examData, err := tools.MachineToolLoadExamData(examTempPath)
	if err != nil {
		fmt.Println("加载题库失败:", err)
		return widget.NewLabel("加载题库失败: " + err.Error())
	}

	// 加载正确答案
	correctAnswers, err := tools.MachineToolLoadCorrectAnswers(examTempPath)
	if err != nil {
		fmt.Println("加载正确答案失败:", err)
		return widget.NewLabel("加载正确答案失败: " + err.Error())
	}

	// 获取所有题目（排除有hook的题目）
	questions := tools.MachineToolGetAllQuestions(examData)

	// 分组材料阅读题
	drGroups := tools.MachineToolGroupDRQuestions(examData)

	// 创建考试状态
	examState := tools.NewExamState(examName, questions, correctAnswers)
	examState.DRGroups = drGroups

	// 创建导航索引映射
	navIndexMap := make(map[string]int)
	currentNavIndex := 0

	// 为DR组分配导航索引
	for idx, drGroup := range drGroups {
		navIndexMap[drGroup.DRQuestion.ID] = idx
		currentNavIndex++
	}

	// 为独立题目分配导航索引
	for _, q := range questions {
		if q.Type == "DR" {
			continue
		}
		// 跳过有hook的题目（它们属于DR组）
		if q.Hook != "" {
			continue
		}
		navIndexMap[q.ID] = currentNavIndex
		currentNavIndex++
	}

	// 清理之前的题目映射表
	clearQuestionMap()

	// === 手动设置所有组件的大小和位置 ===

	// 主容器 - 透明背景板（维持坐标空间，不遮挡控件）
	background := canvas.NewRectangle(color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}) // 透明背景
	background.Resize(fyne.NewSize(5000, 5000))                                        // 设置足够大的背景

	mainContainer := container.NewWithoutLayout()
	mainContainer.Add(background) // 先添加背景

	// 1. 分割线 - 固定在左侧
	divider := canvas.NewRectangle(color.NRGBA{R: 0x87, G: 0xCE, B: 0xEB, A: 0xFF}) // 天蓝色分割线
	divider.Resize(fyne.NewSize(2, 2000))                                           // 设置足够高的分割线
	divider.Move(fyne.NewPos(resConfig.DividerX, 0))
	mainContainer.Add(divider)

	// 2. 倒计时标签 - 占据右侧顶部
	countdownBinding := binding.NewString()
	countdownLabel := widget.NewLabelWithData(countdownBinding)
	countdownLabel.Resize(fyne.NewSize(resConfig.ScrollWidth, 50)) // 使用配置的宽度
	countdownLabel.Move(fyne.NewPos(resConfig.DividerX+2, 0))      // 从分割线右侧开始
	// 设置倒计时标签样式
	countdownLabel.TextStyle = fyne.TextStyle{Bold: true}
	countdownLabel.Alignment = fyne.TextAlignCenter
	mainContainer.Add(countdownLabel)

	// 3. 题目导航区域 - 固定在左侧
	// 计算导航按钮数量：DR组 + 独立题目（排除DR题和有hook的题目）
	navButtonCount := len(drGroups)
	for _, q := range questions {
		if q.Type != "DR" && q.Hook == "" {
			navButtonCount++
		}
	}
	navButtons := make([]*widget.Button, navButtonCount)

	// 导航容器
	navContainer := container.NewWithoutLayout()
	navContainer.Resize(fyne.NewSize(resConfig.DividerX-50, 300)) // 使用配置的宽度
	navContainer.Move(fyne.NewPos(25, 60))                        // 固定在左侧位置
	mainContainer.Add(navContainer)

	// 4. 分页按钮
	prevBtn := createStyledButton(tools.GetLocalized("machine_prev"), func() {})
	prevBtn.Resize(fyne.NewSize(100, 40))
	prevBtn.Move(fyne.NewPos(25, 320)) // 在导航区域下方
	mainContainer.Add(prevBtn)

	nextBtn := createStyledButton(tools.GetLocalized("machine_next"), func() {})
	nextBtn.Resize(fyne.NewSize(100, 40))
	nextBtn.Move(fyne.NewPos(resConfig.DividerX-150, 320)) // 在导航区域下方右侧
	mainContainer.Add(nextBtn)

	// 5. 交卷和退出按钮 - 使用自适应主题的对话框
	submitBtn := createStyledButton(tools.GetLocalized("machine_submit"), func() {
		// 使用自适应主题的确认对话框
		d := createConfirmDialog(
			tools.GetLocalized("machine_submit_confirm_title"),
			tools.GetLocalized("machine_submit_confirm_content"),
			tools.GetLocalized("machine_submit_confirm_yes"),
			tools.GetLocalized("machine_submit_confirm_no"),
			func(confirmed bool) {
				if confirmed {
					submitExam(window, examState, contentStack, mainContent, finalContent, examData, examTempPath)
				}
			},
			window,
		)
		d.Show()
	})
	submitBtn.Resize(fyne.NewSize(100, 40))
	submitBtn.Move(fyne.NewPos(25, 380))
	mainContainer.Add(submitBtn)

	giveUpBtn := createStyledButton(tools.GetLocalized("machine_giveup"), func() {
		// 使用自适应主题的确认对话框
		d := createConfirmDialog(
			tools.GetLocalized("machine_giveup_confirm_title"),
			tools.GetLocalized("machine_giveup_confirm_content"),
			tools.GetLocalized("machine_giveup_confirm_yes"),
			tools.GetLocalized("machine_giveup_confirm_no"),
			func(confirmed bool) {
				if confirmed {
					tools.MachineToolClearExamTemp(examTempPath)
					window.SetFullScreen(false)
					window.Resize(fyne.NewSize(1280, 720))
					// 恢复金黄色主题
					fyne.CurrentApp().Settings().SetTheme(tools.NewGoldTheme())
					window.SetContent(finalContent)
					contentStack.Objects = []fyne.CanvasObject{MachinePage(window, contentStack, mainContent, finalContent)}
				}
			},
			window,
		)
		d.Show()
	})
	giveUpBtn.Resize(fyne.NewSize(100, 40))
	giveUpBtn.Move(fyne.NewPos(resConfig.DividerX-150, 380))
	mainContainer.Add(giveUpBtn)

	// 6. 图片展示区域 - 在交卷按钮下方
	imageDisplay := createImageDisplayArea(resConfig, window)
	imageDisplay.Move(fyne.NewPos(resConfig.ImageDisplayX, resConfig.ImageDisplayY))
	mainContainer.Add(imageDisplay)

	// 7. 右侧题目滚动区域 - 手动设置固定大小
	cardContainer := container.NewVBox()

	// 当前题目序号
	currentQuestionNumber := 1

	// 先添加材料阅读题组
	for idx, drGroup := range drGroups {
		drIdx := idx // 复制索引避免闭包问题
		// 构建图片路径
		imgPaths := make([]string, len(drGroup.DRQuestion.Images))
		for i, img := range drGroup.DRQuestion.Images {
			imgPaths[i] = filepath.Join(examTempPath, img)
		}

		// 构建子题的图片路径
		for i := range drGroup.ChildQuestions {
			childImgPaths := make([]string, len(drGroup.ChildQuestions[i].Images))
			for j, img := range drGroup.ChildQuestions[i].Images {
				childImgPaths[j] = filepath.Join(examTempPath, img)
			}
			drGroup.ChildQuestions[i].Images = childImgPaths
		}

		card := tools.DRGroupCard(drGroup.DRQuestion, drGroup.ChildQuestions, func(questionID string, answers []string) {
			// 直接保存小题答案到考试状态
			examState.SetAnswer(questionID, answers)
			// 更新导航按钮状态 - 使用DR题ID
			if navIdx, exists := navIndexMap[drGroup.DRQuestion.ID]; exists {
				updateNavButtonStyle(navButtons[navIdx], examState.IsDRGroupAnswered(examState.DRGroups[drIdx]))
			}
		}, window, tools.GetLocalized, resConfig.CardWidth)

		// 手动设置卡片大小 - 使用配置值
		card.Resize(fyne.NewSize(resConfig.CardWidth, resConfig.CardHeight))

		// 创建带背景的卡片容器
		cardBackground := canvas.NewRectangle(color.NRGBA{R: 0xF5, G: 0xF5, B: 0xDC, A: 0x80}) // 半透明米色
		cardBackground.Resize(fyne.NewSize(resConfig.CardWidth+20, resConfig.CardHeight+20))

		cardWithBg := container.NewWithoutLayout()
		cardWithBg.Add(cardBackground)
		cardWithBg.Add(card)
		cardWithBg.Resize(fyne.NewSize(resConfig.CardWidth+20, resConfig.CardHeight+20))

		cardContainer.Add(cardWithBg)

		// 注册卡片到全局映射表
		registerQuestionCard(currentQuestionNumber, cardWithBg)
		currentQuestionNumber++
	}

	// 添加其他独立题目（这些题目没有hook，在原来的区域显示）
	for _, q := range questions {
		// 跳过DR题（它们已经在上面处理了）
		if q.Type == "DR" {
			continue
		}
		// 跳过有hook的题目（它们在DR组中显示）
		if q.Hook != "" {
			continue
		}

		var card fyne.CanvasObject

		// 为每个卡片创建柔和绿色背景
		cardBackground := canvas.NewRectangle(color.NRGBA{R: 0xF5, G: 0xF5, B: 0xDC, A: 0x80}) // 半透明米色
		cardBackground.Resize(fyne.NewSize(resConfig.CardWidth+20, resConfig.CardHeight+20))

		// 构建图片路径
		imgPaths := make([]string, len(q.Images))
		for j, img := range q.Images {
			imgPaths[j] = filepath.Join(examTempPath, img)
		}

		switch q.Type {
		case "SC":
			if len(q.Options) > 0 {
				if len(q.Images) > 0 {
					card = tools.NewCardSingleChoiceWithImg(currentQuestionNumber, q.Question, q.Options, imgPaths, func(selected int) {
						if selected == -1 {
							examState.SetAnswer(q.ID, []string{})
						} else {
							ans := []string{q.Options[selected]}
							examState.SetAnswer(q.ID, ans)
						}
						if navIdx, exists := navIndexMap[q.ID]; exists {
							updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
						}
					}, window, tools.GetLocalized, resConfig.CardWidth)
				} else {
					card = tools.NewCardSingleChoice(currentQuestionNumber, q.Question, q.Options, func(selected int) {
						if selected == -1 {
							examState.SetAnswer(q.ID, []string{})
						} else {
							ans := []string{q.Options[selected]}
							examState.SetAnswer(q.ID, ans)
						}
						if navIdx, exists := navIndexMap[q.ID]; exists {
							updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
						}
					}, tools.GetLocalized)
				}
			}
		case "MC":
			if len(q.Options) > 0 {
				if len(q.Images) > 0 {
					card = tools.NewCardMultiChoiceWithImg(currentQuestionNumber, q.Question, q.Options, imgPaths, func(selected []int) {
						var ans []string
						for _, idx := range selected {
							ans = append(ans, q.Options[idx])
						}
						examState.SetAnswer(q.ID, ans)
						if navIdx, exists := navIndexMap[q.ID]; exists {
							updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
						}
					}, window, tools.GetLocalized, resConfig.CardWidth)
				} else {
					card = tools.NewCardMultiChoice(currentQuestionNumber, q.Question, q.Options, func(selected []int) {
						var ans []string
						for _, idx := range selected {
							ans = append(ans, q.Options[idx])
						}
						examState.SetAnswer(q.ID, ans)
						if navIdx, exists := navIndexMap[q.ID]; exists {
							updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
						}
					}, tools.GetLocalized)
				}
			}
		case "FL":
			// 填空题处理 - 使用支持换行的标签
			blanks := q.BlankCount
			entries := make([]*widget.Entry, blanks)
			box := container.NewVBox()
			// 设置填空题标题样式（多语言）
			titleLabel := widget.NewLabel(fmt.Sprintf(tools.GetLocalized("machine_card_title"), currentQuestionNumber))
			titleLabel.TextStyle = fyne.TextStyle{Bold: true}
			box.Add(titleLabel)

			// 使用支持换行的题目标签
			questionLabel := widget.NewLabel(q.Question)
			questionLabel.Wrapping = fyne.TextWrapWord
			box.Add(questionLabel)

			for k := 0; k < blanks; k++ {
				entries[k] = widget.NewEntry()
				entries[k].SetPlaceHolder(fmt.Sprintf(tools.GetLocalized("machine_fill_placeholder"), k+1))
				// 设置输入框样式
				entries[k].Validator = nil
				box.Add(entries[k])
				entries[k].OnChanged = func(s string) {
					ans := make([]string, blanks)
					for m, e := range entries {
						ans[m] = e.Text
					}
					examState.SetAnswer(q.ID, ans)
					if navIdx, exists := navIndexMap[q.ID]; exists {
						updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
					}
				}
			}
			card = container.NewBorder(nil, nil, nil, nil, box)

		// ===== 新题型处理 =====
		case "SCIMG":
			// 题干是图单选
			if len(q.Options) > 0 {
				card = tools.NewCardSCIMG(currentQuestionNumber, q.Question, q.Options, imgPaths, func(selected int) {
					if selected == -1 {
						examState.SetAnswer(q.ID, []string{})
					} else {
						ans := []string{q.Options[selected]}
						examState.SetAnswer(q.ID, ans)
					}
					if navIdx, exists := navIndexMap[q.ID]; exists {
						updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
					}
				}, window, tools.GetLocalized, resConfig.CardWidth)
			}
		case "MCIMG":
			// 题干是图多选
			if len(q.Options) > 0 {
				card = tools.NewCardMultiChoiceWithImg(currentQuestionNumber, q.Question, q.Options, imgPaths, func(selected []int) {
					var ans []string
					for _, idx := range selected {
						ans = append(ans, q.Options[idx])
					}
					examState.SetAnswer(q.ID, ans)
					if navIdx, exists := navIndexMap[q.ID]; exists {
						updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
					}
				}, window, tools.GetLocalized, resConfig.CardWidth)
			}
		case "FLIMG":
			// 题干是图填空
			card = tools.NewCardFLIMG(currentQuestionNumber, q.Question, q.BlankCount, imgPaths, func(answers []string) {
				examState.SetAnswer(q.ID, answers)
				if navIdx, exists := navIndexMap[q.ID]; exists {
					updateNavButtonStyle(navButtons[navIdx], examState.IsAnswered(q.ID))
				}
			}, window, tools.GetLocalized, resConfig.CardWidth)
		default:
			card = widget.NewLabel("未知题目类型: " + q.Type)
		}

		// 手动设置卡片大小 - 使用配置值
		card.Resize(fyne.NewSize(resConfig.CardWidth, resConfig.CardHeight))

		// 创建带背景的卡片容器
		cardWithBg := container.NewWithoutLayout()
		cardWithBg.Add(cardBackground)
		cardWithBg.Add(card)
		cardWithBg.Resize(fyne.NewSize(resConfig.CardWidth+20, resConfig.CardHeight+20))

		cardContainer.Add(cardWithBg)

		// 注册卡片到全局映射表 - 广播定位的关键！
		registerQuestionCard(currentQuestionNumber, cardWithBg)
		currentQuestionNumber++
	}

	// 设置卡片容器总高度，最小为2000
	totalCards := len(drGroups) + (len(questions) - len(drGroups)) // DR组数 + 独立题目数
	calculatedHeight := float32(totalCards) * resConfig.CardSpacing
	if calculatedHeight < 2000 {
		calculatedHeight = 2000
	}
	cardContainer.Resize(fyne.NewSize(resConfig.ScrollWidth, calculatedHeight))

	// 滚动容器 - 手动设置固定位置和大小
	scrollContainer := container.NewScroll(cardContainer)
	scrollContainer.Resize(fyne.NewSize(resConfig.ScrollWidth, resConfig.ScrollHeight)) // 使用配置的尺寸
	scrollContainer.Move(fyne.NewPos(resConfig.DividerX+2, 50))                         // 从分割线右侧开始，在倒计时下方
	mainContainer.Add(scrollContainer)

	// 创建导航按钮 - 总题目数 = DR组数 + 独立题目数（排除DR题和有hook的题目）
	independentQuestions := 0
	for _, q := range questions {
		if q.Type != "DR" && q.Hook == "" {
			independentQuestions++
		}
	}
	totalQuestions := len(drGroups) + independentQuestions
	buttonsPerPage := 20
	totalPages := (totalQuestions + buttonsPerPage - 1) / buttonsPerPage
	currentPage := 0

	for i := 0; i < totalQuestions; i++ {
		btn := createStyledButton(fmt.Sprintf("%d", i+1), func(idx int) func() {
			return func() {
				// 使用修复后的广播定位机制
				showQuestion(idx+1, scrollContainer)
			}
		}(i))
		btn.Resize(fyne.NewSize(resConfig.NavBtnWidth, resConfig.NavBtnHeight)) // 使用配置的按钮尺寸
		navButtons[i] = btn
	}

	// 分页按钮功能
	prevBtn.OnTapped = func() {
		if currentPage > 0 {
			currentPage--
			updateNavPage(navContainer, navButtons, currentPage, buttonsPerPage, resConfig)
		}
	}
	nextBtn.OnTapped = func() {
		if currentPage < totalPages-1 {
			currentPage++
			updateNavPage(navContainer, navButtons, currentPage, buttonsPerPage, resConfig)
		}
	}

	// 初始显示第一页
	updateNavPage(navContainer, navButtons, currentPage, buttonsPerPage, resConfig)

	// 倒计时goroutine
	go func() {
		endTime := examState.StartTime.Add(2 * time.Hour)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			remaining := time.Until(endTime)
			if remaining <= 0 {
				fyne.DoAndWait(func() {
					submitExam(window, examState, contentStack, mainContent, finalContent, examData, examTempPath)
				})
				return
			}
			hours := int(remaining.Hours())
			minutes := int(remaining.Minutes()) % 60
			seconds := int(remaining.Seconds()) % 60
			countdownBinding.Set(fmt.Sprintf(tools.GetLocalized("machine_time_left"), hours, minutes, seconds))
		}
	}()

	// 创建底层透明背景板
	transBg := canvas.NewRectangle(color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00})
	transBg.Resize(fyne.NewSize(5000, 5000))

	// 根据分辨率与主题选择背景图片
	bgImg := selectWallpaperByConfig(resConfig)
	// 设置背景图片在全屏(0,0)并填充
	bgImg.FillMode = canvas.ImageFillStretch
	bgImg.Move(fyne.NewPos(0, 0))
	bgImg.Resize(fyne.NewSize(float32(resConfig.Width), float32(resConfig.Height)))

	// 使用 Stack：透明背景在底层，背景图片其上，前景控件再上层
	window.SetFullScreen(true)
	return container.NewStack(transBg, bgImg, mainContainer)
}

// createStyledButton 创建样式化的按钮
func createStyledButton(text string, tapped func()) *widget.Button {
	btn := widget.NewButton(text, tapped)
	// 使用标准库颜色，设置按钮重要性
	btn.Importance = widget.HighImportance
	return btn
}

// updateNavButtonStyle 更新导航按钮样式
func updateNavButtonStyle(btn *widget.Button, answered bool) {
	if answered {
		// 已回答的题目使用暗灰色
		btn.Importance = widget.MediumImportance
	} else {
		// 未回答的题目使用海绿色（通过高重要性）
		btn.Importance = widget.HighImportance
	}
	btn.Refresh()
}

// updateNavPage 更新导航按钮分页 - 使用分辨率配置
func updateNavPage(navContainer *fyne.Container, navButtons []*widget.Button, currentPage, buttonsPerPage int, resConfig tools.ResolutionConfig) {
	navContainer.Objects = nil
	start := currentPage * buttonsPerPage
	end := start + buttonsPerPage
	if end > len(navButtons) {
		end = len(navButtons)
	}
	for i := start; i < end; i++ {
		row := (i - start) / 5 // 5列
		col := (i - start) % 5
		btn := navButtons[i]
		btn.Move(fyne.NewPos(float32(col)*resConfig.BtnSpacingX, float32(row)*resConfig.BtnSpacingY)) // 使用配置的间距
		navContainer.Add(btn)
	}
	navContainer.Refresh()
}

// submitExam 交卷处理 - 已修改以支持自适应主题
func submitExam(window fyne.Window, examState *tools.ExamState, contentStack *fyne.Container, mainContent fyne.CanvasObject, finalContent fyne.CanvasObject, originalExamData *tools.ExamData, examTempPath string) {
	finalScore, totalScore, correctCount, wrongCount := examState.CalculateScore()

	// 显示等待对话框 - 使用自适应主题
	var waitLabel fyne.CanvasObject
	if app := fyne.CurrentApp(); app != nil && app.Settings().ThemeVariant() == theme.VariantDark {
		waitLabel = canvas.NewText(tools.GetLocalized("machine_wait_wrong_set_label"), color.White)
	} else {
		waitLabel = canvas.NewText(tools.GetLocalized("machine_wait_wrong_set_label"), color.Black)
	}
	waitDialog := dialog.NewCustom(tools.GetLocalized("machine_wait_wrong_set"), tools.GetLocalized("machine_wait_wrong_set_cancel"), waitLabel, window)
	waitDialog.Show()

	// 在goroutine中执行保存和错题集生成
	go func() {
		// 保存考试结果和生成错题集
		err := examState.SaveExamResult("data/grades", originalExamData, examTempPath)

		// 在主线程中处理结果
		fyne.DoAndWait(func() {
			waitDialog.Hide()

			if err != nil {
				// 错误对话框使用自适应主题
				errorDialog := createThemeDialog(tools.GetLocalized("machine_error_title"),
					tools.GetLocalized("machine_save_result_failed")+err.Error(),
					window)
				errorDialog.Show()
				return
			}

			// 成功生成错题集后清理临时文件
			tools.MachineToolClearExamTemp("data/Question/ExamTemp")

			window.SetFullScreen(false)
			window.Resize(fyne.NewSize(1280, 720))

			// 恢复金黄色主题
			fyne.CurrentApp().Settings().SetTheme(tools.NewGoldTheme())

			fmt.Printf("考试结果: 得分 %.1f/%.1f, 正确题数 %d, 错题数 %d\n", finalScore, totalScore, correctCount, wrongCount)

			// 创建考试结果对话框 - 使用自适应主题
			resultContent := fmt.Sprintf(tools.GetLocalized("machine_exam_result_content"),
				finalScore, totalScore, correctCount, wrongCount)

			resultDialog := createThemeDialog(
				tools.GetLocalized("machine_exam_result_title"),
				resultContent,
				window,
			)

			resultDialog.Show()

			window.SetContent(finalContent)
			contentStack.Objects = []fyne.CanvasObject{MachinePage(window, contentStack, mainContent, finalContent)}
		})
	}()
}

// createImageDisplayArea 创建图片展示区域
func createImageDisplayArea(resConfig tools.ResolutionConfig, window fyne.Window) fyne.CanvasObject {
	// 随机选择图片
	imgPath, err := getRandomImageFromAssets()
	if err != nil {
		fmt.Printf("图片加载错误: %v\n", err)

		// 如果找不到图片，显示默认占位符
		placeholder := canvas.NewRectangle(color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF})
		placeholder.Resize(fyne.NewSize(resConfig.ImageSize, resConfig.ImageSize))
		placeholderLabel := widget.NewLabel(fmt.Sprintf("无图片\n%s", err))
		placeholderLabel.Alignment = fyne.TextAlignCenter

		placeholderContainer := container.NewStack(
			placeholder,
			container.NewCenter(placeholderLabel),
		)
		placeholderContainer.Resize(fyne.NewSize(resConfig.ImageSize, resConfig.ImageSize))
		return placeholderContainer
	}

	fmt.Printf("成功加载图片: %s\n", imgPath)

	// 加载并显示图片
	img := loadAndResizeImage(imgPath, resConfig.ImageSize, resConfig.ImageSize)

	// 创建可点击的图片容器
	clickableImg := createClickableImage(img, imgPath, window, resConfig.ImageSize)
	clickableImg.Resize(fyne.NewSize(resConfig.ImageSize, resConfig.ImageSize))
	return clickableImg
}

// createClickableImage 创建可点击的图片（使用透明按钮覆盖）
func createClickableImage(img *canvas.Image, imgPath string, window fyne.Window, size float32) fyne.CanvasObject {
	// 创建透明按钮覆盖在图片上
	btn := widget.NewButton("", nil)
	btn.Importance = widget.LowImportance

	// 设置按钮点击事件
	btn.OnTapped = func() {
		// 初始化 Lightbox 查看器（如果尚未初始化）
		if lightboxViewer == nil {
			lightboxViewer = tools.NewLightboxViewer(window)
		}

		// 显示图片查看器
		lightboxViewer.Show(imgPath)
	}

	// 设置按钮大小与图片相同
	btn.Resize(fyne.NewSize(size, size))

	// 创建堆叠容器：图片在下层，透明按钮在上层
	stack := container.NewWithoutLayout(img, btn)
	stack.Resize(fyne.NewSize(size, size))

	return stack
}

// getRandomImageFromAssets 从 assest/other 目录随机选择图片
func getRandomImageFromAssets() (string, error) {
	dirPath := "Pages/assest/other"

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return "", fmt.Errorf("目录不存在: %s", dirPath)
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("无法读取图片目录: %v", err)
	}

	// 过滤图片文件
	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() {
			name := strings.ToLower(file.Name())
			if strings.HasSuffix(name, ".png") ||
				strings.HasSuffix(name, ".jpg") ||
				strings.HasSuffix(name, ".jpeg") ||
				strings.HasSuffix(name, ".webp") ||
				strings.HasSuffix(name, ".gif") ||
				strings.HasSuffix(name, ".bmp") {
				fullPath := filepath.Join(dirPath, file.Name())
				imageFiles = append(imageFiles, fullPath)
			}
		}
	}

	if len(imageFiles) == 0 {
		return "", fmt.Errorf("目录 '%s' 中未找到图片文件", dirPath)
	}

	// 使用当前时间作为随机数种子，确保每次不同
	rand.New(rand.NewSource(time.Now().UnixNano()))
	selected := imageFiles[rand.Intn(len(imageFiles))]

	// 检查选中的文件是否存在
	if _, err := os.Stat(selected); os.IsNotExist(err) {
		return "", fmt.Errorf("选中的图片文件不存在: %s", selected)
	}

	return selected, nil
}

// loadAndResizeImage 加载并调整图片尺寸
func loadAndResizeImage(path string, width, height float32) *canvas.Image {
	img := tools.LoadImage(path)
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(width, height))
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}

// selectWallpaperByConfig 根据分辨率与主题选择对应的壁纸
func selectWallpaperByConfig(res tools.ResolutionConfig) *canvas.Image {
	variant := fyne.CurrentApp().Settings().ThemeVariant()
	var filename string
	if variant == theme.VariantDark {
		if res.Name == "1920x1080" {
			filename = "MachineB1920.png"
		} else {
			filename = "MachineB2560.png"
		}
	} else {
		if res.Name == "1920x1080" {
			filename = "Machinel1920.png"
		} else {
			filename = "Machinel2560.png"
		}
	}
	path := filepath.Join("images", filename)
	img := tools.LoadImage(path)
	return img
}
