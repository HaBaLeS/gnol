
var pages;
var currentPage = 0;
var  pageUrl;
function loadComic(comicId,numPages) {
    //TODO rename to initi or something

    pages = new Array(pages)
    console.log("Should load:" + comicId + " with " + numPages + " pages");
    pageUrl = "http://192.168.1.248:6969/read2/" + comicId
    loadAndCacheImage(pageUrl, 0)

    document.addEventListener('keydown', handleKeyboardInput);

    initModal();
}

function initModal(){

// Get the modal
    var modal = document.getElementById("myModal");

// Get the button that opens the modal
    var btn = document.getElementById("myBtn");

// Get the <span> element that closes the modal
    var span = document.getElementsByClassName("close")[0];

// When the user clicks the button, open the modal
    btn.onclick = function() {
        modal.style.display = "block";
    }

// When the user clicks on <span> (x), close the modal
    span.onclick = function() {
        modal.style.display = "none";
    }

// When the user clicks anywhere outside of the modal, close it
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    }
}

function handleKeyboardInput(e){
    if(e.code == "ArrowRight"){
        next();
    }
    if(e.code == "ArrowLeft"){
        prev();
    }
    if(e.code == "KeyF"){
        enableFullScreen();
    }
}

function next(){
    currentPage++;
    if(pages[currentPage]){
        replaceImage(currentPage)
    } else {
        loadAndCacheImage(pageUrl, currentPage)
    }
    document.body.scrollTop = 0; // For Safari
    document.documentElement.scrollTop = 0; // For Chrome, Firefox, IE and Opera
}

function prev(){
    currentPage--;
    if(currentPage <0 ){
        currentPage = 0;
    }
    if(pages[currentPage]){
        replaceImage(currentPage)
    } else {
        loadAndCacheImage(pageUrl, currentPage)
    }
}

function loadAndCacheImage(pageUrl, pageNum){
      pageUrl = pageUrl + "/" +pageNum
      jQuery.ajax({
        url: pageUrl,
        cache:false,
        xhr:function(){// Seems like the only way to get access to the xhr object
            var xhr = new XMLHttpRequest();
            xhr.responseType= 'blob'
            return xhr;
        },
        success: function(data){
            var url = window.URL || window.webkitURL;
            pages[pageNum] = url.createObjectURL(data);
            replaceImage(pageNum)
        },
        error:function(){
            alert("Error loading image")
        }
    });
}

function replaceImage(pageNum){
    var img = document.getElementById('cv');
    img.src = pages[pageNum]
    //img.height = window.innerHeight;

    vfit = false;
    if(vfit){
        //fit height
        img.height = window.visualViewport.height-8;
    } else {
        //fit width
        img.width = window.innerWidth-16;
    }




    //var info = document.getElementById("info")
    //info.innerText="Hallo " + pageNum + " ->" + img.width + ":" + img.height;
}

function enableFullScreen(){
    /* When the openFullscreen() function is executed, open the video in fullscreen.
    Note that we must include prefixes for different browsers, as they don't support the requestFullscreen method yet */
    var elem = document.documentElement; //document.getElementById("view");
    if (elem.requestFullscreen) {
        elem.requestFullscreen();
    } else if (elem.mozRequestFullScreen) { /* Firefox */
        elem.mozRequestFullScreen();
    } else if (elem.webkitRequestFullscreen) { /* Chrome, Safari and Opera */
        elem.webkitRequestFullscreen();
    } else if (elem.msRequestFullscreen) { /* IE/Edge */
        elem.msRequestFullscreen();
    }
}
