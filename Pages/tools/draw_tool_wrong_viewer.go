// 该文件专门写操作逻辑
/*



                                                    $$\     $$\
                                                    $$ |    \__|
 $$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\  $$$$$$\   $$\  $$$$$$\  $$$$$$$\
$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\  \____$$\ \_$$  _|  $$ |$$  __$$\ $$  __$$\
$$ /  $$ |$$ /  $$ |$$$$$$$$ |$$ |  \__| $$$$$$$ |  $$ |    $$ |$$ /  $$ |$$ |  $$ |
$$ |  $$ |$$ |  $$ |$$   ____|$$ |      $$  __$$ |  $$ |$$\ $$ |$$ |  $$ |$$ |  $$ |
\$$$$$$  |$$$$$$$  |\$$$$$$$\ $$ |      \$$$$$$$ |  \$$$$  |$$ |\$$$$$$  |$$ |  $$ |
 \______/ $$  ____/  \_______|\__|       \_______|   \____/ \__| \______/ \__|  \__|
          $$ |
          $$ |
          \__|


*/

package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// LoadWrongQuestionSet 加载错题集数据
func LoadWrongQuestionSet(jsonPath string) (*WrongQuestionSet, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("打开错题集文件失败: %v", err)
	}
	defer file.Close()

	var wrongSet WrongQuestionSet
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&wrongSet); err != nil {
		return nil, fmt.Errorf("解析错题集JSON失败: %v", err)
	}

	return &wrongSet, nil
}

// CreateWrongQuestionViewer 创建错题查看界面
func CreateWrongQuestionViewer(window fyne.Window, gradeFolderPath string, returnCallback func()) fyne.CanvasObject {
	// 查找并解压错题集ZIP文件
	sourseZipPath := filepath.Join(gradeFolderPath, "sourse.zip")
	wrongDir := "wrong"

	// 检查ZIP文件是否存在
	if _, err := os.Stat(sourseZipPath); os.IsNotExist(err) {
		return createErrorPage("未找到错题集文件: " + sourseZipPath)
	}

	// 清空并重新创建wrong目录
	if err := MachineToolClearExamTemp(wrongDir); err != nil {
		return createErrorPage("清理错题目录失败: " + err.Error())
	}

	// 解压错题集
	if err := MachineToolUnzip(sourseZipPath, wrongDir); err != nil {
		return createErrorPage("解压错题集失败: " + err.Error())
	}

	// 查找错题JSON文件
	wrongJSONPath, err := findWrongQuestionJSON(wrongDir)
	if err != nil {
		return createErrorPage("查找错题JSON文件失败: " + err.Error())
	}

	// 加载错题数据
	wrongSet, err := LoadWrongQuestionSet(wrongJSONPath)
	if err != nil {
		return createErrorPage("加载错题数据失败: " + err.Error())
	}

	// 创建错题查看界面
	return createWrongQuestionViewerContent(window, wrongSet, wrongDir, gradeFolderPath, returnCallback)
}

// createErrorPage 创建错误页面
func createErrorPage(errorMsg string) fyne.CanvasObject {
	errorLabel := widget.NewLabel(errorMsg)
	errorLabel.Alignment = fyne.TextAlignCenter

	return container.NewCenter(errorLabel)
}

