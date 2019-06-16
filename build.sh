echo 'Going to game directory'
cd $GOPATH/src/github.com/troyspencer/launch-pixelgl
echo 'Checking dependencies'
go get
echo 'Building...'
go build
echo 'Built!'
echo 'Copying build for distribution'
cp ./launch.exe $GOPATH/src/troy/gcloud/grpc_playground/server/game
echo 'Copied!'