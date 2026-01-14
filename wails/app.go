package main

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用结构
type App struct {
	ctx context.Context
}

// NewApp 创建新的应用实例
func NewApp() *App {
	return &App{}
}

// startup 在应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 启动时清理 fix 目录（仅清空内部内容，保留目录本身）
	fixDir := "fix"
	if err := os.MkdirAll(fixDir, 0755); err != nil {
		fmt.Printf("警告: 创建 fix 目录失败: %v\n", err)
	} else {
		if err := clearDirectoryContents(fixDir); err != nil {
			fmt.Printf("警告: 清空 fix 目录失败: %v\n", err)
		} else {
			fmt.Println("✓ 已清空 fix 目录内容")
		}
	}

	// 清理上次运行遗留的临时图片文件（仅清空内容，保留目录结构）
	if err := a.ClearTempImages(); err != nil {
		fmt.Printf("警告: 清理临时图片失败: %v\n", err)
	} else {
		fmt.Println("✓ 已清理临时图片目录")
	}
}

// clearDirectoryContents 递归删除 dir 下的所有文件和子目录，但保留 dir 本身
func clearDirectoryContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		entryPath := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return err
		}
	}
	return nil
}

// ========== 数据结构定义 ==========

// QuestionBank 题库结构
type QuestionBank struct {
	Name      string      `json:"name"`
	Version   string      `json:"version"`
	Metadata  Metadata    `json:"metadata"`
	Images    []ImageInfo `json:"images"`
	Questions Questions   `json:"questions"`
	Errors    []string    `json:"errors"`
}

// Metadata 元数据
type Metadata struct {
	TotalQuestions  int `json:"totalQuestions"`
	SingleChoice    int `json:"singleChoice"`
	MultipleChoice  int `json:"multipleChoice"`
	FillBlank       int `json:"fillBlank"`
	DocumentReading int `json:"documentReading"`
	TotalImages     int `json:"totalImages"`
}

// ImageInfo 图片信息
type ImageInfo struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

// Questions 题目集合
type Questions struct {
	SingleChoice    []SingleChoiceQuestion    `json:"singleChoice"`
	MultipleChoice  []MultipleChoiceQuestion  `json:"multipleChoice"`
	FillBlank       []FillBlankQuestion       `json:"fillBlank"`
	DocumentReading []DocumentReadingQuestion `json:"documentReading"`
}

// SingleChoiceQuestion 单选题
type SingleChoiceQuestion struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"` // "SC" 或 "SCIMG"
	Enabled  bool     `json:"enabled"`
	Question string   `json:"question"`
	Images   []string `json:"images"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
	Hook     string   `json:"hook,omitempty"` // 关联材料题的钩子
}

// MultipleChoiceQuestion 多选题
type MultipleChoiceQuestion struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"` // "MC" 或 "MCIMG"
	Enabled  bool     `json:"enabled"`
	Question string   `json:"question"`
	Images   []string `json:"images"`
	Options  []string `json:"options"`
	Answers  []string `json:"answers"`
	Hook     string   `json:"hook,omitempty"`
}

// FillBlankQuestion 填空题
type FillBlankQuestion struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"` // "FL" 或 "FLIMG"
	Enabled    bool              `json:"enabled"`
	Question   string            `json:"question"`
	Images     []string          `json:"images"`
	Template   string            `json:"template"`
	BlankCount int               `json:"blankCount"`
	Answers    []FillBlankAnswer `json:"answers"`
	HasExtra   bool              `json:"hasExtra"`
	ExtraKey   string            `json:"extraKey"`
	Hook       string            `json:"hook,omitempty"`
}

// FillBlankAnswer 填空题答案
type FillBlankAnswer struct {
	BlankIndex int      `json:"blankIndex"`
	Answers    []string `json:"answers"`
}

// DocumentReadingQuestion 材料阅读题
type DocumentReadingQuestion struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"` // "DR"
	Enabled   bool     `json:"enabled"`
	Question  string   `json:"question"`
	Images    []string `json:"images"`
	Materials []string `json:"materials"`
	Hooks     []string `json:"hooks"`
}

