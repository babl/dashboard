#!/bin/bash

#-------------------------------------------------------
# VARS
#-------------------------------------------------------
LOG_MESSAGE=""

selected_day=`date "+%Y-%m-%d"`
if [ ! -z "$1" ] ; then
  selected_day=$1
fi

cluster=$cluster'production%5C.babl%5C.sh%7C';
cluster=$cluster'sandbox%5C.babl%5C.sh';

modules=$modules'loyalist%5C%2Fdesigner-uploads%7C ';
modules=$modules'loyalist%5C%2Fprep-prints%7C ';
modules=$modules'loyalist%5C%2Fprint-images%7C ';
modules=$modules'loyalist%5C%2Fproduct-image%7C ';
modules=$modules'loyalist%5C%2Fproduct-image-uploader%7C ';
modules=$modules'loyalist%5C%2Fteam-banner%7C ';
modules=$modules'loyalist%5C%2Fstatus';

from='%27'${selected_day}'T00%3A00%3A00.000000000Z%27';
to='%27'${selected_day}'T23%3A59%3A59.999999999Z%27';

#-------------------------------------------------------
# Babl
#-------------------------------------------------------
function getData
{
  m=$1
  # BABL REQUEST TOTAL (code='req-enqueued')
  total_enqueued=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F++AND+%22module%22+%3D\~+%2F%5E\(${m}\)%24%2F+AND+code%3D%27req-enqueued%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to}\&db\=logs_qa);
  [ "$total_enqueued" == '{"results":[{}]}' ] && total_enqueued=0 || total_enqueued=$(echo $total_enqueued | jq -c '.results[]? | .series[]? | .values[]? | .[1:2] | .[]');

  re='^[0-9]+$'
  if ! [[ $total_enqueued =~ $re ]] ; then
     total_enqueued=0
  fi

  # BABL REQUESTS SUCCESS (code='completed' AND status='SUCCESS')
  total_success=$(curl --silent -u babl:qWzBwrWYcvUxiRtLvNuH7uEgLHiqMrVwRthUYndHWBLkc4hFzH https://influxdb.admin.babl.sh:18086/query\?q\=SELECT+count\(duration_ms\)+FROM+logs_qa..kafka_consumer_logs_qa+WHERE+%22cluster%22+%3D\~+%2F%5E\(${cluster}\)%24%2F++AND+%22module%22+%3D\~+%2F%5E\(${m}\)%24%2F+AND+code%3D%27completed%27+AND+status%3D%27SUCCESS%27+AND+time+%3E%3D+${from}+AND+time+%3C%3D+${to}\&db\=logs_qa);
  [ "$total_success" == '{"results":[{}]}' ] && total_success=0 || total_success=$(echo $total_success | jq -c '.results[]? | .series[]? | .values[]? | .[1:2] | .[]');
  
  re='^[0-9]+$'
  if ! [[ $total_success =~ $re ]] ; then
     total_success=0
  fi
  
  total_error=$(expr $total_enqueued - $total_success);

  payload=$payload'{"value": '$total_enqueued',"error": '$total_error',"l":'$total_success' ,"u":'$total_enqueued'}'
  echo $payload  
}
for module in $modules; do
  data=$(getData $module)
  out=$out'{"module":"'$module'","data":'$data'},'
done
out="${out%?}"
echo '{"date": "'$selected_day'","data":['$out']}'


# payload=$payload'{"date": "'$selected_day'","value": '$total_enqueued',"error": '$total_error',"l":'$total_success' ,"u":'$total_enqueued'}'
  # echo $payload  