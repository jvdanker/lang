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
        what: 'Time controlled',
        who: 'New game',
        notes: "Get as many words correct in 30 seconds",
        url: '/view2'
      },
      {
        what: 'Dictionary',
        who: 'Add a new word',
        notes: "Add a new word to your dictionary",
        url: '/view4'
      },
      /*
      {
        what: 'Dictionary',
        who: 'Word list',
        notes: "List all of the words in your dictionary"
      },
      */
    ];

    $scope.go = function(item) {
        $location.path(item.url);
    };
}]);