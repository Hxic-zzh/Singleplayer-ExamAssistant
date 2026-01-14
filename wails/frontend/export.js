// 题库导出功能

// 从所有卡片收集数据并生成题库JSON
function collectQuestionBankData(bankName) {
    // 只使用 allQuestions，按全局顺序存储所有题目
    const allQuestions = [];

    const images = [];
    const imageSet = new Set(); // 去重
    let errors = [];

    // 优先使用 previewItems 的顺序（拖动排序后的顺序）
    // 如果没有 previewItems（旧版本），则回退到 DOM 顺序
    let cards = [];
    if (typeof window.getCardOrderedByPreview === 'function') {
        cards = window.getCardOrderedByPreview();
        console.log('📊 使用 previewItems 顺序收集数据，共', cards.length, '张卡片');
    } else {
        cards = Array.from(document.querySelectorAll('#cardList > div'));
        console.log('📊 使用 DOM 顺序收集数据，共', cards.length, '张卡片');
    }
    console.log('📊 卡片顺序:', cards.map((c, i) => {
        const label = c.querySelector('.card-type-label');
        const globalId = c.dataset.globalQuestionId || '?';
        return `${i + 1}. [全局#${globalId}]${label ? label.textContent : '未知'}`;
    }).join(', '));
    
    // 题型计数器（不再需要 globalQuestionId，因为从卡片读取）
    const counters = { SC: 0, SCIMG: 0, MC: 0, MCIMG: 0, FL: 0, FLIMG: 0, DR: 0 };
    cards.forEach(card => {
        try {
            const typeLabel = card.querySelector('.card-type-label');
            if (!typeLabel) return;

            const classList = typeLabel.classList;
            let questionData = null;

            // 单选题
            if (classList.contains('sc') || classList.contains('scimg')) {
                const type = classList.contains('scimg') ? 'SCIMG' : 'SC';
                counters[type]++;
                questionData = collectSingleChoiceData(card, type, counters[type]);
                if (questionData) {
                    questionData.globalId = parseInt(card.dataset.globalQuestionId) || 0;
                    allQuestions.push(questionData);  // 只添加到 allQuestions
                }
            }
            // 多选题
            else if (classList.contains('mc') || classList.contains('mcimg')) {
                const type = classList.contains('mcimg') ? 'MCIMG' : 'MC';
                counters[type]++;
                questionData = collectMultipleChoiceData(card, type, counters[type]);
                if (questionData) {
                    questionData.globalId = parseInt(card.dataset.globalQuestionId) || 0;
                    allQuestions.push(questionData);  // 只添加到 allQuestions
                }
            }
            // 填空题
            else if (classList.contains('fl') || classList.contains('flimg')) {
                const type = classList.contains('flimg') ? 'FLIMG' : 'FL';
                counters[type]++;
                questionData = collectFillBlankData(card, type, counters[type]);
                if (questionData) {
                    questionData.globalId = parseInt(card.dataset.globalQuestionId) || 0;
                    allQuestions.push(questionData);  // 只添加到 allQuestions
                }
            }
            // 材料题（自身 + 内部子题）
            else if (classList.contains('mt')) {
                const drGlobalId = parseInt(card.dataset.globalQuestionId) || 0;

                // 1) 先把材料题内部子题，按当前顺序当作普通题写入 allQuestions
                const innerCards = card.querySelectorAll('.mt-inner-card');
                console.log('📎 DR 材料题发现内部子题数量:', innerCards.length, ' (globalId=', drGlobalId, ')');

                innerCards.forEach((innerCard, idx) => {
                    const innerTypeLabel = innerCard.querySelector('.card-type-label');
                    if (!innerTypeLabel) {
                        console.warn('⚠️ 材料子题缺少 card-type-label, index=', idx + 1);
                        return;
                    }
                    const innerClassList = innerTypeLabel.classList;
                    let innerData = null;

                    if (innerClassList.contains('sc') || innerClassList.contains('scimg')) {
                        const t = innerClassList.contains('scimg') ? 'SCIMG' : 'SC';
                        counters[t]++;
                        innerData = collectSingleChoiceData(innerCard, t, counters[t]);
                    } else if (innerClassList.contains('mc') || innerClassList.contains('mcimg')) {
                        const t = innerClassList.contains('mcimg') ? 'MCIMG' : 'MC';
                        counters[t]++;
                        innerData = collectMultipleChoiceData(innerCard, t, counters[t]);
                    } else if (innerClassList.contains('fl') || innerClassList.contains('flimg')) {
                        const t = innerClassList.contains('flimg') ? 'FLIMG' : 'FL';
                        counters[t]++;
                        innerData = collectFillBlankData(innerCard, t, counters[t]);
                    } else {
                        console.warn('⚠️ 未识别的材料子题类型, classList=', Array.from(innerClassList).join(' '));
                    }

                    if (innerData) {
                        innerData.globalId = parseInt(innerCard.dataset.globalQuestionId) || 0;
                        console.log('✅ 已收集材料子题为普通题:', innerData.type, '#', innerData.id, 'hook=', innerData.hook || '(无)', 'globalId=', innerData.globalId);
                        allQuestions.push(innerData);
                    }
                });

                // 2) 再收集材料题本身（DR），保持原来逻辑
                counters.DR++;
                questionData = collectDocumentReadingData(card, counters.DR);
                if (questionData) {
                    questionData.globalId = drGlobalId;
                    allQuestions.push(questionData);
                }
            }

            // 收集图片（用于顶层 images 数组）
            const imagesFromCard = questionData ? questionData.images : collectImagesFromCard(card);
            if (imagesFromCard && imagesFromCard.length) {
                imagesFromCard.forEach(img => {
                    if (img && !imageSet.has(img)) {
                        imageSet.add(img);
                        const filename = img.split('/').pop();
                        images.push({
                            filename: filename,
                            path: img
                        });
                    }
                });
            }

        } catch (err) {
            errors.push(`处理卡片时出错: ${err.message}`);
            console.error(err);
        }
    });

    // 统计各题型数量（从 allQuestions 计算）
    const singleChoiceCount = allQuestions.filter(q => q.type === 'SC' || q.type === 'SCIMG').length;
    const multipleChoiceCount = allQuestions.filter(q => q.type === 'MC' || q.type === 'MCIMG').length;
    const fillBlankCount = allQuestions.filter(q => q.type === 'FL' || q.type === 'FLIMG').length;
    const documentReadingCount = allQuestions.filter(q => q.type === 'DR').length;

    // 生成元数据
    const metadata = {
        totalQuestions: allQuestions.length,
        singleChoice: singleChoiceCount,
        multipleChoice: multipleChoiceCount,
        fillBlank: fillBlankCount,
        documentReading: documentReadingCount,
        totalImages: images.length
    };
    
    // 添加调试日志
    console.log('📊 收集完成！');
    console.log('  - 总题数:', allQuestions.length);
    console.log('  - 全局顺序:', allQuestions.map(q => `[${q.type}#${q.id}]=全局#${q.globalId}`).join(' → '));

    // === 把扁平 allQuestions 映射成后端需要的 typed 结构 ===
    const typedQuestions = {
        singleChoice: [],
        multipleChoice: [],
        fillBlank: [],
        documentReading: []
    };

    allQuestions.forEach(q => {
        if (q.type === 'SC' || q.type === 'SCIMG') {
            typedQuestions.singleChoice.push(q);
        } else if (q.type === 'MC' || q.type === 'MCIMG') {
            typedQuestions.multipleChoice.push(q);
        } else if (q.type === 'FL' || q.type === 'FLIMG') {
            typedQuestions.fillBlank.push(q);
        } else if (q.type === 'DR') {
            typedQuestions.documentReading.push(q);
        }
    });

    return {
        name: bankName || '未命名题库',
        version: '1.0',
        metadata: metadata,
        images: images,
        // 只返回 typed 结构，避免 DR 在 flatQuestions 中重复出现
        questions: typedQuestions,
        errors: errors
    };
}

