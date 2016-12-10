function addComment(photoId) {
    var comment={};
    comment.photoId = photoId;
    comment.comment = document.getElementById("userComment").value;
    
    /********* CHANGE TO ACTUAL USERID **********/
    photo.uid = 1;

    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
        if (xhr.status == 200) {
          location.reload(); // reloads page
          alert("Successfully commented!");
        } else {
          alert("Error! Comment failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Comment failed." + xhr.status);
      };
        
      xhr.open('POST', 'http://130.240.170.62:1026/comment', true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(JSON.stringify(photo));
}