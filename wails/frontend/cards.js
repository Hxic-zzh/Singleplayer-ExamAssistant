// 公共工具函数 - 简化版本，图片管理由 cards-enhancement.js 接管
function createImgUploadArea({listClass, addBtnClass, listBtnClass, dropdownClass, maxCount = 4}) {
    let imgCount = 1;
    
    function renderImgs(listElem) {
        // 不要清空整个列表！只检查并添加缺少的图片框
        const existingBoxes = listElem.querySelectorAll('.img-upload-btn');
        const currentCount = existingBoxes.length;
        
        // 只添加新增的图片框
        for(let i = currentCount; i < imgCount; i++) {
            const imgBox = document.createElement('label');
            imgBox.className = 'img-upload-btn';
            imgBox.style = 'display:inline-block;width:120px;height:80px;border:2px dashed #bbb;border-radius:6px;cursor:pointer;text-align:center;line-height:80px;color:#aaa;font-size:15px;background:#f8f8f8;position:relative;user-select:none;';
            imgBox.innerHTML = `+ 添加图片<input type="file" accept="image/*" style="display:none;">`;
            listElem.appendChild(imgBox);
        }
        
        // 如果需要减少图片框（删除操作），只删除空白的框
        if (currentCount > imgCount) {
            for(let i = currentCount - 1; i >= imgCount; i--) {
                const box = existingBoxes[i];
                // 只删除没有上传图片的空白框
                if (!box.dataset.imagePath) {
                    box.remove();
                }
            }
        }
    }
      return {
        setup: function(card) {
            const listElem = card.querySelector(listClass);
            const addBtn = card.querySelector(addBtnClass);
            const listBtn = card.querySelector(listBtnClass);
            const dropdownElem = card.querySelector(dropdownClass);
            
            renderImgs(listElem);
            
            addBtn.onclick = function() {
                if(imgCount<maxCount) { 
                    imgCount++; 
                    renderImgs(listElem);
                }
            };
            
            // 修复：使用独立的事件处理器，避免全局冲突
            listBtn.onclick = function(e) {
                e.preventDefault();
                e.stopPropagation();
                
                // 先关闭其他所有下拉菜单
                document.querySelectorAll('.img-dropdown, .stem-img-dropdown').forEach(dd => {
                    if (dd !== dropdownElem) {
                        dd.style.display = 'none';
                    }
                });
                
                // 切换当前下拉菜单
                dropdownElem.style.display = dropdownElem.style.display === 'block' ? 'none' : 'block';
            };
            
            // 点击下拉菜单内部不关闭
            dropdownElem.onclick = function(e) {
                e.stopPropagation();
            };
            
            // 使用命名函数，方便移除监听器（避免重复添加）
            const closeDropdown = function(e) {
                if (!listBtn.contains(e.target) && !dropdownElem.contains(e.target)) {
                    dropdownElem.style.display = 'none';
                }
            };
            
            // 延迟添加全局点击监听（避免立即触发）
            setTimeout(() => {
                document.addEventListener('click', closeDropdown);
            }, 0);
        }
    };
}

