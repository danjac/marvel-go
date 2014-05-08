'use strict';

/* App Module */

var comicsApp =  angular.module('comicsApp', [
    'ngRoute',
    'comicsControllers',
    'comicsServices'
]).config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
            when('/comics', {
                templateUrl: 'partials/comic-list.html',
                controller: 'ComicListCtrl'
            }).
            when('/comics/:comicId', {
                templateUrl: 'partials/comic-detail.html',
                controller: 'ComicDetailCtrl'
            }).
            otherwise({
                redirectTo: '/comics'
            });
    }]);
