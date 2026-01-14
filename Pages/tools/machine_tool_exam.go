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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExamQuestion 题目结构体
type ExamQuestion struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Enabled    bool          `json:"enabled"`
	Question   string        `json:"question"`
	Images     []string      `json:"images"`
	Template   string        `json:"template,omitempty"`
	BlankCount int           `json:"blankCount,omitempty"`
	Options    []string      `json:"options,omitempty"`
	Answer     interface{}   `json:"answer,omitempty"` // for SC
	Answers    []interface{} `json:"answers"`          // for MC/FL
	HasExtra   bool          `json:"hasExtra,omitempty"`
	ExtraKey   string        `json:"extraKey,omitempty"`
	Hook       string        `json:"hook,omitempty"`      // 新增：钩子字段
	Materials  []string      `json:"materials,omitempty"` // 新增：材料内容
	Hooks      []string      `json:"hooks,omitempty"`     // 新增：钩子数组（用于DR）
}

// AnswerItem 答案项结构体
type AnswerItem struct {
	BlankIndex int      `json:"blankIndex,omitempty"`
	Answers    []string `json:"answers,omitempty"`
	Options    []string `json:"options,omitempty"`
	Correct    []int    `json:"correct,omitempty"`
	Answer     string   `json:"answer,omitempty"`
}

// ExamData 题库数据结构体
type ExamData struct {
	Errors []interface{} `json:"errors"`
	Images []struct {
		Filename string `json:"filename"`
		Path     string `json:"path"`
	} `json:"images"`
	Metadata struct {
		FillBlank       int `json:"fillBlank"`
		MultipleChoice  int `json:"multipleChoice"`
		SingleChoice    int `json:"singleChoice"`
		DocumentReading int `json:"documentReading"` // 新增
		TotalImages     int `json:"totalImages"`
		TotalQuestions  int `json:"totalQuestions"`
	} `json:"metadata"`
	Name      string `json:"name"`
	Questions struct {
		FillBlank       []ExamQuestion `json:"fillBlank"`
		MultipleChoice  []ExamQuestion `json:"multipleChoice"`
		SingleChoice    []ExamQuestion `json:"singleChoice"`
		DocumentReading []ExamQuestion `json:"documentReading"` // 新增
	} `json:"questions"`
	Version string `json:"version"`
}

// GroupedExamData 分组考试数据结构体
type GroupedExamData struct {
	SingleChoice    []GroupedQuestion `json:"singleChoice"`
	MultipleChoice  []GroupedQuestion `json:"multipleChoice"`
	FillBlank       []GroupedQuestion `json:"fillBlank"`
	DocumentReading []GroupedQuestion `json:"documentReading"` // 新增
}

// GroupedQuestion 分组题目结构体
type GroupedQuestion struct {
	ID         string   `json:"id"`
	Index      int      `json:"index"`
	TotalIndex int      `json:"totalIndex"`
	Question   string   `json:"question"`
	Options    []string `json:"options,omitempty"`
	Images     []string `json:"images,omitempty"`
	Answer     string   `json:"answer,omitempty"`    // for SC/MC, the option content
	Blanks     []Blank  `json:"blanks,omitempty"`    // for FL
	Hook       string   `json:"hook,omitempty"`      // 新增：钩子字段
	Materials  []string `json:"materials,omitempty"` // 新增：材料内容
}

// Blank 填空结构体
type Blank struct {
	BlankIndex int      `json:"blankIndex"`
	Answers    []string `json:"answers"`
}

// CorrectAnswer 正确答案结构体
type CorrectAnswer struct {
	Questions map[string][]string `json:"questions"`
}

// DRGroup 材料阅读题组
type DRGroup struct {
	DRQuestion     ExamQuestion
	ChildQuestions []ExamQuestion
}

// findExamJSONFile 在目录中查找题库JSON文件
func findExamJSONFile(examTempPath string) (string, error) {
	files, err := os.ReadDir(examTempPath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			// 排除正确答案文件
			if strings.ToLower(file.Name()) != "correctanswer.json" {
				return file.Name(), nil
			}
		}
	}

	return "", fmt.Errorf("未找到题库JSON文件")
}

// MachineToolLoadExamData 加载题库数据
func MachineToolLoadExamData(examTempPath string) (*ExamData, error) {
	// 动态查找JSON文件
	jsonFileName, err := findExamJSONFile(examTempPath)
	if err != nil {
		return nil, fmt.Errorf("查找题库文件失败: %v", err)
	}

	jsonPath := filepath.Join(examTempPath, jsonFileName)
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("打开题库文件失败: %v", err)
	}
	defer file.Close()

	var data ExamData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("解析题库json失败: %v", err)
	}
	return &data, nil
}

