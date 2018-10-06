# SELect PaGes (selpg)
## 功能
selpg从标准输入或从作为命令行参数给出的文件名读取文本输入。它允许用户指定来自该输入并随后将被输出的页面范围。例如，如果输入含有 100 页，则用户可指定只打印第 35 至 65 页。这种特性有实际价值，因为在打印机上打印选定的页面避免了浪费纸张。另一个实例是，原始文件很大而且以前已打印过，但某些页面由于打印机卡住或其它原因而没有被正确打印。在这样的情况下，则可用该工具来只打印需要打印的页面。
## 用法
> ```selpg -s Number -e Number [-f|-l Number] [fileName]```

如前面所说的那样，selpg 是从文本输入选择页范围的实用程序。该输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入。

selpg 首先处理所有的命令行参数。在扫描了所有的选项参数（也就是那些以连字符为前缀的参数）后，如果 selpg 发现还有一个参数，则它会接受该参数为输入文件的名称并尝试打开它以进行读取。如果没有其它参数，则 selpg 假定输入来自标准输入。

参数处理
“-sNumber”和“-eNumber”强制选项：
selpg 要求用户用两个命令行参数“-sNumber”（例如，“-s10”表示从第 10 页开始）和“-eNumber”（例如，“-e20”表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。

```$ selpg -s10 -e20 ...```

（... 是命令的余下部分，下面对它们做了描述）。

“-lNumber”和“-f”可选选项：
selpg 可以处理两种输入文本：

### 类型 1
该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-lNumber”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。该缺省值可以用“-lNumber”选项覆盖，如下所示：

```$ selpg -s10 -e20 -l66 ...```

这表明页有固定长度，每页为 66 行。

### 类型 2
该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用“\f”表示）定界。

类型 2 格式由“-f”选项表示，如下所示：

```$ selpg -s10 -e20 -f ...```

该命令告诉 selpg 在输入中寻找换页符，并将其作为页定界符处理。

注：“-lNumber”和“-f”选项是互斥的。

### “-dDestination”可选选项
selpg 还允许用户使用“-dDestination”选项将选定的页直接发送至打印机。这里，“Destination”应该是 lp 命令“-d”选项（请参阅“man lp”）可接受的打印目的地名称。该目的地应该存在。在运行了带“-d”选项的 selpg 命令后，若要验证该选项是否已生效，请运行命令“lpstat -t”。该命令应该显示添加到“Destination”打印队列的一项打印作业。如果当前有打印机连接至该目的地并且是启用的，则打印机应打印该输出。在下面的示例中，我们打开到命令

```$ lp -dDestination```

的管道以便输出，并写至该管道而不是标准输出：

```selpg -s10 -e20 -dlp1```

该命令将选定的页作为打印作业发送至 lp1 打印目的地。您应该可以看到类似“request id is lp1-6”的消息。该消息来自 lp 命令；它显示打印作业标识。如果在运行 selpg 命令之后立即运行命令 lpstat -t | grep lp1 ，您应该看见 lp1 队列中的作业。如果在运行 lpstat 命令前耽搁了一些时间，那么您可能看不到该作业，因为它一旦被打印就从队列中消失了。

输入处理
一旦处理了所有的命令行参数，就使用这些指定的选项以及输入、输出源和目标来开始输入的实际处理。

selpg 通过以下方法记住当前页号：如果输入是每页行数固定的，则 selpg 统计新行数，直到达到页长度后增加页计数器。如果输入是换页定界的，则 selpg 改为统计换页符。这两种情况下，只要页计数器的值在起始页和结束页之间这一条件保持为真，selpg 就会输出文本（逐行或逐字）。当那个条件为假（也就是说，页计数器的值小于起始页或大于结束页）时，则 selpg 不再写任何输出。瞧！您得到了想输出的那些页。
## 使用样例
为了演示最终用户可以如何应用我们所介绍的一些原则，下面给出了可使用的 selpg 命令字符串示例：

1. ```$ selpg -s1 -e1 input_file```

> 该命令将把“input_file”的第 1 页写至标准输出（也就是屏幕），因为这里没有重定向或管道。

2. ```$ selpg -s1 -e1 < input_file```

> 该命令与示例 1 所做的工作相同，但在本例中，selpg 读取标准输入，而标准输入已被 shell／内核重定向为来自“input_file”而不是显式命名的文件名参数。输入的第 1 页被写至屏幕。

3. ```$ other_command | selpg -s10 -e20```

> “other_command”的标准输出被 shell／内核重定向至 selpg 的标准输入。将第 10 页到第 20 页写至 selpg 的标准输出（屏幕）。

4. ```$ selpg -s10 -e20 input_file >output_file```

> selpg 将第 10 页到第 20 页写至标准输出；标准输出被 shell／内核重定向至“output_file”。

5. ```$ selpg -s10 -e20 input_file 2>error_file```

> selpg 将第 10 页到第 20 页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“error_file”。请注意：在“2”和“>”之间不能有空格；这是 shell 语法的一部分（请参阅“man bash”或“man sh”）。

6. ```$ selpg -s10 -e20 input_file >output_file 2>error_file```

> selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至“error_file”。当“input_file”很大时可使用这种调用；您不会想坐在那里等着 selpg 完成工作，并且您希望对输出和错误都进行保存。

