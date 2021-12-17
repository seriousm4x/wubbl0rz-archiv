String.prototype.toHHMMSS = function () {
    var sec_num = parseInt(this, 10);
    var hours   = Math.floor(sec_num / 3600);
    var minutes = Math.floor((sec_num - (hours * 3600)) / 60);
    var seconds = sec_num - (hours * 3600) - (minutes * 60);

    if (hours   < 10) {hours   = hours;}
    if (minutes < 10) {minutes = "0"+minutes;}
    if (seconds < 10) {seconds = "0"+seconds;}
    return hours+':'+minutes+':'+seconds;
}

function load() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const shareLabel = document.querySelector("#share-at-text");
    let shareURL = document.querySelector(".share-input");
    const shareCheckbox = document.querySelector("#share-checkbox");
    const btnClipboard = document.querySelector("#btn-copy-clipboard");

    if (urlParams.get("t")) {
        player.currentTime(urlParams.get("t"));
        shareURL.value = window.location.href;
        shareLabel.innerHTML = "Teilen bei " + urlParams.get("t").toHHMMSS();
        shareCheckbox.checked = true;
    } else {
        shareURL.value = location.protocol + '//' + location.host + location.pathname;
        shareLabel.innerHTML = "Teilen bei 0:00:00";
        shareCheckbox.checked = false;
    }

    // prevent dropdown close on click
    document.querySelector('.dropdown-menu').addEventListener("click", (e) => {
        e.stopPropagation();
    })

    // set share time on player time update
    player.on('timeupdate', () => {
        let timeRounded = Math.round(player.currentTime())
        if (shareCheckbox.checked) {
            shareURL.value = location.protocol + '//' + location.host + location.pathname + "?t=" + timeRounded;
        }
        shareLabel.innerHTML = "Teilen bei " + String(timeRounded).toHHMMSS();
    });

    // set share url on radio click
    shareCheckbox.addEventListener("click", () => {
        let shareURL = document.querySelector(".share-input");
        if (shareURL.value.includes("?t=")) {
            shareURL.value = location.protocol + '//' + location.host + location.pathname;
        } else {
            shareURL.value = location.protocol + '//' + location.host + location.pathname + "?t=" + Math.round(player.currentTime());
        }
    })

    // set copy clipboard button
    btnClipboard.addEventListener("click", () => {
        let shareURL = document.querySelector(".share-input");
        const iconClip = document.querySelector("#icon-clipboard");
        const iconClipChecked = document.querySelector("#icon-clipboard-checked");
        navigator.clipboard.writeText(shareURL.value);
        iconClip.style.display = "none";
        iconClipChecked.style.display = "block";
        setTimeout(() => {
            iconClip.style.display = "block";
            iconClipChecked.style.display = "none";
        }, 1000);
    })
}

window.addEventListener("DOMContentLoaded", load);
