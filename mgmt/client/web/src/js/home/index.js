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
    .controller('HomeCtrl', function($scope, Restangular) {
      var hooks = Restangular.all('hooks');

      var refresh = function() {
        hooks.getList().then(function(hooks) {
          $scope.hooks = hooks;
        });
      };

      refresh();

      $scope.addhook = function() {
        var newhook = {
          url: "hehe",
          rate: 100
        };
        hooks.post(newhook).then(refresh);
      };

      $scope.delhook = function(hook) {
        Restangular.one('hooks', hook).remove().then(refresh);
      };
    });


})();
