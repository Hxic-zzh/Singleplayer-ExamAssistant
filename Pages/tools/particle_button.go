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

// tools/particle_button.go
package tools

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
)

// Particle 粒子结构体
type Particle struct {
	X, Y        float64 // 位置
	VX, VY      float64 // 速度
	StartSize   float64 // 起始大小
	CurrentSize float64 // 当前大小
	MaxSize     float64 // 最大大小
	Color       color.RGBA
	Life        float64 // 生命周期 (0-1)
	Decay       float64 // 衰减速度
	GrowthRate  float64 // 生长速度（先增大后减小）
	IsGrowing   bool    // 是否在生长阶段
}

// ParticleButton 自定义粒子按钮
type ParticleButton struct {
	widget.BaseWidget
	OnClick func()

	// 按钮配置
	Text   string
	Width  float32
	Height float32

	// 颜色配置
	BaseColor      color.RGBA   // 基础颜色
	GradientTop    color.RGBA   // 渐变顶部颜色
	GradientBottom color.RGBA   // 渐变底部颜色
	ShadowColor    color.RGBA   // 阴影颜色
	ParticleColors []color.RGBA // 粒子颜色数组

	// 粒子系统
	particles   []*Particle
	isAnimating bool
	isPressed   bool
	isHovered   bool

	// 按钮位置和尺寸
	width, height float32
	offsetY       float32 // 垂直偏移（用于悬停/按下效果）

	// 粒子动画控制
	animationTicker *time.Ticker
	stopAnimation   chan bool
}

// NewParticleButton 创建粒子按钮（默认绿色）
func NewParticleButton(onClick func()) *ParticleButton {
	baseColor := color.RGBA{R: 143, G: 196, B: 0, A: 255} // #8fc400
	return NewParticleButtonWithColor(onClick, baseColor, "按钮")
}

// NewParticleButtonWithColor 创建带自定义颜色的粒子按钮
func NewParticleButtonWithColor(onClick func(), baseColor color.RGBA, text string) *ParticleButton {
	btn := &ParticleButton{
		OnClick:       onClick,
		Text:          text,
		Width:         120, // 默认宽度
		Height:        40,  // 默认高度
		BaseColor:     baseColor,
		particles:     make([]*Particle, 0),
		isAnimating:   false,
		isPressed:     false,
		isHovered:     false,
		width:         120,
		height:        40,
		offsetY:       0,
		stopAnimation: make(chan bool),
	}

	// 自动生成其他颜色
	btn.generateColorsFromBase()

	btn.ExtendBaseWidget(btn)

	// 启动粒子动画循环
	go btn.startAnimationLoop()

	return btn
}

// generateColorsFromBase 从基础颜色生成其他所需颜色
func (btn *ParticleButton) generateColorsFromBase() {
	r, g, b, a := btn.BaseColor.R, btn.BaseColor.G, btn.BaseColor.B, btn.BaseColor.A

	// 计算渐变顶部颜色（比基础色亮一些）
	btn.GradientTop = color.RGBA{
		R: min(255, r+30),
		G: min(255, g+30),
		B: min(255, b+30),
		A: a,
	}

	// 计算渐变底部颜色（比基础色暗一些）
	btn.GradientBottom = color.RGBA{
		R: max(0, r-30),
		G: max(0, g-30),
		B: max(0, b-30),
		A: a,
	}

	// 计算阴影颜色（比基础色暗很多）
	btn.ShadowColor = color.RGBA{
		R: max(0, r-50),
		G: max(0, g-50),
		B: max(0, b-50),
		A: a,
	}

	// 生成粒子颜色数组（基于基础颜色的不同明暗度）
	btn.ParticleColors = []color.RGBA{
		// 基础色
		btn.BaseColor,
		// 稍微亮一点
		{R: min(255, r+20), G: min(255, g+20), B: min(255, b+20), A: a},
		// 稍微暗一点
		{R: max(0, r-20), G: max(0, g-20), B: max(0, b-20), A: a},
		// 更亮
		{R: min(255, r+40), G: min(255, g+40), B: min(255, b+40), A: a},
		// 更暗
		{R: max(0, r-40), G: max(0, g-40), B: max(0, b-40), A: a},
		// 中等亮度
		{R: min(255, r+10), G: min(255, g+10), B: min(255, b+10), A: a},
	}
}

// startAnimationLoop 启动粒子动画循环
func (btn *ParticleButton) startAnimationLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			btn.UpdateParticles()
			if btn.isAnimating {
				fyne.Do(func() {
					btn.Refresh()
				})
			}
		case <-btn.stopAnimation:
			return
		}
	}
}

// Destroy 销毁按钮时停止动画
func (btn *ParticleButton) Destroy() {
	select {
	case btn.stopAnimation <- true:
	default:
	}
}

