
> 前提说明：作者是纯菜鸟，还没有对windows做全面适配，然后linux和mac没用过，不知道怎么修改软件的IO，见谅见谅。有啥问题issues
# 使用说明

## 软件简介

> 基于 Go 语言和 Fyne 框架开发的题库考试与管理软件，支持多种题型（单选、多选、填空、材料阅读、图片题等）、多语言界面、成绩与错题管理、题库导入导出等功能
> 适用于学校、培训机构、个人自学等场景

---

## 启动方法

**1. 点击Hxic-Windows64-GUI.exe或者Hxic-Windows64-Console.exe启动**
2. -Console是携带控制带的版本但是输出信息不全是英文，谨慎调试

## 主界面功能导航

软件左侧为菜单栏，点击可切换不同功能页面：

### 1. 消息公告页
- 展示欢迎信息、公告、轮播图片等内容

<img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/610ad71b7b2d4f2efa777f25aa1e8abaee175a81/githubPng/1.png" width="500">

### 2. 上机考试
- 进入考试主界面，自动加载题库
- 支持单选、多选、填空、材料阅读、图片题等多种题型
- 题目支持图片、材料内容展示
- 支持题目导航、分页、倒计时等功能
- 答题自动保存，支持交卷与自动判分
- 交卷后自动生成错题集
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/2.png" width="500">
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/3.png" width="500">

### 3. 题库提取与管理
- 支持从 Excel 格式批量导入题库
- 支持题库名称编辑、保存
- 获取的题库需要手动导出
- 可预览题库内容，支持清空当前编辑进度
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/4.png" width="500">

### 4. 成绩与错题管理
- 自动保存每次考试成绩，支持成绩列表浏览
- 支持错题集生成与回顾，便于查漏补缺
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/5.png" width="500">

### 5. 设置中心
- 支持多语言切换，目前只支持中英切换
- 支持主题切换、界面自定义等功能
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/6.png" width="500">

---

## 详细使用方法

### 1. 首先了解如何写题库
- 在 软件文件：`data/excample/` 目录下 你会看见8个XLSX文件，这8个文件是“标准题库示例文件”，是我预先设置好的格式文件
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/1116ae326f0b9611f8235f366f8c666a9f50552d/githubPng/7.png" width="500">
- [📖 XLSX读写手册](https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/5427085906d8a98d54ed0f4b19e3d4113ec5ff62/XLSXrules.md) ↗
- 补充说明一下：
> bro，看一下

<img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/237b2728be8633f63de4e3392b752229e95461cb/githubPng/10.gif" width="200">

  - 题干图片和附带图片在同一个文件里面不能重名。
  - 在你这次要上传的文件群内，也不要重名。
  - 不同题目要使用同一张图片也不允许使用同一个文件，必须使用不同图片，重新命名。

### 2. 了解怎么创建题库压缩包，供上机使用
- 在 主界面，右侧的列表中 “Extact” 界面 中完成我们的任务
  1. 先去选择你的xlsx文件和图片文件放置的地方，记住，将所有图片放在add文件夹里面，就叫做“add”
  2. 添加主文件，包括 `SC` `MC` `FL` `SCIMG` `MCIMG` `FLIMG` 点击“Open main file”上传 
  > 注意不要多传
  3. 添加辅助文件，只有 `FE` 点击“Open auxiliary file”上传
  4. 右侧上方输入框中，填写题库名称，可以是数字，不要有特殊符号；点击“save question bank name”保存
  5. 若是要加图片，很简单，先将你的add文件夹打包，然后点击“Import imgames ZIP”上传你的ZIP就行，成功会跳提示
  6. 全部上传完成以后，点击“Gererate preview”按钮，开始生成预览
  > 若是题库超过20题会卡一段时间，请耐心等待
  7. 在左下的json预览，右上的md预览，右下的config通知中，可以看到预览的结果
  > 这三个界面是预览，一般来说修改里面的内容不会导致最终题库错误，但是还是要小心，因为我是“面对字符编程”的编程菜鸟
  8. 全部准备好以后，点击“Save question bank”等待打包，会生成题库文件在 `/data/output` 内，记得注意一下
  9. 然后要使用该题库上机，就将这个ZIP放在 `/data/Question` 目录下面，注意不要出现相同的ZIP文件名，bro
  10. 特别注意！！！部分XLSX读写软件（中国用户的WPS，以及外国游人的office）在写年月日时候，xlsx文件会识别成 例如“1900/0/0”的形式，会导致Go系统识别失败乱码，注意在生成的ZIP文件中，自行定位Json文件中，该题的位置，然后手动修改答案和题干，实在抱歉

### 3. 开始上机
- 在 主界面，右侧的列表中 “Exam” 界面 中完成我们的任务
  1. 在下拉菜单中选择你的题库
  > 若是没找到你的题库，检查一下是不是在 `/data/Question` 目录下面
  > 若是还没找到就刷新一下软件，刷新方法是点击 菜单 里面的其他项，再切换回来
  2. 点击“machine_start_exam”按钮，开始考试，等一下进度条
  3. 进去以后就是正常的上机考试，可能有几个控件翻译不太对，哥们先提前尝试一下，不要作题

### 4. 检查错题
- 在 主界面，右侧的列表中 “Detail” 界面 中完成我们的任务
- 简单易懂，自己看界面
- 错题文件在 `/data/grades/` 这个目录下，删不掉可以手动去删

### 5. 设置
- 按照自己电脑的分辨率调整上机考试的“试卷界面”
- 调整语言


---

## 注意事项

- 软件运行需保证目录结构完整
- 多语言文本配置在 `data/Language.json`
- 考试及成绩数据默认保存在 `data/grades/` 目录下
- 生成的题库文件默认在`data/output/`目录下
- 需要使用的题库在`data/Question`目录下

---

如有更多问题，请联系开发者，或者用的人多了以后，作者会开交流群
