// // function onWindowLoad() {
// // 	alert('window loads');
// // 	 // var message = document.querySelector('#message');

// // 	chrome.tabs.executeScript(null, 
// // 		{file: "autolog.js"}, 
// // 	  	function() {
// // 	    // If you try and inject into an extensions page or the webstore/NTP you'll get an error
// // 	 //   if (chrome.runtime.lastError) {
// // 	  //    message.innerText = 'There was an error injecting script : \n' + chrome.runtime.lastError.message;
// // 	  //  }
// // 	 	 });

	
// // }

// // window.onload = onWindowLoad;


// function autolog() {
// 	chrome.tabs.getSelected(null, function(tab){
// 		var url = tab.url;
// 		console.log(url);
// 		if(url.indexOf("https://www.facebook.com/") > -1)
// 		{
// 			chrome.tabs.executeScript(null, {
// 				code: 'document.getElementById("email").value = "jay08jog@gmail.com"'
// 			});
// 		}
// 		else if (url.indexOf("https://shb.ais.ucla.edu/") > -1)
// 			chrome.tabs.executeScript(null, {
// 				code: 'document.getElementById("logon").value = "jayendra"'
// 			});
// 	});

// 	//alert('something happens');
// 	// chrome.tabs.executeScript(null, {
// 	// 	code: 'document.getElementById("email").value = "jay08jog@gmail.com"'
// 	// });
// }

// window.onload = function() {
// 	document.getElementById('login').addEventListener('click', autolog);
	
// }