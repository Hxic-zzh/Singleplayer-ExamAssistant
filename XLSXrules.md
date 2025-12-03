# Excel 文件数据结构文档
> 一定要看完口牙
> English is below, please check it out
> <img src="https://github.com/Hxic-zzh/Singleplayer-ExamAssistant/blob/bec8e23c2bd8d10a34c774f222072644533a0f66/githubPng/8.png" width="400">

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
