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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

// 题目类型
type QuestionType string

const (
	SingleChoice      QuestionType = "SC"
	MultipleChoice    QuestionType = "MC"
	FillBlank         QuestionType = "FL"
	FillExtra         QuestionType = "FE"
	SingleChoiceImg   QuestionType = "SCIMG"
	MultipleChoiceImg QuestionType = "MCIMG"
	FillBlankImg      QuestionType = "FLIMG"
	DocumentReading   QuestionType = "DR"
)

// 基础题目结构
type BaseQuestion struct {
	ID       string       `json:"id"`
	Type     QuestionType `json:"type"`
	Enabled  bool         `json:"enabled"`
	Question string       `json:"question"`
	Images   []string     `json:"images"`
	Hook     string       `json:"hook,omitempty"` // 钩子字段
}

// 单选题
type SingleChoiceQuestion struct {
	BaseQuestion
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

// 多选题
type MultipleChoiceQuestion struct {
	BaseQuestion
	Options []string `json:"options"`
	Answers []string `json:"answers"`
}

// 填空题答案项
type FillBlankAnswer struct {
	BlankIndex int      `json:"blankIndex"` // 第几个空
	Answers    []string `json:"answers"`    // 可接受的答案
}

// 填空题
type FillBlankQuestion struct {
	BaseQuestion
	Template   string            `json:"template"`   // 题干模板
	BlankCount int               `json:"blankCount"` // 空的数量
	Answers    []FillBlankAnswer `json:"answers"`    // 答案
	HasExtra   bool              `json:"hasExtra"`   // 是否有额外答案
	ExtraKey   string            `json:"extraKey"`   // 额外答案标识 (x%x)
}

// 材料阅读题
type DocumentReadingQuestion struct {
	BaseQuestion
	Materials []string `json:"materials"` // 资料内容
	Hooks     []string `json:"hooks"`     // 关联的钩子列表
}

// 文件解析结果
type ParseResult struct {
	SingleChoice    []SingleChoiceQuestion    `json:"singleChoice"`
	MultipleChoice  []MultipleChoiceQuestion  `json:"multipleChoice"`
	FillBlank       []FillBlankQuestion       `json:"fillBlank"`
	DocumentReading []DocumentReadingQuestion `json:"documentReading"`
	Errors          []string                  `json:"errors"`
}

// 从文件名识别题目类型
func GetQuestionType(filename string) (QuestionType, string, error) {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	parts := strings.Split(nameWithoutExt, "_")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("文件名格式错误: %s", filename)
	}

	fileType := parts[len(parts)-1]
	prefix := strings.Join(parts[:len(parts)-1], "_")

	// 修复：使用 prefix 变量，避免 "declared and not used"
	_ = prefix // 这行先加上，后面会用到前缀匹配

	switch fileType {
	case "SC":
		return SingleChoice, prefix, nil
	case "MC":
		return MultipleChoice, prefix, nil
	case "FL":
		return FillBlank, prefix, nil
	case "FE":
		return FillExtra, prefix, nil
	case "SCIMG":
		return SingleChoiceImg, prefix, nil
	case "MCIMG":
		return MultipleChoiceImg, prefix, nil
	case "FLIMG":
		return FillBlankImg, prefix, nil
	case "DR":
		return DocumentReading, prefix, nil
	default:
		return "", "", fmt.Errorf("未知的文件类型: %s", fileType)
	}
}

// 解析单选题文件
func parseSingleChoiceFile(filePath string) ([]SingleChoiceQuestion, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return nil, fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]
	var questions []SingleChoiceQuestion

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		// 检查是否启用 (A列)
		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取选项 (B-K列，索引1-10)
		var options []string
		for j := 1; j <= 10; j++ {
			option := strings.TrimSpace(row.GetCell(j).String())
			if option != "" && option != "nan" && option != "None" {
				options = append(options, option)
			}
		}

		// 读取题干和答案
		question := strings.TrimSpace(row.GetCell(12).String()) // M列
		answer := strings.TrimSpace(row.GetCell(13).String())   // N列

		if question == "" {
			continue
		}

		// 读取图片名称 (O-Y列，索引14-24)
		var images []string
		for j := 14; j <= 24; j++ {
			imageName := strings.TrimSpace(row.GetCell(j).String())
			if imageName != "" {
				// 修复：构建正确的图片路径
				imagePath := buildImagePath(imageName)
				if imagePath != "" {
					images = append(images, imagePath)
				}
			}
		}

		// 读取钩子 (Z列，索引25)
		hook := strings.TrimSpace(row.GetCell(25).String())

		questions = append(questions, SingleChoiceQuestion{
			BaseQuestion: BaseQuestion{
				ID:       strconv.Itoa(i + 1),
				Type:     SingleChoice,
				Enabled:  true,
				Question: question,
				Images:   images,
				Hook:     hook,
			},
			Options: options,
			Answer:  answer,
		})
	}

	return questions, nil
}

