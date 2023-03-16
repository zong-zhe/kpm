# kpm
KCL Package Manager

···
package main

import (
 "fmt"
 "os"

 "github.com/go-git/go-git/v5"
 "github.com/go-git/go-git/v5/plumbing"
 "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
 //仓库地址
 repoURL := "https://github.com/username/repo.git"
 //本地存储路径
 localPath := "/path/to/local/repo"
 // 授权信息
 auth := &http.BasicAuth{
 Username: "your_username",
 Password: "your_password",
 }

 // 克隆仓库到本地
 _, err := git.PlainClone(localPath, false, &git.CloneOptions{
 URL: repoURL,
 Auth: auth,
 Progress: os.Stdout,
 })
 if err != nil {
 fmt.Println("Clone error:", err)
 return
 }

 // 打开本地仓库
 r, err := git.PlainOpen(localPath)
 if err != nil {
 fmt.Println("Open error:", err)
 return
 }

 // 获取所有分支
 branches, err := r.Branches()
 if err != nil {
 fmt.Println("Branches error:", err)
 return
 }
 // 遍历所有分支
 err = branches.ForEach(func(ref *plumbing.Reference) error {
 fmt.Println("Branch:", ref.Name())
 return nil
 })
 if err != nil {
 fmt.Println("Branches error:", err)
 return
 }

 // 获取所有提交
 commits, err := r.CommitObjects()
 if err != nil {
 fmt.Println("Commits error:", err)
 return
 }
 // 遍历所有提交
 err = commits.ForEach(func(commit *object.Commit) error {
 fmt.Println("Commit:", commit.Hash.String())
 return nil
 })
 if err != nil {
 fmt.Println("Commits error:", err)
 return
 }

 // 获取所有标签
 tags, err := r.Tags()
 if err != nil {
 fmt.Println("Tags error:", err)
 return
 }
 // 遍历所有标签
 err = tags.ForEach(func(ref *plumbing.Reference) error {
 fmt.Println("Tag:", ref.Name())
 return nil
 })
 if err != nil {
 fmt.Println("Tags error:", err)
 return
 }
}

···