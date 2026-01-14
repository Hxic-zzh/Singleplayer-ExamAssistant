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
	"strings"
)

// 生成JSON预览
func GenerateJSONPreview(result *ParseResult, bankName string) (string, error) {
	// 构建完整的题库结构
	questionBank := map[string]interface{}{
		"name":    bankName,
		"version": "1.0",
		"metadata": map[string]interface{}{
			"totalQuestions":  len(result.SingleChoice) + len(result.MultipleChoice) + len(result.FillBlank) + len(result.DocumentReading),
			"singleChoice":    len(result.SingleChoice),
			"multipleChoice":  len(result.MultipleChoice),
			"fillBlank":       len(result.FillBlank),
			"documentReading": len(result.DocumentReading),
			"totalImages":     GetTempImageCount(), // 添加图片统计
		},
		"questions": map[string]interface{}{
			"singleChoice":    result.SingleChoice,
			"multipleChoice":  result.MultipleChoice,
			"fillBlank":       result.FillBlank,
			"documentReading": result.DocumentReading,
		},
		"errors": result.Errors,
	}

	jsonData, err := json.MarshalIndent(questionBank, "", "  ")
	if err != nil {
		return "", fmt.Errorf("生成JSON失败: %v", err)
	}

	return string(jsonData), nil
}

