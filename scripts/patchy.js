var ctime = 0
var stime = 0
var songProg
var playing = false

$(document).ready(function(){
    var loc = window.location
   
    //Initialize Websocket
    conn = new WebSocket("ws://" + loc.host + "/ws");

    //Load now playing
    $.get("/np", function(data) {
        var song = JSON.parse(data)
        console.log(song)
        $("#npArt").attr("src", "/art/" + song["Cover"])
        $("#npSong").text(song["Title"])
        $("#npArtist").text(song["Artist"])
        $("#npAlbum").text(song["Album"])
        $("#songTime").text(secToMin(song["Time"]))
        stime = parseInt(song["Time"])
        $("#curTime").text(secToMin(song["ctime"]))
        ctime = parseInt(song["ctime"])
        var listeners = parseInt(song["listeners"])
        $("#nl").text(listeners)
        
        if(listeners == 1) {
            $("#ltext").text("listener")
        }

        //Load init jplayer
        $("#player-1").jPlayer({
          supplied: "oga",
          preload: "auto",
          volume: 1.0
        });

        console.log("Initialized Player1")

        if(stime != 0 && ctime < stime) {
          playing = true
          
          $("#player-1").jPlayer("setMedia", {
            oga: song["File"]
          });

          $("#player-1").jPlayer("play", ctime)
        }

        $("#songProgress").css("width", (100 * parseInt(song["ctime"])/parseInt(song["Time"])).toString() + "%")

        songProg = window.setInterval(updateSong, 1000);
    });

    //Load queue
    $.get("/curQueue", function(data) {
        var songs = JSON.parse(data)
        $.each(songs, function(index, song) {
            $("#queue").append('<div class="item"><h4><strong>' + song["Title"] + '</strong></h4><p>by <strong>' + song["Artist"] + '</strong></p></div>')
        });

        if($(".item").length > 9) {
            $(".req-button").prop("disabled", true);
        }
    });


    $(".result").slice(20).hide();
    $(".item").slice(10).hide();

    conn.onclose = function(evt) {
        console.log("WS closed")
        
        //conn = new WebSocket("ws://" + loc.host + "/ws");
    }

    conn.onmessage = function(evt) {
        console.log(evt.data)
        var cmd = JSON.parse(evt.data)
        //Update now playing
        if(cmd["cmd"] == "done"){
            endSong(cmd)
        }
        //Pause over, start next song
        if(cmd["cmd"] == "NS"){
            startSong(cmd)
        }

        if(cmd["cmd"] == "queue"){
            updateQueue(cmd)
        }

        //Someone joined
        if(cmd["cmd"] == "ljoin"){
            var listeners = parseInt($("#nl").text()) + 1
            $("#nl").text(listeners)
            if(listeners != 1) {
                $("#ltext").text("listeners")
            }else{
                $("#ltext").text("listener")
            }
        }

        //Someone left
        if(cmd["cmd"] == "lleave"){
            var listeners = parseInt($("#nl").text()) - 1
            $("#nl").text(listeners)
            if(listeners == 1) {
                $("#ltext").text("listener")
            }else{
                $("#ltext").text("listeners")
            }
        }

    }

    $("#search-form").on("submit", function(event) {
      event.preventDefault();

      var dat = $(this).serialize();

      $.get(
        "/search",
        dat,
        function(data) {
          fillSearchRes(data)
        }
      )
    });

    $("#search-clear").click(function() {
      $(".search-results").empty()
    });

    /*
    var form = document.getElementById("ulform");
    var fs = document.getElementById("upload-file");
    form.onsubmit = function(event) {
        var file = fs.files[0]
        var formData = new FormData();
        formData.append("file", file, file.name);
        var request = $.ajax({
            url: "/upload",
            type: "POST",
            async: true,
            data: formData,
            processData: false,
            contentType: false,
            complete: function(jqXHR, status) {
                var resp = jqXHR.responseText
                alert(resp)
            }
        });
    }
    $("#droplink").click(function() {
        $("#upload-file").click()
    });
    $("#upload-button").click(function() {
        form.submit();
        $("#ulform").submit();
    })
    */

    //Initialize Volume stuff
    var slider = $('#slider');
    slider.slider();
        slider.slider({
        range: "min",
        min: 1,
        value: 100,
 
        slide: function(event, ui) {
 
            var value = slider.slider('value') - 2,
                volume = $('.volume');
            var perc = value/100
            $("#player-1").jPlayer("volume", perc)

            if(value <= 5) { 
                volume.css('background-position', '0 0')
            } 
            else if (value <= 25) {
                volume.css('background-position', '0 -25px')
            } 
            else if (value <= 75) {
                volume.css('background-position', '0 -50px')
            } 
            else {
                volume.css('background-position', '0 -75px')
            };
 
        },
    });
});


