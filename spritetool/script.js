const ZOOM_FACTOR = 4; 
const LINE_WIDTH_MAGNIFIED = ZOOM_FACTOR; 

const canvas = document.getElementById('spritesheetCanvas');
const ctx = canvas.getContext('2d');
const magnifiedCanvas = document.getElementById('magnifiedCanvas');
const magnifiedCtx = magnifiedCanvas.getContext('2d');
const originDot = document.getElementById('originDot');

const fileInput = document.getElementById('spritesheetFile');
const gridColsInput = document.getElementById('gridCols');
const gridRowsInput = document.getElementById('gridRows');
const originXInput = document.getElementById('originX');
const originYInput = document.getElementById('originY');
const spriteSizeDisplay = document.getElementById('spriteSizeDisplay');
const selectedSpriteInfo = document.getElementById('selectedSpriteInfo');
const hitboxCoordinatesDiv = document.getElementById('hitboxCoordinates');

let spritesheetImage = null;
let gridCols = 4;
let gridRows = 4;
let spriteWidth = 0;
let spriteHeight = 0;
let originX = 8; 
let originY = 8;

let isDrawing = false;
let startX = 0;
let startY = 0;
let hitbox = { x: 0, y: 0, w: 0, h: 0 }; 
let selectedSpriteIndex = -1; 
let selectedSpriteCol = -1;
let selectedSpriteRow = -1;

// --- Initialization and Update Functions (Unchanged) ---

fileInput.addEventListener('change', (e) => {
    const file = e.target.files[0];
    if (file) {
        spritesheetImage = new Image();
        spritesheetImage.onload = () => {
            canvas.width = spritesheetImage.width;
            canvas.height = spritesheetImage.height;
            updateGrid(); 
        };
        spritesheetImage.src = URL.createObjectURL(file);
    }
});

function updateGrid() {
    if (!spritesheetImage) return;

    gridCols = parseInt(gridColsInput.value) || 1;
    gridRows = parseInt(gridRowsInput.value) || 1;

    spriteWidth = spritesheetImage.width / gridCols;
    spriteHeight = spritesheetImage.height / gridRows;

    magnifiedCanvas.width = spriteWidth * ZOOM_FACTOR;
    magnifiedCanvas.height = spriteHeight * ZOOM_FACTOR;
    
    selectedSpriteIndex = -1;
    hitbox = { x: 0, y: 0, w: 0, h: 0 };
    clearMagnifiedCanvas();

    spriteSizeDisplay.textContent = `${spriteWidth.toFixed(0)}px (W) x ${spriteHeight.toFixed(0)}px (H)`;
    drawSpritesheet();
    updateOrigin(); 
}

function updateOrigin() {
    originX = parseInt(originXInput.value) || 0;
    originY = parseInt(originYInput.value) || 0;
    drawSpritesheet();
    drawMagnifiedSprite(); 
}

function clearMagnifiedCanvas() {
    magnifiedCtx.clearRect(0, 0, magnifiedCanvas.width, magnifiedCanvas.height);
    originDot.style.display = 'none';
}

function drawSpritesheet() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    if (!spritesheetImage) return;

    ctx.drawImage(spritesheetImage, 0, 0);

    // Grid Overlay
    ctx.strokeStyle = 'rgba(255, 0, 0, 0.5)'; 
    ctx.lineWidth = 1;
    for (let i = 1; i < gridCols; i++) {
        ctx.beginPath();
        ctx.moveTo(i * spriteWidth, 0);
        ctx.lineTo(i * spriteWidth, canvas.height);
        ctx.stroke();
    }
    for (let i = 1; i < gridRows; i++) {
        ctx.beginPath();
        ctx.moveTo(0, i * spriteHeight);
        ctx.lineTo(canvas.width, i * spriteHeight);
        ctx.stroke();
    }

    // Highlight Selected Sprite
    if (selectedSpriteIndex !== -1) {
        const x = selectedSpriteCol * spriteWidth;
        const y = selectedSpriteRow * spriteHeight;
        ctx.strokeStyle = 'blue';
        ctx.lineWidth = 3;
        ctx.strokeRect(x, y, spriteWidth, spriteHeight);
    }
}

function drawMagnifiedSprite() {
    clearMagnifiedCanvas();

    if (selectedSpriteIndex === -1 || !spritesheetImage) return;

    const srcX = selectedSpriteCol * spriteWidth;
    const srcY = selectedSpriteRow * spriteHeight;
    const srcW = spriteWidth;
    const srcH = spriteHeight;

    const destX = 0;
    const destY = 0;
    const destW = spriteWidth * ZOOM_FACTOR;
    const destH = spriteHeight * ZOOM_FACTOR;

    // 1. Draw the magnified sprite
    magnifiedCtx.imageSmoothingEnabled = false; 
    magnifiedCtx.drawImage(
        spritesheetImage, 
        srcX, srcY, srcW, srcH, 
        destX, destY, destW, destH
    );

    // 2. Draw the Origin Marker
    const originDotX = originX * ZOOM_FACTOR;
    const originDotY = originY * ZOOM_FACTOR;
    
    originDot.style.left = `${originDotX}px`;
    originDot.style.top = `${originDotY}px`;
    originDot.style.display = 'block';

    // 3. Draw Current Hitbox
    if (hitbox.w !== 0 || hitbox.h !== 0) {
        magnifiedCtx.strokeStyle = 'lime';
        magnifiedCtx.lineWidth = LINE_WIDTH_MAGNIFIED; 
        
        // Pixel Alignment
        const halfLineWidth = LINE_WIDTH_MAGNIFIED / 2;

        magnifiedCtx.strokeRect(
            hitbox.x + halfLineWidth, 
            hitbox.y + halfLineWidth, 
            hitbox.w - LINE_WIDTH_MAGNIFIED, 
            hitbox.h - LINE_WIDTH_MAGNIFIED
        );
    }

    calculateHitboxRelativeCoords();
}

