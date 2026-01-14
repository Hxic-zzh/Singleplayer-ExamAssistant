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
	"log"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"image"
	"image/color"

	"fmt"

	"github.com/fogleman/gg"
)

// 使用 gg 绘制加载环，尽量贴近原 CSS：
// - 三个 150px 直径的圆环，只有部分边（top/left）着色
// - 三个小点随环旋转，带发光（用多层半透明近似）
// - 文本 "Loading....." 位于底部，Times 字体近似
func generateSplashFrame(width, height int, a1, a2, a3 float64) image.Image {
	dc := gg.NewContext(width, height)
	// 背景 #111	dc.SetColor(color.RGBA{0x11, 0x11, 0x11, 0xFF})
	dc.Clear()

	cx := float64(width) / 2
	cy := float64(height) / 2
	// 原 CSS ring: width=150, border=4 -> 半径约 75
	r := 75.0
	stroke := 4.0

	toRad := func(deg float64) float64 { return deg * math.Pi / 180 }

	// 发光点绘制 - 三层固定颜色，无透明度
	drawGlowDot := func(x, y float64, base color.RGBA) {
		// 三层光晕，使用固定颜色，从外到内依次变亮
		layers := []struct {
			r     float64
			color color.RGBA
		}{
			// 外层：深色（基础颜色的 30% 亮度）
			{12, color.RGBA{
				R: uint8(float64(base.R) * 0.3),
				G: uint8(float64(base.G) * 0.3),
				B: uint8(float64(base.B) * 0.3),
				A: 255,
			}},
			// 中层：中等亮度（基础颜色的 60% 亮度）
			{8, color.RGBA{
				R: uint8(float64(base.R) * 0.6),
				G: uint8(float64(base.G) * 0.6),
				B: uint8(float64(base.B) * 0.6),
				A: 255,
			}},
			// 核心：原色（100% 亮度）
			{5, color.RGBA{
				R: base.R,
				G: base.G,
				B: base.B,
				A: 255,
			}},
		}

		// 从外到内绘制
		for _, l := range layers {
			dc.SetColor(l.color)
			dc.DrawCircle(x, y, l.r)
			dc.Fill()
		}
	}

	// ==================== 圆环配置区域 ====================
	// 三个独立的圆环，每个有自己的圆心和半径
	// 环1（蓝色）的圆心和半径
	ring1CenterX := cx - 60.0 // 左侧环的圆心 X 坐标（可修改）
	ring1CenterY := cy        // 左侧环的圆心 Y 坐标（可修改）
	ring1Radius := 75.0       // 左侧环的半径（可修改）

	// 环2（绿色）的圆心和半径
	ring2CenterX := cx - 5.0 // 中间环的圆心 X 坐标（可修改）
	ring2CenterY := cy       // 中间环的圆心 Y 坐标（可修改）
	ring2Radius := 65.0      // 中间环的半径（可修改）

	// 环3（洋红）的圆心和半径
	ring3CenterX := cx + 60.0  // 右侧环的圆心 X 坐标（可修改）
	ring3CenterY := cy - 66.66 // 右侧环的圆心 Y 坐标，上移（可修改）
	ring3Radius := 75.0        // 右侧环的半径（可修改）
	// ====================================================	// 环1：蓝色（#24ecff），border-top 着色
	// 发光点沿着环1的圆（圆心: ring1CenterX, ring1CenterY，半径: ring1Radius）旋转
	{
		// 顺时针旋转：头部（发光点）在运动方向前端
		// 🔧 修改圆弧起始位置：调整下面的 -90 这个值
		//    -90 = 从顶部开始（12点钟方向）
		//    0   = 从右侧开始（3点钟方向）
		//    90  = 从底部开始（6点钟方向）
		//    180 = 从左侧开始（9点钟方向）
		start := toRad(a1 - 90)  // CSS 动画从顶部 (270度) 开始
		end := start + toRad(90) // 1/4圆弧（90度），顺时针画到终点
		dc.Push()
		dc.SetLineWidth(stroke)
		dc.SetColor(color.RGBA{0x24, 0xEC, 0xFF, 0xFF})
		dc.DrawArc(ring1CenterX, ring1CenterY, ring1Radius, start, end)
		dc.Stroke()
		// 发光点（头部）在弧的终点，带着尾部（圆弧）移动
		dotAngle := end
		dx := ring1CenterX + ring1Radius*math.Cos(dotAngle)
		dy := ring1CenterY + ring1Radius*math.Sin(dotAngle)
		drawGlowDot(dx, dy, color.RGBA{0x24, 0xEC, 0xFF, 0xFF})
		dc.Pop()
	}

	// 环2：绿色（#93ff2d），border-left 着色
	// 发光点沿着环2的圆（圆心: ring2CenterX, ring2CenterY，半径: ring2Radius）旋转
	{
		// 逆时针旋转：头部（发光点）在运动方向前端
		// 🔧 修改圆弧起始位置：调整下面的 +180 这个值
		//    0   = 从右侧开始（3点钟方向）
		//    90  = 从底部开始（6点钟方向）
		//    180 = 从左侧开始（9点钟方向）
		//    270 = 从顶部开始（12点钟方向）
		start := toRad(a2 + 180) // 从左侧开始
		end := start - toRad(90) // 1/4圆弧（90度），逆时针画到终点
		dc.Push()
		dc.SetLineWidth(stroke)
		dc.SetColor(color.RGBA{0x93, 0xFF, 0x2D, 0xFF})
		dc.DrawArc(ring2CenterX, ring2CenterY, ring2Radius, start, end)
		dc.Stroke()
		// 发光点（头部）在弧的终点，带着尾部（圆弧）移动
		dotAngle := end
		dx := ring2CenterX + ring2Radius*math.Cos(dotAngle)
		dy := ring2CenterY + ring2Radius*math.Sin(dotAngle)
		drawGlowDot(dx, dy, color.RGBA{0x93, 0xFF, 0x2D, 0xFF})
		dc.Pop()
	}

	// 环3：洋红（#e41cf8），border-left 着色
	// 发光点沿着环3的圆（圆心: ring3CenterX, ring3CenterY，半径: ring3Radius）旋转
	{
		// 逆时针旋转：头部（发光点）在运动方向前端
		// 🔧 修改圆弧起始位置：调整下面的 +180 这个值
		//    0   = 从右侧开始（3点钟方向）
		//    90  = 从底部开始（6点钟方向）
		//    180 = 从左侧开始（9点钟方向）
		//    270 = 从顶部开始（12点钟方向）
		start := toRad(a3 + 180) // 从左侧开始
		end := start - toRad(90) // 1/4圆弧（90度），逆时针画到终点
		dc.Push()
		dc.SetLineWidth(stroke)
		dc.SetColor(color.RGBA{0xE4, 0x1C, 0xF8, 0xFF})
		dc.DrawArc(ring3CenterX, ring3CenterY, ring3Radius, start, end)
		dc.Stroke()
		// 发光点（头部）在弧的终点，带着尾部（圆弧）移动
		dotAngle := end
		dx := ring3CenterX + ring3Radius*math.Cos(dotAngle)
		dy := ring3CenterY + ring3Radius*math.Sin(dotAngle)
		drawGlowDot(dx, dy, color.RGBA{0xE4, 0x1C, 0xF8, 0xFF})
		dc.Pop()
	}

	// 文本：底部居中
	dc.SetColor(color.RGBA{0xF5, 0xF5, 0xF5, 0xFF})
	_ = dc.LoadFontFace("Times New Roman", 24)
	dc.DrawStringAnchored("Loading.....", cx, cy+r+40, 0.5, 0.5)

	return dc.Image()
}

