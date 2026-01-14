// 增强卡片功能 - 添加正确答案输入等必要字段

// 在页面加载完成后为所有卡片添加答案输入框
document.addEventListener('DOMContentLoaded', function() {
    // 监听卡片添加事件
    observeCardAdditions();
});

// 监听卡片列表的变化
function observeCardAdditions() {
    const cardList = document.getElementById('cardList');
    if (!cardList) {
        console.warn('⚠️ [Hook] 未找到 #cardList，材料题 Hook 增强无法启用');
        return;
    }

    console.log('🔍 [Hook] observeCardAdditions 初始化，开始监听题卡添加');

    const observer = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            mutation.addedNodes.forEach(function(node) {
                if (node.nodeType !== 1) return; // 只处理元素节点

                // 顶层题卡：单选、多选、填空、材料
                if (node.classList.contains('single-card') ||
                    node.classList.contains('multiple-card') ||
                    node.classList.contains('fill-card') ||
                    node.classList.contains('material-card')) {
                    console.log('➕ [Hook] 检测到新题卡加入 cardList:', node.className);
                    enhanceCard(node);
                }
            });
        });
    });

    observer.observe(cardList, { childList: true, subtree: false });

    // 增强已存在的卡片
    document.querySelectorAll('#cardList > div').forEach(card => {
        console.log('🔁 [Hook] 初始化已存在题卡:', card.className);
        enhanceCard(card);
    });
}

// 增强单个卡片 - 不添加额外输入框，答案直接从UI组件读取
function enhanceCard(card) {
    if (card.dataset.enhanced === 'true') return;
    card.dataset.enhanced = 'true';

    const typeLabel = card.querySelector('.card-type-label');
    if (!typeLabel) {
        console.warn('⚠️ [Hook] 题卡缺少 .card-type-label，class=', card.className);
        return;
    }

    const classList = typeLabel.classList;
    console.log('🔎 [Hook] enhanceCard 处理题卡, typeLabel.classList =', Array.from(classList).join(' '));

    if (classList.contains('sc') || classList.contains('scimg')) {
        // 答案从 radio 勾选状态获取
    } else if (classList.contains('mc') || classList.contains('mcimg')) {
        // 答案从 checkbox 勾选状态获取
    } else if (classList.contains('fl') || classList.contains('flimg')) {
        // 填空题已有答案配置区域，无需额外处理
    } else if (classList.contains('mt')) {
        console.log('🔧 [Hook] 检测到材料题卡片，初始化 Hook 管理器');
        addMaterialHookManager(card);
        assignHooksForMaterialCard(card);
    } else {
        console.log('ℹ️ [Hook] 未识别为标准题型的卡片，classList=', Array.from(classList).join(' '));
    }
}

// 为材料题添加Hook管理
function addMaterialHookManager(card) {
    if (card.querySelector('.hook-manager')) {
        console.log('ℹ️ [Hook] 材料题已存在 Hook 管理器，跳过重复创建');
        return;
    }

    const materialInnerToolbar = card.querySelector('.material-inner-toolbar');
    if (!materialInnerToolbar) {
        console.warn('⚠️ [Hook] 材料题缺少 .material-inner-toolbar，无法插入 Hook 管理器');
        return;
    }

    const hookManager = document.createElement('div');
    hookManager.className = 'hook-manager';
    hookManager.style.cssText = 'margin-top:12px;padding:12px;background:#fff9c4;border-radius:8px;border-left:4px solid #ffa726;';
    hookManager.innerHTML = `
        <div style="margin-bottom:8px;">
            <label style="font-weight:600;color:#ef6c00;display:block;margin-bottom:4px;">🔗 子题关联 (Hooks)</label>
            <div style="font-size:13px;color:#666;margin-bottom:8px;">
                为材料题添加子题时，系统会自动生成Hook标识。例如: SC.A1, MC.B1, FL.C1
            </div>
            <div class="hooks-list" style="display:flex;flex-wrap:wrap;gap:6px;min-height:30px;">
                <span style="color:#999;font-size:13px;">暂无子题关联</span>
            </div>
        </div>
    `;

    materialInnerToolbar.parentNode.insertBefore(hookManager, materialInnerToolbar.nextSibling);
    
    // 监听子题添加，自动更新hooks
    const materialInnerList = card.querySelector('.material-inner-list');
    if (materialInnerList) {
        console.log('👀 [Hook] 开始监听材料题内部子题列表变化');
        const observer = new MutationObserver(function(mutations) {
            let changed = false;
            mutations.forEach(m => {
                if (m.addedNodes.length || m.removedNodes.length) changed = true;
            });
            if (!changed) return;

            console.log('🔁 [Hook] 材料题内部子题发生变更，重新分配 hooks ...');
            assignHooksForMaterialCard(card);
            updateMaterialHooks(card);
        });
        observer.observe(materialInnerList, { childList: true, subtree: false });
    } else {
        console.warn('⚠️ [Hook] 未找到 .material-inner-list，无法监听子题变化');
    }

    console.log('✓ [Hook] 已为材料题卡片添加 Hook 管理器');
    assignHooksForMaterialCard(card);
    updateMaterialHooks(card);
}

