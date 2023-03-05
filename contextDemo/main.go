package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// ============= context.WithValue ============ //
	//ctx := context.WithValue(context.Background(), "k1", "v1")
	//// value 是一个interface 类型，所以在取值时需要对其进行断言
	//value := ctx.Value("k1")
	//print(value.(string))

	//UseWithCancel()

	//UseWithDeadline()

	//ParentSonValue()

	//ParentGetSonValue()

	//ParentSonControlTimeout()

	//TimeoutControl()

}

// TimeoutControl 超时控制
func TimeoutControl() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	bsChan := make(chan struct{})
	go func() {
		Business()
		bsChan <- struct{}{}
	}()
	// 监听业务处理和超时控制
	select {
	case <-ctx.Done():
		fmt.Println("超时啦。。。")
	case <-bsChan:
		fmt.Println("业务处理完了")
	}
}
func Business() {
	fmt.Println("处理业务")
	time.Sleep(time.Second * 2)
}

// ParentSonControlTimeout 父子对超时的控制区别
func ParentSonControlTimeout() {
	ctx := context.Background()
	parentCtx, cancel1 := context.WithTimeout(ctx, time.Second)
	sonCtx, cancel2 := context.WithTimeout(parentCtx, 3*time.Second)
	go func() {
		<-sonCtx.Done()
		fmt.Printf("超时啦。。。\n")
	}()
	time.Sleep(2 * time.Second)
	fmt.Printf("程序结束\n")
	cancel2()
	cancel1()
}

// ParentGetSonValue 父设置子的值（成功版本，不推荐使用map，因为这将使得context变成可变对象）
func ParentGetSonValue() {
	ctx := context.Background()
	parentCtx := context.WithValue(ctx, "map", map[string]string{})
	sonCtx := context.WithValue(parentCtx, "k2", "v2")
	value := sonCtx.Value("map")
	if val, ok := value.(map[string]string); ok {
		val["k3"] = "v3"
	}
	fmt.Printf("get k2: %v \n", parentCtx.Value("k2"))
	fmt.Printf("get k3: %v \n", parentCtx.Value("map").(map[string]string)["k3"])
}

// ParentSonValue 设置子的值，无效版本
func ParentSonValue() {
	ctx := context.Background()
	parentCtx := context.WithValue(ctx, "k1", "v1")
	sonCtx := context.WithValue(parentCtx, "k2", "v2")
	fmt.Printf("%v \n", parentCtx.Value("k2"))
	fmt.Printf("%v \n", sonCtx.Value("k2"))
}

// UseWithDeadline WithDeadline方法使用
func UseWithDeadline() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(500*time.Millisecond))
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("开始执行")
	case <-ctx.Done():
		fmt.Println("任务结束了")
	}
}

// UseWithCancel WithCancel方法的使用
func UseWithCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	for i := 1; i < 4; i++ {
		go worker(ctx, i)
	}
	time.Sleep(time.Second * 3)
	cancel()
	// 阻塞，看worker是否还在工作
	time.Sleep(time.Second * 3)
}

func worker(ctx context.Context, i int) {
	// 持续处理业务
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("【%d】 业务关闭了。。。\n", i)
			return
		default:
			fmt.Printf("【%d】 开始处理业务....\n", i)
			time.Sleep(time.Second)
		}
	}
}
