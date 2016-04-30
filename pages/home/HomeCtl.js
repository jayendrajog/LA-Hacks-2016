marktai.controller("HomeCtl", ["$scope", "$rootScope", "$http", "$location", function($scope, $rootScope, $http, $location) {
	$rootScope.info = ""
	$rootScope.error = ""
	$rootScope.page = "home";
	$scope.posts = [];
	$scope.postsMap = {}


	var getPosts = function() {
		$http.get($rootScope.apiLocation + '/posts').then(function(result){
			$scope.posts = result.data["Posts"]
			_ = function(post) {
				$http.get($rootScope.apiLocation + '/posts/' + post).then(function(result){ 
					$scope.postsMap[post] = result.data["Post"]
				})
			}
			for (var post of $scope.posts) {
				_(post)
			}
		})
	}

	$scope.generateTag = function(post) {
		if (typeof($scope.postsMap[post]) === 'undefined') {
			return ""
		}

		var retStr = "Posted"
		if ($scope.postsMap[post]["Author"] !== "") {
			retStr += " by " + $scope.postsMap[post]["Author"]
		}
		if ($scope.postsMap[post]["Created"] !== "") {
			retStr += " on " + new Date($scope.postsMap[post]["Created"]).toDateString()
		}
		if (retStr == "Posted") {
			return ""
		}
		return retStr
	}

}])
