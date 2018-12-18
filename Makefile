#DATE = $$(date +"%Y_%m_%d_%H-%M")
DATE = $$(date +"%Y_%m_%d")

build_linux:
	env GOOS=linux go build -o releases/linux/freeboxctl
build_raspberry:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o releases/raspberry/freeboxctl
build_macosx:
	go build -o releases/macosx/freeboxctl
build_windows:
	env GOOS=windows go build -o releases/windows/freeboxctl.exe

releases: build_linux build_raspberry build_macosx build_windows
	tar -zcvf releases/freeboxctl_$(DATE)_linux.tgz releases/linux
	tar -zcvf releases/freeboxctl_$(DATE)_raspberry.tgz releases/raspberry
	tar -zcvf releases/freeboxctl_$(DATE)_macosx.tgz releases/macosx
	tar -zcvf releases/freeboxctl_$(DATE)_windows.tgz releases/windows

