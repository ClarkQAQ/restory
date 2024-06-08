<p align="center">
  <h3 align="center">Restory</h3>
  <p align="center">
    Re (store + (hi) story) VSCode 历史记录文件恢复工具
    <br />
    VSCode History Recovery Tool
  </p>
</p>

[English README](README.en.md)

> 你在使用 VSCode，你的文件夹不小心删掉了? Restory 就能帮你恢复!

### 原理

这是一个很简单的工具. Restory 通过扫描 VSCode 的文件修改历史文件获取信息, 并在你指定的位置重建文件夹目录结构和文件内容。
很不幸的是, 你没有经过 VSCode 无法使用 Restore 恢复, 因为 VSCode 并没有产生对应的修改历史文件。 


### 使用

前往 Release 下载您平台对应的 `restory` 的二进制文件然后运行即可.

例如在 Linux 平台:

```bash
$ curl -L -o ./restory_linux_amd64 https://github.com/ClarkQAQ/restory/releases/latest/download/restory_linux_amd64
$ chmod +x ./restory_linux_amd64
$ ./restory_linux_amd64 ~/.config/Code/User/History /tmp/code-history
```

VSCode 历史文件夹目录在 Linux 平台通常是 `~/.config/Code/User/History`, Windows 平台则是 `C:\Users\[你的用户名]\AppData\Local\Code\User\History`.

[Windows 下载链接](https://github.com/ClarkQAQ/restory/releases/latest/download/restory_windows_amd64.exe)
[Linux 下载链接](https://github.com/ClarkQAQ/restory/releases/latest/download/restory_linux_amd64)
[MacOS 下载链接](https://github.com/ClarkQAQ/restory/releases/latest/download/restory_darwin_amd64)

### 最后

以后请多备份自己的重要文件! 及时上传到 Git! 
或许我的这个工具能帮助你一次, 但是不能保证每一次都能帮助你...