// 解析题干是图片的单选题文件
func parseSingleChoiceImgFile(filePath string) ([]SingleChoiceQuestion, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return nil, fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]
	var questions []SingleChoiceQuestion

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		// 检查是否启用 (A列)
		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取选项 (B-K列，索引1-10)
		var options []string
		for j := 1; j <= 10; j++ {
			option := strings.TrimSpace(row.GetCell(j).String())
			if option != "" && option != "nan" && option != "None" {
				options = append(options, option)
			}
		}

		// 读取题干图片和文字描述
		questionImages := []string{}
		// 读取题干图片 (M-P列，索引12-15)
		for j := 12; j <= 15; j++ {
			imageName := strings.TrimSpace(row.GetCell(j).String())
			if imageName != "" {
				imagePath := buildImagePath(imageName)
				if imagePath != "" {
					questionImages = append(questionImages, imagePath)
				}
			}
		}

		// 读取文字描述 (Q列，索引16)
		description := strings.TrimSpace(row.GetCell(16).String())
		answer := strings.TrimSpace(row.GetCell(17).String()) // R列

		// 读取钩子 (S列，索引18)
		hook := strings.TrimSpace(row.GetCell(18).String())

		if description == "" && len(questionImages) == 0 {
			continue
		}

		questions = append(questions, SingleChoiceQuestion{
			BaseQuestion: BaseQuestion{
				ID:       strconv.Itoa(i + 1),
				Type:     SingleChoiceImg,
				Enabled:  true,
				Question: description,
				Images:   questionImages,
				Hook:     hook,
			},
			Options: options,
			Answer:  answer,
		})
	}

	return questions, nil
}

// 解析多选题文件
func parseMultipleChoiceFile(filePath string) ([]MultipleChoiceQuestion, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return nil, fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]
	var questions []MultipleChoiceQuestion

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取选项 (B-L列，索引1-11)
		var options []string
		for j := 1; j <= 11; j++ {
			option := strings.TrimSpace(row.GetCell(j).String())
			if option != "" {
				options = append(options, option)
			}
		}

		// 读取题干和答案
		question := strings.TrimSpace(row.GetCell(12).String()) // M列
		answer := strings.TrimSpace(row.GetCell(13).String())   // N列

		if question == "" {
			continue
		}

		// 读取图片名称 (O-Y列，索引14-24)
		var images []string
		for j := 14; j <= 24; j++ {
			imageName := strings.TrimSpace(row.GetCell(j).String())
			if imageName != "" {
				// 修复：构建正确的图片路径
				imagePath := buildImagePath(imageName)
				if imagePath != "" {
					images = append(images, imagePath)
				}
			}
		}

		// 读取钩子 (Z列，索引25)
		hook := strings.TrimSpace(row.GetCell(25).String())

		// 将答案字符串转换为答案数组
		var answers []string
		for _, char := range answer {
			answers = append(answers, string(char))
		}

		questions = append(questions, MultipleChoiceQuestion{
			BaseQuestion: BaseQuestion{
				ID:       strconv.Itoa(i),
				Type:     MultipleChoice,
				Enabled:  true,
				Question: question,
				Images:   images,
				Hook:     hook,
			},
			Options: options,
			Answers: answers,
		})
	}

	return questions, nil
}

