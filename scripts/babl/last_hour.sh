#!/bin/bash

#-------------------------------------------------------
# VARS
#-------------------------------------------------------

selected_day=`date "+%Y-%m-%d"`
if [ ! -z "$1" ] ; then
  selected_day=$1
fi

selected_hour=`date "+%H"`
if [ ! -z "$2" ] ; then
  selected_hour=$(printf "%02d" $2)
fi

selected_minute=`date "+%M"`
if [ ! -z "$3" ] ; then
  selected_minute=$(printf "%02d" $3)
fi

cluster=$cluster'production%5C.babl%5C.sh%7C';
cluster=$cluster'sandbox%5C.babl%5C.sh';

minus_hour=$(printf "%02d" $(expr $selected_hour - 1));
[ "$minus_hour" == "-1" ] && minus_hour="23" || minus_hour=$minus_hour;

minus_minute=$(printf "%02d" $(expr $selected_minute - 1));
[ "$minus_minute" == "-1" ] && minus_minute="23" || minus_minute=$minus_minute;

from='%27'${selected_day}'T'${minus_hour}'%3A'${selected_minute}'%3A00.000000000Z%27';
to='%27'${selected_day}'T'${selected_hour}'%3A'${selected_minute}'%3A00.000000000Z%27';
to_enq='%27'${selected_day}'T'${selected_hour}'%3A'${minus_minute}'%3A00.000000000Z%27';
now='"'${selected_day}'T'${selected_hour}'%3A'${selected_minute}'"';

#-------------------------------------------------------
# Babl
#-------------------------------------------------------
# BABL REQUEST TOTAL (code='req-enqueued')

total_enqueued=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F+AND+code%3D%27req-enqueued%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to_enq}\&db\=logs_qa);
[ "$total_enqueued" == '{"results":[{}]}' ] && total_enqueued=0 || total_enqueued=$(echo $total_enqueued | jq -c '.results[]? | .series[]? | .values[]? | .[1:2] | .[]');

re='^[0-9]+$'
if ! [[ $total_enqueued =~ $re ]] ; then
   total_enqueued=0
fi

#echo 'total:'$total_enqueued
# BABL REQUESTS SUCCESS (code='completed' AND status='SUCCESS')
total_success=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F+AND+code%3D%27completed%27+AND+status%3D%27SUCCESS%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to}\&db\=logs_qa);
[ "$total_success" == '{"results":[{}]}' ] && total_success=0 || total_success=$(echo $total_success | jq -c '.results[]? | .series[]? | .values[]? | .[1:2] | .[]');

re='^[0-9]+$'
if ! [[ $total_success =~ $re ]] ; then
   total_success=0
fi

#total_error=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F++AND+%22module%22+%3D\~+%2F%5E\(${modules}\)%24%2F+AND+code%3D%27completed%27+AND+status%3C%3E%27SUCCESS%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to}\&db\=logs_qa | jq -c '.results[]? | .series[]? | .values[]? | .[1:2] | .[]');
#echo 'error:'$total_error

total_error=$(expr $total_enqueued - $total_success);
payload=$payload'{"date": '$now',"total": '$total_enqueued',"error": '$total_error'}'
echo $payload

