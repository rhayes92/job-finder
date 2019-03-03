
jQuery(document).ready(function () {
$( "#submit" ).click(function() {
  url = "http://localhost:8080/jobs/homepage.html?Username=" + $( "#email").val();
  window.open(url,"_self");
});
});
