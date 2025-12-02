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

// ExamState 考试状态结构体
type ExamState struct {
	ExamName       string
	StartTime      time.Time
	Answers        map[string][]string // 题目ID -> 用户答案列表
	Questions      []ExamQuestion
	CorrectAnswers *GroupedExamData
	DRGroups       []DRGroup // 新增：材料阅读题分组
}

// NewExamState 创建考试状态
func NewExamState(examName string, questions []ExamQuestion, correctAnswers *GroupedExamData) *ExamState {
	if correctAnswers == nil {
		fmt.Println("警告: GroupedExamData 为空，将使用备用方案获取正确答案")
	}

	return &ExamState{
		ExamName:       examName,
		StartTime:      time.Now(),
		Answers:        make(map[string][]string),
		Questions:      questions,
		CorrectAnswers: correctAnswers,
	}
}

// SetAnswer 设置用户答案
func (es *ExamState) SetAnswer(questionID string, answers []string) {
	es.Answers[questionID] = answers
}

// GetAnswer 获取用户答案
func (es *ExamState) GetAnswer(questionID string) []string {
	return es.Answers[questionID]
}

// IsAnswered 检查题目是否已答
func (es *ExamState) IsAnswered(questionID string) bool {
	answers, exists := es.Answers[questionID]
	if !exists {
		return false
	}
	if len(answers) == 0 {
		return false
	}
	// 对于填空题，检查是否所有答案都是空的
	for _, ans := range answers {
		if strings.TrimSpace(ans) != "" {
			return true
		}
	}
	return false
}

// IsDRGroupAnswered 检查材料阅读题组是否已答
func (es *ExamState) IsDRGroupAnswered(drGroup DRGroup) bool {
	// 检查DR题的子题中是否有任意一个已答
	for _, child := range drGroup.ChildQuestions {
		if es.IsAnswered(child.ID) {
			return true
		}
	}
	return false
}