// 供主流程触发 GIF 播放的信号
var splashReadyChan = make(chan struct{}, 1)

// 播放完成后的回调（由主流程注册）
var splashOnFinished func()

// SetSplashOnFinished 注册 GIF 播放完成后的回调
func SetSplashOnFinished(fn func()) { splashOnFinished = fn }

// NotifyReady 通知闪屏：主界面已准备好
func NotifyReady() {
	select {
	case splashReadyChan <- struct{}{}:
	default:
	}
}

// ShowStartupSplash 创建并展示启动子窗口（闪屏）。
func ShowStartupSplash(app fyne.App) fyne.Window {
	w := app.NewWindow("")
	// 初始帧 - 模拟 CSS 动画延迟效果
	// ring1: 无延迟，从 0 度开始
	// ring2: -1s 延迟 (在 4s 周期中相当于从 -90 度开始)	// ring3: -3s 延迟 (在 4s 周期中相当于从 -270 度开始)
	angle1, angle2, angle3 := 0.0, -90.0, -270.0
	frame := generateSplashFrame(512, 512, angle1, angle2, angle3)
	fyImg := canvas.NewImageFromImage(frame)
	fyImg.FillMode = canvas.ImageFillContain
	fyImg.SetMinSize(fyne.NewSize(512, 512))

	// 使用 Stack 容器包裹画布，Stack 比画布大 2px
	content := container.NewStack(fyImg)
	w.SetContent(content) // 定时刷新实现动画（使用 fyne.Do 包裹UI操作）
	// CSS 动画 4s 一圈，50ms 刷新 => 每帧旋转 360/(4000/50) = 4.5 度
	stop := make(chan struct{})
	ticker := time.NewTicker(50 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				angle1 = math.Mod(angle1+4.5, 360) // 顺时针 (animate)
				angle2 = math.Mod(angle2-4.5, 360) // 逆时针 (animate2)
				angle3 = math.Mod(angle3-4.5, 360) // 逆时针 (animate2)
				img := generateSplashFrame(512, 512, angle1, angle2, angle3)
				fyne.Do(func() {
					fyImg.Image = img
					canvas.Refresh(fyImg)
				})
			case <-stop:
				return
			}
		}
	}()

	// 使用标志位防止重复关闭 channel
	var stopClosed bool
	w.SetOnClosed(func() {
		if !stopClosed {
			close(stop)
			stopClosed = true
		}
	})

	// 阶段切换：主流程就绪 -> 继续转圈3秒 -> 播放GIF(动画系统) -> 完成后进入主界面
	go func() {
		<-splashReadyChan // 等待主流程构建完成的通知
		// 继续保持加载动画 3 秒
		time.Sleep(3 * time.Second) // 切换到序列帧动画，居中显示，尺寸 180x180
		fyne.Do(func() {
			// 停止旧的加载动画（确保不重复关闭）
			if !stopClosed {
				close(stop)
				stopClosed = true
			}
			fyImg.Hide()

			// 加载第一帧
			gifImg := canvas.NewImageFromFile("images/frame/frame_1.png")
			gifImg.FillMode = canvas.ImageFillContain
			gifImg.SetMinSize(fyne.NewSize(180, 180))
			center := container.NewCenter(gifImg)
			w.SetContent(center)

			log.Printf("[Splash] 开始播放序列帧动画 (100帧)")

			// 使用 Ticker 手动控制帧率：100帧 / 5.07秒 ≈ 50.7ms/帧
			const totalFrames = 100
			const frameDuration = 5070 * time.Millisecond / totalFrames // ≈50.7ms
			frameTicker := time.NewTicker(frameDuration)
			currentFrame := 1

			go func() {
				defer frameTicker.Stop()
				for range frameTicker.C {
					currentFrame++
					if currentFrame > totalFrames {
						// 播放完成
						log.Printf("[Splash] 序列帧动画播放完成")
						fyne.Do(func() {
							if splashOnFinished != nil {
								splashOnFinished()
							}
							w.Close()
						})
						return
					}
					// 更新帧
					framePath := fmt.Sprintf("images/frame/frame_%d.png", currentFrame)
					fyne.Do(func() {
						gifImg.File = framePath
						gifImg.Refresh()
					})
				}
			}()
		})
	}()

	// 固定大小 + 居中，去掉内边距
	// Stack 容器比画布大 2px (512 + 2 = 514)
	w.Resize(fyne.NewSize(514, 514))
	w.SetFixedSize(true)
	w.SetPadded(false)
	w.CenterOnScreen()
	w.SetCloseIntercept(func() {})
	if d, ok := any(w).(interface{ SetDecorated(bool) }); ok {
		d.SetDecorated(false)
	}
	log.Printf("[Splash] 显示闪屏窗口")
	w.Show()
	return w
}
