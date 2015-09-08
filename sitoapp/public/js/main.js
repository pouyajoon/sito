(function() {
  'use strict';
  /*global WebSocket, angular*/

  var app, ws;
  ws = new WebSocket('ws://' + window.location.host + '/ws');

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
    document.getElementById('player').value = 'Guest ' + Math.floor(Math.random() * 1e3);
    sendMousePosition();
    window.onmousemove = sendMousePosition;
  };

  ws.onclose = function(e) {
    console.log('ws close', e.data);
  };

  app = angular.module('app', []);

  function bin2String(array) {
    var result = '',
      i;
    for (i = 0; i < array.length; i += 1) {
      result += String.fromCharCode(parseInt(array[i], 2));
    }
    return result;
  }

  app.controller('playerController', function($scope) {
    $scope.title = 'sito';
    ws.onmessage = function(e) {
      $scope.m = JSON.parse(e.data);
      // console.log('ws message', $scope.m);
      $scope.$apply();
      return e;
    };
  });



  console.log('hello', window.location.host);
}());