// 解析题干是图片的多选题文件
func parseMultipleChoiceImgFile(filePath string) ([]MultipleChoiceQuestion, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return nil, fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]
	var questions []MultipleChoiceQuestion

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取选项 (B-L列，索引1-11)
		var options []string
		for j := 1; j <= 11; j++ {
			option := strings.TrimSpace(row.GetCell(j).String())
			if option != "" {
				options = append(options, option)
			}
		}

		// 读取题干图片和文字描述
		questionImages := []string{}
		// 读取题干图片 (M-P列，索引12-15)
		for j := 12; j <= 15; j++ {
			imageName := strings.TrimSpace(row.GetCell(j).String())
			if imageName != "" {
				imagePath := buildImagePath(imageName)
				if imagePath != "" {
					questionImages = append(questionImages, imagePath)
				}
			}
		}

		// 读取文字描述 (Q列，索引16)
		description := strings.TrimSpace(row.GetCell(16).String())
		answer := strings.TrimSpace(row.GetCell(17).String()) // R列

		// 读取钩子 (S列，索引18)
		hook := strings.TrimSpace(row.GetCell(18).String())

		if description == "" && len(questionImages) == 0 {
			continue
		}

		// 将答案字符串转换为答案数组
		var answers []string
		for _, char := range answer {
			answers = append(answers, string(char))
		}

		questions = append(questions, MultipleChoiceQuestion{
			BaseQuestion: BaseQuestion{
				ID:       strconv.Itoa(i),
				Type:     MultipleChoiceImg,
				Enabled:  true,
				Question: description,
				Images:   questionImages,
				Hook:     hook,
			},
			Options: options,
			Answers: answers,
		})
	}

	return questions, nil
}

// 解析填空题文件
func parseFillBlankFile(filePath string) ([]FillBlankQuestion, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return nil, fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]
	var questions []FillBlankQuestion

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取选项 (B-U列)
		var options []string
		var extraKey string
		hasExtra := false

		for j := 1; j <= 20; j++ {
			option := strings.TrimSpace(row.GetCell(j).String())
			if option != "" {
				// 检查是否是 (x%x) 格式
				if strings.HasPrefix(option, "(") && strings.Contains(option, "%") && strings.HasSuffix(option, ")") {
					hasExtra = true
					extraKey = option
				}
				options = append(options, option)
			}
		}

		// 读取题干和序号
		question := strings.TrimSpace(row.GetCell(21).String()) // V列
		id := strings.TrimSpace(row.GetCell(22).String())       // W列

		if question == "" {
			continue
		}

		// 读取图片名称 (X-AH列，索引23-33)
		var images []string
		for j := 23; j <= 33; j++ {
			imageName := strings.TrimSpace(row.GetCell(j).String())
			if imageName != "" {
				// 修复：构建正确的图片路径
				imagePath := buildImagePath(imageName)
				if imagePath != "" {
					images = append(images, imagePath)
				}
			}
		}

		// 读取钩子 (AI列，索引34)
		hook := strings.TrimSpace(row.GetCell(34).String())

		// 计算空的数量
		blankCount := strings.Count(question, "(%___%)")

		// 创建基础答案
		var answers []FillBlankAnswer
		for i := 1; i <= blankCount; i++ {
			if i-1 < len(options) {
				answers = append(answers, FillBlankAnswer{
					BlankIndex: i,
					Answers:    []string{options[i-1]},
				})
			}
		}

		questions = append(questions, FillBlankQuestion{
			BaseQuestion: BaseQuestion{
				ID:       id,
				Type:     FillBlank,
				Enabled:  true,
				Question: question,
				Images:   images,
				Hook:     hook,
			},
			Template:   question,
			BlankCount: blankCount,
			Answers:    answers,
			HasExtra:   hasExtra,
			ExtraKey:   extraKey,
		})
	}

	return questions, nil
}

// 解析题干是图片的填空题文件
func parseFillBlankImgFile(filePath string) ([]FillBlankQuestion, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return nil, fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]
	var questions []FillBlankQuestion

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取选项 (B-U列)
		var options []string
		var extraKey string
		hasExtra := false

		for j := 1; j <= 20; j++ {
			option := strings.TrimSpace(row.GetCell(j).String())
			if option != "" {
				// 检查是否是 (x%x) 格式
				if strings.HasPrefix(option, "(") && strings.Contains(option, "%") && strings.HasSuffix(option, ")") {
					hasExtra = true
					extraKey = option
				}
				options = append(options, option)
			}
		}

		// 读取题干图片和文字描述
		questionImages := []string{}
		// 读取题干图片 (V-Y列，索引21-24)
		for j := 21; j <= 24; j++ {
			imageName := strings.TrimSpace(row.GetCell(j).String())
			if imageName != "" {
				imagePath := buildImagePath(imageName)
				if imagePath != "" {
					questionImages = append(questionImages, imagePath)
				}
			}
		}

		// 读取文字描述 (Z列，索引25)
		description := strings.TrimSpace(row.GetCell(25).String())
		id := strconv.Itoa(i) // 使用行号作为ID

		// 读取钩子 (AA列，索引26)
		hook := strings.TrimSpace(row.GetCell(26).String())

		if description == "" && len(questionImages) == 0 {
			continue
		}

		// 计算空的数量
		blankCount := strings.Count(description, "(%___%)")

		// 创建基础答案
		var answers []FillBlankAnswer
		for i := 1; i <= blankCount; i++ {
			if i-1 < len(options) {
				answers = append(answers, FillBlankAnswer{
					BlankIndex: i,
					Answers:    []string{options[i-1]},
				})
			}
		}

		questions = append(questions, FillBlankQuestion{
			BaseQuestion: BaseQuestion{
				ID:       id,
				Type:     FillBlankImg,
				Enabled:  true,
				Question: description,
				Images:   questionImages,
				Hook:     hook,
			},
			Template:   description,
			BlankCount: blankCount,
			Answers:    answers,
			HasExtra:   hasExtra,
			ExtraKey:   extraKey,
		})
	}

	return questions, nil
}

