function placeUserPhotos(data){
    console.log("placePhotos");
    myPhotos = data.Thumbnails;
  
    if (myPhotos != null) {
        for (i = 0; i < myPhotos.length; i++){
            document.getElementById("userPhotos").innerHTML += "<div class='col-lg-2 col-sm-4 col-xs-6 well'><a title="
                                            + myPhotos[i].ImgName 
                                            + " href='/photo/" 
                                            + myPhotos[i].Id 
                                            + "'>"
                                            + "</span><img id=" 
                                            + myPhotos[i].Id 
                                            + " class='thumbnail centerBlock img-responsive' src="
                                            + myPhotos[i].Thumbnail 
                                            + "></a>"
                                            + "<span><button type='button' class='btn btn-sm btn-danger' id='deleteCat' onclick='deletePhoto("
                                            + myPhotos[i].Id
                                            + ")'>Delete cat</button></div>";

        }
    } else {
        document.getElementById("userPhotos").innerHTML = "No cats to be found yet, perhaps a dog scared them away. Or perhaps you haven't started a new project yet. Go to Cat Magic to get started!";
    }

            
}


function deletePhoto(pid) {
    console.log(pid);
    var uid = 1;        //CHANGE USERID
    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
        if (xhr.status == 200) {
          location.reload(); // reloads page
          alert("Success! Your cat photo is deleted");
        } else {
          alert("Error! Delete failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Delete failed." + xhr.status);
      };
        
      xhr.open('DELETE', 'http://130.240.170.62:1026/photo/' + pid + "/" + uid, true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(JSON.stringify(null));
}