// MachineToolLoadCorrectAnswers 加载正确答案（现在加载grouped data）
func MachineToolLoadCorrectAnswers(examTempPath string) (*GroupedExamData, error) {
	jsonPath := filepath.Join(examTempPath, "CorrectAnswer.json")
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("打开正确答案文件失败: %v", err)
	}
	defer file.Close()

	var grouped GroupedExamData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&grouped); err != nil {
		return nil, fmt.Errorf("解析正确答案json失败: %v", err)
	}
	return &grouped, nil
}

// MachineToolGetAllQuestions 获取所有题目列表（按顺序）- 修复版：排除有hook的题目
func MachineToolGetAllQuestions(data *ExamData) []ExamQuestion {
	var questions []ExamQuestion

	// 处理材料阅读题（DR）
	for _, q := range data.Questions.DocumentReading {
		q.ID = "DR_" + q.ID
		questions = append(questions, q)
	}

	// 处理单选题（SC + SCIMG）- 排除有hook的题目
	for _, q := range data.Questions.SingleChoice {
		// 跳过有hook的题目（它们会在DR卡片中显示）
		if q.Hook != "" {
			continue
		}
		prefix := "SC"
		if q.Type == "SCIMG" {
			prefix = "SCIMG"
		}
		q.ID = prefix + "_" + q.ID
		questions = append(questions, q)
	}

	// 处理多选题（MC + MCIMG）- 排除有hook的题目
	for _, q := range data.Questions.MultipleChoice {
		// 跳过有hook的题目（它们会在DR卡片中显示）
		if q.Hook != "" {
			continue
		}
		prefix := "MC"
		if q.Type == "MCIMG" {
			prefix = "MCIMG"
		}
		q.ID = prefix + "_" + q.ID
		questions = append(questions, q)
	}

	// 处理填空题（FL + FLIMG）- 排除有hook的题目
	for _, q := range data.Questions.FillBlank {
		// 跳过有hook的题目（它们会在DR卡片中显示）
		if q.Hook != "" {
			continue
		}
		prefix := "FL"
		if q.Type == "FLIMG" {
			prefix = "FLIMG"
		}
		q.ID = prefix + "_" + q.ID
		questions = append(questions, q)
	}

	return questions
}

// MachineToolGroupDRQuestions 分组材料阅读题 - 修复版：根据DR的hooks数组匹配子题
func MachineToolGroupDRQuestions(data *ExamData) []DRGroup {
	var drGroups []DRGroup

	// 收集所有有钩子的题目
	var allHookQuestions []ExamQuestion

	// 收集单选题的钩子题目
	for _, q := range data.Questions.SingleChoice {
		if q.Hook != "" {
			// 设置正确的ID前缀
			prefix := "SC"
			if q.Type == "SCIMG" {
				prefix = "SCIMG"
			}
			q.ID = prefix + "_" + q.ID
			allHookQuestions = append(allHookQuestions, q)
		}
	}

	// 收集多选题的钩子题目
	for _, q := range data.Questions.MultipleChoice {
		if q.Hook != "" {
			// 设置正确的ID前缀
			prefix := "MC"
			if q.Type == "MCIMG" {
				prefix = "MCIMG"
			}
			q.ID = prefix + "_" + q.ID
			allHookQuestions = append(allHookQuestions, q)
		}
	}

	// 收集填空题的钩子题目
	for _, q := range data.Questions.FillBlank {
		if q.Hook != "" {
			// 设置正确的ID前缀
			prefix := "FL"
			if q.Type == "FLIMG" {
				prefix = "FLIMG"
			}
			q.ID = prefix + "_" + q.ID
			allHookQuestions = append(allHookQuestions, q)
		}
	}

	// 为每个DR题匹配子题
	for _, dr := range data.Questions.DocumentReading {
		// 设置DR题的ID前缀
		dr.ID = "DR_" + dr.ID

		var childQuestions []ExamQuestion
		// 根据DR的hooks数组匹配子题
		for _, q := range allHookQuestions {
			for _, h := range dr.Hooks {
				if q.Hook == h {
					childQuestions = append(childQuestions, q)
					break
				}
			}
		}

		if len(childQuestions) > 0 {
			drGroups = append(drGroups, DRGroup{
				DRQuestion:     dr,
				ChildQuestions: childQuestions,
			})
		}
	}

	return drGroups
}

