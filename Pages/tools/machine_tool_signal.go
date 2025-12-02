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
	"time"
)

type MachineExamProgress struct {
	Stage   string // 当前阶段描述
	Percent int    // 进度百分比
	Error   string // 错误信息（如有）
}

type MachineExamSignalBridge struct {
	ProgressChan chan MachineExamProgress
	DoneChan     chan bool
}

func NewMachineExamSignalBridge() *MachineExamSignalBridge {
	return &MachineExamSignalBridge{
		ProgressChan: make(chan MachineExamProgress, 10),
		DoneChan:     make(chan bool, 1),
	}
}

// 示例：发送进度
func (b *MachineExamSignalBridge) SendProgress(stage string, percent int, err string) {
	b.ProgressChan <- MachineExamProgress{Stage: stage, Percent: percent, Error: err}
}

// 示例：发送完成
func (b *MachineExamSignalBridge) SendDone() {
	b.DoneChan <- true
}

// 示例：模拟进度动画（实际动画在Animation/startLoder.go）
func (b *MachineExamSignalBridge) SimulateProgress() {
	stages := []string{"题库解压中", "试卷生成中", "错误检查中", "准备进入考试"}
	for i, s := range stages {
		time.Sleep(500 * time.Millisecond)
		b.SendProgress(s, (i+1)*25, "")
	}
	b.SendDone()
}
