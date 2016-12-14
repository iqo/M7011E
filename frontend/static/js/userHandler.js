//handels googlesing in and send auth token and first //last name on succefull login 

function onSignIn(googleUser) {
	var auth2 = gapi.auth2.init();
	var user={};
	var profile = auth2.currentUser.get().getBasicProfile();
	var xhr = new XMLHttpRequest();
	if (auth2.isSignedIn.get()) {
    location.reload();
		//get basic info and the google token 
		user.firstname = profile.getGivenName();
		user.lastname = profile.getFamilyName();
		user.googletoken = profile.getId();
		user.authtoken = googleUser.getAuthResponse().id_token;
		xhr.open('POST', 'http://130.240.170.62:1026/user');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.onload = function() {
			console.log('Signed in as: ' + xhr.responseText);
      getCurrentUserId();
		};
		xhr.send(JSON.stringify(user));
	}
}

function logOut() {
	var auth2 = gapi.auth2.getAuthInstance();
	auth2.signOut().then(function() {
    document.cookie
		console.log('User signed out.');
    location.reload();
  });
}

function onLoad() {
	gapi.load('auth2,signin2', function() {
		var auth2 = gapi.auth2.init();
		auth2.then(function() {
          // Current values
          var isSignedIn = auth2.isSignedIn.get();
          var currentUser = auth2.currentUser.get();

          if (!isSignedIn) {
            // Rendering g-signin2 button.
            gapi.signin2.render('google-signin-button', {
            	'onsuccess': 'onSignIn', 
            	'theme': 'dark',
            	'onfailure': 'onFailure'
            });

          } else {
           //document.getElementById("id").innerHTML = getUser(profile.getId());

         }
       });
	});
}


function getUser(token) {
  if (token == null){
    token = document.cookie.split("; ")[1].split("=")[1];
    console.log(token)
  }
  var xhr = new XMLHttpRequest();
  xhr.onload = function(event) {
    //console.log("console log: " + xhr.status);
    if (xhr.status == 200) {
      var usr = event.target.response;
          usr = JSON.parse(usr);
          //console.log(usr);
          //window.open(photo.Image);
          //document.getElementById("test_user").innerHTML = usr.GoogleToken;
        } else {
          alert("Error! Get user token failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Get user token failed. Cannot connect to server.");
      };

      xhr.open('GET', 'http://130.240.170.62:1026/user/' + token, false);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);

    }
    function checkAuth(){
      gapi.load('auth2,signin2', function() {
      var auth2 = gapi.auth2.init();
     
      auth2.then(function() {
          // Current values
          var isSignedIn = auth2.isSignedIn.get();
          if (!isSignedIn) {

            location.href       
          } else {
            var profile = auth2.currentUser.get().getBasicProfile();
           getCurrentUserId(profile.getId());
         }
       });
  });
    }
    //returns the database id of the currentlty logged in user 
    function getCurrentUserId(token) {
      var xhr = new XMLHttpRequest();
      xhr.onload = function(event) {
        console.log(xhr.status);
        if (xhr.status == 200) {
          var usr = event.target.response;
          usr = JSON.parse(usr);
          //return the usertoken of currently logged in user 
          //returnUserId(usr.Id);
          document.cookie = "id="+usr.Id;
          //returnUserId
          console.log(document.cookie)
        } else {
          alert("Error! Get user id failed");
        }
      };
      xhr.onerror = function() {
        alert("Error! Get user id failed. Cannot connect to server.");
      };

      xhr.open('GET', 'http://130.240.170.62:1026/user/' + token, false);
      xhr.setRequestHeader("Content-Type", "application/json");
      xhr.send(null);

    }
//return the user token
function returnUserId(userId){
  token = document.cookie.split("; ")[1].split("=")[1];
  console.log(token)
  return token;
} 