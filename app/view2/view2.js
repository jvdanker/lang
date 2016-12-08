'use strict';

angular.module('myApp.view2', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view2', {
    templateUrl: 'view2/view2.html',
    controller: 'View2Ctrl'
  });
}])

.controller('View2Ctrl', ['$scope', '$location', 'game', 'words', function($scope, $location, game, words) {
    game.newGame().then(function (game) {
        $scope.game = game;
    });

    $scope.save = function(gameId, word, answer) {
        $scope.message = '';
        delete $scope.answer;
        var result = game.saveAnswer(gameId, word, answer);

        if (answer == word.Correct) {
            $scope.game.word = game.nextWord(gameId).then(function(word) {
                if (word === undefined) {
                    $location.path( '/view3' );
                } else {
                    $scope.game.word = word;
                }
            });
        } else {
            $scope.message = 'incorrect';
        }
    };

    $scope.back = function() {
        $location.path("/");
    };
}]);