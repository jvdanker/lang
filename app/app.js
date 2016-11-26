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
config(['$mdThemingProvider', '$locationProvider', '$routeProvider', '$mdIconProvider', function($mdThemingProvider, $locationProvider, $routeProvider, $mdIconProvider) {
  //$mdThemingProvider.theme('default').dark();

  $locationProvider.hashPrefix('!');

  $routeProvider.otherwise({redirectTo: '/view1'});

  $mdIconProvider
         .iconSet('social', 'img/icons/sets/social-icons.svg', 24)
         .defaultIconSet('img/icons/sets/core-icons.svg', 24);
}]);
