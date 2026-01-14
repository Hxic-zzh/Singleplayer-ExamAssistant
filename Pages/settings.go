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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func SettingsPage(window fyne.Window) fyne.CanvasObject {
	// 加载当前设置
	currentConfig, err := tools.LoadSetConfig()
	if err != nil {
		currentConfig = tools.GetDefaultSetConfig()
	}

	// === 分辨率设置区块 ===
	resolutionLabel := widget.NewLabel(tools.GetLocalized("screen_resolution_settings"))
	resolutionLabel.TextStyle = fyne.TextStyle{Bold: true}

	resolutionSelect := widget.NewSelect([]string{"2560x1440", "1920x1080"}, func(selected string) {
		currentConfig.Resolution = selected
		currentConfig.CardWidth = tools.GetDefaultCardWidth(selected)
	})
	resolutionSelect.SetSelected(currentConfig.Resolution)
	resolutionSelect.PlaceHolder = tools.GetLocalized("select_resolution")

	resolutionDesc := widget.NewLabel(tools.GetLocalized("resolution_description"))
	resolutionDesc.Wrapping = fyne.TextWrapWord

	resolutionBox := container.NewVBox(
		resolutionLabel,
		resolutionSelect,
		resolutionDesc,
	)

	// === 语言设置区块 ===
	languageLabel := widget.NewLabel(tools.GetLocalized("language_settings"))
	languageLabel.TextStyle = fyne.TextStyle{Bold: true}

	languageSelect := widget.NewSelect([]string{"zh", "en"}, func(selected string) {
		currentConfig.Language = selected
		// 不再这里弹提示
	})
	languageSelect.SetSelected(currentConfig.Language)
	languageSelect.PlaceHolder = tools.GetLocalized("select_language")

	languageDesc := widget.NewLabel(tools.GetLocalized("language_description"))
	languageDesc.Wrapping = fyne.TextWrapWord

	languageBox := container.NewVBox(
		languageLabel,
		languageSelect,
		languageDesc,
	)

	// === 操作按钮 ===
	saveBtn := widget.NewButton(tools.GetLocalized("save_settings"), func() {
		if err := tools.SaveSetConfig(currentConfig); err != nil {
			widget.ShowPopUp(widget.NewLabel(tools.GetLocalized("save_failed")+err.Error()), window.Canvas())
		} else {
			widget.ShowPopUp(widget.NewLabel(tools.GetLocalized("save_success")), window.Canvas())
			dialog.ShowInformation(tools.GetLocalized("notice"), tools.GetLocalized("language_restart_notice"), window)
		}
	})

	resetBtn := widget.NewButton(tools.GetLocalized("reset_to_default"), func() {
		defaultConfig := tools.GetDefaultSetConfig()
		resolutionSelect.SetSelected(defaultConfig.Resolution)
		currentConfig.Resolution = defaultConfig.Resolution
		currentConfig.CardWidth = defaultConfig.CardWidth
		languageSelect.SetSelected(defaultConfig.Language)
		currentConfig.Language = defaultConfig.Language
	})

	buttonsBox := container.NewHBox(saveBtn, resetBtn)

	// === 主容器 - 使用滚动容器 ===
	content := container.NewVBox(
		// 分辨率设置
		resolutionBox,
		widget.NewSeparator(), // 分割线

		// 语言设置
		languageBox,
		widget.NewSeparator(), // 分割线

		// 操作按钮
		buttonsBox,
	)

	// 返回滚动容器
	return container.NewScroll(content)
}