// ========== API 方法 ==========

// SaveImage 保存上传的图片，使用规范命名: {题型}_{题型序号}_{图片序号}.{扩展名}
// 例如: SCIMG_2_1.png (单选题干有图，第2题，第1张图片)
// 同时保存两份：tempwails/图片名.ext (预览用) 和 tempwails/add/图片名.ext (导出用)
func (a *App) SaveImage(filename, base64Data, questionType string, typeIndex, imageIndex int) (string, error) {
	// 创建临时图片目录
	tempwailsDir := "tempwails"
	tempAddDir := filepath.Join(tempwailsDir, "add")

	if err := os.MkdirAll(tempwailsDir, 0755); err != nil {
		return "", fmt.Errorf("创建tempwails目录失败: %v", err)
	}
	if err := os.MkdirAll(tempAddDir, 0755); err != nil {
		return "", fmt.Errorf("创建tempwails/add目录失败: %v", err)
	}

	// 获取文件扩展名
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png" // 默认扩展名
	}

	// 生成规范文件名: {题型}_{题型序号}_{图片序号}.{扩展名}
	newFilename := fmt.Sprintf("%s_%d_%d%s", questionType, typeIndex, imageIndex, ext)

	// 解码base64数据
	// base64Data格式: "data:image/png;base64,xxxxx"
	parts := strings.Split(base64Data, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("无效的base64数据格式")
	}

	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %v", err)
	}

	// 1. 保存到 tempwails/ (预览用)
	previewPath := filepath.Join(tempwailsDir, newFilename)
	if err := os.WriteFile(previewPath, imageData, 0644); err != nil {
		return "", fmt.Errorf("写入预览图片失败: %v", err)
	}

	// 2. 保存到 tempwails/add/ (导出ZIP用)
	addPath := filepath.Join(tempAddDir, newFilename)
	if err := os.WriteFile(addPath, imageData, 0644); err != nil {
		return "", fmt.Errorf("写入导出图片失败: %v", err)
	}

	// 返回相对于tempwails的路径，供前端预览使用
	return newFilename, nil
}

// ExportQuestionBank 导出题库为ZIP
// 新逻辑:
// 1. 先在 tempwails/ 生成 JSON 文件
// 2. 然后将 JSON 和 add/ 目录一起打包成 ZIP
// 3. 输出到 data/output/add.zip
func (a *App) ExportQuestionBank(bankData string) (string, error) {
	// 解析题库数据
	var bank QuestionBank
	if err := json.Unmarshal([]byte(bankData), &bank); err != nil {
		return "", fmt.Errorf("解析题库数据失败: %v", err)
	}

	// === 步骤1: 在 tempwails/ 目录生成 JSON 文件 ===
	jsonFilename := bank.Name + ".json"
	jsonPath := filepath.Join("tempwails", jsonFilename)

	jsonData, err := json.MarshalIndent(bank, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化JSON失败: %v", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("写入JSON文件失败: %v", err)
	}
	fmt.Printf("✓ 已生成JSON文件: %s\n", jsonPath)
	// === 步骤2: 创建输出目录 ===
	outputDir := filepath.Join("data", "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %v", err)
	}

	// === 步骤3: 创建 ZIP 文件 (使用题库名称命名) ===
	zipFilename := bank.Name + ".zip"
	zipPath := filepath.Join(outputDir, zipFilename)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("创建ZIP文件失败: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// === 步骤4: 添加 JSON 文件到 ZIP 根目录 ===
	jsonWriter, err := zipWriter.Create(jsonFilename)
	if err != nil {
		return "", fmt.Errorf("创建ZIP中的JSON条目失败: %v", err)
	}
	if _, err := jsonWriter.Write(jsonData); err != nil {
		return "", fmt.Errorf("写入ZIP中的JSON失败: %v", err)
	}
	fmt.Printf("✓ 已添加JSON到ZIP: %s\n", jsonFilename)

	// === 步骤5: 添加 add/ 目录到 ZIP ===
	addDir := filepath.Join("tempwails", "add")
	if _, err := os.Stat(addDir); os.IsNotExist(err) {
		fmt.Println("警告: add目录不存在，跳过图片打包")
		return zipPath, nil
	}

	// 遍历 add/ 目录的所有文件
	addedCount := 0
	err = filepath.Walk(addDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 计算ZIP中的相对路径 (保持 add/xxx.png 结构)
		relPath, err := filepath.Rel("tempwails", path)
		if err != nil {
			return fmt.Errorf("计算相对路径失败: %v", err)
		}

		// 统一使用正斜杠 (ZIP标准)
		zipPath := filepath.ToSlash(relPath)

		// 创建ZIP条目
		zipEntry, err := zipWriter.Create(zipPath)
		if err != nil {
			fmt.Printf("警告: 创建ZIP条目失败 %s: %v\n", zipPath, err)
			return nil // 继续处理其他文件
		}

		// 复制文件内容
		srcFile, err := os.Open(path)
		if err != nil {
			fmt.Printf("警告: 打开文件失败 %s: %v\n", path, err)
			return nil
		}
		defer srcFile.Close()

		if _, err := io.Copy(zipEntry, srcFile); err != nil {
			fmt.Printf("警告: 复制文件失败 %s: %v\n", path, err)
			return nil
		}

		addedCount++
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("遍历add目录失败: %v", err)
	}

	fmt.Printf("✓ 已添加 %d 个图片文件到ZIP\n", addedCount)
	fmt.Printf("✓ 导出完成: %s\n", zipPath)

	return zipPath, nil
}