// SetSize 设置按钮尺寸
func (btn *ParticleButton) SetSize(width, height float32) {
	btn.Width = width
	btn.Height = height
	btn.width = width
	btn.height = height
	btn.Refresh()
}

// SetText 设置按钮文字
func (btn *ParticleButton) SetText(text string) {
	btn.Text = text
	btn.Refresh()
}

// SetBaseColor 设置基础颜色并重新生成其他颜色
func (btn *ParticleButton) SetBaseColor(baseColor color.RGBA) {
	btn.BaseColor = baseColor
	btn.generateColorsFromBase()
	btn.Refresh()
}

// SetColors 直接设置所有颜色
func (btn *ParticleButton) SetColors(base, top, bottom, shadow color.RGBA, particles []color.RGBA) {
	btn.BaseColor = base
	btn.GradientTop = top
	btn.GradientBottom = bottom
	btn.ShadowColor = shadow
	if particles != nil {
		btn.ParticleColors = particles
	}
	btn.Refresh()
}

// CreateParticles 创建新粒子
func (btn *ParticleButton) CreateParticles(x, y float64, count int) {
	for i := 0; i < count; i++ {
		// 随机方向角度
		angle := rand.Float64() * 2 * math.Pi

		// 适中速度
		speed := 1.5 + rand.Float64()*3.0

		// 计算速度分量
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		// 随机起始大小
		startSize := 2.0 + rand.Float64()*4.0

		// 随机最大大小
		maxSize := 10 + rand.Float64()*12

		// 随机颜色（从粒子颜色数组中选取）
		colorIdx := rand.Intn(len(btn.ParticleColors))

		particle := &Particle{
			X:           x,
			Y:           y,
			VX:          vx,
			VY:          vy,
			StartSize:   startSize,
			CurrentSize: startSize,
			MaxSize:     maxSize,
			Color:       btn.ParticleColors[colorIdx],
			Life:        1.0,
			Decay:       0.01 + rand.Float64()*0.02,
			GrowthRate:  0.2 + rand.Float64()*0.3,
			IsGrowing:   true,
		}

		btn.particles = append(btn.particles, particle)
	}
	btn.isAnimating = true
}

// UpdateParticles 更新粒子状态
func (btn *ParticleButton) UpdateParticles() {
	aliveParticles := make([]*Particle, 0)

	for _, p := range btn.particles {
		p.X += p.VX
		p.Y += p.VY

		// 轻微阻力
		p.VX *= 0.98
		p.VY *= 0.98

		// 添加重力
		p.VY += 0.15

		if p.IsGrowing {
			p.CurrentSize += p.GrowthRate
			if p.CurrentSize >= p.MaxSize {
				p.CurrentSize = p.MaxSize
				p.IsGrowing = false
			}
		} else {
			p.Life -= p.Decay
			p.CurrentSize = math.Max(0.5, p.CurrentSize*p.Life)
		}

		if p.Life > 0.05 && p.CurrentSize > 0.5 {
			aliveParticles = append(aliveParticles, p)
		}
	}

	btn.particles = aliveParticles
	btn.isAnimating = len(btn.particles) > 0
}

// DrawParticles 绘制粒子到指定上下文
func (btn *ParticleButton) DrawParticles(dc *gg.Context, offsetX, offsetY float64) {
	for _, p := range btn.particles {
		// 根据生命周期调整透明度
		alpha := uint8(p.Life * 255)
		colorWithAlpha := color.RGBA{
			R: p.Color.R,
			G: p.Color.G,
			B: p.Color.B,
			A: alpha,
		}
		dc.SetColor(colorWithAlpha)
		dc.DrawCircle(p.X+offsetX, p.Y+offsetY, p.CurrentSize)
		dc.Fill()
	}
}

// Tapped 处理点击事件
func (btn *ParticleButton) Tapped(e *fyne.PointEvent) {
	btn.isPressed = true
	btn.offsetY = 4 // CSS: margin: 8px 0 0 0
	btn.Refresh()

	// 在按钮中心位置创建粒子
	btn.CreateParticles(float64(btn.width/2), float64(btn.height/2)+float64(btn.offsetY), 25+rand.Intn(10))

	if btn.OnClick != nil {
		btn.OnClick()
	}

	// 延迟恢复按钮状态
	go func() {
		time.Sleep(150 * time.Millisecond)
		fyne.Do(func() {
			btn.isPressed = false
			if btn.isHovered {
				btn.offsetY = 2 // CSS: margin: 4px 0 0 0
			} else {
				btn.offsetY = 0 // 正常状态
			}
			btn.Refresh()
		})
	}()
}