// findWrongQuestionJSON 查找错题JSON文件
func findWrongQuestionJSON(wrongDir string) (string, error) {
	files, err := os.ReadDir(wrongDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") && strings.HasPrefix(file.Name(), "wrong_") {
			return filepath.Join(wrongDir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("未找到错题JSON文件")
}

// createWrongQuestionViewerContent 创建错题查看器界面内容
func createWrongQuestionViewerContent(window fyne.Window, wrongSet *WrongQuestionSet, wrongDir string, gradeFolderPath string, returnCallback func()) fyne.CanvasObject {
	// 获取当前分辨率配置
	resConfig := GetCurrentResolutionConfig()

	// 主容器
	mainContainer := container.NewVBox()

	// 1. 标题区域
	titleLabel := widget.NewLabel(fmt.Sprintf(GetLocalized("wrong_viewer_title"), wrongSet.SourceExam, wrongSet.WrongCount))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	mainContainer.Add(titleLabel)

	// 添加按钮区域
	buttonContainer := container.NewHBox()

	// 返回按钮
	backBtn := widget.NewButton(GetLocalized("wrong_viewer_back_btn"), func() {
		returnToGradeList(returnCallback)
	})
	buttonContainer.Add(backBtn)

	// 添加删除按钮
	deleteBtn := widget.NewButton(GetLocalized("wrong_viewer_delete_btn"), func() {
		deleteWrongQuestionSet(gradeFolderPath, window, returnCallback)
	})
	buttonContainer.Add(deleteBtn)

	mainContainer.Add(buttonContainer)

	// 2. 题目卡片区域
	for i, wrongQ := range wrongSet.WrongQuestions {
		card := createWrongQuestionCard(i+1, wrongQ, wrongDir, window, resConfig)
		mainContainer.Add(card)
	}

	// 创建滚动容器
	scrollContainer := container.NewScroll(mainContainer)
	scrollContainer.SetMinSize(fyne.NewSize(800, 600))

	// 设置ESC快捷键
	setupEscapeShortcut(window, returnCallback)

	return scrollContainer
}

// returnToGradeList 返回到成绩列表界面
func returnToGradeList(returnCallback func()) {
	// 调用回调函数返回到成绩列表
	if returnCallback != nil {
		returnCallback()
	}
}

// setupEscapeShortcut 设置ESC快捷键
func setupEscapeShortcut(window fyne.Window, returnCallback func()) {
	// 监听键盘事件
	window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyEscape {
			returnToGradeList(returnCallback)
		}
	})
}

// deleteWrongQuestionSet 删除错题集
func deleteWrongQuestionSet(gradeFolderPath string, window fyne.Window, returnCallback func()) {
	// 显示确认对话框
	confirmDialog := dialog.NewConfirm(
		GetLocalized("wrong_viewer_confirm_delete_title"),
		fmt.Sprintf(GetLocalized("wrong_viewer_confirm_delete_message"), filepath.Base(gradeFolderPath)),
		func(confirmed bool) {
			if confirmed {
				// 执行删除操作
				err := os.RemoveAll(gradeFolderPath)
				if err != nil {
					dialog.ShowError(fmt.Errorf(GetLocalized("wrong_viewer_delete_error_message"), err), window)
					return
				}

				// 显示成功消息
				dialog.ShowInformation(GetLocalized("wrong_viewer_delete_success_title"), GetLocalized("wrong_viewer_delete_success_message"), window)

				// 返回成绩列表
				returnToGradeList(returnCallback)
			}
		},
		window,
	)
	confirmDialog.Show()
}

// createWrongQuestionCard 创建错题卡片 - 修复版
func createWrongQuestionCard(questionNumber int, wrongQ WrongQuestion, wrongDir string, window fyne.Window, resConfig ResolutionConfig) fyne.CanvasObject {
	// 构建图片路径
	imgPaths := make([]string, len(wrongQ.Images))
	for j, img := range wrongQ.Images {
		imgPaths[j] = filepath.Join(wrongDir, "add", filepath.Base(img))
	}

	// 根据题目类型创建不同的显示卡片
	var content *fyne.Container
	switch wrongQ.Type {
	case "SC", "SCIMG":
		content = createSCWrongCardContent(wrongQ, imgPaths, window)
	case "MC", "MCIMG":
		content = createMCWrongCardContent(wrongQ, imgPaths, window)
	case "FL", "FLIMG":
		content = createFLWrongCardContent(wrongQ, imgPaths, window)
	case "DR":
		content = createDRWrongCardContent(wrongQ, imgPaths, window, resConfig, wrongDir)
	default:
		content = container.NewVBox(widget.NewLabel(GetLocalized("machine_unknown_type") + ": " + wrongQ.Type))
	}

	// 创建卡片
	card := widget.NewCard(
		fmt.Sprintf(GetLocalized("wrong_viewer_card_title"), questionNumber, getQuestionTypeName(wrongQ.Type)),
		"",
		content,
	)

	return card
}