// 根据题型 class 映射到 Hook 字母前缀
function getHookPrefixByClassList(classList) {
    if (classList.contains('sc')) return 'A';       // SC
    if (classList.contains('scimg')) return 'B';    // SCIMG
    if (classList.contains('mc')) return 'C';       // MC
    if (classList.contains('mcimg')) return 'D';    // MCIMG
    if (classList.contains('fl')) return 'E';       // FL
    if (classList.contains('flimg')) return 'F';    // FLIMG
    return '';
}

// 为单个材料题卡片内的所有子题重新分配 data-hook
function assignHooksForMaterialCard(materialCard) {
    const innerCards = materialCard.querySelectorAll('.mt-inner-card');
    console.log('🔎 [Hook] assignHooksForMaterialCard 被调用, 找到子题数量 =', innerCards.length);
    if (!innerCards || innerCards.length === 0) {
        console.log('ℹ️ [Hook] 材料题暂无内部子题，跳过 Hook 分配');
        return;
    }

    console.log('🔧 [Hook] 开始为材料题重新分配内部子题 hooks，子题数量:', innerCards.length);

    // 为每种题型维护一个递增计数器
    const typeCounters = {
        'SC': 0,
        'SCIMG': 0,
        'MC': 0,
        'MCIMG': 0,
        'FL': 0,
        'FLIMG': 0,
    };

    innerCards.forEach((innerCard, index) => {
        const typeLabel = innerCard.querySelector('.card-type-label');
        if (!typeLabel) {
            console.warn('⚠️ [Hook] 内部子题缺少 card-type-label，跳过:', index + 1);
            return;
        }

        const classList = typeLabel.classList;
        let typeCode = '';
        let letter = '';

        if (classList.contains('sc')) {
            typeCode = 'SC';
            letter = 'A';
        } else if (classList.contains('scimg')) {
            typeCode = 'SCIMG';
            letter = 'B';
        } else if (classList.contains('mc')) {
            typeCode = 'MC';
            letter = 'C';
        } else if (classList.contains('mcimg')) {
            typeCode = 'MCIMG';
            letter = 'D';
        } else if (classList.contains('fl')) {
            typeCode = 'FL';
            letter = 'E';
        } else if (classList.contains('flimg')) {
            typeCode = 'FLIMG';
            letter = 'F';
        } else {
            console.warn('⚠️ [Hook] 未知子题类型 classList=', Array.from(classList).join(' '));
            innerCard.dataset.hook = '';
            return;
        }

        // 当前类型序号 +1
        typeCounters[typeCode] = (typeCounters[typeCode] || 0) + 1;
        const n = typeCounters[typeCode];

        // 生成 hook: 例如 SC.A1, FLIMG.F1
        const hook = typeCode + '.' + letter + String(n);
        innerCard.dataset.hook = hook;

        console.log(`✅ [Hook] 为材料题内部子题分配 hook: index=${index + 1}, type=${typeCode}, n=${n}, hook=${hook}`);
    });
}

// 更新材料题的Hook列表
function updateMaterialHooks(materialCard) {
    const hooksList = materialCard.querySelector('.hooks-list');
    if (!hooksList) {
        console.warn('⚠️ [Hook] 未找到 hooks-list，无法更新材料题 hooks 展示');
        return;
    }

    const innerCards = materialCard.querySelectorAll('.mt-inner-card');
    console.log('🔎 [Hook] updateMaterialHooks 被调用, 子题数量 =', innerCards.length);

    const hooks = [];

    innerCards.forEach((innerCard, index) => {
        const typeLabel = innerCard.querySelector('.card-type-label');
        if (!typeLabel) {
            console.warn('⚠️ [Hook] 更新 hooks 时发现子题缺少 card-type-label，index=', index + 1);
            return;
        }

        // 直接使用 assignHooksForMaterialCard 写入的 data-hook
        const hookAttr = innerCard.dataset.hook;
        if (!hookAttr) {
            console.warn('⚠️ [Hook] 子题尚未分配 hook，index=', index + 1);
            return;
        }

        hooks.push(hookAttr);
    });

    console.log('📌 [Hook] 材料题当前 hooks 列表:', hooks.join(', ') || '(空)');

    // 更新显示
    if (hooks.length === 0) {
        hooksList.innerHTML = '<span style="color:#999;font-size:13px;">暂无子题关联</span>';
    } else {
        hooksList.innerHTML = hooks.map(hook => 
            `<span style="background:#fff;padding:4px 8px;border-radius:4px;border:1px solid #ffa726;font-size:12px;font-weight:600;color:#ef6c00;">${hook}</span>`
        ).join('');
    }
}

