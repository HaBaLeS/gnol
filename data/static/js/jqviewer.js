
let pages
let currentPage = 0
let pageUrl
let preloadCursor =0
let preloadNum = 7 //number of images to preload
let singleMaximized = false //default false
let doubleImage = false

function loadComic(comicId, numPages, di) {
    //TODO rename to init or something
    doubleImage = di
    pages = new Array(pages)
    console.log("Should load:" + comicId + " with " + numPages + " pages");
    pageUrl = "/read2/" + comicId
    loadAndCacheImage(0)

    document.addEventListener('keydown', handleKeyboardInput);
    document.addEventListener('cmxLoadComplete', handleLoadComplete);
    window.addEventListener('resize', reportWindowSize);

    document.getElementById("orientationStyleBtn").onclick =function () {
        singleMaximized =!singleMaximized;
        enableSinglePage();
    }

    /*
    let mcb = document.getElementById("closeHelp");
    mcb.onclick =function () {
        let mdh = document.getElementById("modalHelp")
        mdh.classList.remove("is-active")
    }
    */

    var hammertime = new Hammer(document);
    hammertime.on('swipe', function(ev) {
        console.log(ev);
        alert(ev);
    });

    enableSinglePage();
}

function showHelp(){
    let mdh = document.getElementById("modalHelp")
    mdh.classList.add("is-active")
}

function reportWindowSize(e){
    updateScreen(currentPage)
}

function enableDualPage(){
    document.getElementById("dualPage").hidden=false
    document.getElementById("singleFullPage").hidden=true
    document.getElementById("singleMaximized").hidden=true
    document.getElementById("toggleToSingle").hidden=false
    document.getElementById("toggleToDual").hidden=true
    document.getElementById("toggleOrientation").hidden=true
    doubleImage = true
    updateScreen(currentPage)
}

function enableSinglePage(){
    doubleImage=false
    document.getElementById("dualPage").hidden=true
    document.getElementById("toggleOrientation").hidden=false
    if(singleMaximized){
        document.getElementById("singleFullPage").hidden=true
        document.getElementById("singleMaximized").hidden=false
        document.getElementById("orientationIcon").classList.remove("fa-arrows-alt-h")
        document.getElementById("orientationIcon").classList.add("fa-arrows-alt-v")
    } else {
        document.getElementById("singleFullPage").hidden=false
        document.getElementById("singleMaximized").hidden=true
        document.getElementById("orientationIcon").classList.add("fa-arrows-alt-h")
        document.getElementById("orientationIcon").classList.remove("fa-arrows-alt-v")
    }

    document.getElementById("toggleToSingle").hidden=true
    document.getElementById("toggleToDual").hidden=false
    updateScreen(currentPage);
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
    if(doubleImage) {
        currentPage++
    }
    if(pages[currentPage]){
        updateScreen(currentPage)
    } else {
        loadAndCacheImage(currentPage)
    }
    document.body.scrollTop = 0; // For Safari
    document.documentElement.scrollTop = 0; // For Chrome, Firefox, IE and Opera
    checkPreload();
}

function prev(){
    currentPage--;
    if(doubleImage) {
        currentPage--
    }
    if(currentPage <0 ){
        currentPage = 0;
    }
    if(pages[currentPage]){
        updateScreen(currentPage)
    } else {
        loadAndCacheImage(currentPage)
    }
}

// -- stolen from SO https://stackoverflow.com/questions/2264072/detect-a-finger-swipe-through-javascript-on-the-iphone-and-android
var xDown = null, yDown = null, xUp = null, yUp = null;
document.addEventListener('touchstart', touchstart, false);
document.addEventListener('touchmove', touchmove, false);
document.addEventListener('touchend', touchend, false);
function touchstart(evt) { const firstTouch = (evt.touches || evt.originalEvent.touches)[0]; xDown = firstTouch.clientX; yDown = firstTouch.clientY; }
function touchmove(evt) { if (!xDown || !yDown ) return; xUp = evt.touches[0].clientX; yUp = evt.touches[0].clientY; }
function touchend(evt) {
    var xDiff = xUp - xDown, yDiff = yUp - yDown;
    if ((Math.abs(xDiff) > Math.abs(yDiff)) && (Math.abs(xDiff) > 0.20 * document.body.clientWidth)) {
        if (xDiff < 0)
            next()
        else
            prev()
    }
    xDown = null, yDown = null;
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
            updateScreen(pageNum)
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

function updateScreen(pageNum){
    if(doubleImage){
        document.getElementById('dpl').src = pages[pageNum]
        document.getElementById('dpr').src = pages[pageNum+1]
    }else {
        if(singleMaximized){
            document.getElementById("smi").src = pages[pageNum]
        } else {
            document.getElementById("fpi").src = pages[pageNum]
        }
    }

}

/*
function updateScreen(pageNum){
    let oldImg = document.getElementById('cv');
    let newImage = document.createElement("img");
    oldImg.parentNode.replaceChild(newImage, oldImg)
    newImage.src = pages[pageNum];
    newImage.id = "cv"
   //calc W,H and then the Resized Versions

    if(verticalFit){
        newImage.height = window.innerHeight-8;
        document.getElementById("orientationIcon").className = "fa fa-arrows-alt-h my-float";
    } else {
        newImage.width = window.innerWidth-16;
        document.getElementById("orientationIcon").className = "fa fa-arrows-alt-v my-float";
    }
}*/

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
