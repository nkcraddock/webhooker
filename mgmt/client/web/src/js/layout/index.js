(function() {
  'use strict';

  var angular = require('angular');

  var ngModule = angular.module('app.layout', []);

  ngModule
    .config(function($stateProvider, $urlRouterProvider) {
      $stateProvider
        .state('app', {
          abstract: true,
          views: {
            'navbar': {
              templateUrl: "layout/navbar.html",
              controller: 'NavCtrl'
            },
            'footer': {
              templateUrl: "layout/footer.html"
            }
          }
        });
    })
    .controller('NavCtrl', function($scope) {
    });


})();
