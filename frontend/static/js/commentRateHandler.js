function addComment(photoId) {
    var comment={};
    comment.photoId = parseInt(photoId);
    comment.comment = document.getElementById("userComment").value;
    
    /********* CHANGE TO ACTUAL USERID **********/
    comment.uid = 1;

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
      xhr.send(JSON.stringify(comment));
}


function getComments(pid) {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
        console.log(xhr.status);
        if (xhr.status == 200) {
          var c = event.target.response;
          
          c = JSON.parse(c);
          displayComments(c.Comments);

        } else {
          alert("Error! Get files failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Get file failed. Cannot connect to server.");
      };
        
      xhr.open('GET', 'http://130.240.170.62:1026/comments/' + pid, false);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);
}

function displayComments(comments){
    if (comments != null) {
        console.log(comments);
        comments.forEach(function(c){
            document.getElementById('comments').innerHTML += "<div class='row' ><div class='col-sm-10 col-lg-offset-1'><div class='panel panel-default'><div class='panel-heading'><strong>"+ c.Firstname + " " + c.Lastname +"</strong> <span class='text-muted'>commented " + c.Timestamp + "</span></div><div class='panel-body'>"+ c.Comment + "</div></div></div></div> "
            
    }); 
  }
}
