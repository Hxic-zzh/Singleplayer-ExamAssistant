// 该文件专门写样式逻辑
/*


            $$\               $$\
            $$ |              $$ |
 $$$$$$$\ $$$$$$\   $$\   $$\ $$ | $$$$$$\
$$  _____|\_$$  _|  $$ |  $$ |$$ |$$  __$$\
\$$$$$$\    $$ |    $$ |  $$ |$$ |$$$$$$$$ |
 \____$$\   $$ |$$\ $$ |  $$ |$$ |$$   ____|
$$$$$$$  |  \$$$$  |\$$$$$$$ |$$ |\$$$$$$$\
\_______/    \____/  \____$$ |\__| \_______|
                    $$\   $$ |
                    \$$$$$$  |
                     \______/


*/

package tools

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// === ⓵轮播图项 ===
// 控件的第一步一定是准备好我要修改什么，要修改什么就放到结构体里面
// 在 Fyne 中，属性设置顺序不严格影响结果,在任何地方适用
type CarouselItem struct {
	Title     string      //标题文本
	Content   string      //内容文本
	Color     color.Color //背景颜色
	ImagePath string      //图片路径
}

// === ⓶创建轮播图项 ===
// 要显示的写清楚
func NewCarouselItem(title, content string, colorValue uint32, imagePath string) *CarouselItem {
	return &CarouselItem{
		Title:     title,
		Content:   content,
		ImagePath: imagePath,
		Color: color.NRGBA{
			R: uint8((colorValue >> 16) & 0xFF), // 提取红色分量
			G: uint8((colorValue >> 8) & 0xFF),  // 提取绿色分量
			B: uint8(colorValue & 0xFF),         // 提取蓝色分量
			A: 255,                              // 不透明度
		},
	}
}

// === ⓷创建轮播图项的单个可视化组件 ===
// 有什么属性要改?有什么功能要实现? -> 改好 -> 然后就要去将内容显示出来（内容为何+区域大小+区域位置）
func createCarouselItem(item *CarouselItem, width, height float32) fyne.CanvasObject {
	var bg fyne.CanvasObject // 定义成画布
	if item.ImagePath != "" {
		img := canvas.NewImageFromFile(item.ImagePath) // 先把路径给画板
		img.FillMode = canvas.ImageFillStretch         // ✔控制图片如何填充，我这里就是不做任何的操作，我本来就是要做海报的，就把海报的大小做对就可以了
		img.Resize(fyne.NewSize(width, height))        // ✔控制需要填充的大小
		bg = img                                       // bg是“通用接口”，所以不能直接对bg来回使用方法进行配置，需要另定义一个专门的“图像对象”，然后再转给它
	} else {
		bgRect := canvas.NewRectangle(item.Color)
		bgRect.Resize(fyne.NewSize(width, height))
		bg = bgRect // 可以看出来这个“通用接口”的用处，就是让一个控件适应各种情况，所以为了防止意外的错误，一般“通用接口”只能传送赋值，不能用具体的函数改
	}

	// ==| 标题文本 |==
	title := canvas.NewText(item.Title, color.Black)
	title.TextSize = 24
	title.TextStyle.Bold = true
	title.Move(fyne.NewPos(40, 40))

	// ==| 内容文本 |==
	content := widget.NewLabel(item.Content)
	content.Move(fyne.NewPos(40, 80))

	return container.NewWithoutLayout(bg, title, content) // 使用没有布局的容器，手动定位 NewWithoutLayout记牢了！
}

// === ⓸创建轮播图的整体架构 ===
func NewCarouselWithItems(config *CarouselConfig, items []*CarouselItem) *Carousel {
	canvasItems := make([]fyne.CanvasObject, len(items))
	for i, item := range items {
		canvasItems[i] = createCarouselItem(item, config.Width, config.Height) //用上面的创建单个项目的函数，一个一个将项目转换到CanvasObject类型
	}

	return NewCarousel(config, canvasItems) // 返回carousel.go的函数
}
