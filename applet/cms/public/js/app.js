'use strict';

// 获取页面上服务器下发的数据
var appName = document.getElementById("app.name").value; // redefine in index.html
var appVersion = document.getElementById("app.version").value; // redefine in index.html

// 重构数组方法以便可以删除
Array.prototype.remove = function(dx) {
  if (isNaN(dx) || dx > this.length) {
    return false;
  }
  for (var i = 0, n = 0; i < this.length; i++) {
    if (this[i] != this[dx]) {
      this[n++] = this[i]
    }
  }
  this.length -= 1　
}


// Declare app level module which depends on filters, and services
var app = angular
  .module('app', [
    'ngAnimate',
    'ngCookies',
    'ngStorage',
    'ui.router',
    'ui.bootstrap',
    'ui.load',
    'ui.jq',
    'ui.validate',
    'oc.lazyLoad',
    'pascalprecht.translate',

    'validation',
    'validation.rule',
    'textAngular',
    'datePicker',

    'app.filters',
    'app.directives',
    'app.services',
    'app.main',
    'app.access',
    'app.cms',
    'app.log'
  ])
  .constant('cms', {
    name: appName,
    version: appVersion,
    // for chart colors
    color: {
      primary: '#7266ba',
      info: '#23b7e5',
      success: '#27c24c',
      warning: '#fad733',
      danger: '#f05050',
      light: '#e8eff0',
      dark: '#3a3f51',
      black: '#1c2b36'
    },
    settings: {
      themeID: 1,
      navbarHeaderColor: 'bg-black',
      navbarCollapseColor: 'bg-white-only',
      asideColor: 'bg-black',
      headerFixed: true,
      asideFixed: false,
      asideFolded: false,
      asideDock: false,
      container: false
    }
  })
  .run([
    '$rootScope',
    '$state',
    '$stateParams',
    function($rootScope, $state, $stateParams) {
      $rootScope.$state = $state;
      $rootScope.$stateParams = $stateParams;
    }
  ]);

// config
app.config([
    '$stateProvider',
    '$urlRouterProvider',
    '$controllerProvider',
    '$compileProvider',
    '$filterProvider',
    '$provide',
    '$logProvider',
    'cms',
    function($stateProvider, $urlRouterProvider, $controllerProvider, $compileProvider, $filterProvider, $provide, $logProvider, cms) {
      $logProvider.debugEnabled(false);

      // controller, directive and service
      app.controller = $controllerProvider.register;
      app.directive = $compileProvider.directive;
      app.filter = $filterProvider.register;
      app.factory = $provide.factory;
      app.service = $provide.service;
      app.constant = $provide.constant;

      $urlRouterProvider.otherwise('/access/signin');
      $stateProvider
        .state('access', {
          url: '/access',
          template: '<div ui-view class="fade-in-right-big smooth"></div>'
        })
        .state('access.signin', {
          url: '/signin',
          templateUrl: '/tpl/access_signin.html' + '?v=' + cms.version
        })

        .state('app', {
          abstract: true,
          url: '/app',
          templateUrl: '/tpl/app.html' + '?v=' + cms.version
        })
        .state('app.dashboard', {
          url: '/dashboard',
          templateUrl: '/tpl/app_dashboard.html' + '?v=' + cms.version
        })

        .state('app.cms', {
          abstract: true,
          url: '/cms',
          template: '<div ui-view class="fade-in-right-big"></div>'
        })

        .state('app.cms.user', {
          abstract: true,
          url: '/user',
          template: '<div ui-view class="fade-in-right-big"></div>'
        })
        .state('app.cms.user.info', {
          url: '/info',
          templateUrl: '/tpl/app_cms_user_info.html' + '?v=' + cms.version,
          controller: 'CmsUserInfoCtrl',
          controllerAs: 'vm'
        })

        .state('app.cms.group', {
          abstract: true,
          url: '/group',
          template: '<div ui-view class="fade-in-right-big"></div>'
        })
        .state('app.cms.group.info', {
          url: '/info',
          templateUrl: '/tpl/app_cms_group_info.html' + '?v=' + cms.version,
          controller: 'CmsGroupInfoCtrl',
          controllerAs: 'vm'
        })

        .state('app.cms.priv', {
          url: '/priv',
          templateUrl: '/tpl/app_cms_priv.html' + '?v=' + cms.version,
          controller: 'PrivCtrl',
          controllerAs: 'vm'
        })
        .state('app.cms.privtpl', {
          url: '/privtpl',
          templateUrl: '/tpl/app_cms_privtpl.html' + '?v=' + cms.version,
          controller: 'PrivTplCtrl',
          controllerAs: 'vm'
        })
        .state('app.cms.log', {
          url: '/log',
          templateUrl: '/tpl/app_cms_log.html' + '?v=' + cms.version,
          controller: 'CmsLogCtrl',
          controllerAs: 'vm'
        })

        .state('app.log', {
          abstract: true,
          url: '/log',
          template: '<div ui-view class="fade-in-right-big"></div>'
        })
        .state('app.log.info', {
          url: '/info',
          templateUrl: '/tpl/app_log_info.html' + '?v=' + cms.version,
          controller: 'LogInfoCtrl',
          controllerAs: 'vm'
        })
        .state('app.log.alertor', {
          url: '/alertor',
          templateUrl: '/tpl/app_log_alertor.html' + '?v=' + cms.version,
          controller: 'LogAlertorCtrl',
          controllerAs: 'vm'
        })
        .state('app.log.mail', {
          url: '/mail',
          templateUrl: '/tpl/app_log_mail.html' + '?v=' + cms.version,
          controller: 'LogMailCtrl',
          controllerAs: 'vm'
        })

    }
  ])

  .config([
    '$locationProvider',
    function($locationProvider) {
      // mode
      $locationProvider.html5Mode(true);
    }
  ])

  .config(['$validationProvider',
    function($validationProvider) {
      $validationProvider.showSuccessMessage = false;
      $validationProvider
        .setExpression({
          username: function(value, scope, element, attrs) {
            if (scope.$parent.edit_user && value.length === 0) {
              return true;
            }
            var result = /^[a-zA-Z]{1}([a-zA-Z0-9]|[._]){4,19}$/.test(value);
            return result;
          },
          password: function(value, scope, element, attrs) {
            if (scope.$parent.edit_user && value.length === 0) {
              return true;
            }
            // 需要8-12位密码
            if (value.length < 8 || value.length > 20) {
              return false;
            }
            // 需要含有大小写数字等
            var regex = new RegExp(/^(?![^a-z]+$)(?![^A-Z]+$)(?!\D+$)/);
            return regex.test(value);
          },
          money: /^(([1-9]\d{0,9})|0)(\.\d{1,2})?$/,
          fee: /^(([1-9]\d{0,9}|0))$/,
          descr: function(value) {
            if (value.length > 200) {
              return false;
            }
            return true;
          },
          phonenumber: /^(13[0-9]|14[157]|15[0-35-9]|17[0-9]|18[0-9])\d{8}$/
        })
        .setDefaultMsg({
          username: {
            error: '只能输入5-20个以字母开头、可带数字、“_”、“.”的字串',
          },
          password: {
            error: '8~20位，需包含大小写字母和数字'
          },
          money: {
            error: '格式有误,例如0.00'
          },
          fee: {
            error: '金额格式错误'
          },
          descr: {
            error: '描述不能超过200字'
          },
          phonenumber: {
            error: '号码格式错误'
          }
        });
    }
  ]);
