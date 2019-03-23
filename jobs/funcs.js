// jQuery(document).ready(function () {
//
//
// });
var jsonVal = ""
var catVal = 0
var divVal = 0
function getCat() {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      //console.log(this.responseText);
      var  CatAndDiv = JSON.parse( this.responseText );
      jsonVal = CatAndDiv
      console.log(jsonVal.categories)
    //  auto()
    }
  };
  xhttp.open("GET", "http://localhost:8080/cat");
  xhttp.send();
}

function auto() {
  $( function() {
    $(".autocomplete1").autocomplete({source:jsonVal.categories});
    $(".autocomplete2").autocomplete({source:jsonVal.divisions});
   } );
}

jQuery(document).ready(function () {
  $( "#catAdd" ).click(function() {

    $( "#catCol" ).append( '<br/><input class="autocomplete1" name="catText"> Rating:</input> <select name="catSelect"> <option value="1">1</option>  <option value="2">2</option>  <option value="3">3</option>  <option value="4">4</option>  <option value="5">5</option>  <option value="6">6</option>  <option value="7">7</option></select>' );
    catVal = catVal + 1
    auto();
  });
  $( "#divisionAdd" ).click(function() {
    $( "#divCol" ).append( '<br/><input  class="autocomplete2" name="divText"> Rating:</input> <select name="divSelect">  <option value="1">1</option>  <option value="2">2</option>  <option value="3">3</option>  <option value="4">4</option>  <option value="5">5</option>  <option value="6">6</option>  <option value="7">7</option></select>' );
    divVal = divVal + 1
    auto();
  });

  $( "#eval" ).click(function() {



     cat =$('input[name=catText]')
     catRank =$('select[name=catSelect]')
     cats = []
     for(var i = 0; i <cat.length; i++){
       cats.push({category:cat[i].value, rank:parseFloat(catRank[i].value)})
     }

     div =$('input[name=divText]')
     divRank =$('select[name=divSelect]')
     divs =[]
     for(var i = 0; i <div.length; i++){
       divs.push({category:div[i].value, rank:parseFloat(divRank[i].value)})
     }
    var eval = {categories:cats, divisions:divs};

    //
    var xhr = new XMLHttpRequest();
    var url = "http://localhost:8080/ScoreEval";
    xhr.open("POST", url);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
        var json = JSON.parse(xhr.responseText);
        console.log(json)

    }
  };
  var data = JSON.stringify(eval);
  xhr.send(data);
  });
});