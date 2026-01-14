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
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 处理图片ZIP文件
func ProcessImageZip(zipPath string) (int, error) {
	// 检查data\temp路径下是否有add文件夹，有则删除
	tempAddPath := filepath.Join("data", "temp", "add")
	if _, err := os.Stat(tempAddPath); err == nil {
		os.RemoveAll(tempAddPath)
	}

	// 解压ZIP文件
	extractPath := filepath.Join("data", "temp")
	err := extractZip(zipPath, extractPath)
	if err != nil {
		return 0, fmt.Errorf("解压ZIP失败: %v", err)
	}

	// 检查是否有add文件夹
	addFolderPath := filepath.Join(extractPath, "add")
	if _, err := os.Stat(addFolderPath); os.IsNotExist(err) {
		// 清理并返回错误
		os.RemoveAll(addFolderPath)
		return 0, fmt.Errorf("ZIP文件中必须包含add文件夹")
	}

	// 检查图片文件名冲突
	tempImagesPath := filepath.Join("data", "temp", "tempImages")
	err = checkImageNameConflicts(addFolderPath, tempImagesPath)
	if err != nil {
		// 清理add文件夹并返回错误
		os.RemoveAll(addFolderPath)
		return 0, err
	}

	// 复制图片到tempImages
	imageCount, err := copyImagesToTemp(addFolderPath, tempImagesPath)
	if err != nil {
		os.RemoveAll(addFolderPath)
		return 0, fmt.Errorf("复制图片失败: %v", err)
	}

	return imageCount, nil
}

