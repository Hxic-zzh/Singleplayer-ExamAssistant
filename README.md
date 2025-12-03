# Singleplayer-ExamAssistant
A cross-platform exam software built with Go and Fyne. It manages question banks, tests, scores, and wrong-answer reviews. Ideal for education and self-study, it features a clean interface, multi-language support, and handles diverse question types like images and reading materials.

**Read this in other languages: [English](README.md), [中文](README_zh.md).**


# User Manual

## Prerequisites

> **Important Note**: The author is a complete beginner and has not fully adapted the software for Windows. Additionally, I have not used Linux and macOS, so I don't know how to modify the software's I/O for those systems. My apologies. If you have any issues, please submit them via GitHub Issues.

---

## Software Introduction

> A question bank examination and management software developed based on Go language and Fyne framework. Supports multiple question types (single choice, multiple choice, fill-in-the-blank, reading comprehension, picture-based questions, etc.), multi-language interfaces, score and mistake tracking, question bank import/export functionality, etc.
> Suitable for schools, training institutions, personal self-study, and other scenarios.

---

## Launch Methods

**1. Launch by clicking `Hxic-Windows64-GUI.exe` or `Hxic-Windows64-Console.exe`**
2. `-Console` version includes a console window but output information is not fully in English; use with caution for debugging.

## Main Interface Navigation

The left side of the software is the menu bar; click to switch between different functional pages:

### 1. Announcements Page
- Displays welcome messages, announcements, carousel images, etc.

<img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/610ad71b7b2d4f2efa777f25aa1e8abaee175a81/githubPng/1.png" width="500">

### 2. Online Exam
- Enter the main exam interface, automatically loads question bank
- Supports multiple question types: single choice, multiple choice, fill-in-the-blank, reading comprehension, picture-based questions, etc.
- Supports images and material content display for questions
- Features question navigation, pagination, countdown timer, etc.
- Auto-saves answers, supports submitting exams and automatic scoring
- Automatically generates mistake collection after submission
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/2.png" width="500">
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/3.png" width="500">

### 3. Question Bank Extraction and Management
- Supports batch import of question banks from Excel format
- Supports question bank name editing and saving
- Extracted question banks need to be manually exported
- Preview question bank content, supports clearing current editing progress
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/4.png" width="500">

### 4. Score and Mistake Management
- Automatically saves each exam score, supports score list browsing
- Supports mistake collection generation and review for targeted practice
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/5.png" width="500">

### 5. Settings Center
- Supports multi-language switching, currently only Chinese/English
- Supports theme switching, interface customization, etc.
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/8357210025b07f5289e19985edf9aa235077cb18/githubPng/6.png" width="500">

---

## Detailed Usage Instructions

### 1. First, Learn How to Create Question Banks
- In the software directory: `data/example/`, you'll find 8 XLSX files. These are "Standard Question Bank Example Files" with pre-configured formats I've set up.
- <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/1116ae326f0b9611f8235f366f8c666a9f50552d/githubPng/7.png" width="500">
- [📖 XLSX Read/Write Manual](https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/5427085906d8a98d54ed0f4b19e3d4113ec5ff62/XLSXrules.md){:target="_blank"} ↗
- Additional notes:
> Bro, check this out

<img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/237b2728be8633f63de4e3392b752229e95461cb/githubPng/10.gif" width="200">

  - Question stem images and supplementary images cannot have duplicate names within the same file.
  - Within your current upload file group, also avoid duplicate names.
  - Different questions wanting to use the same image cannot use the same file; must use different image files with renamed copies.

### 2. Learn How to Create Question Bank ZIP Files for Exam Use
- Complete this task in the "Extract" interface on the right side list of the main interface:
  1. First select the location where your XLSX files and image files are stored. Remember: place all images in a folder named "add" (exactly "add").
  2. Add main files, including `SC`, `MC`, `FL`, `SCIMG`, `MCIMG`, `FLIMG`. Click "Open main file" to upload.
  > Note: Don't upload extra files.
  3. Add auxiliary files, only `FE`. Click "Open auxiliary file" to upload.
  4. In the input box on the right top, enter the question bank name (can be numbers, no special characters); click "save question bank name" to save.
  5. To add images: simply compress your "add" folder into a ZIP file, then click "Import images ZIP" to upload your ZIP. Success will show a prompt.
  6. After all uploads are complete, click the "Generate preview" button to start generating preview.
  > If the question bank has over 20 questions, it may lag for a while; please be patient.
  7. You can see preview results in: JSON preview (bottom left), MD preview (top right), and config notifications (bottom right).
  > These three interfaces are previews. Generally, modifying content here won't cause final question bank errors, but still be careful because I'm a "character-oriented programming" beginner.
  8. After everything is ready, click "Save question bank" and wait for packaging. The question bank file will be generated in `/data/output` directory; remember to check.
  9. To use this question bank for exams, place this ZIP file in the `/data/Question` directory. Note: Don't have duplicate ZIP filenames, bro.
  10. **SPECIAL ATTENTION!!!** Some XLSX read/write software (WPS for Chinese users, and Office for international users) may recognize dates like "1900/0/0" format when writing year/month/day, causing Go system recognition failure and garbled text. In the generated ZIP file, manually locate the question position in the JSON file and manually modify answers and question stems if needed. My sincere apologies.

  <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/188aed76e8a186c4e56574ea6d23b8270fcc5df7/githubPng/11.gif" width="200">

### 3. Start Online Exam
- Complete this task in the "Exam" interface on the right side list of the main interface:
  1. Select your question bank from the dropdown menu.
  > If you can't find your question bank, check if it's in the `/data/Question` directory.
  > If still not found, refresh the software by clicking other items in the menu and switching back.
  2. Click the "machine_start_exam" button to start the exam. Wait for the progress bar.
  3. Once inside, it's a normal online exam interface. Some control translations might be inaccurate, bro. Try it out first before answering questions.

### 4. Review Mistakes
- Complete this task in the "Detail" interface on the right side list of the main interface.
- Simple and intuitive; just look at the interface.
- Mistake files are in the `/data/grades/` directory. If they can't be deleted, delete manually.

### 5. Settings
- Adjust the "exam paper interface" for online exams according to your computer's screen resolution.
- Adjust language.

---

## Important Notes

- Ensure complete directory structure for software operation.
- Multi-language text configuration is in `data/Language.json`.
- Exam and score data are automatically saved in `data/grades/` directory.
- Generated question bank files default to `data/output/` directory.
- Question banks for use should be placed in `data/Question` directory.

---

If you have more questions, contact the developer. If enough users emerge, the author will open a communication group.
