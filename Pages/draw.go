package pages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"TestFyne-1119/Pages/tools" // 添加这个导入

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// GradeData 成绩数据结构
type GradeData struct {
	ExamName     string  `json:"examName"`
	FinalScore   float64 `json:"finalScore"`
	TotalScore   float64 `json:"totalScore"`
	CorrectCount int     `json:"correctCount"`
	WrongCount   int     `json:"wrongCount"`
	StartTime    string  `json:"startTime"`
	EndTime      string  `json:"endTime"`
	Duration     string  `json:"duration"`
}

// 添加全局变量来保存主界面引用
var (
	mainAppContent         fyne.CanvasObject
	contentStack           *fyne.Container
	currentScrollContainer *container.Scroll
)

// DrawPage 考试详细界面
func DrawPage(window fyne.Window, appContent fyne.CanvasObject, stack *fyne.Container) fyne.CanvasObject {
	// 保存主界面引用
	mainAppContent = appContent
	contentStack = stack

	// 主滚动容器
	scrollContainer := container.NewScroll(container.NewVBox())
	scrollContainer.SetMinSize(fyne.NewSize(800, 600))
	currentScrollContainer = scrollContainer

	// 刷新按钮
	refreshBtn := widget.NewButton(tools.GetLocalized("draw_refresh_btn"), func() {
		refreshGradeCards(scrollContainer, window)
	})

	// 顶部工具栏
	toolbar := container.NewHBox(
		refreshBtn,
		layout.NewSpacer(),
		widget.NewLabel(tools.GetLocalized("draw_toolbar_title")),
	)

	// 主布局
	content := container.NewBorder(
		toolbar,
		nil, nil, nil,
		scrollContainer,
	)

	// 初始加载成绩卡片
	refreshGradeCards(scrollContainer, window)

	return content
}

// refreshGradeCards 刷新成绩卡片
func refreshGradeCards(scrollContainer *container.Scroll, window fyne.Window) {
	gradesDir := "data/grades"

	// 清空现有卡片
	if scrollContainer.Content != nil {
		scrollContainer.Content.(*fyne.Container).Objects = nil
	}

	// 创建新的卡片容器
	cardContainer := container.NewVBox()

	// 添加标题卡片
	titleCard := createTitleCard()
	cardContainer.Add(titleCard)

	// 读取grades目录
	folders, err := getGradeFolders(gradesDir)
	if err != nil {
		errorLabel := widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_error_read_grades"), err))
		errorCard := widget.NewCard(tools.GetLocalized("draw_error_title"), "", errorLabel)
		cardContainer.Add(errorCard)
		scrollContainer.Content = cardContainer
		scrollContainer.Refresh()
		return
	}

	if len(folders) == 0 {
		emptyLabel := widget.NewLabel(tools.GetLocalized("draw_empty_label"))
		emptyCard := widget.NewCard(tools.GetLocalized("draw_empty_title"), "", emptyLabel)
		cardContainer.Add(emptyCard)
		scrollContainer.Content = cardContainer
		scrollContainer.Refresh()
		return
	}

	// 为每个成绩文件夹创建卡片
	for _, folder := range folders {
		gradeCard := createGradeCard(folder, window, scrollContainer)
		if gradeCard != nil {
			cardContainer.Add(gradeCard)
		}
	}

	scrollContainer.Content = cardContainer
	scrollContainer.Refresh()
}

// getGradeFolders 获取成绩文件夹列表
func getGradeFolders(gradesDir string) ([]string, error) {
	var folders []string

	// 检查目录是否存在
	if _, err := os.Stat(gradesDir); os.IsNotExist(err) {
		return folders, nil
	}

	entries, err := os.ReadDir(gradesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, filepath.Join(gradesDir, entry.Name()))
		}
	}

	return folders, nil
}