// 为图片上传添加路径存储和规范命名
function enhanceImageUpload() {
    document.addEventListener('change', async function(e) {
        if (e.target.type !== 'file' || !e.target.accept || !e.target.accept.includes('image')) {
            return;
        }

        const file = e.target.files[0];
        if (!file) return;

        console.log('📷 检测到图片上传:', file.name);

        try {
            // 找到当前图片所属的题卡：可能是普通题卡，也可能是材料题或材料子题
            let card = e.target.closest('.single-card, .multiple-card, .fill-card, .material-card');
            if (!card) {
                console.warn('未找到父级卡片');
                return;
            }

            // 如果是在材料题内部子题中，card 会是内部的 single-card/multiple-card/fill-card，
            // 此时我们需要找到所属的外层材料题 .material-card
            const parentMaterialCard = card.closest('.material-card');
            let questionType = '';
            let typeIndex = 1;

            if (parentMaterialCard) {
                // 所有材料题以及其内部子题的图片，统一按 DR_材料题序号_X 命名
                const mtLabel = parentMaterialCard.querySelector('.card-type-label');
                const mtIndexSpan = parentMaterialCard.querySelector('.card-index');
                const drIndex = mtIndexSpan ? parseInt(mtIndexSpan.textContent) || 1 : 1;
                questionType = 'DR';
                typeIndex = drIndex;
                console.log(`🧩 检测到材料题或子题图片，使用 DR 命名: DR_${drIndex}_X`);

                // 后续保存时，card 统一用外层材料题，保证 collectImagesFromCard 能拿到所有图片
                card = parentMaterialCard;
            } else {
                // 普通题卡：按各自题型+题号命名
                const typeLabel = card.querySelector('.card-type-label');
                if (!typeLabel) {
                    console.warn('未找到题型标签');
                    return;
                }
                const classList = typeLabel.classList;
                if (classList.contains('sc')) questionType = 'SC';
                else if (classList.contains('scimg')) questionType = 'SCIMG';
                else if (classList.contains('mc')) questionType = 'MC';
                else if (classList.contains('mcimg')) questionType = 'MCIMG';
                else if (classList.contains('fl')) questionType = 'FL';
                else if (classList.contains('flimg')) questionType = 'FLIMG';
                else if (classList.contains('mt')) questionType = 'DR';

                if (!questionType) {
                    console.warn('未识别题型');
                    return;
                }

                const cardIndexSpan = card.querySelector('.card-index');
                typeIndex = cardIndexSpan ? parseInt(cardIndexSpan.textContent) || 1 : 1;
                console.log(`📝 普通题图片，命名前缀: ${questionType}_${typeIndex}_X`);
            }

            // 计算该题卡下已有的图片数量，作为当前图片序号
            const allImgBoxes = card.querySelectorAll('.img-upload-btn[data-image-path]');
            const imageIndex = allImgBoxes.length + 1;
            console.log(`📋 题型: ${questionType}, 序号: ${typeIndex}, 图片序号: ${imageIndex} (已有${allImgBoxes.length}张图片)`);

            const reader = new FileReader();
            reader.onload = async function(event) {
                const imgData = event.target.result;
                const imgBox = e.target.closest('label.img-upload-btn');
                if (!imgBox) {
                    console.warn('未找到图片容器');
                    return;
                }

                try {
                    console.log('🚀 调用后端API保存图片...');
                    const imagePath = await window.go.main.App.SaveImage(
                        file.name,
                        imgData,
                        questionType,
                        typeIndex,
                        imageIndex
                    );
                    console.log('✅ 后端返回路径:', imagePath);

                    const previewPath = '../tempwails/' + imagePath;
                    const exportPath = 'add/' + imagePath;
                    imgBox.dataset.imagePath = exportPath;
                    imgBox.dataset.previewPath = previewPath;
                    imgBox.dataset.imageData = imgData;

                    imgBox.innerHTML = `
                        <img src="${previewPath}" style="width:100%;height:100%;object-fit:cover;border-radius:4px;">
                        <input type="file" accept="image/*" style="display:none;">
                        <span class="img-filename" style="position:absolute;left:0;bottom:0;width:100%;font-size:11px;color:#fff;line-height:1.2;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;background:rgba(0,0,0,0.6);padding:2px 4px;">${imagePath}</span>
                    `;

                    const newInput = imgBox.querySelector('input[type="file"]');
                    if (newInput) {
                        newInput.addEventListener('change', function(newE) {
                            document.dispatchEvent(new Event('change', { target: newE.target }));
                        });
                    }

                    console.log(`✓ 图片已保存: ${imagePath} (${questionType}_${typeIndex}_${imageIndex})`);
                } catch (error) {
                    console.error('❌ 保存图片失败:', error);
                    alert('保存图片失败: ' + error.message);
                }
            };
            reader.readAsDataURL(file);
        } catch (error) {
            console.error('❌ 读取文件失败:', error);
            alert('读取文件失败: ' + error.message);
        }
    }, true);
}

