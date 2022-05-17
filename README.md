sudo spctl --master-disable
sudo xattr -r -d com.apple.quarantine /xxx/xxx.app
codesign --force --deep --sign - /xxx/xxx.app

