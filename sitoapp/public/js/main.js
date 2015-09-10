(function() {
  'use strict';
  /*global WebSocket, angular, PbZone*/

  var app, ws;

  function sendMousePosition(e) {
    if (e === undefined) {
      e = {
        pageX: 0,
        pageY: 0
      };
    }
    var d, ds, p = document.getElementById('player');
    d = {
      'player': p !== null ? p.value : 'P.',
      'x': e.pageX,
      'y': e.pageY,
      's': 20
    };
    ds = JSON.stringify(d);
    ws.send(ds);
  }


  app = angular.module('app', []);
  app.controller('playerController', function($scope) {

    ws = new WebSocket('ws://' + window.location.host + '/ws');

    ws.onopen = function() {
      console.log('ws open');
      $scope.myname = 'Guest ' + Math.floor(Math.random() * 1e3);
      sendMousePosition();
      window.onmousemove = sendMousePosition;
    };

    ws.onclose = function(e) {
      console.log('ws close', e.data);
    };

    // var z = new PbZone();
    $scope.title = 'sito';
    $scope.m = [];
    ws.onmessage = function(e) {
      $scope.m = JSON.parse(e.data);
      // console.log($scope.m);
      // z.c.animate({
      //   cx: $scope.m.x
      // }, 200);
      $scope.$apply();
      return e;
    };
  });



  console.log('hello', window.location.host);
}());
