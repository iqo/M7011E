var canvas;
var ctx;
var canvasLeft;
var canvasTop;
var startOffsetX, startOffsetY;

var pi2 = Math.PI * 2;
var resizerRadius = 12;
var rr = resizerRadius * resizerRadius;
var draggingResizer = {
    x: 0,
    y: 0
};

var startX;
var startY;
var dragId = "";

var trash = new Image();
trash.src = '/static/img/icon/trash.png';
//

function allowDrop(ev) {
    ev.preventDefault();
}

function mousedown(ev) {
    startOffsetX = ev.offsetX;
    startOffsetY = ev.offsetY;
}

function dragstart(ev) {
    ev.dataTransfer.setData("Text", ev.target.id);
}

function drop(ev) {
    ev.preventDefault();

    //var dropX = ev.clientX - canvasLeft - startOffsetX;
    //var dropY = ev.clientY - canvasTop - startOffsetY;
    var dropX = ev.offsetX - 75;
    var dropY = ev.offsetY-75;
    var id = ev.dataTransfer.getData("Text");
    
    var dropElement = document.getElementById(id);

    // draw the drag image at the drop coordinates
    var image = new DragImage(dropElement.src, dropX, dropY, id);
    ctx.drawImage(dropElement, dropX, dropY, 150, 150);
    
    ctx.imageList.push(image);
}



var mouseX = 0, mouseY = 0;
var mousePressed = false;

function canvasInit(img) {
    canvas = document.getElementById('catCanvas');
    ctx = canvas.getContext('2d');
    ctx.backgroundImg = img;
    canvasLeft = canvas.offsetLeft;
    canvasTop = canvas.offsetTop;
    canvas.ondrop = drop;
    canvas.ondragover = allowDrop;
    ctx.imageList = [];
    ctx.backgroundImg = img;
    var loop = setInterval(function() {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
        //drawImageProp(ctx, ctx.backgroundImg,0,0, canvas.width, canvas.height, 0.5, 0.5);

        for (var i = 0; i < ctx.imageList.length; i++) {
            ctx.imageList[i].update();
        }
        if (ctx.loop == false) {clearInterval(loop)};

    }, 20);

    
    canvas.addEventListener('mousemove', function(e) {
      mousePos = getMousePos(canvas,e);
      mouseX = mousePos.x;
      mouseY = mousePos.y;

    });
     
    $(document).mousedown(function(){
        mousePressed = true;
    }).mouseup(function(){
        mousePressed = false;
    });
}

function getMousePos(canvas, evt) {
        var rect = canvas.getBoundingClientRect();
        return {
          x: evt.clientX - rect.left,
          y: evt.clientY - rect.top
        };
      }

 function DragImage(src, x, y, id) {
        var imageX = 100;
        var imageY = 100;
        var imageWidth = 100;
        var imageHeight = 100;
        var that = this;
        var startX = 0, startY = 0;
        var drag = false;
        var remove = false;
        var dragResize = -1;
        this.x = x;
        this.y = y;
        var img = new Image();
        img.src = src;
        img.name = id;
        this.update = function() {
                var left = that.x;
                var right = that.x + imageWidth;
                var top = that.y;
                var bottom = that.y + imageHeight;
            if (mousePressed) {

                if (!drag){
                  startX = mouseX - that.x;
                  startY = mouseY - that.y;

                }
                if (hitImage(left, right, top, bottom) && (dragId == "" || dragId == id)){
                   drag = true;
                   dragId = id;
                   putCorners(left,right,top, bottom);
                   dragResize = checkIfResize(startX, startY, right, bottom);
                   drawTrash();
                }
            } else {
                dragId = "";
                drag = false;
                if (hitImage(left, right, top, bottom)){
                   putCorners(left,right,top, bottom);
                }
                checkIfremoveHat(this, left, bottom);
            }
            // Sets new position, drags object
            if (drag && dragResize == -1){
                that.x = mouseX - startX;
                that.y = mouseY - startY;

            } 
            if (drag && (dragResize == 0 || dragResize == 1 || dragResize == 2 || dragResize == 3)) {
                putCorners(left,right,top, bottom);
                switch (dragResize) {
                    case 0:
                        //top-left
                        that.x = mouseX - startX;
                        imageWidth = right - that.x;
                        that.y = mouseY - startY;
                        imageHeight = bottom - that.y;
                        
                        break;
                    case 1:
                        //top-right
                        //that.x = mouseX - startX;
                        imageWidth = mouseX - that.x;
                        that.y = mouseY - startY;
                        imageHeight = bottom-that.y;
                        break;
                    case 2:
                        //bottom-right
                        imageWidth = mouseX - that.x;
                        imageHeight = mouseY - that.y;
                        break;
                    case 3:
                        //bottom-left
                        that.x = mouseX - startX;
                        imageWidth = right - that.x;
                        //that.y = mouseY - startY;
                        imageHeight = mouseY - that.y;
                        
                        break;
                }


                if(imageWidth<25){imageWidth=25;}
                if(imageHeight<25){imageHeight=25;}
                imageX = imageWidth;
                imageY = imageHeight;

                    // set the image right and bottom
            };
            ctx.drawImage(img, that.x, that.y, imageWidth,imageHeight);

        }
    }

