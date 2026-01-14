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
)

var translations map[string]map[string]string
var currentLanguage string

// InitLocalization 初始化本地化，加载翻译文件
func InitLocalization(lang string) error {
	currentLanguage = lang

	// 构建Language.json路径
	langFile := filepath.Join("data", "Language.json")

	// 读取文件
	data, err := os.ReadFile(langFile)
	if err != nil {
		return fmt.Errorf("failed to read language file: %v", err)
	}

	// 解析JSON
	err = json.Unmarshal(data, &translations)
	if err != nil {
		return fmt.Errorf("failed to parse language file: %v", err)
	}

	return nil
}

// GetLocalized 获取本地化字符串
func GetLocalized(key string) string {
	if translations == nil {
		return key // 如果未初始化，返回键
	}
	if langMap, exists := translations[currentLanguage]; exists {
		if value, exists := langMap[key]; exists {
			return value
		}
	}
	// 如果找不到，返回键
	return key
}

// SetLanguage 设置当前语言
func SetLanguage(lang string) {
	currentLanguage = lang
}

// GetCurrentLanguage 获取当前语言
func GetCurrentLanguage() string {
	if currentLanguage != "" {
		return currentLanguage
	}
	// 尝试从SetConfig.json读取
	config, err := LoadSetConfig()
	if err == nil && config.Language != "" {
		return config.Language
	}
	return "zh"
}
