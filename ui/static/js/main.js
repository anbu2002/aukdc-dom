var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

document.getElementById("print").addEventListener("load", window.print());
document.getElementById("openbtn").addEventListener("click", function(){
	  var x = document.getElementById("rightside");
	  if (x.style.display === "none") {
	    x.style.display = "block";
	  } else {
	    x.style.display = "none";
	  }
});
document.getElementById("closebtn").addEventListener("click", function(){
	document.getElementById("sidebar").style.width = "0px";
});
$(document).ready(function() {
    $('input[type="number"]').val('');
  });