// getQuestionTypeName 获取题型名称
func getQuestionTypeName(qType string) string {
	switch qType {
	case "SC", "SCIMG":
		return GetLocalized("machine_single_choice_subtitle")
	case "MC", "MCIMG":
		return GetLocalized("machine_multi_choice_subtitle")
	case "FL", "FLIMG":
		return GetLocalized("machine_fill_subtitle")
	case "DR":
		return GetLocalized("machine_material_title")
	default:
		return GetLocalized("machine_unknown_type")
	}
}

// createSCWrongCardContent 创建单选题错题卡片内容 - 修复版
func createSCWrongCardContent(wrongQ WrongQuestion, imgPaths []string, window fyne.Window) *fyne.Container {
	content := container.NewVBox()

	// 题目描述
	questionLabel := widget.NewLabel(wrongQ.Question)
	questionLabel.Wrapping = fyne.TextWrapWord
	content.Add(questionLabel)

	// 显示选项
	if len(wrongQ.Options) > 0 {
		optionsLabel := widget.NewLabel(GetLocalized("wrong_viewer_options_label"))
		optionsLabel.TextStyle = fyne.TextStyle{Bold: true}
		content.Add(optionsLabel)

		for i, option := range wrongQ.Options {
			optionText := fmt.Sprintf("%c. %s", 'A'+i, option)
			optionLabel := widget.NewLabel(optionText)
			content.Add(optionLabel)
		}
	}

	// 显示正确答案 - 修复答案显示逻辑
	correctAnswer := "未知"
	switch correctAns := wrongQ.CorrectAnswer.(type) {
	case []string:
		if len(correctAns) > 0 {
			correctAnswer = correctAns[0]
		}
	case []interface{}:
		if len(correctAns) > 0 {
			if ans, ok := correctAns[0].(string); ok {
				correctAnswer = ans
			}
		}
	case string:
		correctAnswer = correctAns
	}

	correctAnswerLabel := widget.NewLabel(fmt.Sprintf(GetLocalized("wrong_viewer_correct_answer_label"), correctAnswer))
	correctAnswerLabel.TextStyle = fyne.TextStyle{Bold: true}
	content.Add(correctAnswerLabel)

	// 显示图片（如果有）
	if len(imgPaths) > 0 {
		imgGrid := createWrongViewerImageGrid(imgPaths, window)
		content.Add(imgGrid)
	}

	return content
}

// createMCWrongCardContent 创建多选题错题卡片内容 - 修复版
func createMCWrongCardContent(wrongQ WrongQuestion, imgPaths []string, window fyne.Window) *fyne.Container {
	content := container.NewVBox()

	// 题目描述
	questionLabel := widget.NewLabel(wrongQ.Question)
	questionLabel.Wrapping = fyne.TextWrapWord
	content.Add(questionLabel)

	// 显示选项
	if len(wrongQ.Options) > 0 {
		optionsLabel := widget.NewLabel(GetLocalized("wrong_viewer_options_label"))
		optionsLabel.TextStyle = fyne.TextStyle{Bold: true}
		content.Add(optionsLabel)

		for i, option := range wrongQ.Options {
			optionText := fmt.Sprintf("%c. %s", 'A'+i, option)
			optionLabel := widget.NewLabel(optionText)
			content.Add(optionLabel)
		}
	}

	// 显示正确答案 - 修复答案显示逻辑
	correctAnswer := "未知"
	switch correctAns := wrongQ.CorrectAnswer.(type) {
	case []string:
		correctAnswer = strings.Join(correctAns, ", ")
	case []interface{}:
		var answers []string
		for _, ans := range correctAns {
			if str, ok := ans.(string); ok {
				answers = append(answers, str)
			}
		}
		correctAnswer = strings.Join(answers, ", ")
	}

	correctAnswerLabel := widget.NewLabel(fmt.Sprintf(GetLocalized("wrong_viewer_correct_answer_label"), correctAnswer))
	correctAnswerLabel.TextStyle = fyne.TextStyle{Bold: true}
	content.Add(correctAnswerLabel)

	// 显示图片（如果有）
	if len(imgPaths) > 0 {
		imgGrid := createWrongViewerImageGrid(imgPaths, window)
		content.Add(imgGrid)
	}

	return content
}

