(function() {
  'use strict';

  angular
    .module('app.cms', ['tm.pagination'])

    .controller('CmsUserInfoCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        var vm = this;

        // 网络查询数据
        vm.tableData = [];
        vm.currentPage = 1;
        vm.getInfo = function(currentPage) {
          if (vm.userName == undefined) {
            vm.userName = "";
          }

          var deferred = $q.defer();

          vm.tableData = [];
          vm.paginationConf.totalItems = undefined;
          vm.currentPage = currentPage;
          vm.paginationConf.currentPage = currentPage;

          // 请求数据
          HttpService.post(
            "/app/cms/user/info", {
              userName: vm.userName,
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

        // 弹框:创建
        vm.cmsCreate = function(data) {
          MsgService.openDialog({
            id: 'CmsUserCreateDialog',
            template: '/tpl/app_cms_user_create.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'CmsUserCreateCtrl',
            data: data,
            callback: function() {
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };
        // 弹框:修改密码
        vm.cmsPwd = function(data) {
          MsgService.openDialog({
            id: 'CmsUserPwdDialog',
            template: '/tpl/app_cms_user_pwd.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'CmsUserPwdCtrl',
            data: data,
            callback: function() {
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };

        // 弹框:分组
        vm.cmsGroup = function(data) {
          MsgService.openDialog({
            id: 'CmsUserGroupDialog',
            template: '/tpl/app_cms_user_group.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'CmsUserGroupCtrl',
            data: data,
            callback: function() {
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };

        // 弹框:删除
        vm.cmsDelete = function(data) {
          MsgService.openDialog({
            id: 'CmsUserDeleteDialog',
            template: '/tpl/app_cms_user_delete.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'CmsUserDeleteCtrl',
            data: data,
            callback: function() {
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };

        // end function
      }
      // end controller
    ])

    .controller('CmsUserCreateCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
      function($filter, $q, $scope, PwdService, HttpService, MsgService) {
        $scope.createPwd = function() {
          $scope.userPwd = PwdService.createPwd();
        }
        $scope.commit = function() {
          if ($scope.userName == undefined || $scope.userName.length == 0) {
            MsgService.openConfirm("请填写帐号");
            return;
          }
          if ($scope.nickName == undefined || $scope.nickName.length == 0) {
            MsgService.openConfirm("请填写昵称");
            return;
          }
          if ($scope.userPwd == undefined || $scope.userPwd.length == 0) {
            MsgService.openConfirm("请填写登录密码");
            return;
          }
          if ($scope.gid == undefined || $scope.gid.length == 0) {
            MsgService.openConfirm("请填写组ID");
            return;
          }
          if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
            MsgService.openConfirm("请填写操作密码");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/user/create", {
              userName: $scope.userName,
              userPwd: sha256($scope.userPwd),
              nickName: $scope.nickName,
              gid: $scope.gid,
              authPwd: sha256($scope.authPwd)
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

    .controller('CmsUserPwdCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
      function($filter, $q, $scope, PwdService, HttpService, MsgService) {
        $scope.userName = $scope.ngDialogData[0];
        $scope.createPwd = function() {
          $scope.userPwd = PwdService.createPwd();
        }
        $scope.commit = function() {
          if ($scope.userPwd == undefined || $scope.userPwd.length == 0) {
            MsgService.openConfirm("请填写用户密码");
            return;
          }
          if ($scope.memo == undefined || $scope.memo.length == 0) {
            MsgService.openConfirm("请填写备注");
            return;
          }
          if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
            MsgService.openConfirm("请填写操作密码");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/user/pwd", {
              userName: $scope.userName,
              userPwd: sha256($scope.userPwd),
              authPwd: sha256($scope.authPwd)
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

    .controller('CmsUserGroupCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
      function($filter, $q, $scope, PwdService, HttpService, MsgService) {
        $scope.userName = $scope.ngDialogData[0];
        $scope.gid = $scope.ngDialogData[2];
        $scope.commit = function() {
          if ($scope.gid == undefined || $scope.gid.length == 0) {
            MsgService.openConfirm("请填写分组ID");
            return;
          }
          if ($scope.memo == undefined || $scope.memo.length == 0) {
            MsgService.openConfirm("请填写备注");
            return;
          }
          if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
            MsgService.openConfirm("请填写操作密码");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/user/group", {
              userName: $scope.userName,
              gid: $scope.gid,
              memo: $scope.memo,
              authPwd: sha256($scope.authPwd)
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

    .controller('CmsUserDeleteCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        $scope.userName = $scope.ngDialogData[0];
        $scope.commit = function() {
          if ($scope.memo == undefined || $scope.memo.length == 0) {
            MsgService.openConfirm("请填写备注");
            return;
          }
          if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
            MsgService.openConfirm("请填写操作密码");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/user/delete", {
              userName: $scope.userName,
              memo: $scope.memo,
              authPwd: sha256($scope.authPwd)
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

    .controller('CmsGroupInfoCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        var vm = this;

        // 网络查询数据
        vm.tableData = [];
        vm.currentPage = 1;
        vm.getInfo = function(currentPage) {
          if (vm.name == undefined) {
            vm.name = "";
          }

          var deferred = $q.defer();

          vm.tableData = [];
          vm.paginationConf.totalItems = undefined;
          vm.currentPage = currentPage;
          vm.paginationConf.currentPage = currentPage;

          // 请求数据
          HttpService.post(
            "/app/cms/group/info", {
              name: vm.name,
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

        // 弹框:创建
        vm.cmsCreate = function(data) {
          MsgService.openDialog({
            id: 'CmsGroupCreateDialog',
            template: '/tpl/app_cms_group_create.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'CmsGroupCreateCtrl',
            data: data,
            callback: function() {
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };
        // 弹框:修改
        vm.cmsSet = function(data) {
          MsgService.openDialog({
            id: 'CmsGroupSetDialog',
            template: '/tpl/app_cms_group_set.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'CmsGroupSetCtrl',
            data: data,
            callback: function() {
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };

        // 弹框:删除
        vm.cmsDelete = function(data) {
          MsgService.openDialog({
            id: 'CmsGroupDeleteDialog',
            template: '/tpl/app_cms_group_delete.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'CmsGroupDeleteCtrl',
            data: data,
            callback: function() {
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };

        // end function
      }
      // end controller
    ])

    .controller('CmsGroupCreateCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
      function($filter, $q, $scope, PwdService, HttpService, MsgService) {
        $scope.commit = function() {
          if ($scope.name == undefined || $scope.name.length == 0) {
            MsgService.openConfirm("请填写分组名称");
            return;
          }
          if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
            MsgService.openConfirm("请填写操作密码");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/group/create", {
              name: $scope.name,
              memo: $scope.memo,
              authPwd: sha256($scope.authPwd)
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

    .controller('CmsGroupSetCtrl', ['$filter', '$q', '$scope', 'PwdService', 'HttpService', 'MsgService',
      function($filter, $q, $scope, PwdService, HttpService, MsgService) {
        $scope.gid = $scope.ngDialogData[0];
        $scope.name = $scope.ngDialogData[1];
        $scope.commit = function() {
          if ($scope.name == undefined || $scope.name.length == 0) {
            MsgService.openConfirm("请填写组名称");
            return;
          }
          if ($scope.memo == undefined || $scope.memo.length == 0) {
            MsgService.openConfirm("请填写备注");
            return;
          }
          if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
            MsgService.openConfirm("请填写操作密码");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/group/set", {
              gid: $scope.gid,
              name: $scope.name,
              memo: $scope.memo,
              authPwd: sha256($scope.authPwd)
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

    .controller('CmsGroupDeleteCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        $scope.gid = $scope.ngDialogData[0];
        $scope.name = $scope.ngDialogData[1];
        $scope.commit = function() {
          if ($scope.memo == undefined || $scope.memo.length == 0) {
            MsgService.openConfirm("请填写备注");
            return;
          }
          if ($scope.authPwd == undefined || $scope.authPwd.length == 0) {
            MsgService.openConfirm("请填写操作密码");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/group/delete", {
              gid: $scope.gid,
              memo: $scope.memo,
              authPwd: sha256($scope.authPwd)
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

    .controller('CmsLogCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        var vm = this;

        // 默认时间控制
        vm.start_time = moment().startOf('day').toDate();
        vm.end_time = moment().add(1, 'days').startOf('day').toDate();


        // 网络查询数据
        vm.tableData = [];
        vm.currentPage = 1;
        vm.getInfo = function(currentPage) {
          if (moment(vm.start_time).isAfter(vm.end_time)) {
            MsgService.openConfirm('起始日期不能晚于截止日期');
            return;
          }
          if (vm.userName == undefined) {
            vm.userName = "";
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
            "/app/cms/log", {
              beginTime: startTime,
              endTime: endTime,
              userName: vm.userName,
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
        // end function
      }
      // end controller
    ])

    .controller('PrivCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        var vm = this;

        // 状态
        vm.myStats = [{
          Name: "已开通",
          Value: '1'
        }, {
          Name: "未开通",
          Value: '2'
        }, {
          Name: "全部",
          Value: ''
        }];
        vm.myStat = vm.myStats[0];


        // 网络查询数据
        vm.tableData = [];
        vm.currentPage = 1;
        vm.getInfo = function(currentPage) {
          if (vm.gid == undefined || vm.gid.length == 0) {
            MsgService.openConfirm("请输入组ID");
            return;
          }

          var deferred = $q.defer();

          vm.tableData = [];
          vm.paginationConf.totalItems = undefined;
          vm.currentPage = currentPage;
          vm.paginationConf.currentPage = currentPage;

          // 请求数据
          HttpService.post(
            "/app/cms/priv", {
              gid: vm.gid,
              menuName: vm.menuName,
              status: vm.myStat.Value,
              pageId: currentPage - 1
            },
            function(data) {
              vm.tmpUserName = vm.gid;
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

        vm.privOn = function(data) {
          // 请求数据
          HttpService.post(
            "/app/cms/priv/on", {
              gid: vm.gid,
              menuId: data[0],
            },
            function(data) {
              MsgService.openConfirm("操作成功", function() {
                vm.getInfo(vm.currentPage)
              });
            });
        };

        vm.privOff = function(data) {
          // 请求数据
          HttpService.post(
            "/app/cms/priv/off", {
              gid: vm.gid,
              menuId: data[0],
            },
            function(data) {
              MsgService.openConfirm("操作成功", function() {
                vm.getInfo(vm.currentPage)
              });
            });
        };

        // 弹框:绑定模板
        vm.privBind = function(data) {
          MsgService.openDialog({
            id: 'PrivBindDialog',
            template: '/tpl/app_cms_priv_bind.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'PrivBindCtrl',
            data: vm.gid,
            callback: function(data) {
              vm.gid = data;
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };
        // end function
      }
      // end controller
    ])

    .controller('PrivBindCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        $scope.gid = $scope.ngDialogData;
        // 创建数据
        $scope.myOptions = [
          ['标准模板']
        ];
        $scope.myOption = $scope.myOptions[0];
        $scope.getTplList = function() {
          // 请求在线数据
          HttpService.post(
            "/app/cms/priv/tpl/list", {},
            function(data) {
              if (data.data && data.data.length > 0) {
                $scope.myOptions = data.data;
                $scope.myOption = $scope.myOptions[0];
              }
            });
        };
        // 拉取模板数据
        $scope.getTplList();

        // 提交数据
        $scope.commit = function() {
          if ($scope.gid == undefined || $scope.gid.length == 0) {
            MsgService.openConfirm("请填写组ID");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/priv/bind", {
              tplName: $scope.myOption[0],
              gid: $scope.gid
            },
            function(data) {
              MsgService.openConfirm("操作成功", function() {
                $scope.closeThisDialog($scope.gid);
              });
            });
          return
        };

        // end function
      }
      // end controller
    ])

    .controller('PrivTplCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        var vm = this;

        // 创建数据
        vm.myOptions = [
          ['标准模板']
        ];
        vm.myOption = vm.myOptions[0];

        // 状态
        vm.myStats = [{
          Name: "已开通",
          Value: '1',
        }, {
          Name: "未开通",
          Value: '2',
        }, {
          Name: "全部",
          Value: '',
        }];
        vm.myStat = vm.myStats[0];

        vm.getTplList = function() {
          // 请求在线数据
          HttpService.post(
            "/app/cms/priv/tpl/list", {},
            function(data) {
              if (data.data && data.data.length > 0) {
                vm.myOptions = data.data;
                vm.myOption = vm.myOptions[0];
              }
            });
        };
        // 更新数据
        vm.getTplList();

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
            "/app/cms/privtpl", {
              tplName: vm.myOption[0],
              menuName: vm.menuName,
              status: vm.myStat.Value,
              pageId: currentPage - 1
            },
            function(data) {
              vm.tmpTplName = vm.myOption[0];
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

        vm.privTplOn = function(data) {
          // 请求数据
          HttpService.post(
            "/app/cms/priv/tpl/on", {
              tplName: vm.myOption[0],
              menuId: data[0],
            },
            function(data) {
              MsgService.openConfirm("操作成功", function() {
                // 刷新页面
                vm.getInfo(vm.currentPage)
              });
            });
        };

        vm.privTplOff = function(data) {
          // 请求数据
          HttpService.post(
            "/app/cms/priv/tpl/off", {
              tplName: vm.myOption[0],
              menuId: data[0]
            },
            function(data) {
              MsgService.openConfirm("操作成功", function() {
                // 刷新页面
                vm.getInfo(vm.currentPage)
              });
            });
        };

        // 弹框:删除模板
        vm.privTplDelete = function(data) {
          MsgService.openDialog({
            id: 'PrivTplDeleteDialog',
            template: '/tpl/app_cms_privtpl_delete.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'PrivTplDeleteCtrl',
            data: vm.myOptions,
            callback: function(data) {
              // 更新数据
              vm.getTplList();
            }
          });
        };

        // 弹框:新建模板
        vm.privTplCreate = function(data) {
          MsgService.openDialog({
            id: 'PrivTplCreateDialog',
            template: '/tpl/app_cms_privtpl_create.html',
            className: 'ngdialog-theme-flat',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            controller: 'PrivTplCreateCtrl',
            data: vm.myOptions,
            callback: function(data) {
              vm.myOptions.push(data);
              vm.myOption = data;
              // 刷新页面
              vm.getInfo(vm.currentPage)
            }
          });
        };
        // end function
      }
      // end controller
    ])

    .controller('PrivTplCreateCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        $scope.myOptions = $scope.ngDialogData;
        $scope.myOption = $scope.ngDialogData[0];

        // 提交数据
        $scope.commit = function() {
          if ($scope.toTplName == undefined || $scope.toTplName.length == 0) {
            MsgService.openConfirm("请填写新模板名称");
            return;
          }

          // 请求数据
          HttpService.post(
            "/app/cms/priv/tpl/new", {
              aTplName: $scope.myOption[0],
              toTplName: $scope.toTplName
            },
            function(data) {
              MsgService.openConfirm("操作成功", function() {
                $scope.closeThisDialog(new Array($scope.toTplName));
              });
            });
          return
        };
        // end function
      }
      // end controller
    ])

    .controller('PrivTplDeleteCtrl', ['$filter', '$q', '$scope', 'HttpService', 'MsgService',
      function($filter, $q, $scope, HttpService, MsgService) {
        $scope.myOptions = $scope.ngDialogData;
        $scope.myOption = $scope.ngDialogData[0];

        // 提交数据
        $scope.commit = function() {
          // 请求数据
          HttpService.post(
            "/app/cms/priv/tpl/delete", {
              tplName: $scope.myOption[0],
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
