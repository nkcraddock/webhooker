(function() {
  'use strict';

  var angular = require('angular');

  var ngModule = angular.module('app.hook', []);

  ngModule
    .config(function($stateProvider, $urlRouterProvider) {
      $stateProvider
        .state('app.hook', {
          abstract: true,
          url: '/hook'
        })
        .state('app.hook.new', {
          url: '/new',
          views: {
            '@': {
              templateUrl: "hook/detail.html",
              controller: 'NewHookCtrl'
            }
          }
        })
        .state('app.hook.view', {
          url: '/:hook',
          views: {
            '@': {
              templateUrl: "hook/detail.html",
              controller: 'HookDetailCtrl'
            }
          }
        });
    })
    .controller('HookDetailCtrl', function($scope, $stateParams, Restangular) {
      var hook = Restangular.one('hooks', $stateParams.hook);
      var filters = hook.all('filters');

      var refreshFilters = function() {
        filters.getList().then(function(f) {
          $scope.filters = Restangular.stripRestangular(f);
        });
      };

      hook.get().then(function(h) {
        $scope.hook = Restangular.stripRestangular(h);
      });

      refreshFilters();

      $scope.save = function() {
        $scope.hook.save();
      };

      $scope.addFilter = function() {
        var f = {
          "src": "test source",
          "evt": "something happened event",
          "key": "3"
        };

        filters.post(f).then(function(newfilter) {
          refreshFilters();
        });
      };
    })
    .controller('NewHookCtrl', function($scope, $state, Restangular) {
      $scope.isNew = true;
      $scope.hook = {
        url: "",
        rate: 100
      };
      $scope.save = function() {
        Restangular.all('hooks').post($scope.hook).then(function(newhook) {
          $state.go("app.home");
        });
      };
    });


})();