// 单选题干无图
function createSingleChoiceCard(cardIndex) {
    const card = document.createElement('div');
    card.className = 'single-card sc';
    card.style = '';
    card.innerHTML = `
        <div style="margin-bottom:12px;display:flex;align-items:center;justify-content:space-between;">
            <div><span class="card-type-label sc">[SC]</span><span style="font-weight:bold;">单选题干无图</span></div>
            <div style="font-size:15px;color:#888;">题型序号：<span class="card-index">${cardIndex}</span></div>
        </div>
        <textarea class="stem-input" style="width:90%;height:60px;padding:6px 8px;font-size:15px;resize:vertical;margin-bottom:12px;" placeholder="请输入题干内容"></textarea>
        <div style="margin:10px 0 16px 0;">
            <button class="option-add">+ 选项</button>
            <button class="option-remove">- 选项</button>
        </div>
        <div class="options-area"></div>        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">附带的题目参考图片：</div>
            <div class="img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="img-list-btn">图片列表 ▲</button>
                    <div class="img-dropdown" style="display:none;position:absolute;left:0;bottom:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;max-width:200px;max-height:300px;overflow-y:auto;z-index:1000;"></div>
                </div>
            </div>
        </div>
        <button class="card-delete-btn">🗑 删除本题</button>
    `;
    // 选项逻辑
    const optionsArea = card.querySelector('.options-area');
    let optionCount = 2;
    const optionLabels = 'ABCDEFGHIJK'.split('');
    function renderOptions() {
        optionsArea.innerHTML = '';
        const radioName = 'single-choice-' + Date.now() + Math.random();
        for(let i=0;i<optionCount;i++) {
            const optDiv = document.createElement('div');
            optDiv.style = 'margin-bottom:8px;display:flex;align-items:center;gap:8px;';
            optDiv.innerHTML = `
                <input type="radio" name="${radioName}" style="margin-right:4px;">
                <span style="width:22px;display:inline-block;text-align:center;font-weight:bold;">${optionLabels[i]}</span>
                <input type="text" class="option-input" style="width:60%;padding:5px 8px;font-size:15px;" placeholder="请输入选项内容">
            `;
            optionsArea.appendChild(optDiv);
        }
    }
    renderOptions();
    card.querySelector('.option-add').onclick = function() {
        if(optionCount<11) { optionCount++; renderOptions(); }
    };    card.querySelector('.option-remove').onclick = function() {
        if(optionCount>2) { optionCount--; renderOptions(); }
    };
    // 附带图片逻辑
    createImgUploadArea({
        listClass: '.img-list',
        addBtnClass: '.img-add-btn',
        listBtnClass: '.img-list-btn',
        dropdownClass: '.img-dropdown'
    }).setup(card);

    const deleteBtn = card.querySelector('.card-delete-btn');
    deleteBtn.onclick = function () {
        const evt = new CustomEvent('card-delete', { detail: { card } });
        window.dispatchEvent(evt);
    };
    return card;
}
window.createSingleChoiceCard = createSingleChoiceCard;

// 多选题干无图
function createMultipleChoiceCard(cardIndex) {
    const card = document.createElement('div');
    card.className = 'multiple-card mc';
    card.style = '';
    card.innerHTML = `
        <div style="margin-bottom:12px;display:flex;align-items:center;justify-content:space-between;">
            <div><span class="card-type-label mc">[MC]</span><span style="font-weight:bold;">多选题干无图</span></div>
            <div style="font-size:15px;color:#888;">题型序号：<span class="card-index">${cardIndex}</span></div>
        </div>
        <textarea class="stem-input" style="width:90%;height:60px;padding:6px 8px;font-size:15px;resize:vertical;margin-bottom:12px;" placeholder="请输入题干内容"></textarea>
        <div style="margin:10px 0 16px 0;">
            <button class="option-add">+ 选项</button>
            <button class="option-remove">- 选项</button>
        </div>
        <div class="options-area"></div>        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">附带的题目参考图片：</div>
            <div class="img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="img-list-btn">图片列表 ▲</button>
                    <div class="img-dropdown" style="display:none;position:absolute;left:0;bottom:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;max-width:200px;max-height:300px;overflow-y:auto;z-index:1000;"></div>
                </div>
            </div>
        </div>
        <button class="card-delete-btn">🗑 删除本题</button>
    `;
    // 选项逻辑
    const optionsArea = card.querySelector('.options-area');
    let optionCount = 2;
    const optionLabels = 'ABCDEFGHIJK'.split('');
    function renderOptions() {
        optionsArea.innerHTML = '';
        const checkboxName = 'multiple-choice-' + Date.now() + Math.random();
        for(let i=0;i<optionCount;i++) {
            const optDiv = document.createElement('div');
            optDiv.style = 'margin-bottom:8px;display:flex;align-items:center;gap:8px;';
            optDiv.innerHTML = `
                <input type="checkbox" name="${checkboxName}" style="margin-right:4px;">
                <span style="width:22px;display:inline-block;text-align:center;font-weight:bold;">${optionLabels[i]}</span>
                <input type="text" class="option-input" style="width:60%;padding:5px 8px;font-size:15px;" placeholder="请输入选项内容">
            `;
            optionsArea.appendChild(optDiv);
        }
    }
    renderOptions();
    card.querySelector('.option-add').onclick = function() {
        if(optionCount<11) { optionCount++; renderOptions(); }
    };    card.querySelector('.option-remove').onclick = function() {
        if(optionCount>2) { optionCount--; renderOptions(); }
    };
    // 附带图片逻辑
    createImgUploadArea({
        listClass: '.img-list',
        addBtnClass: '.img-add-btn',
        listBtnClass: '.img-list-btn',
        dropdownClass: '.img-dropdown'
    }).setup(card);

    const deleteBtn = card.querySelector('.card-delete-btn');
    deleteBtn.onclick = function () {
        const evt = new CustomEvent('card-delete', { detail: { card } });
        window.dispatchEvent(evt);
    };
    return card;
}
window.createMultipleChoiceCard = createMultipleChoiceCard;