// createTitleCard 创建标题卡片
func createTitleCard() fyne.CanvasObject {
	titleLabel := widget.NewLabel(tools.GetLocalized("draw_title_label"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter

	return widget.NewCard("", "", container.NewCenter(titleLabel))
}

// createGradeCard 创建成绩卡片
func createGradeCard(folderPath string, window fyne.Window, scrollContainer *container.Scroll) fyne.CanvasObject {
	// 查找JSON文件
	jsonFile, err := findGradeJSONFile(folderPath)
	if err != nil {
		fmt.Printf("在文件夹 %s 中找不到成绩JSON文件: %v\n", folderPath, err)
		return nil
	}

	// 读取JSON数据
	gradeData, err := loadGradeData(jsonFile)
	if err != nil {
		fmt.Printf("读取成绩数据失败 %s: %v\n", jsonFile, err)
		return nil
	}

	// 计算是否及格 (60分及格)
	isPassed := gradeData.FinalScore >= 60
	passStatus := tools.GetLocalized("draw_status_fail")
	if isPassed {
		passStatus = tools.GetLocalized("draw_status_pass")
	}

	// 计算正确率
	accuracy := 0.0
	totalQuestions := gradeData.CorrectCount + gradeData.WrongCount
	if totalQuestions > 0 {
		accuracy = float64(gradeData.CorrectCount) / float64(totalQuestions) * 100
	}

	// 解析考试时间，提取更友好的显示格式
	displayTime := formatDisplayTime(gradeData.EndTime)

	// 创建卡片内容
	content := container.NewVBox(
		// 第一行：考试名称和时间
		container.NewHBox(
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_exam_name_label"), gradeData.ExamName)),
			layout.NewSpacer(),
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_exam_time_label"), displayTime)),
		),

		// 第二行：分数信息
		container.NewHBox(
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_score_label"), gradeData.FinalScore, gradeData.TotalScore)),
			layout.NewSpacer(),
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_status_label"), passStatus)),
		),

		// 第三行：题目统计
		container.NewHBox(
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_total_questions_label"), totalQuestions)),
			layout.NewSpacer(),
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_correct_label"), gradeData.CorrectCount)),
			layout.NewSpacer(),
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_wrong_label"), gradeData.WrongCount)),
		),

		// 第四行：正确率和考试时长
		container.NewHBox(
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_accuracy_label"), accuracy)),
			layout.NewSpacer(),
			widget.NewLabel(fmt.Sprintf(tools.GetLocalized("draw_duration_label"), gradeData.Duration)),
		),

		// 分隔线
		widget.NewSeparator(),

		// 按钮区域
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton(tools.GetLocalized("draw_wrong_btn"), func() {
				// 进入错题查看界面
				wrongViewer := tools.CreateWrongQuestionViewer(window, folderPath, func() {
					window.SetContent(mainAppContent)
					contentStack.Objects = []fyne.CanvasObject{DrawPage(window, mainAppContent, contentStack)}
					contentStack.Refresh()
				})
				window.SetContent(wrongViewer)
			}),
			widget.NewButton(tools.GetLocalized("draw_detail_btn"), func() {
				showExamDetails(gradeData, window)
			}),
			widget.NewButton(tools.GetLocalized("draw_delete_btn"), func() {
				deleteGradeFolder(folderPath, window, func() {
					refreshGradeCards(scrollContainer, window)
				})
			}),
		),
	)

	// 根据是否及格设置卡片标题
	cardTitle := fmt.Sprintf(tools.GetLocalized("draw_card_title"), passStatus)
	return widget.NewCard(cardTitle, "", content)
}

// findGradeJSONFile 在文件夹中查找成绩JSON文件
func findGradeJSONFile(folderPath string) (string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			return filepath.Join(folderPath, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("未找到JSON文件")
}

// loadGradeData 加载成绩数据
func loadGradeData(jsonFile string) (*GradeData, error) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	var gradeData GradeData
	if err := json.Unmarshal(data, &gradeData); err != nil {
		return nil, err
	}

	return &gradeData, nil
}

// formatDisplayTime 格式化显示时间
func formatDisplayTime(timeStr string) string {
	// 尝试解析时间
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		// 如果解析失败，返回原始字符串
		return timeStr
	}

	// 格式化为更友好的显示
	return t.Format("2006年01月02日 15:04")
}

// showExamDetails 显示考试详情
func showExamDetails(gradeData *GradeData, window fyne.Window) {
	totalQuestions := gradeData.CorrectCount + gradeData.WrongCount
	accuracy := 0.0
	if totalQuestions > 0 {
		accuracy = float64(gradeData.CorrectCount) / float64(totalQuestions) * 100
	}

	details := fmt.Sprintf(
		tools.GetLocalized("draw_detail_content"),
		gradeData.ExamName,
		gradeData.StartTime,
		gradeData.EndTime,
		gradeData.Duration,
		gradeData.FinalScore,
		gradeData.TotalScore,
		gradeData.CorrectCount,
		gradeData.WrongCount,
		totalQuestions,
		accuracy,
		getPassStatus(gradeData.FinalScore),
	)

	dialog.ShowCustom(tools.GetLocalized("draw_detail_title"), tools.GetLocalized("draw_close_btn"), widget.NewLabel(details), window)
}

// getPassStatus 获取考试状态
func getPassStatus(score float64) string {
	if score >= 60 {
		return tools.GetLocalized("draw_status_pass")
	}
	return tools.GetLocalized("draw_status_fail")
}

// deleteGradeFolder 删除成绩文件夹
func deleteGradeFolder(folderPath string, window fyne.Window, refreshFunc func()) {
	// 显示确认对话框
	confirmDialog := dialog.NewConfirm(
		tools.GetLocalized("draw_delete_confirm_title"),
		fmt.Sprintf(tools.GetLocalized("draw_delete_confirm_content"), filepath.Base(folderPath)),
		func(confirmed bool) {
			if confirmed {
				// 执行删除操作
				err := os.RemoveAll(folderPath)
				if err != nil {
					dialog.ShowError(fmt.Errorf(tools.GetLocalized("draw_delete_failed"), err), window)
					return
				}

				// 显示成功消息
				dialog.ShowInformation(tools.GetLocalized("draw_delete_success_title"), tools.GetLocalized("draw_delete_success_content"), window)

				// 刷新成绩列表
				refreshFunc()
			}
		},
		window,
	)
	confirmDialog.Show()
}
