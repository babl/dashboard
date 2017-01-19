var theme = 'light';
var moduleUser;


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

          //#TODO, EXCHANGE BROADCAST WITH GROUP CHANNELS!!!!
          if(data["Group"] == moduleUser) {
            data = data.Last
            
            //#error calculation can be negative...
            if (data.error < 0){
              data.error = 0
            }
            var noCache="?"+(new Date()).getTime();
            var lastHourPath = "../../data/"+moduleUser+"/hour_max.json"+noCache
            d3.json(lastHourPath,function(max){

              if(data.total == 0 && data.error ==0){
                $('#no-progress').show()
                var errorRate = 0
              }else{
                $('#no-progress').hide()
                var errorRate = (((data.error)/data.total)*100).toFixed(2)
              }

              var successRateMax = (((data.total-data.error)/max.total)*100).toFixed(2)
              var errorRateMax = (((data.error)/max.total)*100).toFixed(2)
              $('#success-bar').attr("style","width:"+successRateMax+"%")
              $('#error-bar').attr("style","width:"+errorRateMax+"%")

              //info
              $("#cur-date").text(moment(decodeURIComponent(data.date)).calendar())
              $("#total-req").text(data.total)
              $("#error-rate").text(errorRate+'%')
              $("#max").text(max.total)
              $("#max-date").text(moment(decodeURIComponent(max.date)).fromNow())  
            })
          }
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
            websocket = new WebSocket("ws://" + host + "/ws/"+moduleUser);
            websocket.onopen = function(evt) {
              console.log("Websocket connection opened");
              // $('#ws').text("ws connected, refresh every min")
            }
            websocket.onclose = function(evt) {
              console.log("Websocket connection closed");
              // $('#ws').text("ws closed... no data refresh!")
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
          console.log()
          $.post('/lasthour','user='+moduleUser,function(data){
            data = JSON.parse(data);
            updateLastHour(data)
          })
        }

        $(document).ready(function() {
          moduleUser = location.pathname.substring(1)
          connectWebsocket();
          update()
        });
      }())