// 单选题干有图
function createSingleChoiceWithStemImgCard(cardIndex) {
    const card = document.createElement('div');
    card.className = 'single-card scimg';
    card.style = '';
    card.innerHTML = `
        <div style="margin-bottom:12px;display:flex;align-items:center;justify-content:space-between;">
            <div><span class="card-type-label scimg">[SCIMG]</span><span style="font-weight:bold;">单选题干有图</span></div>
            <div style="font-size:15px;color:#888;">题型序号：<span class="card-index">${cardIndex}</span></div>
        </div>
        <textarea class="stem-input" style="width:90%;height:60px;padding:6px 8px;font-size:15px;resize:vertical;margin-bottom:12px;" placeholder="请输入题干内容"></textarea>
        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">题干图片：</div>
            <div class="stem-img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="stem-img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="stem-img-list-btn">图片列表 ▼</button>
                    <div class="stem-img-dropdown" style="display:none;position:absolute;left:0;top:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;z-index:10;"></div>
                </div>
            </div>
        </div>
        <div style="margin:10px 0 16px 0;">
            <button class="option-add">+ 选项</button>
            <button class="option-remove">- 选项</button>
        </div>
        <div class="options-area"></div>
        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">附带的题目参考图片：</div>
            <div class="img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="img-list-btn">图片列表 ▼</button>
                    <div class="img-dropdown" style="display:none;position:absolute;left:0;top:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;z-index:10;"></div>
                </div>
            </div>
        </div>
        <button class="card-delete-btn">🗑 删除本题</button>    `;
    // 题干图片逻辑
    createImgUploadArea({
        listClass: '.stem-img-list',
        addBtnClass: '.stem-img-add-btn',
        listBtnClass: '.stem-img-list-btn',
        dropdownClass: '.stem-img-dropdown'
    }).setup(card);
    // 选项逻辑
    const optionsArea = card.querySelector('.options-area');
    let optionCount = 2;
    const optionLabels = 'ABCDEFGHIJK'.split('');
    function renderOptions() {
        optionsArea.innerHTML = '';
        const radioName = 'single-choice-' + Date.now() + Math.random();
        for(let i=0;i<optionCount;i++) {
            const optDiv = document.createElement('div');
            optDiv.style = 'margin-bottom:8px;display:flex;align-items:center;gap:8px;';
            optDiv.innerHTML = `
                <input type="radio" name="${radioName}" style="margin-right:4px;">
                <span style="width:22px;display:inline-block;text-align:center;font-weight:bold;">${optionLabels[i]}</span>
                <input type="text" class="option-input" style="width:60%;padding:5px 8px;font-size:15px;" placeholder="请输入选项内容">
            `;
            optionsArea.appendChild(optDiv);
        }
    }
    renderOptions();
    card.querySelector('.option-add').onclick = function() {
        if(optionCount<11) { optionCount++; renderOptions(); }
    };
    card.querySelector('.option-remove').onclick = function() {
        if(optionCount>2) { optionCount--; renderOptions(); }
    };    // 附带图片逻辑
    createImgUploadArea({
        listClass: '.img-list',
        addBtnClass: '.img-add-btn',
        listBtnClass: '.img-list-btn',
        dropdownClass: '.img-dropdown'
    }).setup(card);

    const deleteBtn = card.querySelector('.card-delete-btn');
    deleteBtn.onclick = function () {
        const evt = new CustomEvent('card-delete', { detail: { card } });
        window.dispatchEvent(evt);
    };
    return card;
}
window.createSingleChoiceWithStemImgCard = createSingleChoiceWithStemImgCard;

