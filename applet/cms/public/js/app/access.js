(function() {
  'use strict';

  angular
    .module('app.access', [])

  .controller('SigninCtrl', ['$state', '$localStorage', 'AuthService', 'HttpService', 'MsgService',
    function($state, $localStorage, AuthService, HttpService, MsgService) {
      var vm = this;

      // 获取验证码
      vm.vcodeId = "";
      vm.vcodeImg = "";
      vm.reloadVcode = function() {
        HttpService.post("/api/vcode", {
            "id": vm.vcodeId
          },
          function(data) {
            vm.vcodeId = data.id;
            vm.vcodeImg = data.img;
          });
      };
      // 获取第一张验证码
      vm.reloadVcode();


      // 执行登录
      vm.login = function() {
        if (vm.acc == undefined || vm.acc.length == 0) {
          MsgService.openConfirm("输入账户");
          return;
        };
        if (vm.pwd == undefined || vm.pwd.length == 0) {
          MsgService.openConfirm("输入密码");
          return;
        };
        if (vm.vcodeData == undefined || vm.vcodeData.length == 0) {
          MsgService.openConfirm("输入验证码");
          return;
        };


        HttpService.post("/access/signin", {
            "acc": vm.acc,
            "pwd": sha256(vm.pwd),
            "vcodeId": vm.vcodeId,
            "vcodeData": vm.vcodeData
          },
          // 正确的响应
          function(data) {
            // 更新登录信息
            AuthService.setAuth(data);
            $state.go('app.dashboard');
          },

          // 错误的响应
          function(data, status) {
            if (data) {
              MsgService.openConfirm(data, function() {
                vm.reloadVcode();
              });
            } else {
              MsgService.openConfirm("网络异常:" + status);
            }
          });
      };
    }
  ])

  .controller('SignoutCtrl', ['$window', 'HttpService', 'MsgService',
    function($window, HttpService, MsgService) {
      var vm = this;
      vm.logout = function() {
        HttpService.post("/access/signup", {},
          function(data) {
            $window.location.href = '/access/signin';
          })
      };
      // 弹框:修改密码
      vm.setPwd = function(data) {
        MsgService.openDialog({
          id: 'setPwdDialog',
          template: '/tpl/access_pwd.html',
          className: 'ngdialog-theme-flat',
          showClose: true,
          closeByEscape: false,
          closeByDocument: false,
          controller: 'setPwdCtrl',
          data: data,
          callback: function() {
            // 关闭后不需要做什么
          }
        });
      };
    }
  ])

  .controller('setPwdCtrl', ['$scope', 'MsgService', 'HttpService',
    function($scope, MsgService, HttpService) {
      $scope.commit = function() {
        if ($scope.oldPwd == undefined || $scope.oldPwd.length == 0) {
          MsgService.openConfirm("请输入原密码");
          return;
        }
        if ($scope.newPwd == undefined || $scope.newPwd.length == 0) {
          MsgService.openConfirm("请输入新密码");
          return;
        }
        if ($scope.rePwd == undefined || $scope.rePwd.length == 0) {
          MsgService.openConfirm("请重复新密码");
          return;
        }
        if ($scope.newPwd != $scope.rePwd) {
          MsgService.openConfirm("新密码不一致");
          return;
        }
        HttpService.post("/access/pwd", {
            oldPwd: sha256($scope.oldPwd),
            newPwd: sha256($scope.newPwd)
          },
          function(data) {
            MsgService.openConfirm("修改成功", function() {
              $scope.closeThisDialog();
            });
          })
      };
    }
  ]);
})();
