# kpm
KCL Package Manager

```
package main

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Name string
	Age  int
	Map  map[string]string
}

// MarshalTOML 实现了 toml.Marshaler 接口，以自定义 Config 的 TOML 编码。
func (c Config) MarshalTOML() ([]byte, error) {
	// 创建一个缓冲区以保存编码后的数据
	var buf bytes.Buffer

	// 将 Name 和 Age 字段编码为 TOML 并写入缓冲区
	if err := toml.NewEncoder(&buf).Encode(struct {
		Name string
		Age  int
		Map  string
	}{
		Name: c.Name,
		Age:  c.Age,
		Map:  "getMapValues(c.Map)",
	}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	config := Config{
		Name: "Tom",
		Age:  25,
		Map: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	data, err := config.MarshalTOML()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}

```