�ù�����Ҫ���޸���Intel��ң�����е�һ��process������ò������tukey's test�������쳣��Ⲣ�����쳣ֵ��
����˼��Ƚϼ򵥣�����bufLength������Ƭ�����������ķ�λQ3�����ķ�λQ1������[q1 - factor * (q3 - q1)��q3 + factor * (q3 - q1)]������ľ���Ϊ���쳣ֵ��
1.�ù����޸�ԭ�������ļ��е�bug��tukey��testʱû�жԴ��������������������ᵼ��q1����q3������������˵���������������x�ᣬɸѡ�쳣ֵ�Ĺ���Ҳ��ʧЧ�ˡ�
2.������https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection��ȡԭ��Intel�Ĵ��룬�����·��$GOPATH/src/github.com/intelsdi-x/��git clone����ֱ��go get
3.������snap��ܵķ�����cd $GOPATH
3.# cd src/github.com/intelsdi-x/snap-plugin-processor-anomalydetection/
4.��snap-plugin-processor-anomalydetection/anomalydetection�µ�����go�ļ��滻Ϊ�ù���anomalydetection/�µ�����go�ļ���
4.cd.. ��make�ɵö������ļ�
5.����snap��ܵķ�����������Ҫ�����·���ṹ���������Ƿ�����/optĿ¼�¡�
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
����snap-plugin-processor-anomalydetection��һ��Ĳ�����԰���֮ǰ�ķ�ʽֱ�ӱ��룬
Ҳ�������عٷ��ṩ�Ķ������ļ����Լ�����ĺô����ڱ����������ϵͳ��һЩ��Ӧ���Ż����Ѿ��ѱ���õĿ�ִ���ļ����ڹ������ˡ�
6.��֮ǰ����õ��Ķ������ļ���������Ӧ·���£�
cp build/linux/x86_64/snap-plugin-processor-anomalydetection /opt/snap_plugin_repos_binary/
7.�л�����Ӧ·��
# cd /opt/snap_plugin_repos_binary/
8.ִ�������ű�
# ./anomaly.sh
9.Լ20���Ӻ󣬵�/tmp/��ȡlog�ļ���windows������
10.����Python����parseSnapAnomaly��ͼ��
Python�����µ�����Python�ļ���
10.1 anomalyDataParseAndContrast.py  ������snapteldһ���ײ���󣬴�linux�ϵ�/tmp�ļ�����ȡ�����Աȵ�log�ļ���windows������
�����ļ���������task Manifest�н��ж��塣���������Ա��ļ������Ի�������https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection�ĶԱ�ͼ��
10.2 tukeyContrastOnlyFromRawData.py  ����tukey's test���쳣���ԭ��ֱ�Ӹ���ԭʼ���ݵ�log�ļ���������Ҫanomaly-processor������log�ļ���
ֱ�ӻ����Ա�ͼ��
ͼƬ�� ͼƬ�����������壬����Go����ζ���Ǹ�������log�ļ���ȡ�õ��ġ