// createFLWrongCardContent 创建填空题错题卡片内容 - 修复版
func createFLWrongCardContent(wrongQ WrongQuestion, imgPaths []string, window fyne.Window) *fyne.Container {
	content := container.NewVBox()

	// 题目描述
	questionLabel := widget.NewLabel(wrongQ.Question)
	questionLabel.Wrapping = fyne.TextWrapWord
	content.Add(questionLabel)

	// 显示正确答案 - 修复填空题答案显示逻辑
	correctAnswers := make([]string, wrongQ.BlankCount)

	// 初始化正确答案
	for i := 0; i < wrongQ.BlankCount; i++ {
		correctAnswers[i] = "未知"
	}

	// 获取正确答案 - 修复逻辑
	switch correctAns := wrongQ.CorrectAnswer.(type) {
	case [][]string:
		for i, blankAnswers := range correctAns {
			if i < len(correctAnswers) {
				correctAnswers[i] = strings.Join(blankAnswers, " 或 ")
			}
		}
	case []interface{}:
		for i, blankAnswerInterface := range correctAns {
			if i < len(correctAnswers) {
				switch blankAnswer := blankAnswerInterface.(type) {
				case []interface{}:
					var options []string
					for _, opt := range blankAnswer {
						if str, ok := opt.(string); ok {
							options = append(options, str)
						}
					}
					correctAnswers[i] = strings.Join(options, " 或 ")
				case string:
					correctAnswers[i] = blankAnswer
				case []string:
					correctAnswers[i] = strings.Join(blankAnswer, " 或 ")
				}
			}
		}
	case []string:
		// 处理单个空的填空题
		if len(correctAns) > 0 && wrongQ.BlankCount == 1 {
			correctAnswers[0] = strings.Join(correctAns, " 或 ")
		}
	}

	// 显示每个空的正确答案
	for i := 0; i < wrongQ.BlankCount; i++ {
		blankLabel := widget.NewLabel(fmt.Sprintf(GetLocalized("wrong_viewer_blank_label"), i+1))
		blankLabel.TextStyle = fyne.TextStyle{Bold: true}
		content.Add(blankLabel)

		correctAnswerText := widget.NewLabel(fmt.Sprintf(GetLocalized("wrong_viewer_blank_correct_label"), correctAnswers[i]))
		content.Add(correctAnswerText)
	}

	// 显示图片（如果有）
	if len(imgPaths) > 0 {
		imgGrid := createWrongViewerImageGrid(imgPaths, window)
		content.Add(imgGrid)
	}

	return content
}

