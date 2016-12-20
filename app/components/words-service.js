'use strict';

angular.module('myApp').
factory('game', ['$http', '$q', 'words', function($http, $q, words) {
    var i = 1;

    return {
        newGame: function() {
            return $http.get('/v1/game/new').then(function (response) {
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
            // console.log('save answer', gameId, word, answer);
            return;
        },
        nextWord: function(gameId) {
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
            return $http.get('/v1/game/word').then(function(data) {
                return data.data;
            });
        }
    };
}]);
