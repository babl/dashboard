var theme = 'light';

(function() {
  'use strict';

    // replace all SVG images with inline SVG
    // http://stackoverflow.com/questions/11978995/how-to-change-color-of-svg
    // -image-using-css-jquery-svg-image-replacement
    $('img.svg').each(function() {
      var $img = jQuery(this);
      var imgID = $img.attr('id');
      var imgClass = $img.attr('class');
      var imgURL = $img.attr('src');

      $.get(imgURL, function(data) {
        // Get the SVG tag, ignore the rest
        var $svg = jQuery(data).find('svg');

        // Add replaced image's ID to the new SVG
        if (typeof imgID !== 'undefined') {
          $svg = $svg.attr('id', imgID);
        }
        // Add replaced image's classes to the new SVG
        if (typeof imgClass !== 'undefined') {
          $svg = $svg.attr('class', imgClass + ' replaced-svg');
        }

        // Remove any invalid XML tags as per http://validator.w3.org
        $svg = $svg.removeAttr('xmlns:a');

        // Replace image with new SVG
        $img.replaceWith($svg);

      }, 'xml');
    });    
  })();


  (function(){
    'use strict';

      var websocket;
      
        function updateLastHour(data){

          d3.json('../../data/hour_max.json',function(max){
          //bar
          var successRate = (((data.total-data.error)/max.value)*100).toFixed(2)
          var errorRate = (((data.error)/max.value)*100).toFixed(2)
          $('#success-bar').attr("style","width:"+successRate+"%")
          $('#error-bar').attr("style","width:"+errorRate+"%")

          if(data.value == 0 && data.error ==0){
            $('#no-progress').show()
          }else{
            $('#no-progress').hide()
          }
          //info
          console.log('max',max)
          $("#cur-date").text(moment(decodeURIComponent(data.date)).calendar())
          $("#total-req").text(data.total)
          $("#error-rate").text((((data.error)/data.total)*100).toFixed(2))
          $("#max").text(max.value)
          $("#max-date").text(moment(decodeURIComponent(max.date)).fromNow())
        })
      }

      function connectWebsocket() {

        function resetWs() {
          websocket.onmessage = function () {};
          websocket.onclose = function () {};
          websocket.onopen = function () {};
          websocket.close();
          websocket = null;
        }

        if (websocket != null) {
          resetWs();
          console.log('web sockets reseted!')
        }

        if (window["WebSocket"]) {
          var host = window.location.host;
          websocket = new WebSocket("ws://" + host + "/ws");
          websocket.onopen = function(evt) {
            console.log("Websocket connection opened");
            $('#ws').text("ws connected, refresh every min")
          }
          websocket.onclose = function(evt) {
            console.log("Websocket connection closed");
            $('#ws').text("ws closed... no data refresh!")
          resetWs();
        };
        websocket.onmessage = function(evt) {
          var data = JSON.parse(evt.data)
          updateLastHour(data)
        };
      } else {
        tableDataAddInfo("Your browser does not support WebSockets");
      }
    }

    function update(){
      $.get('/lasthour',function(data){
        console.log('data',data)
        data = JSON.parse(data);
        updateLastHour(data)
      })
    }

    $(document).ready(function() {
      console.log("document loaded");
      
      connectWebsocket();
      update()
      $('#refresh').on('click', update)

      d3.json('../../data/daily.json',function(data){
        
        var newData = []
        data.forEach(function(d){
          var suc = d.value-d.error
          var i = {
            "Date":d.date,
            "Total":d.value,
            "Success":suc,
            "Error": d.error,
            "Sucess Rate": (((suc)/d.value)*100).toFixed(2)+'%',
            "Error Rate":(((d.error)/d.value)*100).toFixed(2)+'%'
          }
          newData.push(i)
        })
        var jsonHtmlTable = ConvertJsonToTable(newData, 'table', 'table', 'Download');

        $('#table').html(jsonHtmlTable)
        
      })
    });
}())
