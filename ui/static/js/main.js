var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

var openbtn = document.getElementById("openbtn");
if (openbtn) {
	openbtn.addEventListener('click', toggle);
}

function toggle(){
	  var x = document.getElementById("rightside");
	  if (x.style.display === "none") {
		x.style.display = "block";
	  } else {
		x.style.display = "none";
	  }
}

var uploads = document.querySelectorAll(".imageInput")
var displays = document.querySelectorAll(".imagePreview")
if (uploads){
	for (let i = 0; i < uploads.length; i++){
		uploads[i].addEventListener("change", function () {
		container = displays[i]
		displayImage(this, container);
	  });
	}
}


function displayImage(input, container) {
	if (input.files && input.files[0]) {
	  const reader = new FileReader();
	  reader.onload = function (e) {
		const img = container.querySelector("img");
		img.src = e.target.result;
	  };
	  reader.readAsDataURL(input.files[0]);
	}
}

var stepbutton = document.querySelectorAll('.stepButton')
if (stepbutton) {
	for (let i = 0; i < stepbutton.length; i++){
		stepbutton[i].addEventListener("click", function () {
		showStep(i);
	  });
	}
}
function showStep(from) {
  const steps = document.querySelectorAll('.tab');
  var active = from^1
  steps[active].classList.add('active'); 
  steps[from].classList.remove('active');
}

var qc = document.getElementById('qc');
if (qc) {
	qc.addEventListener("input", function() {
		inputValue = parseInt(qc.value);
		if (inputValue == 2){
			none = document.querySelectorAll('.none');
			none[0].classList.add('active');
			none[0].classList.remove('none');
			document.getElementById('kc').disabled = true;
		}else{
			none[0].classList.remove('active');
			none[0].classList.add('none');
			document.getElementById('kc').value = 0;
			document.getElementById('kc').disabled = false;
		}
	});
}
