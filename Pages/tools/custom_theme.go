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
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// GoldTheme 自定义金黄色主题
type GoldTheme struct {
	defaultTheme fyne.Theme
}

var _ fyne.Theme = (*GoldTheme)(nil)

// NewGoldTheme 创建金黄色主题
func NewGoldTheme() fyne.Theme {
	return &GoldTheme{
		defaultTheme: theme.DefaultTheme(),
	}
}

// Color 覆盖颜色设置
func (t *GoldTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// 覆盖选中颜色为浅金黄色半透明
	if name == theme.ColorNameSelection {
		return &color.NRGBA{
			R: 255, // 红色
			G: 236, // 绿色 - 偏黄
			B: 139, // 蓝色 - 偏黄
			A: 32,
		}
	}

	// 覆盖悬停颜色（如果需要）
	if name == theme.ColorNameHover {
		return &color.NRGBA{
			R: 255, // 红色
			G: 245, // 更浅的绿色
			B: 157, // 更浅的蓝色
			A: 128,
		}
	}

	// 其他颜色使用默认主题
	return t.defaultTheme.Color(name, variant)
}

// Font 字体设置（使用默认）
func (t *GoldTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.defaultTheme.Font(style)
}

// Icon 图标设置（使用默认）
func (t *GoldTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.defaultTheme.Icon(name)
}

// Size 尺寸设置（使用默认）
func (t *GoldTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.defaultTheme.Size(name)
}
