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
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// LightboxViewer 图片查看器
type LightboxViewer struct {
	window    fyne.Window
	overlay   *canvas.Rectangle
	image     *canvas.Image
	closeBtn  *widget.Button
	container *fyne.Container
	isVisible bool
}

// NewLightboxViewer 创建新的图片查看器
func NewLightboxViewer(window fyne.Window) *LightboxViewer {
	// 创建半透明黑色遮罩
	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 200})

	// 创建图片组件
	image := canvas.NewImageFromResource(nil)
	image.FillMode = canvas.ImageFillContain

	// 创建关闭按钮
	closeBtn := widget.NewButton("×", nil)
	closeBtn.Importance = widget.HighImportance

	viewer := &LightboxViewer{
		window:    window,
		overlay:   overlay,
		image:     image,
		closeBtn:  closeBtn,
		isVisible: false,
	}

	// 设置关闭按钮点击事件
	closeBtn.OnTapped = viewer.Hide

	// 创建容器
	viewer.container = container.NewWithoutLayout(overlay, image, closeBtn)
	viewer.container.Hide()

	return viewer
}

// Show 显示图片查看器
func (l *LightboxViewer) Show(imgPath string) {
	if l.isVisible {
		return
	}

	fmt.Printf("正在显示大图: %s\n", imgPath)

	// 直接使用文件路径设置图片
	l.image.File = imgPath
	l.image.FillMode = canvas.ImageFillContain

	// 强制刷新图片
	l.image.Refresh()

	// 获取窗口尺寸
	windowSize := l.window.Canvas().Size()
	fmt.Printf("窗口尺寸: %.0fx%.0f\n", windowSize.Width, windowSize.Height)

	// 更新遮罩大小和位置
	l.overlay.Resize(windowSize)
	l.overlay.Move(fyne.NewPos(0, 0))

	// 计算图片显示尺寸（最大为窗口的80%）
	maxWidth := windowSize.Width * 0.8
	maxHeight := windowSize.Height * 0.8

	// 使用固定尺寸，因为自动检测可能有问题
	displayWidth := maxWidth
	displayHeight := maxHeight

	fmt.Printf("大图显示尺寸: %.0fx%.0f\n", displayWidth, displayHeight)

	// 设置图片位置（居中）
	posX := (windowSize.Width - displayWidth) / 2
	posY := (windowSize.Height - displayHeight) / 2

	l.image.Resize(fyne.NewSize(displayWidth, displayHeight))
	l.image.Move(fyne.NewPos(posX, posY))

	// 设置关闭按钮位置（右上角）
	closeBtnX := windowSize.Width - 60
	l.closeBtn.Resize(fyne.NewSize(40, 40))
	l.closeBtn.Move(fyne.NewPos(closeBtnX, 20))

	// 显示查看器
	l.container.Show()
	l.isVisible = true

	// 使用 overlays 显示
	l.window.Canvas().Overlays().Add(l.container)
	l.window.Canvas().Refresh(l.container)

	fmt.Println("Lightbox 大图查看器已显示")
}

// Hide 隐藏图片查看器
func (l *LightboxViewer) Hide() {
	if !l.isVisible {
		return
	}

	l.container.Hide()
	l.window.Canvas().Overlays().Remove(l.container)
	l.isVisible = false
	fmt.Println("Lightbox 查看器已隐藏")
}

// IsVisible 检查查看器是否可见
func (l *LightboxViewer) IsVisible() bool {
	return l.isVisible
}
