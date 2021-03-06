var rated = false;
var ip = getIP();

function addComment(photoId) {
    if (isSignedIn()) {
        var comment={};
        comment.photoId = parseInt(photoId);
        comment.comment = document.getElementById("userComment").value;
        
        comment.uid = parseInt(returnUserId());

        var xhr = new XMLHttpRequest();
        xhr.onload = function() {
            if (xhr.status == 200) {
                console.log("new comment added")
                location.reload(); // reloads page

          } else {
              alert("Error! Comment failed");
          }
        };
        xhr.onerror = function() {
            alert("Error! Comment failed." + xhr.status);
        };

        xhr.open('POST', ip + '/comment', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.send(JSON.stringify(comment));

    } else {
        alert("You need to be signed in to add a comment!")
    }
}


function getComments(photoId) {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
        if (xhr.status == 200) {
          var c = event.target.response;
          if (c.length == 0) {
            return;
        };
        c = JSON.parse(c);
        displayComments(c.Comments);

        } else if (xhr.status == 404){
            console.log ("No comments found");
        } else {}
        };
        xhr.onerror = function() {
            alert("Error! Get comments failed. Cannot connect to server.");
      };
        
      xhr.open('GET', ip + '/comments/' + photoId, true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);
}

function displayComments(comments){
    if (comments != null) {
        comments.forEach(function(c){
            document.getElementById('comments').innerHTML += "<div class='row' ><div class='col-sm-10 col-lg-offset-1'><div class='panel panel-default'><div class='panel-heading'><strong>"+ c.Firstname + " " + c.Lastname +"</strong> <span class='text-muted'>commented " + c.Timestamp + "</span></div><div class='panel-body'>"+ c.Comment + "</div></div></div></div> "
            
    }); 
  }
}

function getRate(pid) {
    var xhr = new XMLHttpRequest();
    var uid = parseInt(returnUserId());
    xhr.onload = function(event) {
        if (xhr.status == 200) {
            var r = event.target.response;
            console.log(r);
            if (r.length == 0) {
                return;
            };
            r = JSON.parse(r);
            rated = true;
            changeColor(r.Rate);
            } else {

        }
      };
      xhr.onerror = function() {
            alert("Error! Get rate failed. Cannot connect to server.");
      };
      xhr.open('GET', ip + '/rating/' + pid + "/" + uid, true);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);
}


function rate(photoId, rate) {
    var r={};
    r.photoId = parseInt(photoId);
    r.rate = parseInt(rate);

    
    r.uid = parseInt(returnUserId());
    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
    if (xhr.status == 200) {
        updateVote(r.rate);


    } else {
        alert("Error! Rating failed");
        }
    };
    xhr.onerror = function() {
        alert("Error! Rating failed." + xhr.status);
    };

    if (!rated) {
        xhr.open('POST', ip + '/rating', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.send(JSON.stringify(r));

    } 
    if (rated) {
        xhr.open('PUT', ip + '/rating/update', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.send(JSON.stringify(r));
    }
    
}

// updating vote without reloading window
function updateVote(rate){
        var sum = document.getElementById("sum").innerHTML;
        sum = parseInt(sum) + rate;
        if (sum == 0) {sum = sum + rate;};
          voted = true;
          document.getElementById("sum").innerHTML = sum;
          changeColor(rate);
    }

function checkIfVoted(vote, id){
    if (isSignedIn()) {
        if (vote == "up") {
            if (document.getElementById("upvote").style.color == "green") {
                return;
            } else {
                rate(id, 1);
            }
        };
        if (vote == "down") {
            if (document.getElementById("downvote").style.color == "red") {
                return;
            } else {
                rate(id, -1);
            }
        };
        return;
    } else {
        alert("You need to be signed in to vote!");
    }
}


function changeColor(rate){
    console.log(isSignedIn());
    if (isSignedIn() || parseInt(returnUserId()) != 0) {
        if (rate == 1) {
            document.getElementById("upvote").style.color = "green";
            document.getElementById("downvote").style.color = "black";
        } else if (rate == -1) {
            document.getElementById("upvote").style.color = "black";
            document.getElementById("downvote").style.color = "red";
        } else {
            document.getElementById("upvote").style.color = "black";
            document.getElementById("downvote").style.color = "black";
        }
    }
}