// MachineToolGenerateCorrectAnswers 从ExamData生成正确答案（现在生成grouped JSON）- 修复版
func MachineToolGenerateCorrectAnswers(data *ExamData) *GroupedExamData {
	grouped := &GroupedExamData{}
	totalIndex := 1

	// DocumentReading (DR)
	for idx, q := range data.Questions.DocumentReading {
		gq := GroupedQuestion{
			ID:         "DR_" + q.ID,
			Index:      idx + 1,
			TotalIndex: totalIndex,
			Question:   q.Question,
			Images:     q.Images,
			Materials:  q.Materials,
			Hook:       q.Hook,
		}
		grouped.DocumentReading = append(grouped.DocumentReading, gq)
		totalIndex++
	}

	// SingleChoice (SC + SCIMG) - 包含所有题目，包括有hook的
	for idx, q := range data.Questions.SingleChoice {
		prefix := "SC"
		if q.Type == "SCIMG" {
			prefix = "SCIMG"
		}

		gq := GroupedQuestion{
			ID:         prefix + "_" + q.ID,
			Index:      idx + 1,
			TotalIndex: totalIndex,
			Question:   q.Question,
			Options:    q.Options,
			Images:     q.Images,
			Hook:       q.Hook,
		}

		// 处理单选题答案
		if ansStr, ok := q.Answer.(string); ok && ansStr != "" {
			index := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", ansStr)
			if index >= 0 && index < len(q.Options) {
				gq.Answer = q.Options[index]
			} else {
				gq.Answer = ansStr
			}
		}
		grouped.SingleChoice = append(grouped.SingleChoice, gq)
		totalIndex++
	}

	// MultipleChoice (MC + MCIMG) - 包含所有题目，包括有hook的
	for idx, q := range data.Questions.MultipleChoice {
		prefix := "MC"
		if q.Type == "MCIMG" {
			prefix = "MCIMG"
		}

		gq := GroupedQuestion{
			ID:         prefix + "_" + q.ID,
			Index:      idx + 1,
			TotalIndex: totalIndex,
			Question:   q.Question,
			Options:    q.Options,
			Images:     q.Images,
			Hook:       q.Hook,
		}
		var answers []string
		for _, a := range q.Answers {
			if s, ok := a.(string); ok {
				parts := strings.Split(strings.TrimSpace(s), ";")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if part != "" {
						index := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", part)
						if index >= 0 && index < len(q.Options) {
							answers = append(answers, q.Options[index])
						} else {
							answers = append(answers, part)
						}
					}
				}
			}
		}
		gq.Answer = strings.Join(answers, "; ")
		grouped.MultipleChoice = append(grouped.MultipleChoice, gq)
		totalIndex++
	}

	// FillBlank (FL + FLIMG) - 包含所有题目，包括有hook的
	for idx, q := range data.Questions.FillBlank {
		prefix := "FL"
		if q.Type == "FLIMG" {
			prefix = "FLIMG"
		}

		gq := GroupedQuestion{
			ID:         prefix + "_" + q.ID,
			Index:      idx + 1,
			TotalIndex: totalIndex,
			Question:   q.Question,
			Images:     q.Images,
			Hook:       q.Hook,
		}
		// 按blankIndex排序
		var blankMap = make(map[int][]string)
		var blankIndexes []int
		for _, ansItem := range q.Answers {
			if item, ok := ansItem.(map[string]interface{}); ok {
				if idx, ok := item["blankIndex"].(float64); ok {
					blankIndex := int(idx)
					if !containsInt(blankIndexes, blankIndex) {
						blankIndexes = append(blankIndexes, blankIndex)
					}
					if ansList, exists := item["answers"]; exists {
						if list, ok := ansList.([]interface{}); ok {
							for _, v := range list {
								if str, ok := v.(string); ok {
									blankMap[blankIndex] = append(blankMap[blankIndex], str)
								}
							}
						}
					}
				}
			}
		}
		sort.Ints(blankIndexes)
		for _, bidx := range blankIndexes {
			gq.Blanks = append(gq.Blanks, Blank{
				BlankIndex: bidx,
				Answers:    blankMap[bidx],
			})
		}
		grouped.FillBlank = append(grouped.FillBlank, gq)
		totalIndex++
	}

	return grouped
}

// containsInt 检查slice是否包含int
func containsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// MachineToolSaveCorrectAnswers 保存正确答案到JSON文件（现在保存grouped data）
func MachineToolSaveCorrectAnswers(examTempPath string, grouped *GroupedExamData) error {
	jsonPath := filepath.Join(examTempPath, "CorrectAnswer.json")
	file, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("创建正确答案文件失败: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(grouped); err != nil {
		return fmt.Errorf("写入正确答案json失败: %v", err)
	}
	return nil
}

func GetOriginalExamData(examTempPath string) (*ExamData, error) {
	return MachineToolLoadExamData(examTempPath)
}

// MachineToolGetAllQuestions
