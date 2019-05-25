import pandas as pd
from parseSnapAnomaly.anomalyDataParseAndContrast import *


def parse_values(values, q1, q3, factor):
    fence1 = q1 - factor * (q3 - q1)
    fence2 = q3 + factor * (q3 - q1)
    # print(fence1, fence2)
    # assert fence1 <= fence2
    outliers = []
    value = 0
    for (i, v) in enumerate(values):
        if v < fence1 or v > fence2:
            value = value + v
            outliers.append(i)
    if len(outliers) != 0:
        return value / float(len(outliers)), outliers

    return 0.0, outliers


def get_outliers(values, factor, sort_index):
    length = len(values)
    l = sort_index
    if length % 2 == 0:
        # print("q1 index:{} q3 index:{}".format(l[length // 4],l[3 * length // 4]))
        q1 = values[l[length // 4]]
        q3 = values[l[3 * length // 4]]
        # assert q1 <= q3
        print("q1: " + str(q1) + " q3: " + str(q3))
        return parse_values(values, q1, q3, factor)
    i = values[l[length // 4]]
    j = values[l[length // 4 + 1]]
    q1 = i + (j - i) * 0.25
    i = values[l[3 * length // 4 - 1]]
    j = values[l[3 * length // 4]]
    q3 = i + (j - i) * 0.75
    # assert q1 <= q3
    return parse_values(values, q1, q3, factor)


def anomaly_detection(values, factor, bufLength, sort_index):
    if len(values) != bufLength:
        raise RuntimeError("the length of values is not bufLength!")

    _, anomaly = get_outliers(values, factor, sort_index)
    print("anomaly_index:{}".format(anomaly))
    return anomaly


def csv_file_to_df(csv_path):
    csv_file = csv_path
    csv_data = pd.read_csv(csv_file, low_memory=False)
    csv_df = pd.DataFrame(csv_data)
    return csv_df


# 获取异常值元素在切片中的索引时，根据sort_flag来决定是否需要对传入的数据列表进行排序
def calc_tukey_get_anomaly_data(raw_total_value, buf_length, sort_flag=True):
    i = 0
    # anomaly_values=[]
    anomaly_values = [0 for l in range(len(raw_total_value))]
    while i + buf_length <= len(raw_total_value):
        values = raw_total_value[i:i + buf_length]
        if sort_flag == True:
            sort_index = sorted(range(len(values)), key=lambda k: values[k])
        else:
            sort_index = [x for x in range(len(values))]
        print("sort_index:{}".format(sort_index))
        # print("sorted values:{}".format(sorted(values)))
        # values = sorted(values)
        length = len(values)
        anomaly_index = anomaly_detection(values, bufLength=length, factor=3.1, sort_index=sort_index)

        for n in anomaly_index:
            anomaly_values[i + n] = values[n]
        i += buf_length
    print(anomaly_values)
    return anomaly_values


if __name__ == '__main__':
    # values = [1, 100, 4, 5, 6, 7, 9, 1, 20, 2]
    # iperf_100_df = csv_file_to_df(
    #     r'C:\Users\***\ReceiveFile\data_serverlef4_700.csv')
    # cpu_value = list(iperf_100_df['CpuUsage'])
    #
    # bins1 = bins2 = [i for i in range(len(cpu_value))]
    # anomaly_data = calc_tukey_get_anomaly_data(cpu_value, 10)

    # plot_contrast_bar(bins1, cpu_value, bins2, anomaly_data, "CpuUsage", "DCN")
    # calc_compression(cpu_value, anomaly_data)

    bins1, data1, time_map = get_data_from_log(r"D:\snapTelemetry\data\published_anomaly_contrast_vm_used_3.1.log")
    data2 = calc_tukey_get_anomaly_data(data1, 9, sort_flag=False)
    bins2 = bins1
    plot_contrast_bar(bins1, data1, bins2, data2, "VmUsed", "TestPython_unsorted")
    calc_compression(data1, data2)