// 多选题干有图
function createMultipleChoiceWithStemImgCard(cardIndex) {
    const card = document.createElement('div');
    card.className = 'multiple-card mcimg';
    card.style = '';
    card.innerHTML = `
        <div style="margin-bottom:12px;display:flex;align-items:center;justify-content:space-between;">
            <div><span class="card-type-label mcimg">[MCIMG]</span><span style="font-weight:bold;">多选题干有图</span></div>
            <div style="font-size:15px;color:#888;">题型序号：<span class="card-index">${cardIndex}</span></div>
        </div>
        <textarea class="stem-input" style="width:90%;height:60px;padding:6px 8px;font-size:15px;resize:vertical;margin-bottom:12px;" placeholder="请输入题干内容"></textarea>
        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">题干图片：</div>
            <div class="stem-img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="stem-img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="stem-img-list-btn">图片列表 ▼</button>
                    <div class="stem-img-dropdown" style="display:none;position:absolute;left:0;top:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;z-index:10;"></div>
                </div>
            </div>
        </div>
        <div style="margin:10px 0 16px 0;">
            <button class="option-add">+ 选项</button>
            <button class="option-remove">- 选项</button>
        </div>
        <div class="options-area"></div>
        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">附带的题目参考图片：</div>
            <div class="img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="img-list-btn">图片列表 ▼</button>
                    <div class="img-dropdown" style="display:none;position:absolute;left:0;top:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;z-index:10;"></div>
                </div>
            </div>
        </div>
        <button class="card-delete-btn">🗑 删除本题</button>
    `;    // 题干图片逻辑
    createImgUploadArea({
        listClass: '.stem-img-list',
        addBtnClass: '.stem-img-add-btn',
        listBtnClass: '.stem-img-list-btn',
        dropdownClass: '.stem-img-dropdown'
    }).setup(card);

    // 选项逻辑（多选）
    const optionsArea = card.querySelector('.options-area');
    let optionCount = 2;
    const optionLabels = 'ABCDEFGHIJK'.split('');
    function renderOptions() {
        optionsArea.innerHTML = '';
        const checkboxName = 'multiple-choice-img-' + Date.now() + Math.random();
        for (let i = 0; i < optionCount; i++) {
            const optDiv = document.createElement('div');
            optDiv.style = 'margin-bottom:8px;display:flex;align-items:center;gap:8px;';
            optDiv.innerHTML = `
                <input type="checkbox" name="${checkboxName}" style="margin-right:4px;">
                <span style="width:22px;display:inline-block;text-align:center;font-weight:bold;">${optionLabels[i]}</span>
                <input type="text" class="option-input" style="width:60%;padding:5px 8px;font-size:15px;" placeholder="请输入选项内容">
            `;
            optionsArea.appendChild(optDiv);
        }
    }
    renderOptions();
    card.querySelector('.option-add').onclick = function() {
        if (optionCount < 11) { optionCount++; renderOptions(); }
    };
    card.querySelector('.option-remove').onclick = function() {
        if (optionCount > 2) { optionCount--; renderOptions(); }
    };    // 附带图片逻辑
    createImgUploadArea({
        listClass: '.img-list',
        addBtnClass: '.img-add-btn',
        listBtnClass: '.img-list-btn',
        dropdownClass: '.img-dropdown'
    }).setup(card);

    const deleteBtn = card.querySelector('.card-delete-btn');
    deleteBtn.onclick = function () {
        const evt = new CustomEvent('card-delete', { detail: { card } });
        window.dispatchEvent(evt);
    };
    return card;
}
window.createMultipleChoiceWithStemImgCard = createMultipleChoiceWithStemImgCard;

