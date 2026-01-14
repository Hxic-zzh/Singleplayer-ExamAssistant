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
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ===================== 自定义特效控件 =====================
// EffectImage 特效图片自定义控件
type EffectImage struct {
	widget.BaseWidget
	image     *canvas.Image
	resource  fyne.Resource
	isVisible bool
	position  fyne.Position
	size      fyne.Size
	dirty     bool // 脏标记，避免无效刷新
}

// NewEffectImage 创建特效图片控件
func NewEffectImage() *EffectImage {
	e := &EffectImage{
		image:     canvas.NewImageFromResource(nil),
		isVisible: false,
		dirty:     false,
	}
	e.ExtendBaseWidget(e)
	e.image.FillMode = canvas.ImageFillContain
	e.image.Hide()
	return e
}

// SetResource 设置图片资源
func (e *EffectImage) SetResource(res fyne.Resource) {
	if e.resource != res {
		e.resource = res
		e.dirty = true
	}
}

// SetPosition 设置位置
func (e *EffectImage) SetPosition(pos fyne.Position) {
	if e.position != pos {
		e.position = pos
		e.dirty = true
	}
}

// SetSize 设置尺寸
func (e *EffectImage) SetSize(size fyne.Size) {
	if e.size != size {
		e.size = size
		e.dirty = true
	}
}

// Show 显示特效
func (e *EffectImage) Show() {
	if !e.isVisible {
		e.isVisible = true
		e.dirty = true
	}
}

// Hide 隐藏特效
func (e *EffectImage) Hide() {
	if e.isVisible {
		e.isVisible = false
		e.dirty = true
	}
}

// IsVisible 是否可见
func (e *EffectImage) IsVisible() bool {
	return e.isVisible
}

// CreateRenderer 创建渲染器
func (e *EffectImage) CreateRenderer() fyne.WidgetRenderer {
	return &effectImageRenderer{effectImage: e}
}

// effectImageRenderer 特效图片渲染器
type effectImageRenderer struct {
	effectImage *EffectImage
}

func (r *effectImageRenderer) Layout(size fyne.Size) {
	// 这里直接设置image的位置和大小
	if r.effectImage.isVisible {
		r.effectImage.image.Resize(r.effectImage.size)
		r.effectImage.image.Move(r.effectImage.position)
	}
}

func (r *effectImageRenderer) MinSize() fyne.Size {
	if r.effectImage.isVisible {
		return r.effectImage.size
	}
	return fyne.NewSize(0, 0)
}

func (r *effectImageRenderer) Refresh() {
	// Refresh是由Fyne在主线程调用的，所以这里直接操作UI是安全的
	if r.effectImage.dirty {
		if r.effectImage.isVisible && r.effectImage.resource != nil {
			r.effectImage.image.Resource = r.effectImage.resource
			r.effectImage.image.Resize(r.effectImage.size)
			r.effectImage.image.Move(r.effectImage.position)
			r.effectImage.image.Show()
			// ⚠️ 关键修复：这里不需要canvas.Refresh，只需要image.Refresh
			r.effectImage.image.Refresh()
		} else {
			r.effectImage.image.Hide()
		}
		r.effectImage.dirty = false
	}
}

func (r *effectImageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.effectImage.image}
}

func (r *effectImageRenderer) Destroy() {}

// ===================== 自定义图片控件 =====================
// TopLeftContainImage 左上角对齐的等比例缩放图片
type TopLeftContainImage struct {
	widget.BaseWidget
	resource fyne.Resource
	original fyne.Size
	image    *canvas.Image
}

// NewTopLeftContainImage 创建左上角对齐的等比例缩放图片
func NewTopLeftContainImage(res fyne.Resource, originalSize fyne.Size) *TopLeftContainImage {
	img := &TopLeftContainImage{
		resource: res,
		original: originalSize,
		image:    canvas.NewImageFromResource(res),
	}
	img.ExtendBaseWidget(img)
	// 关键：使用原始填充，通过Layout控制尺寸和位置
	img.image.FillMode = canvas.ImageFillOriginal
	return img
}

