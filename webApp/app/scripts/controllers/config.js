'use strict';

/**
 * @ngdoc function
 * @name mgmApp.controller:ConfigCtrl
 * @description
 * # ConfigCtrl
 * Controller of the mgmApp
 */
angular.module('mgmApp')
  .controller('ConfigCtrl', function ($scope, $routeParams, $location, mgm) {

    if ($scope.auth === undefined) {
      mgm.pushLocation($location.url());
      $location.url("/");
    }

    $scope.regions = [];
    $scope.estates = mgm.estates;
    $scope.currentEstate = "";
    $scope.currentRegion = "";

    //list regions when estate is selected
    $scope.displayEstate = function () {
      populateRegions();
      if ($routeParams['estate'] !== $scope.currentEstate.ID) {
        console.log("redirecting to estate");
        $location.url('/config/' + $scope.currentEstate.ID);
      }
    }

    function populateRegions() {
      var regions = [];
      for (var i = 0; i < $scope.currentEstate.Regions.length; i++) {
        regions.push(mgm.regions[$scope.currentEstate.Regions[i]]);
      }
      $scope.regions = regions;
    }

    $scope.displayConfig = function () {
      console.log('/config/' + $scope.currentEstate.ID + "/" + $scope.currentRegion.UUID);
      $location.url('/config/' + $scope.currentEstate.ID + "/" + $scope.currentRegion.UUID);
    }

    //assign variables from url, where possible
    if ($routeParams['estate'] !== undefined) {
      $scope.currentEstate = mgm.estates[$routeParams['estate']];
      populateRegions()
    }
    if ($routeParams['region'] !== undefined) {
      $scope.currentRegion = mgm.regions[$routeParams['region']];
    }
  });
