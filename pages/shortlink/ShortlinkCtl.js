marktai.controller("ShortlinkCtl", ["$scope", "$rootScope", "$http", "$location", "$q", function($scope, $rootScope, $http, $location, $q) {
	$rootScope.info = ""
	$rootScope.error = ""
	$rootScope.page = "home";

	$scope.oldLonglink = "";

	$scope.longlink = "";

	$scope.shortlink = "";	
	$scope.id = "";

	var createShortlink = function(link) {
		var data = {"Link":link};
		return $http.post('/s', data).then(function(result){ 
			var data = result.data;
			return $q.resolve(data["ID"]);
		}, function(error){
			return $q.reject(error);
		});			
	}

	$scope.createShortlink = function() {

		// prepends with http protocol if not given
		if (!$scope.longlink.includes("://")) {
			$scope.longlink = "http://" + $scope.longlink;
		}
		if ($scope.oldLonglink !== $scope.longlink) {
			createShortlink($scope.longlink).then(function(id){
				$scope.shortlink = "https://www.marktai.com/s/" + id
				$scope.id = id;
				console.log(id);
			}, function(error){
				$rootScope.error = error;
			})
		}
		$scope.oldLonglink = $scope.longlink;
	}

}])
