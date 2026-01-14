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
	"image"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/chai2010/webp"
)

// LoadImage 支持加载多种格式的图片，包括 WebP
func LoadImage(path string) *canvas.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("无法打开图片文件: %s, 错误: %v\n", path, err)
		return canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	}
	defer file.Close()

	var img image.Image
	if strings.HasSuffix(strings.ToLower(path), ".webp") {
		img, err = webp.Decode(file)
		if err != nil {
			fmt.Printf("解码WebP图片失败: %s, 错误: %v\n", path, err)
		}
	} else {
		img, _, err = image.Decode(file)
		if err != nil {
			fmt.Printf("解码图片失败: %s, 错误: %v\n", path, err)
		}
	}
	if err != nil {
		return canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	}

	canvasImg := canvas.NewImageFromImage(img)
	return canvasImg
}

// 创建支持换行的标签（重命名以避免冲突）
func createWrappedLabelForImage(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

// 创建支持换行的单选选项 - 使用自定义容器（修复单选问题）
func createWrappedRadioGroupForImage(options []string, onSelected func(string)) fyne.CanvasObject {
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
func createWrappedCheckboxesForImage(options []string, onChecked func(int, bool)) fyne.CanvasObject {
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

// createClickableImageGrid 创建可点击的图片网格（带注释）- 横向排列，用于附带图片
func createClickableImageGrid(imgPaths []string, window fyne.Window) fyne.CanvasObject {
	if len(imgPaths) == 0 {
		return container.NewVBox()
	}

	// 创建图片和注释的容器
	var rows []fyne.CanvasObject
	imgObjs := []fyne.CanvasObject{}

	// 为每张图片创建带注释的容器
	for i, p := range imgPaths {
		img := LoadImage(p)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(200, 120))

		// 创建可点击的图片
		clickableImg := createClickableImageContainer(img, p, window)

		// 创建图片注释标签
		annotation := widget.NewLabel(fmt.Sprintf("图%d", i+1))
		annotation.Alignment = fyne.TextAlignCenter

		// 创建单个图片的垂直容器：图片 + 注释
		imgWithAnnotation := container.NewVBox(
			clickableImg,
			annotation,
		)

		imgObjs = append(imgObjs, imgWithAnnotation)
	}

	// 如果图片数量少于等于4个，直接在一行显示
	if len(imgObjs) <= 4 {
		return container.NewHBox(imgObjs...)
	}

	// 如果图片数量多于4个，分组显示，每行最多4个
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

// createClickableImageContainer 创建可点击的图片容器
func createClickableImageContainer(img *canvas.Image, imgPath string, window fyne.Window) fyne.CanvasObject {
	// 创建透明按钮
	btn := widget.NewButton("", nil)
	btn.Importance = widget.LowImportance

	// 设置按钮点击事件
	btn.OnTapped = func() {
		// 初始化 Lightbox 查看器（如果尚未初始化）
		var lightboxViewer *LightboxViewer
		lightboxViewer = NewLightboxViewer(window)

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

// createImageGrid 创建图片网格，每行最多4个图片（保留原有函数用于兼容）
func createImageGrid(imgPaths []string) fyne.CanvasObject {
	if len(imgPaths) == 0 {
		return container.NewVBox()
	}

	imgObjs := []fyne.CanvasObject{}
	for _, p := range imgPaths {
		img := LoadImage(p)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(200, 120))
		imgObjs = append(imgObjs, img)
	}

	if len(imgObjs) <= 4 {
		return container.NewHBox(imgObjs...)
	}

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

// calculateImageSize 计算图片显示尺寸
func calculateImageSize(img *canvas.Image, cardWidth float32) fyne.Size {
	// 获取图片原始尺寸
	origSize := img.MinSize()

	// 最小尺寸512x512
	minSize := fyne.NewSize(512, 512)

	// 计算缩放比例
	scale := float32(1.0)
	if origSize.Width > cardWidth-40 { // 留出边距
		scale = (cardWidth - 40) / origSize.Width
	}

	// 应用缩放
	width := origSize.Width * scale
	height := origSize.Height * scale

	// 确保不小于最小尺寸
	if width < minSize.Width {
		width = minSize.Width
		// 等比例调整高度
		if origSize.Width > 0 {
			height = origSize.Height * (width / origSize.Width)
		}
	}
	if height < minSize.Height {
		height = minSize.Height
		// 等比例调整宽度
		if origSize.Height > 0 {
			width = origSize.Width * (height / origSize.Height)
		}
	}

	return fyne.NewSize(width, height)
}

// createQuestionImageLayout 创建题干图片布局（竖向排列，带注释）- 用于题干带图片
func createQuestionImageLayout(imgPaths []string, cardWidth float32, window fyne.Window) fyne.CanvasObject {
	if len(imgPaths) == 0 {
		return container.NewVBox()
	}

	var imgContainers []fyne.CanvasObject

	// 为每张图片创建带注释的容器
	for i, path := range imgPaths {
		img := LoadImage(path)
		img.FillMode = canvas.ImageFillContain

		// 计算图片尺寸
		imgSize := calculateImageSize(img, cardWidth)
		img.SetMinSize(imgSize)

		// 创建可点击的图片容器
		clickableImg := createClickableImageContainer(img, path, window)

		// 创建图片注释标签
		annotation := widget.NewLabel(fmt.Sprintf("图%d", i+1))
		annotation.Alignment = fyne.TextAlignCenter

		// 创建单个图片的垂直容器：图片 + 注释
		imgWithAnnotation := container.NewVBox(
			clickableImg,
			annotation,
		)

		imgContainers = append(imgContainers, imgWithAnnotation)
	}

	// 题干图片竖向排列，每张图片单独一行
	return container.NewVBox(imgContainers...)
}

// 单选题卡片（有图片）- 附带图片，横向排列
func NewCardSingleChoiceWithImg(questionNumber int, question string, options []string, imgPaths []string, onAnswered func(int), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {
	radioContainer := createWrappedRadioGroupForImage(options, func(selected string) {
		for i, opt := range options {
			if opt == selected {
				onAnswered(i)
				break
			}
		}
	})

	// 使用容器包装，确保换行
	content := container.NewVBox(
		createWrappedLabelForImage(question), // 题目单独显示
		radioContainer,
	)

	// 创建可点击的图片网格（带注释）- 横向排列
	if len(imgPaths) > 0 {
		imgGrid := createClickableImageGrid(imgPaths, window)
		content.Add(imgGrid)
	}

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_single_img_subtitle"),
		content,
	)
}

// 多选题卡片（有图片）- 附带图片，横向排列
func NewCardMultiChoiceWithImg(questionNumber int, question string, options []string, imgPaths []string, onAnswered func([]int), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {
	selected := make([]bool, len(options))

	checkContainer := createWrappedCheckboxesForImage(options, func(idx int, checked bool) {
		selected[idx] = checked
		var ans []int
		for j, v := range selected {
			if v {
				ans = append(ans, j)
			}
		}
		onAnswered(ans)
	})

	// 使用容器包装，确保换行
	content := container.NewVBox(
		createWrappedLabelForImage(question), // 题目单独显示
		checkContainer,
	)

	// 创建可点击的图片网格（带注释）- 横向排列
	if len(imgPaths) > 0 {
		imgGrid := createClickableImageGrid(imgPaths, window)
		content.Add(imgGrid)
	}

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_multi_img_subtitle"),
		content,
	)
}

// 填空题卡片（有图片）- 附带图片，横向排列
func NewCardFillBlankWithImg(questionNumber int, question string, blanks int, imgPaths []string, onAnswered func([]string), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {
	entries := make([]*widget.Entry, blanks)
	box := container.NewVBox()

	// 先添加题目
	box.Add(createWrappedLabelForImage(question))

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

	// 创建可点击的图片网格（带注释）- 横向排列
	if len(imgPaths) > 0 {
		imgGrid := createClickableImageGrid(imgPaths, window)
		box.Add(imgGrid)
	}

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_fill_img_subtitle"),
		box,
	)
}

// ===== 新题型：题干为图片的卡片 =====

// NewCardSCIMG 题干是图单选 - 题干图片竖向排列
func NewCardSCIMG(questionNumber int, question string, options []string, imgPaths []string, onAnswered func(int), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {
	radioContainer := createWrappedRadioGroupForImage(options, func(selected string) {
		for i, opt := range options {
			if opt == selected {
				onAnswered(i)
				break
			}
		}
	})

	// 创建题干图片布局（竖向排列，带注释）
	var imgLayout fyne.CanvasObject
	if len(imgPaths) > 0 {
		imgLayout = createQuestionImageLayout(imgPaths, cardWidth, window) // 使用动态宽度
	} else {
		imgLayout = container.NewVBox()
	}

	content := container.NewVBox(
		createWrappedLabelForImage(question), // 文字题目描述
	)

	if len(imgPaths) > 0 {
		content.Add(imgLayout) // 题干图片（竖向排列，带注释）
	}
	content.Add(radioContainer) // 选项

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_single_img_subtitle"),
		content,
	)
}

// NewCardMCIMG 题干是图多选 - 题干图片竖向排列
func NewCardMCIMG(questionNumber int, question string, options []string, imgPaths []string, onAnswered func([]int), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {
	selected := make([]bool, len(options))

	checkContainer := createWrappedCheckboxesForImage(options, func(idx int, checked bool) {
		selected[idx] = checked
		var ans []int
		for j, v := range selected {
			if v {
				ans = append(ans, j)
			}
		}
		onAnswered(ans)
	})

	// 创建题干图片布局（竖向排列，带注释）
	var imgLayout fyne.CanvasObject
	if len(imgPaths) > 0 {
		imgLayout = createQuestionImageLayout(imgPaths, cardWidth, window) // 使用动态宽度
	} else {
		imgLayout = container.NewVBox()
	}

	content := container.NewVBox(
		createWrappedLabelForImage(question), // 文字题目描述
	)

	if len(imgPaths) > 0 {
		content.Add(imgLayout) // 题干图片（竖向排列，带注释）
	}
	content.Add(checkContainer) // 选项

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_multi_img_subtitle"),
		content,
	)
}

// NewCardFLIMG 题干是图填空 - 题干图片竖向排列
func NewCardFLIMG(questionNumber int, question string, blanks int, imgPaths []string, onAnswered func([]string), window fyne.Window, getLocalized func(string) string, cardWidth float32) fyne.CanvasObject {
	entries := make([]*widget.Entry, blanks)
	box := container.NewVBox()

	// 创建题干图片布局（竖向排列，带注释）
	var imgLayout fyne.CanvasObject
	if len(imgPaths) > 0 {
		imgLayout = createQuestionImageLayout(imgPaths, cardWidth, window) // 使用动态宽度
	} else {
		imgLayout = container.NewVBox()
	}

	// 添加组件顺序：题目描述 → 题干图片 → 填空输入框
	box.Add(createWrappedLabelForImage(question)) // 文字题目描述

	if len(imgPaths) > 0 {
		box.Add(imgLayout) // 题干图片（竖向排列，带注释）
	}

	for i := 0; i < blanks; i++ {
		entries[i] = widget.NewEntry()
		entries[i].SetPlaceHolder(fmt.Sprintf(getLocalized("machine_fill_placeholder"), i+1))
		box.Add(entries[i])

		// 实时监听输入变化
		entryIndex := i // 创建局部变量避免闭包问题
		entries[i].OnChanged = func(text string) {
			ans := make([]string, blanks)
			for j, e := range entries {
				ans[j] = e.Text
			}
			onAnswered(ans)
		}
		_ = entryIndex // 避免未使用变量警告
	}

	return widget.NewCard(
		fmt.Sprintf(getLocalized("machine_card_title"), questionNumber),
		getLocalized("machine_fill_img_subtitle"),
		box,
	)
}
