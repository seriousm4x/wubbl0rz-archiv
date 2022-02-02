function getVodProgress(uuid) {
    if (localStorage.getItem("watched")) {
        let watched = JSON.parse(localStorage.getItem("watched"));
        if (uuid in watched) {
            return watched[uuid];
        }
    }
    return false
}

function load() {
    // show progress bars in grid
    if (localStorage.getItem("watched")) {
        const progressBars = document.querySelectorAll("#watched-progress");
        let watched = JSON.parse(localStorage.getItem("watched"));
        progressBars.forEach(element => {
            if (element.dataset.id in watched) {
                element.classList.remove("d-none");
                let bar = element.children[0];
                bar.ariaValueNow = watched[element.dataset.id];
                bar.style.width = `${Math.round(bar.ariaValueNow*100/bar.ariaValueMax)}%`;
            }
        });
    }
}

if (document.readyState !== "loading") {
    load();
} else {
    document.addEventListener("DOMContentLoaded", load);
}