// ===== 填空题干无图 FL =====
function createFillBlankCard(cardIndex) {
    const card = document.createElement('div');
    card.className = 'fill-card fl';
    card.style = '';
    card.innerHTML = `
        <div style="margin-bottom:12px;display:flex;align-items:center;justify-content:space-between;">
            <div><span class="card-type-label fl">[FL]</span><span style="font-weight:bold;">填空题干无图</span></div>
            <div style="font-size:15px;color:#888;">题型序号：<span class="card-index">${cardIndex}</span></div>
        </div>
        <div style="margin-bottom:8px;">
            <textarea class="stem-input" style="width:90%;height:60px;padding:6px 8px;font-size:15px;resize:vertical;" placeholder="请输入题干内容，空用(%___%)表示，例如：1+1等于(%___%)"></textarea>
        </div>
        <div style="margin-bottom:10px;display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
            <button class="blank-sync-btn">同步空位</button>
            <button class="blank-check-btn">检查填空</button>
            <span style="font-size:12px;color:#888;">识别符：(%___%)，按出现顺序对应空1、空2…</span>
        </div>
        <div class="blank-config-area" style="border-top:1px dashed #ddd;padding-top:8px;margin-bottom:12px;"></div>
        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">附带的题目参考图片：</div>
            <div class="img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="img-list-btn">图片列表 ▼</button>
                    <div class="img-dropdown" style="display:none;position:absolute;left:0;top:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;z-index:10;"></div>
                </div>
            </div>
        </div>
        <button class="card-delete-btn">🗑 删除本题</button>
    `;    setupFillBlankLogic(card);

    // 附带图片逻辑
    createImgUploadArea({
        listClass: '.img-list',
        addBtnClass: '.img-add-btn',
        listBtnClass: '.img-list-btn',
        dropdownClass: '.img-dropdown'
    }).setup(card);

    const deleteBtn = card.querySelector('.card-delete-btn');
    deleteBtn.onclick = function () {
        const evt = new CustomEvent('card-delete', { detail: { card } });
        window.dispatchEvent(evt);
    };
    return card;
}
window.createFillBlankCard = createFillBlankCard;

// ===== 填空题干有图 FLIMG =====
function createFillBlankWithStemImgCard(cardIndex) {
    const card = document.createElement('div');
    card.className = 'fill-card flimg';
    card.style = '';
    card.innerHTML = `
        <div style="margin-bottom:12px;display:flex;align-items:center;justify-content:space-between;">
            <div><span class="card-type-label flimg">[FLIMG]</span><span style="font-weight:bold;">填空题干有图</span></div>
            <div style="font-size:15px;color:#888;">题型序号：<span class="card-index">${cardIndex}</span></div>
        </div>
        <div style="margin-bottom:8px;">
            <textarea class="stem-input" style="width:90%;height:60px;padding:6px 8px;font-size:15px;resize:vertical;" placeholder="请输入题干内容，空用(%___%)表示，例如：1+1等于(%___%)"></textarea>
        </div>
        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">题干图片：</div>
            <div class="stem-img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="stem-img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="stem-img-list-btn">图片列表 ▼</button>
                    <div class="stem-img-dropdown" style="display:none;position:absolute;left:0;top:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;z-index:10;"></div>
                </div>
            </div>
        </div>
        <div style="margin-bottom:10px;display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
            <button class="blank-sync-btn">同步空位</button>
            <button class="blank-check-btn">检查填空</button>
            <span style="font-size:12px;color:#888;">识别符：(%___%)，按出现顺序对应空1、空2…</span>
        </div>
        <div class="blank-config-area" style="border-top:1px dashed #ddd;padding-top:8px;margin-bottom:12px;"></div>
        <div style="margin-bottom:16px;">
            <div style="font-weight:bold;margin-bottom:6px;">附带的题目参考图片：</div>
            <div class="img-list" style="display:flex;gap:12px;flex-wrap:wrap;"></div>
            <div style="margin-top:8px;display:flex;align-items:center;gap:10px;">
                <button class="img-add-btn">+ 添加图片</button>
                <div style="position:relative;display:inline-block;">
                    <button class="img-list-btn">图片列表 ▼</button>
                    <div class="img-dropdown" style="display:none;position:absolute;left:0;top:110%;background:#fff;border:1px solid #ccc;box-shadow:0 2px 8px #0002;border-radius:4px;min-width:120px;z-index:10;"></div>
                </div>
            </div>
        </div>
        <button class="card-delete-btn">🗑 删除本题</button>
    `;    // 题干图片逻辑
    createImgUploadArea({
        listClass: '.stem-img-list',
        addBtnClass: '.stem-img-add-btn',
        listBtnClass: '.stem-img-list-btn',
        dropdownClass: '.stem-img-dropdown'
    }).setup(card);

    setupFillBlankLogic(card);    // 附带图片逻辑
    createImgUploadArea({
        listClass: '.img-list',
        addBtnClass: '.img-add-btn',
        listBtnClass: '.img-list-btn',
        dropdownClass: '.img-dropdown'
    }).setup(card);

    const deleteBtn = card.querySelector('.card-delete-btn');
    deleteBtn.onclick = function () {
        const evt = new CustomEvent('card-delete', { detail: { card } });
        window.dispatchEvent(evt);
    };
    return card;
}
window.createFillBlankWithStemImgCard = createFillBlankWithStemImgCard;

