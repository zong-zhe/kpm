# 关于 kcl 编译器的编译入口

编译入口是指 KCL 每次编译活动使用那些 KCL 代码进行编译。

KCL 编译器支持编译来自远端和本地文件系统三个来源的 KCL 程序。

## 编译来自远端的 KCL 程序

KCL 目前支持通过 URL 编译来自 git 仓库和 oci 仓库的包。

以 git 为例：
```
kcl run git://github.com/test/helloworld
```

在编译来自远端的包时，

1. URL 对应的内容必须是一个完整的包含 kcl.mod 文件的 KCL 包。
2. 远端包的 kcl.mod 文件中如果没有指定 entiries，那么就默认使用 kcl.mod 同级目录的全部 *.k 文件进行编译。
3. 远端包的 kcl.mod 文件中如果指定了 entiries，那么就默认使用 kcl.mod 文件中指定的文件进行编译。

这样设计的主要考虑是 打包好的 KCL 包可以通过里面的 kcl.mod 文件指定特定的编译入口，如果不指定，默认就是和 kcl.mod 同一个目录的 *.k 文件。

## 编译本地的 KCL 程序

本地的 KCL 程序分为如下情况：

1. KCL 文件 *.k
2. 包含 kcl.mod 的文件目录
3. 不包含 kcl.mod 的文件目录 

kcl 编译器在编译本地的 KCL 程序时候，会优先寻找 RootPath, 

1. 如果编译一个文件目录 a，会现在 a 目录内寻找 kcl.mod, 如果能找到，a 目录就是 RootPath; 如果 a 目录中没有，就递归的向上级目录寻找，直到找到 kcl.mod，包含 kcl.mod 的目录就是 RootPath; 如果一直找不到，a 目录本身就是 RootPath。
2. 如果编译一个 KCL 文件 b.k, 就递归的想上级目录寻找 kcl.mod, 如果找到 kcl.mod，包含 kcl.mod 的目录就是 RootPath; 如果一直找不到，b.k 的上级目录就是 RootPath。
3. 如果编译一个多个文件路径，包含文件目录和 *.k 文件，则会针对每个文件路径，按照 1 和 2 规则寻找 RootPath，如果找到了多个 RootPath 则以 WorkDir 作为 RootPath。WorkDir 是 kcl.yaml 或者 pwd。





其中：
- git/oci 来源的 KCL 程序必须是一个完整的 KCL 包，在包的根目录下必须包含 kcl.mod。
- 本地文件系统的 KCL 程序主要包括：*.k 文件，包含 kcl.mod 文件的文件目录，不包含 kcl.mod 文件的文件目录。

KCL 目前不支持同时编译多个 KCL 包，即如果从 KCL 编译入口中找到了多个 kcl.mod 文件，将无法通过编译。

因此：

- kcl run 传入多个 git url 或者 oci url 将会得到编译错误。因为，git/oci 来源的 KCL 程序必须是一个完整的 KCL 包，在包的根目录下必须包含 kcl.mod。
- kcl run 传入多个文件目录，如果每个文件目录下都有 kcl.mod，编译也将会报错。
