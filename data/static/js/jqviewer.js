
let pages;
let currentPage = 0;
let pageUrl;
let preloadCursor =0;
let preloadNum = 7; //number of images to preload
let verticalFit = true; //default false

function loadComic(comicId,numPages) {
    //TODO rename to init or something

    pages = new Array(pages)
    console.log("Should load:" + comicId + " with " + numPages + " pages");
    pageUrl = "/read2/" + comicId
    loadAndCacheImage(0)

    document.addEventListener('keydown', handleKeyboardInput);
    document.addEventListener('cmxLoadComplete', handleLoadComplete);
    window.addEventListener('resize', reportWindowSize);

    let fsb = document.getElementById("orientationStyleBtn")
    fsb.onclick =function () {
        verticalFit =!verticalFit;
        replaceImage(currentPage);
    }

    let mcb = document.getElementById("closeHelp");
    mcb.onclick =function () {
        let mdh = document.getElementById("modalHelp")
        mdh.classList.remove("is-active")
    }


    var hammertime = new Hammer(document);
    hammertime.on('swipe', function(ev) {
        console.log(ev);
        alert(ev);
    });
}

function showHelp(){
    let mdh = document.getElementById("modalHelp")
    mdh.classList.add("is-active")
}

function reportWindowSize(e){
    replaceImage(currentPage)
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

function handleLoadComplete(e){
    checkPreload();
}

function checkPreload(){
    if (preloadCursor <= currentPage + preloadNum) {
        for(i=currentPage; i<currentPage + preloadNum; i++ ){
            if(!pages[i]){
                preloadImage(i);
                return
            }
        }
    }
}

function next(){
    currentPage++;
    if(pages[currentPage]){
        replaceImage(currentPage)
    } else {
        loadAndCacheImage(currentPage)
    }
    document.body.scrollTop = 0; // For Safari
    document.documentElement.scrollTop = 0; // For Chrome, Firefox, IE and Opera
    checkPreload();
}

function prev(){
    currentPage--;
    if(currentPage <0 ){
        currentPage = 0;
    }
    if(pages[currentPage]){
        replaceImage(currentPage)
    } else {
        loadAndCacheImage(currentPage)
    }
}

function loadAndCacheImage(pageNum){
      let pgu = pageUrl + "/" +pageNum
      jQuery.ajax({
        url: pgu,
        cache:false,
        xhr:function(){// Seems like the only way to get access to the xhr object
            let xhr = new XMLHttpRequest();
            xhr.responseType= 'blob'
            return xhr;
        },
        success: function(data){
            let url = window.URL || window.webkitURL;
            pages[pageNum] = url.createObjectURL(data);
            document.dispatchEvent(new Event("cmxLoadComplete"));
            replaceImage(pageNum)
        },
        error:function(){
            alert("Error loading image")
        }
    });
}

function preloadImage(i){
    let pgu = pageUrl + "/" +i
    if(pages[i]){
        return
    }
    jQuery.ajax({
        url: pgu,
        cache:false,
        xhr:function(){// Seems like the only way to get access to the xhr object
            let xhr = new XMLHttpRequest();
            xhr.responseType= 'blob'
            return xhr;
        },
        success: function(data){
            let url = window.URL || window.webkitURL;
            pages[i] = url.createObjectURL(data);
            document.dispatchEvent(new Event("cmxLoadComplete"));
        },
        error:function(){
            alert("Error loading image")
        }
    });
}

function replaceImage(pageNum){
    let oldImg = document.getElementById('cv');
    let newImage = document.createElement("img");
    oldImg.parentNode.replaceChild(newImage, oldImg)
    newImage.src = pages[pageNum];
    newImage.id = "cv"
    newImage.classList.add("has-ratio")

   //calc W,H and then the Resized Versions

    if(verticalFit){
        newImage.height = window.innerHeight-8;
        document.getElementById("orientationIcon").className = "fa fa-arrows-alt-h my-float";
    } else {
        newImage.width = window.innerWidth-16;
        document.getElementById("orientationIcon").className = "fa fa-arrows-alt-v my-float";
    }
}

function enableFullScreen(){
    /* When the openFullscreen() function is executed, open the video in fullscreen.
    Note that we must include prefixes for different browsers, as they don't support the requestFullscreen method yet */
    let elem = document.documentElement;
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