7. ```$ selpg -s10 -e20 input_file >output_file 2>/dev/null```

> selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至 /dev/null（空设备），这意味着错误消息被丢弃了。设备文件 /dev/null 废弃所有写至它的输出，当从该设备文件读取时，会立即返回 EOF。

8. ```$ selpg -s10 -e20 input_file >/dev/null```

> selpg 将第 10 页到第 20 页写至标准输出，标准输出被丢弃；错误消息在屏幕出现。这可作为测试 selpg 的用途，此时您也许只想（对一些测试情况）检查错误消息，而不想看到正常输出。

9. ```$ selpg -s10 -e20 input_file | other_command```

> selpg 的标准输出透明地被 shell／内核重定向，成为“other_command”的标准输入，第 10 页到第 20 页被写至该标准输入。“other_command”的示例可以是 lp，它使输出在系统缺省打印机上打印。“other_command”的示例也可以 wc，它会显示选定范围的页中包含的行数、字数和字符数。“other_command”可以是任何其它能从其标准输入读取的命令。错误消息仍在屏幕显示。

10. ```$ selpg -s10 -e20 input_file 2>error_file | other_command```

> 与上面的示例 9 相似，只有一点不同：错误消息被写至“error_file”。

在以上涉及标准输出或标准错误重定向的任一示例中，用“>>”替代“>”将把输出或错误数据附加在目标文件后面，而不是覆盖目标文件（当目标文件存在时）或创建目标文件（当目标文件不存在时）。

以下所有的示例也都可以（有一个例外）结合上面显示的重定向或管道命令。我没有将这些特性添加到下面的示例，因为我认为它们在上面示例中的出现次数已经足够多了。例外情况是您不能在任何包含“-dDestination”选项的 selpg 调用中使用输出重定向或管道命令。实际上，您仍然可以对标准错误使用重定向或管道命令，但不能对标准输出使用，因为没有任何标准输出 — 正在内部使用 popen() 函数由管道将它输送至 lp 进程。

11. ```$ selpg -s10 -e20 -l66 input_file```

> 该命令将页长设置为 66 行，这样 selpg 就可以把输入当作被定界为该长度的页那样处理。第 10 页到第 20 页被写至 selpg 的标准输出（屏幕）。

12. ```$ selpg -s10 -e20 -f input_file```

> 假定页由换页符定界。第 10 页到第 20 页被写至 selpg 的标准输出（屏幕）。

13. ```$ selpg -s10 -e20 -dlp1 input_file```

> 第 10 页到第 20 页由管道输送至命令“lp -dlp1”，该命令将使输出在打印机 lp1 上打印。

最后一个示例将演示 Linux shell 的另一特性：

14. ```$ selpg -s10 -e20 input_file > output_file 2>error_file &```

> 该命令利用了 Linux 的一个强大特性，即：在“后台”运行进程的能力。在这个例子中发生的情况是：“进程标识”（pid）如 1234 将被显示，然后 shell 提示符几乎立刻会出现，使得您能向 shell 输入更多命令。同时，selpg 进程在后台运行，并且标准输出和标准错误都被重定向至文件。这样做的好处是您可以在 selpg 运行时继续做其它工作。

## 代码设计
### 参数的读取
我们在代码中使用pflag包进行参数的读取。这个包的功能十分完善，参数处理十分方便，此处不多加阐述。
```golang
flag.IntVarP(&arg_start, "start", "s", 0, "the start page number(include)")
...
```
### Reader和Writer的使用
在程序设计时，我们的输入有两种途径：键盘输入和文件输入，同样，输出也有两种途径：输出到屏幕，输出到lp程序。为了简化设计，我们使用Reader和Writer规范输入输出，两种输入途径被统一为一个Reader，输入途径被统一成一个Writer,这样子简化了设计。如输出的代码：
```golang
printPage := func(page_ctr int, page string) {
  if page_ctr >= arg_start && page_ctr <= arg_end {
    _, err := writer.WriteString(page)
    if err != nil {
      println(err.Error())
      os.Exit(2)
    }
  }
}
```
### 与lp程序的通信（pipe）
在代码中，我们用到了多个技巧：
- 使用```var wg sync.WaitGroup```使主线程等待lp子线程
- 使用io.Pipe把exec.Cmd类的Stdin接口接入到向外的bufio.Writer上
- 由于Pipe产生的Writer最后需要调用Close()方法标识输入结束但bufio的Writer并没有提供Close方法，因此我们使用_closeWriter函数来统一处理输入完毕的情况。

```golang
var writer *bufio.Writer
_closeWriter := func() {

}
var wg sync.WaitGroup
defer wg.Wait()
if arg_destination != "" {
  r, w := io.Pipe()
  _closeWriter = func() {
    w.Close()
  }
  cmd := exec.Command("lp", "-d"+arg_destination)
  cmd.Stdin = r
  cmd.Stdout = os.Stderr
  cmd.Stderr = os.Stderr
  writer = bufio.NewWriter(w)
  wg.Add(1)
  go func() {
    cmd.Run()
    wg.Done()
  }()
} else {
  writer = bufio.NewWriter(os.Stdout)
}
```
## 参考链接
[开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)