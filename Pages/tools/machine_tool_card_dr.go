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
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// DRMaterialCard 材料阅读题材料卡片
func DRMaterialCard(materials []string, materialImages []string, window fyne.Window, getLocalized func(string) string) fyne.CanvasObject {
	content := container.NewVBox()

	// 添加材料标题
	titleLabel := widget.NewLabel(getLocalized("machine_material_title"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	content.Add(titleLabel)

	// 添加材料文本
	for i, material := range materials {
		if material != "" {
			materialLabel := widget.NewLabel(fmt.Sprintf(getLocalized("machine_material_item"), i+1, material))
			materialLabel.Wrapping = fyne.TextWrapWord
			content.Add(materialLabel)
		}
	}

	// 添加材料图片
	if len(materialImages) > 0 {
		imgGrid := createClickableImageGrid(materialImages, window)
		content.Add(imgGrid)
	}

	return widget.NewCard("", "", content)
}

// DRQuestionCard 材料阅读题的子题卡片
func DRQuestionCard(questionNumber int, question ExamQuestion, onAnswered func(string, []string), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {

	// 构建图片路径
	imgPaths := make([]string, len(question.Images))
	copy(imgPaths, question.Images)

	switch question.Type {
	case "SC":
		if len(question.Options) > 0 {
			if len(question.Images) > 0 {
				return NewCardSingleChoiceWithImg(questionNumber, question.Question, question.Options, imgPaths, func(selected int) {
					if selected == -1 {
						onAnswered(question.ID, []string{})
					} else {
						onAnswered(question.ID, []string{question.Options[selected]})
					}
				}, window, getLocalized, cardWidth)
			} else {
				return NewCardSingleChoice(questionNumber, question.Question, question.Options, func(selected int) {
					if selected == -1 {
						onAnswered(question.ID, []string{})
					} else {
						onAnswered(question.ID, []string{question.Options[selected]})
					}
				}, getLocalized)
			}
		}
	case "SCIMG":
		if len(question.Options) > 0 {
			return NewCardSCIMG(questionNumber, question.Question, question.Options, imgPaths, func(selected int) {
				if selected == -1 {
					onAnswered(question.ID, []string{})
				} else {
					onAnswered(question.ID, []string{question.Options[selected]})
				}
			}, window, getLocalized, cardWidth)
		}
	case "MC":
		if len(question.Options) > 0 {
			if len(question.Images) > 0 {
				return NewCardMultiChoiceWithImg(questionNumber, question.Question, question.Options, imgPaths, func(selected []int) {
					var ans []string
					for _, idx := range selected {
						ans = append(ans, question.Options[idx])
					}
					onAnswered(question.ID, ans)
				}, window, getLocalized, cardWidth)
			} else {
				return NewCardMultiChoice(questionNumber, question.Question, question.Options, func(selected []int) {
					var ans []string
					for _, idx := range selected {
						ans = append(ans, question.Options[idx])
					}
					onAnswered(question.ID, ans)
				}, getLocalized)
			}
		}
	case "MCIMG":
		if len(question.Options) > 0 {
			return NewCardMCIMG(questionNumber, question.Question, question.Options, imgPaths, func(selected []int) {
				var ans []string
				for _, idx := range selected {
					ans = append(ans, question.Options[idx])
				}
				onAnswered(question.ID, ans)
			}, window, getLocalized, cardWidth)
		}
	case "FL":
		// 填空题处理
		entries := make([]*widget.Entry, question.BlankCount)
		box := container.NewVBox()

		// 设置填空题标题样式
		titleLabel := widget.NewLabel(fmt.Sprintf(getLocalized("machine_fill_title"), questionNumber))
		titleLabel.TextStyle = fyne.TextStyle{Bold: true}
		box.Add(titleLabel)

		// 使用支持换行的题目标签
		questionLabel := widget.NewLabel(question.Question)
		questionLabel.Wrapping = fyne.TextWrapWord
		box.Add(questionLabel)

		for i := 0; i < question.BlankCount; i++ {
			entries[i] = widget.NewEntry()
			entries[i].SetPlaceHolder(fmt.Sprintf(getLocalized("machine_fill_placeholder"), i+1))
			box.Add(entries[i])

			// 实时监听输入变化
			entryIndex := i
			entries[i].OnChanged = func(text string) {
				ans := make([]string, question.BlankCount)
				for j, e := range entries {
					ans[j] = e.Text
				}
				onAnswered(question.ID, ans)
			}
			_ = entryIndex
		}

		return container.NewBorder(nil, nil, nil, nil, box)
	case "FLIMG":
		return NewCardFLIMG(questionNumber, question.Question, question.BlankCount, imgPaths, func(answers []string) {
			onAnswered(question.ID, answers)
		}, window, getLocalized, cardWidth)
	default:
		return widget.NewLabel(fmt.Sprintf(getLocalized("machine_unknown_child_type"), question.Type))
	}

	return widget.NewLabel(getLocalized("machine_question_data_error"))
}

// DRGroupCard 完整的材料阅读题组卡片 - 修复版：直接保存小题答案
func DRGroupCard(drQuestion ExamQuestion, childQuestions []ExamQuestion, onChildAnswered func(string, []string), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {
	// 按照题型顺序排序子题：SC, SCIMG, MC, MCIMG, FL, FLIMG
	typeOrder := map[string]int{
		"SC":    1,
		"SCIMG": 2,
		"MC":    3,
		"MCIMG": 4,
		"FL":    5,
		"FLIMG": 6,
	}
	sort.Slice(childQuestions, func(i, j int) bool {
		return typeOrder[childQuestions[i].Type] < typeOrder[childQuestions[j].Type]
	})

	// 创建主容器
	mainContainer := container.NewVBox()

	// 添加材料卡片
	materialCard := DRMaterialCard(drQuestion.Materials, drQuestion.Images, window, getLocalized)
	mainContainer.Add(materialCard)

	// 添加子题卡片 - 按照排序后的顺序
	for i, childQuestion := range childQuestions {
		childCard := DRQuestionCard(i+1, childQuestion, func(questionID string, answers []string) {
			// ✅ 修复：直接保存小题答案，而不是通过DR题ID
			onChildAnswered(questionID, answers)
		}, window, getLocalized, cardWidth)
		mainContainer.Add(childCard)
	}

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_material_card_title"), drQuestion.ID),
		fmt.Sprintf(getLocalized("machine_material_card_subtitle"), len(childQuestions)),
		mainContainer,
	)
}