// CreateRenderer 创建渲染器
func (t *TopLeftContainImage) CreateRenderer() fyne.WidgetRenderer {
	return &topLeftImageRenderer{img: t}
}

// topLeftImageRenderer 自定义图片渲染器
type topLeftImageRenderer struct {
	img *TopLeftContainImage
}

// Layout 布局：实现左上角对齐的等比例缩放
func (r *topLeftImageRenderer) Layout(size fyne.Size) {
	if r.img.original.Width <= 0 || size.Width <= 0 {
		return
	}

	// 基于容器宽度等比例缩放
	scale := size.Width / r.img.original.Width
	scaledWidth := size.Width
	scaledHeight := r.img.original.Height * scale

	// 关键：始终左上角对齐 (0, 0)
	r.img.image.Move(fyne.NewPos(0, 0))
	r.img.image.Resize(fyne.NewSize(scaledWidth, scaledHeight))
}

// MinSize 最小尺寸
func (r *topLeftImageRenderer) MinSize() fyne.Size {
	// 返回原始尺寸的1/4，避免太大
	if r.img.original.Width > 0 {
		return fyne.NewSize(
			r.img.original.Width/4,
			r.img.original.Height/4,
		)
	}
	return fyne.NewSize(50, 50)
}

// Refresh 刷新
func (r *topLeftImageRenderer) Refresh() {
	r.img.image.Refresh()
}

// Objects 返回画布对象
func (r *topLeftImageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.img.image}
}

// Destroy 销毁
func (r *topLeftImageRenderer) Destroy() {
	// 清理资源
}

// ===================== 列表相关结构体 =====================
// ListItem 表示列表项的数据结构
type ListItem struct {
	Title string
	Icon  fyne.Resource
}

// CustomListConfig 列表配置
type CustomListConfig struct {
	Items           []ListItem   // 列表项数据
	OnSelected      func(id int) // 选中事件回调
	InitialSelected int          // 初始选中项ID，-1表示无选中
}

// ImageListConfig 带图片背景的列表配置
type ImageListConfig struct {
	ListConfig CustomListConfig // 列表配置
	Background *ListBackground  // 背景配置（可选）
}

// ListBackground 背景配置
type ListBackground struct {
	Image     fyne.Resource // 背景图片
	ImageSize fyne.Size     // 图片原始尺寸
	Color     color.Color   // 背景颜色（如果图片为空则使用）
	Opacity   float32       // 透明度 0.0-1.0
}

// SelectionEffectConfig 特效层配置
type SelectionEffectConfig struct {
	Images     [6]fyne.Resource
	ImageSizes [6]fyne.Size
	Positions  [6]*EffectPosition
	DefaultPos EffectPosition
}

// EffectPosition 特效位置
type EffectPosition struct {
	XOffset float32
	YOffset float32
	Scale   float32
}

// NewCustomList 创建自定义列表控件（应用主题颜色）
func NewCustomList(config CustomListConfig) *widget.List {
	list := widget.NewList(
		func() int {
			return len(config.Items)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(nil)
			richText := widget.NewRichTextWithText("template")

			if len(richText.Segments) > 0 {
				if textSeg, ok := richText.Segments[0].(*widget.TextSegment); ok {
					textSeg.Style = widget.RichTextStyle{
						SizeName:  theme.SizeNameSubHeadingText,
						ColorName: theme.ColorNameForeground,
					}
				}
			}

			return container.NewHBox(
				icon,
				container.NewPadded(richText),
			)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			hbox := o.(*fyne.Container)
			icon := hbox.Objects[0].(*widget.Icon)
			paddedContainer := hbox.Objects[1].(*fyne.Container)
			richText := paddedContainer.Objects[0].(*widget.RichText)

			if config.Items[id].Icon != nil {
				icon.SetResource(config.Items[id].Icon)
			} else {
				icon.SetResource(nil)
			}

			if len(richText.Segments) > 0 {
				if textSeg, ok := richText.Segments[0].(*widget.TextSegment); ok {
					textSeg.Text = config.Items[id].Title
					textSeg.Style = widget.RichTextStyle{
						SizeName:  theme.SizeNameSubHeadingText,
						ColorName: theme.ColorNameForeground,
					}
				}
			} else {
				segment := &widget.TextSegment{
					Text: config.Items[id].Title,
					Style: widget.RichTextStyle{
						SizeName:  theme.SizeNameSubHeadingText,
						ColorName: theme.ColorNameForeground,
					},
				}
				richText.Segments = []widget.RichTextSegment{segment}
			}
			richText.Refresh()
		})

	if config.OnSelected != nil {
		list.OnSelected = func(id widget.ListItemID) {
			config.OnSelected(int(id))
		}
	}

	// 设置初始选中项
	if config.InitialSelected >= 0 && config.InitialSelected < len(config.Items) {
		list.Select(config.InitialSelected)
	}

	return list
}

