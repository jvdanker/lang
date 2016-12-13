'use strict';

angular.module('myApp.view2', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view2', {
    templateUrl: 'view2/view2.html',
    controller: 'View2Ctrl'
  });
}])

.controller('View2Ctrl', ['$scope', '$location', '$interval', 'game', 'words', function($scope, $location, $interval, game, words) {
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

    var stopTime = $interval(updateTime, 1000);
    var startTime = new Date();
    var seconds = 0;

    // used to update the UI
    function updateTime() {
        seconds = Math.round((new Date() - startTime) / 1000);
        $scope.seconds = seconds;

        if (seconds > 30) {
            $interval.cancel(stopTime);
            $location.path( '/view3' );
        }
    }

    // listen on DOM destroy (removal) event, and cancel the next UI update
    // to prevent updating time after the DOM element was removed.
    $scope.$on('$destroy', function() {
        $interval.cancel(stopTime);
    });

}]);