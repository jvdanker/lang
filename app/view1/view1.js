'use strict';

angular.module('myApp.view1', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view1', {
    templateUrl: 'view1/view1.html',
    controller: 'View1Ctrl'
  });
}])

.controller('View1Ctrl', ['$scope', '$location', function($scope, $location) {
    var imagePath = 'img/list/60.jpeg';

    $scope.todos = [
      {
        face : imagePath,
        what: 'Time controlled',
        who: 'New game',
        when: '3:08PM',
        notes: "Get as many words correct in 30 seconds",
        url: '/view2'
      },
      {
        face : imagePath,
        what: 'Dictionary',
        who: 'Add a new word',
        when: '3:08PM',
        notes: "Add a new word to your dictionary",
        url: '/view4'
      },
      {
        face : imagePath,
        what: 'Dictionary',
        who: 'Word list',
        when: '3:08PM',
        notes: "List all of the words in your dictionary"
      },
    ];

    $scope.go = function(item) {
        $location.path(item.url);
    };
}]);