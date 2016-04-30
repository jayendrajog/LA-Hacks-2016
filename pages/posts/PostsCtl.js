marktai.controller("PostsCtl", ["$scope", "$rootScope", "$http", "$location", "$sce",  function($scope, $rootScope, $http, $location, $sce) {
	$rootScope.page = "posts";

	$scope.post = $location.hash()

	$scope.postData = {}

	var getPost = function(post) {
		$http.get($rootScope.apiLocation + '/posts/' + post).then(function(result){ 
			$scope.postData = result.data["Post"]
            for (var par of $scope.postData.Content) {
                $scope.pars.push($sce.trustAsHtml(par))
            }

		}, function(error) {
			if (error.status === 404) {

				var url = "/#" + $location.url();
				console.log(url)
				$location.path("/404");
				$location.hash(url);
			}
				
		})
	}

    $scope.pars = []

	$scope.generateTag = function() {
		var retStr = "Posted"
		if ($scope.postData["Author"] !== "") {
			retStr += " by " + $scope.postData["Author"]
		}
		if ($scope.postData["Created"] !== "") {
			retStr += " on " + new Date($scope.postData["Created"]).toDateString()
		}
		if (retStr == "Posted") {
			return ""
		}
		return retStr
	}

	getPost($scope.post)

}])
