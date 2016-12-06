var canvasOffset;
var offsetX;
var offsetY;

var startX;
var startY;
var isDown = false;


var pi2 = Math.PI * 2;
var resizerRadius = 8;
var rr = resizerRadius * resizerRadius;
var draggingResizer = {
    x: 0,
    y: 0
};
var imageX = 50;
var imageY = 50;
var imageWidth, imageHeight, imageRight, imageBottom;
var draggingImage = false;
var startX;
var startY;

var mouseX = 0, mouseY = 0;
var mousePressed = false;
var c;


function canvasInit(img) {
    canvas = document.getElementById('catCanvas');
    ctx = canvas.getContext('2d');
    canvasOffset = $("#canvas").offset();
    offsetX = canvasOffset.left;
    offsetY = canvasOffset.top;
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


    ctx.drawImage(dropElement, dropX, dropY, 150, 150);
    var image = new DragImage(dropElement.src, dropX, dropY);

    var img = new Image();
    img.onload = function () {
        imageWidth = img.width;
        imageHeight = img.height;
        imageRight = imageX + imageWidth;
        imageBottom = imageY + imageHeight
        draw(true, false);
}
    img.src = dropElement.src;
    ctx.imageList.push(img);
}


function draw(withAnchors, withBorders) {

    // clear the canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // draw the image
    ctx.drawImage(img, 0, 0, img.width, img.height, imageX, imageY, imageWidth, imageHeight);

    // optionally draw the draggable anchors
    if (withAnchors) {
        drawDragAnchor(imageX, imageY);
        drawDragAnchor(imageRight, imageY);
        drawDragAnchor(imageRight, imageBottom);
        drawDragAnchor(imageX, imageBottom);
    }

    // optionally draw the connecting anchor lines
    if (withBorders) {
        ctx.beginPath();
        ctx.moveTo(imageX, imageY);
        ctx.lineTo(imageRight, imageY);
        ctx.lineTo(imageRight, imageBottom);
        ctx.lineTo(imageX, imageBottom);
        ctx.closePath();
        ctx.stroke();
    }

}

function drawDragAnchor(x, y) {
    ctx.beginPath();
    ctx.arc(x, y, resizerRadius, 0, pi2, false);
    ctx.closePath();
    ctx.fill();
}

function anchorHitTest(x, y) {

    var dx, dy;

    // top-left
    dx = x - imageX;
    dy = y - imageY;
    if (dx * dx + dy * dy <= rr) {
        return (0);
    }
    // top-right
    dx = x - imageRight;
    dy = y - imageY;
    if (dx * dx + dy * dy <= rr) {
        return (1);
    }
    // bottom-right
    dx = x - imageRight;
    dy = y - imageBottom;
    if (dx * dx + dy * dy <= rr) {
        return (2);
    }
    // bottom-left
    dx = x - imageX;
    dy = y - imageBottom;
    if (dx * dx + dy * dy <= rr) {
        return (3);
    }
    return (-1);

}


function hitImage(x, y) {
    return (x > imageX && x < imageX + imageWidth && y > imageY && y < imageY + imageHeight);
}


function handleMouseDown(e) {
    startX = parseInt(e.clientX - offsetX);
    startY = parseInt(e.clientY - offsetY);
    draggingResizer = anchorHitTest(startX, startY);
    draggingImage = draggingResizer < 0 && hitImage(startX, startY);
}

function handleMouseUp(e) {
    draggingResizer = -1;
    draggingImage = false;
    draw(true, false);
}

function handleMouseOut(e) {
    handleMouseUp(e);
}

function handleMouseMove(e) {

    if (draggingResizer > -1) {

        mouseX = parseInt(e.clientX - offsetX);
        mouseY = parseInt(e.clientY - offsetY);

        // resize the image
        switch (draggingResizer) {
            case 0:
                //top-left
                imageX = mouseX;
                imageWidth = imageRight - mouseX;
                imageY = mouseY;
                imageHeight = imageBottom - mouseY;
                break;
            case 1:
                //top-right
                imageY = mouseY;
                imageWidth = mouseX - imageX;
                imageHeight = imageBottom - mouseY;
                break;
            case 2:
                //bottom-right
                imageWidth = mouseX - imageX;
                imageHeight = mouseY - imageY;
                break;
            case 3:
                //bottom-left
                imageX = mouseX;
                imageWidth = imageRight - mouseX;
                imageHeight = mouseY - imageY;
                break;
        }

        if(imageWidth<25){imageWidth=25;}
        if(imageHeight<25){imageHeight=25;}

        // set the image right and bottom
        imageRight = imageX + imageWidth;
        imageBottom = imageY + imageHeight;

        // redraw the image with resizing anchors
        draw(true, true);

    } else if (draggingImage) {

        imageClick = false;

        mouseX = parseInt(e.clientX - offsetX);
        mouseY = parseInt(e.clientY - offsetY);

        // move the image by the amount of the latest drag
        var dx = mouseX - startX;
        var dy = mouseY - startY;
        imageX += dx;
        imageY += dy;
        imageRight += dx;
        imageBottom += dy;
        // reset the startXY for next time
        startX = mouseX;
        startY = mouseY;

        // redraw the image with border
        draw(false, true);

    }
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


$("#canvas").mousedown(function (e) {
    handleMouseDown(e);
});
$("#canvas").mousemove(function (e) {
    handleMouseMove(e);
});
$("#canvas").mouseup(function (e) {
    handleMouseUp(e);
});
$("#canvas").mouseout(function (e) {
    handleMouseOut(e);
});