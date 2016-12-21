//handels googlesing in and send auth token and first //last name on succefull login 
var loggedIn = false;


function onSignIn(googleUser) {
	var auth2 = gapi.auth2.init();
	var user={};
	var profile = auth2.currentUser.get().getBasicProfile();
	var xhr = new XMLHttpRequest();
	if (auth2.isSignedIn.get()) {
        
		//get basic info and the google token 
		user.firstname = profile.getGivenName();
		user.lastname = profile.getFamilyName();
		user.googletoken = profile.getId();
		user.authtoken = googleUser.getAuthResponse().id_token;
		xhr.open('POST', 'http://130.240.170.62:1026/user');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.onload = function() {
            location.reload();
            loggedIn = true;

        };
		xhr.send(JSON.stringify(user));
	}
}

function logOut() {
	var auth2 = gapi.auth2.getAuthInstance();
	auth2.signOut().then(function() {
    document.cookie = "id=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
		console.log('User signed out.');
    loggedIn = false;
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
            loggedIn = true;
            var uid = returnUserId();
           //document.getElementById("id").innerHTML = getUser(profile.getId());
            document.getElementById("logout").innerHTML = '<button href="#" onclick="logOut();">Sign out</button>';
            document.getElementById('mypageMenu').innerHTML = "<a href='/mypage/" + uid + "''>My page</a>";
        
         }
       });
	});
}


function getUser(token) {
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
        if (xhr.status == 200) {
          var usr = event.target.response;
          usr = JSON.parse(usr);
          //return the usertoken of currently logged in user 
          //returnUserId(usr.Id);
          //cookies expires after a certain time 
          var oldDateObj = new Date();
          // diff is the time in minutes before cookie expires 
          var diff = 30;
          var newDateObj = new Date(oldDateObj.getTime() + diff*60000);
          document.cookie = "id="+usr.Id + "; expires=" + newDateObj.toUTCString() +"; path=/;";  //path makes sure where to save cookies 
          //returnUserId
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
  //token = document.cookie.split("; ")[1].split("=")[1];
  if (document.cookie.split("id=")[1] != undefined) {
    token = document.cookie.split("id=")[1][0];
  } else {
    token = 0;
  }
  return token;
} 

function isSignedIn(){
  return loggedIn;
}