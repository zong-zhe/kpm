<h1 align="center">Kpm: KCL 包管理器</h1>

<p align="center">
<a href="./README.md">English</a> | <a href="./README-zh.md">简体中文</a>
</p>
<p align="center">
<a href="#introduction">介绍</a> | <a href="#installation">安装</a> | <a href="#quick-start">快速开始</a>
</p>

<p align="center">
<img src="https://coveralls.io/repos/github/KusionStack/kpm/badge.svg">
<img src="https://img.shields.io/badge/license-Apache--2.0-green">
<img src="https://img.shields.io/badge/PRs-welcome-brightgreen">
</p>

## 介绍

`kpm` 是 KCL 包管理工具， `kpm` 作为 [kcl cli](https://github.com/kcl-lang/cli) 的三方库，负责 [`kcl mod`](https://kcl-lang.io/docs/tools/cli/package-management/command-reference/init) 功能的实现，`kpm` 会下载您的 KCL 包的依赖项、编译您的 KCL 包、制作可分发的包并将其上传到 KCL 包的仓库中。

## 了解更多

- [OCI registry 支持](./docs/kpm_oci-zh.md).
- [如何使用 kpm 与他人分享您的 kcl 包](./docs/publish_your_kcl_packages-zh.md)
- [如何使用 kpm 在 docker.io 上与他人分享您的 kcl 包](./docs/publish_to_docker_reg-zh.md)
- [kpm 命令参考](./docs/command-reference-zh/index.md)
- [kcl.mod: KCL 包清单文件](./docs/kcl_mod-zh.md)
- [如何使用 kpm 通过 github action 来推送您的 kcl 包](./docs/push_by_github_action-zh.md)
- [发布 KCL 包到官方 Registry](./docs/publish_pkg_to_ah-zh.md)