// SelectSaveFolder 选择保存文件夹
func (a *App) SelectSaveFolder() (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择保存位置",
	})
	if err != nil {
		return "", err
	}
	return folder, nil
}

// PreviewQuestionBank 预览题库数据(返回格式化的JSON)
func (a *App) PreviewQuestionBank(bankData string) (string, error) {
	var bank QuestionBank
	if err := json.Unmarshal([]byte(bankData), &bank); err != nil {
		return "", fmt.Errorf("解析题库数据失败: %v", err)
	}

	// 格式化输出
	formatted, err := json.MarshalIndent(bank, "", "  ")
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

// GetImageList 获取已上传的图片列表
func (a *App) GetImageList() ([]ImageInfo, error) {
	tempDir := filepath.Join("tempwails", "add")
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return []ImageInfo{}, nil
	}

	var images []ImageInfo
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel("tempwails", path)
			images = append(images, ImageInfo{
				Filename: info.Name(),
				Path:     filepath.ToSlash(relPath),
			})
		}
		return nil
	})

	return images, err
}

// ClearTempImages 清理临时图片
// 调整为：
// 1. 如果 tempwails/ 不存在，则创建 tempwails/ 和 tempwails/add/ 并返回
// 2. 如果存在，则仅清空 tempwails/ 与 tempwails/add/ 内部的文件和子目录，但保留目录本身
// 3. 兼容旧版本 data/temp/add/：如果存在也清空其内容
func (a *App) ClearTempImages() error {
	tempDir := "tempwails"
	addDir := filepath.Join(tempDir, "add")

	// 确保目录存在
	if err := os.MkdirAll(addDir, 0755); err != nil {
		return fmt.Errorf("创建 tempwails/add 目录失败: %v", err)
	}

	// 仅清空 tempwails 根目录下内容（但保留 add 子目录，稍后单独处理）
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("读取 tempwails 目录失败: %v", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		// 保留 add 目录本身，后面单独清理其内容
		if name == "add" {
			continue
		}
		entryPath := filepath.Join(tempDir, name)
		if err := os.RemoveAll(entryPath); err != nil {
			return fmt.Errorf("删除 tempwails 下条目失败 (%s): %v", entryPath, err)
		}
	}

	// 清空 tempwails/add/ 下所有内容
	if err := clearDirectoryContents(addDir); err != nil {
		return fmt.Errorf("清空 tempwails/add 目录失败: %v", err)
	}

	// 兼容旧版本 data/temp/add/：如存在则清空其内容
	legacyAddDir := filepath.Join("data", "temp", "add")
	if _, err := os.Stat(legacyAddDir); err == nil {
		if err := clearDirectoryContents(legacyAddDir); err != nil {
			return fmt.Errorf("清空 data/temp/add 目录失败: %v", err)
		}
	}

	return nil
}