// NewImageBackgroundList 创建带背景的自定义列表控件
func NewImageBackgroundList(config ImageListConfig) fyne.CanvasObject {
	return NewImageBackgroundListWithEffect(config, nil).CanvasObject
}

// ImageBackgroundListWithEffect 带背景和特效的列表控件
type ImageBackgroundListWithEffect struct {
	fyne.CanvasObject
	SetOnSelected      func(func(int))
	PauseEffectUpdate  func()
	ResumeEffectUpdate func()
}

// NewImageBackgroundListWithEffect 创建带背景和特效的自定义列表控件
func NewImageBackgroundListWithEffect(config ImageListConfig, effectConfig *SelectionEffectConfig) *ImageBackgroundListWithEffect {
	// 保存原始的回调函数
	var currentOnSelected func(int) = config.ListConfig.OnSelected

	// 创建列表
	list := NewCustomList(config.ListConfig)

	// === 1. 创建独立的背景层 ===
	var background fyne.CanvasObject
	if config.Background != nil && config.Background.Image != nil {
		// 使用自定义图片控件：左上角对齐，等比例缩放
		background = NewTopLeftContainImage(config.Background.Image, config.Background.ImageSize)

		// 背景容器
		backgroundContainer := container.NewStack(background)

		// 不再需要轮询，背景控件会自动在Layout中适应尺寸
		background = backgroundContainer

	} else if config.Background != nil && config.Background.Color != nil {
		bgColor := config.Background.Color
		if config.Background.Opacity > 0 && config.Background.Opacity < 1.0 {
			bgColor = applyOpacity(bgColor, config.Background.Opacity)
		}
		rect := canvas.NewRectangle(bgColor)
		background = rect
	} else {
		// 使用主题背景色
		rect := canvas.NewRectangle(theme.Color(theme.ColorNameBackground))
		background = rect
	}

	// 如果没有特效配置，直接返回背景+列表
	if effectConfig == nil {
		stack := container.NewStack(
			background,
			list,
		)

		return &ImageBackgroundListWithEffect{
			CanvasObject: stack,
			SetOnSelected: func(onSelected func(int)) {
				currentOnSelected = onSelected
			},
		}
	}

	// === 2. 创建自定义特效控件 ===
	effectControl := NewEffectImage()
	var currentSelectedID int = config.ListConfig.InitialSelected
	estimatedItemHeight := float32(60)

	// === 3. 特效刷新函数 ===
	updateEffectPosition := func() {
		if currentSelectedID < 0 || currentSelectedID >= 6 {
			effectControl.Hide()
			return
		}

		var pos EffectPosition
		if effectConfig.Positions[currentSelectedID] != nil {
			pos = *effectConfig.Positions[currentSelectedID]
		} else {
			pos = effectConfig.DefaultPos
		}

		listSize := list.Size()
		itemTop := float32(currentSelectedID) * estimatedItemHeight
		targetHeight := estimatedItemHeight * pos.Scale
		var targetWidth float32

		if currentSelectedID < len(effectConfig.ImageSizes) {
			imgSize := effectConfig.ImageSizes[currentSelectedID]
			if imgSize.Width > 0 && imgSize.Height > 0 {
				scaleFactor := targetHeight / float32(imgSize.Height)
				targetWidth = float32(imgSize.Width) * scaleFactor
			} else {
				targetWidth = listSize.Width * pos.Scale
			}
		} else {
			targetWidth = listSize.Width * pos.Scale
		}

		y := itemTop + pos.YOffset
		centerX := (listSize.Width - targetWidth) / 2
		x := centerX + pos.XOffset

		// 边界检查
		if y < 0 {
			y = 0
		}
		if y+targetHeight > listSize.Height {
			scale := (listSize.Height - y) / targetHeight
			targetHeight *= scale
			targetWidth *= scale
		}
		if x < 0 {
			x = 0
		}
		if x+targetWidth > listSize.Width {
			scale := (listSize.Width - x) / targetWidth
			targetWidth *= scale
			targetHeight *= scale
		}

		if currentSelectedID < len(effectConfig.Images) {
			img := effectConfig.Images[currentSelectedID]
			if img != nil {
				effectControl.SetResource(img)
				effectControl.SetSize(fyne.NewSize(targetWidth, targetHeight))
				effectControl.SetPosition(fyne.NewPos(x, y))
				effectControl.Show()
				// ⚠️ 关键：这里需要刷新控件
				effectControl.Refresh()
				return
			}
		}
		effectControl.Hide()
		effectControl.Refresh()
	}

	// === 4. 列表事件处理 ===
	list.OnSelected = func(id widget.ListItemID) {
		currentSelectedID = int(id)
		// ⚠️ 关键：用fyne.Do包装
		fyne.Do(updateEffectPosition)
		if currentOnSelected != nil {
			currentOnSelected(int(id))
		}
	}

	list.OnUnselected = func(id widget.ListItemID) {
		if int(id) == currentSelectedID {
			currentSelectedID = -1
			fyne.Do(func() {
				effectControl.Hide()
				effectControl.Refresh()
			})
		}
	}

	// === 5. 特效定时器（智能刷新）===
	var paused bool
	var pauseCh = make(chan bool, 1)
	var resumeCh = make(chan bool, 1)

	// ⚠️ 关键修复：简化定时器逻辑，避免频繁刷新
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-pauseCh:
				paused = true
			case <-resumeCh:
				paused = false
			case <-ticker.C:
				if !paused && currentSelectedID >= 0 {
					// ⚠️ 关键：定时器中的UI操作必须用fyne.Do包装
					fyne.Do(updateEffectPosition)
				}
			}
		}
	}()

	// === 6. 创建三层堆栈 ===
	stack := container.NewStack(
		background,    // 底层：静态背景
		effectControl, // 中层：自定义特效控件
		list,          // 上层：透明列表
	)

	return &ImageBackgroundListWithEffect{
		CanvasObject: stack,
		SetOnSelected: func(onSelected func(int)) {
			currentOnSelected = onSelected
		},
		PauseEffectUpdate: func() {
			select {
			case pauseCh <- true:
			default:
			}
		},
		ResumeEffectUpdate: func() {
			select {
			case resumeCh <- true:
			default:
			}
		},
	}
}

// NewEffectPosition 创建特效位置配置
func NewEffectPosition(xOffset, yOffset, scale float32) *EffectPosition {
	return &EffectPosition{
		XOffset: xOffset,
		YOffset: yOffset,
		Scale:   scale,
	}
}

// applyOpacity 辅助函数：应用透明度到颜色
func applyOpacity(c color.Color, opacity float32) color.Color {
	r, g, b, a := c.RGBA()
	newAlpha := uint8(float32(uint8(a>>8)) * opacity)
	return &color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: newAlpha,
	}
}

// posChanged 检查位置是否显著变化
func posChanged(old, new fyne.Position) bool {
	return math.Abs(float64(old.X-new.X)) > 0.5 ||
		math.Abs(float64(old.Y-new.Y)) > 0.5
}

// sizeChanged 检查尺寸是否显著变化
func sizeChanged(old, new fyne.Size) bool {
	return math.Abs(float64(old.Width-new.Width)) > 0.5 ||
		math.Abs(float64(old.Height-new.Height)) > 0.5
}
