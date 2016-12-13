//document.getElementById("catFile").onchange = uploadPhoto;
function uploadPhoto(){
    
    var file = document.getElementById("catFile").files[0];
    var readFile = new FileReader();
    //console.log(file);

    function success(e){
        var img = e.target.result;
        loadCanvas(img);
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
    //console.log(hats);
    /*if (hats.length == 0){
        document.getElementById('hatHolder').insertAdjacentHTML('afterbegin', "<img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/strawhat.png' id='strawhat' width='10%' height='5%'> <img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/santa.png' id='santa' width='10%' height='5%'><img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/tophat.png' id='tophat' width='10%' height='5%'><img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/yellow_hat.png' id='yellow_hat' width='10%' height='5%'><img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/pirate.png' id='pirate' width='10%' height='5%'><img draggable='true' class='draggable' onmousedown=mousedown ondragstart=dragstart src='/static/img/hats/propeller.png' id='propeller' width='10%' height='5%'>");
    }*/
    document.getElementById('hatHolder').style.display = "block";
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
            context.loop = true;
            canvas.id = "catCanvas";
            canvas.width=window.innerHeight*(2/3)*(16/9);
            canvas.height=window.innerHeight *(2/3);          
            //canvas.style.width  = div.width; 


            if (div.firstElementChild != null) {
                context.loop = false;
                img.name = "if";
                div.removeChild(div.getElementsByTagName('canvas')[0]);
                context.backgroundImg = img;
                context.imageList = [];
            } 
            div.appendChild(canvas);

            canvasInit(img);
            
        };
        img.src = src;
    }

function canvasToPhoto(){
    var canvas = document.getElementById("catCanvas");
    var img    = canvas.toDataURL("image/png");
    var imgName = document.getElementById("imgName").value;

    var thumbnail = getThumbnail(img, canvas.width*(1/5), canvas.height*(1/5));

}

function photoToCanvas(src){
    canvas = document.getElementById('photoCanvas');
    ctx = canvas.getContext('2d');
    canvas.width=window.innerHeight*(2/3)*(16/9);
    canvas.height=window.innerHeight *(2/3); 
    var img = new Image();
    img.onload = function() {
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
    }
    img.src = src;
    
}


function savePhoto(img, thumbnail) {
    var imgName = document.getElementById("imgName").value;
    if (imgName == "") {imgName = "untitled"};
    var photo={};
    photo.imgName = imgName;
    photo.imgDesc = document.getElementById("imgDesc").value;
    photo.image = img;//img;
    photo.thumbnail = thumbnail;
    
    /********* CHANGE TO ACTUAL USERID **********/
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

          
          photo = JSON.parse(photo);
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

function getLatestPhotos(page) {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
        console.log(xhr.status);
        if (xhr.status == 200) {
          var thumbnail = event.target.response;
          
          thumbnail = JSON.parse(thumbnail);
          placeLatestPhotos(thumbnail.Thumbnails);

        } else {
          alert("Error! Get files failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Get file failed. Cannot connect to server.");
      };
        
      xhr.open('GET', 'http://130.240.170.62:1026/latest/' + page, false);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);
}

function placeLatestPhotos(thumbnails){
    if (thumbnails != null) {
        thumbnails.forEach(function(thumbnail){
        document.getElementById('latestPhotos').innerHTML += "<div class='col-lg-2 col-sm-4 col-xs-6'><a title="
                                                + thumbnail.ImgName 
                                                + " href='/photo/" 
                                                + thumbnail.Id 
                                                + "'><img id=" 
                                                + thumbnail.Id 
                                                + " class='thumbnail img-responsive' src="
                                                + thumbnail.Thumbnail 
                                                + "></a></div>";
    }); 
  }
}



function addSaveButton(){
    var div = document.getElementById('buttonHolder');
    if (div.childNodes.length == 3) {
        div.innerHTML += "<div id='savePhoto'><button type='button' class='btn btn-sm btn-primary' id='saveCat' data-toggle='modal' data-target='#saveImageModal'>Save cat</button></div>";
    };
}

function getThumbnail(src, width, height) {
    var canvas = document.createElement("canvas");
    var ctx = canvas.getContext("2d"); 
    canvas.width = width;
    canvas.height = height;

    var original = new Image();
    original.onload = function(){
        ctx.drawImage(original, 0, 0, canvas.width, canvas.height);
        var thumbnail = canvas.toDataURL("image/png");
        savePhoto(src, thumbnail);

    }
    original.src = src;
    
}


function getToplist() {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
        console.log(xhr.status);
        if (xhr.status == 200) {
          var top = event.target.response;
          
          top = JSON.parse(top);
          placeToplist(top.Toplist);

        } else {
          alert("Error! Get toplist failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Get file failed. Cannot connect to server.");
      };
        
      xhr.open('GET', 'http://130.240.170.62:1026/top', false);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);
}

function placeToplist(toplist){
    if (toplist != null) {
        for (i = 0; i < toplist.length; i++){
            if (i == 0) {
                document.getElementById('toplist-alltime').innerHTML += "<div class='row'><div class='col-lg-4 col-sm-4 col-xs-4'><a title="
                                                + toplist[0].ImgName 
                                                + " href='/photo/" 
                                                + toplist[0].Id 
                                                + "'>"
                                                + [i]
                                                +"<img id=" 
                                                + toplist[0].Id 
                                                + " class='toplist img-responsive' src="
                                                + toplist[0].Thumbnail 
                                                + "></a></div></div>";
            } else{
                document.getElementById('toplist-alltime').innerHTML += "<div class='col-lg-4 col-sm-4 col-xs-4'><a title="
                                                + toplist[i].ImgName 
                                                + " href='/photo/" 
                                                + toplist[i].Id 
                                                + "'>"
                                                + [i]
                                                + "<img id=" 
                                                + toplist[i].Id 
                                                + " class='toplist img-responsive' src="
                                                + toplist[i].Thumbnail 
                                                + "></a></div>";

    }; 
  }
}
}

