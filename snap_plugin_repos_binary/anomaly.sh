set -e
set -u
set -o pipefail

[ -f "/tmp/published_anomalydetection_vm_used.log" ] && rm -rf /tmp/published_anomalydetection_vm_used.log
[ -f "/tmp/published_anomaly_contrast_vm_used.log" ] && rm -rf /tmp/published_anomaly_contrast_vm_used.log

plugin_dir="$(find / -name snap-plugin-collector-psutil | grep -v tmp|grep -v src)"
plugin__dir="$(cd "$(dirname "$plugin_dir")" && pwd)"

service snap-telemetry restart
cd $plugin__dir 
sleep 5s
snaptel plugin load snap-plugin-collector-psutil && snaptel plugin load snap-plugin-processor-anomalydetection && snaptel plugin load snap-plugin-processor-tag && snaptel plugin load snap-plugin-publisher-file

manifest_dir="$(find / -name task_manifest_file)"

snaptel task create -t $manifest_dir"/psutil-anomalydetection-file-vm-used.json" &
snaptel task create -t $manifest_dir"/psutil-tag-file-vm-used.json" &
wait