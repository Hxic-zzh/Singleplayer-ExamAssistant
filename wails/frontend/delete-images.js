// 图片删除功能 - 通过 JSON 桥接与后端通信

/**
 * 删除卡片中的某张图片
 * @param {HTMLElement} imgBox - 图片框元素 (label.img-upload-btn)
 */
async function deleteCardImage(imgBox) {
    const imagePath = imgBox.dataset.imagePath;
    
    if (!imagePath) {
        console.warn('图片框没有 imagePath，无需删除');
        return;
    }
    
    console.log(`🗑 准备删除图片: ${imagePath}`);
    
    try {
        // 1. 读取现有的删除列表
        let deleteList = [];
        try {
            const response = await fetch('../tempwails/delete_images.json');
            if (response.ok) {
                const data = await response.json();
                deleteList = data.images || [];
            }
        } catch (e) {
            console.log('delete_images.json 不存在或为空，创建新的');
        }
          // 2. 添加新的删除路径
        if (!deleteList.includes(imagePath)) {
            deleteList.push(imagePath);
        }
        
        // 3. 构造 JSON 数据
        const jsonData = {
            timestamp: new Date().toISOString(),
            images: [imagePath]
        };
        
        // 4. 调用后端 API 删除文件
        try {
            await window.go.main.App.ProcessPendingDeleteImages(JSON.stringify(jsonData));
            console.log('✓ 后端已成功删除图片文件');
        } catch (backendError) {
            console.warn('后端删除图片时出错:', backendError);
            // 不阻止继续执行，因为 DOM 已经清理
        }
          // 5. 清除 DOM 中的图片
        imgBox.dataset.imagePath = '';
        imgBox.dataset.previewPath = '';
        imgBox.dataset.imageData = '';
        imgBox.innerHTML = `+ 添加图片<input type="file" accept="image/*" style="display:none;">`;
        
        console.log('✓ 已清除 DOM 中的图片数据');
        
        return true;
    } catch (error) {
        console.error('❌ 删除图片失败:', error);
        alert('删除图片失败: ' + error.message);
        return false;
    }
}

/**
 * 删除卡片的所有图片
 * @param {HTMLElement} card - 卡片元素
 */
async function deleteAllCardImages(card) {
    const imgBoxes = card.querySelectorAll('.img-upload-btn[data-image-path]');
    
    if (imgBoxes.length === 0) {
        alert('该卡片没有图片');
        return;
    }
    
    const confirmMsg = `确定要删除该卡片的所有 ${imgBoxes.length} 张图片吗？\n\n此操作将同时删除：\n1. DOM 中的图片记录\n2. 磁盘上的图片文件`;
    
    if (!confirm(confirmMsg)) {
        return;
    }
    
    let successCount = 0;
    
    for (const imgBox of imgBoxes) {
        const success = await deleteCardImage(imgBox);
        if (success) successCount++;
    }
    
    alert(`成功删除 ${successCount}/${imgBoxes.length} 张图片\n\n图片文件将在后端处理后彻底删除`);
}

/**
 * 获取待删除的图片列表（从 localStorage）
 */
function getPendingDeleteImages() {
    try {
        const data = localStorage.getItem('pendingDeleteImages');
        if (!data) return [];
        
        const parsed = JSON.parse(data);
        return parsed.images || [];
    } catch (e) {
        console.error('读取待删除图片列表失败:', e);
        return [];
    }
}

/**
 * 清空待删除列表（后端删除完成后调用）
 */
function clearPendingDeleteImages() {
    localStorage.removeItem('pendingDeleteImages');
    console.log('✓ 已清空待删除图片列表');
}

/**
 * 为图片框添加删除按钮
 * @param {HTMLElement} imgBox - 图片框元素
 */
function addDeleteButtonToImageBox(imgBox) {
    // 避免重复添加
    if (imgBox.querySelector('.img-delete-btn')) {
        return;
    }
    
    // 只为已上传的图片添加删除按钮
    if (!imgBox.dataset.imagePath) {
        return;
    }
      const deleteBtn = document.createElement('button');
    deleteBtn.className = 'img-delete-btn';
    deleteBtn.innerHTML = '🗑️';
    deleteBtn.title = '删除这张图片';
    deleteBtn.style.cssText = `
        position: absolute;
        top: 2px;
        right: 2px;
        width: 32px;
        height: 32px;
        border-radius: 6px;
        background: linear-gradient(135deg, #ff5252, #f44336);
        color: white;
        border: 2px solid white;
        cursor: pointer;
        font-size: 16px;
        line-height: 1;
        padding: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 2px 8px rgba(0,0,0,0.4);
        z-index: 100;
        transition: all 0.2s ease;
    `;
    
    // 添加悬停效果
    deleteBtn.onmouseenter = function() {
        deleteBtn.style.transform = 'scale(1.15)';
        deleteBtn.style.boxShadow = '0 4px 12px rgba(255,0,0,0.5)';
    };
    deleteBtn.onmouseleave = function() {
        deleteBtn.style.transform = 'scale(1)';
        deleteBtn.style.boxShadow = '0 2px 8px rgba(0,0,0,0.4)';
    };
    
    deleteBtn.onclick = async function(e) {
        e.preventDefault();
        e.stopPropagation();
        
        const confirmMsg = '确定要删除这张图片吗？';
        if (!confirm(confirmMsg)) {
            return;
        }
        
        await deleteCardImage(imgBox);
        deleteBtn.remove(); // 删除按钮本身
    };
    
    imgBox.appendChild(deleteBtn);
}

/**
 * 监听图片上传，自动为新图片添加删除按钮
 */
function initImageDeleteButtons() {
    // 使用 MutationObserver 监听 DOM 变化
    const observer = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            mutation.addedNodes.forEach(function(node) {
                if (node.nodeType === 1) { // Element node
                    // 检查是否是图片框
                    if (node.classList && node.classList.contains('img-upload-btn')) {
                        addDeleteButtonToImageBox(node);
                    }
                    
                    // 检查子元素中的图片框
                    const imgBoxes = node.querySelectorAll('.img-upload-btn[data-image-path]');
                    imgBoxes.forEach(imgBox => {
                        addDeleteButtonToImageBox(imgBox);
                    });
                }
            });
            
            // 监听属性变化（图片上传完成时会设置 data-image-path）
            if (mutation.type === 'attributes' && 
                mutation.attributeName === 'data-image-path') {
                const target = mutation.target;
                if (target.classList.contains('img-upload-btn')) {
                    addDeleteButtonToImageBox(target);
                }
            }
        });
    });
    
    // 监听整个 cardList
    const cardList = document.getElementById('cardList');
    if (cardList) {
        observer.observe(cardList, {
            childList: true,
            subtree: true,
            attributes: true,
            attributeFilter: ['data-image-path']
        });
    }
    
    // 为已存在的图片添加删除按钮
    document.querySelectorAll('.img-upload-btn[data-image-path]').forEach(imgBox => {
        addDeleteButtonToImageBox(imgBox);
    });
}

// 导出函数到全局
window.deleteCardImage = deleteCardImage;
window.deleteAllCardImages = deleteAllCardImages;
window.getPendingDeleteImages = getPendingDeleteImages;
window.clearPendingDeleteImages = clearPendingDeleteImages;

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    initImageDeleteButtons();
    console.log('✓ 图片删除功能已加载');
});

console.log('✓ delete-images.js 已加载');
