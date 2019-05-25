package anomalydetection

import (
	"bytes"
	"encoding/gob"
	log "github.com/sirupsen/logrus"
	"sort"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	name                = "anomalydetection"
	version             = 1
	pluginType          = plugin.ProcessorPluginType
	defaultBufferLength = 30
	defaultFactor       = 3.0
)

// Buffer struct stores []plugin.MetricType for specific namespace
//根据bufLength形成的数组
type Buffer struct {
	Metrics []plugin.MetricType
}

// BufferMetric struct, stores all Buffers
//key是join后的namespace
type BufferMetric struct {
	Buffer map[string]*Buffer
}

// Meta returns a plugin meta data
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// NewAnomalydetectionProcessor creates new processor
func NewAnomalydetectionProcessor() *anomalyDetectionProcessor {
	buffer := make(map[string]*Buffer, 30)
	return &anomalyDetectionProcessor{
		BufferMetric: BufferMetric{
			Buffer: buffer,
		},
	}
}

type anomalyDetectionProcessor struct {
	BufferMetric BufferMetric
}

func (p *anomalyDetectionProcessor) addToBuffer(m plugin.MetricType, logger *log.Logger) error {
	ns := m.Namespace().String()
	//把Metric里面的data项转换为float
	m, err := dataToFloat(m)
	if err != nil {
		return err
	}
	if _, ok := p.BufferMetric.Buffer[ns]; ok {
		p.BufferMetric.Buffer[ns].Metrics = append(p.BufferMetric.Buffer[ns].Metrics, m)

	} else {
		vMet := []plugin.MetricType{m}
		p.BufferMetric.Buffer[ns] = &Buffer{
			Metrics: vMet,
		}
	}
	logger.Debug("Buffer length: ", len(p.BufferMetric.Buffer[ns].Metrics))
	return nil

}

func (p *anomalyDetectionProcessor) clearBuffer(ns string) {

	vMet := []plugin.MetricType{}
	p.BufferMetric.Buffer[ns] = &Buffer{
		Metrics: vMet,
	}
}

func (p *anomalyDetectionProcessor) getBuffer(ns string) []plugin.MetricType {

	return p.BufferMetric.Buffer[ns].Metrics
}

func (p *anomalyDetectionProcessor) getBufferLength(ns string) int {

	return len(p.BufferMetric.Buffer[ns].Metrics)
}

type Slice struct {
	sort.Float64Slice
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

//go语言的函数封装不如Python丰富，返回排序后的切片索引，Python里面一句话就够了sort_index = sorted(range(len(values)), key=lambda k: values[k])
//go需要自己写结构体
func NewSlice(n []float64) *Slice {
	s := &Slice{Float64Slice: sort.Float64Slice(n), idx: make([]int, len(n))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func (p *anomalyDetectionProcessor) calculateTukeyMethod(m plugin.MetricType, factor float64, logger *log.Logger) ([]plugin.MetricType, error) {

	ns := m.Namespace().String()
	//获取bufLength长的metric切片
	metrics := p.getBuffer(ns)
	//从metrics里面把数据解开形成新的数组
	values, err := unpackData(metrics)
	//deepcopy
	values_bak := make([]float64, len(values))
	copy(values_bak, values)
	if err != nil {
		return nil, err
	}
	s := NewSlice(values)
	sort.Sort(s)

	_, outliersIndex := getOutliers(values_bak, factor, s.idx)

	//这里的用处不大，似乎log没有返回到snapteld 守护进程的log
	//log.WithFields(log.Fields{
	//	"outliersIndex": outliersIndex,
	//}).Warning("outliersIndex:")

	ret := []plugin.MetricType{}
	//直接获取到了异常值，那么null值的索引如何记录？与gdb编码相关？？
	for _, v := range outliersIndex {
		ret = append(ret, metrics[v])
	}
	return ret, nil

}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

//获取配置文件，检查参数项配置是否正确
func (p *anomalyDetectionProcessor) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()
	r1, err := cpolicy.NewIntegerRule("BufLength", true, defaultBufferLength)
	handleErr(err)
	r1.Description = "Buffer Length for tukey method "
	config.Add(r1)
	r2, err := cpolicy.NewFloatRule("Factor", false, defaultFactor)
	handleErr(err)
	r2.Description = "Sensitivity Factor for scaling"
	config.Add(r2)
	cp.Add([]string{""}, config)
	return cp, nil
}

func (p *anomalyDetectionProcessor) Process(contentType string, content []byte, config map[string]ctypes.ConfigValue) (string, []byte, error) {
	var (
		metrics, metricsTemp []plugin.MetricType
		bufferLength         int
		factor               float64
	)

	logger := log.New()
	logger.Level = log.DebugLevel
	logger.Debug("anomalyDetection Processor started")

	//0表示默认值。这里表示先断言后取值
	if config["BufLength"].(ctypes.ConfigValueInt).Value > 0 {
		bufferLength = config["BufLength"].(ctypes.ConfigValueInt).Value

	}
	if config["Factor"].(ctypes.ConfigValueFloat).Value > 0 {
		factor = config["Factor"].(ctypes.ConfigValueFloat).Value

	}
	//Decodes the content into MetricType

	//从collector收集的来的数据，之前应该经过了序列化。这里是执行反序列化的操作？应该是这样的，这是gob解码的方式，将从collector接收来的数据解码成metric
	//buf1:=bytes.NewBufferString("hello")
	//buf2:=bytes.NewBuffer([]byte("hello"))
	//buf3:=bytes.NewBuffer([]byte{"h","e","l","l","o"})
	//以上三者等效。bytes.NewBuffer表示创建一个字节流的Buffer
	//??？？疑问，content传进来的是多少个buffer的字节流？？
	dec := gob.NewDecoder(bytes.NewBuffer(content))
	if err := dec.Decode(&metrics); err != nil {
		logger.Printf("Error decoding: error=%v content=%v", err, content)
		return "", nil, err
	}

	for _, m := range metrics {

		ns := m.Namespace().String()
		if _, ok := p.BufferMetric.Buffer[ns]; ok {
			if p.getBufferLength(ns) == bufferLength-1 {
				mVal, err := p.calculateTukeyMethod(m, factor, logger)
				if err != nil {
					return "", nil, err
				}
				//将Tukey得到的异常值切片append到metricsTemp，for结束后将metricsTemp通过gob编码传出
				metricsTemp = append(metricsTemp, mVal...)
				p.clearBuffer(ns)

			} else {
				p.addToBuffer(m, logger)
			}
		} else {
			p.addToBuffer(m, logger)
		}

	}
	//下面这一串代表的是gob编码
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(metricsTemp)
	return contentType, buf.Bytes(), nil
}