// CalculateScore 计算得分 - 修复版
func (es *ExamState) CalculateScore() (finalScore float64, totalScore float64, correctCount int, wrongCount int) {
	// 调试：打印所有题目
	fmt.Println("=== 所有题目 ===")
	for i, q := range es.Questions {
		fmt.Printf("%d. ID: %s, Type: %s, Hook: %s\n", i+1, q.ID, q.Type, q.Hook)
	}
	for i, drGroup := range es.DRGroups {
		fmt.Printf("DR Group %d: %s\n", i+1, drGroup.DRQuestion.Question)
		for j, child := range drGroup.ChildQuestions {
			fmt.Printf("  子题 %d: ID: %s, Type: %s\n", j+1, child.ID, child.Type)
		}
	}

	// 计算总分和每题分值
	totalScore = 100.0 // 总分100分
	questionScores := es.calculateQuestionScores(totalScore)

	correctMap := es.GetAllCorrectAnswers()

	// 调试：打印所有正确答案
	fmt.Println("=== 正确答案映射 ===")
	for id, correct := range correctMap {
		fmt.Printf("题目 %s -> 正确答案: %v\n", id, correct)
	}

	earnedScore := 0.0
	correctCount = 0
	wrongCount = 0

	for _, q := range es.Questions {
		// 跳过材料阅读题（DR题本身不计分）
		if q.Type == "DR" {
			continue
		}

		userAns := es.GetAnswer(q.ID)
		correctAns := correctMap[q.ID]
		if correctAns == nil {
			fmt.Printf("警告: 题目 %s 没有正确答案\n", q.ID)
			continue
		}

		questionScore := questionScores[q.ID]
		earned := 0.0
		isFullyCorrect := false

		// 调试信息
		fmt.Printf("批改题目 %s: 用户答案=%v, 正确答案=%v, 选项=%v\n", q.ID, userAns, correctAns, q.Options)

		switch q.Type {
		case "FL", "FLIMG":
			// 填空题处理 - 按空计分，但只有全对才算正确题目
			if blanksCorrect, ok := correctAns.([][]string); ok {
				if len(userAns) == len(blanksCorrect) {
					// 计算每个空的得分
					blankScore := questionScore / float64(len(blanksCorrect))
					allCorrect := true

					// 检查每个用户输入是否在对应blank的答案中
					for i, u := range userAns {
						userInput := strings.TrimSpace(u)
						if userInput == "" {
							allCorrect = false
							continue
						}

						found := false
						for _, c := range blanksCorrect[i] {
							if userInput == strings.TrimSpace(c) {
								found = true
								earned += blankScore
								break
							}
						}
						if !found {
							allCorrect = false
						}
					}

					// 检查唯一性约束
					if allCorrect {
						hasUniqueMarker := false
						for _, blank := range blanksCorrect {
							for _, ans := range blank {
								if strings.TrimSpace(ans) == "(x%x)" {
									hasUniqueMarker = true
									break
								}
							}
							if hasUniqueMarker {
								break
							}
						}
						if hasUniqueMarker {
							// 检查用户输入是否各不相同
							seen := make(map[string]bool)
							allUnique := true
							for _, u := range userAns {
								trimmed := strings.TrimSpace(u)
								if seen[trimmed] {
									allUnique = false
									break
								}
								seen[trimmed] = true
							}
							if !allUnique {
								// 如果不唯一，扣除所有空的分数
								earned = 0
								allCorrect = false
							}
						}
					}

					isFullyCorrect = allCorrect
				}
			}
		default:
			// 单选题和多选题处理 - 全对才得分
			if correctList, ok := correctAns.([]string); ok {
				// 将正确答案从选项索引转换为选项内容
				correctAnswersContent := convertToOptionContent(correctList, q.Options)

				if compareAnswers(userAns, correctAnswersContent, q.Type, q.Options) {
					earned = questionScore
					isFullyCorrect = true
				}

				fmt.Printf("转换后比较: 用户答案=%v, 正确答案内容=%v\n", userAns, correctAnswersContent)
			}
		}

		earnedScore += earned

		// 统计正确和错误题目数量
		if isFullyCorrect {
			correctCount++
			fmt.Printf("题目 %s: 正确\n", q.ID)
		} else {
			wrongCount++
			fmt.Printf("题目 %s: 错误\n", q.ID)
		}
	}

	// 处理DR子题
	for _, drGroup := range es.DRGroups {
		for _, q := range drGroup.ChildQuestions {
			userAns := es.GetAnswer(q.ID)
			correctAns := correctMap[q.ID]
			if correctAns == nil {
				fmt.Printf("警告: 题目 %s 没有正确答案\n", q.ID)
				continue
			}

			questionScore := questionScores[q.ID]
			earned := 0.0
			isFullyCorrect := false

			// 调试信息
			fmt.Printf("批改题目 %s: 用户答案=%v, 正确答案=%v, 选项=%v\n", q.ID, userAns, correctAns, q.Options)

			switch q.Type {
			case "FL", "FLIMG":
				// 填空题处理 - 按空计分，但只有全对才算正确题目
				if blanksCorrect, ok := correctAns.([][]string); ok {
					if len(userAns) == len(blanksCorrect) {
						// 计算每个空的得分
						blankScore := questionScore / float64(len(blanksCorrect))
						allCorrect := true

						// 检查每个用户输入是否在对应blank的答案中
						for i, u := range userAns {
							userInput := strings.TrimSpace(u)
							if userInput == "" {
								allCorrect = false
								continue
							}

							found := false
							for _, c := range blanksCorrect[i] {
								if userInput == strings.TrimSpace(c) {
									found = true
									earned += blankScore
									break
								}
							}
							if !found {
								allCorrect = false
							}
						}

						// 检查唯一性约束
						if allCorrect {
							hasUniqueMarker := false
							for _, blank := range blanksCorrect {
								for _, ans := range blank {
									if strings.TrimSpace(ans) == "(x%x)" {
										hasUniqueMarker = true
										break
									}
								}
								if hasUniqueMarker {
									break
								}
							}
							if hasUniqueMarker {
								// 检查用户输入是否各不相同
								seen := make(map[string]bool)
								allUnique := true
								for _, u := range userAns {
									trimmed := strings.TrimSpace(u)
									if seen[trimmed] {
										allUnique = false
										break
									}
									seen[trimmed] = true
								}
								if !allUnique {
									// 如果不唯一，扣除所有空的分数
									earned = 0
									allCorrect = false
								}
							}
						}

						isFullyCorrect = allCorrect
					}
				}
			default:
				// 单选题和多选题处理 - 全对才得分
				if correctList, ok := correctAns.([]string); ok {
					// 将正确答案从选项索引转换为选项内容
					correctAnswersContent := convertToOptionContent(correctList, q.Options)

					if compareAnswers(userAns, correctAnswersContent, q.Type, q.Options) {
						earned = questionScore
						isFullyCorrect = true
					}

					fmt.Printf("转换后比较: 用户答案=%v, 正确答案内容=%v\n", userAns, correctAnswersContent)
				}
			}

			earnedScore += earned

			// 统计正确和错误题目数量
			if isFullyCorrect {
				correctCount++
				fmt.Printf("题目 %s: 正确\n", q.ID)
			} else {
				wrongCount++
				fmt.Printf("题目 %s: 错误\n", q.ID)
			}
		}
	}

	// 保留一位小数
	finalScore = roundToOneDecimal(earnedScore)
	totalScore = roundToOneDecimal(totalScore)

	fmt.Printf("考试结果: 得分 %.1f/%.1f, 正确题数 %d, 错题数 %d\n",
		finalScore, totalScore, correctCount, wrongCount)

	return finalScore, totalScore, correctCount, wrongCount
}

