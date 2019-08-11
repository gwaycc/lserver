(function() {
  'use strict';

  angular
    .module('app.log', ['tm.pagination'])

  .controller('LogInfoCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService', 'AuthService',
    function($filter, $q, $scope, HttpService, MsgService, AuthService) {
      $scope.Auth = AuthService.checksum;
      var vm = this;

      // 默认时间控制
      vm.start_time = moment().startOf('day').toDate();
      vm.end_time = moment().add(1, 'days').startOf('day').toDate();
      vm.level = 2;

      $scope.$watch("vm.md5", function() {
        if (vm.md5 != undefined && vm.md5.length > 0) {
          vm.platform = "";
          vm.level = "";
          vm.logger = "";
          vm.msg = "";
        }
      });
      $scope.$watch("vm.platform", function() {
        if (vm.platform != undefined && vm.platform.length > 0) {
          vm.md5 = "";
        }
      });
      $scope.$watch("vm.level", function() {
        if (vm.level != undefined && vm.level.length > 0) {
          vm.md5 = "";
        }
      });
      $scope.$watch("vm.logger", function() {
        if (vm.logger != undefined && vm.logger.length > 0) {
          vm.md5 = "";
        }
      });
      $scope.$watch("vm.msg", function() {
        if (vm.msg != undefined && vm.msg.length > 0) {
          vm.md5 = "";
        }
      });

      // 网络查询数据
      vm.tableData = [];
      vm.currentPage = 1;
      vm.getInfo = function(currentPage) {
        if (moment(vm.start_time).isAfter(vm.end_time)) {
          MsgService.openConfirm('起始日期不能晚于截止日期');
          return;
        }

        var deferred = $q.defer();
        var startTime = $filter('date')(vm.start_time, 'yyyy-MM-dd HH:mm:ss'),
          endTime = $filter('date')(vm.end_time, 'yyyy-MM-dd HH:mm:ss');

        vm.tableData = [];
        vm.paginationConf.totalItems = undefined;
        vm.currentPage = currentPage;
        vm.paginationConf.currentPage = currentPage;

        // 请求数据
        HttpService.post(
          "/app/log/info", {
            beginTime: startTime,
            endTime: endTime,
            md5: vm.md5,
            platform: vm.platform,
            level: vm.level,
            logger: vm.logger,
            msg: vm.msg,
            pageId: currentPage - 1
          },
          function(data) {
            vm.tableName = data.names;
            vm.tableData = data.data;
            vm.paginationConf.totalItems = data.total;
            deferred.resolve(data.total);

            if (!data.data || data.data.length == 0) {
              MsgService.openConfirm("无可用数据");
            }
          });
        return deferred.promise;
      };

      // 分页基本参数
      vm.paginationConf = {
        currentPage: 1,
        totalItems: 0,
        itemsPerPage: 10,
        pagesLength: 15,
        onChange: function() {
          // currentPage==0为假
          if (vm.paginationConf.totalItems && vm.currentPage != vm.paginationConf.currentPage) {
            vm.getInfo(vm.paginationConf.currentPage);
          }
        }
      };

      // 日志详情
      vm.detail = function(data) {
        MsgService.openDialog({
          id: 'LogDetailDialog',
          template: '/tpl/app_log_detail.html',
          className: 'ngdialog-theme-flat',
          showClose: true,
          closeByEscape: false,
          closeByDocument: false,
          controller: 'LogDetailCtrl',
          data: data,
          callback: function() {
            // 什么者不需要操作
          }
        });
      };
    }
  ])

  .controller('LogDetailCtrl', ['$scope', 'MsgService',
    function($scope, MsgService) {
      // 不需做什么
      $scope.myFunction = function() {
        var result = JSON.stringify($scope.ngDialogData, null, 4);
        return result;
      }　
    }
  ])

  .controller('LogAlertorCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService', 'AuthService',
    function($filter, $q, $scope, HttpService, MsgService, AuthService) {
      $scope.Auth = AuthService.checksum;
      var vm = this;

      // 网络查询数据
      vm.tableData = [];
      vm.currentPage = 1;
      vm.getInfo = function(currentPage) {
        var deferred = $q.defer();

        vm.tableData = [];
        vm.paginationConf.totalItems = undefined;
        vm.currentPage = currentPage;
        vm.paginationConf.currentPage = currentPage;

        // 请求数据
        HttpService.post(
          "/app/log/alertor", {
            pageId: currentPage - 1
          },
          function(data) {
            vm.tableName = data.names;
            vm.tableData = data.data;
            vm.paginationConf.totalItems = data.total;
            deferred.resolve(data.total);

            if (!data.data || data.data.length == 0) {
              MsgService.openConfirm("无可用数据");
            }
          });
        return deferred.promise;
      };

      // 分页基本参数
      vm.paginationConf = {
        currentPage: 1,
        totalItems: 0,
        itemsPerPage: 10,
        pagesLength: 15,
        onChange: function() {
          // currentPage==0为假
          if (vm.paginationConf.totalItems && vm.currentPage != vm.paginationConf.currentPage) {
            vm.getInfo(vm.paginationConf.currentPage);
          }
        }
      };

      // 首次使用时初始化数据
      vm.getInfo(0);

      // 增加数据
      vm.add = function(data) {
        MsgService.openDialog({
          id: 'LogAlertorAddDialog',
          template: '/tpl/app_log_alertor_add.html',
          className: 'ngdialog-theme-flat',
          showClose: true,
          closeByEscape: false,
          closeByDocument: false,
          controller: 'LogAlertorAddCtrl',
          data: data,
          callback: function() {
            // 刷新页面
            vm.getInfo(vm.currentPage)
          }
        });
      };

      // 修改数据
      vm.set = function(data) {
        MsgService.openDialog({
          id: 'LogAlertorSetDialog',
          template: '/tpl/app_log_alertor_set.html',
          className: 'ngdialog-theme-flat',
          showClose: true,
          closeByEscape: false,
          closeByDocument: false,
          controller: 'LogAlertorSetCtrl',
          data: data,
          callback: function() {
            // 刷新页面
            vm.getInfo(vm.currentPage)
          }
        });
      };

      // 删除数据
      vm.del = function(data) {
        MsgService.openDialog({
          id: 'LogAlertorDelDialog',
          template: '/tpl/app_log_alertor_del.html',
          className: 'ngdialog-theme-flat',
          showClose: true,
          closeByEscape: false,
          closeByDocument: false,
          controller: 'LogAlertorDelCtrl',
          data: data,
          callback: function() {
            // 刷新页面
            vm.getInfo(vm.currentPage)
          }
        });
      };
    }
  ])

  .controller('LogAlertorAddCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
    function($filter, $q, $scope, PwdService, HttpService, MsgService) {
      $scope.commit = function() {
        if ($scope.nickName == undefined || $scope.nickName.length == 0) {
          MsgService.openConfirm("请填写昵称");
          return;
        }
        if ($scope.mobile == undefined || $scope.mobile.length == 0) {
          MsgService.openConfirm("请填写手机号");
          return;
        }
        if ($scope.email == undefined || $scope.email.length == 0) {
          MsgService.openConfirm("请填写邮件");
          return;
        }
        if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
          MsgService.openConfirm("请填写操作密码");
          return;
        }

        // 请求数据
        HttpService.post(
          "/app/log/alertor/add", {
            nickName: $scope.nickName,
            mobile: $scope.mobile,
            email: $scope.email,
            authPwd: $scope.authPwd
          },
          function(data) {
            MsgService.openConfirm("操作成功", function() {
              $scope.closeThisDialog();
            });
          });
        return
      };
      // end function
    }
    // end controller
  ])

  .controller('LogAlertorSetCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
    function($filter, $q, $scope, PwdService, HttpService, MsgService) {
      $scope.nickName = $scope.ngDialogData[0];
      $scope.mobile = $scope.ngDialogData[1];
      $scope.email = $scope.ngDialogData[2];
      $scope.commit = function() {
        if ($scope.nickName == undefined || $scope.nickName.length == 0) {
          MsgService.openConfirm("请填写昵称");
          return;
        }
        if ($scope.mobile == undefined || $scope.mobile.length == 0) {
          MsgService.openConfirm("请填写手机号");
          return;
        }
        if ($scope.email == undefined || $scope.email.length == 0) {
          MsgService.openConfirm("请填写邮件");
          return;
        }
        if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
          MsgService.openConfirm("请填写操作密码");
          return;
        }

        // 请求数据
        HttpService.post(
          "/app/log/alertor/set", {
            nickName: $scope.nickName,
            mobile: $scope.mobile,
            email: $scope.email,
            authPwd: $scope.authPwd
          },
          function(data) {
            MsgService.openConfirm("操作成功", function() {
              $scope.closeThisDialog();
            });
          });
        return
      };
      // end function
    }
    // end controller
  ])

  .controller('LogAlertorDelCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
    function($filter, $q, $scope, PwdService, HttpService, MsgService) {
      $scope.nickName = $scope.ngDialogData[0];
      $scope.commit = function() {
        if ($scope.nickName == undefined || $scope.nickName.length == 0) {
          MsgService.openConfirm("请填写昵称");
          return;
        }
        if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
          MsgService.openConfirm("请填写操作密码");
          return;
        }

        // 请求数据
        HttpService.post(
          "/app/log/alertor/del", {
            nickName: $scope.nickName,
            authPwd: $scope.authPwd
          },
          function(data) {
            MsgService.openConfirm("操作成功", function() {
              $scope.closeThisDialog();
            });
          });
        return
      };
      // end function
    }
    // end controller
  ])

  .controller('LogMailCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService', 'AuthService',
    function($filter, $q, $scope, HttpService, MsgService, AuthService) {
      $scope.Auth = AuthService.checksum;
      var vm = this;

      // 网络查询数据
      vm.tableData = [];
      vm.currentPage = 1;
      vm.getInfo = function(currentPage) {
        var deferred = $q.defer();

        vm.tableData = [];
        vm.paginationConf.totalItems = undefined;
        vm.currentPage = currentPage;
        vm.paginationConf.currentPage = currentPage;

        // 请求数据
        HttpService.post(
          "/app/log/mail", {
            pageId: currentPage - 1
          },
          function(data) {
            vm.tableName = data.names;
            vm.tableData = data.data;
            vm.paginationConf.totalItems = data.total;
            deferred.resolve(data.total);

            if (!data.data || data.data.length == 0) {
              MsgService.openConfirm("无可用数据");
            }
          });
        return deferred.promise;
      };

      // 分页基本参数
      vm.paginationConf = {
        currentPage: 1,
        totalItems: 0,
        itemsPerPage: 10,
        pagesLength: 15,
        onChange: function() {
          // currentPage==0为假
          if (vm.paginationConf.totalItems && vm.currentPage != vm.paginationConf.currentPage) {
            vm.getInfo(vm.paginationConf.currentPage);
          }
        }
      };
      vm.getInfo(0);

      // 修改数据
      vm.set = function(data) {
        MsgService.openDialog({
          id: 'LogMailSetDialog',
          template: '/tpl/app_log_mail_set.html',
          className: 'ngdialog-theme-flat',
          showClose: true,
          closeByEscape: false,
          closeByDocument: false,
          controller: 'LogMailSetCtrl',
          data: data,
          callback: function() {
            // 刷新页面
            vm.getInfo(vm.currentPage)
          }
        });
      };
    }
  ])

  .controller('LogMailSetCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
    function($filter, $q, $scope, PwdService, HttpService, MsgService) {
      $scope.smtpHost = $scope.ngDialogData[0];
      $scope.stmpPort = $scope.ngDialogData[1];
      if ($scope.smtpPort == undefined || $scope.smtpPort.length == 0) {
        $scope.smtpPort = "25"
      }
      $scope.mAuthName = $scope.ngDialogData[2];
      $scope.commit = function() {
        if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
          MsgService.openConfirm("请填写操作密码");
          return;
        }

        // 请求数据
        HttpService.post(
          "/app/log/mail/set", {
            smtpHost: $scope.smtpHost,
            smtpPort: $scope.smtpPort,
            mAuthName: $scope.mAuthName,
            mAuthPwd: $scope.mAuthPwd,
            authPwd: $scope.authPwd
          },
          function(data) {
            MsgService.openConfirm("操作成功", function() {
              $scope.closeThisDialog();
            });
          });
        return
      };
      // end function
    }
    // end controller
  ])

})();