// 解析材料阅读题文件
func parseDocumentReadingFile(filePath string) ([]DocumentReadingQuestion, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return nil, fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]
	var questions []DocumentReadingQuestion

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取资料 (B-E列，索引1-4)
		var materials []string
		for j := 1; j <= 4; j++ {
			material := strings.TrimSpace(row.GetCell(j).String())
			if material != "" {
				materials = append(materials, material)
			}
		}

		// 读取题干图片 (F-I列，索引5-8)
		var images []string
		for j := 5; j <= 8; j++ {
			imageName := strings.TrimSpace(row.GetCell(j).String())
			if imageName != "" {
				imagePath := buildImagePath(imageName)
				if imagePath != "" {
					images = append(images, imagePath)
				}
			}
		}

		// 读取钩子列表 (J-AM列，索引9-38)
		var hooks []string
		for j := 9; j <= 38; j++ {
			hook := strings.TrimSpace(row.GetCell(j).String())
			if hook != "" {
				hooks = append(hooks, hook)
			}
		}

		if len(materials) == 0 {
			continue
		}

		questions = append(questions, DocumentReadingQuestion{
			BaseQuestion: BaseQuestion{
				ID:       strconv.Itoa(i),
				Type:     DocumentReading,
				Enabled:  true,
				Question: "材料阅读题", // 材料题没有具体题干，用固定描述
				Images:   images,
				Hook:     "", // 材料题本身没有钩子，钩子在hooks字段中
			},
			Materials: materials,
			Hooks:     hooks,
		})
	}

	return questions, nil
}

// addImagePathsToSingleChoice
// 解析填空补充文件
func parseFillExtraFile(filePath string, fillQuestions *[]FillBlankQuestion) error {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}

	if len(file.Sheets) == 0 {
		return fmt.Errorf("文件没有工作表")
	}

	sheet := file.Sheets[0]

	for i := 1; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			continue
		}

		enabled := strings.ToUpper(row.GetCell(0).String()) == "TURE"
		if !enabled {
			continue
		}

		// 读取序号、第几空和答案
		id := strings.TrimSpace(row.GetCell(1).String())            // B列
		blankIndexStr := strings.TrimSpace(row.GetCell(2).String()) // C列
		blankIndex, err := strconv.Atoi(blankIndexStr)
		if err != nil {
			continue
		}

		// 读取答案 (D-M列)
		var extraAnswers []string
		for j := 3; j <= 12; j++ {
			answer := strings.TrimSpace(row.GetCell(j).String())
			if answer != "" {
				extraAnswers = append(extraAnswers, answer)
			}
		}

		// 找到对应的填空题并附加答案
		for i, q := range *fillQuestions {
			if q.ID == id {
				// 找到对应的空
				for j, answer := range q.Answers {
					if answer.BlankIndex == blankIndex {
						// 附加FE文件的答案到FL文件的答案后面，而不是替换
						combinedAnswers := append(answer.Answers, extraAnswers...)
						(*fillQuestions)[i].Answers[j].Answers = combinedAnswers
						break
					}
				}
				break
			}
		}
	}

	return nil
}

