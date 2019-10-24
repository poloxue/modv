 > [本文](http://poloxue.com/go/go-module-relationship-visible-tool)首发于[我的博客](http://poloxue.com)，如果觉得有用，欢迎点赞收藏，让更多的朋友看到。

最近，我开发了一个非常简单的小工具，总的代码量 200 行不到。今天，简单介绍下它。这是个什么工具呢？它是一个用于可视化展示 Go Module 依赖关系的工具。

# 为何开发

为什么会想到开发这个工具？主要有两点原因：

一是最近经常看到大家在社区讨论 Go Module。于是，我也花了一些时间研究了下。期间，遇到了一个需求，如何清晰地识别模块中依赖项之间的关系。一番了解后，发现了 `go mod graph`。

效果如下：

```
$ go mod graph
github.com/poloxue/testmod golang.org/x/text@v0.3.2
github.com/poloxue/testmod rsc.io/quote/v3@v3.1.0
github.com/poloxue/testmod rsc.io/sampler@v1.3.1
golang.org/x/text@v0.3.2 golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e
rsc.io/quote/v3@v3.1.0 rsc.io/sampler@v1.3.0
rsc.io/sampler@v1.3.1 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c
rsc.io/sampler@v1.3.0 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c
```

每一行的格式是 `模块 依赖模块`，基本能满足要求，但总觉得还是不那么直观。

二是我之前手里有一个项目，包管理一直用的是 dep。于是，我也了解了下它，把官方文档仔细读了一遍。其中的[某个章节](https://golang.github.io/dep/docs/daily-dep.html)介绍了依赖项可视化展示的方法。

文档中给出的包关系图：

<center>
<img src="https://blogimg.poloxue.com/0014-go-mod-graph-visible-02.png">
</center>

看到这张图的时候，眼睛瞬间就亮了，图形化就是优秀，不同依赖之间的关系一目了然。这不就是我想要的效果吗？666，点个赞。

但 ...，随之而来的问题是，go mod 没这个能力啊。怎么办？

# 如何实现

先看看是不是已经有人做了这件事了。网上搜了下，没找到。那是不是能自己实现？应该可以借鉴下 dep 的思路吧？

如下是 dep 依赖实现可视化的方式：

```bash
# linux
$ sudo apt-get install graphviz
$ dep status -dot | dot -T png | display

# macOS
$ brew install graphviz
$ dep status -dot | dot -T png | open -f -a /Applications/Preview.app

# Windows
> choco install graphviz.portable
> dep status -dot | dot -T png -o status.png; start status.png
```

这里展示了三大系统下的使用方式，它们都安装了一个软件包，graphviz。从名字上看，这应该是一个用来实现可视化的软件，即用来画图的。事实也是这样，可以看看它的[官网](http://www.graphviz.org/documentation/)。

再看下它的使用，发现都是通过管道命令组合的方式，而且前面的部分基本相同，都是 `dep status -dot | dot -T png`。后面的部分在不同的系统就不同了，Linux 是 `display`，MacOS 是 `open -f -a /Applications/Preview.app`，Window 是 `start status.png`。

稍微分析下就会明白，前面是生成图片，后面是显示图片。因为不同系统的图片展示命令不同，所以后面的部分也就不同了。

现在关心的重点在前面，即 `dep status -dot | dot -T png` 干了啥，它究竟是如何实现绘图的？大致猜测，dot -T png 是由 dep status -dot 提供的数据生成图片。继续看下 `dep status -dot` 的执行效果吧。

```bash
$ dep status -dot
digraph {
	node [shape=box];
	2609291568 [label="github.com/poloxue/hellodep"];
	953278068 [label="rsc.io/quote\nv3.1.0"];
	3852693168 [label="rsc.io/sampler\nv1.0.0"];
	2609291568 -> 953278068;
	953278068 -> 3852693168;
}
```

咋一看，输出的是一段看起来不知道是啥的代码，这应该是 graphviz 用于绘制图表的语言。那是不是还有学习下？当然不用啊，这里用的很简单，直接套用就行了。

试着分析一下吧，前面两行可以不用关心，这应该是 graphviz 特定的写法，表示要画的是什么图。我们主要关心如何将数据以正确形式提供出来。

```bash
2609291568 [label="github.com/poloxue/hellodep"];
953278068 [label="rsc.io/quote\nv3.1.0"];
3852693168 [label="rsc.io/sampler\nv1.0.0"];
2609291568 -> 953278068;
953278068 -> 3852693168;
```

一看就知道，这里有两种结构，分别是为依赖项关联 ID ，和通过 ID 和 `->` 表示依赖间的关系。

按上面的猜想，我们可以试着画出一个简单的图, 用于表示 a 模块依赖 b 模块。执行命令如下，将绘图代码通过 `each` 管道的方式发送给 `dot` 命令。

```
$ echo 'digraph {
node [shape=box];
1 [label="a"];
2 [label="b"];
1 -> 2;
}' | dot -T png | open -f -a /Applications/Preview.app 
```

效果如下：

<center>
<img src="https://blogimg.poloxue.com/0014-go-mod-graph-visible-03.png"/>
</center>

绘制一个依赖关系图竟然这么简单。

看到这里，是不是发现问题已经变得非常简单了。我们只要将 `go mod graph` 的输出转化为类似的结构就能实现可视化了。

# 开发流程介绍

接下来，开发这个小程序吧，我将这个小程序命名为 `modv`，即 module visible 的意思。项目源码位于 [poloxue/modv](https://github.com/poloxue/modv)。

# 接收管道的输入

先要检查数据输入管道是否正常。

我们的目标是使用类似 `dep` 中作图的方式，`go mod graph` 通过管道将数据传递给 `modv`。因此，要先检查 `os.Stdin`，即检查标准输入状态是否正常， 以及是否是管道传输。

下面是 main 函数的代码，位于 [main.go](https://github.com/poloxue/modv/blob/master/main.go) 中。

```go
func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("os.Stdin.Stat:", err)
		PrintUsage()
		os.Exit(1)
	}

	// 是否是管道传输
	if info.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("command err: command is intended to work with pipes.")
		PrintUsage()
		os.Exit(1)
	}
```

一旦确认输入设备一切正常，我们就可以进入到数据读取、解析与渲染的流程了。

```go
	mg := NewModuleGraph(os.Stdin)
	mg.Parse()
	mg.Render(os.Stdout)
}
```

接下来，开始具体看看如何实现数据的处理流程。

# 抽象实现结构

先定义一个结构体，并大致定义整个流程。

```go
type ModGraph struct {
	Reader io.Reader  // 读取数据流
}

func NewModGraph(r io.Reader) *ModGraph {
    return &ModGraph{Reader: r}
}

// 执行数据的处理转化
func (m *ModGraph) Parse() error {}

// 结果渲染与输出
func (m *ModGraph) Render(w io.Writer) error {}
```

再看下 `go mod graph` 的输出吧，如下：

```bash
github.com/poloxue/testmod golang.org/x/text@v0.3.2
github.com/poloxue/testmod rsc.io/quote/v3@v3.1.0
...
```

每一行的结构是 `模块 依赖项`。现在的目标是要它解析成下面这样的结构：

```bash
digraph {
    node [shape=box];
    1 github.com/poloxue/testmod;
    2 golang.org/x/text@v0.3.2;
    3 rsc.io/quote/v3@v3.1.0;
    1 -> 2;
    1 -> 3;
}
```

前面说过，这里包含了两种不同的结构，分别是模块与 ID 关联关系，以及模块 ID 表示模块间的依赖关联。为 `ModGraph` 结构体增加两个成员表示它们。

```go
type ModGraph struct {
	r io.Reader  // 数据流读取实例，这里即 os.Stdin
 
	// 每一项名称与 ID 的映射
	Mods         map[string]int
	// ID 和依赖 ID 关系映射，一个 ID 可能依赖多个项
	Dependencies map[int][]int
}
```

要注意的是，增加了两个 map 成员后，记住要在 `NewModGraph` 中初始化下它们。

# mod graph 输出解析

如何进行解析？

介绍到这里，目标已经很明白了。就是要将输入数据解析到 `Mods` 和 `Dependencies` 两个成员中，实现代码都在 `Parse` 方法中。

为了方便进行数据读取，首先，我们利用 `bufio` 基于 `reader` 创建一个新的 `bufReader`，

```go
func (m *ModGraph) Parse() error {
	bufReader := bufio.NewReader(m.Reader)
	...
```

为便于按行解析数据，我们通过 bufReader 的 `ReadBytes()` 方法循环一行一行地读取 os.Stdin 中的数据。然后，对每一行数据按空格切分，获取到依赖关系的两项。代码如下：

```go
for {
	relationBytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

    relation := bytes.Split(relationBytes, []byte(" "))
    // module and dependency
    mod, depMod := strings.TrimSpace(string(relation[0])), strings.TrimSpace(string(relation[1]))

    ...
}
```

接下来，就是将解析出来的依赖关系组织到 `Mods` 和 `Dependencies` 两个成员中。模块 ID 是生成规则采用的是最简单的实现方式，从 1 自增。实现代码如下：

```go
modId, ok := m.Mods[mod]
if !ok {
	modId = serialID
	m.Mods[mod] = modId
	serialID += 1
}

depModId, ok := m.Mods[depMod]
if !ok {
	depModId = serialID
	m.Mods[depMod] = depModId
	serialID += 1
}

if _, ok := m.Dependencies[modId]; ok {
	m.Dependencies[modId] = append(m.Dependencies[modId], depModId)
} else {
	m.Dependencies[modId] = []int{depModId}
}
```

解析的工作到这里就结束了。

# 渲染解析的结果

这个小工具还剩下最后一步工作要做，即将解析出来的数据渲染出来，以满足 `graphviz` 工具的作图要求。实现代码是 `Render`部分：

首先，定义一个模板，以生成满足要求的输出格式。

```go
var graphTemplate = `digraph {
node [shape=box];
{{ range $mod, $modId := .mods -}}
{{ $modId }} [label="{{ $mod }}"];
{{ end -}}
{{- range $modId, $depModIds := .dependencies -}}
{{- range $_, $depModId := $depModIds -}}
{{ $modId }} -> {{ $depModId }};
{{  end -}}
{{- end -}}
}
`
```

这一块没啥好介绍的，主要是要熟悉 Go 中的 `text/template` 模板的语法规范。为了展示友好，这里通过 `-` 实现换行的去除，整体而言不影响阅读。

接下来，看 `Render` 方法的实现，把前面解析出来的 `Mods` 和 `Dependencies` 放入模板进行渲染。

```go
func (m *ModuleGraph) Render(w io.Writer) error {
	templ, err := template.New("graph").Parse(graphTemplate)
	if err != nil {
		return fmt.Errorf("templ.Parse: %v", err)
	}

	if err := templ.Execute(w, map[string]interface{}{
		"mods":         m.Mods,
		"dependencies": m.Dependencies,
	}); err != nil {
		return fmt.Errorf("templ.Execute: %v", err)
	}

	return nil
}
```

现在，全部工作都完成了。最后，将这个流程整合到 main 函数。接下来就是使用了。

# 使用体验

开始体验下吧。补充一句，这个工具，我现在只测试了 Mac 下的使用，如有问题，欢迎提出来。

首先，要先安装一下 `graphviz`，安装的方式在本文开头已经介绍了，选择你的系统安装方式。

接着是安装 `modv`，命令如下：

```bash
$ go get github.com/poloxue/modv
```

安装完成！简单测试下它的使用。

以 MacOS 为例。先下载测试库，github.com/poloxue/testmod。 进入 testmod 目录执行命令：

```bash
$ go mod graph | modv | dot -T png | open -f -a /Applications/Preview.app
```

如果执行成功，将看到如下的效果：

![](https://blogimg.poloxue.com/0014-go-mod-graph-visible-04.png)

完美地展示了各个模块之间的依赖关系。

# 一些思考

本文是篇实践性的文章，从一个简单想法到成功呈现出一个可以使用的工具。虽然，开发起来并不难，从开发到完成，仅仅花了一两个小时。但我的感觉，这确实是个有实际价值的工具。

还有一些想法没有实现和验证，比如一旦项目较大，是否可以方便的展示某个指定节点的依赖树，而非整个项目。还有，在其他项目向 Go Module 迁移的时候，这个小工具是否能产生一些价值。
