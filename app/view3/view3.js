'use strict';

angular.module('myApp.view3', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view3', {
    templateUrl: 'view3/view3.html',
    controller: 'View3Ctrl'
  });
}])

.controller('View3Ctrl', ['$scope', '$location', 'game', function($scope, $location, game) {
    $scope.results = game.getResults();

    $scope.back = function() {
        $location.path("/");
    };
}]);