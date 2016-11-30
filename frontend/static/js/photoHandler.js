document.getElementById("catPhotoUpload").onchange = uploadPhoto;

function uploadPhoto(){
  var file = document.getElementById("catPhotoUpload").files[0];
  console.log(file);
  /*
  if (file.type.match('image.*')){
    document.getElementById('catPhoto').src = ""; 
    return;
  }
  */
  var readFile = new FileReader();
    
  function success(e){
    document.getElementById('catPhoto').src = e.target.result; 

  };

  readFile.readAsDataURL(file);
  // function call to success after file is read, no () needed
  readFile.onload = success;
}

function downloadPhoto(text, name, type) {
  var a = document.getElementById("a");
  var file = new Blob([text], {type: type});
  a.href = URL.createObjectURL(file);
  a.download = name;
}

function test() {
    window.html2canvas([ document.getElementById('catPhotoHolder') ], {
            onrendered: function( canvas ) {
                document.body.appendChild(canvas);
                window.open(canvas.toDataURL());

                                }
                            });
    /*
    html2canvas( [ document.getElementById('catPhotoHolder') ], {
        onrendered: function( canvas ) {
            saveData(canvas.toDataURL());
        }
    });
*/
}
