echo "install tuntap driver"

deps\drivers\devcon.exe install deps\drivers\OemVista.inf tap0901

reg query  "HKLM\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"

if %errorlevel% neq 0 goto NO
if %errorlevel% equ 0 goto YES

:NO

reg query  "HKCU\Software\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"

if %errorlevel% neq 0 goto INSTALL
if %errorlevel% equ 0 goto YES

:YES
echo Webview2Runtime installed
goto END

:INSTALL
echo install webview2runtime
deps\MicrosoftEdgeWebview2Setup.exe
goto END

:END



