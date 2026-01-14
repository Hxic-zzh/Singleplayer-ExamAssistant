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
	"errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 页面状态
type ExtractPageState struct {
	mainFiles   []string
	extraFiles  []string
	bankName    string
	parseResult *tools.ParseResult
}

func ExtractPage(window fyne.Window) fyne.CanvasObject {
	// 状态管理
	state := &ExtractPageState{}

	// 根据主题设置按钮颜色
	var baseColor color.RGBA
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantDark {
		baseColor = color.RGBA{R: 60, G: 120, B: 60, A: 255} // 深色主题用深绿色
	} else {
		baseColor = color.RGBA{R: 143, G: 196, B: 0, A: 255} // 浅色主题用原色
	}

	// === 名称区块 ===
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder(tools.GetLocalized("extract_input_bank_name"))
	saveNameBtn := tools.NewParticleButtonWithColor(func() {
		state.bankName = nameEntry.Text
		if state.bankName != "" {
			fmt.Println(tools.GetLocalized("extract_bank_name_saved"), state.bankName)
		}
	}, baseColor, tools.GetLocalized("extract_save_bank_name"))
	saveNameBtn.SetSize(150, 40)

	// === 只读文本区域 for JSON ===
	jsonText := widget.NewMultiLineEntry()
	jsonText.SetPlaceHolder(tools.GetLocalized("extract_json_preview_placeholder"))

	// === 只读文本区域 for Markdown ===
	markdownText := widget.NewMultiLineEntry()
	markdownText.SetPlaceHolder(tools.GetLocalized("extract_markdown_preview_placeholder"))

	// === 只读文本区域 for status and progress ===
	statusTextEntry := widget.NewMultiLineEntry()
	statusTextEntry.SetText(tools.GetLocalized("extract_please_select_file"))

	// === 将保存json到什么文件夹 ===
	var selectFolderBtn *tools.ParticleButton

	selectFolderBtn = tools.NewParticleButtonWithColor(func() {
		tools.SelectFolder(window, func(selectedPath string) {
			displayText := tools.TruncatePathSmart(selectedPath, 35)
			selectFolderBtn.SetText(tools.GetLocalized("extract_selected_folder") + displayText)
			fmt.Println(tools.GetLocalized("extract_selected_folder_path"), selectedPath)
		})
	}, baseColor, tools.GetLocalized("extract_select_folder"))
	selectFolderBtn.SetSize(400, 40)

	go func() {
		if tempData, err := tools.LoadTempData(); err == nil && tempData.SelectedFolder != "" {
			displayText := tools.TruncatePathSmart(tempData.SelectedFolder, 35)
			selectFolderBtn.SetText(tools.GetLocalized("extract_selected_folder") + displayText)
		}
	}()

	// === 清空数据按钮 ===
	// 这个函数一定要放在后面，有用到界面
	clearDataBtn := tools.NewParticleButtonWithColor(func() {
		// 清空所有状态
		state.mainFiles = []string{}
		state.extraFiles = []string{}
		state.bankName = ""
		state.parseResult = nil

		// 清空界面显示
		nameEntry.SetText("")
		jsonText.SetText("")
		markdownText.SetText("")
		statusTextEntry.SetText(tools.GetLocalized("extract_please_select_file"))

		// 重置按钮显示
		selectFolderBtn.SetText(tools.GetLocalized("extract_select_folder"))

		// 新增：删除所有缓存的图片和临时文件（不删除输出ZIP）
		clearAllCacheData()

		// 显示成功提示
		dialog.ShowInformation(tools.GetLocalized("extract_clear_complete"), tools.GetLocalized("extract_all_cache_cleared"), window)
	}, baseColor, tools.GetLocalized("extract_clear_data"))
	clearDataBtn.SetSize(120, 40)

	// 更新状态显示的函数
	updateStatusDisplay := func() {
		var statusText string

		if len(state.mainFiles) > 0 {
			statusText += "📁 主文件:\n"
			for _, file := range state.mainFiles {
				filename := filepath.Base(file)
				// 显示文件类型
				fileType, _, err := tools.GetQuestionType(file)
				typeDesc := ""
				if err == nil {
					switch fileType {
					case tools.SingleChoice:
						typeDesc = " (单选题)"
					case tools.MultipleChoice:
						typeDesc = " (多选题)"
					case tools.FillBlank:
						typeDesc = " (填空题)"
					case tools.SingleChoiceImg:
						typeDesc = " (题干是图单选)"
					case tools.MultipleChoiceImg:
						typeDesc = " (题干是图多选)"
					case tools.FillBlankImg:
						typeDesc = " (题干是图填空)"
					case tools.DocumentReading:
						typeDesc = " (材料阅读题)"
					}
				}
				statusText += fmt.Sprintf("  ✅ %s%s\n", filename, typeDesc)
			}
		}

		if len(state.extraFiles) > 0 {
			statusText += "\n📎 辅助文件:\n"
			for _, file := range state.extraFiles {
				filename := filepath.Base(file)
				statusText += fmt.Sprintf("  🔗 %s\n", filename)
			}
		}

		if state.parseResult != nil {
			statusText += "\n📊 解析统计:\n"

			// 统计各种类型的题目数量
			scCount, scImgCount := countQuestionTypes(state.parseResult.SingleChoice)
			mcCount, mcImgCount := countQuestionTypesMC(state.parseResult.MultipleChoice)
			flCount, flImgCount := countQuestionTypesFL(state.parseResult.FillBlank)

			statusText += fmt.Sprintf("  单选题: %d (普通: %d, 题干是图: %d)\n", len(state.parseResult.SingleChoice), scCount, scImgCount)
			statusText += fmt.Sprintf("  多选题: %d (普通: %d, 题干是图: %d)\n", len(state.parseResult.MultipleChoice), mcCount, mcImgCount)
			statusText += fmt.Sprintf("  填空题: %d (普通: %d, 题干是图: %d)\n", len(state.parseResult.FillBlank), flCount, flImgCount)
			statusText += fmt.Sprintf("  材料阅读题: %d\n", len(state.parseResult.DocumentReading))
		}

		if statusText == "" {
			statusText = tools.GetLocalized("extract_please_select_file")
		}

		statusTextEntry.SetText(statusText)
	}

	// === xlsx文件选择区块 ===
	openMainBtn := tools.NewParticleButtonWithColor(func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}

			filePath := reader.URI().Path()
			// 检查文件类型 - 现在支持所有主文件类型
			fileType, _, err := tools.GetQuestionType(filePath)
			if err != nil || (fileType != tools.SingleChoice && fileType != tools.MultipleChoice && fileType != tools.FillBlank &&
				fileType != tools.SingleChoiceImg && fileType != tools.MultipleChoiceImg && fileType != tools.FillBlankImg &&
				fileType != tools.DocumentReading) {
				dialog.ShowError(errors.New(tools.GetLocalized("extract_select_valid_file")), window)
				return
			}

			state.mainFiles = append(state.mainFiles, filePath)
			updateStatusDisplay()
		}, window)

		// 设置文件过滤器
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fileDialog.Show()
	}, baseColor, tools.GetLocalized("extract_open_main_file"))
	openMainBtn.SetSize(190, 40)

	openAuxBtn := tools.NewParticleButtonWithColor(func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}

			filePath := reader.URI().Path()
			// 检查文件类型
			fileType, _, err := tools.GetQuestionType(filePath)
			if err != nil || fileType != tools.FillExtra {
				dialog.ShowError(errors.New(tools.GetLocalized("extract_select_fe_file")), window)
				return
			}

			state.extraFiles = append(state.extraFiles, filePath)
			updateStatusDisplay()
		}, window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fileDialog.Show()
	}, baseColor, tools.GetLocalized("extract_open_aux_file"))
	openAuxBtn.SetSize(190, 40)

	// === 预览和保存按钮区块 ===
	generatePreviewBtn := tools.NewParticleButtonWithColor(func() {
		if len(state.mainFiles) == 0 {
			dialog.ShowError(errors.New(tools.GetLocalized("extract_select_main_file_first")), window)
			return
		}

		// 解析文件
		result, err := tools.ParseQuestionFiles(state.mainFiles, state.extraFiles)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		state.parseResult = result

		// 生成预览
		bankName := nameEntry.Text
		if bankName == "" {
			bankName = tools.GetLocalized("extract_unnamed_bank")
		}

		// JSON预览
		jsonPreview, err := tools.GenerateJSONPreview(result, bankName)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			jsonText.SetText(jsonPreview)
		}

		// Markdown预览
		mdPreview := tools.GenerateMarkdownPreview(result, bankName)
		markdownText.SetText(mdPreview)

		// 更新状态
		updateStatusDisplay()
	}, baseColor, tools.GetLocalized("extract_generate_preview"))
	generatePreviewBtn.SetSize(190, 40)

	saveBtn := tools.NewParticleButtonWithColor(func() {
		if state.parseResult == nil {
			dialog.ShowError(errors.New(tools.GetLocalized("extract_generate_preview_first")), window)
			return
		}

		// 获取保存路径
		tempData, err := tools.LoadTempData()
		if err != nil || tempData.SelectedFolder == "" {
			dialog.ShowError(errors.New(tools.GetLocalized("extract_select_save_folder_first")), window)
			return
		}

		// 生成文件名
		bankName := nameEntry.Text
		var fileName string
		if bankName == "" {
			fileName = tools.GenerateFileName("")
		} else {
			fileName = bankName + ".json"
		}

		jsonPath := filepath.Join(tempData.SelectedFolder, fileName)

		// 使用新的保存函数（包含图片打包）
		err = tools.SaveQuestionBankWithImages(state.parseResult, bankName, jsonPath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("%s: %v", tools.GetLocalized("extract_save_failed"), err), window)
			return
		}

		// 保存成功后清理导入的临时文件
		importTempPath := filepath.Join("data", "temp", "imported_files")
		if _, err := os.Stat(importTempPath); err == nil {
			os.RemoveAll(importTempPath)
			fmt.Println(tools.GetLocalized("extract_import_temp_files_cleared"))
		}

		dialog.ShowInformation(tools.GetLocalized("extract_save_success"), fmt.Sprintf(tools.GetLocalized("extract_bank_and_images_saved_to"), strings.TrimSuffix(jsonPath, ".json")), window)
	}, baseColor, tools.GetLocalized("extract_save_bank"))
	saveBtn.SetSize(190, 40)

	// === 图片和题库导入按钮 ===
	importImagesBtn := tools.NewParticleButtonWithColor(func() {
		// 打开ZIP文件选择对话框
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}

			zipPath := reader.URI().Path()
			if filepath.Ext(zipPath) != ".zip" {
				dialog.ShowError(errors.New(tools.GetLocalized("extract_select_zip_file")), window)
				return
			}

			// 处理ZIP文件
			imageCount, err := tools.ProcessImageZip(zipPath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("%s: %v", tools.GetLocalized("extract_import_images_failed"), err), window)
				return
			}

			// 更新状态显示
			dialog.ShowInformation(tools.GetLocalized("extract_import_success"), fmt.Sprintf(tools.GetLocalized("extract_images_imported_successfully"), imageCount), window)
			updateStatusDisplay()
		}, window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".zip"}))
		fileDialog.Show()
	}, baseColor, tools.GetLocalized("extract_import_images"))
	importImagesBtn.SetSize(190, 40)

	importQuestionBankBtn := tools.NewParticleButtonWithColor(func() {
		// 打开ZIP文件选择对话框
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}

			zipPath := reader.URI().Path()
			if filepath.Ext(zipPath) != ".zip" {
				dialog.ShowError(errors.New(tools.GetLocalized("extract_select_zip_file")), window)
				return
			}

			// 显示加载中对话框
			progressDialog := dialog.NewProgress(tools.GetLocalized("extract_import_bank"), tools.GetLocalized("extract_parsing_bank_file"), window)
			progressDialog.Show()

			// 在goroutine中处理耗时操作
			go func() {
				// 处理题库ZIP文件 - 现在返回解析结果和文件路径
				result, mainFiles, extraFiles, err := tools.ImportQuestionBankZip(zipPath)

				// 使用 fyne.Do 在主线程安全地更新UI
				fyne.Do(func() {
					progressDialog.Hide()

					if err != nil {
						dialog.ShowError(fmt.Errorf("%s: %v", tools.GetLocalized("extract_import_bank_failed"), err), window)
						return
					}

					// 更新状态 - 完全模拟用户手动操作
					state.parseResult = result
					state.mainFiles = mainFiles   // 设置实际的文件路径
					state.extraFiles = extraFiles // 设置实际的文件路径

					// 生成预览
					bankName := nameEntry.Text
					if bankName == "" {
						bankName = tools.GetLocalized("extract_imported_bank")
					}

					// JSON预览
					jsonPreview, err := tools.GenerateJSONPreview(result, bankName)
					if err != nil {
						dialog.ShowError(err, window)
					} else {
						jsonText.SetText(jsonPreview)
					}

					// Markdown预览
					mdPreview := tools.GenerateMarkdownPreview(result, bankName)
					markdownText.SetText(mdPreview)

					// 更新状态显示
					updateStatusDisplay()

					dialog.ShowInformation(tools.GetLocalized("extract_import_success"),
						fmt.Sprintf(tools.GetLocalized("extract_import_summary"),
							len(result.SingleChoice),
							len(result.MultipleChoice),
							len(result.FillBlank),
							len(result.DocumentReading),
							tools.GetTempImageCount(),
							len(result.Errors)),
						window)
				})
			}()
		}, window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".zip"}))
		fileDialog.Show()
	}, baseColor, tools.GetLocalized("extract_import_bank"))
	importQuestionBankBtn.SetSize(190, 40)

	// === 手动布局区块 ===
	content := container.NewWithoutLayout(
		selectFolderBtn,
		openMainBtn,
		openAuxBtn,
		generatePreviewBtn,
		saveBtn,
		importImagesBtn,
		importQuestionBankBtn,
		jsonText,
		nameEntry,
		saveNameBtn,
		clearDataBtn,
		markdownText,
		statusTextEntry,
	)

	// 选择文件夹按钮
	selectFolderBtn.Move(fyne.NewPos(-100, 10))
	selectFolderBtn.Resize(fyne.NewSize(600, 45))

	// 打开主文件按钮
	openMainBtn.Move(fyne.NewPos(-20, 60))
	openMainBtn.Resize(fyne.NewSize(240, 45))

	// 打开辅助文件按钮
	openAuxBtn.Move(fyne.NewPos(180, 60))
	openAuxBtn.Resize(fyne.NewSize(240, 45))

	// 生成预览按钮
	generatePreviewBtn.Move(fyne.NewPos(-20, 110))
	generatePreviewBtn.Resize(fyne.NewSize(240, 45))

	// 保存题库按钮
	saveBtn.Move(fyne.NewPos(180, 110))
	saveBtn.Resize(fyne.NewSize(240, 45))

	importImagesBtn.Move(fyne.NewPos(-20, 160))
	importImagesBtn.Resize(fyne.NewSize(240, 45))

	importQuestionBankBtn.Move(fyne.NewPos(180, 160))
	importQuestionBankBtn.Resize(fyne.NewSize(240, 45))

	// JSON 显示文本区域
	jsonText.Move(fyne.NewPos(10, 210))
	jsonText.Resize(fyne.NewSize(400, 450))

	// 题库名称输入框
	nameEntry.Move(fyne.NewPos(420, 10))
	nameEntry.Resize(fyne.NewSize(200, 40))

	// 保存题库名称按钮
	saveNameBtn.Move(fyne.NewPos(600, 10))
	saveNameBtn.Resize(fyne.NewSize(240, 45))

	// 清空数据按钮
	clearDataBtn.Move(fyne.NewPos(800, 10))
	clearDataBtn.Resize(fyne.NewSize(120, 45))

	// Markdown 预览文本区域
	markdownText.Move(fyne.NewPos(420, 60))
	markdownText.Resize(fyne.NewSize(500, 300))

	// 状态和进度文本区域
	statusTextEntry.Move(fyne.NewPos(420, 370))
	statusTextEntry.Resize(fyne.NewSize(500, 290))

	return content
}

