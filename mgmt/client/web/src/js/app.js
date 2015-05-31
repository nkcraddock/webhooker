
(function() {
  'use strict';

  var angular = require('angular');

  var app = angular.module('app', [
    'ui.router',
    'restangular',
    'templates-main',
    'app.home',
    'app.layout'
  ]);

  app.config(function($locationProvider, RestangularProvider) {
    $locationProvider.html5Mode({
      enabled: true,
      requireBase: false
    });

    RestangularProvider.setBaseUrl('/api');
  });

  app.config(function($urlRouterProvider) {
    $urlRouterProvider.otherwise('/');
  });


})();
