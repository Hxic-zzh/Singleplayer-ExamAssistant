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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type TempData struct { // 定义的结构都放上面，找起来方便
	SelectedFolder string `json:"selectedFolder"` // 这后面的叫做“标签映射”，因为Go里面喜欢用大写开头，但是json一般用小写。只是告诉系统，在json里面是叫做别的
}

// === 打开文件夹主程序 ===
func SelectFolder(window fyne.Window, callback func(string)) { // 所谓回调函数，就是我将部分子功能放在外面，然后用回调函数接入
	dialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) { // NewFolderOpen是一个方法，所以要自赋值咯
		if err != nil { // error 类型就是专门用来保存错误返回值的
			fmt.Println("选择文件夹出错:", err)
			return
		}
		if uri == nil {
			return
		}

		selectedPath := uri.Path()
		fmt.Println("选择的文件夹:", selectedPath)

		if err := savePathToJSON(selectedPath); err != nil { // 上面都通过了下面这个err当然还能接着用
			fmt.Println("保存路径到 JSON 出错:", err)
			return
		}

		if callback != nil {
			callback(selectedPath)
		}
	}, window)

	dialog.Show()
}

// === data/tempData.json保存方案 ===
func savePathToJSON(folderPath string) error {
	data := TempData{
		SelectedFolder: folderPath, // 赋值
	}

	// 转换为 JSON
	jsonData, err := json.MarshalIndent(data, "", "  ") // 带缩进的转换，第三个变量是缩进；第二个变量是前缀
	if err != nil {
		return fmt.Errorf("JSON 编码失败: %v", err)
	}

	// 确保 data 目录存在
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil { // 0755表示“完全控制” ，错误报告为空就是没有问题
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	filePath := filepath.Join(dataDir, "tempData.json") // 我这里突然有个领悟：Go语言 里面第一次创建变量 总是喜欢直接就连带上之后的操作
	err = os.WriteFile(filePath, jsonData, 0644)        // 只允许该go读写文件；用err的时候还是要小心，一个函数里面经常用它做错误判断，难免会有错漏
	if err != nil {                                     // 提醒一下：是先完成方法再赋值的，所以err当然是可以是空的，这是Go的一种底层逻辑，不用怀疑
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Println("路径已保存到:", filePath)
	return nil
}

// === 从 JSON 文件加载数据 ===
func LoadTempData() (*TempData, error) {
	filePath := filepath.Join("data", "tempData.json")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) { // "_"意味着“不需要文件的详细信息”，就检查路径下的文件有没有 ； IsNotExist用上一步的err判断是否是"文件不存在"错误
		return &TempData{}, nil // 这里还是要记一下，返回结构体是 xxx{}，是指针的话是 &xxx{}
	}

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析 JSON
	var tempData TempData
	if err := json.Unmarshal(data, &tempData); err != nil { // Unmarshal用来将json内容转换成Go内容
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	return &tempData, nil
}

// === 智能路径截断 === AI写的，我没仔细看，算法就懒得看了，能用就行，复杂度高也没事，就计算一次而已
func TruncatePathSmart(path string, maxLength int) string {
	if len(path) <= maxLength || path == "" {
		return path
	}

	dir, file := filepath.Split(path) // Split专门用来做路径切割的，go还挺贴心的
	if dir == "" {                    // 如果本身就是文件名，Split不会切割，Split主要是识别"/"有没有
		if len(file) > maxLength {
			return file[:maxLength-3] + "..."
		}
		return file
	}

	if file != "" {
		if len(file) >= maxLength-3 { // 如果文件夹名本身就很长
			return "..." + file[len(file)-(maxLength-3):]
		}

		// 计算剩余可用长度
		remaining := maxLength - len(file) - 3 // -3 是为了 "..."
		if remaining > 5 {
			// 从目录开头取一部分
			start := dir
			if len(start) > remaining {
				start = start[:remaining]
			}
			return start + "..." + file
		} else {
			return "..." + file
		}
	}

	// 如果是纯目录路径（以分隔符结尾）
	if len(dir) > maxLength {
		// 保留开头和结尾
		keepEachSide := (maxLength - 3) / 2
		if keepEachSide < 3 {
			return dir[:maxLength-3] + "..."
		}
		start := dir[:keepEachSide]
		end := dir[len(dir)-keepEachSide:]
		return start + "..." + end
	}

	return path
}
