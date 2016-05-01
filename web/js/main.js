$(document).ready(function() {
    "use strict";
    console.log("main.js loaded");
    myoTest();
});

function myoTest() {
	Myo.connect('com.myojs.main');

	//Start talking with Myo Connect
	//Myo.connect('com.stolksdorf.myAwesomeApp');
	console.log('Hello Myo!');
	//this.vibrate();

	Myo.on('fist', function(){
		console.log('fist!');
		//this.vibrate();
	});

	Myo.on('wave_out', function(){
		console.log('wave_out!');
		//this.vibrate();
	});

	Myo.on("fingers_spread", function() {
		console.log('fingers_spread!');
		//this.vibrate();
	});

	Myo.on("wave_in", function() {
		console.log('wave_in!');
		//this.vibrate();
	});

	Myo.on("double_tab", function() {
		console.log('double_tab!');
		//this.vibrate();
	});

}