// 收集单选题数据
function collectSingleChoiceData(card, type, id) {
    const stemInput = card.querySelector('.stem-input');
    const question = stemInput ? stemInput.value.trim() : '';

    if (!question) {
        console.warn(`${type}题${id}: 题干为空`);
        return null;
    }

    const options = [];
    const optionInputs = card.querySelectorAll('.option-input');
    optionInputs.forEach(input => {
        const text = input.value.trim();
        if (text) options.push(text);
    });

    if (options.length === 0) {
        console.warn(`${type}题${id}: 没有选项`);
        return null;
    }

    // 从radio勾选状态获取答案
    let answer = '';
    const checkedRadio = card.querySelector('input[type="radio"]:checked');
    if (checkedRadio) {
        const optionDiv = checkedRadio.closest('div');
        const label = optionDiv.querySelector('span[style*="width:22px"]');
        answer = label ? label.textContent.trim() : '';
    }

    if (!answer) {
        console.warn(`${type}题${id}: 未勾选正确答案`);
    }

    const images = collectImagesFromCard(card);

    // 检查是否有Hook（材料题子题）
    let hook = '';
    const parentWrapper = card.closest('.mt-inner-card');
    if (parentWrapper && parentWrapper.dataset.hook) {
        hook = parentWrapper.dataset.hook;
    }

    const result = {
        id: String(id),
        type: type,
        enabled: true,
        question: question,
        images: images.length > 0 ? images : null,
        options: options,
        answer: answer || 'A' // 默认A（如果未勾选）
    };

    if (hook) {
        result.hook = hook;
    }

    return result;
}

