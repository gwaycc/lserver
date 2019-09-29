# /bin/sh

# # run docker mode

pwd_dir=`pwd`

# build bin data
export CGO_ENABLED=0 GOOS=linux GOARCH=amd64

sup build -o $GOBIN/supd github.com/gwaycc/supd/cmd/supd||exit 1
sup build -o $GOBIN/bsa github.com/davidpersson/bsa||exit 1

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

