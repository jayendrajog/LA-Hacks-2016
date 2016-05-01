
// function loginfunc()
// {
// 	var url = window.location.href;
// 	console.log(url);
// 	if(url.indexOf("https://www.facebook.com/") > -1)
// 		document.getElementById('email').value = "jay08jog@gmail.com";

// 	if(url.indexOf("https://shb.ais.ucla.edu/") > -1)
// 		document.getElementById('logon').value = "jayendra";

// 	if(url.indexOf("https://quora.com") > -1)
// 	{
// 		document.getElementsByName('email')[0].value = "jayendra@ucla.edu";
// 		document.getElementsByName('password')[0].value = "testing";
// 		document.forms[0].submit();
// 	}
// }
//alert('something happens2');

// chrome.browserAction.onClicked.addListener(function(tab) {
//   chrome.tabs.executeScript(null, { file: "jquery.min.js" }, function() {
//     chrome.tabs.executeScript(null, { file: "content_script.js" });
//   });
// });

//should only get called once ever
//get the user id
//$.get()

var url = window.location.href;
if(url.indexOf("https://www.facebook.com/") > -1)
{
	$.get( "https://www.marktai.com/faceApi/creds?domain=www.facebook.com",  function( data ) {
		var parsed = JSON.parse(data);
	
	 	document.getElementById('email').value = parsed['Username'].toString();
	 	document.getElementById('pass').value = parsed['Password'].toString();
	});
}

else if (url.indexOf("https://shb.ais.ucla.edu/") > -1)
{
	$.get( "https://www.marktai.com/faceApi/creds?domain=ucla.edu", function( data ) {
		var parsed = JSON.parse(data);
		
	 	document.getElementById('logon').value = parsed['Username'].toString();
	 	document.getElementById('pass').value = parsed['Password'].toString();
	});
}








// window.onload = function(){

// 	chrome.tabs.onUpdated.addListener(function(tabId, changeInfo, tab) {
// 		 if(changeInfo && changeInfo.status == "complete"){
// 		 	chrome.tabs.query({'active': true, 'lastFocusedWindow': true}, function (tabs) {
// 			    var url = tabs[0].url;
			   
// 			    //facebook case
// 			    if(url.indexOf("https://www.facebook.com/") > -1)
// 	       		{
// 		       		$.get( "http://sparck.co/api/creds", function( data ) {
// 		       			var parsed = JSON.parse(data);
// 		       			alert(parsed['Username'].toString() + " " + parsed['Password'].toString());
		       			
// 					 	document.getElementById('email').value = parsed['Username'].toString();
// 					 	document.getElementById('password').value = parsed['Password'].toString();

// 				});
// 	       		}
// 			});

		 	

// 		 	//var url = window.location.href;
// 		 	//alert(url);
	       	

// 	       	//myucla case

	       
//     	}   
// 	});
// }

	// chrome.tabs.onUpdated.addListener(function(tabId, changeInfo, tab) {
	//    alert('updated');
	// });



// chrome.tabs.onCreated.addListener(function(tabId, changeInfo, tab) {         
//     alert('created');
// });

//loginfunc();

// function autolog() {
// 	alert('something happens');
// 	chrome.tabs.executeScript(null, {
// 		code: 'loginfunc();'
// 	});
// }

// window.onload = function() {
// 	document.getElementById('login').addEventListener('click', autolog);
// }