// ValidateQuestionBank 验证题库数据
func (a *App) ValidateQuestionBank(bankData string) ([]string, error) {
	var bank QuestionBank
	if err := json.Unmarshal([]byte(bankData), &bank); err != nil {
		return nil, fmt.Errorf("解析题库数据失败: %v", err)
	}

	var errors []string

	// 验证题库名称
	if bank.Name == "" {
		errors = append(errors, "题库名称不能为空")
	}

	// 验证题目数量
	if len(bank.Questions.SingleChoice) == 0 &&
		len(bank.Questions.MultipleChoice) == 0 &&
		len(bank.Questions.FillBlank) == 0 &&
		len(bank.Questions.DocumentReading) == 0 {
		errors = append(errors, "题库至少需要包含一道题目")
	}

	// 验证材料题的hooks
	for _, dr := range bank.Questions.DocumentReading {
		for _, hook := range dr.Hooks {
			found := false
			// 检查是否存在对应的子题
			if strings.HasPrefix(hook, "SC.") || strings.HasPrefix(hook, "SCIMG.") {
				for _, q := range bank.Questions.SingleChoice {
					if q.Hook == hook {
						found = true
						break
					}
				}
			} else if strings.HasPrefix(hook, "MC.") || strings.HasPrefix(hook, "MCIMG.") {
				for _, q := range bank.Questions.MultipleChoice {
					if q.Hook == hook {
						found = true
						break
					}
				}
			} else if strings.HasPrefix(hook, "FL.") || strings.HasPrefix(hook, "FLIMG.") {
				for _, q := range bank.Questions.FillBlank {
					if q.Hook == hook {
						found = true
						break
					}
				}
			}
			if !found {
				errors = append(errors, fmt.Sprintf("材料题关联的子题不存在: %s", hook))
			}
		}
	}

	return errors, nil
}

// ========== 图片删除功能（通过 JSON 桥接） ==========

// DeleteImageRequest 删除图片请求结构
type DeleteImageRequest struct {
	Timestamp string   `json:"timestamp"`
	Images    []string `json:"images"`
}

// ProcessPendingDeleteImages 处理待删除的图片（读取 localStorage 数据）
func (a *App) ProcessPendingDeleteImages(jsonData string) error {
	if jsonData == "" {
		return fmt.Errorf("没有待删除的图片")
	}

	var req DeleteImageRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		return fmt.Errorf("解析删除请求失败: %v", err)
	}

	if len(req.Images) == 0 {
		return fmt.Errorf("删除列表为空")
	}

	deletedCount := 0
	var errors []string

	for _, imgPath := range req.Images {
		// imgPath 格式: "add/SC_1_1.png"
		// 需要删除两个位置的文件:
		// 1. tempwails/SC_1_1.png (预览用)
		// 2. tempwails/add/SC_1_1.png (导出用)

		// 提取文件名
		filename := filepath.Base(imgPath)

		// 删除预览文件
		previewPath := filepath.Join("tempwails", filename)
		if err := os.Remove(previewPath); err != nil {
			if !os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("删除预览文件失败 %s: %v", previewPath, err))
			}
		} else {
			fmt.Printf("✓ 已删除预览文件: %s\n", previewPath)
		}

		// 删除导出文件
		exportPath := filepath.Join("tempwails", imgPath)
		if err := os.Remove(exportPath); err != nil {
			if !os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("删除导出文件失败 %s: %v", exportPath, err))
			}
		} else {
			fmt.Printf("✓ 已删除导出文件: %s\n", exportPath)
			deletedCount++
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("删除完成，但有错误: %s", strings.Join(errors, "; "))
	}

	fmt.Printf("✓ 成功删除 %d 个图片文件\n", deletedCount)
	return nil
}

