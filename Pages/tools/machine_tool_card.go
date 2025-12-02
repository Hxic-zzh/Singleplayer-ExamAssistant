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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 创建支持换行的标签
func createWrappedLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

// 创建支持换行的单选选项 - 使用自定义容器（修复单选问题）
func createWrappedRadioGroup(options []string, onSelected func(string)) fyne.CanvasObject {
	var checkItems []fyne.CanvasObject
	var checks []*widget.Check
	selectedIndex := -1

	for i, option := range options {
		idx := i
		check := widget.NewCheck(option, func(checked bool) {
			if checked {
				// 取消其他选中
				for j, chk := range checks {
					if j != idx {
						chk.SetChecked(false)
					}
				}
				selectedIndex = idx
				onSelected(option)
			} else {
				if selectedIndex == idx {
					selectedIndex = -1
					onSelected("")
				}
			}
		})
		checks = append(checks, check)
		// 为每个选项创建单独的容器，确保换行
		itemContainer := container.NewVBox(
			check,
		)
		checkItems = append(checkItems, itemContainer)
	}

	return container.NewVBox(checkItems...)
}

// 创建支持换行的多选选项 - 使用自定义容器
func createWrappedCheckboxes(options []string, onChecked func(int, bool)) fyne.CanvasObject {
	var checkItems []fyne.CanvasObject

	for i, option := range options {
		idx := i
		check := widget.NewCheck(option, func(checked bool) {
			onChecked(idx, checked)
		})
		// 为每个选项创建单独的容器
		itemContainer := container.NewVBox(
			check,
		)
		checkItems = append(checkItems, itemContainer)
	}

	return container.NewVBox(checkItems...)
}

// 单选题卡片（无图片）
func NewCardSingleChoice(questionNumber int, question string, options []string, onAnswered func(int), getLocalized func(string) string) fyne.CanvasObject {
	radioContainer := createWrappedRadioGroup(options, func(selected string) {
		for i, opt := range options {
			if opt == selected {
				onAnswered(i)
				break
			}
		}
	})

	content := container.NewVBox(
		createWrappedLabel(question),
		radioContainer,
	)

	card := widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_single_choice_subtitle"),
		content,
	)

	return card
}

// 多选题卡片（无图片）
func NewCardMultiChoice(questionNumber int, question string, options []string, onAnswered func([]int), getLocalized func(string) string) fyne.CanvasObject {
	selected := make([]bool, len(options))

	checkContainer := createWrappedCheckboxes(options, func(idx int, checked bool) {
		selected[idx] = checked
		var ans []int
		for j, v := range selected {
			if v {
				ans = append(ans, j)
			}
		}
		onAnswered(ans)
	})

	content := container.NewVBox(
		createWrappedLabel(question),
		checkContainer,
	)

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_multi_choice_subtitle"),
		content,
	)
}

// 填空题卡片（无图片）
func NewCardFillBlank(questionNumber int, question string, blanks int, onAnswered func([]string), getLocalized func(string) string) fyne.CanvasObject {
	entries := make([]*widget.Entry, blanks)
	box := container.NewVBox()

	box.Add(createWrappedLabel(question))

	for i := 0; i < blanks; i++ {
		entries[i] = widget.NewEntry()
		entries[i].SetPlaceHolder(fmt.Sprintf(getLocalized("machine_fill_placeholder"), i+1))
		box.Add(entries[i])
	}
	btn := widget.NewButton(getLocalized("machine_submit_btn"), func() {
		ans := make([]string, blanks)
		for i, e := range entries {
			ans[i] = e.Text
		}
		onAnswered(ans)
	})
	box.Add(btn)

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_fill_subtitle"),
		box,
	)
}
