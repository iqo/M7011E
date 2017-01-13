var hasFavorite = false;
var ip = getIP();

function changeFavoriteColor(id){
    if (document.getElementById("favorite").style.color == "red") {
		document.getElementById("favorite").style.color = "black";
		hasFavorite = false;
    } else {
		document.getElementById("favorite").style.color = "red";
		hasFavorite = true;
	}
}

function changeFavorite(id){
	if (isSignedIn()) {
        updateFavorite(id);

    } else {
        alert("You need to be signed in to add as favorite!");
    }
}

function updateFavorite(photoId) {
    var f={};
    f.photoId = parseInt(photoId);

    f.uid = parseInt(returnUserId());
    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
    if (xhr.status == 200) {
        changeFavoriteColor(photoId);

    } else {
        alert("Error! Update favorite failed");
        }
    };
    xhr.onerror = function() {
        alert("Error! Update favorite failed." + xhr.status);
    };

    if (!hasFavorite) {
        xhr.open('POST', ip + '/favorite', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.send(JSON.stringify(f));

    } 
    if (hasFavorite) {
        xhr.open('DELETE', ip + '/favorite/' + photoId + '/' + f.uid, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.send(JSON.stringify(null));
    }
    
}

function getFavorite(pid) {
    var xhr = new XMLHttpRequest();
    var uid = parseInt(returnUserId());
    xhr.onload = function(event) {

        if (xhr.status == 200) {
          	var f = event.target.response;
            if (f.length == 0) {
                return;
            }
	        f = JSON.parse(f);
	        hasFavorite = true;
	        changeFavoriteColor(pid);
        } else {}
      };
      xhr.onerror = function() {
            alert("Error! Get rate failed. Cannot connect to server.");
      };
      xhr.open('GET', ip + '/favorite/' + pid + "/" + uid, true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);
}