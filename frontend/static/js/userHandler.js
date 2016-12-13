//handels googlesing in and send auth token and first //last name on succefull login 

function onSignIn(googleUser) {
	var auth2 = gapi.auth2.init();
	var user={};
	var profile = auth2.currentUser.get().getBasicProfile();
	var xhr = new XMLHttpRequest();
	if (auth2.isSignedIn.get()) {
		//get basic info and the google token 
		user.firstname = profile.getName();
		user.lastname = profile.getFamilyName();
		user.googletoken = profile.getId();
		user.authtoken = googleUser.getAuthResponse().id_token;
		xhr.open('POST', 'http://130.240.170.62:1026/user');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.onload = function() {
			console.log('Signed in as: ' + xhr.responseText);
		};
		xhr.send(JSON.stringify(user));
	}
}

function logOut() {
	var auth2 = gapi.auth2.getAuthInstance();
	auth2.signOut().then(function() {
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
          	'<h1>you need to be loged in to view this</h1>';
            // Rendering g-signin2 button.
            gapi.signin2.render('google-signin-button', {
            	'onsuccess': 'onSignIn', 
            	'theme': 'dark',
            	'onfailure': 'onFailure'
            });

        }else{
        	var profile = auth2.currentUser.get().getBasicProfile();
        	console.log('Image URL: ' + profile.getImageUrl());
        	document.getElementById("logout").innerHTML = "<button onclick='logOut()'>Sign out</button>";
         	document.getElementById("user").innerHTML ="Good day " + profile.getName() + " are you ready for some cats in hats? ";
        	document.getElementById("image").src = profile.getImageUrl();
        }
    });
	});
}


function getUser(token) {
    var xhr = new XMLHttpRequest();
    xhr.onload = function(event) {
        console.log(xhr.status);
        if (xhr.status == 200) {
          var usr = event.target.response;
          //console.log(photo[0].Image);
          usr = JSON.parse(usr);
          //window.open(photo.Image);
          document.getElementById("test_user").src = usr.GoogleToken;
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