// 主解析函数
func ParseQuestionFiles(mainFiles []string, extraFiles []string) (*ParseResult, error) {
	result := &ParseResult{
		SingleChoice:    []SingleChoiceQuestion{},
		MultipleChoice:  []MultipleChoiceQuestion{},
		FillBlank:       []FillBlankQuestion{},
		DocumentReading: []DocumentReadingQuestion{},
		Errors:          []string{},
	}

	// 先解析主文件
	for _, filePath := range mainFiles {
		fileType, prefix, err := GetQuestionType(filePath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("文件 %s: %v", filePath, err))
			continue
		}

		switch fileType {
		case SingleChoice:
			questions, err := parseSingleChoiceFile(filePath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("解析单选题文件 %s: %v", filePath, err))
			} else {
				result.SingleChoice = append(result.SingleChoice, questions...)
			}

		case SingleChoiceImg:
			questions, err := parseSingleChoiceImgFile(filePath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("解析题干是图单选题文件 %s: %v", filePath, err))
			} else {
				result.SingleChoice = append(result.SingleChoice, questions...)
			}

		case MultipleChoice:
			questions, err := parseMultipleChoiceFile(filePath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("解析多选题文件 %s: %v", filePath, err))
			} else {
				result.MultipleChoice = append(result.MultipleChoice, questions...)
			}

		case MultipleChoiceImg:
			questions, err := parseMultipleChoiceImgFile(filePath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("解析题干是图多选题文件 %s: %v", filePath, err))
			} else {
				result.MultipleChoice = append(result.MultipleChoice, questions...)
			}

		case FillBlank:
			questions, err := parseFillBlankFile(filePath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("解析填空题文件 %s: %v", filePath, err))
			} else {
				result.FillBlank = append(result.FillBlank, questions...)

				// 为每个填空题文件寻找对应的补充文件
				matchedExtraFile := findMatchingExtraFile(prefix, extraFiles)
				if matchedExtraFile != "" {
					err := parseFillExtraFile(matchedExtraFile, &result.FillBlank)
					if err != nil {
						result.Errors = append(result.Errors, fmt.Sprintf("解析填空补充文件 %s: %v", matchedExtraFile, err))
					}
				}
			}

		case FillBlankImg:
			questions, err := parseFillBlankImgFile(filePath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("解析题干是图填空题文件 %s: %v", filePath, err))
			} else {
				result.FillBlank = append(result.FillBlank, questions...)
				// FLIMG类型没有对应的FE文件，所以不处理补充文件
			}

		case DocumentReading:
			questions, err := parseDocumentReadingFile(filePath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("解析材料阅读题文件 %s: %v", filePath, err))
			} else {
				result.DocumentReading = append(result.DocumentReading, questions...)
			}
		}
	}

	return result, nil
}

// 根据前缀找到匹配的补充文件
func findMatchingExtraFile(prefix string, extraFiles []string) string {
	for _, filePath := range extraFiles {
		fileType, filePrefix, err := GetQuestionType(filePath)
		if err == nil && fileType == FillExtra && filePrefix == prefix {
			return filePath
		}
	}
	return ""
}