// 生成Markdown预览
func GenerateMarkdownPreview(result *ParseResult, bankName string) string {
	var md strings.Builder

	md.WriteString("# " + bankName + "\n\n")
	md.WriteString("## 题库概览\n\n")
	md.WriteString(fmt.Sprintf("- 总题数: %d\n", len(result.SingleChoice)+len(result.MultipleChoice)+len(result.FillBlank)+len(result.DocumentReading)))
	md.WriteString(fmt.Sprintf("- 单选题: %d\n", len(result.SingleChoice)))
	md.WriteString(fmt.Sprintf("- 多选题: %d\n", len(result.MultipleChoice)))
	md.WriteString(fmt.Sprintf("- 填空题: %d\n", len(result.FillBlank)))
	md.WriteString(fmt.Sprintf("- 材料阅读题: %d\n", len(result.DocumentReading)))
	md.WriteString(fmt.Sprintf("- 解析错误: %d\n", len(result.Errors)))
	md.WriteString(fmt.Sprintf("- 图片数量: %d\n\n", GetTempImageCount())) // 添加图片统计

	// 单选题（包括普通和题干是图）
	if len(result.SingleChoice) > 0 {
		md.WriteString("## 单选题\n\n")
		for i, q := range result.SingleChoice {
			md.WriteString(fmt.Sprintf("### 第%d题\n", i+1))

			// 显示题目类型
			if q.Type == SingleChoiceImg {
				md.WriteString("**类型**: 题干是图单选题\n")
			} else {
				md.WriteString("**类型**: 普通单选题\n")
			}

			md.WriteString(fmt.Sprintf("**题干**: %s\n\n", q.Question))

			// 显示钩子信息
			if q.Hook != "" {
				md.WriteString(fmt.Sprintf("**钩子**: %s\n\n", q.Hook))
			}

			// 显示图片信息
			if len(q.Images) > 0 {
				md.WriteString("**题干图片**:\n")
				for _, img := range q.Images {
					md.WriteString(fmt.Sprintf("- %s\n", img))
				}
				md.WriteString("\n")
			}

			md.WriteString("**选项**:\n")
			for j, opt := range q.Options {
				md.WriteString(fmt.Sprintf("- %s. %s\n", string(rune('A'+j)), opt))
			}
			md.WriteString(fmt.Sprintf("**答案**: %s\n\n", q.Answer))
		}
	}

	// 多选题（包括普通和题干是图）
	if len(result.MultipleChoice) > 0 {
		md.WriteString("## 多选题\n\n")
		for i, q := range result.MultipleChoice {
			md.WriteString(fmt.Sprintf("### 第%d题\n", i+1))

			// 显示题目类型
			if q.Type == MultipleChoiceImg {
				md.WriteString("**类型**: 题干是图多选题\n")
			} else {
				md.WriteString("**类型**: 普通多选题\n")
			}

			md.WriteString(fmt.Sprintf("**题干**: %s\n\n", q.Question))

			// 显示钩子信息
			if q.Hook != "" {
				md.WriteString(fmt.Sprintf("**钩子**: %s\n\n", q.Hook))
			}

			// 显示图片信息
			if len(q.Images) > 0 {
				md.WriteString("**题干图片**:\n")
				for _, img := range q.Images {
					md.WriteString(fmt.Sprintf("- %s\n", img))
				}
				md.WriteString("\n")
			}

			md.WriteString("**选项**:\n")
			for j, opt := range q.Options {
				md.WriteString(fmt.Sprintf("- %s. %s\n", string(rune('A'+j)), opt))
			}
			md.WriteString(fmt.Sprintf("**答案**: %s\n\n", strings.Join(q.Answers, ", ")))
		}
	}

	// 填空题（包括普通和题干是图）
	if len(result.FillBlank) > 0 {
		md.WriteString("## 填空题\n\n")
		for i, q := range result.FillBlank {
			md.WriteString(fmt.Sprintf("### 第%d题 (ID: %s)\n", i+1, q.ID))

			// 显示题目类型
			if q.Type == FillBlankImg {
				md.WriteString("**类型**: 题干是图填空题\n")
			} else {
				md.WriteString("**类型**: 普通填空题\n")
			}

			md.WriteString(fmt.Sprintf("**题干**: %s\n\n", q.Question))

			// 显示钩子信息
			if q.Hook != "" {
				md.WriteString(fmt.Sprintf("**钩子**: %s\n\n", q.Hook))
			}

			// 显示图片信息
			if len(q.Images) > 0 {
				md.WriteString("**题干图片**:\n")
				for _, img := range q.Images {
					md.WriteString(fmt.Sprintf("- %s\n", img))
				}
				md.WriteString("\n")
			}

			md.WriteString(fmt.Sprintf("**空的数量**: %d\n", q.BlankCount))
			if q.HasExtra {
				md.WriteString(fmt.Sprintf("**特殊标记**: %s\n", q.ExtraKey))
			}
			md.WriteString("**答案**:\n")
			for _, answer := range q.Answers {
				md.WriteString(fmt.Sprintf("- 第%d空: %s\n", answer.BlankIndex, strings.Join(answer.Answers, ", ")))
			}
			md.WriteString("\n")
		}
	}

	// 材料阅读题
	if len(result.DocumentReading) > 0 {
		md.WriteString("## 材料阅读题\n\n")
		for i, q := range result.DocumentReading {
			md.WriteString(fmt.Sprintf("### 第%d题\n", i+1))

			md.WriteString("**类型**: 材料阅读题\n")
			md.WriteString(fmt.Sprintf("**题干**: %s\n\n", q.Question))

			// 显示资料内容
			if len(q.Materials) > 0 {
				md.WriteString("**资料内容**:\n")
				for j, material := range q.Materials {
					md.WriteString(fmt.Sprintf("#### 资料%d\n", j+1))
					md.WriteString(fmt.Sprintf("%s\n\n", material))
				}
			}

			// 显示钩子列表
			if len(q.Hooks) > 0 {
				md.WriteString("**关联题目钩子**:\n")
				for _, hook := range q.Hooks {
					md.WriteString(fmt.Sprintf("- %s\n", hook))
				}
				md.WriteString("\n")
			}

			// 显示图片信息
			if len(q.Images) > 0 {
				md.WriteString("**题干图片**:\n")
				for _, img := range q.Images {
					md.WriteString(fmt.Sprintf("- %s\n", img))
				}
				md.WriteString("\n")
			}
		}
	}

	// 错误信息
	if len(result.Errors) > 0 {
		md.WriteString("## 解析错误\n\n")
		for _, err := range result.Errors {
			md.WriteString(fmt.Sprintf("- %s\n", err))
		}
	}

	return md.String()
}

