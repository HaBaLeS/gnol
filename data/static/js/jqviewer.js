
let pages
let currentPage = 0
let pageUrl
let preloadCursor =0
let preloadNum = 7 //number of images to preload
let singleMaximized = false //default false
let doubleImage = false
let lastPage
let cId

function loadComic(comicId, numPages,cp, di) {
    //TODO rename to init or something
    currentPage = cp
    cId = comicId
    lastPage = numPages-1;
    doubleImage = di
    pages = new Array(pages)
    console.log("Should load:" + comicId + " with " + numPages + " pages");
    pageUrl = "/comics/" + comicId
    loadAndCacheImage(currentPage)

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
    /*if(e.code == "Space"){
        y =window.scrollY
        x = window.scrollX
        document.head
    }*/
}

function handleLoadComplete(e){
    checkPreload();
}

function checkPreload(){
    if (preloadCursor <= currentPage + preloadNum) {
        for(let i=currentPage; i<currentPage + preloadNum; i++ ){
            if(!pages[i] && i <= lastPage){
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
    if (currentPage > lastPage) {
        currentPage = lastPage
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


document.addEventListener('swiped-left', function(e) {
    next();
});

document.addEventListener('swiped-right', function(e) {
    prev();
});


function loadAndCacheImage(pageNum){
    if (pageNum > lastPage) {
        return
    }
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
    if(pages[i] || i > lastPage){
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

    $.ajax({
        url: '/comics/last/'+cId+ '/' + pageNum,
        type: 'PUT',
        success: function(result) {
            console.log("PUT:" + this.url + " success")
        }
    });

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
    document.getElementById("toggleFullScreen").hidden=true
    document.getElementById("toggleExitFullScreen").hidden=false
}


function exitFullScreen(){
    /* When the openFullscreen() function is executed, open the video in fullscreen.
    Note that we must include prefixes for different browsers, as they don't support the requestFullscreen method yet */
    if (document.exitFullscreen) {
        document.exitFullscreen();
    } else if (document.webkitExitFullscreen) {
        document.webkitExitFullscreen();
    } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen();
    } else if (document.msExitFullscreen) {
        document.msExitFullscreen();
    }
    document.getElementById("toggleFullScreen").hidden=false
    document.getElementById("toggleExitFullScreen").hidden=true
}



/*!
 * swiped-events.js - v@version@
 * Pure JavaScript swipe events
 * https://github.com/john-doherty/swiped-events
 * @inspiration https://stackoverflow.com/questions/16348031/disable-scrolling-when-touch-moving-certain-element
 * @author John Doherty <www.johndoherty.info>
 * @license MIT
 */
(function (window, document) {

    'use strict';

    // patch CustomEvent to allow constructor creation (IE/Chrome)
    if (typeof window.CustomEvent !== 'function') {

        window.CustomEvent = function (event, params) {

            params = params || { bubbles: false, cancelable: false, detail: undefined };

            var evt = document.createEvent('CustomEvent');
            evt.initCustomEvent(event, params.bubbles, params.cancelable, params.detail);
            return evt;
        };

        window.CustomEvent.prototype = window.Event.prototype;
    }

    document.addEventListener('touchstart', handleTouchStart, false);
    document.addEventListener('touchmove', handleTouchMove, false);
    document.addEventListener('touchend', handleTouchEnd, false);

    var xDown = null;
    var yDown = null;
    var xDiff = null;
    var yDiff = null;
    var timeDown = null;
    var startEl = null;

    /**
     * Fires swiped event if swipe detected on touchend
     * @param {object} e - browser event object
     * @returns {void}
     */
    function handleTouchEnd(e) {

        // if the user released on a different target, cancel!
        if (startEl !== e.target) return;

        var swipeThreshold = parseInt(getNearestAttribute(startEl, 'data-swipe-threshold', '20'), 10); // default 20px
        var swipeTimeout = parseInt(getNearestAttribute(startEl, 'data-swipe-timeout', '500'), 10);    // default 500ms
        var timeDiff = Date.now() - timeDown;
        var eventType = '';
        var changedTouches = e.changedTouches || e.touches || [];

        if (Math.abs(xDiff) > Math.abs(yDiff)) { // most significant
            if (Math.abs(xDiff) > swipeThreshold && timeDiff < swipeTimeout) {
                if (xDiff > 0) {
                    eventType = 'swiped-left';
                }
                else {
                    eventType = 'swiped-right';
                }
            }
        }
        else if (Math.abs(yDiff) > swipeThreshold && timeDiff < swipeTimeout) {
            if (yDiff > 0) {
                eventType = 'swiped-up';
            }
            else {
                eventType = 'swiped-down';
            }
        }

        if (eventType !== '') {

            var eventData = {
                dir: eventType.replace(/swiped-/, ''),
                touchType: (changedTouches[0] || {}).touchType || 'direct',
                xStart: parseInt(xDown, 10),
                xEnd: parseInt((changedTouches[0] || {}).clientX || -1, 10),
                yStart: parseInt(yDown, 10),
                yEnd: parseInt((changedTouches[0] || {}).clientY || -1, 10)
            };

            // fire `swiped` event event on the element that started the swipe
            startEl.dispatchEvent(new CustomEvent('swiped', { bubbles: true, cancelable: true, detail: eventData }));

            // fire `swiped-dir` event on the element that started the swipe
            startEl.dispatchEvent(new CustomEvent(eventType, { bubbles: true, cancelable: true, detail: eventData }));
        }

        // reset values
        xDown = null;
        yDown = null;
        timeDown = null;
    }

    /**
     * Records current location on touchstart event
     * @param {object} e - browser event object
     * @returns {void}
     */
    function handleTouchStart(e) {

        // if the element has data-swipe-ignore="true" we stop listening for swipe events
        if (e.target.getAttribute('data-swipe-ignore') === 'true') return;

        startEl = e.target;

        timeDown = Date.now();
        xDown = e.touches[0].clientX;
        yDown = e.touches[0].clientY;
        xDiff = 0;
        yDiff = 0;
    }

    /**
     * Records location diff in px on touchmove event
     * @param {object} e - browser event object
     * @returns {void}
     */
    function handleTouchMove(e) {

        if (!xDown || !yDown) return;

        var xUp = e.touches[0].clientX;
        var yUp = e.touches[0].clientY;

        xDiff = xDown - xUp;
        yDiff = yDown - yUp;
    }

    /**
     * Gets attribute off HTML element or nearest parent
     * @param {object} el - HTML element to retrieve attribute from
     * @param {string} attributeName - name of the attribute
     * @param {any} defaultValue - default value to return if no match found
     * @returns {any} attribute value or defaultValue
     */
    function getNearestAttribute(el, attributeName, defaultValue) {

        // walk up the dom tree looking for attributeName
        while (el && el !== document.documentElement) {

            var attributeValue = el.getAttribute(attributeName);

            if (attributeValue) {
                return attributeValue;
            }

            el = el.parentNode;
        }

        return defaultValue;
    }

}(window, document));