// 导入题库ZIP文件（包含所有类型的xlsx和add文件夹）
func ImportQuestionBankZip(zipPath string) (*ParseResult, []string, []string, error) {
	// 临时解压目录 - 使用持久的目录
	tempExtractPath := filepath.Join("data", "temp", "imported_files")

	// 确保目录存在，先清空旧文件（但保留到用户保存题库后）
	os.RemoveAll(tempExtractPath)
	if err := os.MkdirAll(tempExtractPath, 0755); err != nil {
		return nil, nil, nil, fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 解压ZIP
	err := extractZip(zipPath, tempExtractPath)
	if err != nil {
		os.RemoveAll(tempExtractPath)
		return nil, nil, nil, fmt.Errorf("解压题库ZIP失败: %v", err)
	}

	// 查找所有xlsx文件并分类
	mainFiles, extraFiles, addFolderPath, err := classifyFilesInZip(tempExtractPath)
	if err != nil {
		os.RemoveAll(tempExtractPath)
		return nil, nil, nil, err
	}

	if len(mainFiles) == 0 {
		os.RemoveAll(tempExtractPath)
		return nil, nil, nil, fmt.Errorf("ZIP文件中未找到主文件(SC/MC/FL/SCIMG/MCIMG/FLIMG/DR类型的xlsx文件)")
	}

	// 处理图片
	if addFolderPath != "" {
		// 检查图片文件名冲突
		tempImagesPath := filepath.Join("data", "temp", "tempImages")
		err = checkImageNameConflicts(addFolderPath, tempImagesPath)
		if err != nil {
			os.RemoveAll(tempExtractPath)
			return nil, nil, nil, err
		}

		// 复制图片到tempImages
		_, err = copyImagesToTemp(addFolderPath, tempImagesPath)
		if err != nil {
			os.RemoveAll(tempExtractPath)
			return nil, nil, nil, fmt.Errorf("复制图片失败: %v", err)
		}
	}

	// 解析xlsx文件
	result, err := ParseQuestionFiles(mainFiles, extraFiles)
	if err != nil {
		os.RemoveAll(tempExtractPath)
		return nil, nil, nil, fmt.Errorf("解析xlsx文件失败: %v", err)
	}

	// 重要：不要在这里删除临时目录！保留到用户保存题库后
	// os.RemoveAll(tempExtractPath) // 注释掉这行

	return result, mainFiles, extraFiles, nil
}

// 在解压目录中分类文件
func classifyFilesInZip(rootPath string) ([]string, []string, string, error) {
	var mainFiles []string  // SC, MC, FL, SCIMG, MCIMG, FLIMG, DR
	var extraFiles []string // FE
	var addFolderPath string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 查找add文件夹
		if info.IsDir() && info.Name() == "add" {
			addFolderPath = path
			return nil
		}

		// 只处理xlsx文件
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".xlsx" {
			// 识别文件类型
			fileType, _, err := GetQuestionType(path)
			if err != nil {
				// 如果不是标准格式的文件，跳过
				return nil
			}

			switch fileType {
			case SingleChoice, MultipleChoice, FillBlank,
				SingleChoiceImg, MultipleChoiceImg, FillBlankImg, DocumentReading:
				mainFiles = append(mainFiles, path)
			case FillExtra:
				extraFiles = append(extraFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, "", err
	}

	return mainFiles, extraFiles, addFolderPath, nil
}

// 解压ZIP文件
func extractZip(zipPath, destPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 确保目标目录存在
	os.MkdirAll(destPath, 0755)

	for _, file := range reader.File {
		filePath := filepath.Join(destPath, file.Name)

		// 防止路径遍历攻击
		if !strings.HasPrefix(filePath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("无效的文件路径: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		// 创建文件目录
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		// 创建目标文件
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		// 打开源文件
		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		// 复制文件内容
		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// 检查图片文件名冲突
func checkImageNameConflicts(sourcePath, targetPath string) error {
	// 确保目标目录存在
	os.MkdirAll(targetPath, 0755)

	// 获取目标目录中现有的所有文件名（不含扩展名）
	existingFiles := make(map[string]bool)

	targetFiles, err := os.ReadDir(targetPath)
	if err == nil {
		for _, file := range targetFiles {
			if !file.IsDir() {
				nameWithoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				existingFiles[nameWithoutExt] = true
			}
		}
	}

	// 检查源目录中的文件
	sourceFiles, err := os.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	for _, file := range sourceFiles {
		if !file.IsDir() && isImageFile(file.Name()) {
			nameWithoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			if existingFiles[nameWithoutExt] {
				return fmt.Errorf("图片文件名冲突: %s (已存在同名文件)", nameWithoutExt)
			}
		}
	}

	return nil
}

// 复制图片到tempImages目录
func copyImagesToTemp(sourcePath, targetPath string) (int, error) {
	// 确保目标目录存在
	os.MkdirAll(targetPath, 0755)

	sourceFiles, err := os.ReadDir(sourcePath)
	if err != nil {
		return 0, err
	}

	imageCount := 0
	for _, file := range sourceFiles {
		if !file.IsDir() && isImageFile(file.Name()) {
			srcPath := filepath.Join(sourcePath, file.Name())
			dstPath := filepath.Join(targetPath, file.Name())

			// 复制文件
			err := copyFile(srcPath, dstPath)
			if err != nil {
				return imageCount, err
			}
			imageCount++
		}
	}

	return imageCount, nil
}

// 判断是否为图片文件
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExtensions := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".tiff", ".tif", ".webp"}

	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// 获取当前tempImages中的图片数量
func GetTempImageCount() int {
	tempImagesPath := filepath.Join("data", "temp", "tempImages")

	if _, err := os.Stat(tempImagesPath); os.IsNotExist(err) {
		return 0
	}

	files, err := os.ReadDir(tempImagesPath)
	if err != nil {
		return 0
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() && isImageFile(file.Name()) {
			count++
		}
	}
	return count
}

// 为显示目的分类ZIP中的文件（不实际解压）
func ClassifyFilesInZipForDisplay(zipPath string) ([]string, []string, string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, nil, "", err
	}
	defer reader.Close()

	var mainFiles []string
	var extraFiles []string
	var hasAddFolder bool

	for _, file := range reader.File {
		// 检查是否有add文件夹
		if file.FileInfo().IsDir() && (file.Name == "add/" || file.Name == "add") {
			hasAddFolder = true
			continue
		}

		// 只处理xlsx文件
		if !file.FileInfo().IsDir() && strings.ToLower(filepath.Ext(file.Name)) == ".xlsx" {
			// 从文件名识别类型
			baseName := filepath.Base(file.Name)
			nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
			parts := strings.Split(nameWithoutExt, "_")

			if len(parts) < 2 {
				continue
			}

			fileType := parts[len(parts)-1]

			switch fileType {
			case "SC", "MC", "FL", "SCIMG", "MCIMG", "FLIMG", "DR":
				mainFiles = append(mainFiles, baseName)
			case "FE":
				extraFiles = append(extraFiles, baseName)
			}
		}
	}

	addFolderPath := ""
	if hasAddFolder {
		addFolderPath = "add/"
	}

	return mainFiles, extraFiles, addFolderPath, nil
}
