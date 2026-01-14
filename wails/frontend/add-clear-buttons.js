// 动态为所有卡片的图片操作区域添加"清空图片"按钮

/**
 * 为图片操作工具栏添加清空按钮
 */
function addClearButtonsToCards() {
    console.log('🔧 开始为卡片添加清空图片按钮...');
    
    // 查找所有包含 "img-add-btn" 的容器（图片操作区域）
    const toolbars = document.querySelectorAll('.img-add-btn');    toolbars.forEach(addBtn => {
        const toolbar = addBtn.parentElement;
        
        // 检查是否已经添加过清空按钮（检查两个按钮中的任意一个）
        if (toolbar.querySelector('.img-clear-local-btn') || toolbar.querySelector('.img-clear-all-btn')) {
            return; // 已存在，跳过
        }
        
        // 添加标记，防止重复处理
        if (toolbar.dataset.clearButtonsAdded === 'true') {
            return;
        }
        toolbar.dataset.clearButtonsAdded = 'true';
          // 1. 创建"清空本区域图片"按钮
        const clearLocalBtn = document.createElement('button');
        clearLocalBtn.className = 'img-clear-local-btn';
        clearLocalBtn.textContent = '🗑️ 清空本区域';
        clearLocalBtn.title = '删除当前图片区域的所有图片';
        clearLocalBtn.style.cssText = `
            background: linear-gradient(135deg, #ff9800, #f57c00);
            color: white;
            padding: 6px 12px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 13px;
            font-weight: 500;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
            transition: all 0.2s ease;
        `;
        
        // 悬停效果
        clearLocalBtn.onmouseenter = function() {
            clearLocalBtn.style.transform = 'scale(1.05)';
            clearLocalBtn.style.boxShadow = '0 4px 8px rgba(255,152,0,0.5)';
        };
        clearLocalBtn.onmouseleave = function() {
            clearLocalBtn.style.transform = 'scale(1)';
            clearLocalBtn.style.boxShadow = '0 2px 4px rgba(0,0,0,0.2)';
        };
        
        // 点击事件：清空当前区域
        clearLocalBtn.onclick = async function(e) {
            e.preventDefault();
            e.stopPropagation();
            
            // 找到当前图片列表
            const imgList = toolbar.parentElement.querySelector('.img-list, .stem-img-list, .material-img-list');
            if (!imgList) {
                alert('未找到图片列表');
                return;
            }
            
            const imgBoxes = imgList.querySelectorAll('.img-upload-btn[data-image-path]');
            if (imgBoxes.length === 0) {
                alert('当前区域没有图片');
                return;
            }
            
            const confirmMsg = `确定要清空当前区域的 ${imgBoxes.length} 张图片吗？`;
            if (!confirm(confirmMsg)) {
                return;
            }
            
            let successCount = 0;
            for (const imgBox of imgBoxes) {
                if (typeof window.deleteCardImage === 'function') {
                    const success = await window.deleteCardImage(imgBox);
                    if (success) successCount++;
                }
            }
            
            alert(`成功删除 ${successCount}/${imgBoxes.length} 张图片`);
        };
        
        // 2. 创建"清空所有图片"按钮
        const clearAllBtn = document.createElement('button');
        clearAllBtn.className = 'img-clear-all-btn';
        clearAllBtn.textContent = '🗑️ 清空所有图片';
        clearAllBtn.title = '删除整个卡片的所有图片（包括题干图片、附带图片等）';
        clearAllBtn.style.cssText = `
            background: linear-gradient(135deg, #ff5252, #f44336);
            color: white;
            padding: 6px 12px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 13px;
            font-weight: 500;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
            transition: all 0.2s ease;
        `;
        
        // 悬停效果
        clearAllBtn.onmouseenter = function() {
            clearAllBtn.style.transform = 'scale(1.05)';
            clearAllBtn.style.boxShadow = '0 4px 8px rgba(255,0,0,0.5)';
        };
        clearAllBtn.onmouseleave = function() {
            clearAllBtn.style.transform = 'scale(1)';
            clearAllBtn.style.boxShadow = '0 2px 4px rgba(0,0,0,0.2)';
        };
        
        // 点击事件：清空整个卡片的所有图片
        clearAllBtn.onclick = async function(e) {
            e.preventDefault();
            e.stopPropagation();
            
            // 找到所属的卡片
            const card = toolbar.closest('.single-card, .multiple-card, .fill-card, .material-card');
            if (!card) {
                alert('未找到所属卡片');
                return;
            }
            
            if (typeof window.deleteAllCardImages === 'function') {
                await window.deleteAllCardImages(card);
            } else {
                alert('删除功能未加载，请刷新页面');
            }
        };
        
        // 将两个按钮插入到"添加图片"按钮之后
        addBtn.insertAdjacentElement('afterend', clearLocalBtn);
        clearLocalBtn.insertAdjacentElement('afterend', clearAllBtn);
    });
    
    console.log(`✓ 已为 ${toolbars.length} 个图片区域添加清空按钮`);
}

/**
 * 使用 MutationObserver 监听新卡片的创建
 */
function observeNewCards() {
    const cardList = document.getElementById('cardList');
    if (!cardList) {
        console.warn('未找到 cardList，无法监听新卡片');
        return;
    }
    
    const observer = new MutationObserver(function(mutations) {
        let hasNewCards = false;
        
        mutations.forEach(function(mutation) {
            // 只监听 cardList 的直接子节点添加
            if (mutation.target === cardList && mutation.type === 'childList') {
                mutation.addedNodes.forEach(function(node) {
                    if (node.nodeType === 1 && 
                        (node.classList.contains('single-card') || 
                         node.classList.contains('multiple-card') || 
                         node.classList.contains('fill-card') || 
                         node.classList.contains('material-card'))) {
                        hasNewCards = true;
                    }
                });
            }
        });
        
        if (hasNewCards) {
            // 延迟一小段时间，确保卡片完全渲染
            setTimeout(addClearButtonsToCards, 100);
        }
    });
    
    // 只监听 cardList 的直接子节点变化
    observer.observe(cardList, {
        childList: true,
        subtree: false  // 不监听子树，避免重复触发
    });
    
    console.log('✓ 已启动新卡片监听器');
}

/**
 * 初始化
 */
function init() {
    console.log('🚀 初始化清空图片按钮模块...');
    
    // 为现有卡片添加按钮
    addClearButtonsToCards();
    
    // 监听新卡片
    observeNewCards();
    
    console.log('✓ 清空图片按钮模块初始化完成');
}

// 页面加载完成后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}

console.log('✓ add-clear-buttons.js 已加载');