// 生成配置预览
func GenerateConfigPreview(result *ParseResult, bankName string) string {
	var config strings.Builder

	config.WriteString("题库配置信息\n")
	config.WriteString("============\n\n")
	config.WriteString(fmt.Sprintf("题库名称: %s\n", bankName))
	config.WriteString(fmt.Sprintf("生成时间: %s\n", "自动生成"))
	config.WriteString("题目统计:\n")

	// 统计各种类型的题目数量
	scCount, scImgCount := countQuestionTypes(result.SingleChoice)
	mcCount, mcImgCount := countQuestionTypes(result.MultipleChoice)
	flCount, flImgCount := countQuestionTypesFill(result.FillBlank)

	config.WriteString(fmt.Sprintf("  - 单选题: %d 题 (普通: %d, 题干是图: %d)\n", len(result.SingleChoice), scCount, scImgCount))
	config.WriteString(fmt.Sprintf("  - 多选题: %d 题 (普通: %d, 题干是图: %d)\n", len(result.MultipleChoice), mcCount, mcImgCount))
	config.WriteString(fmt.Sprintf("  - 填空题: %d 题 (普通: %d, 题干是图: %d)\n", len(result.FillBlank), flCount, flImgCount))
	config.WriteString(fmt.Sprintf("  - 材料阅读题: %d 题\n", len(result.DocumentReading)))
	config.WriteString(fmt.Sprintf("  - 含特殊标记: %d 题\n", countSpecialFillBlanks(result.FillBlank)))
	config.WriteString(fmt.Sprintf("  - 解析错误: %d 个\n", len(result.Errors)))
	config.WriteString(fmt.Sprintf("  - 图片数量: %d 张\n", GetTempImageCount()))

	return config.String()
}

// 统计单选题和多选题的类型数量
func countQuestionTypes(questions interface{}) (int, int) {
	normalCount := 0
	imgCount := 0

	switch qs := questions.(type) {
	case []SingleChoiceQuestion:
		for _, q := range qs {
			if q.Type == SingleChoiceImg {
				imgCount++
			} else {
				normalCount++
			}
		}
	case []MultipleChoiceQuestion:
		for _, q := range qs {
			if q.Type == MultipleChoiceImg {
				imgCount++
			} else {
				normalCount++
			}
		}
	}

	return normalCount, imgCount
}

// 统计填空题的类型数量
func countQuestionTypesFill(questions []FillBlankQuestion) (int, int) {
	normalCount := 0
	imgCount := 0

	for _, q := range questions {
		if q.Type == FillBlankImg {
			imgCount++
		} else {
			normalCount++
		}
	}

	return normalCount, imgCount
}

// 统计有特殊标记的填空题
func countSpecialFillBlanks(questions []FillBlankQuestion) int {
	count := 0
	for _, q := range questions {
		if q.HasExtra {
			count++
		}
	}
	return count
}

// 更新状态信息
func UpdateStatusText(result *ParseResult, statusText *string) {
	*statusText = "解析完成!\n\n"

	// 统计各种类型的题目数量
	scCount, scImgCount := countQuestionTypes(result.SingleChoice)
	mcCount, mcImgCount := countQuestionTypes(result.MultipleChoice)
	flCount, flImgCount := countQuestionTypesFill(result.FillBlank)

	*statusText += fmt.Sprintf("✅ 单选题: %d 题 (普通: %d, 题干是图: %d)\n", len(result.SingleChoice), scCount, scImgCount)
	*statusText += fmt.Sprintf("✅ 多选题: %d 题 (普通: %d, 题干是图: %d)\n", len(result.MultipleChoice), mcCount, mcImgCount)
	*statusText += fmt.Sprintf("✅ 填空题: %d 题 (普通: %d, 题干是图: %d)\n", len(result.FillBlank), flCount, flImgCount)
	*statusText += fmt.Sprintf("📚 材料阅读题: %d 题\n", len(result.DocumentReading))
	*statusText += fmt.Sprintf("🖼️  图片数量: %d 张\n", GetTempImageCount())

	if len(result.Errors) > 0 {
		*statusText += fmt.Sprintf("\n❌ 解析错误: %d 个\n", len(result.Errors))
		for _, err := range result.Errors {
			*statusText += fmt.Sprintf("   - %s\n", err)
		}
	} else {
		*statusText += "\n🎉 所有文件解析成功!"
	}
}