function hitImage(left, right, top, bottom){
    return (mouseX < right && mouseX > left && mouseY < bottom && mouseY > top)
}
//image x = left, imagey = top, imagerignt= right, imagebottom=bottom 
// startx, starty = x,y = that.x, that.y
function checkIfResize (x, y, right, bottom) {
    var dx, dy;
    // top-left
    dx = x;
    dy = y;
    if (dx * dx + dy * dy <= rr) {
        return (0);
    }

    // top-right
    dx = mouseX - right;
    dy = y;
    if (dx * dx + dy * dy <= rr) {
        return (1);
    }
    // bottom-right
    dx = mouseX - right;
    dy = mouseY - bottom;
    if (dx * dx + dy * dy <= rr) {
        return (2);
    }
    // bottom-left
    dx = x;
    dy = mouseY - bottom;
    if (dx * dx + dy * dy <= rr) {
        return (3);
    }
    return (-1);

}



function checkIfremoveHat(img, left, bottom){
    if (left < 30 && bottom > (canvas.height-40)) {
        var index = ctx.imageList.indexOf(img);
        if (index > -1) {
            ctx.imageList.splice(index, 1);
        };
    };

}

function putCorners (left, right, top, bottom) {
    drawCorners(left+5, top+5);
    drawCorners(right-5, top+5);
    drawCorners(left+5, bottom-5);
    drawCorners(right-5, bottom-5);
}

function drawCorners(x, y) {
    ctx.beginPath();
    ctx.arc(x, y, resizerRadius, 0, pi2, false);
    ctx.closePath();
    ctx.fill();
}

function drawTrash() {
    ctx.drawImage(trash, 0, canvas.height-60 , 50, 50);
}

/*
// Draws background image on canvas
function drawImageProp(ctx, img, x, y, w, h, offsetX, offsetY) {

    if (arguments.length === 2) {
        x = y = 0;
        w = ctx.canvas.width;
        h = ctx.canvas.height;
    }

    /// default offset is center
    offsetX = typeof offsetX === 'number' ? offsetX : 0.5;
    offsetY = typeof offsetY === 'number' ? offsetY : 0.5;

    /// keep bounds [0.0, 1.0]
    if (offsetX < 0) offsetX = 0;
    if (offsetY < 0) offsetY = 0;
    if (offsetX > 1) offsetX = 1;
    if (offsetY > 1) offsetY = 1;

    var iw = img.width,
        ih = img.height,
        r = Math.min(w / iw, h / ih),
        nw = iw * r,   /// new prop. width
        nh = ih * r,   /// new prop. height
        cx, cy, cw, ch, ar = 1;

    /// decide which gap to fill    
    if (nw < w) ar = w / nw;
    if (nh < h) ar = h / nh;
    nw *= ar;
    nh *= ar;

    /// calc source rectangle
    cw = iw / (nw / w);
    ch = ih / (nh / h);

    cx = (iw - cw) * offsetX;
    cy = (ih - ch) * offsetY;

    /// make sure source rectangle is valid
    if (cx < 0) cx = 0;
    if (cy < 0) cy = 0;
    if (cw > iw) cw = iw;
    if (ch > ih) ch = ih;

    /// fill image in dest. rectangle
    ctx.drawImage(img, cx, cy, cw, ch,  x, y, w, h);
}

*/
