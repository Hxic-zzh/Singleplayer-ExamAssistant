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
	"time"
)

// SetConfig 设置配置结构体
type SetConfig struct {
	Resolution string `json:"resolution"` // 分辨率模式："2560x1440" 或 "1920x1080"
	Language   string `json:"language"`   // 语言："zh" 或 "en"
	CardWidth  int    `json:"card_width"` // 卡片宽度
	LastUpdate string `json:"lastUpdate"` // 最后更新时间
}

// ResolutionConfig 分辨率配置结构体
type ResolutionConfig struct {
	Name          string
	Width         int
	Height        int
	DividerX      float32 // 分割线X位置
	CardWidth     float32 // 卡片宽度
	CardHeight    float32 // 卡片高度
	ScrollWidth   float32 // 滚动区域宽度
	ScrollHeight  float32 // 滚动区域高度
	NavBtnWidth   float32 // 导航按钮宽度
	NavBtnHeight  float32 // 导航按钮高度
	BtnSpacingX   float32 // 按钮水平间距
	BtnSpacingY   float32 // 按钮垂直间距
	CardSpacing   float32 // 卡片间距
	ImageDisplayX float32 // 图片显示区域X位置
	ImageDisplayY float32 // 图片显示区域Y位置
	ImageSize     float32 // 图片尺寸
}

// 分辨率配置映射
var ResolutionConfigs = map[string]ResolutionConfig{
	"2560x1440": {
		Name:          "2560x1440",
		Width:         2560,
		Height:        1440,
		DividerX:      500,
		CardWidth:     980,
		CardHeight:    280,
		ScrollWidth:   1000,
		ScrollHeight:  900,
		NavBtnWidth:   80,
		NavBtnHeight:  40,
		BtnSpacingX:   85,
		BtnSpacingY:   45,
		CardSpacing:   320,
		ImageDisplayX: 25,
		ImageDisplayY: 440,
		ImageSize:     256,
	},
	"1920x1080": {
		Name:          "1920x1080",
		Width:         1920,
		Height:        1080,
		DividerX:      400,
		CardWidth:     700,
		CardHeight:    250,
		ScrollWidth:   750,
		ScrollHeight:  700,
		NavBtnWidth:   70,
		NavBtnHeight:  35,
		BtnSpacingX:   75,
		BtnSpacingY:   40,
		CardSpacing:   280,
		ImageDisplayX: 25,
		ImageDisplayY: 440,
		ImageSize:     256,
	},
}

// GetCurrentResolutionConfig 获取当前分辨率配置
func GetCurrentResolutionConfig() ResolutionConfig {
	config, err := LoadSetConfig()
	if err != nil {
		// 默认使用2560x1440
		return ResolutionConfigs["2560x1440"]
	}

	if resConfig, exists := ResolutionConfigs[config.Resolution]; exists {
		// 使用设置中的 CardWidth
		resConfig.CardWidth = float32(config.CardWidth)
		return resConfig
	}

	return ResolutionConfigs["2560x1440"]
}

// SaveSetConfig 保存设置到文件
func SaveSetConfig(config *SetConfig) error {
	config.LastUpdate = time.Now().Format("2006-01-02 15:04:05")

	// 确保data目录存在
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("创建data目录失败: %v", err)
	}

	filePath := filepath.Join(dataDir, "SetConfig.json")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建设置文件失败: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("写入设置文件失败: %v", err)
	}

	return nil
}

// LoadSetConfig 从文件加载设置
func LoadSetConfig() (*SetConfig, error) {
	filePath := filepath.Join("data", "SetConfig.json")
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，返回默认配置
			return GetDefaultSetConfig(), nil
		}
		return nil, fmt.Errorf("打开设置文件失败: %v", err)
	}
	defer file.Close()

	var config SetConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析设置文件失败: %v", err)
	}

	return &config, nil
}

// GetDefaultSetConfig 获取默认设置
func GetDefaultSetConfig() *SetConfig {
	return &SetConfig{
		Resolution: "2560x1440",
		Language:   "zh",
		CardWidth:  GetDefaultCardWidth("2560x1440"),
		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// GetDefaultCardWidth 根据分辨率获取默认卡片宽度
func GetDefaultCardWidth(resolution string) int {
	switch resolution {
	case "2560x1440":
		return 980
	case "1920x1080":
		return 735 // 按比例缩放 980 * (1920/2560) ≈ 735
	default:
		return 980
	}
}
