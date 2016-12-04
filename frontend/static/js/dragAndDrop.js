
  function dragMoveListener (event) {
    var target = event.target,
        // keep the dragged position in the data-x/data-y attributes
        x = (parseFloat(target.getAttribute('data-x')) || 0) + event.dx,
        y = (parseFloat(target.getAttribute('data-y')) || 0) + event.dy;

    // translate the element
    target.style.webkitTransform =
    target.style.transform =
      'translate(' + x + 'px, ' + y + 'px)';

    // update the posiion attributes
    target.setAttribute('data-x', x);
    target.setAttribute('data-y', y);
  }

 /*   function drawOnCanvas (event) {
        var canvas = document.getElementById("catCanvas");
        if (canvas) {
            var ctx = canvas.getContext("2d");
            var target = event.target,
            // keep the dragged position in the data-x/data-y attributes
            x = (parseFloat(target.getAttribute('data-x')) || 0) + event.dx,
            y = (parseFloat(target.getAttribute('data-y')) || 0) + event.dy;
            //ctx.drawImage(target, 33, 71, 104, 124, 21, 20, 87, 104);
            console.log("x: " + x + ", y: " + y);
            var rect = canvas.getBoundingClientRect();
            console.log(rect.top, rect.right, rect.bottom, rect.left);
        } else {
            console.log("droped hat, no canvas")
        }
        
    }
*/

  window.dragMoveListener = dragMoveListener;

interact('.draggable')
  .draggable({
    onmove: window.dragMoveListener
    //onend: window.drawOnCanvas
  })
  .resizable({
    preserveAspectRatio: true,
    edges: { left: true, right: true, bottom: true, top: true }
  })
  .on('resizemove', function (event) {
    var target = event.target,
        x = (parseFloat(target.getAttribute('data-x')) || 0),
        y = (parseFloat(target.getAttribute('data-y')) || 0);

    // update the element's style
    target.style.width  = event.rect.width + 'px';
    target.style.height = event.rect.height + 'px';

    // translate when resizing from top or left edges
    x += event.deltaRect.left;
    y += event.deltaRect.top;

    target.style.webkitTransform = target.style.transform =
        'translate(' + x + 'px,' + y + 'px)';

    target.setAttribute('data-x', x);
    target.setAttribute('data-y', y);
    target.textContent = Math.round(event.rect.width) + '×' + Math.round(event.rect.height);
});

var insideDropzone = "unset";

interact('.dropzone').dropzone({
  // only accept elements matching this CSS selector
  //accept: '#yes-drop',
  // Require a 75% element overlap for a drop to be possible
  //overlap: 0.75,

  
// listen for drop related events:
    ondropactivate: function (event) {
        var hat = event.relatedTarget;
        var id = hat.id;
        console.log(id);
        // add active dropzone feedback
        //event.target.classList.add('drop-active');
        console.log("ondropactive");
        console.log(insideDropzone);
      },
    ondragenter: function (event) {
        insideDropzone = true;
      },
    ondragleave: function (event) {
        insideDropzone = false;
        console.log("leaves dropzone");
        var hat = event.relatedTarget

      },
    ondrop: function (event) {
        console.log(event);
        var hat = event.relatedTarget,
            dropzone = event.target.id;
        if (document.getElementById(hat.id).parentNode.id != dropzone){
            addHatInDiv(dropzone, hat.id);
        }
        
        //removeHatFromDiv(hat.id);
        
        console.log("dropped");
        
      },
    ondropdeactivate: function (event) {
        console.log("not dragging anymore");
        var hat = event.relatedTarget;

      }
});




