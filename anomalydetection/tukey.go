package anomalydetection

import (
	"errors"
	"strconv"

	"github.com/intelsdi-x/snap/control/plugin"
)

//估计这个地方是把buf数据传入，然后将异常值的索引加到outliers里面
func parseValues(values []float64, q1 float64, q3 float64, factor float64) (float64, []int) {
	var (
		outliers []int
		value    float64
	)
	fence1 := q1 - factor*(q3-q1)
	fence2 := q3 + factor*(q3-q1)
	for i, v := range values {
		if v < fence1 || v > fence2 {
			value = value + v
			outliers = append(outliers, i)

		}
	}

	if len(outliers) != 0 {
		//返回异常值的平均值和异常数组的长度
		return value / float64(len(outliers)), outliers
	}
	return 0.0, outliers
}

func getOutliers(values []float64, factor float64,s_idx []int) (float64, []int) {

	l := len(values)
	//如果传入的buf为偶数长度，计算对应的四分位值
	if l%2 == 0 {
		q1 := values[s_idx[l/4]]
		q3 := values[s_idx[3*l/4]]

		err := lessThan(q1, q3)
		handleErr(err)


		return parseValues(values, q1, q3, factor)

	}
	i := values[s_idx[l/4]]
	j := values[s_idx[l/4+1]]
	q1 := i + (j-i)*0.25
	i = values[s_idx[3*l/4-1]]
	j = values[s_idx[3*l/4]]
	q3 := i + (j-i)*0.75

	err := lessThan(q1, q3)
	handleErr(err)

	return parseValues(values, q1, q3, factor)

}
func lessThan(a float64, b float64) error {
	if a > b {
		return errors.New("q1>q3")
	} else {
		return nil
	}

}

func interfaceToFloat(face interface{}) (float64, error) {
	var (
		ret float64
		err error
	)
	//i.(type)
	switch val := face.(type) {
	case string:
		ret, err = strconv.ParseFloat(val, 64)
	case int:
		ret = float64(val)
	case int16:
		ret = float64(val)
	case int32:
		ret = float64(val)
	case int64:
		ret = float64(val)
	case uint:
		ret = float64(val)
	case uint16:
		ret = float64(val)
	case uint32:
		ret = float64(val)
	case uint64:
		ret = float64(val)
	case float32:
		ret = float64(val)
	case float64:
		ret = val

	default:
		err = errors.New("unsupported type")
	}
	return ret, err
}

func unpackData(values []plugin.MetricType) ([]float64, error) {
	//float类型数组，初始化为空
	metrics := []float64{}
	for _, v := range values {
		//.(dataType)表示的是类型断言
		metrics = append(metrics, v.Data_.(float64))
	}
	return metrics, nil
}

func dataToFloat(m plugin.MetricType) (plugin.MetricType, error) {
	var err error
	m.Data_, err = interfaceToFloat(m.Data_)
	return m, err
}