// DeleteSingleImage 删除单张图片
func (a *App) DeleteSingleImage(imagePath string) error {
	if imagePath == "" {
		return fmt.Errorf("图片路径为空")
	}

	// 提取文件名
	filename := filepath.Base(imagePath)

	// 删除预览文件
	previewPath := filepath.Join("tempwails", filename)
	if err := os.Remove(previewPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除预览文件失败: %v", err)
	}

	// 删除导出文件
	exportPath := filepath.Join("tempwails", imagePath)
	if err := os.Remove(exportPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除导出文件失败: %v", err)
	}

	fmt.Printf("✓ 已删除图片: %s\n", imagePath)
	return nil
}

// ImportQuestionBank 从ZIP导入题库：
// 1. 让用户选择ZIP
// 2. 解压到 wails/fix 目录
// 3. 将其中 add/ 下图片复制到 wails/tempwails/add
// 4. 读取JSON文件内容并原样返回给前端
func (a *App) ImportQuestionBank() (string, error) {
	// 1. 选择 ZIP 文件
	dialogOpts := runtime.OpenDialogOptions{
		Title: "选择题库ZIP文件",
		Filters: []runtime.FileFilter{{
			DisplayName: "题库ZIP",
			Pattern:     "*.zip",
		}},
	}
	zipPath, err := runtime.OpenFileDialog(a.ctx, dialogOpts)
	if err != nil {
		return "", fmt.Errorf("打开文件选择对话框失败: %v", err)
	}
	if zipPath == "" {
		// 用户取消
		return "", nil
	}

	fmt.Println("选择的题库ZIP:", zipPath)

	// 2. 解压到 wails/fix 目录
	fixRoot := "fix"
	if err := os.MkdirAll(fixRoot, 0755); err != nil {
		return "", fmt.Errorf("创建 fix 目录失败: %v", err)
	}

	// 为本次导入创建子目录，避免相互覆盖
	baseName := strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath))
	targetDir := filepath.Join(fixRoot, baseName)

	// 先清空同名目录
	_ = os.RemoveAll(targetDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("创建 fix 子目录失败: %v", err)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("打开ZIP失败: %v", err)
	}
	defer r.Close()

	var jsonFilePath string

	for _, f := range r.File {
		relPath := filepath.ToSlash(f.Name)

		// 记录根目录 JSON 文件路径
		if !f.FileInfo().IsDir() && strings.HasSuffix(strings.ToLower(relPath), ".json") && !strings.Contains(relPath, "/") {
			jsonFilePath = filepath.Join(targetDir, filepath.Base(relPath))
		}

		outPath := filepath.Join(targetDir, relPath)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, 0755); err != nil {
				return "", fmt.Errorf("创建解压目录失败: %v", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return "", fmt.Errorf("创建解压父目录失败: %v", err)
		}

		rc, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("打开ZIP内部文件失败: %v", err)
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			return "", fmt.Errorf("创建解压文件失败: %v", err)
		}

		if _, err := io.Copy(outFile, rc); err != nil {
			rc.Close()
			outFile.Close()
			return "", fmt.Errorf("写入解压文件失败: %v", err)
		}
		rc.Close()
		outFile.Close()
	}

	if jsonFilePath == "" {
		return "", fmt.Errorf("在ZIP中未找到题库JSON文件")
	}

	// 3. 将 add/ 下图片复制到 tempwails/add，同时确保存在 tempwails 根目录（预览用）
	addSrc := filepath.Join(targetDir, "add")
	tempRoot := "tempwails"
	addDstRoot := filepath.Join(tempRoot, "add")
	if err := os.MkdirAll(addDstRoot, 0755); err != nil {
		return "", fmt.Errorf("创建 tempwails/add 目录失败: %v", err)
	}

	filepath.Walk(addSrc, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(addSrc, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(addDstRoot, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			dstFile.Close()
			return err
		}
		dstFile.Close()
		return nil
	})

	// 4. 读取JSON内容并返回
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("读取题库JSON失败: %v", err)
	}

	return string(data), nil
}
