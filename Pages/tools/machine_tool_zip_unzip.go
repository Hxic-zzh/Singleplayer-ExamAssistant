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
	"io"
	"os"
	"path/filepath"
)

// 解压 zip 文件到指定目录
func MachineToolUnzip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		inFile, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, inFile)
		outFile.Close()
		inFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
