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
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// === ⓵CarouselConfig（轮播图配置）结构体 ===
// 新创建任何控件都需要先确定好大小和位置，然后考虑该自定义控件继承谁？有什么属性需要从外部引入。在这里定义外部引入的属性
// 这里并非不能使用指针，而是只写属性是有利于逻辑的
type CarouselConfig struct {
	X, Y, Width, Height float32
	Interval            time.Duration // 图片切换所需间隔
	AutoPlay            bool          // 是否自动播放判断标志
}

// === ⓶Carousel （轮播图组件） ===
// 然后做继承，将控件的父类的功能继承过来；再是将第一步的内容拷贝到该结构中；最后将需要用到的数组以及指针：“子项”“容器”“外部方法” 还有关键的“信号通道”建立起来
type Carousel struct {
	widget.BaseWidget                     // 只要是Fyne的基础Widget类，都要写这个。用来继承组件的行为方法
	config            *CarouselConfig     // 每一个轮播图应该都有这个结构体的属性
	items             []fyne.CanvasObject // Canvas肯定是与绘制有关的。在这里意味用“切片”存放所有需要渲染的子项
	currentIndex      int                 // 常见的定位当前 显示何子项，用来上下翻页
	container         *fyne.Container     // 放items的地方
	ticker            *time.Ticker        // 帧，服务于自动计时器，下面有一个专门的方法用来发送信号；也相当于从time.Ticker里面继承必要的函数，总之就是要写
	done              chan bool           // 布尔值类型的通道，用来发信号
}

// === ⓷创建轮播图逻辑 ===
// 变量写好了就开始写逻辑
// 创建的时候是需要控件的位置，大小，样式...这些东西不在这个文件里面，所以这个函数是被调用的，在carousel_item.go里面调用
// ps：这是构造函数  func 函数名(参数) 返回值类型 {
func NewCarousel(config *CarouselConfig, items []fyne.CanvasObject) *Carousel { // 既然要创建控件，就要返回自定义好的控件类型
	c := &Carousel{
		config:       config,
		items:        items,
		currentIndex: 0,
		done:         make(chan bool), // chan bool 是引用类型，所以要用上了，就要make一下，再结构体里面初始定义的时候就不用make
	}
	c.ExtendBaseWidget(c) // 用BaseWidget来保证控件拥有基本功能和标准行为

	c.container = container.NewStack() // 在容器中先创建堆栈
	if len(items) > 0 {
		c.container.Objects = []fyne.CanvasObject{items[0]}
	}

	c.container.Resize(fyne.NewSize(config.Width, config.Height))
	c.container.Move(fyne.NewPos(config.X, config.Y))

	if config.AutoPlay && len(items) > 1 {
		c.startAutoPlay() // 是内部方法，不用多想。启动定时器
	}

	return c
}

// === ⓸创建渲染器 ===
// 逻辑写好了就让控件渲染到界面上
// ps：这是结构体方法，所以语法不一样  func (接收者) 方法名(参数) 返回值类型 {
func (c *Carousel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.container)
}

// === ⓹其他功能制作 ===
// ==| 开始自动播放 |==
func (c *Carousel) startAutoPlay() {
	c.ticker = time.NewTicker(c.config.Interval)

	go func() { // go func语法，用来做并发执行
		for { // 啥也没写的for是无限循环
			select { // select语句是用来监听通道的
			case <-c.ticker.C: // ticker.C是默认的方法，专门用来给通道发送“计时达到”的信号
				fyne.Do(c.Next) // 一定要用主线程！这个Next方法在go里面写在哪都可以，不用向上说明
			case <-c.done: // 这里也说明白了，done这个通道是用来发送停止信号的，而且不关心发送了什么信号，所以安全性很差，但是我懒得管
				return
			}
		}
	}()
}

// ==| 下一张 |==
func (c *Carousel) Next() {
	if len(c.items) <= 1 { // 其实在这里搞这么一句，没啥意思，真有逻辑问题，也不是在逻辑的末尾处理
		return
	}

	c.currentIndex = (c.currentIndex + 1) % len(c.items)               // 模运算是为了循环（0-4）
	c.container.Objects = []fyne.CanvasObject{c.items[c.currentIndex]} // 显示逻辑
	c.container.Refresh()                                              // 别忘了刷新
}

// ==| 上一张 |==
func (c *Carousel) Previous() {
	if len(c.items) <= 1 {
		return
	}

	c.currentIndex = (c.currentIndex - 1 + len(c.items)) % len(c.items)
	c.container.Objects = []fyne.CanvasObject{c.items[c.currentIndex]}
	c.container.Refresh()
}

// ==| 停止自动播放 |==
func (c *Carousel) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
		select {
		case c.done <- true:
		default:
		}
	}
}

// 那个啥进度条要用，和上面的没关系，因为是操作逻辑，所以放在这个文件里面，管理方便点
// ==| 获取当前索引 |==
func (c *Carousel) GetCurrentIndex() int {
	return c.currentIndex
}

// ==| 获取项目数量 |==
func (c *Carousel) GetItemsCount() int {
	return len(c.items)
}
