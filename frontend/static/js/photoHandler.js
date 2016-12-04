document.getElementById("catPhotoUpload").onchange = uploadPhoto;

function uploadPhoto(){
    displayHats();
    var file = document.getElementById("catPhotoUpload").files[0];
    var readFile = new FileReader();

    function success(e){
        loadCanvas(e.target.result);
        //document.getElementById('catPhoto').src = e.target.result;
        
    //document.getElementById('hatHolder').innerHTML += "<img class = 'draggable' src='/static/img/hats/strawhat.png' id='strawhat' width='20%' height='10%'> <img class = 'draggable' src='/static/img/hats/tophat.png' id='tophat' width='20%' height='10%'>" ;
    
    };

  readFile.readAsDataURL(file);
  // function call to success after file is read, no () needed
  readFile.onload = success;
}

function downloadPhoto(text, name, type) {
  var a = document.getElementById("a");
  var file = new Blob([text], {type: type});
  a.href = URL.createObjectURL(file);
  a.download = name;
}

function displayHats(){
    var hats = document.getElementsByClassName('draggable');
    console.log(hats);
    if (hats.length == 0){
        document.getElementById('hatHolder').insertAdjacentHTML('afterbegin', "<img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/strawhat.png' id='strawhat' width='20%' height='10%'> <img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='https://upload.wikimedia.org/wikipedia/commons/thumb/b/b2/Green_square.svg/1024px-Green_square.svg.png' id='santa' width='20%' height='10%'>");
    }
    imgs = document.getElementById('hatHolder').getElementsByTagName('img');
    for(var i =0; i<imgs.length; i++){
        imgs[i].onmousedown = mousedown;
        imgs[i].ondragstart = dragstart;
    }
}


function loadCanvas(src) {
        var canvas = document.createElement('canvas'),
            div = document.getElementById('catPhotoHolder');
        var img = new Image();
        var context = canvas.getContext('2d');

        

        img.onload = function() {
            canvas.id = "catCanvas";
            canvas.height=window.innerHeight;
            canvas.width=window.innerHeight*(16/9);
            //canvas.style.width  = div.width; 


            if (div.firstElementChild != null) {
                img.name = "if";
                div.removeChild(div.getElementsByTagName('canvas')[0]);
                context.backgroundImg = img;
                context.imageList = [];
                console.log(context);
            } 
                //context.drawImage(img,0,0);
            div.appendChild(canvas);

            canvasInit(img);
            
        };
        img.src = src;
    }


function savePhoto() {
    var c=document.getElementById("catCanvas");
    var w=window.open(c.toDataURL('image/png'));

    /*html2canvas(document.getElementById('catPhotoHolder'), {
        onrendered: function(canvas) {
            document.body.appendChild(canvas);
  }
});*/

}

function addHatInDiv(divId, hatId) {
    var canvas = document.getElementById('catCanvas');
    var context = canvas.getContext('2d');

    hat = document.getElementById(hatId);
    newHat = hat.cloneNode();
    hat.style.transform = "translate(0,0)";
    hat.setAttribute("data-x", "0");
    hat.setAttribute("data-y", "0");

    var image = new DragImage(newHat.src, 0, 0);
    context.imageList.push(image);
    console.log("imageList" + context.imageList);
    //newParent.appendChild(newHat);
}

function removeHatFromDiv(hatId) {
    var hat = document.getElementById(hatId);
    hat.parentNode.removeChild(hat);
}

