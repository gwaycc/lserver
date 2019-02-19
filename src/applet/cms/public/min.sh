#!/bin/sh

# npm安装(已安装忽略) http://www.iteblog.com/archives/1313
# wget https://npmjs.org/install.sh --no-check-certificate
# chmod 777 install.sh
# sudo ./install.sh
# npm -v
#
# 工具安装
# sudo npm install clean-css-cli -g
# sudo npm install uglify-es -g # 或者 sudo npm install uglify-js -g


#################################3
# cleancss
cat /dev/null > app.min.css|exit 1

cat css/bootstrap.css >> app.min.css|exit 1
cat css/animate.css >> app.min.css|exit 1
cat css/font-awesome.min.css >> app.min.css|exit 1
cat css/simple-line-icons.css >> app.min.css|exit 1
cat css/font.css >> app.min.css|exit 1
cat css/app.css >> app.min.css|exit 1
cat css/ngDialog.css >> app.min.css|exit 1
cat css/ngDialog-flat.css >> app.min.css|exit 1
cat css/ngDialog-theme-default.css >> app.min.css|exit 1
cat css/datepicker.css >> app.min.css|exit 1

cleancss -o app.min.css app.min.css|exit 1
mv app.min.css css/|exit 1

#################################3
# uglilyjs
cat /dev/null > app.min.js|exit 1
cat vendor/jquery/jquery.min.js >> app.min.js|exit 1
cat vendor/angular/angular-file-upload-shim.min.js >> app.min.js|exit 1
cat vendor/angular/angular.js >> app.min.js|exit 1
cat vendor/angular/angular-file-upload.min.js >> app.min.js|exit 1
cat vendor/angular/angular-cookies.min.js >> app.min.js|exit 1
cat vendor/angular/angular-animate.min.js >> app.min.js|exit 1
cat vendor/angular/angular-ui-router.js >> app.min.js|exit 1
cat vendor/angular/angular-translate.js >> app.min.js|exit 1
cat vendor/angular/ngStorage.min.js >> app.min.js|exit 1
cat vendor/angular/ocLazyLoad.min.js >> app.min.js|exit 1
cat vendor/angular/ui-load.js >> app.min.js|exit 1
cat vendor/angular/ui-jq.js >> app.min.js|exit 1
cat vendor/angular/ui-validate.js >> app.min.js|exit 1
cat vendor/angular/ui-bootstrap-tpls.min.js >> app.min.js|exit 1
cat vendor/angular/textAngular.min.js >> app.min.js|exit 1
cat vendor/angular/textAngular-rangy.min.js >> app.min.js|exit 1
cat vendor/angular/textAngular-sanitize.min.js >> app.min.js|exit 1
cat vendor/angular/angular-spinner.min.js >> app.min.js|exit 1
cat vendor/angular/spin.min.js >> app.min.js|exit 1
cat js/libs/moment.js >> app.min.js|exit 1
cat js/libs/area.js >> app.min.js|exit 1
cat js/libs/ngDialog.js >> app.min.js|exit 1
cat js/libs/tm.pagination.js >> app.min.js|exit 1
cat js/libs/ichart.1.2.min.js >> app.min.js|exit 1
cat js/libs/bootstrap.min.js >> app.min.js|exit 1
cat js/libs/angular-datepicker.js >> app.min.js|exit 1
cat js/libs/angular-validation.js >> app.min.js|exit 1
cat js/libs/angular-validation-rule.js >> app.min.js|exit 1
cat js/app/services.js >> app.min.js|exit 1
cat js/app/filters.js >> app.min.js|exit 1
cat js/app/directives.js >> app.min.js|exit 1
cat js/app/main.js >> app.min.js|exit 1
cat js/app/access.js >> app.min.js|exit 1
cat js/app/cms.js >> app.min.js|exit 1
cat js/app/log.js >> app.min.js|exit 1
cat js/app.js >> app.min.js|exit 1

uglifyjs -o app.min.js app.min.js|exit 1
mv app.min.js js/|exit 1
