'use strict';

// Declare app level module which depends on views, and components
angular.module('myApp', [
  'ngMaterial',
  'ngRoute',
  'myApp.view1',
  'myApp.view2',
  'myApp.view3',
  'myApp.view4',
  'myApp.version'
]).
config(['$mdThemingProvider', '$locationProvider', '$routeProvider', function($mdThemingProvider, $locationProvider, $routeProvider) {
  $mdThemingProvider.theme('default').dark();

  $locationProvider.hashPrefix('!');

  $routeProvider.otherwise({redirectTo: '/view1'});
}]);
