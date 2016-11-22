'use strict';

angular.module('myApp').
factory('game', ['$http', '$q', 'words', function($http, $q, words) {
    var i = 1;

    return {
        newGame: function() {
            console.log('new game');
            i = 1;

            return $http.get('http://localhost:8080/v1/game/new').then(function (response) {
                var id = response.data.Id;
                return words.getNext(id).then(function(data) {
                    return {
                        gameId: id,
                        word: data
                    };
                });
            });
        },
        saveAnswer: function(gameId, word, answer) {
            console.log('save answer', gameId, word, answer);
            return;
        },
        nextWord: function(gameId) {
            if (i === 3) {
                return $q.when();
            }

            i++;
            return words.getNext(gameId);
        },
        getResults: function(gameId) {
            return 'results';
        }
    };
}]).
factory('words', ['$http', function($http) {
    var i = 1;

    return {
        reset: function() {
            console.log('reset');
        },
        getNext: function(gameId) {
            return $http.get('http://localhost:8080/v1/game/word').then(function(data) {
                return data.data;
            });
        }
    };
}]);
