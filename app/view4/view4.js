'use strict';

angular.module('myApp.view4', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view4', {
    templateUrl: 'view4/view4.html',
    controller: 'View4Ctrl'
  });
}])

.controller('View4Ctrl', ['$scope', '$location', '$http', 'words', function($scope, $location, $http, words) {
    $scope.words = {};

    $scope.save = function(words) {
        $http.post('http://localhost:8080/v1/game/word/add', {
            Word1: words.word1,
            Word2: words.word2
        }).then(function() {
            $location.path( '/view1' );
        });
    };

}]);