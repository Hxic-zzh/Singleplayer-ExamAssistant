// 为卡片图片区域添加删除按钮事件处理

/**
 * 初始化卡片的图片删除按钮
 */
function initCardImageButtons() {
    console.log('🔧 开始初始化卡片图片按钮...');
    
    // 使用事件委托监听所有 "清空图片" 按钮
    document.addEventListener('click', async function(e) {
        // 检查是否点击了"清空图片"按钮
        if (e.target.classList.contains('img-clear-btn') || 
            e.target.closest('.img-clear-btn')) {
            
            const btn = e.target.classList.contains('img-clear-btn') ? 
                        e.target : e.target.closest('.img-clear-btn');
            
            // 找到所属的卡片
            const card = btn.closest('.single-card, .multiple-card, .fill-card, .material-card');
            
            if (!card) {
                console.warn('未找到所属卡片');
                return;
            }
            
            // 调用删除所有图片功能
            if (typeof window.deleteAllCardImages === 'function') {
                await window.deleteAllCardImages(card);
            } else {
                console.error('❌ deleteAllCardImages 函数未定义');
                alert('删除功能未加载，请刷新页面');
            }
        }
    });
    
    console.log('✓ 卡片图片按钮已初始化（使用事件委托）');
}

// 页面加载完成后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initCardImageButtons);
} else {
    initCardImageButtons();
}

console.log('✓ card-image-buttons.js 已加载');
