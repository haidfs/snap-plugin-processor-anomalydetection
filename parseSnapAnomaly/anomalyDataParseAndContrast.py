import matplotlib.pyplot as plt


def get_data_from_anomaly_detection_log(raw_data, filepath, time_map_dyct):
    # i = 0
    data_lyst2 = [0 for x in range(len(raw_data))]
    # with open(r"D:\snapTelemetry\published_anomalydetection.log") as file:
    with open(filepath) as file:
        for ten_metric in file.readlines():
            if ten_metric == "null\n":
                continue
            for metric in eval(ten_metric):
                minute_second = ''.join(metric["timestamp"].split('.')[0].split(':')[-2:])
                anomaly_index = time_map_dyct[minute_second]
                data_lyst2[anomaly_index] = metric["data"]
                # if time_map_dyct[minute_second] == i:
                #     data_lyst2.append(metric["data"])
                #     i += 1
                # else:
                #     while i < time_map_dyct[minute_second]:
                #         data_lyst2.append(0)
                #         i += 1
                #     data_lyst2.append(metric["data"])
                #     i += 1
    # print(len(data_lyst2))
    bins = [i for i in range(len(data_lyst2))]
    # X = [i for i in range(10)]
    # Y = data_lyst[:10]
    data = data_lyst2
    return bins, data


def get_data_from_log(filepath):
    metric_dyct = {}
    time_map_dyct = {}
    data_lyst = []
    i = 0
    # with open(r"D:\snapTelemetry\published_anomaly_contrast.log") as file:
    with open(filepath) as file:
        for metric in file.readlines():
            metric_raw_dyct = dict(eval(metric)[0])
            metric_dyct[i] = metric_raw_dyct["data"]
            data_lyst.append(metric_raw_dyct["data"])
            minute_second = ''.join(metric_raw_dyct["timestamp"].split('.')[0].split(':')[-2:])
            time_map_dyct[minute_second] = i
            i += 1
    bins = [i for i in range(len(data_lyst))]
    data = data_lyst
    return bins, data, time_map_dyct


def plot_bar(bins, data, xlabel, ylabel, title):
    plt.bar(bins, data, 0.5, color='green')

    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.title(title)


def plot_contrast_bar(bins1, data1, bins2, data2, metric, suffix=''):
    plt.figure(figsize=(30, 20))
    plt.subplot(212)
    plot_bar(bins=bins1, data=data1, xlabel="Continuous Time Series/s", ylabel=metric,
             title=metric + " Without AnomalyDetection")
    plt.subplot(211)
    plot_bar(bins=bins2, data=data2, xlabel="Continuous Time Series/s", ylabel=metric,
             title=metric + " With AnomalyDetection")
    plt.savefig(metric + "AnomalyDetectionContrast" + "_" + suffix + ".png")
    plt.show()


def calc_compression(raw_data, anomaly_data):
    raw_number = len(raw_data)
    reduced_number = len([1 for i in anomaly_data if i != 0])
    # print(len(anomaly_data))
    print("\nThe raw data number: %d" % raw_number)
    print("The anomaly data number: %d" % reduced_number)
    print("The compression ratio: %.3f" % (raw_number / reduced_number))


def main():
    bins1, data1, time_map = get_data_from_log(
        r"D:\snapTelemetry\data\sorted_myself\published_anomaly_contrast_vm_used.log")
    bins2, data2 = get_data_from_anomaly_detection_log(data1,
                                                       r"D:\snapTelemetry\data\sorted_myself\published_anomalydetection_vm_used.log",
                                                       time_map)
    plot_contrast_bar(bins1, data1, bins2, data2, "VmUsed", "SortedGo_3.1")
    calc_compression(data1, data2)


if __name__ == '__main__':
    main()