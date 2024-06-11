<p align="center">
  <h3 align="center">Restory</h3>
  <p align="center">
    Re (store + (hi) story) Based on VSCode History File Recovery Tool (and directory structure)
    <br />
    基于 VSCode 历史记录的文件 (包括目录结构) 恢复工具
  </p>
</p>

[中文 README](README.md)

> Are you using VSCode but accidentally deleted some of your files and folders? Too lazy to find each historical record one by one? Restory can help you restore the entire directory with one click!

### Principle

This is a very simple tool. Restory scans the file modification history files of VSCode to obtain information, and rebuilds the directory structure and file content at the location you specify.

**Unfortunately, files that were not modified through VSCode cannot be recovered using Restory, as VSCode does not generate the corresponding modification history files.**

### Usage

Download the binary file of `restory` for your platform from the Release page and run it.

For example, on the Linux platform:

```bash
$ curl -L -o ./restory_linux_amd64 https://github.com/ClarkQAQ/restory/releases/latest/download/restory_linux_amd64
$ chmod +x ./restory_linux_amd64
$ ./restory_linux_amd64 ~/.config/Code/User/History /tmp/code-history
```

The VSCode history folder is usually located at `~/.config/Code/User/History` on Linux and `C:\Users\[your username]\AppData\Local\Code\User\History` on Windows.

### Quick Release Downloads

- [Windows Download Link](https://github.com/ClarkQAQ/restory/releases/latest/download/restory_windows_amd64.exe)
- [Linux Download Link](https://github.com/ClarkQAQ/restory/releases/latest/download/restory_linux_amd64)
- [MacOS Download Link](https://github.com/ClarkQAQ/restory/releases/latest/download/restory_darwin_amd64)

### Finally

Please remember to back up your important files regularly! Upload them to Git timely!
This tool might help you once, but it cannot guarantee to help you every time...