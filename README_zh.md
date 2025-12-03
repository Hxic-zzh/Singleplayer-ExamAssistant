

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




---

## 常见问题



---

## 注意事项

- 软件运行需保证目录结构完整
- 多语言文本配置在 `data/Language.json`
- 考试及成绩数据默认保存在 `data/grades/` 目录下
- 生成的题库文件默认在`data/output/`目录下
- 需要使用的题库在`data/Question`目录下

---

如有更多问题，请联系开发者或查阅源码注释。
