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
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// PigProgressBar 自定义进度条控件
type PigProgressBar struct {
	widget.BaseWidget
	progress       float64 // 0.0 到 1.0
	grass          *canvas.Image
	pig            *canvas.Image
	pigFrames      []string
	currentFrame   int
	isAnimating    bool
	stopChan       chan struct{}
	manualPosition *fyne.Position // 手动设置的位置，nil表示自动居中
}

// NewPigProgressBar 创建新的 PigProgressBar
func NewPigProgressBar() *PigProgressBar {
	fmt.Println("[DEBUG] NewPigProgressBar 初始化开始")
	p := &PigProgressBar{
		progress:     0,
		currentFrame: 0,
		isAnimating:  false,
		stopChan:     make(chan struct{}),
	}
	p.ExtendBaseWidget(p)

	// 加载 grass 背景图片
	grassPath := filepath.Join("Pages", "tools", "animationAsset", "pigProcessBar", "background", "grass.png")
	fmt.Println("[DEBUG] grassPath:", grassPath)
	p.grass = canvas.NewImageFromFile(grassPath)
	p.grass.FillMode = canvas.ImageFillContain
	p.grass.Resize(fyne.NewSize(800*4/5, 190*4/5)) // 640 x 152
	fmt.Printf("[DEBUG] grass image size: %v\n", p.grass.Size())

	// 加载猪动画帧
	p.pigFrames = make([]string, 0, 76)
	for i := 1; i <= 76; i++ {
		framePath := filepath.Join("Pages", "tools", "animationAsset", "pigProcessBar", "pigRunPng", fmt.Sprintf("frame_%d.png", i))
		p.pigFrames = append(p.pigFrames, framePath)
	}
	fmt.Printf("[DEBUG] pigFrames loaded: %d\n", len(p.pigFrames))

	// 初始化猪图片
	if len(p.pigFrames) > 0 {
		fmt.Println("[DEBUG] 使用第一帧初始化猪图片:", p.pigFrames[0])
		p.pig = canvas.NewImageFromFile(p.pigFrames[0])
	} else {
		fmt.Println("[DEBUG] 没有找到猪动画帧，使用空图片")
		p.pig = canvas.NewImageFromResource(nil)
	}
	p.pig.FillMode = canvas.ImageFillContain
	p.pig.Resize(fyne.NewSize(400/2, 400/2)) // 200 x 200
	fmt.Printf("[DEBUG] pig image size: %v\n", p.pig.Size())

	fmt.Println("[DEBUG] NewPigProgressBar 初始化完成")
	return p
}

// SetProgress 设置进度（外部信号，内部平滑移动）
func (p *PigProgressBar) SetProgress(target float64) {
	if target < 0 {
		target = 0
	} else if target > 1 {
		target = 1
	}

	// 启动动画
	if !p.isAnimating {
		p.startAnimation()
	}

	// 平滑动画到目标进度
	go func() {
		step := 0.0025 // 步长减半，更丝滑
		if target < p.progress {
			step = -step
		}
		for {
			if (step > 0 && p.progress >= target) || (step < 0 && p.progress <= target) {
				p.progress = target
				break
			}
			p.progress += step
			time.Sleep(4 * time.Millisecond) // 间隔减半，更高帧率
			fyne.Do(func() {
				p.Refresh()
			})
		}
		fyne.Do(func() {
			p.progress = target
			p.Refresh()
			// 到达目标后如果进度为1则停止动画
			if p.progress >= 1.0 {
				p.stopAnimation()
			}
		})
	}()
}

// SetPosition 手动设置进度条的位置（禁用自动居中）
func (p *PigProgressBar) SetPosition(x, y float32) {
	p.manualPosition = &fyne.Position{
		X: x,
		Y: y,
	}
	p.Refresh()
}

// SetAutoCenter 启用或禁用自动居中
func (p *PigProgressBar) SetAutoCenter(center bool) {
	if center {
		p.manualPosition = nil
	}
	p.Refresh()
}

