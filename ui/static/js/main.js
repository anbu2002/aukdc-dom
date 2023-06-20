var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

var openbtn = document.getElementById("openbtn");
openbtn.addEventListener('click', toggle);

function toggle(){
	  var x = document.getElementById("rightside");
	  if (x.style.display === "none") {
		x.style.display = "block";
	  } else {
		x.style.display = "none";
	  }
}