// 收集多选题数据
function collectMultipleChoiceData(card, type, id) {
    const stemInput = card.querySelector('.stem-input');
    const question = stemInput ? stemInput.value.trim() : '';

    if (!question) return null;

    const options = [];
    const optionInputs = card.querySelectorAll('.option-input');
    optionInputs.forEach(input => {
        const text = input.value.trim();
        if (text) options.push(text);
    });    if (options.length === 0) return null;

    // 从checkbox勾选状态获取答案
    const answers = [];
    const checkedBoxes = card.querySelectorAll('input[type="checkbox"]:checked');
    checkedBoxes.forEach(cb => {
        const optionDiv = cb.closest('div');
        const label = optionDiv.querySelector('span[style*="width:22px"]');
        if (label) {
            answers.push(label.textContent.trim());
        }
    });

    if (answers.length === 0) {
        console.warn(`${type}题${id}: 未勾选正确答案`);
    }

    const images = collectImagesFromCard(card);

    // 检查是否有Hook（材料题子题）
    let hook = '';
    const parentWrapper = card.closest('.mt-inner-card');
    if (parentWrapper && parentWrapper.dataset.hook) {
        hook = parentWrapper.dataset.hook;
    }

    const result = {
        id: String(id),
        type: type,
        enabled: true,
        question: question,
        images: images.length > 0 ? images : null,
        options: options,
        answers: answers.length > 0 ? answers : ['A', 'B'] // 默认AB
    };

    if (hook) {
        result.hook = hook;
    }

    return result;
}

// 收集填空题数据
function collectFillBlankData(card, type, id) {
    const stemInput = card.querySelector('.stem-input');
    const question = stemInput ? stemInput.value.trim() : '';

    if (!question) return null;

    // 统计空位数量 - 修正：使用 3 个下划线的正则，与 Go 后端和机器考试一致
    const blankMatches = question.match(/\(%___%\)/g);
    const blankCount = blankMatches ? blankMatches.length : 0;

    if (blankCount === 0) {
        console.warn(`${type}题${id}: 没有空位`);
        return null;
    }    // 收集答案 - 需要从 blank-config-area 的 DOM 结构中正确读取
    const answers = [];
    let hasExtra = false;
    const extraKey = '(x%x)';
    
    const blankConfigArea = card.querySelector('.blank-config-area');
    if (blankConfigArea) {
        // 遍历每个空的配置块
        const blankBlocks = blankConfigArea.children;
        Array.from(blankBlocks).forEach((block, blockIndex) => {
            const blankIndex = blockIndex + 1;
            
            // 检查该空是否勾选了"学生答案不可重复使用"
            const uniqueCheckbox = block.querySelector('.blank-unique');
            const isUnique = uniqueCheckbox ? uniqueCheckbox.checked : false;
              // 收集该空的所有答案输入框
            const answerInputs = block.querySelectorAll('.blank-answer-input');
            const blankAnswers = [];
            
            // 如果该空启用了唯一性约束，先添加 (x%x) 标记到最前面
            if (isUnique) {
                blankAnswers.push('(x%x)');
                hasExtra = true;
            }
            
            // 然后添加实际答案
            answerInputs.forEach(input => {
                const val = input.value.trim();
                if (val) {
                    blankAnswers.push(val);
                }
            });
            
            // 只有当该空有至少一个答案时才添加
            if (blankAnswers.length > 0) {
                answers.push({
                    blankIndex: blankIndex,
                    answers: blankAnswers
                });
            }
        });
    }

    const images = collectImagesFromCard(card);

    // 检查是否有Hook（材料题子题）
    let hook = '';
    const parentWrapper = card.closest('.mt-inner-card');
    if (parentWrapper && parentWrapper.dataset.hook) {
        hook = parentWrapper.dataset.hook;
    }

    const result = {
        id: String(id),
        type: type,
        enabled: true,
        question: question,
        images: images.length > 0 ? images : null,
        template: question,
        blankCount: blankCount,
        answers: answers,
        hasExtra: hasExtra,
        extraKey: hasExtra ? extraKey : ''
    };

    if (hook) {
        result.hook = hook;
    }

    return result;
}

// 收集材料题数据
function collectDocumentReadingData(card, id) {
    const materialInput = card.querySelector('.material-input');
    const materials = materialInput ? [materialInput.value.trim()] : [];

    if (materials[0] === '') return null;

    // 收集hooks (从内部子题)
    const hooks = [];
    const innerCards = card.querySelectorAll('.mt-inner-card');
    innerCards.forEach(innerCard => {
        const hookAttr = innerCard.dataset.hook;
        if (hookAttr) hooks.push(hookAttr);
    });

    // 只收集材料题自身题干图片，不收集子题图片
    const images = collectImagesFromDocumentReadingCard(card);

    return {
        id: String(id),
        type: 'DR',
        enabled: true,
        question: '材料阅读题',
        images: images.length > 0 ? images : null,
        materials: materials,
        hooks: hooks
    };
}

// 专门为材料题收集图片：仅材料区/材料题干图片
function collectImagesFromDocumentReadingCard(card) {
    const images = [];

    // 限制：只统计“直接属于材料卡片本身”的图片，排除所有子题(.mt-inner-card)内部的上传按钮
    const uploadLabels = card.querySelectorAll('label.img-upload-btn[data-image-path]');
    uploadLabels.forEach(label => {
        // 如果这个 label 在任意一个子题容器内部，则跳过
        if (label.closest('.mt-inner-card')) {
            return;
        }
        const path = label.dataset.imagePath;
        if (path) {
            images.push(path);
        }
    });

    return images;
}

