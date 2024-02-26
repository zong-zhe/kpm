<h1 align="center">Kpm: KCL Package Manager</h1>

<p align="center">
<a href="./README.md">English</a> | <a href="./README-zh.md">简体中文</a>
</p>
<p align="center">
<a href="#introduction">Introduction</a> | <a href="#installation">Installation</a> | <a href="#quick-start">Quick start</a> 
</p>

<p align="center">
<img src="https://coveralls.io/repos/github/KusionStack/kpm/badge.svg">
<img src="https://img.shields.io/badge/license-Apache--2.0-green">
<img src="https://img.shields.io/badge/PRs-welcome-brightgreen">
</p>

## Introduction

`kpm` is the KCL package manager. `kpm` downloads your KCL package's dependencies, compiles your KCL packages, makes packages, and uploads them to the kcl package registry.

`kpm` is a third-party library of [kcl cli](https://github.com/kcl-lang/cli), which is responsible for the implementation of [`kcl mod`](https://kcl-lang.io/docs/tools/cli/package-management/command-reference/init), `kpm` will download your KCL package's dependencies, compile your KCL package, make a distributable package and upload it to the KCL package's repository. 

## Contributing

- See [contribution guideline](https://kcl-lang.io/docs/community/contribute/).

## Learn More

- [OCI registry support](./docs/kpm_oci-zh.md).
- [How to share your kcl package with others using kpm](./docs/publish_your_kcl_packages-zh.md).
- [How to use kpm to share your kcl package with others on docker.io](./docs/publish_to_docker_reg.md)
- [kpm command reference](./docs/command-reference/index.md)
- [kcl.mod: The KCL package Manifest File](./docs/kcl_mod.md)
- [How to use kpm to push your kcl package by github action](./docs/push_by_github_action.md)
- [How to publish KCL package to official Registry](./docs/publish_pkg_to_ah.md)