// calculateQuestionScores 计算每道题的分值
func (es *ExamState) calculateQuestionScores(totalScore float64) map[string]float64 {
	// 只计算需要计分的题目（排除DR题，但包括DR子题）
	scoredQuestions := 0
	for _, q := range es.Questions {
		if q.Type != "DR" {
			scoredQuestions++
		}
	}
	for _, drGroup := range es.DRGroups {
		scoredQuestions += len(drGroup.ChildQuestions)
	}

	questionScores := make(map[string]float64)

	// 先计算基础分值（平均分配）
	baseScore := totalScore / float64(scoredQuestions)

	// 为每道题分配分值，考虑填空题的空数
	for _, q := range es.Questions {
		// DR题不计分
		if q.Type == "DR" {
			questionScores[q.ID] = 0
			continue
		}

		score := baseScore

		// 如果是填空题，根据空数调整分值权重
		if q.Type == "FL" || q.Type == "FLIMG" {
			if q.BlankCount > 0 {
				// 填空题分值略高于选择题，因为包含多个空
				score = baseScore * 1.2
			}
		}

		questionScores[q.ID] = score
	}

	// 为DR子题分配分值
	for _, drGroup := range es.DRGroups {
		for _, child := range drGroup.ChildQuestions {
			score := baseScore

			// 如果是填空题，根据空数调整分值权重
			if child.Type == "FL" || child.Type == "FLIMG" {
				if child.BlankCount > 0 {
					// 填空题分值略高于选择题，因为包含多个空
					score = baseScore * 1.2
				}
			}

			questionScores[child.ID] = score
		}
	}

	// 调整总分确保正好是100分
	actualTotal := 0.0
	for _, score := range questionScores {
		actualTotal += score
	}

	// 如果总分不等于100，按比例调整
	if actualTotal != totalScore {
		ratio := totalScore / actualTotal
		for id := range questionScores {
			questionScores[id] = questionScores[id] * ratio
		}
	}

	return questionScores
}

// roundToOneDecimal 保留一位小数
func roundToOneDecimal(score float64) float64 {
	return float64(int(score*10+0.5)) / 10
}

