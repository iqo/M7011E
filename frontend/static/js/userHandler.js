//handels googlesing in and send auth token and first //last name on succefull login 

function successfulSignIn(googleUser) {
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
	});
}