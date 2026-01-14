package color

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// MachineTheme 机器主题
type MachineTheme struct {
	defaultTheme fyne.Theme
}

// NewMachineTheme 创建机器主题实例
func NewMachineTheme() fyne.Theme {
	return &MachineTheme{
		defaultTheme: theme.DefaultTheme(),
	}
}

// Color 返回自定义颜色
func (t *MachineTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// 深色柔和绿色 - 常态
	darkGreen := color.NRGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0xFF}
	// 浅绿色 - 高亮状态
	lightGreen := color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	// 更浅的绿色 - 悬停状态
	lighterGreen := color.NRGBA{R: 0x66, G: 0xBB, B: 0x6A, A: 0xFF}
	// 背景色 - 非常浅的绿色
	backgroundColor := color.NRGBA{R: 0xF1, G: 0xF8, B: 0xE9, A: 0xFF}
	// 卡片背景色 - 白色
	cardBackground := color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	switch name {
	case theme.ColorNamePrimary:
		return darkGreen // 主色调使用深色柔和绿色
	case theme.ColorNameFocus:
		return lightGreen // 焦点状态使用浅绿色
	case theme.ColorNameButton:
		return darkGreen // 按钮常态使用深色柔和绿色
	case theme.ColorNameBackground:
		return backgroundColor // 页面背景色
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xFF} // 深灰色文字
	case theme.ColorNameInputBackground:
		return cardBackground // 输入框背景色
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF} // 占位符文字颜色
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0x40} // 使用主题色的浅阴影
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x43, G: 0xA0, B: 0x47, A: 0xFF} // 成功绿色
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0xFF, G: 0xB7, B: 0x4D, A: 0xFF} // 柔和的警告橙色
	case theme.ColorNameError:
		return color.NRGBA{R: 0xEF, G: 0x53, B: 0x50, A: 0xFF} // 适度的错误红色
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0x80} // 使用主题色的滚动条
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0x40} // 选中背景色
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0xB0, G: 0xB0, B: 0xB0, A: 0xFF} // 禁用状态颜色
	case theme.ColorNameHover:
		return lighterGreen // 悬停状态使用更浅的绿色
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0x80} // 输入框边框
	case theme.ColorNameOverlayBackground:
		return cardBackground // 覆盖层背景（卡片等）
	default:
		return t.defaultTheme.Color(name, variant)
	}
}

// Font 返回字体
func (t *MachineTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.defaultTheme.Font(style)
}

// Icon 返回图标
func (t *MachineTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.defaultTheme.Icon(name)
}

// Size 返回尺寸
func (t *MachineTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 10
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 5
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 18
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInputBorder:
		return 1
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputRadius:
		return 5
	default:
		return t.defaultTheme.Size(name)
	}
}
