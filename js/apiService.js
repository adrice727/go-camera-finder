
(function(){

  var apiService = {};

  apiService.get = function(origin, value, target){

    var endpoint = '/' + origin.id + '/' + value;
    var apiCall = new XMLHttpRequest();    
    apiCall.open("GET", endpoint, false);
    apiCall.send();
    var models = document.createElement("div");
    models.setAttribute("id", "models");
    models.innerHTML = apiCall.responseText;
    document.getElementById("main").appendChild(models);
  }

  return apiService;

  console.log('HERE');
})();