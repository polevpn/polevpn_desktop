; Script generated by the HM NIS Edit Script Wizard.

; HM NIS Edit Wizard helper defines
!define PRODUCT_NAME "PoleVPN"
!define PRODUCT_VERSION "1.1.11"
!define PRODUCT_PUBLISHER "PoleVPN.COM"
!define PRODUCT_WEB_SITE "http://www.polevpn.com"
!define PRODUCT_DIR_REGKEY "Software\Microsoft\Windows\CurrentVersion\App Paths\polevpn.exe"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!define PRODUCT_UNINST_ROOT_KEY "HKLM"

; MUI 1.67 compatible ------
!include "MUI.nsh"

; MUI Settings
!define MUI_ABORTWARNING
!define MUI_ICON "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\polevpn.ico"
!define MUI_UNICON "${NSISDIR}\Contrib\Graphics\Icons\modern-uninstall.ico"

; Language Selection Dialog Settings
!define MUI_LANGDLL_REGISTRY_ROOT "${PRODUCT_UNINST_ROOT_KEY}"
!define MUI_LANGDLL_REGISTRY_KEY "${PRODUCT_UNINST_KEY}"
!define MUI_LANGDLL_REGISTRY_VALUENAME "NSIS:Language"

; Welcome page
!insertmacro MUI_PAGE_WELCOME
; License page
!insertmacro MUI_PAGE_LICENSE "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\license.txt"
; Directory page
!insertmacro MUI_PAGE_DIRECTORY
; Instfiles page
!insertmacro MUI_PAGE_INSTFILES
; Finish page
!define MUI_FINISHPAGE_RUN "$INSTDIR\PoleVPN.exe"
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_INSTFILES

; Language files
!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "SimpChinese"

; MUI end ------

Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "Setup.exe"
InstallDir "$PROGRAMFILES\PoleVPN"
InstallDirRegKey HKLM "${PRODUCT_DIR_REGKEY}" ""
ShowInstDetails show
ShowUnInstDetails show

Function .onInit
  !insertmacro MUI_LANGDLL_DISPLAY
FunctionEnd

Section "MainSection" SEC01

  nsExec::Exec "TaskKill /IM PoleVPN.exe /F"
  nsExec::Exec "TaskKill /IM polevpn_service.exe /F"

  SetOutPath "$INSTDIR"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\init.bat"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\PoleVPN.exe"
  CreateShortCut "$DESKTOP\PoleVPN.lnk" "$INSTDIR\PoleVPN.exe"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\polevpn.ico"
  SetOutPath "$INSTDIR\service"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\service\polevpn_service.exe"
  SetOutPath "$INSTDIR"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\Webview2Loader.dll"
  SetOutPath "$INSTDIR\deps\drivers"
  SetOverwrite try
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\deps\drivers\devcon.exe"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\deps\drivers\OemVista.inf"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\deps\drivers\tap0901.cat"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\deps\drivers\tap0901.sys"
  SetOutPath "$INSTDIR\deps"
  File "C:\Users\Administrator\work\polevpn_desktop\build\win\PoleVPN\deps\MicrosoftEdgeWebview2Setup.exe"
  SetOutPath "$INSTDIR"
  ExecWait "$INSTDIR\init.bat"
SectionEnd

Section -AdditionalIcons
  WriteIniStr "$INSTDIR\${PRODUCT_NAME}.url" "InternetShortcut" "URL" "${PRODUCT_WEB_SITE}"
  CreateDirectory "$SMPROGRAMS\PoleVPN"
  CreateShortCut "$SMPROGRAMS\PoleVPN\Website.lnk" "$INSTDIR\${PRODUCT_NAME}.url"
  CreateShortCut "$SMPROGRAMS\PoleVPN\Uninstall.lnk" "$INSTDIR\uninst.exe"
SectionEnd

Section -Post
  WriteUninstaller "$INSTDIR\uninst.exe"
  WriteRegStr HKLM "${PRODUCT_DIR_REGKEY}" "" "$INSTDIR\PoleVPN.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayName" "$(^Name)"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "UninstallString" "$INSTDIR\uninst.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayIcon" "$INSTDIR\PoleVPN.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "URLInfoAbout" "${PRODUCT_WEB_SITE}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
SectionEnd


Function un.onUninstSuccess
  HideWindow
  MessageBox MB_ICONINFORMATION|MB_OK "$(^Name) �ѳɹ��ش���ļ�����Ƴ���"
FunctionEnd

Function un.onInit
!insertmacro MUI_UNGETLANGUAGE
  MessageBox MB_ICONQUESTION|MB_YESNO|MB_DEFBUTTON2 "��ȷʵҪ��ȫ�Ƴ� $(^Name) ���估���е������" IDYES +2
  Abort
FunctionEnd

Section Uninstall

  nsExec::Exec "TaskKill /IM PoleVPN.exe /F"
  nsExec::Exec "TaskKill /IM polevpn_service.exe /F"
  
  Delete "$INSTDIR\${PRODUCT_NAME}.url"
  Delete "$INSTDIR\uninst.exe"
  Delete "$INSTDIR\Webview2Loader.dll"
  Delete "$INSTDIR\service\polevpn_service.exe"
  Delete "$INSTDIR\polevpn.ico"
  Delete "$INSTDIR\PoleVPN.exe"
  Delete "$INSTDIR\init.bat"
  Delete "$INSTDIR\deps\MicrosoftEdgeWebview2Setup.exe"
  Delete "$INSTDIR\deps\drivers\tap0901.sys"
  Delete "$INSTDIR\deps\drivers\tap0901.cat"
  Delete "$INSTDIR\deps\drivers\OemVista.inf"
  Delete "$INSTDIR\deps\drivers\devcon.exe"

  Delete "$SMPROGRAMS\PoleVPN\Uninstall.lnk"
  Delete "$SMPROGRAMS\PoleVPN\Website.lnk"
  Delete "$DESKTOP\PoleVPN.lnk"

  RMDir "$SMPROGRAMS\PoleVPN"
  RMDir "$INSTDIR\service"
  RMDir "$INSTDIR\deps\drivers"
  RMDir "$INSTDIR\deps"
  RMDir "$INSTDIR"

  DeleteRegKey ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}"
  DeleteRegKey HKLM "${PRODUCT_DIR_REGKEY}"
  SetAutoClose true
SectionEnd