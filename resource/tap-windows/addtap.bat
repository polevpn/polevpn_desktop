rem Add a new TAP virtual ethernet adapter
cd /D %~dp0
tapinstall.exe install driver\OemVista.inf tap0901
pause
