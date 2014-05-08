'use strict';

/* Controllers */

var comicsControllers = angular.module('comicsControllers', [])
    .controller('ComicListCtrl', ['$scope', 'Comic',
        function ($scope, Comic) {
            $scope.comics = Comic.query();
        }]).controller('ComicDetailCtrl', [
        '$scope',
        '$routeParams',
        'Comic',
        function ($scope, $routeParams, Comic) {
            $scope.comic = Comic.get({comicId: $routeParams.comicId},
                function (comic) {
                    $scope.images = [];
                    for (var i=0; i < comic.Images.length; i++) {
                        var image = comic.Images[0].Path + "/portrait_xlarge." + comic.Images[0].Extension;
                        $scope.images.push(image);
                    }
                });
        }]);
