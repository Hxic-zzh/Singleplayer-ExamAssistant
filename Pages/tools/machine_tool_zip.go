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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 获取 data/Question 目录下所有 zip 文件名（不含扩展名）
func MachineToolListZips(questionDir string) ([]string, error) {
	files, err := ioutil.ReadDir(questionDir)
	if err != nil {
		return nil, err
	}
	var zips []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".zip") {
			zips = append(zips, strings.TrimSuffix(f.Name(), ".zip"))
		}
	}
	return zips, nil
}

// 清空 ExamTemp 目录
func MachineToolClearExamTemp(examTempDir string) error {
	files, err := ioutil.ReadDir(examTempDir)
	if err != nil {
		return nil // 目录不存在也算清空
	}
	for _, f := range files {
		os.RemoveAll(filepath.Join(examTempDir, f.Name()))
	}
	return nil
}