// 公共：填空题逻辑（识别(%___%)、配置空答案、检查）
function setupFillBlankLogic(card, initialBlanks = null) {
    const stemInput = card.querySelector('.stem-input');
    const blankArea = card.querySelector('.blank-config-area');
    const syncBtn = card.querySelector('.blank-sync-btn');
    const checkBtn = card.querySelector('.blank-check-btn');

    // 内部状态：每个空一个配置对象
    let blanks = initialBlanks || []; // [{ answers: ['a','b'], unique: false }, ...]

    function parseBlankCountFromStem() {
        const text = stemInput.value || '';
        const matches = text.match(/\(%___%\)/g);
        return matches ? matches.length : 0;
    }

    function rebuildBlankUI() {
        blankArea.innerHTML = '';
        blanks.forEach((blank, idx) => {
            const index = idx + 1;
            const block = document.createElement('div');
            block.style = 'margin-bottom:10px;padding:8px 10px;border-radius:8px;background:#f9fafb;border:1px solid #e5e7eb;';            block.innerHTML = `
                <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:6px;">
                    <div style="font-size:13px;font-weight:600;">空${index}</div>
                    <label style="font-size:12px;color:#555;display:flex;align-items:center;gap:4px;">
                        <input type="checkbox" class="blank-unique" ${blank.unique ? 'checked' : ''}>
                        <span>学生答案不可重复使用</span>
                    </label>
                </div>
                <div class="blank-answers"></div>
                <button class="blank-add-answer-btn" style="margin-top:6px;font-size:12px;">+ 添加备选答案</button>
            `;const answersContainer = block.querySelector('.blank-answers');
            function renderAnswers() {
                answersContainer.innerHTML = '';
                if (!blank.answers || blank.answers.length === 0) {
                    blank.answers = [''];
                }
                // 只渲染非 (x%x) 标记的答案
                const visibleAnswers = blank.answers.filter(a => a !== '(x%x)');
                if (visibleAnswers.length === 0) {
                    visibleAnswers.push('');
                }
                visibleAnswers.forEach((ans, displayIdx) => {
                    // 找到实际在 blank.answers 中的索引
                    const aIdx = blank.answers.indexOf(ans, displayIdx > 0 ? blank.answers.indexOf(visibleAnswers[displayIdx - 1]) + 1 : 0);
                    const row = document.createElement('div');
                    row.style = 'display:flex;align-items:center;gap:6px;margin-bottom:4px;';
                    row.innerHTML = `
                        <input type="text" class="blank-answer-input" style="flex:1;padding:4px 8px;font-size:14px;" placeholder="空${index} 的一个可能答案" value="${ans.replace(/"/g, '&quot;')}">
                        <button class="blank-del-answer-btn" style="font-size:12px;">-</button>
                    `;
                    const input = row.querySelector('.blank-answer-input');
                    const delBtn = row.querySelector('.blank-del-answer-btn');
                    input.oninput = function() {
                        blank.answers[aIdx] = input.value;
                    };
                    delBtn.onclick = function() {
                        if (visibleAnswers.length > 1) {
                            blank.answers.splice(aIdx, 1);
                            renderAnswers();
                        } else {
                            blank.answers[aIdx] = '';
                            renderAnswers();
                        }
                    };
                    answersContainer.appendChild(row);
                });
            }
            renderAnswers();

            block.querySelector('.blank-add-answer-btn').onclick = function() {
                blank.answers.push('');
                renderAnswers();
            };            const uniqueCheckbox = block.querySelector('.blank-unique');
            uniqueCheckbox.onchange = function() {
                blank.unique = uniqueCheckbox.checked;
                // 添加/移除 (x%x) 标记
                if (uniqueCheckbox.checked) {
                    if (!blank.answers.includes('(x%x)')) {
                        // 将 (x%x) 标记添加到数组开头，与导出逻辑保持一致
                        blank.answers.unshift('(x%x)');
                    }
                } else {
                    blank.answers = blank.answers.filter(a => a !== '(x%x)');
                }
                renderAnswers();
            };

            blankArea.appendChild(block);
        });
    }

    function syncBlanksToStem() {
        const count = parseBlankCountFromStem();
        if (count === 0) {
            alert('题干中未检测到 (%___%)，请先在题干中用该符号标记空位。');
            return;
        }
        if (count === blanks.length) {
            alert('题干中检测到 ' + count + ' 个空，与当前配置数量一致。');
            return;
        }
        const newBlanks = [];
        for (let i = 0; i < count; i++) {
            newBlanks.push(blanks[i] || { answers: [''], unique: false });
        }
        blanks = newBlanks;
        rebuildBlankUI();
        alert('已根据题干中的空位数量同步为空 ' + count + ' 个。');
    }

    function checkBlanks() {
        const countInStem = parseBlankCountFromStem();
        if (countInStem === 0) {
            alert('题干中没有任何 (%___%)，请先在题干中标记空位。');
            return;
        }
        if (blanks.length !== countInStem) {
            alert('题干中有 ' + countInStem + ' 个空，但仅配置了 ' + blanks.length + ' 个，请先点击"同步空位"。');
            return;
        }
        // 检查每个空至少有一个非空答案(排除 (x%x) 标记)
        for (let i = 0; i < blanks.length; i++) {
            const b = blanks[i];
            const hasNonEmpty = (b.answers || [])
                .filter(a => a !== '(x%x)') // 排除标记
                .some(a => (a || '').trim() !== '');
            if (!hasNonEmpty) {
                alert('空 ' + (i + 1) + ' 未设置任何有效答案。');
                return;
            }
        }
        alert('填空检查通过：题干空位数量与配置一致，所有空均有答案。');
    }    syncBtn.onclick = syncBlanksToStem;
    checkBtn.onclick = checkBlanks;
    
    // 如果提供了初始数据，立即重建UI
    if (initialBlanks && initialBlanks.length > 0) {
        rebuildBlankUI();
    }
}
// 将 setupFillBlankLogic 暴露到全局作用域，供 export.js 使用
window.setupFillBlankLogic = setupFillBlankLogic;