// CreateRenderer 创建渲染器
func (btn *ParticleButton) CreateRenderer() fyne.WidgetRenderer {
	// 创建主画布用于按钮背景
	mainCanvas := canvas.NewImageFromImage(nil)
	mainCanvas.FillMode = canvas.ImageFillOriginal

	// 创建文本对象（使用Fyne的字体系统）
	text := canvas.NewText(btn.Text, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 14

	renderer := &particleButtonRenderer{
		btn:        btn,
		mainCanvas: mainCanvas,
		text:       text,
		objects:    []fyne.CanvasObject{mainCanvas, text},
	}

	// 关键修复：立即刷新一次，确保初始显示
	renderer.Refresh()

	return renderer
}

// IsAnimating 获取动画状态
func (btn *ParticleButton) IsAnimating() bool {
	return btn.isAnimating
}

// particleButtonRenderer 渲染器
type particleButtonRenderer struct {
	btn        *ParticleButton
	mainCanvas *canvas.Image
	text       *canvas.Text
	objects    []fyne.CanvasObject
}

func (r *particleButtonRenderer) Destroy() {
	r.btn.Destroy()
}

func (r *particleButtonRenderer) Layout(size fyne.Size) {
	r.btn.width = size.Width
	r.btn.height = size.Height
	r.mainCanvas.Resize(size)

	// 文本居中
	textSize := r.text.MinSize()
	textX := (size.Width - textSize.Width) / 2
	textY := (size.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(textX, textY))
	r.text.Resize(textSize)

	// 布局改变后也需要刷新
	r.Refresh()
}

func (r *particleButtonRenderer) MinSize() fyne.Size {
	// 返回合适的最小尺寸
	return fyne.NewSize(r.btn.Width, r.btn.Height+10) // 增加高度以容纳阴影
}

func (r *particleButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// drawCSSButton 绘制CSS按钮（只绘制背景和粒子）
func (r *particleButtonRenderer) drawCSSButton() *gg.Context {
	width := int(r.btn.width)
	height := int(r.btn.height)

	// 创建足够大的画布以容纳阴影
	canvasWidth := width + 10   // 左右留边距
	canvasHeight := height + 20 // 上下留边距

	dc := gg.NewContext(canvasWidth, canvasHeight)

	// 透明背景（按钮是独立的）
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	// 按钮在画布中的位置（居中）
	buttonX := float64(canvasWidth-width) / 2
	buttonY := float64(canvasHeight-height) / 2

	// 根据按钮状态设置偏移
	buttonOffsetY := float64(r.btn.offsetY)

	// 绘制按钮阴影（多层CSS阴影效果）
	if r.btn.isPressed {
		// 按下状态：没有外阴影，只有内阴影
		dc.SetRGBA(1, 1, 1, 0.2)
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), 2, 7)
		dc.Fill()

	} else if r.btn.isHovered {
		// 悬停状态：较少阴影
		// 黑色外阴影（2px模糊）
		dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+2, float64(width), float64(height), 8)
		dc.Fill()

		// 两层阴影
		for i := 0; i < 2; i++ {
			dc.SetColor(r.btn.ShadowColor)
			dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+float64(i+1), float64(width), float64(height), 8)
			dc.Stroke()
		}

		// 内阴影（白色高光）
		dc.SetRGBA(1, 1, 1, 0.2)
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), 2, 7)
		dc.Fill()

	} else {
		// 正常状态：完整的多层阴影
		// 黑色外阴影
		dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+4, float64(width), float64(height), 8)
		dc.Fill()

		// 四层阴影
		for i := 0; i < 4; i++ {
			dc.SetColor(r.btn.ShadowColor)
			dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY+float64(i+1), float64(width), float64(height), 8)
			dc.Stroke()
		}

		// 内阴影（白色高光）
		dc.SetRGBA(1, 1, 1, 0.2)
		dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), 2, 7)
		dc.Fill()
	}

	// 绘制按钮主体（渐变）
	gradient := gg.NewLinearGradient(
		buttonX, buttonY+buttonOffsetY,
		buttonX, buttonY+buttonOffsetY+float64(height),
	)
	gradient.AddColorStop(0, r.btn.GradientTop)
	gradient.AddColorStop(1, r.btn.GradientBottom)

	dc.SetFillStyle(gradient)
	dc.DrawRoundedRectangle(buttonX, buttonY+buttonOffsetY, float64(width), float64(height), 8)
	dc.Fill()

	// 绘制粒子（在按钮之上）
	if len(r.btn.particles) > 0 {
		r.btn.DrawParticles(dc, buttonX, buttonY+buttonOffsetY)
	}

	return dc
}

func (r *particleButtonRenderer) Refresh() {
	// 更新文本内容
	r.text.Text = r.btn.Text

	// 绘制按钮图像
	dc := r.drawCSSButton()
	r.mainCanvas.Image = dc.Image()
	r.mainCanvas.Refresh()
	r.text.Refresh()
}

// 辅助函数
func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
