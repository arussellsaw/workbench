define([
    './controller/dashboard',
],
function (dashboardController) {

    var app = angular.module('workbench', ['ngRoute', 'ngMaterial'])
    .config(function($mdThemingProvider) {
        $mdThemingProvider.theme('default')
            .primaryPalette('pink')
            .accentPalette('orange');
    });

    //module config
    app.config(['$routeProvider', function($routeProvider){
        $routeProvider
            .when('/dashboard', {
                templateUrl: '/ui/src/app/view/dashboard.html',
                controller: 'dashboardController',
            })
            .otherwise({redirectTo: '/dashboard'});
    }]);

    //register controllers
    app.controller('dashboardController', dashboardController);

});
