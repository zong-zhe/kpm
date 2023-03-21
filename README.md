# kpm
KCL Package Manager

```
import (
 "bytes"
 "fmt"
 "github.com/BurntSushi/toml"
)

// 定义一个结构体类型
type Person struct {
 Name string
 Age int
 Gender string `toml:"sex"`
}

// 实现自定义序列化
func (p *Person) MarshalTOML() ([]byte, error) {
 // 创建一个buffer
 buffer := new(bytes.Buffer)
 // 向buffer中写入自定义的TOML格式数据
 fmt.Fprintf(buffer, "name = %q
", p.Name)
 fmt.Fprintf(buffer, "age = %d
", p.Age)
 fmt.Fprintf(buffer, "sex = %q
", p.Gender)
 // 返回buffer中的数据
 return buffer.Bytes(), nil
}

// 测试代码
func main() {
 // 创建Person结构体实例
 person := Person{"Tom",20, "male"}
 // 序列化Person结构体实例为TOML格式数据
 data, _ := toml.Marshal(&person)
 // 输出序列化结果
 fmt.Println(string(data))
}


```