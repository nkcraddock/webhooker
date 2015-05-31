(function() {
  'use strict';

  var angular = require('angular');

  var ngModule = angular.module('app.home', []);

  ngModule
    .config(function($stateProvider, $urlRouterProvider) {
      $stateProvider
        .state('app.home', {
          url: '/',
          views: {
            '@': {
              templateUrl: "home/home.html",
              controller: 'HomeCtrl'
            }
          }
        });
    })
    .controller('HomeCtrl', function($scope) {
      console.log("fuck");

    });


})();