// 初始化增强功能
enhanceImageUpload();

// 提供获取卡片答案的工具函数
window.getCardAnswer = function(card) {
    const typeLabel = card.querySelector('.card-type-label');
    if (!typeLabel) return '';

    const classList = typeLabel.classList;

    // 单选题 - 从radio获取答案
    if (classList.contains('sc') || classList.contains('scimg')) {
        const checkedRadio = card.querySelector('input[type="radio"]:checked');
        if (checkedRadio) {
            const optionDiv = checkedRadio.closest('div');
            const label = optionDiv.querySelector('span[style*="width:22px"]');
            return label ? label.textContent.trim() : '';
        }
        return '';
    }
    
    // 多选题 - 从checkbox获取答案（返回数组）
    if (classList.contains('mc') || classList.contains('mcimg')) {
        const checkedBoxes = card.querySelectorAll('input[type="checkbox"]:checked');
        const answers = [];
        checkedBoxes.forEach(cb => {
            const optionDiv = cb.closest('div');
            const label = optionDiv.querySelector('span[style*="width:22px"]');
            if (label) answers.push(label.textContent.trim());
        });
        return answers;
    }

    return '';
};

// 提供获取卡片图片路径的工具函数
window.getCardImages = function(card) {
    const images = [];
    const imgBoxes = card.querySelectorAll('.img-upload-btn[data-image-path]');
    imgBoxes.forEach(box => {
        const path = box.dataset.imagePath;
        if (path) images.push(path);
    });
    return images;
};

// ===== 统一处理题卡删除：先删图片再删卡片本身 =====
window.addEventListener('card-delete', function(e) {
    const card = e.detail && e.detail.card;
    if (!card) return;

    // 如果图片删除模块已加载，则优先删除该卡片下的所有图片
    if (window.deleteAllCardImages) {
        // 注意：deleteAllCardImages 内部已包含用户确认弹窗
        window.deleteAllCardImages(card).then(() => {
            // 图片删除流程结束后，再从 DOM 中移除整张题卡
            if (card.parentNode) {
                card.parentNode.removeChild(card);
            }
        }).catch(() => {
            // 即使图片删除出错，也允许用户继续删卡片，避免卡片无法移除
            if (card.parentNode) {
                card.parentNode.removeChild(card);
            }
        });
    } else {
        // 兜底：未加载 delete-images.js 时，保持原有行为，仅从 DOM 中移除
        if (card.parentNode) {
            card.parentNode.removeChild(card);
        }
    }
});

// ===== 统一处理材料题内部子题删除：先删图片再删DOM =====
window.addEventListener('mt-inner-delete', function(e) {
    const wrapper = e.detail && e.detail.wrapper;
    if (!wrapper) return;

    const doRemove = () => {
        if (wrapper.parentNode) {
            wrapper.parentNode.removeChild(wrapper);
        }
    };

    if (window.deleteAllCardImages) {
        // deleteAllCardImages 会提示并逐张删图（含后端文件）
        window.deleteAllCardImages(wrapper).then(() => {
            doRemove();
        }).catch(() => {
            // 即使删图出错，也不要卡死 UI，仍然允许移除子题
            doRemove();
        });
    } else {
        doRemove();
    }
});

console.log('✓ 卡片增强功能已加载');
