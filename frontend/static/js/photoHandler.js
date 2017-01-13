var ip = getIP();

function uploadPhoto(){
    
    var file = document.getElementById("catFile").files[0];
    var readFile = new FileReader();
    if (checkFileType(file.type)){

        function success(e){
            var img = e.target.result;
            loadCanvas(img);
            displayHats();
            addSaveButton();
        };

      readFile.readAsDataURL(file);
      // function call to success after file is read, no () needed
      readFile.onload = success;
    } else {
        alert("Not a valid file format");
    }
}

function checkFileType(fileType){
    type = fileType.split("/")[0];
    if (type == 'image'){
        return true;
    } else {
        return false;
    }
}



function clickUploadButton(){
    document.getElementById('catFile').click();
    document.getElementById("catFile").onchange = uploadPhoto;
    return false;
}


function displayHats(){
    var hats = document.getElementsByClassName('draggable');
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
        console.log(div.firstElementChild);

        img.onload = function() {
            
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
                context.loop = true;

            } 
            context.loop = true;
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

function photoToCanvas(src, canvas){
    canvas = document.getElementById(canvas);
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
    if (isSignedIn()) {
        var imgName = document.getElementById("imgName").value;
        if (imgName == "") {imgName = "untitled"};
        var photo={};
        photo.imgName = imgName;
        photo.imgDesc = document.getElementById("imgDesc").value;
        photo.image = img;//img;
        photo.thumbnail = thumbnail;

        photo.uid = parseInt(returnUserId());

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
            
          xhr.open('POST', ip + '/photo', true);
          xhr.setRequestHeader("Content-Type", "application/json");
          xhr.send(JSON.stringify(photo));
      }
}

function getPhoto(pid, canvas) {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
        if (xhr.status == 200) {
          var photo = event.target.response;

          
          photo = JSON.parse(photo);
          //window.open(photo.Image);
          photoToCanvas(photo.Image, canvas);
          //document.getElementById("test").src = photo.Image;
        } else {
          alert("Error! Get files failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Get file failed. Cannot connect to server.");
      };
        
      xhr.open('GET', ip + '/photo/get/' + pid, true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);
}

function getLatestPhotos(page) {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
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
        
      xhr.open('GET', ip + '/photo/latest/' + page, true);
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
    var div = document.getElementById('saveCheck');
    if (isSignedIn()) {
        div.innerHTML = "<div id='savePhoto'><button type='button' class='btn btn-sm btn-primary' id='saveCat' data-toggle='modal' data-target='#saveImageModal'>Save cat</button></div>";
    } else {
        div.innerHTML = "You need to sign in to save cat";
    }
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


function placeToplist(data, div){
    toplist = data.Toplist;
    if (toplist != null) {
        for (i = 0; i < toplist.length; i++){
                document.getElementById(div).innerHTML += "<div class='col-lg-4 col-sm-4 col-xs-4'><a title="
                                                + toplist[i].ImgName 
                                                + " href='/photo/" 
                                                + toplist[i].Id 
                                                + "'>"
                                                + "<span class='rate'>" 
                                                + [i+1] 
                                                +"</span><img id=" 
                                                + toplist[i].Id 
                                                + " class='toplist centerBlock img-responsive' src="
                                                + toplist[i].Thumbnail 
                                                + "></a></div>";

        }
    }
}

