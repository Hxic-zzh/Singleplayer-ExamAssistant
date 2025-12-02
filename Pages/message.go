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
	"time"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MessagePage(window fyne.Window) fyne.CanvasObject { // 使用控件类型控制控件的位置和逻辑
	carouselConfig := &tools.CarouselConfig{
		X:        0, // 实际上的“自动居中”不会生效，因为被我放进layout里面去了，忽略了子控件的自操作
		Y:        0, // 那我还是觉得留着比较好，学习参考使用
		Width:    900,
		Height:   400,
		Interval: 4 * time.Second, // 差不多1s
		AutoPlay: true,
	}
	carouselItems := []*tools.CarouselItem{
		tools.NewCarouselItem("🚀 "+tools.GetLocalized("message_carousel_welcome_title"), tools.GetLocalized("message_carousel_welcome_subtitle"), 0x4169E1, ""),
		tools.NewCarouselItem("🔔 "+tools.GetLocalized("message_carousel_notice_title"), fmt.Sprintf(tools.GetLocalized("message_carousel_notice_subtitle"), "2025-11-23-1"), 0x228B22, ""),
		tools.NewCarouselItem("", "", 0, "Pages/assest/carousel1.png"),
		tools.NewCarouselItem("", "", 0, "Pages/assest/carousel2.png"),
		tools.NewCarouselItem("", "", 0, "Pages/assest/carousel3.png"),
	}
	carousel := tools.NewCarouselWithItems(carouselConfig, carouselItems)

	title := canvas.NewText(tools.GetLocalized("message_title"), color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	title.TextSize = 20
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	statusText := widget.NewLabel(tools.GetLocalized("message_status_autoplay"))
	progressText := widget.NewLabel("")

	// ！！！如果要在函数内部定义一个子函数来访问外部函数的局部变量（如 carousel 和 progressText），必须使用匿名函数（闭包）
	updateProgress := func() { // func 匿名函数 直接赋给updateProgress，然后updateProgress()马上去调用这个匿名函数
		current := carousel.GetCurrentIndex() + 1
		total := carousel.GetItemsCount()
		progressText.SetText(fmt.Sprintf("%d / %d", current, total))
	}
	updateProgress()

	prevBtn := widget.NewButton("◀️ "+tools.GetLocalized("message_prev_btn"), func() {
		carousel.Previous() // 来自carousel.go -> 关键代码：c.container.Objects = []fyne.CanvasObject{c.items[c.currentIndex]}
		statusText.SetText(tools.GetLocalized("message_status_prev"))
		updateProgress()
	})
	nextBtn := widget.NewButton(tools.GetLocalized("message_next_btn")+" ▶️", func() {
		carousel.Next()
		statusText.SetText(tools.GetLocalized("message_status_next"))
		updateProgress()
	})
	stopBtn := widget.NewButton("⏹️ "+tools.GetLocalized("message_stop_btn"), func() {
		carousel.Stop()
		statusText.SetText(tools.GetLocalized("message_status_stop"))
	})

	controls := container.NewHBox(prevBtn, nextBtn, stopBtn)                        // 水平 -- 按钮 -- 控制轮盘
	rouletteMass := container.NewHBox(statusText, layout.NewSpacer(), progressText) // 水平 -- 杂类 -- 轮盘附属

	// 使用 WithoutLayout 手动定位
	container := container.NewWithoutLayout(
		title,
		carousel,
		controls,
		rouletteMass,
	)

	// 手动设置位置和大小
	title.Move(fyne.NewPos(400, 10))
	title.Resize(fyne.NewSize(200, 30))

	carousel.Move(fyne.NewPos(50, 50)) // 轮播内容
	carousel.Resize(fyne.NewSize(900, 400))

	controls.Move(fyne.NewPos(350, 470))
	controls.Resize(fyne.NewSize(200, 40))

	rouletteMass.Move(fyne.NewPos(0, 550))
	rouletteMass.Resize(fyne.NewSize(900, 30))

	return container
}
