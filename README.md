## 动态链接命令含打包工具

本程序实现了一个简单的功能，用来打包动态链接二进制文件。除内核API限制以外，可以使得不同系统版本之间复制命令行更加 “easy and dirty”

### 使用说明

```
sucker -export-dir ./destdir -exec-prefix /home/test/runtest /bin/df
```

本例将 df 工具打包到当前目录下的 __destdir__
目录中，此后将该目录中内容拷贝到任意机器的 /home/test/runtest
下都可以运行，不依赖目标机器除内核以外的任何文件

