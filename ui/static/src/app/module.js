define([
    './controller/dashboard',
],
function (dashboardController) {

    var app = angular.module('workbench', ['ngRoute', 'ngMaterial'])
    .config(function($mdThemingProvider) {
        $mdThemingProvider.theme('default');
    });

    app.config(['$routeProvider', function($routeProvider){
        $routeProvider
            .when('/dashboard', {
                templateUrl: '/ui/src/app/view/dashboard.html',
                controller: 'dashboardController',
            })
            .otherwise({redirectTo: '/dashboard'});
    }]);

    app.controller('dashboardController', dashboardController);

});