// --- Helper Function (Unchanged) ---
function getMousePos(canvas, evt) {
    const rect = canvas.getBoundingClientRect();
    return {
        x: evt.clientX - rect.left,
        y: evt.clientY - rect.top
    };
}

// --- Sprite Selection on Full Sheet Click (Unchanged) ---
canvas.addEventListener('click', (e) => {
    if (!spritesheetImage || spriteWidth === 0) return; 

    const pos = getMousePos(canvas, e);

    selectedSpriteCol = Math.floor(pos.x / spriteWidth);
    selectedSpriteRow = Math.floor(pos.y / spriteHeight);
    
    if (selectedSpriteCol >= 0 && selectedSpriteCol < gridCols && 
        selectedSpriteRow >= 0 && selectedSpriteRow < gridRows) {
        
        selectedSpriteIndex = selectedSpriteRow * gridCols + selectedSpriteCol;
        
        hitbox = { x: 0, y: 0, w: 0, h: 0 }; 

        selectedSpriteInfo.textContent = `Selected Sprite: Index ${selectedSpriteIndex} (Col: ${selectedSpriteCol}, Row: ${selectedSpriteRow})`;
        
        drawSpritesheet();
        drawMagnifiedSprite();
    }
});


// --- Hitbox Drawing on Magnified Canvas (Unchanged) ---

function snapToPixelGrid(coord) {
    return Math.round(coord / ZOOM_FACTOR) * ZOOM_FACTOR;
}

magnifiedCanvas.addEventListener('mousedown', (e) => {
    if (selectedSpriteIndex === -1) {
        alert("Please select a sprite from the full sheet first.");
        return;
    }

    isDrawing = true;
    const pos = getMousePos(magnifiedCanvas, e);
    
    startX = snapToPixelGrid(pos.x);
    startY = snapToPixelGrid(pos.y);
    
    hitbox = { x: startX, y: startY, w: 0, h: 0 };
});

magnifiedCanvas.addEventListener('mousemove', (e) => {
    if (!isDrawing) return;

    const pos = getMousePos(magnifiedCanvas, e);
    
    const currentX = snapToPixelGrid(pos.x);
    const currentY = snapToPixelGrid(pos.y);

    const width = currentX - startX;
    const height = currentY - startY;

    hitbox.x = Math.min(startX, currentX);
    hitbox.y = Math.min(startY, currentY);
    hitbox.w = Math.abs(width);
    hitbox.h = Math.abs(height);
    
    drawMagnifiedSprite(); 
});

magnifiedCanvas.addEventListener('mouseup', () => {
    if (!isDrawing) return;
    isDrawing = false;
    
    calculateHitboxRelativeCoords();
});

// --- Coordinate Calculation (Simplified Output) ---

function calculateHitboxRelativeCoords() {
    if (selectedSpriteIndex === -1 || (hitbox.w === 0 && hitbox.h === 0)) {
        hitboxCoordinatesDiv.innerHTML = "Hitbox: N/A";
        return;
    }
    
    // Convert Magnified Hitbox Coords (4x) back to Sprite Coords (1x)
    const spriteX_1x = hitbox.x / ZOOM_FACTOR;
    const spriteY_1x = hitbox.y / ZOOM_FACTOR;
    const spriteW_1x = hitbox.w / ZOOM_FACTOR;
    const spriteH_1x = hitbox.h / ZOOM_FACTOR;

    // Calculate Relative Coords using Origin
    const relativeLeft = spriteX_1x - originX;
    const relativeTop = spriteY_1x - originY;
    
    const relativeRight = relativeLeft + spriteW_1x;
    const relativeBottom = relativeTop + spriteH_1x;

    // Final integer pixel values
    const finalLeft = Math.round(relativeLeft);
    const finalTop = Math.round(relativeTop);
    const finalRight = Math.round(relativeRight);
    const finalBottom = Math.round(relativeBottom);
    const finalW = Math.round(spriteW_1x);
    const finalH = Math.round(spriteH_1x);

    // --- Display Output ---

    hitboxCoordinatesDiv.innerHTML = `
        Boundaries: L: ${finalLeft}, T: ${finalTop}, R: ${finalRight}, B: ${finalBottom}<br>
        Dimensions: W: ${finalW}, H: ${finalH}<br><br>

        (Rect): { Left: ${finalLeft}, Top: ${finalTop}, Right: ${finalRight}, Bottom: ${finalBottom} }<br>
        (XYWH): { x: ${finalLeft}, y: ${finalTop}, w: ${finalW}, h: ${finalH} }
    `;
}

// Attach initial event listeners to controls
gridColsInput.addEventListener('change', updateGrid);
gridRowsInput.addEventListener('change', updateGrid);
originXInput.addEventListener('change', updateOrigin);
originYInput.addEventListener('change', updateOrigin);

// Initial setup
updateOrigin();