// 保存题库到JSON文件
func SaveQuestionBank(result *ParseResult, bankName string, filePath string) error {
	// 获取图片信息
	imageFiles := getImageFilesInfo()

	questionBank := map[string]interface{}{
		"name":    bankName,
		"version": "1.0",
		"metadata": map[string]interface{}{
			"totalQuestions":  len(result.SingleChoice) + len(result.MultipleChoice) + len(result.FillBlank) + len(result.DocumentReading),
			"singleChoice":    len(result.SingleChoice),
			"multipleChoice":  len(result.MultipleChoice),
			"fillBlank":       len(result.FillBlank),
			"documentReading": len(result.DocumentReading),
			"totalImages":     len(imageFiles),
		},
		"questions": map[string]interface{}{
			"singleChoice":    result.SingleChoice,    // 直接使用，图片路径已在解析时设置
			"multipleChoice":  result.MultipleChoice,  // 直接使用
			"fillBlank":       result.FillBlank,       // 直接使用
			"documentReading": result.DocumentReading, // 材料阅读题
		},
		"images": imageFiles, // 所有图片文件列表
		"errors": result.Errors,
	}

	jsonData, err := json.MarshalIndent(questionBank, "", "  ")
	if err != nil {
		return fmt.Errorf("生成JSON失败: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// 生成带时间戳的文件名
func GenerateFileName(bankName string) string {
	// 获取当前时间
	now := time.Now()
	// 格式化为: 20060102_150405 (年月日_时分秒)
	timestamp := now.Format("20060102_150405")

	if bankName == "" {
		bankName = "未命名题库"
	}

	return fmt.Sprintf("%s_%s.json", bankName, timestamp)
}

// 在文件末尾添加以下函数：

// 保存题库并打包ZIP
func SaveQuestionBankWithImages(result *ParseResult, bankName string, filePath string) error {
	// 创建输出临时目录
	outputTempDir := filepath.Join("data", "output", "outputTemp")
	if err := os.MkdirAll(outputTempDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 生成JSON文件路径
	jsonFileName := bankName + ".json"
	if bankName == "" {
		jsonFileName = "未命名题库.json"
	}
	jsonFilePath := filepath.Join(outputTempDir, jsonFileName)

	// 保存JSON文件（使用修改后的函数）
	err := SaveQuestionBank(result, bankName, jsonFilePath)
	if err != nil {
		return err
	}

	// 创建add文件夹
	addFolderPath := filepath.Join(outputTempDir, "add")
	if err := os.MkdirAll(addFolderPath, 0755); err != nil {
		return fmt.Errorf("创建add文件夹失败: %v", err)
	}

	// 移动图片到add文件夹
	tempImagesPath := filepath.Join("data", "temp", "tempImages")
	if _, err := os.Stat(tempImagesPath); err == nil {
		files, err := os.ReadDir(tempImagesPath)
		if err != nil {
			return fmt.Errorf("读取图片目录失败: %v", err)
		}

		for _, file := range files {
			if !file.IsDir() && isImageFile(file.Name()) {
				srcPath := filepath.Join(tempImagesPath, file.Name())
				dstPath := filepath.Join(addFolderPath, file.Name())

				if err := os.Rename(srcPath, dstPath); err != nil {
					// 如果移动失败，尝试复制
					if err := copyFile(srcPath, dstPath); err != nil {
						return fmt.Errorf("移动图片失败: %v", err)
					}
				}
			}
		}
	}

	// 打包ZIP
	zipFilePath := filepath.Join("data", "output", bankName+".zip")
	if bankName == "" {
		zipFilePath = filepath.Join("data", "output", "未命名题库.zip")
	}

	err = createZip(outputTempDir, zipFilePath)
	if err != nil {
		return fmt.Errorf("打包ZIP失败: %v", err)
	}

	// 清理临时add文件夹
	os.RemoveAll(addFolderPath)

	return nil
}

// 创建ZIP文件
func createZip(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 创建ZIP文件头
		zipHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		zipHeader.Name = relPath

		if info.IsDir() {
			zipHeader.Name += "/"
		} else {
			zipHeader.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		// 如果是文件，写入内容
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// 获取图片文件信息
func getImageFilesInfo() []map[string]string {
	tempImagesPath := filepath.Join("data", "temp", "tempImages")
	var imageFiles []map[string]string

	if _, err := os.Stat(tempImagesPath); os.IsNotExist(err) {
		return imageFiles
	}

	files, err := os.ReadDir(tempImagesPath)
	if err != nil {
		return imageFiles
	}

	for _, file := range files {
		if !file.IsDir() && isImageFile(file.Name()) {
			imageInfo := map[string]string{
				"filename": file.Name(),
				"path":     filepath.ToSlash(filepath.Join("./add", file.Name())), // 修复：使用正斜杠
			}
			imageFiles = append(imageFiles, imageInfo)
		}
	}

	return imageFiles
}

// 构建图片路径并检查文件扩展名
func buildImagePath(imageName string) string {
	// 如果图片名已经包含扩展名，直接使用
	if hasImageExtension(imageName) {
		return filepath.ToSlash(filepath.Join("./add", imageName))
	}

	// 如果没有扩展名，尝试查找实际存在的图片文件
	tempImagesPath := filepath.Join("data", "temp", "tempImages")
	actualFileName := findActualImageFile(tempImagesPath, imageName)
	if actualFileName != "" {
		return filepath.ToSlash(filepath.Join("./add", actualFileName))
	}

	// 如果找不到实际文件，返回空字符串（不添加该图片）
	return ""
}

// 检查字符串是否包含图片扩展名
func hasImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExtensions := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".tiff", ".tif", ".webp"}

	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// 在目录中查找实际存在的图片文件
func findActualImageFile(dirPath, baseName string) string {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return ""
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			// 检查文件名是否以baseName开头（忽略扩展名）
			nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			if nameWithoutExt == baseName && isImageFile(fileName) {
				return fileName
			}
		}
	}

	return ""
}
