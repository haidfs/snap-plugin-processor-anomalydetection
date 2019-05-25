package main

//������Linux�ϵĵ��ԣ���Ҫ��ʹ��gdb����dlv�����Ƕ�����������̫��Ϥ��goҲ�Ǹ����֣����԰Ѻ��ĵĴ����߼���ȡ�������һ���򵥵�windows���Դ��룬�Ա���goland��ֱ�ӵ��ԡ�
import (
	"errors"
	"fmt"
	"sort"
)

type Slice struct {
	sort.Float64Slice
	idx []int
}

//ʵ��sort�ӿ������swap�������ڽ���ʵ����ƬԪ��ֵ��ʱ�����ƬԪ�ص�����Ҳ������Ӧ�Ľ�����
func (s Slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n []float64) *Slice {
	s := &Slice{Float64Slice: sort.Float64Slice(n), idx: make([]int, len(n))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}
func lessThan(a float64, b float64) error {
	if a > b {
		return errors.New("q1>q3")
	} else {
		return nil
	}

}
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
		//�����쳣ֵ��ƽ��ֵ���쳣����ĳ���
		return value / float64(len(outliers)), outliers
	}
	return 0.0, outliers
}

func getOutliers(values []float64, factor float64, s_idx []int) (float64, []int) {

	l := len(values)
	//��������bufΪż�����ȣ������Ӧ���ķ�λֵ
	if l%2 == 0 {
		q1 := values[s_idx[l/4]]
		q3 := values[s_idx[3*l/4]]

		fmt.Println("q1_index: ", s_idx[l/4], "q1: ", q1, "q3:", q3, "q3_index: ", s_idx[3*l/4])

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

	fmt.Println("q1_index: ", s_idx[l/4], "\tq1: ", q1, "\n\r q3_index: ", s_idx[3*l/4],
		"\tq3: ", q3)

	err := lessThan(q1, q3)
	handleErr(err)

	return parseValues(values, q1, q3, factor)

}
func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	values := []float64{1.2, 25.3, 3.1, 5.33, 4.9, 1.2, 100.8, 11.0}
	fmt.Println(values)
	//deepcopy
	values_bak := make([]float64, len(values))
	copy(values_bak, values)
	s := NewSlice(values)
	sort.Sort(s)
	fmt.Println("values_bak_:\t", values_bak)
	fmt.Println(s.Float64Slice, s.idx)
	//s.idx����Ƭ����󷵻��������������

	_, outlier := getOutliers(values_bak, 1.5, s.idx)
	fmt.Println(outlier)
}
