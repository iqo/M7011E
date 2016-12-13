function placeUserPhotos(data){
    console.log("placePhotos")
    console.log(data);
    myPhotos = data.Toplist;
    console.log(myPhotos);
    
    if (myPhotos != null) {
        for (i = 0; i < myPhotos.length; i++){
                document.getElementById("myPhotos").innerHTML += "<div class='col-lg-4 col-sm-4 col-xs-4'><a title="
                                                + myPhotos[i].ImgName 
                                                + " href='/photo/" 
                                                + myPhotos[i].Id 
                                                + "'>"
                                                + "<span class='rate'>" 
                                                + [i+1] 
                                                +"</span><img id=" 
                                                + myPhotos[i].Id 
                                                + " class='toplist centerBlock img-responsive' src="
                                                + myPhotos[i].Thumbnail 
                                                + "></a></div>";

        }
    }
}