// createDRWrongCardContent 创建DR题错题卡片内容
func createDRWrongCardContent(wrongQ WrongQuestion, imgPaths []string, window fyne.Window, resConfig ResolutionConfig, wrongDir string) *fyne.Container {
	content := container.NewVBox()

	// DR题目描述
	questionLabel := widget.NewLabel(wrongQ.Question)
	questionLabel.Wrapping = fyne.TextWrapWord
	content.Add(questionLabel)

	// 显示材料（如果有）
	if len(wrongQ.Materials) > 0 {
		materialsLabel := widget.NewLabel(GetLocalized("wrong_viewer_materials_label"))
		materialsLabel.TextStyle = fyne.TextStyle{Bold: true}
		content.Add(materialsLabel)

		for _, material := range wrongQ.Materials {
			materialLabel := widget.NewLabel(material)
			materialLabel.Wrapping = fyne.TextWrapWord
			content.Add(materialLabel)
		}
	}

	// 显示DR图片（如果有）
	if len(imgPaths) > 0 {
		imgGrid := createWrongViewerImageGrid(imgPaths, window)
		content.Add(imgGrid)
	}

	// 显示子错题
	if len(wrongQ.SubWrongQuestions) > 0 {
		subTitleLabel := widget.NewLabel(GetLocalized("wrong_viewer_sub_wrong_label"))
		subTitleLabel.TextStyle = fyne.TextStyle{Bold: true}
		content.Add(subTitleLabel)

		for subIndex, subWrongQ := range wrongQ.SubWrongQuestions {
			// 为子问题构建图片路径
			subImgPaths := make([]string, len(subWrongQ.Images))
			for j, img := range subWrongQ.Images {
				subImgPaths[j] = filepath.Join(wrongDir, "add", filepath.Base(img))
			}

			// 根据子问题类型创建内容
			var subContent *fyne.Container
			switch subWrongQ.Type {
			case "SC", "SCIMG":
				subContent = createSCWrongCardContent(subWrongQ, subImgPaths, window)
			case "MC", "MCIMG":
				subContent = createMCWrongCardContent(subWrongQ, subImgPaths, window)
			case "FL", "FLIMG":
				subContent = createFLWrongCardContent(subWrongQ, subImgPaths, window)
			default:
				subContent = container.NewVBox(widget.NewLabel(GetLocalized("machine_unknown_type") + ": " + subWrongQ.Type))
			}

			// 添加子问题标题
			subCard := widget.NewCard(
				fmt.Sprintf(GetLocalized("wrong_viewer_sub_card_title"), subIndex+1, getQuestionTypeName(subWrongQ.Type)),
				"",
				subContent,
			)
			content.Add(subCard)
		}
	}

	return content
}

// createWrongViewerImageGrid 创建错题查看器的图片网格（避免重复声明）
func createWrongViewerImageGrid(imgPaths []string, window fyne.Window) fyne.CanvasObject {
	if len(imgPaths) == 0 {
		return container.NewVBox()
	}

	imgObjs := []fyne.CanvasObject{}
	for _, p := range imgPaths {
		img := LoadImage(p)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(200, 120))

		// 创建可点击的图片
		clickableImg := createWrongViewerClickableImage(img, p, window)
		imgObjs = append(imgObjs, clickableImg)
	}

	// 如果图片数量少于等于4个，直接在一行显示
	if len(imgObjs) <= 4 {
		return container.NewHBox(imgObjs...)
	}

	// 如果图片数量多于4个，分组显示，每行最多4个
	var rows []fyne.CanvasObject
	for i := 0; i < len(imgObjs); i += 4 {
		end := i + 4
		if end > len(imgObjs) {
			end = len(imgObjs)
		}
		row := container.NewHBox(imgObjs[i:end]...)
		rows = append(rows, row)
	}

	return container.NewVBox(rows...)
}

// createWrongViewerClickableImage 创建错题查看器的可点击图片容器（避免重复声明）
func createWrongViewerClickableImage(img *canvas.Image, imgPath string, window fyne.Window) fyne.CanvasObject {
	// 创建透明按钮
	btn := widget.NewButton("", nil)
	btn.Importance = widget.LowImportance

	// 设置按钮点击事件
	btn.OnTapped = func() {
		// 初始化 Lightbox 查看器（如果尚未初始化）
		lightboxViewer := NewLightboxViewer(window)

		// 显示图片查看器
		lightboxViewer.Show(imgPath)
	}

	// 设置按钮大小与图片相同
	imgSize := img.MinSize()
	btn.Resize(imgSize)

	// 创建堆叠容器：图片在下层，透明按钮在上层
	stack := container.NewStack(img, btn)
	stack.Resize(imgSize)

	return stack
}
