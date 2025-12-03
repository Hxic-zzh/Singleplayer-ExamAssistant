# Excel 文件数据结构文档
> 一定要看完口牙

> English is below, please check it out

> <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/bec8e23c2bd8d10a34c774f222072644533a0f66/githubPng/8.png" width="200">

**[点这里看中文](#中文文档)** | **[Click here for English](#英文文档)**

<hr>

<a id="中文文档"></a>
```
 $$$$$$\  $$\       $$\                                         
$$  __$$\ $$ |      \__|                                        
$$ /  \__|$$$$$$$\  $$\ $$$$$$$\   $$$$$$\   $$$$$$$\  $$$$$$\  
$$ |      $$  __$$\ $$ |$$  __$$\ $$  __$$\ $$  _____|$$  __$$\ 
$$ |      $$ |  $$ |$$ |$$ |  $$ |$$$$$$$$ |\$$$$$$\  $$$$$$$$ |
$$ |  $$\ $$ |  $$ |$$ |$$ |  $$ |$$   ____| \____$$\ $$   ____|
\$$$$$$  |$$ |  $$ |$$ |$$ |  $$ |\$$$$$$$\ $$$$$$$  |\$$$$$$$\ 
 \______/ \__|  \__|\__|\__|  \__| \_______|\_______/  \_______|  here
```
                                                                

## 📁 示例文件文件列表

| 文件名 | 文件描述 | 主要用途 |
|--------|----------|----------|
| **MultipleChoice_MC.xlsx** | 多项选择题题库 | 存储多项选择题题目、选项、答案及相关信息 |
| **SingleChoice_SC.xlsx** | 单项选择题题库 | 存储单项选择题题目、选项、答案及相关信息 |
| **Fill_FL.xlsx** | 填空题题库 | 存储填空题题目、答案及相关信息 |
| **Fill_FE.xlsx** | 填空题答案扩展 | 存储填空题的详细答案选项 |
| **FILL2_FLIMG.xlsx** | 题干是图填空题题库 | 题干主体是图片 可以附带文字信息 |
| **SingleChoice2_SCIMG.xlsx** | 题干是图选择题题库 | 题干主体是图片 可以附带文字信息  |
| **MultipleChoice2_MCIMG.xlsx** | 题干是图多选题选择题题库 | 题干主体是图片 可以附带文字信息  |
| **DocumentReading_DR.xlsx** | 材料阅读题库 | 存储阅读材料及相关题目信息 使用钩子读取其他文件题目 |

---

## 📊 详细数据结构
> 注意：故意使用TURE,不要疑惑,源于作者懒得改
> ** 注意：XLSX文件不允许同名，不允许同名，不允许同名，重要的事情说三遍！ **
> _MC _SC _FL _FE _FLIMG _SCIMG _MCIMG _DR 为文件类型标识，必须写！具体往下看

### 1. **MultipleChoice_MC.xlsx**
**Sheet1: 多项选择题题库**
- **A列**: 是否启用 (TURE/FALSE)
- **B-L列**: 选项A-K (最多11个选项)
- **M列**: 题干描述
- **N列**: 正确答案 (多选，如ABD)
- **O-Y列**: 图片名称1-11 (对应附属图片，在选项的下面)
- **Z列**: 钩子 (题目关联标识符)

### 2. **SingleChoice_SC.xlsx**
**Sheet1: 单项选择题题库**
- **A列**: 是否启用 (TURE/FALSE)
- **B-L列**: 选项A-K (最多11个选项)
- **M列**: 题干描述
- **N列**: 正确答案 (单选，如A/B/C)
- **O-Y列**: 图片名称1-11 (对应题干图片)
- **Z列**: 钩子 (题目关联标识符)

### 3. **Fill_FL.xlsx**
**Sheet1: 填空题题库**
- **A列**: 是否启用 (TURE/FALSE)
- **B-U列**: 答案1-20 (对应填空题的各个空位)
- **V列**: 题干描述 (包含`(%___%)`作为填空标志符,同为占位识别符)
- **W列**: 序号（在FE文件中使用，用于确定FL文件内题目的哪道空需要补充）
- **X-AH列**: 图片名称1-11 (对应题干图片)
- **AI列**: 钩子 (题目关联标识符)

### 4. **Fill_FE.xlsx**
**Sheet1: 填空题答案扩展**
- **A列**: 是否启用 (TURE/FALSE)
- **B列**: 序号
- **C列**: 第几空
- **D-M列**: 答案1-10 (每个空位可接受的多个答案)

### 5. **FILL2_FLIMG.xlsx**
**Sheet1: 图片填空题题库**
- **A列**: 是否启用 (TURE/FALSE)
- **B-U列**: 答案1-20
- **V-Y列**: 题干图片1-4
- **Z列**: 题干描述
- **AA列**: 钩子 (题目关联标识符)

### 6. **SingleChoice2_SCIMG.xlsx**
**Sheet1: 图片单项选择题题库**
- **A列**: 是否启用 (TURE/FALSE)
- **B-K列**: 选项A-J
- **L-O列**: 题干图片1-4
- **P列**: 题干描述
- **Q列**: 正确答案
- **R列**: 钩子 (题目关联标识符)

### 7. **MultipleChoice2_MCIMG.xlsx**
**Sheet1: 图片多项选择题题库**
- **A列**: 是否启用 (TURE/FALSE)
- **B-L列**: 选项A-K
- **K-P列**: 题干图片1-4
- **Q列**: 题干描述
- **R列**: 正确答案 (多选)
- **S列**: 钩子 (题目关联标识符)

### 8. **DocumentReading_DR.xlsx**
**Sheet1: 材料阅读题库**
- **A列**: 是否启用 (TURE/FALSE)
- **B-E列**: 资料1-4
- **F-I列**: 题干图片1-4
- **J-AM列**: 所属题1-30 (题目关联标识符，如SC.A1, MC.C1等)

---

## 🔗 关联关系说明

1. **钩子(Hook)系统**:
   - 用于题目之间的关联
   - 格式示例: `SC.A1`, `MC.C1`, `FL.E1` 等
   - 前缀表示题目类型，后缀为唯一标识
   - 不允许在多个DR文件中使用相同的钩子标识！！！
   - 例如：我有一个XXX_DR.xlsx,还有一个YYY_DR.xlsx，我有一道题一样的，但是我也不能使用相同的钩子标志，必须一个是SC.A1,另一个是SC.A2。差不多这个意思

2. **图片命名规范**:
   - 所有文件共享相同的图片资源
   - 图片名称如: `test1`, `testF1`, `test4` 等
   - 对应实际的图片文件
   - 文件格式支持 png jpg webp

3. **题目类型标识**:
   - `SC`: 单项选择题
   - `MC`: 多项选择题  
   - `FL`: 填空题
   - `IMG`: 图片题
  
4. **重要文件须知**：
   - 不允许出现 ABA_SC.xlsx 和 ABA_MC.xlsx 这样文件类型不同，但是文件名一样的文件。系统会给你报错，没报错的话就自己掂量掂量，别瞎搞
   - Fill_FL.xlsx 与 Fill_FE.xlsx 一一对应，是一次函数，逻辑是 FILL题库本身就是有两个文件，一个是主文件FL，一个数辅助文件FE
   - 所以 FILL_FL.xlsx 用不了 FILLXX_FE.xlsx ，FILL_FL.xlsx 也不会接受两个_FE的补充
   - 详细的看看我的示例文件，命名时有规则的，仔细看看能看懂
   - 然后 FILL2_FLIMG.xlsx 文件没有_FE补充，我的设想是图片是题干，那就是看图题，答案一般情况下也是唯一的
   - 还有最好用“Spanish Naming Convention”这种命名方法，不知道的去查一下，简单的
  
5. **题目展示顺序**：材料->单选->多选->填空（题干是图的在各个层次的前面）

---


## 🏷️ 文件命名约定

| 前缀 | 含义 | 示例 |
|------|------|------|
| `_MC` | 多项选择题 | MultipleChoice_MC.xlsx |
| `_SC` | 单项选择题 | SingleChoice_SC.xlsx |
| `_FL` | 填空题 | Fill_FL.xlsx |
| `_FE` | 填空题扩展 | Fill_FE.xlsx |
| `_IMG` | 图片题 | SingleChoice2_SCIMG.xlsx |
| `_DR` | 材料阅读 | DocumentReading_DR.xlsx |

这个结构支持一个完整的考试系统，包含多种题型、图片支持、题目关联等功能。


<hr>

<a id="英文文档"></a>

```
$$$$$$$$\                     $$\ $$\           $$\       
$$  _____|                    $$ |\__|          $$ |      
$$ |      $$$$$$$\   $$$$$$\  $$ |$$\  $$$$$$$\ $$$$$$$\  
$$$$$\    $$  __$$\ $$  __$$\ $$ |$$ |$$  _____|$$  __$$\ 
$$  __|   $$ |  $$ |$$ /  $$ |$$ |$$ |\$$$$$$\  $$ |  $$ |
$$ |      $$ |  $$ |$$ |  $$ |$$ |$$ | \____$$\ $$ |  $$ |
$$$$$$$$\ $$ |  $$ |\$$$$$$$ |$$ |$$ |$$$$$$$  |$$ |  $$ |
\________|\__|  \__| \____$$ |\__|\__|\_______/ \__|  \__|
                    $$\   $$ |                            
                    \$$$$$$  |                            
                     \______/                              here
```

# 📁 Sample File List

| File Name | File Description | Main Purpose |
|-----------|------------------|--------------|
| **MultipleChoice_MC.xlsx** | Multiple Choice Question Bank | Stores multiple-choice questions, options, answers, and related information |
| **SingleChoice_SC.xlsx** | Single Choice Question Bank | Stores single-choice questions, options, answers, and related information |
| **Fill_FL.xlsx** | Fill-in-the-Blank Question Bank | Stores fill-in-the-blank questions, answers, and related information |
| **Fill_FE.xlsx** | Fill-in-the-Blank Answer Extension | Stores detailed answer options for fill-in-the-blank questions |
| **FILL2_FLIMG.xlsx** | Image-Based Fill-in-the-Blank Question Bank | The question stem is primarily an image, can include text information |
| **SingleChoice2_SCIMG.xlsx** | Image-Based Single Choice Question Bank | The question stem is primarily an image, can include text information |
| **MultipleChoice2_MCIMG.xlsx** | Image-Based Multiple Choice Question Bank | The question stem is primarily an image, can include text information |
| **DocumentReading_DR.xlsx** | Reading Material Question Bank | Stores reading materials and related question information, uses hooks to reference questions from other files |

---

# 📊 Detailed Data Structure
> Note: Deliberately uses "TURE" instead of "TRUE", don't be confused, the author was too lazy to change it
> ** Important: XLSX files cannot have duplicate names, cannot have duplicate names, cannot have duplicate names, important things said three times! **
> _MC _SC _FL _FE _FLIMG _SCIMG _MCIMG _DR are file type identifiers, must be included! See details below

## 1. **MultipleChoice_MC.xlsx**
**Sheet1: Multiple Choice Question Bank**
- **Column A**: Enabled Status (TURE/FALSE)
- **Columns B-L**: Options A-K (up to 11 options)
- **Column M**: Question Stem Description
- **Column N**: Correct Answer (multiple selection, e.g., ABD)
- **Columns O-Y**: Image Names 1-11 (corresponding to supplementary images, displayed below options)
- **Column Z**: Hook (question association identifier)

## 2. **SingleChoice_SC.xlsx**
**Sheet1: Single Choice Question Bank**
- **Column A**: Enabled Status (TURE/FALSE)
- **Columns B-L**: Options A-K (up to 11 options)
- **Column M**: Question Stem Description
- **Column N**: Correct Answer (single selection, e.g., A/B/C)
- **Columns O-Y**: Image Names 1-11 (corresponding to question stem images)
- **Column Z**: Hook (question association identifier)

## 3. **Fill_FL.xlsx**
**Sheet1: Fill-in-the-Blank Question Bank**
- **Column A**: Enabled Status (TURE/FALSE)
- **Columns B-U**: Answers 1-20 (corresponding to each blank space)
- **Column V**: Question Stem Description (contains `(%___%)` as fill-in-the-blank marker and placeholder identifier)
- **Column W**: Serial Number (used in FE files to determine which blank in the FL file needs completion)
- **Columns X-AH**: Image Names 1-11 (corresponding to question stem images)
- **Column AI**: Hook (question association identifier)

## 4. **Fill_FE.xlsx**
**Sheet1: Fill-in-the-Blank Answer Extension**
- **Column A**: Enabled Status (TURE/FALSE)
- **Column B**: Serial Number
- **Column C**: Which Blank (blank position number)
- **Columns D-M**: Answers 1-10 (multiple acceptable answers for each blank)

## 5. **FILL2_FLIMG.xlsx**
**Sheet1: Image-Based Fill-in-the-Blank Question Bank**
- **Column A**: Enabled Status (TURE/FALSE)
- **Columns B-U**: Answers 1-20
- **Columns V-Y**: Question Stem Images 1-4
- **Column Z**: Question Stem Description
- **Column AA**: Hook (question association identifier)

## 6. **SingleChoice2_SCIMG.xlsx**
**Sheet1: Image-Based Single Choice Question Bank**
- **Column A**: Enabled Status (TURE/FALSE)
- **Columns B-K**: Options A-J
- **Columns L-O**: Question Stem Images 1-4
- **Column P**: Question Stem Description
- **Column Q**: Correct Answer
- **Column R**: Hook (question association identifier)

## 7. **MultipleChoice2_MCIMG.xlsx**
**Sheet1: Image-Based Multiple Choice Question Bank**
- **Column A**: Enabled Status (TURE/FALSE)
- **Columns B-L**: Options A-K
- **Columns K-P**: Question Stem Images 1-4
- **Column Q**: Question Stem Description
- **Column R**: Correct Answer (multiple selection)
- **Column S**: Hook (question association identifier)

## 8. **DocumentReading_DR.xlsx**
**Sheet1: Reading Material Question Bank**
- **Column A**: Enabled Status (TURE/FALSE)
- **Columns B-E**: Materials 1-4
- **Columns F-I**: Question Stem Images 1-4
- **Columns J-AM**: Associated Questions 1-30 (question association identifiers, e.g., SC.A1, MC.C1, etc.)

---

# 🔗 Association Relationship Explanation

## 1. **Hook System**:
   - Used for associations between questions
   - Format example: `SC.A1`, `MC.C1`, `FL.E1`, etc.
   - Prefix indicates question type, suffix is a unique identifier
   - **DO NOT use the same hook identifier in multiple DR files!!!**
   - Example: If I have XXX_DR.xlsx and YYY_DR.xlsx, and I have the same question, I still cannot use the same hook identifier; one must be SC.A1 and the other SC.A2. That's the general idea.

## 2. **Image Naming Convention**:
   - All files share the same image resources
   - Image names such as: `test1`, `testF1`, `test4`, etc.
   - Correspond to actual image files
   - Supported file formats: png, jpg, webp

## 3. **Question Type Identifiers**:
   - `SC`: Single Choice
   - `MC`: Multiple Choice
   - `FL`: Fill-in-the-Blank
   - `IMG`: Image-Based Questions

## 4. **Important File Notes**:
   - **Do not create files with the same name but different types**, e.g., ABA_SC.xlsx and ABA_MC.xlsx. The system will report an error; if it doesn't, be careful and don't mess around.
   - Fill_FL.xlsx and Fill_FE.xlsx correspond one-to-one; this is a function relationship. The logic is that the FILL question bank consists of two files: a main file (FL) and an auxiliary file (FE).
   - Therefore, Fill_FL.xlsx **cannot use** FillXX_FE.xlsx, and Fill_FL.xlsx **will not accept supplements from two _FE files**.
   - Check the example files carefully for naming rules; you should be able to understand by looking closely.
   - Additionally, the FILL2_FLIMG.xlsx file **does not have an _FE supplement**. The idea is that if the image is the question stem (a picture-based question), the answer is usually unique.
   - Also, it's best to use the "Spanish Naming Convention" for naming. If you don't know what that is, look it up; it's simple.

## 5. **Question Display Order**: Materials → Single Choice → Multiple Choice → Fill-in-the-Blank (Image-based questions appear first within each category)

---

# 🏷️ File Naming Convention

| Suffix | Meaning | Example |
|--------|---------|---------|
| `_MC` | Multiple Choice | MultipleChoice_MC.xlsx |
| `_SC` | Single Choice | SingleChoice_SC.xlsx |
| `_FL` | Fill-in-the-Blank | Fill_FL.xlsx |
| `_FE` | Fill-in-the-Blank Extension | Fill_FE.xlsx |
| `_IMG` | Image-Based Question | SingleChoice2_SCIMG.xlsx |
| `_DR` | Reading Material | DocumentReading_DR.xlsx |


