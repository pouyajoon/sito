(function() {
  'use strict';
  /*global WebSocket, angular*/

  var ws = new WebSocket('ws://' + window.location.host + '/ws');

  function sendMousePosition(e) {
    if (e === undefined) {
      e = {
        pageX: -1,
        pageY: -1
      };
    }
    var d, ds, p = document.getElementById('player');
    d = {
      'player': p !== null ? p.value : 'P.',
      'x': e.pageX,
      'y': e.pageY
    };
    ds = JSON.stringify(d);
    // console.log('send d', ds);
    ws.send(ds);
  }


  ws.onopen = function() {
    console.log('ws open');
    document.getElementById('player').value = 'Guest ' + Math.random();
    sendMousePosition();
    window.onmousemove = sendMousePosition;
  };

  ws.onclose = function(e) {
    console.log('ws close', e.data);
  };

  var app = angular.module('app', []);

  app.controller('playerController', function($scope) {
    $scope.title = 'sito';
    ws.onmessage = function(e) {
      $scope.m = JSON.parse(e.data);
      console.log('ws message', $scope.m);
      return e;
    };
  });



  console.log('hello', window.location.host);
}());