// clearAllCacheData 清除所有缓存数据（临时图片和临时文件）
func clearAllCacheData() {
	cacheDirs := []string{
		filepath.Join("data", "temp", "tempImages"),
		filepath.Join("data", "temp", "add"),
		filepath.Join("data", "output", "outputTemp"),
		filepath.Join("data", "temp", "import_temp"),
		filepath.Join("data", "temp", "imported_files"), // 新增：清理导入的临时文件
	}

	// 清理临时数据目录
	for _, dir := range cacheDirs {
		if _, err := os.Stat(dir); err == nil {
			files, err := os.ReadDir(dir)
			if err == nil {
				for _, file := range files {
					filePath := filepath.Join(dir, file.Name())
					os.RemoveAll(filePath)
				}
				fmt.Printf(tools.GetLocalized("extract_directory_cleared"), dir)
			}
		}
	}

	// 清理data目录中的tempData.json（保存的文件夹路径配置）
	tempDataFile := filepath.Join("data", "tempData.json")
	if _, err := os.Stat(tempDataFile); err == nil {
		os.Remove(tempDataFile)
		fmt.Println(tools.GetLocalized("extract_temp_data_cleared"))
	}

	fmt.Println(tools.GetLocalized("extract_all_cache_cleared_keep_zip"))
}

// 统计单选题类型数量
func countQuestionTypes(questions []tools.SingleChoiceQuestion) (int, int) {
	normalCount := 0
	imgCount := 0

	for _, q := range questions {
		if q.Type == tools.SingleChoiceImg {
			imgCount++
		} else {
			normalCount++
		}
	}

	return normalCount, imgCount
}

// 统计多选题类型数量
func countQuestionTypesMC(questions []tools.MultipleChoiceQuestion) (int, int) {
	normalCount := 0
	imgCount := 0

	for _, q := range questions {
		if q.Type == tools.MultipleChoiceImg {
			imgCount++
		} else {
			normalCount++
		}
	}

	return normalCount, imgCount
}

// 统计填空题类型数量
func countQuestionTypesFL(questions []tools.FillBlankQuestion) (int, int) {
	normalCount := 0
	imgCount := 0

	for _, q := range questions {
		if q.Type == tools.FillBlankImg {
			imgCount++
		} else {
			normalCount++
		}
	}

	return normalCount, imgCount
}

// SetSize
