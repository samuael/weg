function PlayAudio( host  , resource  , containerid  ){
    var htmls = `<div class="container-audio">
    <audio id="audio" controls  loop autoplay>
    </audio>
</div> 
<div class="audio-section container-audio">
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
    <div class="colum1">
        <div class="raw"></div>
    </div>
</div> 
<hr>

<div style="padding: 10px;margin: 10px;">
    <h4> <strong>  ` +resource.Title + ` </strong> </h4>
    <hr>
    <div class="col-5" >
        <img  style="width:50px; height:50px ; border-radius: 50%;"   src="{{.HOST}}{{.ActiveResource.UploaderImage}}" alt="">
        <strong><small> `+ resource.Uploadedby+`  </small></strong>
        <br>
        <strong><small> Uploaded at : ` + resource.Day +"/"+ resource.Month +"/"+resource.Year +"/"+resource.Hour +":"+resource.Minute+ ` </small> </strong>
        <br>
        <br>
    </div>
    <hr>
    <p>
        <strong>
           <small> `+resource.Description +`</small>
        </strong>
    </p>
</div>  
    <script>
        var audio = document.getElementById('audio');
       if (Hls.isSupported()) {
           var hls = new Hls();   //hls.js
           hls.loadSource(   `+host+resource.HLSDirectory+`" );
           hls.attachMedia(audio);
           hls.on(Hls.Events.MANIFEST_PARSED, function () {
               video.play();
           });
       } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
           video.src = host + audioPath;
           video.addEventListener('loadedmetadata', function () {
               video.play();
           });
       }
   </script>`;
   var mainelement = document.getElementById( containerid );
   mainelement.innerHTML= htmls;
}