var defaultPage = "home";

var marktai = new angular.module("marktai", ["ngTouch","ngResource", 'ngRoute', "ngSanitize", "ngWebsocket"]);


// Sets up default page to be login and redirects every other one to 
marktai.config(function($routeProvider, $locationProvider) {
  $routeProvider

    .when('/404', {
      templateUrl : './pages/404/404.html',
      controller : "404Ctl"
    })

    .when('/home', {
      templateUrl : './pages/home/home.html',
      controller : "HomeCtl"
    })

    .when('/posts', {
      templateUrl : './pages/posts/posts.html',
      controller : "PostsCtl"
    })

    .when('/about', {
         templateUrl : './pages/about/about.html',
         controller : "AboutCtl"
    })

    .when('/contact', {
      templateUrl : './pages/contact/contact.html',
      controller : "ContactCtl"
    })

	.when('/shortlink', {
      templateUrl : './pages/shortlink/shortlink.html',
      controller : "ShortlinkCtl"
    })


    // old url
	.when('/login', {
			
      controller : function() {
        window.location.replace('/meta-tic-tac-toe');
      },
      template : "<div></div>"
	})

    // causes no path to go to default page
    .when('', {
      redirectTo : function () {
        return "/" + defaultPage
      }
    })
    .when('/', {
      redirectTo : function () {
        return "/" + defaultPage
      }
    })


    // causes unrecognized path to go to default page
    // redirects to login if not logged in
    .otherwise({
      redirectTo : function (a, b, c) {
        return "/404#/#"+b;
      }
    })

    // $locationProvider.html5Mode(true);

});


// enables attribute ng-enter which calls a function when enter is pressed
marktai.directive('ngEnter', function () {
    return function (scope, element, attrs) {
        element.bind("keydown keypress", function (event) {
            if(event.which === 13) {
                scope.$apply(function (){
                    scope.$eval(attrs.ngEnter);
                });
 
                event.preventDefault();
            }
        });
    };
});

// inject the $resource dependency here
marktai.controller("MainCtl", ["$scope", "$rootScope", "$resource", "$location",  function($scope, $rootScope, $resource, $location){

  $rootScope.info = "";
  $rootScope.error = "";
  $rootScope.infoDiv = "";
  $rootScope.firstHit = true;
  $rootScope.oldPage = "";

  if ($location.absUrl().indexOf("/beta") != -1) {
     $rootScope.apiLocation = "/apiBeta"
  } else {

	$rootScope.apiLocation = "/api"
  }
}]);