// 材料题容器卡片
function createMaterialCard(cardIndex) {
    const card = document.createElement('div');
    card.className = 'material-card mt';
    card.innerHTML = `
        <div class="card-header-line">
            <div>
                <span class="card-type-label mt">[DR]</span>
                <span class="card-title-text">材料题</span>
            </div>
            <div class="card-index-line">题型序号：<span class="card-index">${cardIndex}</span></div>
        </div>
        <div class="material-block" style="margin-bottom:10px;">
            <div class="material-title">材料内容（必填）</div>
            <textarea class="material-input" placeholder="请在此输入材料全文，支持多行"></textarea>
        </div>
        <div class="material-img-block" style="margin-bottom:10px;">
            <div class="material-title">材料图片（可选）</div>
            <div class="material-img-list"></div>
            <div class="material-img-toolbar">
                <button class="material-img-add-btn">+ 添加图片</button>
                <div class="material-img-dropdown-wrap">
                    <button class="material-img-list-btn">图片列表 ▼</button>
                    <div class="material-img-dropdown"></div>
                </div>
            </div>
        </div>
        <div class="material-check-line">
            <button class="material-check-btn">检查材料</button>
            <span>请先完善材料内容并点击“检查材料”，再添加子题。</span>
        </div>
        <div class="material-inner-toolbar" style="display:none;">
            <span class="inner-toolbar-label">添加子题：</span>
            <button class="mt-inner-btn mt-inner-sc">单选无图</button>
            <button class="mt-inner-btn mt-inner-scimg">单选有图</button>
            <button class="mt-inner-btn mt-inner-mc">多选无图</button>
            <button class="mt-inner-btn mt-inner-mcimg">多选有图</button>
            <button class="mt-inner-btn mt-inner-fl">填空无图</button>
            <button class="mt-inner-btn mt-inner-flimg">填空有图</button>
        </div>
        <div class="material-inner-list"></div>
        <button class="card-delete-btn">🗑 删除整套材料题</button>
    `;    // 材料图片上传逻辑
    createImgUploadArea({
        listClass: '.material-img-list',
        addBtnClass: '.material-img-add-btn',
        listBtnClass: '.material-img-list-btn',
        dropdownClass: '.material-img-dropdown'
    }).setup(card);

    let materialChecked = false;
    const checkBtn = card.querySelector('.material-check-btn');
    const materialInput = card.querySelector('.material-input');
    const innerToolbar = card.querySelector('.material-inner-toolbar');

    checkBtn.onclick = function() {
        if (!materialInput.value.trim()) {
            alert('请先填写材料内容。');
            return;
        }
        materialChecked = true;
        checkBtn.disabled = true;
        checkBtn.textContent = '材料已检查';
        innerToolbar.style.display = 'flex';
    };

    const innerList = card.querySelector('.material-inner-list');
    let innerIndex = 1;

    function wrapAsInnerCard(factoryFn, typeLabelText) {
        if (!materialChecked) {
            alert('请先填写并检查材料，再添加子题。');
            return;
        }
        const innerCard = factoryFn(innerIndex);
        // 去掉外层删除按钮，避免触发全局删除整张题卡
        const delBtn = innerCard.querySelector('.card-delete-btn');
        if (delBtn) delBtn.remove();
        const idxSpan = innerCard.querySelector('.card-index');
        if (idxSpan) idxSpan.textContent = innerIndex;

        const wrapper = document.createElement('div');
        wrapper.className = 'mt-inner-card';
        wrapper.innerHTML = `
            <div class="mt-inner-header">
                <span class="mt-inner-tag">内题 ${innerIndex}</span>
                <span class="mt-inner-title">${typeLabelText}</span>
                <button class="mt-inner-delete-btn" style="margin-left:auto;">🗑 删除子题</button>
            </div>
        `;
        wrapper.appendChild(innerCard);

        // 子题删除按钮：派发专用事件，交由全局处理图片+DOM
        const innerDeleteBtn = wrapper.querySelector('.mt-inner-delete-btn');
        innerDeleteBtn.onclick = function () {
            const evt = new CustomEvent('mt-inner-delete', { detail: { wrapper } });
            window.dispatchEvent(evt);
        };

        innerList.appendChild(wrapper);
        innerIndex++;
    }

    card.querySelector('.mt-inner-sc').onclick = () => wrapAsInnerCard(window.createSingleChoiceCard, '单选题干无图');
    card.querySelector('.mt-inner-scimg').onclick = () => wrapAsInnerCard(window.createSingleChoiceWithStemImgCard, '单选题干有图');
    card.querySelector('.mt-inner-mc').onclick = () => wrapAsInnerCard(window.createMultipleChoiceCard, '多选题干无图');
    card.querySelector('.mt-inner-mcimg').onclick = () => wrapAsInnerCard(window.createMultipleChoiceWithStemImgCard, '多选题干有图');
    card.querySelector('.mt-inner-fl').onclick = () => wrapAsInnerCard(window.createFillBlankCard, '填空题干无图');
    card.querySelector('.mt-inner-flimg').onclick = () => wrapAsInnerCard(window.createFillBlankWithStemImgCard, '填空题干有图');

    const deleteBtn = card.querySelector('.card-delete-btn');
    deleteBtn.onclick = function () {
        const evt = new CustomEvent('card-delete', { detail: { card } });
        window.dispatchEvent(evt);
    };

    return card;
}
window.createMaterialCard = createMaterialCard;
