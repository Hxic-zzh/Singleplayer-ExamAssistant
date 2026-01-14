// 渲染图片下拉菜单，显示已上传的图片列表

/**
 * 更新下拉菜单内容，显示已上传的图片列表
 */
function updateImageDropdown(dropdown, imageListContainer) {
    if (!dropdown || !imageListContainer) return;
    
    // 获取所有已上传的图片
    const imgBoxes = imageListContainer.querySelectorAll('.img-upload-btn[data-image-path]');
    
    if (imgBoxes.length === 0) {
        dropdown.innerHTML = '<div style="padding:8px;color:#999;font-size:13px;">暂无已上传图片</div>';
        return;
    }
    
    dropdown.innerHTML = '';
      imgBoxes.forEach((imgBox, index) => {
        const imagePath = imgBox.dataset.imagePath;
        
        if (!imagePath) return;
        
        const item = document.createElement('div');
        item.className = 'dropdown-item';
        item.style.cssText = `
            padding: 10px 16px;
            border-bottom: 1px solid #eee;
            cursor: pointer;
            transition: background 0.2s;
            font-size: 14px;
            color: #333;
        `;
        
        // 只显示文字：图片序号和文件名
        const fileName = imagePath.split('/').pop() || imagePath.split('\\').pop() || '图片';
        item.textContent = `图片 ${index + 1}: ${fileName}`;
        
        // 悬停效果
        item.onmouseenter = function() {
            item.style.background = '#e3f2fd';
            item.style.color = '#1976d2';
        };
        item.onmouseleave = function() {
            item.style.background = 'transparent';
            item.style.color = '#333';
        };
          // 点击打开 lightbox 预览
        item.onclick = function(e) {
            e.stopPropagation();
            
            // 使用 previewPath (存储在 dataset 中的预览路径)
            const previewPath = imgBox.dataset.previewPath;
            
            // 打开 lightbox
            const lightbox = document.getElementById('img-lightbox');
            const lightboxImg = document.getElementById('lightbox-img');
            
            if (lightbox && lightboxImg) {
                // previewPath 已经是完整路径，例如：../tempwails/SCIMG_2_1.png
                lightboxImg.src = previewPath || '../tempwails/' + imagePath;
                lightbox.style.display = 'flex';
                console.log('打开 lightbox，图片路径:', previewPath || imagePath);
            } else {
                console.warn('Lightbox 元素不存在');
            }
            
            // 关闭下拉菜单
            dropdown.style.display = 'none';
        };
        
        dropdown.appendChild(item);
    });
}

/**
 * 为所有图片区域添加下拉菜单更新功能
 */
function initImageDropdownRenderers() {
    console.log('🔧 初始化图片下拉菜单渲染器...');
    
    // 使用 MutationObserver 监听图片上传
    const observer = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            // 监听 data-image-path 属性变化（图片上传完成）
            if (mutation.type === 'attributes' && mutation.attributeName === 'data-image-path') {
                const imgBox = mutation.target;
                const card = imgBox.closest('.single-card, .multiple-card, .fill-card, .material-card');
                
                if (card) {
                    // 更新对应区域的所有下拉菜单
                    updateAllDropdownsInCard(card);
                }
            }
            
            // 监听图片框的删除
            if (mutation.type === 'childList' && mutation.removedNodes.length > 0) {
                mutation.removedNodes.forEach(node => {
                    if (node.classList && node.classList.contains('img-upload-btn')) {
                        const card = mutation.target.closest('.single-card, .multiple-card, .fill-card, .material-card');
                        if (card) {
                            updateAllDropdownsInCard(card);
                        }
                    }
                });
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
    
    console.log('✓ 图片下拉菜单渲染器已启动');
}

/**
 * 更新卡片中所有下拉菜单
 */
function updateAllDropdownsInCard(card) {
    // 更新附带图片的下拉菜单
    const imgDropdown = card.querySelector('.img-dropdown');
    const imgList = card.querySelector('.img-list');
    if (imgDropdown && imgList) {
        updateImageDropdown(imgDropdown, imgList);
    }
    
    // 更新题干图片的下拉菜单
    const stemImgDropdown = card.querySelector('.stem-img-dropdown');
    const stemImgList = card.querySelector('.stem-img-list');
    if (stemImgDropdown && stemImgList) {
        updateImageDropdown(stemImgDropdown, stemImgList);
    }
}

// 页面加载完成后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initImageDropdownRenderers);
} else {
    initImageDropdownRenderers();
}

console.log('✓ image-dropdown-renderer.js 已加载');
