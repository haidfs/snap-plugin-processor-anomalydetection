该工程主要是修改了Intel的遥测框架中的一个process插件，该插件根据tukey's test来进行异常检测并返回异常值。

基本思想比较简单，就是bufLength长的切片排序后计算上四分位Q3和下四分位Q1，处在[q1 - factor * (q3 - q1)，q3 + factor * (q3 - q1)]区间外的就认为是异常值。

1.该工程修改原本代码文件中的bug，tukey’test时没有对传入的数组进行排序，这样会导致q1大于q3，这样上面所说的条件会包含整个x轴，筛选异常值的功能也就失效了。

2.在这里https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection获取原本Intel的代码，在这个路径$GOPATH/src/github.com/intelsdi-x/下git clone。或直接go get

4.在运行snap框架的服务器cd $GOPATH

5.# cd src/github.com/intelsdi-x/snap-plugin-processor-anomalydetection/

6.将snap-plugin-processor-anomalydetection/anomalydetection下的两个go文件替换为该工程anomalydetection/下的两个go文件。

7.cd.. 再make可得二进制文件

8.运行snap框架的服务器首先需要定义的路径结构，例如我是放在了/opt目录下。
/opt/snap_plugin_repos_binary/
|-- anomaly.sh
|-- bin_bak
|   |-- snap-plugin-processor-anomalydetection
|   `-- snap-plugin-processor-anomalydetection-errorsq1q3
|-- snap-plugin-collector-cpu
|-- snap-plugin-collector-logs
|-- snap-plugin-collector-psutil
|-- snap-plugin-processor-anomalydetection
|-- snap-plugin-processor-change-detector
|-- snap-plugin-processor-logs-regexp
|-- snap-plugin-processor-statistics
|-- snap-plugin-processor-tag
|-- snap-plugin-processor-tags-filter
|-- snap-plugin-publisher-file
|-- task_manifest_file
    |-- psutil-anomalydetection-file.json
    |-- psutil-anomalydetection-file-vm-used.json
    |-- psutil-statistics-file.json
    |-- psutil-tag-file.json
    |-- psutil-tag-file-vm-used.json
    |-- psutil-tags-filter-file.json
    |-- task-config.json
    |-- task.json
    `-- task.yaml
其中snap-plugin-processor-anomalydetection这一类的插件可以按照之前的方式直接编译，
也可以下载官方提供的二进制文件，自己编译的好处在于编译器会根据系统做一些相应的优化。已经把编译好的可执行文件放在工程中了。

9.把之前编译得到的二进制文件拷贝到对应路径下：
cp build/linux/x86_64/snap-plugin-processor-anomalydetection /opt/snap_plugin_repos_binary/

10.切换到对应路径
# cd /opt/snap_plugin_repos_binary/

11.执行启动脚本
# ./anomaly.sh

12.约20分钟后，到/tmp/拉取log文件到windows本机。

13.根据Python工程parseSnapAnomaly绘图。
Python工程下的两个Python文件：

13.1 anomalyDataParseAndContrast.py  在运行snapteld一整套插件后，从linux上的/tmp文件夹拉取两个对比的log文件到windows本机，
具体文件名在两个task Manifest中进行定义。根据两个对比文件，可以画出近似https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection的对比图。

13.2 tukeyContrastOnlyFromRawData.py  根据tukey's test的异常检测原理，直接根据原始数据的log文件，而不需要anomaly-processor处理后的log文件，
直接画出对比图。
图片： 图片命名即本身含义，带有Go的意味着是根据两个log文件获取得到的。
