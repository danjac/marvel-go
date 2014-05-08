'use strict';

/* Services */

var comicsServices = angular.module('comicsServices', ['ngResource']);

comicsServices.factory('Comic', ['$resource',
      function ($resource) {
        return $resource('/api/comics/:comicId', {}, {
            query: {method: 'GET',
                  params: {comicId: ''}, isArray: true}
        });
    }]);
