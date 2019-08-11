'use strict';

/* Ctrls */

var app = angular.module('app.main', [
    'pascalprecht.translate',
    'ngCookies',
    'angularFileUpload',
    'ngStorage'
  ])
  .value('Ctx', {
    qiniu: {
      token: "",
      url_prefix: ""
    }
  })



app.controller('AppCtrl', ['$scope', '$translate', '$localStorage', '$window', 'cms',
  function($scope, $translate, $localStorage, $window, cms) {
    // config
    var vm = $scope;
    $scope.app = cms;

    // add 'ie' classes to html
    var isIE = !!navigator.userAgent.match(/MSIE/i);
    isIE && angular.element($window.document.body).addClass('ie');
    isSmartDevice($window) && angular.element($window.document.body).addClass('smart');

    // save settings to local storage
    if (angular.isDefined($localStorage.settings)) {
      $scope.app.settings = $localStorage.settings;
    } else {
      $localStorage.settings = $scope.app.settings;
    }

    function isSmartDevice($window) {
      // Adapted from http://www.detectmobilebrowsers.com
      var ua = $window['navigator']['userAgent'] || $window['navigator']['vendor'] || $window['opera'];
      // Checks for iOs, Android, Blackberry, Opera Mini, and Windows mobile devices
      return (/iPhone|iPod|iPad|Silk|Android|BlackBerry|Opera Mini|IEMobile/).test(ua);
    }
  }
]);

app.controller('MainCtrl', ['$scope', '$cookies', '$localStorage', 'AuthService', 'cms',
  function($scope, $cookies, $localStorage, AuthService, cms) {
    $scope.app = cms;
    $scope.Auth = AuthService.checksum; // 权限校验方法，页面调ng-show="Auth('keyname')"
    $scope.user = AuthService.user;
  }
]);
