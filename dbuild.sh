# /bin/sh

# # run docker mode

pwd_dir=`pwd`

go get -v github.com/gwaylib/goget||exit 1
mkdir -p $PRJ_ROOT/bin/docker||exit 1

# build bin data
export CGO_ENABLED=0 GOOS=linux GOARCH=amd64

goget -v -d github.com/gwaycc/supd||exit 1
cd $GOLIB/src/github.com/gwaycc/supd||exit 1
go build||exit 1
mv ./supd $PRJ_ROOT/bin/docker||exit 1
cd $pwd_dir

goget -v -d github.com/davidpersson/bsa||exit 1
cd $GOLIB/src/github.com/davidpersson/bsa||exit 1
go build||exit 1
mv ./bsa $PRJ_ROOT/bin/docker||exit 1
cd $pwd_dir

sup publish all

echo "# Building Dockerfile"
# remove old images
sudo docker rmi -f $PRJ_NAME||exit 1
# build images
sudo docker build -t $PRJ_NAME .||exit 1
# rm tmp data
# rm app

# show images build result
sudo docker images $PRJ_NAME||exit 1

