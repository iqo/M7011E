var canvas;
var ctx;
var canvasLeft;
var canvasTop;
//

function canvasInit2(img) {

    canvas = document.getElementById('catCanvas');
    ctx = canvas.getContext('2d');
    ctx.backgroundImg = img;
    canvasLeft = canvas.offsetLeft;
    canvasTop = canvas.offsetTop;
    canvas.ondrop = drop;
    canvas.ondragover = allowDrop;
    canvas.onmousemove = function(ev){
    
    }
    drawImageProp(ctx, ctx.backgroundImg, 0, 0, canvas.width, canvas.height, 0.5, 0.5);
    //ctx.drawImage(ctx.backgroundImg,0,0, canvas.width, canvas.height);
        //c.fillRect(0, 0, 500, 500);
        //image.update();
}


// this is the mouse position within the drag element
var startOffsetX, startOffsetY;

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
    ctx.imageList.push(image);
}

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
        context.clearRect(0, 0, canvas.width, canvas.height);
        drawImageProp(ctx, ctx.backgroundImg,0,0, canvas.width, canvas.height, 0.5, 0.5);
        //c.fillRect(0, 0, 500, 500);
        //image.update();

        
        for (var i = 0; i < ctx.imageList.length; i++) {
            ctx.imageList[i].update();
        }
        console.log(ctx.imageList);

    }, 50);

    
    canvas.addEventListener('mousemove', function(e) {
      //mouseX = e.offsetX;
      //mouseY = e.offsetY;
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
        this.x = x;
        this.y = y;
        var img = new Image();
        img.src = src;
        this.update = function() {

            if (mousePressed){
                //console.log(mouseX + ", " + mouseY);
                //console.log(this.x + ", " + this.y);
                var left = that.x;
                var right = that.x + img.width;
                var top = that.y;
                var bottom = that.y + img.height;
                if (!drag){
                  startX = mouseX - that.x;
                  startY = mouseY - that.y;
                }
                if (mouseX < right && mouseX > left && mouseY < bottom && mouseY > top){
                   drag = true;
                }
            }else{
               drag = false;
            }
            if (drag){
                that.x = mouseX - startX;
                that.y = mouseY - startY;
            }
            ctx.drawImage(img, that.x, that.y, 50,50);

        }
    }
