package main
import ("fmt"
"unsafe"
)

func typesBasic(){
	fmt.Println("----基础数据类型----")
	// bool 类型
	var v_bool bool=true
	fmt.Printf("bool: %t\n",v_bool)
	// 整型 signed
	var v_int int = 25
	fmt.Printf("int: %d\n",v_int)
	// 不同长度的整数
	var v_int8 int8 = 127 // 8位整数，一个字节
	var v_int16 int16 = 30000
	var v_int32 int32 = 2000000000 
	var v_int64 int64 = 9223372036854775807
	fmt.Printf("int8: %d,int16: %d,int32: %d,int64: %d\n",v_int8,v_int16,v_int32,v_int64)
	// 无符号整数
	var v_uint8 uint8 = 220
	var v_uint16 uint16 = 50000
	var v_uint32 uint32 = 4000000000
	var v_uint64 uint64 = 18446744073709551615
	fmt.Printf("uint8: %d,uint16: %d,uint32: %d,uint64: %d\n",v_uint8,v_uint16,v_uint32,v_uint64)
	// 类型别名 uint8 也称为byte int32也称为rune
	var v_byte byte = 243
	var v_rune rune = 2000000000
	fmt.Printf("byte: %d(%c), rune: %d(%c)\n",v_byte,v_byte,v_rune,v_rune)
	// 显示类型占用字节数 unsafe.Sizeof()
	fmt.Printf("int 8 size: %d byte\n",unsafe.Sizeof(v_int8))
	fmt.Printf("int 16 size: %d bytes\n",unsafe.Sizeof(v_int16))
	fmt.Printf("int 32 size: %d bytes\n",unsafe.Sizeof(v_int32))
	fmt.Printf("int 64 size: %d bytes\n",unsafe.Sizeof(v_int64))
	// 浮点型
	var v_float32 float32 = 99.99
	var v_float64 float64 = 3.141592653
	fmt.Printf("float32: %.2f, size of float32: %d, float64: %.10f, size of float64: %d",v_float32,unsafe.Sizeof(v_float32),v_float64,unsafe.Sizeof(v_float64))
	// 字符串
	str:="Hello,你好"
	fmt.Printf("字符串长度: %d (字节数)\n", len(str)) // 注意：len返回字节数，不是字符数
	firstByte:=str[0]
	fmt.Printf("第一个字符：%d --- %c\n",firstByte,firstByte)
	// 字符串字面量
	raw:=`这
是
一
个
多行文本`
	fmt.Println(raw)
	var v_rune_c rune = '中'
	fmt.Printf("rune: %c, %d\n",v_rune_c,v_rune_c)
	// 复数类型
	var v_complex1 complex64 = 1+2i;
	var v_complex2 complex128 = 3.14 + 6.28i
	v_complex3:=complex(5.0, 10.0)
	fmt.Println("复数类型:")
	fmt.Printf("Complex: %v, real: %f, imag: %f, size: %d bytes.\n",v_complex1,real(v_complex1),imag(v_complex1),unsafe.Sizeof(v_complex1))
	fmt.Printf("Complex: %v, real: %f, imag: %f, size: %d bytes.\n",v_complex2,real(v_complex2),imag(v_complex2),unsafe.Sizeof(v_complex2))
	fmt.Printf("Complex: %v, real: %f, imag: %f, size: %d bytes.\n",v_complex3,real(v_complex3),imag(v_complex3),unsafe.Sizeof(v_complex3))
}
func arrayDemo(){
	fmt.Println("=== 数组示例 ===")
	// 初始化数组
	var v_int_arr [5]int;
	v_int_arr=[5]int{1,2,3,4,5}
	fmt.Printf("初始化int array: %v\n",v_int_arr)
	fmt.Printf("数组长度: %d\n",len(v_int_arr))
	// 自动推断数组长度
	v_int_arr2 :=[...]int{1,2,3}
	
	fmt.Printf("数组为: %v, 数组长度: %d\n",v_int_arr2,len(v_int_arr2))
}
func sliceDemo(){
	fmt.Println("=== slice切片示例 ===")
	var v_slice []int //定义一个v_slice切片
	fmt.Printf("空切片: %v, 长度: %d, 容量: %d\n",v_slice,len(v_slice),cap(v_slice))
	v_slice=make([]int, 3,5)
	fmt.Printf("长度3容量5的切片: %v, 长度: %d, 容量: %d\n",v_slice,len(v_slice),cap(v_slice))
	// 长度和容量的区别
	v_slice=[]int{1,2,3,4,5}
	fmt.Printf("初始化之后的切片: %v\n",v_slice)
	// 切片可以追加容量
	v_slice=append(v_slice,6)
	fmt.Printf("append追加后的的切片: %v, 长度: %d, 容量: %d\n",v_slice,len(v_slice),cap(v_slice))
	// 切片截取
	subSlice:=v_slice[1:4]
	fmt.Printf("[1:4]的切片内容: %v\n",subSlice)
	// 当修改subSlice之后观察v_slice的变化
	subSlice[0]=999
	fmt.Printf("subSlice内容: %v, 原slice的内容: %v\n",subSlice,v_slice)
	// 说明切片共享一段内容
}
func mappingDemo(){
	// 声明一个map
	var v_map1 map[string]int
	fmt.Printf("声明但不初始化一个map: %v\n",v_map1)
	// make声明
	v_map2:=make(map[string]int)
	v_map2["apple"]=12
	v_map2["orange"]=15
	fmt.Printf("make声明一个map: %v\n",v_map2)
	// 声明并初始化
	v_map3:=map[string]int{
		"apple":21,
		"orange":32,
	}
	fmt.Printf("声明并初始化一个map: %v\n",v_map3)
	// 检查key是否存在
	value,err:=v_map3["apple"]
	fmt.Printf("apple: %d, %v\n",value,err)
	_,err2:=v_map3["apple1"]
	if(!err2){
		fmt.Println("key不存在")
	}
	// 删除键
	delete(v_map2,"apple")
	fmt.Printf("删除apple键: %v\n",v_map2)
	// 遍历映射
	fmt.Println("遍历映射")
	for k,v := range v_map3{
		fmt.Printf("%v ----- %v\n",k,v)
	}
}
func main(){
	fmt.Println("helloworld")
	typesBasic()
	arrayDemo()
	sliceDemo()
	mappingDemo()
}