// startAnimation 启动猪奔跑动画
func (p *PigProgressBar) startAnimation() {
	if p.isAnimating {
		return
	}

	p.isAnimating = true
	p.stopChan = make(chan struct{})

	go func() {
		ticker := time.NewTicker(40 * time.Millisecond) // 0.04秒/帧
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !p.isAnimating {
					return
				}

				// 只在主线程更新UI
				fyne.Do(func() {
					if len(p.pigFrames) == 0 {
						return
					}

					// 循环播放帧
					p.currentFrame = (p.currentFrame + 1) % len(p.pigFrames)
					p.pig.File = p.pigFrames[p.currentFrame]
					p.pig.Refresh()
				})

			case <-p.stopChan:
				return
			}
		}
	}()
}

// stopAnimation 停止动画（幂等，允许重复调用）
func (p *PigProgressBar) stopAnimation() {
	if !p.isAnimating {
		return
	}

	p.isAnimating = false
	select {
	case <-p.stopChan:
		// 已关闭/已触发
	default:
		close(p.stopChan)
	}
}

// Reset 重置进度条到初始状态（停止动画、回到0进度、回到第一帧）
func (p *PigProgressBar) Reset() {
	p.stopAnimation()
	p.progress = 0
	p.currentFrame = 0
	if len(p.pigFrames) > 0 {
		fyne.Do(func() {
			p.pig.File = p.pigFrames[0]
			p.pig.Refresh()
			p.Refresh()
		})
	} else {
		fyne.Do(func() { p.Refresh() })
	}
}

// CreateRenderer 创建渲染器
func (p *PigProgressBar) CreateRenderer() fyne.WidgetRenderer {
	return &pigProgressBarRenderer{
		pigProgressBar: p,
		// 先 pig 后 grass，grass 在上层
		objects: []fyne.CanvasObject{p.pig, p.grass},
	}
}

type pigProgressBarRenderer struct {
	pigProgressBar *PigProgressBar
	objects        []fyne.CanvasObject
}

func (r *pigProgressBarRenderer) Layout(size fyne.Size) {
	// grass 缩放比例应与初始化一致，4/5
	grassWidth := float32(800) * 4 / 5
	grassHeight := float32(190) * 4 / 5
	pigWidth := float32(400) / 2
	pigHeight := float32(400) / 2

	r.pigProgressBar.grass.Resize(fyne.NewSize(grassWidth, grassHeight))
	r.pigProgressBar.pig.Resize(fyne.NewSize(pigWidth, pigHeight))

	if r.pigProgressBar.manualPosition != nil {
		// 使用手动位置放置
		grassX := r.pigProgressBar.manualPosition.X
		grassY := r.pigProgressBar.manualPosition.Y
		r.pigProgressBar.grass.Move(fyne.NewPos(grassX, grassY))

		// 猪在grass上方位置基于grass
		maxX := grassWidth - pigWidth
		newX := grassX + maxX*float32(r.pigProgressBar.progress)
		newY := grassY + grassHeight - pigHeight - 20
		r.pigProgressBar.pig.Move(fyne.NewPos(newX, newY))
	} else {
		// 自动居中模式
		centerX := (size.Width - grassWidth) / 2
		centerY := (size.Height - grassHeight) / 2
		r.pigProgressBar.grass.Move(fyne.NewPos(centerX, centerY))

		maxX := grassWidth - pigWidth
		newX := centerX + maxX*float32(r.pigProgressBar.progress)
		newY := centerY + grassHeight - pigHeight - 20
		r.pigProgressBar.pig.Move(fyne.NewPos(newX, newY))
	}
}

func (r *pigProgressBarRenderer) MinSize() fyne.Size {
	return r.pigProgressBar.grass.MinSize()
}

func (r *pigProgressBarRenderer) Refresh() {
	// 只需要更新位置
	r.Layout(r.pigProgressBar.Size())
	canvas.Refresh(r.pigProgressBar)
}

func (r *pigProgressBarRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *pigProgressBarRenderer) Destroy() {
	// 停止动画
	r.pigProgressBar.stopAnimation()
}