// 从卡片收集图片路径（普通题 & 子题使用）
function collectImagesFromCard(card) {
    const images = [];
    
    // 只在当前卡片内部查找（不需要区分材料/子题，因为子题自身不会再包含其它 .mt-inner-card）
    const uploadLabels = card.querySelectorAll('label.img-upload-btn[data-image-path]');
    uploadLabels.forEach(label => {
        const path = label.dataset.imagePath;
        if (path) {
            images.push(path);
        }
    });

    return images;
}

// 验证题库数据
function validateQuestionBank(bankData) {
    const errors = [];
    
    // 验证题库名称
    if (!bankData.name || bankData.name.trim() === '') {
        errors.push('题库名称不能为空');
    }
    
    // 验证题目数量
    if (!bankData.metadata || bankData.metadata.totalQuestions === 0) {
        errors.push('题库中没有题目');
    }

    const qGroups = bankData.questions || {};

    // 单选题/多选题验证（SC/SCIMG/MC/MCIMG）
    [...(qGroups.singleChoice || []), ...(qGroups.multipleChoice || [])].forEach((q, index) => {
        const qNum = index + 1;
        if (!q.question) {
            errors.push(`题目 ${qNum} (${q.type}): 题干为空`);
        }
        if (!q.options || q.options.length === 0) {
            errors.push(`题目 ${qNum} (${q.type}): 没有选项`);
        }
        if (q.type && q.type.includes('SC') && !q.answer) {
            errors.push(`题目 ${qNum} (${q.type}): 没有设置答案`);
        }
        if (q.type && q.type.includes('MC') && (!q.answers || q.answers.length === 0)) {
            errors.push(`题目 ${qNum} (${q.type}): 没有设置答案`);
        }
    });

    // 填空题验证（FL/FLIMG）
    (qGroups.fillBlank || []).forEach((q, index) => {
        const qNum = index + 1;
        if (!q.question) {
            errors.push(`填空题 ${qNum} (${q.type}): 题干为空`);
        }
        if (!q.blankCount || q.blankCount === 0) {
            errors.push(`填空题 ${qNum} (${q.type}): 没有空位标记`);
        }
        if (!q.answers || q.answers.length === 0) {
            errors.push(`填空题 ${qNum} (${q.type}): 没有设置答案`);
        }
    });

    // 材料题验证（DR）
    (qGroups.documentReading || []).forEach((q, index) => {
        const qNum = index + 1;
        if (!q.materials || q.materials.length === 0 || !q.materials[0]) {
            errors.push(`材料题 ${qNum} (${q.type}): 材料内容为空`);
        }
        if (!q.hooks || q.hooks.length === 0) {
            errors.push(`材料题 ${qNum} (${q.type}): 没有关联的子题`);
        }
    });

    return errors;
}

// 导出题库
async function exportQuestionBank() {
    try {
        // 获取题库名称
        const bankName = prompt('请输入题库名称:', '我的题库');
        if (!bankName) {
            alert('已取消导出');
            return;
        }

        // 收集数据
        console.log('正在收集题库数据...');
        const bankData = collectQuestionBankData(bankName);

        // 验证数据
        const validationErrors = validateQuestionBank(bankData);
        if (validationErrors.length > 0) {
            const errorMsg = '题库数据验证失败，发现以下问题:\n\n' + 
                           validationErrors.slice(0, 10).join('\n') +
                           (validationErrors.length > 10 ? '\n\n... 还有 ' + (validationErrors.length - 10) + ' 个问题' : '');
            
            if (!confirm(errorMsg + '\n\n是否仍要导出？')) {
                return;
            }
        }

        // 显示预览
        console.log('题库数据:', bankData);

        // 调用后端API导出
        console.log('正在导出ZIP文件...');
        const jsonStr = JSON.stringify(bankData);
        
        // 使用Wails绑定的Go方法
        const zipPath = await window.go.main.App.ExportQuestionBank(jsonStr);
        
        alert(`导出成功！\n文件路径: ${zipPath}\n\n统计信息:\n- 单选题: ${bankData.metadata.singleChoice}\n- 多选题: ${bankData.metadata.multipleChoice}\n- 填空题: ${bankData.metadata.fillBlank}\n- 材料题: ${bankData.metadata.documentReading}\n- 图片数: ${bankData.metadata.totalImages}`);

    } catch (error) {
        console.error('导出失败:', error);
        alert('导出失败: ' + error.message);
    }
}

