// jQuery(document).ready(function () {
//
//
// });
var jsonCat = ""
var jsonEval = ""
var catVal = 0
var divVal = 0
function getCat() {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      //console.log(this.responseText);
      var  CatAndDiv = JSON.parse( this.responseText );
      jsonCat = CatAndDiv
    //  auto()
    }
  };
  xhttp.open("GET", "http://localhost:8080/cat");
  xhttp.send();
}

function auto() {
  $( function() {
    $(".autocomplete1").autocomplete({source:jsonCat.categories});
    $(".autocomplete2").autocomplete({source:jsonCat.divisions});
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
       if (cat[i].value ==""){
         i++
       }
       if (!jsonCat.categories.includes(cat[i].value)){
         alert(cat[i].value+" is invalid\nCategory must be in database(use autocomplete)")
         return
       }
       cats.push({category:cat[i].value, rank:parseFloat(catRank[i].value)})
     }

     div =$('input[name=divText]')
     divRank =$('select[name=divSelect]')
     divs =[]
     for(var i = 0; i <div.length; i++){
       if (div[i].value ==""){
         i++
       }
       if (!jsonCat.divisions.includes(div[i].value)){
         alert(div[i].value+" is invalid\nDivisions must be in database(use autocomplete)")
         return
       }
       divs.push({category:div[i].value, rank:parseFloat(divRank[i].value)})
     }

    weight  = {begin:parseFloat($('#begin').val()), end:parseFloat($('#end').val()), divisions:parseFloat($('#divisions').val()), category:parseFloat($('#category').val())}
    //if (weight.begin != 1 && weight.end != 1  && weight.category != 1 && weight.divisions != 1){
    var weightCount = [0,0,0,0]
    for(var i = 1; i < 5; i++){
      for(var index in weight) {
        if (weight[index] == i) {
          weightCount[i-1]++;
        }
      }
    }
    for(var i = 0; i < 3; i++){
       console.log(weightCount)
        if (weightCount[i] < weightCount[i+1]){
          alert("Weights must go in acsending order \n(example: can't have all weights equal 4 they must all equal 1)");
          return
        }
    }
    var eval = {categories:cats, divisions:divs, weight:weight};

    //
    var xhr = new XMLHttpRequest();
    var url = "http://localhost:8080/ScoreEval";
    xhr.open("POST", url);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
    if (xhr.readyState === 4 && xhr.status === 200) {
        var json = JSON.parse(xhr.responseText);
        //console.log(json)
        $("#result").find("img").remove();
        $("#evals").find("tr").remove();
        jsonEval = json;
        for(var index in jsonEval) {
         $( "#evals" ).append(
           '<tr>'+
          '<td class="item">'+ index + '</td>'+
          '<td class="item">'+ jsonEval[index].Score + '</td>'+
          '<td class="item">'+ jsonEval[index].BusTitle + '</td>'+
          '<td class="item">'+ jsonEval[index].JobCategory + '</td>'+
          '<td class="item">'+ jsonEval[index].DivisionUnit + '</td>'+
          '<td class="item">'+ jsonEval[index].SalaryRangeEnd + '</td>'+
          '<td class="item">'+ jsonEval[index].SalaryRangeBegin + '</td>'+
          '<td class="item"> <input type="button" id="'+index+'" data-href="#entry1"  onclick="getInfo('+index+')"/></td> </tr>'
         );
        //console.log(jsonEval[index])
       }
      $("#result").addSortWidget();
    }
    };
  var data = JSON.stringify(eval);
  console.log(eval)
  xhr.send(data);
  });
});
function getInfo(id){
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      console.log(this.responseText);
      var  info = JSON.parse( this.responseText );
      popup = "Job ID: - "+info.JobID+"\n"+
      "Title: - "+info.BusTitle+"\n"+
      "Job Category: - "+info.JobCategory+"\n"+
      "Division: - "+info.DivisionUnit+"\n"+
      "Salary Cap: $"+info.SalaryRangeEnd+"\n"+
      "Starting Salary:  $"+info.SalaryRangeBegin+"\n"
      alert( popup)


    //  auto()
    }
  };
  xhttp.open("GET", "http://localhost:8080/info?id="+id);
  xhttp.send();
}