// compareAnswers 比较答案（修复版）
func compareAnswers(user, correct []string, qType string, options []string) bool {
	// 多选题需要排序比较（顺序无关）
	if qType == "MC" || qType == "MCIMG" {
		if len(user) != len(correct) {
			return false
		}

		userSorted := make([]string, len(user))
		correctSorted := make([]string, len(correct))
		copy(userSorted, user)
		copy(correctSorted, correct)

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
	if len(user) != len(correct) {
		return false
	}
	for i, u := range user {
		if strings.TrimSpace(u) != strings.TrimSpace(correct[i]) {
			return false
		}
	}
	return true
}

// convertToOptionContent 将选项索引转换为选项内容
func convertToOptionContent(answerIndexes []string, options []string) []string {
	if len(options) == 0 {
		return answerIndexes
	}

	var contents []string
	for _, index := range answerIndexes {
		// 处理单个字符的选项索引（如 "A", "B", "C"）
		if len(index) == 1 {
			idx := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", strings.ToUpper(index))
			if idx >= 0 && idx < len(options) {
				contents = append(contents, options[idx])
			} else {
				contents = append(contents, index) // 如果不是有效索引，保持原样
			}
		} else {
			// 处理多选题的复合答案（如 "A;B;C"）
			parts := strings.Split(index, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				idx := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", strings.ToUpper(part))
				if idx >= 0 && idx < len(options) {
					contents = append(contents, options[idx])
				} else {
					contents = append(contents, part)
				}
			}
		}
	}
	return contents
}

// SaveGrade 保存成绩
func (es *ExamState) SaveGrade(gradesPath string) error {
	os.MkdirAll(gradesPath, 0755)
	finalScore, totalScore, correctCount, wrongCount := es.CalculateScore()
	grade := map[string]interface{}{
		"examName":     es.ExamName,
		"finalScore":   finalScore,
		"totalScore":   totalScore,
		"correctCount": correctCount,
		"wrongCount":   wrongCount,
		"startTime":    es.StartTime.Format("2006-01-02 15:04:05"),
		"endTime":      time.Now().Format("2006-01-02 15:04:05"),
		"duration":     time.Since(es.StartTime).String(),
	}

	// 创建文件名：题库名+考试结束时间
	fileName := fmt.Sprintf("%s_%s.json", es.ExamName, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(gradesPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建成绩文件失败: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(grade)
}

// GetCorrectAnswers 获取正确答案map
func (ged *GroupedExamData) GetCorrectAnswers() map[string]interface{} {
	correct := make(map[string]interface{})

	// 材料阅读题（DR题本身没有答案）
	for _, q := range ged.DocumentReading {
		correct[q.ID] = nil
	}

	// 单选题
	for _, q := range ged.SingleChoice {
		correct[q.ID] = []string{q.Answer}
	}

	// 多选题
	for _, q := range ged.MultipleChoice {
		answers := strings.Split(strings.TrimSpace(q.Answer), "; ")
		correct[q.ID] = answers
	}

	// 填空题
	for _, q := range ged.FillBlank {
		var blanks [][]string
		for _, b := range q.Blanks {
			blanks = append(blanks, b.Answers)
		}
		correct[q.ID] = blanks
	}

	return correct
}

// GetAllCorrectAnswers 获取所有题型的正确答案map（修复版）
func (es *ExamState) GetAllCorrectAnswers() map[string]interface{} {
	correct := make(map[string]interface{})

	// 优先使用 GroupedExamData，但只用于有答案的题目
	if es.CorrectAnswers != nil {
		fmt.Println("使用 GroupedExamData 获取正确答案")
		groupedCorrect := es.CorrectAnswers.GetCorrectAnswers()
		// 复制 GroupedExamData 的答案
		for id, ans := range groupedCorrect {
			correct[id] = ans
		}
	}

	// 总是使用备用方案补充，确保所有题目都有答案
	fmt.Println("使用备用方案补充正确答案")

	// 备用方案：直接从原始题目数据提取所有题型的正确答案
	for _, q := range es.Questions {
		// 如果已经从 GroupedExamData 获取了答案，跳过
		if _, exists := correct[q.ID]; exists {
			continue
		}

		// 合并处理同类型题目
		switch {
		case q.Type == "DR":
			// 材料阅读题没有答案
			correct[q.ID] = nil
		case q.Type == "SC" || q.Type == "SCIMG":
			// 单选题：Answer是string，如"A"
			if ansStr, ok := q.Answer.(string); ok && ansStr != "" {
				// 直接存储选项索引，不转换为内容
				correct[q.ID] = []string{ansStr}
			}
		case q.Type == "MC" || q.Type == "MCIMG":
			// 多选题：Answers是[]interface{}，每个是string如"A;B;C"
			var answers []string
			for _, ansInterface := range q.Answers {
				if ans, ok := ansInterface.(string); ok && ans != "" {
					// 直接存储选项索引
					answers = append(answers, ans)
				}
			}
			correct[q.ID] = answers
		case q.Type == "FL" || q.Type == "FLIMG":
			// 填空题：Answers是[]interface{}，每个是AnswerItem
			var blanks [][]string

			// 使用map按blankIndex收集答案
			blankMap := make(map[int][]string)
			var blankIndexes []int

			for _, ansInterface := range q.Answers {
				if ansItem, ok := ansInterface.(map[string]interface{}); ok {
					if blankIdx, ok := ansItem["blankIndex"].(float64); ok {
						idx := int(blankIdx)
						if !contains(blankIndexes, idx) {
							blankIndexes = append(blankIndexes, idx)
						}
						if answers, exists := ansItem["answers"]; exists {
							if ansList, ok := answers.([]interface{}); ok {
								for _, a := range ansList {
									if str, ok := a.(string); ok {
										blankMap[idx] = append(blankMap[idx], str)
									}
								}
							}
						}
					}
				}
			}

			// 按blankIndex排序
			sort.Ints(blankIndexes)
			for _, idx := range blankIndexes {
				blanks = append(blanks, blankMap[idx])
			}

			// 确保空白数量与BlankCount一致
			for len(blanks) < q.BlankCount {
				blanks = append(blanks, []string{})
			}

			correct[q.ID] = blanks
		default:
			fmt.Printf("警告: 未知题目类型 %s for question %s\n", q.Type, q.ID)
		}
	}

	return correct
}

// contains 检查slice是否包含元素
func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// SaveExamResult 保存考试结果（包含错题集）
func (es *ExamState) SaveExamResult(gradesPath string, originalData *ExamData, examTempPath string) error {
	// 创建成绩目录
	if err := os.MkdirAll(gradesPath, 0755); err != nil {
		return fmt.Errorf("创建成绩目录失败: %v", err)
	}

	// 计算分数
	finalScore, totalScore, correctCount, wrongCount := es.CalculateScore()

	// 使用统一的时间格式
	examTime := time.Now().Format("20060102_150405")

	// 创建专属考试文件夹
	examFolder := fmt.Sprintf("%s_%s", es.ExamName, examTime)
	fullExamPath := filepath.Join(gradesPath, examFolder)
	if err := os.MkdirAll(fullExamPath, 0755); err != nil {
		return fmt.Errorf("创建考试文件夹失败: %v", err)
	}

	// 保存成绩文件
	gradeData := map[string]interface{}{
		"examName":     es.ExamName,
		"finalScore":   finalScore,
		"totalScore":   totalScore,
		"correctCount": correctCount,
		"wrongCount":   wrongCount,
		"startTime":    es.StartTime.Format("2006-01-02 15:04:05"),
		"endTime":      time.Now().Format("2006-01-02 15:04:05"),
		"duration":     time.Since(es.StartTime).String(),
	}

	gradeFileName := fmt.Sprintf("%s_%s.json", es.ExamName, examTime)
	gradeFilePath := filepath.Join(fullExamPath, gradeFileName)

	file, err := os.Create(gradeFilePath)
	if err != nil {
		return fmt.Errorf("创建成绩文件失败: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(gradeData); err != nil {
		return fmt.Errorf("编码成绩文件失败: %v", err)
	}

	// 同步生成错题集
	if err := CreateWrongQuestionCollection(es, originalData, examTempPath); err != nil {
		fmt.Printf("生成错题集失败: %v\n", err)
		return fmt.Errorf("生成错题集失败: %v", err)
	}

	return nil
}