// 预览题库JSON
async function previewQuestionBank() {
    try {
        const bankName = prompt('请输入题库名称(用于预览):', '预览题库');
        if (!bankName) return;

        const bankData = collectQuestionBankData(bankName);
        const jsonStr = JSON.stringify(bankData, null, 2);

        // 创建预览窗口
        const previewWindow = window.open('', '_blank', 'width=800,height=600');
        previewWindow.document.write(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>题库预览 - ${bankName}</title>
                <style>
                    body { font-family: monospace; padding: 20px; background: #f5f5f5; }
                    pre { background: white; padding: 20px; border-radius: 8px; overflow: auto; }
                    .stats { background: #e3f2fd; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
                </style>
            </head>
            <body>
                <h1>📚 ${bankName}</h1>
                <div class="stats">
                    <h3>统计信息</h3>
                    <p>总题数: ${bankData.metadata.totalQuestions}</p>
                    <p>单选题: ${bankData.metadata.singleChoice}</p>
                    <p>多选题: ${bankData.metadata.multipleChoice}</p>
                    <p>填空题: ${bankData.metadata.fillBlank}</p>
                    <p>材料题: ${bankData.metadata.documentReading}</p>
                    <p>图片数: ${bankData.metadata.totalImages}</p>
                </div>
                <h3>JSON数据</h3>
                <pre>${escapeHtml(jsonStr)}</pre>
            </body>
            </html>
        `);
    } catch (error) {
        alert('预览失败: ' + error.message);
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ===== 辅助：显示/隐藏全局加载遮罩 =====
function showGlobalLoading(message) {
    let mask = document.getElementById('globalLoadingMask');
    if (!mask) {
        mask = document.createElement('div');
        mask.id = 'globalLoadingMask';
        mask.style.cssText = `
            position: fixed; left:0; top:0; right:0; bottom:0;
            background: rgba(0,0,0,0.4);
            display:flex;align-items:center;justify-content:center;
            z-index: 9999;
            color:#fff;font-size:18px;font-weight:bold;
        `;
        const inner = document.createElement('div');
        inner.id = 'globalLoadingInner';
        inner.style.cssText = 'padding:16px 24px;background:rgba(0,0,0,0.75);border-radius:8px;';
        inner.textContent = message || '正在处理中...';
        mask.appendChild(inner);
        document.body.appendChild(mask);
    } else {
        const inner = document.getElementById('globalLoadingInner');
        if (inner) inner.textContent = message || '正在处理中...';
        mask.style.display = 'flex';
    }
}

function hideGlobalLoading() {
    const mask = document.getElementById('globalLoadingMask');
    if (mask) mask.style.display = 'none';
}

// 辅助函数：恢复图片显示（创建完整的图片预览 UI）
function restoreImageDisplay(imgList, imagePath) {
    const box = document.createElement('label');
    box.className = 'img-upload-btn';
    box.style = 'display:inline-block;width:120px;height:80px;border:2px solid #4caf50;border-radius:6px;cursor:pointer;text-align:center;position:relative;user-select:none;overflow:hidden;';
    
    const previewPath = '../tempwails/' + imagePath;
    const exportPath = imagePath.startsWith('add/') ? imagePath : 'add/' + imagePath;
    
    box.dataset.imagePath = exportPath;
    box.dataset.previewPath = previewPath;
    
    box.innerHTML = `
        <img src="${previewPath}" style="width:100%;height:100%;object-fit:cover;border-radius:4px;">
        <input type="file" accept="image/*" style="display:none;">
        <span class="img-filename" style="position:absolute;left:0;bottom:0;width:100%;font-size:11px;color:#fff;line-height:1.2;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;background:rgba(0,0,0,0.6);padding:2px 4px;">${imagePath.split('/').pop()}</span>
        <button class="img-delete-btn" style="position:absolute;top:2px;right:2px;width:20px;height:20px;border:none;background:#f44336;color:#fff;border-radius:50%;cursor:pointer;font-size:12px;line-height:1;padding:0;">×</button>
    `;
    
    // 绑定删除按钮事件
    const deleteBtn = box.querySelector('.img-delete-btn');
    if (deleteBtn) {
        deleteBtn.onclick = function(e) {
            e.preventDefault();
            e.stopPropagation();
            if (confirm('确定要删除这张图片吗？')) {
                box.remove();
            }
        };
    }
    
    imgList.appendChild(box);
}

// 根据导入的题库 JSON 重新渲染卡片
function renderQuestionBankFromJson(bank) {
    if (!bank || !bank.questions) {
        alert('导入的题库数据不合法');
        return;
    }
    console.log('🛠 renderQuestionBankFromJson，收到题库:', bank.name || '(未命名)');

    const qGroups = bank.questions;
    const cardList = document.getElementById('cardList');
    if (!cardList) return;

    // 1. 构建 hook -> 子题 的映射表
    const hookToQuestionMap = new Map();
    const allQuestions = [
        ...(qGroups.singleChoice || []),
        ...(qGroups.multipleChoice || []),
        ...(qGroups.fillBlank || [])
    ];
    
    allQuestions.forEach(q => {
        if (q.hook) {
            hookToQuestionMap.set(q.hook, q);
        }
    });
    
    console.log('📋 构建 hook 映射表，共', hookToQuestionMap.size, '个子题');

    // 2. 渲染普通题型（过滤掉有 hook 的子题）
    function addSingleChoiceCard(q) {
        if (q.hook) {
            console.log('⏭️ 跳过子题 SC:', q.id, 'hook=', q.hook);
            return; // 跳过子题，它们会在材料题中渲染
        }
        
        const idx = (window.typeCounters && window.typeCounters.SC + 1) || 1;
        const isImg = q.type === 'SCIMG';
        const card = isImg ? window.createSingleChoiceWithStemImgCard(idx) : window.createSingleChoiceCard(idx);
        // 题干
        const stem = card.querySelector('.stem-input');
        if (stem) stem.value = q.question || '';
        // 选项
        const optionInputs = card.querySelectorAll('.option-input');
        optionInputs.forEach((input, i) => {
            if (q.options && q.options[i]) input.value = q.options[i];
        });
        // 答案
        const radios = card.querySelectorAll('input[type="radio"]');
        if (radios.length && q.answer) {
            radios.forEach(radio => {
                const optDiv = radio.closest('div');
                const label = optDiv && optDiv.querySelector('span[style*="width:22px"]');
                if (label && label.textContent.trim() === q.answer) {
                    radio.checked = true;
                }
            });
        }        // 图片（使用完整的图片预览 UI）
        if (q.images && q.images.length) {
            const imgList = card.querySelector('.img-list') || card.querySelector('.stem-img-list');
            if (imgList) {
                q.images.forEach(path => {
                    restoreImageDisplay(imgList, path);
                });
            }
        }
        cardList.appendChild(card);
    }

    function addMultipleChoiceCard(q) {
        if (q.hook) {
            console.log('⏭️ 跳过子题 MC:', q.id, 'hook=', q.hook);
            return; // 跳过子题
        }
        
        const idx = (window.typeCounters && window.typeCounters.MC + 1) || 1;
        const isImg = q.type === 'MCIMG';
        const card = isImg ? window.createMultipleChoiceWithStemImgCard(idx) : window.createMultipleChoiceCard(idx);
        const stem = card.querySelector('.stem-input');
        if (stem) stem.value = q.question || '';
        const optionInputs = card.querySelectorAll('.option-input');
        optionInputs.forEach((input, i) => {
            if (q.options && q.options[i]) input.value = q.options[i];
        });
        const boxes = card.querySelectorAll('input[type="checkbox"]');
        if (boxes.length && q.answers && q.answers.length) {
            boxes.forEach(box => {
                const optDiv = box.closest('div');
                const label = optDiv && optDiv.querySelector('span[style*="width:22px"]');
                if (label && q.answers.includes(label.textContent.trim())) {
                    box.checked = true;
                }
            });        }
        if (q.images && q.images.length) {
            const imgList = card.querySelector('.img-list') || card.querySelector('.stem-img-list');
            if (imgList) {
                q.images.forEach(path => {
                    restoreImageDisplay(imgList, path);
                });
            }
        }
        cardList.appendChild(card);    }function addFillBlankCard(q) {
        if (q.hook) {
            console.log('⏭️ 跳过子题 FL:', q.id, 'hook=', q.hook);
            return; // 跳过子题
        }
        
        const idx = (window.typeCounters && window.typeCounters.FL + 1) || 1;
        const isImg = q.type === 'FLIMG';
        const card = isImg ? window.createFillBlankWithStemImgCard(idx) : window.createFillBlankCard(idx);
        
        // 1. 恢复题干
        const stem = card.querySelector('.stem-input');
        if (stem) stem.value = q.question || '';
        
        // 2. 准备填空答案数据
        let initialBlanks = null;
        if (q.answers && q.answers.length > 0) {
            // 将导入的答案数据转换为 setupFillBlankLogic 期望的格式
            initialBlanks = q.answers.map(blankData => {
                const answerList = blankData.answers || [];
                
                // 检查是否有 (x%x) 标记
                const hasUnique = answerList.includes('(x%x)');
                
                // 过滤掉 (x%x) 得到实际答案
                const actualAnswers = answerList.filter(a => a !== '(x%x)');
                
                // 返回 setupFillBlankLogic 期望的格式
                return {
                    answers: hasUnique ? ['(x%x)', ...actualAnswers] : actualAnswers,
                    unique: hasUnique
                };
            });
            
            // 重新调用 setupFillBlankLogic 以恢复事件绑定和内部状态
            // 注意：createFillBlankCard 已经调用过一次 setupFillBlankLogic(card)，
            // 但那次是用空数据初始化的。我们需要用导入的数据重新初始化。
            // 为了避免重复绑定事件，我们需要先清理旧的事件监听器。
            // 最简单的方法是只重建 UI，不重新绑定按钮事件。
            // 但这需要访问 setupFillBlankLogic 的内部函数，所以我们采用完全重新初始化的方式。
            
            // 由于按钮事件是通过 onclick 直接赋值的，重新调用 setupFillBlankLogic 会覆盖旧的事件
            // 这实际上是我们想要的，因为它会使用新的 blanks 数据
            window.setupFillBlankLogic(card, initialBlanks);
        }
          // 3. 恢复图片
        if (q.images && q.images.length) {
            const imgList = card.querySelector('.img-list') || card.querySelector('.stem-img-list');
            if (imgList) {
                q.images.forEach(path => {
                    restoreImageDisplay(imgList, path);
                });
            }
        }
        
        cardList.appendChild(card);
    }    function addDocumentReadingCard(dr, hookToQuestionMap) {
        // 创建材料题外卡片
        const idx = (window.typeCounters && window.typeCounters.DR + 1) || 1;
        const card = window.createMaterialCard(idx);
        const matInput = card.querySelector('.material-input');
        if (matInput && dr.materials && dr.materials.length) {
            matInput.value = dr.materials[0] || '';
        }
          // 恢复 DR 自身图片
        if (dr.images && dr.images.length) {
            const imgList = card.querySelector('.material-img-list');
            if (imgList) {
                dr.images.forEach(path => {
                    restoreImageDisplay(imgList, path);
                });
            }
        }
        
        // 模拟"检查材料"已通过，使子题按钮可用
        const checkBtn = card.querySelector('.material-check-btn');
        if (checkBtn) {
            checkBtn.textContent = '材料已通过 ✓';
            checkBtn.disabled = true;
        }
        const innerToolbar = card.querySelector('.material-inner-toolbar');
        if (innerToolbar) {
            innerToolbar.style.display = 'block';
        }
        
        cardList.appendChild(card);
        
        // 根据 hooks 恢复子题
        if (dr.hooks && dr.hooks.length > 0 && hookToQuestionMap) {
            const innerList = card.querySelector('.material-inner-list');
            let innerIndex = 1;
            
            dr.hooks.forEach(hook => {
                const subQuestion = hookToQuestionMap.get(hook);
                if (!subQuestion) {
                    console.warn(`⚠️ 未找到 hook=${hook} 对应的子题`);
                    return;
                }
                
                // 根据子题类型创建对应的卡片
                let innerCard;
                let typeLabelText = '';
                
                switch (subQuestion.type) {
                    case 'SC':
                        innerCard = window.createSingleChoiceCard(innerIndex);
                        typeLabelText = '单选题干无图';
                        break;
                    case 'SCIMG':
                        innerCard = window.createSingleChoiceWithStemImgCard(innerIndex);
                        typeLabelText = '单选题干有图';
                        break;
                    case 'MC':
                        innerCard = window.createMultipleChoiceCard(innerIndex);
                        typeLabelText = '多选题干无图';
                        break;
                    case 'MCIMG':
                        innerCard = window.createMultipleChoiceWithStemImgCard(innerIndex);
                        typeLabelText = '多选题干有图';
                        break;
                    case 'FL':
                        innerCard = window.createFillBlankCard(innerIndex);
                        typeLabelText = '填空题干无图';
                        break;
                    case 'FLIMG':
                        innerCard = window.createFillBlankWithStemImgCard(innerIndex);
                        typeLabelText = '填空题干有图';
                        break;
                    default:
                        console.warn(`⚠️ 未知的子题类型: ${subQuestion.type}`);
                        return;
                }
                
                // 去掉外层删除按钮
                const delBtn = innerCard.querySelector('.card-delete-btn');
                if (delBtn) delBtn.remove();
                
                // 更新序号
                const idxSpan = innerCard.querySelector('.card-index');
                if (idxSpan) idxSpan.textContent = innerIndex;
                
                // 恢复题目数据
                restoreQuestionData(innerCard, subQuestion);
                
                // 创建 wrapper
                const wrapper = document.createElement('div');
                wrapper.className = 'mt-inner-card';
                wrapper.dataset.hook = hook; // 保存 hook 以便导出时使用
                wrapper.innerHTML = `
                    <div class="mt-inner-header">
                        <span class="mt-inner-tag">内题 ${innerIndex}</span>
                        <span class="mt-inner-title">${typeLabelText}</span>
                        <button class="mt-inner-delete-btn" style="margin-left:auto;">🗑 删除子题</button>
                    </div>
                `;
                wrapper.appendChild(innerCard);
                
                // 绑定删除按钮
                const innerDeleteBtn = wrapper.querySelector('.mt-inner-delete-btn');
                innerDeleteBtn.onclick = function () {
                    const evt = new CustomEvent('mt-inner-delete', { detail: { wrapper } });
                    window.dispatchEvent(evt);
                };
                
                innerList.appendChild(wrapper);
                innerIndex++;
            });
        }
    }
    
    // 辅助函数：恢复题目数据到卡片
    function restoreQuestionData(card, q) {
        // 恢复题干
        const stem = card.querySelector('.stem-input');
        if (stem) stem.value = q.question || '';
        
        // 根据题型恢复特定数据
        if (q.type === 'SC' || q.type === 'SCIMG') {
            // 恢复选项
            const optionInputs = card.querySelectorAll('.option-input');
            optionInputs.forEach((input, i) => {
                if (q.options && q.options[i]) input.value = q.options[i];
            });
            // 恢复答案
            const radios = card.querySelectorAll('input[type="radio"]');
            if (radios.length && q.answer) {
                radios.forEach(radio => {
                    const optDiv = radio.closest('div');
                    const label = optDiv && optDiv.querySelector('span[style*="width:22px"]');
                    if (label && label.textContent.trim() === q.answer) {
                        radio.checked = true;
                    }
                });
            }
        } else if (q.type === 'MC' || q.type === 'MCIMG') {
            // 恢复选项
            const optionInputs = card.querySelectorAll('.option-input');
            optionInputs.forEach((input, i) => {
                if (q.options && q.options[i]) input.value = q.options[i];
            });
            // 恢复答案
            const boxes = card.querySelectorAll('input[type="checkbox"]');
            if (boxes.length && q.answers && q.answers.length) {
                boxes.forEach(box => {
                    const optDiv = box.closest('div');
                    const label = optDiv && optDiv.querySelector('span[style*="width:22px"]');
                    if (label && q.answers.includes(label.textContent.trim())) {
                        box.checked = true;
                    }
                });
            }
        } else if (q.type === 'FL' || q.type === 'FLIMG') {
            // 恢复填空题答案
            if (q.answers && q.answers.length > 0) {
                const initialBlanks = q.answers.map(blankData => {
                    const answerList = blankData.answers || [];
                    const hasUnique = answerList.includes('(x%x)');
                    const actualAnswers = answerList.filter(a => a !== '(x%x)');
                    return {
                        answers: hasUnique ? ['(x%x)', ...actualAnswers] : actualAnswers,
                        unique: hasUnique
                    };
                });
                window.setupFillBlankLogic(card, initialBlanks);
            }
        }
          // 恢复图片
        if (q.images && q.images.length) {
            const imgList = card.querySelector('.img-list') || card.querySelector('.stem-img-list');
            if (imgList) {
                q.images.forEach(path => {
                    restoreImageDisplay(imgList, path);
                });
            }
        }
    }    // 依次渲染
    (qGroups.singleChoice || []).forEach(addSingleChoiceCard);
    (qGroups.multipleChoice || []).forEach(addMultipleChoiceCard);
    (qGroups.fillBlank || []).forEach(addFillBlankCard);
    (qGroups.documentReading || []).forEach(dr => addDocumentReadingCard(dr, hookToQuestionMap));

    // 渲染完成后刷新左侧预览
    if (typeof window.refreshAllIndexes === 'function') {
        window.refreshAllIndexes();
    } else {
        console.log('⚠️ 未找到 refreshAllIndexes，左侧预览可能需要手动刷新');
    }
}

// 清空当前编辑区（DOM + 左侧预览），不负责删磁盘文件，磁盘由后端 Reset 处理
function clearEditorCards() {
    const cardList = document.getElementById('cardList');
    const slidePreviewList = document.getElementById('slidePreviewList');
    if (cardList) cardList.innerHTML = '';
    if (slidePreviewList) slidePreviewList.innerHTML = '';
    if (window.previewItems && Array.isArray(window.previewItems)) {
        window.previewItems.length = 0;
    }
}

// 修改题库：清空当前编辑 + 让后端导入ZIP并返回JSON，再渲染
async function modifyQuestionBank() {
    try {
        if (!confirm('该操作将清空当前正在编辑的所有题目，并删除临时图片文件。\n\n是否继续导入并修改新的题库？')) {
            return;
        }
        showGlobalLoading('正在清空并导入题库，请稍候...');

        // 1. 让后端清空 tempwails 目录（图片+JSON 等）
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.ClearTempImages) {
            try {
                await window.go.main.App.ClearTempImages();
                console.log('✓ 已通过后端清空 tempwails 目录');
            } catch (e) {
                console.warn('清空 tempwails 失败:', e);
            }
        }

        // 2. 前端先清空当前编辑区的卡片和预览（DOM）
        clearEditorCards();

        // 3. 调用后端 ImportQuestionBank：弹出文件选择框，解压到 fix，复制图片到 tempwails/add，并返回题库JSON字符串
        if (!window.go || !window.go.main || !window.go.main.App || !window.go.main.App.ImportQuestionBank) {
            hideGlobalLoading();
            alert('后端未实现 ImportQuestionBank 接口，无法导入题库');
            return;
        }

        console.log('📂 正在打开文件选择对话框以导入题库ZIP...');
        const jsonStr = await window.go.main.App.ImportQuestionBank();
        if (!jsonStr) {
            hideGlobalLoading();
            alert('未选择题库文件或导入被取消');
            return;
        }

        let bank;
        try {
            bank = JSON.parse(jsonStr);
        } catch (e) {
            console.error('解析导入题库JSON失败:', e);
            hideGlobalLoading();
            alert('导入题库失败：JSON 解析错误');
            return;
        }

        console.log('📥 已导入题库:', bank.name || '(未命名)', '题目统计:', bank.metadata || {});
        // 4. 根据 JSON 重新渲染所有卡片
        renderQuestionBankFromJson(bank);
        hideGlobalLoading();
    } catch (err) {
        console.error('修改题库流程出错:', err);
        hideGlobalLoading();
        alert('修改题库失败: ' + err.message);
    }
}

// 导出到window对象供外部调用
window.exportQuestionBank = exportQuestionBank;
window.previewQuestionBank = previewQuestionBank;
window.collectQuestionBankData = collectQuestionBankData;
window.modifyQuestionBank = modifyQuestionBank;
