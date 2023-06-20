[package]
name = "kpm"
edition = "0.0.1"
version = "0.0.1"

[dependencies]
sub = "0.0.1"

[profile]
entries = ["main.k", "xxx/xxx/dir"]
options = ["a=b", "b=c"]
overrides = ["a=b", "b=c"]
disable_none = true
sort_key = true
settings = ["xxx/xxx/kcl.yaml"]
