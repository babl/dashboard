#!/bin/sh

#-------------------------------------------------------
# VARS
#-------------------------------------------------------
LOG_MESSAGE=""

selected_day=`date "+%Y-%m-%d"`
if [ ! -z "$1" ] ; then
  selected_day=$1
fi

selected_hour=`date "+%H"`
if [ ! -z "$2" ] ; then
  selected_hour=$2
fi

selected_minute=`date "+%M"`
if [ ! -z "$3" ] ; then
  selected_minute=$3
fi

cluster=$cluster'production%5C.babl%5C.sh%7C';
cluster=$cluster'sandbox%5C.babl%5C.sh';

modules=$modules'loyalist%5C%2Fdesigner-uploads%7C';
modules=$modules'loyalist%5C%2Fprep-prints%7C';
modules=$modules'loyalist%5C%2Fprint-images%7C';
modules=$modules'loyalist%5C%2Fproduct-image%7C';
modules=$modules'loyalist%5C%2Fproduct-image-uploader%7C';
modules=$modules'loyalist%5C%2Fteam-banner%7C';
modules=$modules'loyalist%5C%2Fstatus';

from='%27'${selected_day}'T'$(expr $selected_hour - 1)'%3A'${selected_minute}'%3A00.000000000Z%27';
to='%27'${selected_day}'T'${selected_hour}'%3A'${selected_minute}'%3A00.000000000Z%27';
to_enq='%27'${selected_day}'T'${selected_hour}'%3A'$(expr $selected_minute - 1)'%3A00.000000000Z%27';

echo 'from' $selected_day ':' $from 'to' $to

# from='%27'${selected_day}'T12%3A00%3A00.000000000Z%27';
# to='%27'${selected_day}'T12%3A59%3A59.999999999Z%27';
# to_less_1m='%27'${selected_day}'T12%3A58%3A59.999999999Z%27';
#-------------------------------------------------------
# Babl
#-------------------------------------------------------
# BABL REQUEST TOTAL (code='req-enqueued')

total_enqueued=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F++AND+%22module%22+%3D\~+%2F%5E\(${modules}\)%24%2F+AND+code%3D%27req-enqueued%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to_enq}\&db\=logs_qa | jq -c '.results[]? | .series[]? | .values[] | .[1:2] | .[]');
echo 'total:'$total_enqueued
# BABL REQUESTS SUCCESS (code='completed' AND status='SUCCESS')

total_success=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F++AND+%22module%22+%3D\~+%2F%5E\(${modules}\)%24%2F+AND+code%3D%27completed%27+AND+status%3D%27SUCCESS%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to}\&db\=logs_qa | jq -c '.results[]? | .series[]? | .values[]? | .[1:2] | .[]');
echo 'suc:'$total_success

total_error=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F++AND+%22module%22+%3D\~+%2F%5E\(${modules}\)%24%2F+AND+code%3D%27completed%27+AND+status%3C%3E%27SUCCESS%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to}\&db\=logs_qa | jq -c '.results[]? | .series[]? | .values[]? | .[1:2] | .[]');
echo 'error:'$total_error

total_error=$(expr $total_enqueued - $total_success);
percent_success=$(echo "scale=5;"$total_success"/"$total_enqueued"*100.00" | bc -l);
percent_error=$(echo 'scale=2;' "100.00-"$percent_success | bc -l);
payload=$payload'{"date": "'$selected_day'","value": '$total_enqueued',"error": '$total_error',"l":'$total_success' ,"u":'$total_enqueued'}\n'
echo $payload

