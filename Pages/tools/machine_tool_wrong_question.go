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
	"time"
)

// WrongQuestionSet 错题集结构
type WrongQuestionSet struct {
	SourceExam     string          `json:"sourceExam"`
	ExamTime       string          `json:"examTime"`
	WrongCount     int             `json:"wrongCount"`
	TotalCount     int             `json:"totalCount"`
	WrongQuestions []WrongQuestion `json:"wrongQuestions"`
}

// WrongQuestion 错题结构
type WrongQuestion struct {
	ID                string          `json:"id"`
	Type              string          `json:"type"`
	Question          string          `json:"question"`
	Images            []string        `json:"images"`
	Options           []string        `json:"options,omitempty"`
	BlankCount        int             `json:"blankCount,omitempty"`
	CorrectAnswer     interface{}     `json:"correctAnswer"`
	Hook              string          `json:"hook,omitempty"`              // 新增：钩子字段
	Materials         []string        `json:"materials,omitempty"`         // 新增：材料内容
	SubWrongQuestions []WrongQuestion `json:"subWrongQuestions,omitempty"` // 新增：子错题列表，用于DR题分组
}

// CreateWrongQuestionCollection 创建错题集
func CreateWrongQuestionCollection(examState *ExamState, originalData *ExamData, examTempPath string) error {
	// 使用统一的时间格式
	examTime := time.Now().Format("20060102_150405")
	gradeDir := "data/grades"
	examFolder := fmt.Sprintf("%s_%s", examState.ExamName, examTime)
	fullExamPath := filepath.Join(gradeDir, examFolder)
	sourseDir := filepath.Join(fullExamPath, "sourse")
	addDir := filepath.Join(sourseDir, "add")

	if err := os.MkdirAll(addDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 识别错题
	wrongQuestions, err := findWrongQuestions(examState)
	if err != nil {
		return fmt.Errorf("查找错题失败: %v", err)
	}

	// 生成错题数据
	wrongSet := WrongQuestionSet{
		SourceExam:     examState.ExamName,
		ExamTime:       examTime,
		WrongCount:     len(wrongQuestions),
		TotalCount:     len(examState.Questions),
		WrongQuestions: wrongQuestions,
	}

	// 保存错题JSON
	wrongJSONName := fmt.Sprintf("wrong_%s_%s.json", examState.ExamName, examTime)
	wrongJSONPath := filepath.Join(sourseDir, wrongJSONName)
	if err := saveWrongSetToFile(wrongSet, wrongJSONPath); err != nil {
		return fmt.Errorf("保存错题JSON失败: %v", err)
	}

	// 复制图片文件
	if err := copyQuestionImages(wrongQuestions, examTempPath, addDir); err != nil {
		return fmt.Errorf("复制图片失败: %v", err)
	}

	// 创建ZIP文件
	zipPath := filepath.Join(fullExamPath, "sourse.zip")
	if err := createQuestionZip(sourseDir, zipPath); err != nil {
		return fmt.Errorf("创建ZIP失败: %v", err)
	}

	fmt.Printf("错题集生成成功: %s, 包含 %d 道错题\n", zipPath, len(wrongQuestions))
	return nil
}

// findWrongQuestions 查找错题 - 修复版，支持多个DR题分组和子题批改
func findWrongQuestions(examState *ExamState) ([]WrongQuestion, error) {
	var wrongList []WrongQuestion
	correctAnswers := examState.GetAllCorrectAnswers()

	questionNumber := 1 // 错题序号从1开始

	// 创建DR组的子错题映射
	drWrongSubsMap := make(map[string][]WrongQuestion)

	// 首先处理DR组中的子题
	for _, drGroup := range examState.DRGroups {
		for _, childQ := range drGroup.ChildQuestions {
			userAnswer := examState.GetAnswer(childQ.ID)
			correctAnswer := correctAnswers[childQ.ID]
			if correctAnswer == nil {
				fmt.Printf("警告: 未找到题目 %s 的正确答案\n", childQ.ID)
				continue
			}

			// 调试信息
			fmt.Printf("检查DR子题错题 %s: 用户答案=%v, 正确答案=%v\n",
				childQ.ID, userAnswer, correctAnswer)

			isCorrect := checkAnswerCorrectness(childQ.Type, userAnswer, correctAnswer)

			if !isCorrect {
				wrongQ := buildWrongQuestion(childQ, correctAnswer, 0) // 暂时用0，后续设置序号
				wrongQ.ID = childQ.ID                                  // 子问题用原始ID
				drWrongSubsMap[drGroup.DRQuestion.ID] = append(drWrongSubsMap[drGroup.DRQuestion.ID], wrongQ)
				fmt.Printf("识别为DR子错题: %s (DR: %s)\n", childQ.ID, drGroup.DRQuestion.ID)
			}
		}
	}

	// 然后处理独立题目（examState.Questions中不包含有hook的题目）
	for _, question := range examState.Questions {
		// 跳过DR题本身（不计分）
		if question.Type == "DR" {
			continue
		}

		userAnswer := examState.GetAnswer(question.ID)
		correctAnswer := correctAnswers[question.ID]
		if correctAnswer == nil {
			fmt.Printf("警告: 未找到题目 %s 的正确答案\n", question.ID)
			continue
		}

		// 调试信息
		fmt.Printf("检查独立错题 %s: 用户答案=%v, 正确答案=%v\n",
			question.ID, userAnswer, correctAnswer)

		isCorrect := checkAnswerCorrectness(question.Type, userAnswer, correctAnswer)

		if !isCorrect {
			wrongQ := buildWrongQuestion(question, correctAnswer, 0) // 暂时用0，后续设置序号
			wrongQ.ID = fmt.Sprintf("%d", questionNumber)
			wrongList = append(wrongList, wrongQ)
			questionNumber++
			fmt.Printf("识别为独立错题: %s\n", question.ID)
		}
	}

	// 处理DR题分组
	for _, drGroup := range examState.DRGroups {
		drID := drGroup.DRQuestion.ID
		if subWrongs, exists := drWrongSubsMap[drID]; exists && len(subWrongs) > 0 {
			// 有错误子问题，创建DR错题条目
			drWrongQ := WrongQuestion{
				ID:                fmt.Sprintf("%d", questionNumber),
				Type:              drGroup.DRQuestion.Type,
				Question:          drGroup.DRQuestion.Question,
				Images:            drGroup.DRQuestion.Images,
				CorrectAnswer:     nil, // DR题本身无答案
				Hook:              drGroup.DRQuestion.Hook,
				Materials:         drGroup.DRQuestion.Materials,
				SubWrongQuestions: subWrongs,
			}
			wrongList = append(wrongList, drWrongQ)
			questionNumber++
			fmt.Printf("识别为DR分组错题: %s (子错题数: %d)\n", drID, len(subWrongs))
		}
	}

	fmt.Printf("总共识别出 %d 道错题\n", len(wrongList))
	return wrongList, nil
}

// buildWrongQuestion 构建错题对象 - 修复版
func buildWrongQuestion(q ExamQuestion, correctAnswer interface{}, questionNumber int) WrongQuestion {
	wrongQ := WrongQuestion{
		ID:            fmt.Sprintf("%d", questionNumber), // 重写ID为错题序号
		Type:          q.Type,
		Question:      q.Question,
		Images:        q.Images,
		CorrectAnswer: correctAnswer, // 直接使用传入的正确答案
		Hook:          q.Hook,        // 新增：钩子字段
		Materials:     q.Materials,   // 新增：材料内容
	}

	// 根据题型设置相应字段
	switch {
	case q.Type == "SC" || q.Type == "SCIMG" || q.Type == "MC" || q.Type == "MCIMG":
		wrongQ.Options = q.Options
	case q.Type == "FL" || q.Type == "FLIMG":
		wrongQ.BlankCount = q.BlankCount
		// 确保填空题答案格式正确
		if wrongQ.CorrectAnswer == nil {
			// 如果没有正确答案，创建空答案结构
			var blanks [][]string
			for i := 0; i < q.BlankCount; i++ {
				blanks = append(blanks, []string{})
			}
			wrongQ.CorrectAnswer = blanks
		}
	}

	return wrongQ
}

// checkAnswerCorrectness 检查答案正确性 - 修复版
func checkAnswerCorrectness(qType string, userAnswer []string, correctAnswer interface{}) bool {
	// 合并处理同类型题目
	switch {
	case qType == "FL" || qType == "FLIMG":
		return checkFillBlankCorrectness(userAnswer, correctAnswer)
	case qType == "SC" || qType == "SCIMG" || qType == "MC" || qType == "MCIMG":
		return checkChoiceCorrectness(userAnswer, correctAnswer, qType)
	default:
		fmt.Printf("未知题目类型: %s\n", qType)
		return false
	}
}

// checkFillBlankCorrectness 检查填空题正确性 - 修复版
func checkFillBlankCorrectness(userAnswer []string, correctAnswer interface{}) bool {
	correctBlanks, ok := correctAnswer.([][]string)
	if !ok {
		fmt.Printf("填空题答案格式错误: %T\n", correctAnswer)
		return false
	}

	if len(userAnswer) != len(correctBlanks) {
		fmt.Printf("填空题数量不匹配: 用户%d空, 正确答案%d空\n", len(userAnswer), len(correctBlanks))
		return false
	}

	// 检查每个空是否正确
	for i, userInput := range userAnswer {
		userInput = strings.TrimSpace(userInput)
		if userInput == "" {
			return false // 空答案直接判错
		}

		found := false
		for _, correctOption := range correctBlanks[i] {
			if userInput == strings.TrimSpace(correctOption) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查唯一性约束
	return checkUniqueConstraint(userAnswer, correctBlanks)
}

// checkUniqueConstraint 检查唯一性约束
func checkUniqueConstraint(userAnswer []string, correctBlanks [][]string) bool {
	hasUniqueMarker := false
	for _, blank := range correctBlanks {
		for _, answer := range blank {
			if strings.TrimSpace(answer) == "(x%x)" {
				hasUniqueMarker = true
				break
			}
		}
		if hasUniqueMarker {
			break
		}
	}

	if hasUniqueMarker {
		seen := make(map[string]bool)
		for _, answer := range userAnswer {
			trimmed := strings.TrimSpace(answer)
			if seen[trimmed] {
				return false
			}
			seen[trimmed] = true
		}
	}

	return true
}

// checkChoiceCorrectness 检查选择题正确性 - 修复版
func checkChoiceCorrectness(userAnswer []string, correctAnswer interface{}, qType string) bool {
	correctList, ok := correctAnswer.([]string)
	if !ok {
		fmt.Printf("选择题答案格式错误: %T\n", correctAnswer)
		return false
	}

	// 多选题需要排序比较（顺序无关）
	if qType == "MC" || qType == "MCIMG" {
		if len(userAnswer) != len(correctList) {
			return false
		}

		userSorted := make([]string, len(userAnswer))
		correctSorted := make([]string, len(correctList))
		copy(userSorted, userAnswer)
		copy(correctSorted, correctList)

		sort.Strings(userSorted)
		sort.Strings(correctSorted)

		for i := range userSorted {
			if strings.TrimSpace(userSorted[i]) != strings.TrimSpace(correctSorted[i]) {
				return false
			}
		}
		return true
	}

	// 单选题直接比较
	if len(userAnswer) != 1 || len(correctList) != 1 {
		return false
	}

	return strings.TrimSpace(userAnswer[0]) == strings.TrimSpace(correctList[0])
}

// saveWrongSetToFile 保存错题集到文件
func saveWrongSetToFile(wrongSet WrongQuestionSet, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(wrongSet)
}

// copyQuestionImages 复制题目图片
func copyQuestionImages(wrongQuestions []WrongQuestion, sourceDir, targetDir string) error {
	copyCount := 0
	failCount := 0

	for _, question := range wrongQuestions {
		// 复制当前错题的图片
		for _, imageName := range question.Images {
			sourcePath := filepath.Join(sourceDir, imageName)
			targetPath := filepath.Join(targetDir, filepath.Base(imageName)) // 只取文件名，避免嵌套

			// 确保目标目录存在
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				fmt.Printf("创建目录失败 %s: %v\n", filepath.Dir(targetPath), err)
				failCount++
				continue
			}

			// 检查源文件是否存在
			if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
				fmt.Printf("源图片不存在: %s\n", sourcePath)
				failCount++
				continue
			}

			if err := copySingleFile(sourcePath, targetPath); err != nil {
				fmt.Printf("复制图片失败 %s -> %s: %v\n", sourcePath, targetPath, err)
				failCount++
			} else {
				copyCount++
				fmt.Printf("成功复制图片: %s\n", filepath.Base(imageName))
			}
		}

		// 递归复制子错题的图片
		if len(question.SubWrongQuestions) > 0 {
			if err := copyQuestionImages(question.SubWrongQuestions, sourceDir, targetDir); err != nil {
				fmt.Printf("复制子错题图片失败: %v\n", err)
				failCount++
			}
		}
	}

	fmt.Printf("图片复制完成: 成功 %d 个, 失败 %d 个\n", copyCount, failCount)
	return nil
}

// copySingleFile 复制单个文件
func copySingleFile(source, destination string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return os.WriteFile(destination, data, 0644)
}

// createQuestionZip 创建题目ZIP文件
func createQuestionZip(sourceDir, zipPath string) error {
	return CreateZipArchive(sourceDir, zipPath)
}
