'use strict';

angular.module('myApp.view2', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view2', {
    templateUrl: 'view2/view2.html',
    controller: 'View2Ctrl'
  });
}])

.controller('View2Ctrl', ['$scope', '$location', '$interval', 'game', 'words', function($scope, $location, $interval, game, words) {
    function getScore() {
        var score = localStorage.getItem('score');
        if (score) {
            score = Number(score);
        } else {
            score = 0;
        }
        return score;
    }

    game.newGame().then(function (game) {
        $scope.game = game;
        $scope.score = getScore();
    });

    $scope.save = function(gameId, word, answer) {
        $scope.message = '';
        delete $scope.answer;
        var result = game.saveAnswer(gameId, word, answer);

        var score = localStorage.getItem('score');
        if (score) {
            score = Number(score);
        } else {
            score = 0;
        }

        if (answer == word.Correct) {
            if (score) {
                score += 4;
            } else {
                score = 4;
            }

            game.nextWord(gameId).then(function(word) {
                $scope.game.word = word;
            });
        } else {
            $scope.message = 'incorrect';
            if (score) {
                score -= 10;
            }
        }

        if (score < 0) {
            score = 0;
        }

        console.log(score);
        localStorage.setItem('score', score);
        $scope.score = score;
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
            //$location.path( '/view3' );
        }
    }

    // listen on DOM destroy (removal) event, and cancel the next UI update
    // to prevent updating time after the DOM element was removed.
    $scope.$on('$destroy', function() {
        $interval.cancel(stopTime);
    });

}]);