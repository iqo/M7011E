var canvas;
var ctx;
var canvasLeft;
var canvasTop;
var startOffsetX, startOffsetY;


var pi2 = Math.PI * 2;
var resizerRadius = 10;
var rr = resizerRadius * resizerRadius;
var draggingResizer = {
    x: 0,
    y: 0
};

var imageX = 100;
var imageY = 100;
var imageWidth = 100;
var imageHeight = 100;
var draggingImage = false;
var startX;
var startY;
//

function allowDrop(ev) {
    ev.preventDefault();
}

function mousedown(ev) {
    startOffsetX = ev.offsetX;
    startOffsetY = ev.offsetY;
}

function dragstart(ev) {
    console.log("dragstart");
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

    console.log(dropElement);
    // draw the drag image at the drop coordinates

    var image = new DragImage(dropElement.src, dropX, dropY);
    ctx.drawImage(dropElement, dropX, dropY, 150, 150);
    
    ctx.imageList.push(image);
}



var mouseX = 0, mouseY = 0;
var mousePressed = false;
var c;

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
        drawImageProp(ctx, ctx.backgroundImg,0,0, canvas.width, canvas.height, 0.5, 0.5);

        for (var i = 0; i < ctx.imageList.length; i++) {
            ctx.imageList[i].update();
        }

    }, 50);

    
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

 function DragImage(src, x, y) {
        var that = this;
        var startX = 0, startY = 0;
        var drag = false;
        var dragResize = -1;
        this.x = x;
        this.y = y;
        var img = new Image();
        img.src = src;
        this.update = function() {
            if (mousePressed) {
                var left = that.x;
                var right = that.x + imageWidth;
                console.log(right);
                var top = that.y;
                var bottom = that.y + imageHeight;
                if (!drag){
                  startX = mouseX - that.x;
                  startY = mouseY - that.y;
                }
                if (mouseX < right && mouseX > left && mouseY < bottom && mouseY > top){
                   drag = true;
                   putCorners(left,right,top, bottom);
                   dragResize = checkIfResize(startX, startY, left, right, top, bottom);
                }
            } else {
               drag = false;
            }
            if (drag && dragResize == -1){
                that.x = mouseX - startX;
                that.y = mouseY - startY;
            }
            if (drag && (dragResize == 0 || dragResize == 1 || dragResize == 2 || dragResize == 3)) {
                switch (dragResize) {
                    case 0:
                        //top-left
                        imageWidth = startX - mouseX;
                        imageHeight = startY - mouseY;
                        
                        break;
                    case 1:
                        //top-right
                        imageWidth = mouseX - that.x;
                        imageHeight = bottom - mouseY;
                        break;
                    case 2:
                        //bottom-right
                        imageWidth = mouseX - that.x;
                        imageHeight = mouseY - that.y;
                        break;
                    case 3:
                        //bottom-left
                        imageWidth = right - mouseX;
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

function checkIfResize (x, y, left, right, top, bottom) {

    var dx, dy;
    // top-left
    dx = x;
    dy = y;
    console.log(x, y);
    console.log(dx, dy);
    if (dx * dx + dy * dy <= rr) {
        console.log("top-left");
        return (0);
    }

    // top-right
    dx = x - imageX;
    dy = y;
    if (dx * dx + dy * dy <= rr) {
        console.log("top-right");
        return (1);
    }
    // bottom-right
    dx = x - imageX;
    dy = y - imageY;
    if (dx * dx + dy * dy <= rr) {
        console.log("bottom-right");
        return (2);
    }
    // bottom-left
    dx = x;
    dy = y - imageY;
    if (dx * dx + dy * dy <= rr) {
        console.log("bottom-left");
        return (3);
    }
    return (-1);

}

function putCorners (left, right, top, bottom) {
    drawCorners(left, top);
    drawCorners(right, top);
    drawCorners(left, bottom);
    drawCorners(right, bottom);
}

function drawCorners(x, y) {
    ctx.beginPath();
    ctx.arc(x, y, resizerRadius, 0, pi2, false);
    ctx.closePath();
    ctx.fill();
}


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
