@echo off
cd ../abugo
git pull
cd ../gobase
xcopy "../abugo/abugo.go" "adminapi/abugo" /D /E /I /F /Y
xcopy "../abugo/abugo.go" "clientapi/abugo" /D /E /I /F /Y
xcopy "../abugo/abugo.go" "thirdapi/abugo" /D /E /I /F /Y
pause