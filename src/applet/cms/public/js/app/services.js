'use strict';

/* Services */

// Demonstrate how to register services
angular
  .module('app.services', ['ngDialog'])

.factory('UUIDService', [

    function() {
      return {
        create: function() {
          var s = [];
          var hexDigits = "0123456789abcdef";
          for (var i = 0; i < 36; i++) {
            s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
          }
          s[14] = "4"; // bits 12-15 of the time_hi_and_version field to 0010
          s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1); // bits 6-7 of the clock_seq_hi_and_reserved to 01
          s[8] = s[13] = s[18] = s[23] = "";

          var uuid = s.join("");
          return uuid;
        }
      };
    }
  ])
  .factory('MsgService', ['ngDialog', '$sce',
    function(ngDialog, $sce) {
      // 去除无响应时的重复弹框
      var isOpenDialog = false;
      var isOpenConfirm = false;
      var isOpenMedia = false;
      return {
        openDialog: function(opts) {
          if (isOpenDialog) {
            return
          }
          isOpenDialog = true;
          // 默认不缓存此页面
          if (opts.cache == undefined) {
            opts.cache = false;
            opts.template = opts.template + '?r=' + Math.random();
          }

          ngDialog.open(opts).closePromise.then(function(data) {
            isOpenDialog = false;
            if (opts.callback != undefined && data.value != '$closeButton') {
              opts.callback(data.value);
            }
          });
        },
        openConfirm: function(msg, callback) {
          if (isOpenConfirm) {
            return
          }
          isOpenConfirm = true;

          ngDialog.open({
            template: '/tpl/app_alert.html',
            showClose: false,
            closeByEscape: false,
            closeByDocument: false,
            data: {
              msg: msg
            },
          }).closePromise.then(function(data) {
            isOpenConfirm = false;
            if (callback != undefined) {
              callback(data.value);
            }
          });
          return
        },

        openMedia: function(kind, title, url, callback) {
          if (isOpenMedia) {
            return
          }
          isOpenMedia = true;

          ngDialog.open({
            template: '/tpl/app_media.html',
            showClose: true,
            closeByEscape: false,
            closeByDocument: false,
            data: {
              kind: kind,
              title: title,
              url: $sce.trustAsResourceUrl(url)
            },
          }).closePromise.then(function(data) {
            isOpenMedia = false;
            if (callback != undefined) {
              callback(data.value);
            }
          });
          return
        },
      }
    }
  ])
  .factory('HttpService', ['$state', '$http', 'ngDialog', 'MsgService',
    function($state, $http, ngDialog, MsgService) {
      var posting = false;
      return {
        post: function(url, params, succFn, errFn) {
          if (posting) {
            MsgService.openConfirm("请求中，请稍候...");
            return;
          }
          posting = true
          return $http({
              method: 'POST',
              url: url,
              params: params
            })
            .success(function(data, status, header, config) {
              posting = false;
              succFn(data, status, header, config)
            })
            .error(function(data, status, header, config) {
              posting = false;
              if (errFn) {
                errFn(data, status, header, config)
              } else {
                if (data) {
                  MsgService.openConfirm(data, function() {
                    if (status == 302) {
                      $state.go('access.signin');
                    }
                  });
                } else {
                  MsgService.openConfirm("网络异常:" + status);
                }
              }
            })
        }
      }
    }
  ])
  .factory('PwdService', [
    function() {
      var rand = function(min, max) {
        return min + Math.round(Math.random() * (max - min));
      }

      return {
        createAppPwd: function() {
          // 随机6位数字
          var text = '0123456789';
          var pwd = '';
          for (var i = 0; i < 6; ++i) {
            pwd += text.charAt(Math.floor(Math.random() * 10))
          }
          return pwd;
        },
        createPwd: function(min, max) {
          var text = ['abcdefghijklmnopqrstuvwxyz', 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', '1234567890'];
          if (min == undefined) {
            min = 8
          }
          if (max == undefined) {
            max = 12
          }
          var len = rand(min, max);
          var pw = '';
          for (var i = 0; i < len; ++i) {
            var strpos = rand(0, 2);
            pw += text[strpos].charAt(rand(0, text[strpos].length));
          }
          return pw;
        }
      };
    }
  ])
  .factory('ExportService', ['MsgService',
    function(MsgService) {
      return {
        // 需要对应服务器地址实现导出任务
        // 服务器需接收export与memo两个参数值进行任务判定
        newTask: function(url, params) {
            MsgService.openDialog({
              template: '/tpl/app_export.html',
              className: 'ngdialog-theme-flat',
              showClose: true,
              closeByEscape: false,
              closeByDocument: false,
              controller: ['$scope', 'HttpService', 'MsgService',
                function($scope, HttpService, MsgService) {
                  $scope.commit = function() {
                    if (params == undefined) {
                      params = {}
                    }
                    // 导出数据专用
                    // export=1为标识，memo为备注
                    params.export = 1; // 专用的参数
                    params.memo = $scope.memo; // 专用的数据
                    params.authPwd = $scope.authPwd;
                    // 请求导出
                    HttpService.post(url, params, function(data) {
                      MsgService.openConfirm('请求成功，请在数据导出的"任务列表"中查看状态', function() {
                        $scope.closeThisDialog();
                      })
                    })
                  }
                }
              ]
            });
            return
          } // end newTask
      } // end return
    } // end service function
  ]) // end service

.factory('AuthService', ['$localStorage',
    function($localStorage) {
      var user = {};
      if (angular.isDefined($localStorage.user)) {
        user = $localStorage.user;
      }
      // 权限校验
      // 使用流程
      // 在服务端配置权限key
      // 登录时下发key至本地，更改权限时，需要重新登录生效
      // 在controller中加载此算法
      // 在页上使用ng-show="Auth('keyname')控制页面显示
      // 请求服器时，服务器再次校验key值
      // 
      // key生成规则
      // 请求路径的组合为key，服器器收到请求时，会以请求路径拼接出key值，所以权限的实际粒度为每一个请求路径。
      return {
        user: user,
        setAuth: function(u) {
          $localStorage.user = u;
          user.username = u.username;
          user.nickname = u.nickname;
          user.logo = u.logo;
          user.priv = u.priv;
        },
        checksum: function(key) {
          if (!user.priv || !key) {
            return false;
          }

          // 以下算法请查阅服务器相关说明
          var cp = user.priv;
          var cpLen = user.priv.length;
          var node = key.split(".")
          var nodeLen = node.length;
          var pass = false;
          for (var i = 0; i < cpLen; i++) {
            if (cp[i].length < nodeLen) {
              continue;
            }
            pass = true;
            for (var j = nodeLen - 1; j > -1; j--) {
              if (!(cp[i][j] == node[j] || cp[i][j] == '*')) {
                pass = false;
                break;
              }
            }
            if (pass) {
              return true;
            }
          }
          return false;
        }
      }
    }
  ])
  .factory('UploadService', ['$http', '$upload', 'MsgService', 'UUIDService',
    function($http, $upload, MsgService, UUIDService) {
      return {
        // 上传成功会回调地址
        upload: function($files, group, callback) {
          if (!$files) {
            MsgService.openConfirm('控件错误');
            return
          }
          if ($files.length > 1) {
            MsgService.openConfirm('仅能上传单个文件');
            return
          }
          if ($files.length == 0) {
            MsgService.openConfirm('未选择文件');
            return;
          }

          // 给用户提示文件大小
          if ($files.size > 100 * 1024) {
            MsgService.openConfirm('文件大小超过了100KB, 建议对图片进行优化');
          }

          // 获取上传的token
          $http({
            method: 'POST',
            url: "/api/uploader"
          }).error(function(err) {
            MsgService.openConfirm(err);
            return;
          }).success(function(data) {
            // 上传数据 
            var key = [data.bucket, "/", group, "-", UUIDService.create()].join("");
            $upload.upload({
              url: 'https://up-z2.qiniup.com/',
              data: {
                key: key,
                token: data.token
              },
              file: $files[0]
            }).progress(function(evt) {
              // 控制台显示进茺
              console.log($files[0].name + ' upload: ' + parseInt(100.0 * evt.loaded / evt.total));
            }).success(function() {
              // 上传成功，回调下载地址给调用者
              callback(data.domain + "/" + key);
            }).error(function(err) {
              MsgService.openConfirm(err);
            });

          });
        }
      };
    }
  ]);