function fillSearchRes(data) {
  var songs = JSON.parse(data)

  $(".search-results").empty()

  $.each(songs, function(index, song) {
    var title = song["Title"]
    var artist = song["Artist"]
    var album = song["Album"]

    if($.fn.textWidth(song["Title"], "10pt arial") > 400){
      var i = song["Title"].length-1; 
      while($.fn.textWidth(song["Title"].substring(0, i) + "...", "10pt arial") > 400){
        i--;
      }     
      title = song["Title"].substring(0,i) + "..."
    }

    if($.fn.textWidth("by " + song["Artist"], "10pt arial") > 400){
      var i = song["Artist"].length-1; 
      while($.fn.textWidth("by " + song["Artist"].substring(0, i) + "...", "10pt arial") > 400){
        i--;
      }     
      artist = song["Artist"].substring(0,i) + "..."
    }

    if($.fn.textWidth("from " + song["Album"], "10pt arial") > 400){
      var i = song["Album"].length-1; 
      while($.fn.textWidth("from " + song["Album"].substring(0, i) + "...", "10pt arial") > 400){
        i--;
      }     
      album = song["Album"].substring(0,i) + "..."
    }

    $(".search-results").append(
      '<div file="' + song["file"] + '" class="result">' +
        '<img alt="Album art" src="/art/' + song["Cover"] + '">' +
        '<div>' + 
          '<p><strong>' + title + '</strong>' + ' by <strong>' + artist + '</strong></p>' +
          '<p>from <strong>' + album +' </strong></p>' +
        '</div>' +
        '<button class="req-button btn btn-primary btn-block">Request</button>' +
      '</div>'
    )
  }); 

  $(".req-button").click(function() {
    var req = {}
    var block = $(this).parent()

    console.log(block)

    req["cmd"] = "req"
    req["File"] = $(block).attr("file")

    console.log(JSON.stringify(req))

    conn.send(JSON.stringify(req))
  });
}

function updateQueue(song) {
    $("#queue").append('<div class="item"><h4><strong>' + song["Title"] + '</strong></h4><p>by <strong>' + song["Artist"] + '</strong></p></div>')
    if($(".item").length > 9) {
        $(".req-button").prop("disabled", true);
    }else{
        $(".req-button").prop("disabled", false);
    }
}

function endSong(song) {
    $("#player-1").jPlayer("stop")
 
    $("#player-1").jPlayer("clearMedia")

    $("#player-1").jPlayer("setMedia", {
            oga: song["File"]
    });

    console.log("Set Player1 to load " + song["File"] + " in the background") 

    window.clearInterval(songProg)

    $("#curTime").text(secToMin(stime))
    $("#songProgress").css("width", "100%")

    if($(".item").length > 9) {
        $(".req-button").prop("disabled", true);
    }else{
        $(".req-button").prop("disabled", false);
    }
}

function startSong(song) {
    playing = true
    $("#queue").find("div:first").remove();
    $("#npArt").attr("src", song["Cover"])
    $("#npSong").text(song["Title"])
    $("#npArtist").text(song["Artist"])
    $("#npAlbum").text(song["Album"])
    $("#songTime").text(secToMin(song["Time"]))
    $("#curTime").text(secToMin("0"))
    $("#songProgress").css("width", "0%")

    stime = parseInt(song["Time"])
    ctime = 0
    
    $("#player-1").jPlayer("play")

    console.log("Set Player1 to start playing song")

    songProg = window.setInterval(updateSong, 1000);
}

function secToMin(seconds){
    seconds = parseInt(seconds)
    var min = Math.floor(seconds/60);
    var rsecs = seconds - min * 60
    if (rsecs < 10){
        return min.toString() + ":0" + rsecs.toString();
    }else{
        return min.toString() + ":" + rsecs.toString();
    }
}

function updateSong() {
    if(ctime < stime) {
        ctime++

        $("#curTime").text(secToMin(ctime))

        $("#songProgress").css("width", (100 * parseInt(ctime)/parseInt(stime)).toString() + "%")
    } else {
        playing = false

        if($(".item").length > 9) {
            $(".req-button").prop("disabled", true);
        } else {
            $(".req-button").prop("disabled", false);
        }
    }
}

$.fn.textWidth = function(text, font) {
    if (!$.fn.textWidth.fakeEl) $.fn.textWidth.fakeEl = $('<span>').hide().appendTo(document.body);
    $.fn.textWidth.fakeEl.text(text || this.val() || this.text()).css('font', font || this.css('font'));
    return $.fn.textWidth.fakeEl.width();
};
