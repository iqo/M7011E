//document.getElementById("catFile").onchange = uploadPhoto;
function uploadPhoto(){
    
    var file = document.getElementById("catFile").files[0];
    var readFile = new FileReader();
    console.log(file);

    function success(e){
        loadCanvas(e.target.result);
        displayHats();
        addSaveButton();
    };

  readFile.readAsDataURL(file);
  // function call to success after file is read, no () needed
  readFile.onload = success;
}



function clickUploadButton(){
    document.getElementById('catFile').click();
    document.getElementById("catFile").onchange = uploadPhoto;
    return false;
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
        document.getElementById('hatHolder').insertAdjacentHTML('afterbegin', "<img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/strawhat.png' id='strawhat' width='10%' height='5%'> <img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/santa.png' id='santa' width='10%' height='5%'>");
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
            canvas.width=window.innerHeight*(16/9);
            canvas.height=window.innerHeight;
            
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
    var canvas = document.getElementById("catCanvas");
    var img    = canvas.toDataURL("image/png");
    console.log(img);
    //var w=window.open(c.toDataURL('image/png'));
    var photo={};
    photo.imgName = "testCat";
    photo.imgDesc = "mjaow";
    photo.image = img;//img;
    photo.uid = 1;
    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
        if (xhr.status == 200) {
          location.reload(); // reloads page
          alert("Success! Upload completed");
        } else {
          alert("Error! Upload failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Upload failed." + xhr.status);
      };
        
      xhr.open('POST', 'http://130.240.170.62:1026/photo', true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(JSON.stringify(photo));

}

function getPhoto(id) {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
        console.log(xhr.status);
        if (xhr.status == 200) {
          var photo = event.target.response;
          console.log(photo[0].Image);
          //photo = JSON.parse(photo);
          //window.open(photo.Image);
          document.getElementById("test").src = photo.Image;
        } else {
          alert("Error! Get files failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Get file failed. Cannot connect to server.");
      };
        
      xhr.open('GET', 'http://130.240.170.62:1026/photo/' + id, false);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);

}

function createCORSRequest(method, url) {
  var xhr = new XMLHttpRequest();
  if ("withCredentials" in xhr) {

    // Check if the XMLHttpRequest object has a "withCredentials" property.
    // "withCredentials" only exists on XMLHTTPRequest2 objects.
    xhr.open(method, url, true);

  } else if (typeof XDomainRequest != "undefined") {

    // Otherwise, check if XDomainRequest.
    // XDomainRequest only exists in IE, and is IE's way of making CORS requests.
    xhr = new XDomainRequest();
    xhr.open(method, url);

  } else {

    // Otherwise, CORS is not supported by the browser.
    xhr = null;

  }
  return xhr;
}



function addSaveButton(){
    var div = document.getElementById('buttonHolder');
    console.log(div.childNodes.length);
    if (div.childNodes.length == 3) {
        div.innerHTML += "<div id='savePhoto'><button type='button' class='btn btn-sm btn-primary' id='saveCat' onclick='savePhoto()'>Save cat</button></div>";
    };
}

function addHatInDiv(divId, hatId) {
    var canvas = document.getElementById('catCanvas');
    var context = canvas.getContext('2d');

    hat = document.getElementById(hatId);
    newHat = hat.cloneNode();
    hat.style.transform = "translate(0,0)";
    hat.setAttribute("data-x", "0");
    hat.setAttribute("data-y", "0");
    console.log(hat.src);
    var image = new DragImage(newHat.src, 0, 0);
    context.imageList.push(image);
    console.log("imageList" + context.imageList);
    //newParent.appendChild(newHat);
}


