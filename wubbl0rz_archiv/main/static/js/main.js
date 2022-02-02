function autocomplete(inp, arr) {
    if (arr == undefined) {
        return
    }
    if (arr.length == null) {
        removeActive(x);
        return
    }
    var currentFocus;
    inp.addEventListener("input", function (e) {
        var a, b, i, val = this.value;
        closeAllLists();
        if (!val) {
            return false;
        }
        currentFocus = -1;
        a = document.createElement("DIV");
        a.setAttribute("id", this.id + "autocomplete-list");
        a.setAttribute("class", "autocomplete-items");
        this.parentNode.appendChild(a);
        for (i = 0; i < arr.length; i++) {
            if (document.querySelectorAll("#searchautocomplete-list > div").length < 10) {
                if (arr[i].toUpperCase().includes(val.toUpperCase())) {
                    b = document.createElement("DIV");
                    b.innerHTML = arr[i].substr(0, val.length);
                    b.innerHTML += arr[i].substr(val.length);
                    b.innerHTML += "<input type='hidden' value='" + arr[i] + "'>";
                    b.addEventListener("click", function (e) {
                        inp.value = this.getElementsByTagName("input")[0].value;
                        closeAllLists();
                    });
                    a.appendChild(b);
                }
            }
        }
        if (document.querySelectorAll("#searchautocomplete-list > div").length == 0) {
            document.getElementById("searchautocomplete-list").style.opacity = 0;
        } else {
            document.getElementById("searchautocomplete-list").style.opacity = 1;
        }
    });
    inp.addEventListener("keydown", function (e) {
        var x = document.getElementById(this.id + "autocomplete-list");
        if (x) x = x.getElementsByTagName("div");
        if (e.keyCode == 40) {
            currentFocus++;
            addActive(x);
        } else if (e.keyCode == 38) {
            currentFocus--;
            addActive(x);
        } else if (e.keyCode == 13) {
            if (currentFocus > -1) {
                if (x) x[currentFocus].click();
            }
        }
    });

    function addActive(x) {
        if (!x) return false;
        removeActive(x);
        if (currentFocus >= x.length) currentFocus = 0;
        if (currentFocus < 0) currentFocus = (x.length - 1);
        x[currentFocus].classList.add("autocomplete-active");
    }

    function removeActive(x) {
        for (var i = 0; i < x.length; i++) {
            x[i].classList.remove("autocomplete-active");
        }
    }

    function closeAllLists(elmnt) {
        var x = document.getElementsByClassName("autocomplete-items");
        for (var i = 0; i < x.length; i++) {
            if (elmnt != x[i] && elmnt != inp) {
                x[i].parentNode.removeChild(x[i]);
            }
        }
    }
    document.addEventListener("click", function (e) {
        closeAllLists(e.target);
    });
}

function cssvar(name) {
    return getComputedStyle(document.documentElement).getPropertyValue(name);
}

function drawDoughnut(id, labels, dataset) {
    const ctx = document.getElementById(id).getContext('2d');
    const myChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: dataset,
                backgroundColor: ["#4E79A7", "#F28E2B", "#E15759", "#76B7B2", "#59A14F", "#EDC948", "#B07AA1", "#FF9DA7", "#9C755F", "#BAB0AC"],
            }]
        },
        options: {
            plugins: {
                legend: {
                    position: 'top',
                    labels: {
                        color: cssvar('--main-color'),
                    }
                }
            }
        }
    });
}

function drawBar(id, labels, dataset) {
    const ctx = document.getElementById(id).getContext('2d');
    const myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            color: cssvar('--main-color'),
            datasets: [{
                data: dataset,
                backgroundColor: ["#4E79A7", "#F28E2B", "#E15759", "#76B7B2", "#59A14F", "#EDC948", "#B07AA1", "#FF9DA7", "#9C755F", "#BAB0AC"],
            }]
        },
        options: {
            plugins: {
                legend: {
                    display: false
                },
            },
            scales: {
                y: {
                    ticks: {
                        color: cssvar('--main-color'),
                    },
                    grid: {
                        color: cssvar("--stats-grid"),
                    }
                },
                x: {
                    ticks: {
                        color: cssvar('--main-color'),
                    },
                    grid: {
                        color: cssvar("--stats-grid"),
                    }
                }
            }
        }
    });
};

function drawLine(id, labels, dataset) {
    const ctx = document.getElementById(id).getContext('2d');
    const myChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            color: cssvar('--main-color'),
            datasets: [{
                data: dataset,
                borderColor: "#F28E2B",
                tension: 0.4
            }]
        },
        options: {
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                },
            },
            interaction: {
                intersect: false,
            },
            scales: {
                y: {
                    ticks: {
                        color: cssvar('--main-color'),
                    },
                    grid: {
                        color: cssvar("--stats-grid"),
                    }
                },
                x: {
                    ticks: {
                        color: cssvar('--main-color'),
                    },
                    title: {
                        display: true,
                        text: 'Uhrzeit',
                        color: cssvar('--main-color'),
                        font: {
                            size: 16,
                        }
                    },
                    grid: {
                        color: cssvar("--stats-grid"),
                    }
                }
            }
        }
    })
};

function showPrev(elem) {
    document.getElementById(elem.id + "-sm-avif").srcset = `/media/${elem.dataset.type}/` + elem.id + "-preview.webp";
    document.getElementById(elem.id + "-md-avif").srcset = `/media/${elem.dataset.type}/` + elem.id + "-preview.webp";
}

function hidePrev(elem) {
    document.getElementById(elem.id + "-sm-avif").srcset = `/media/${elem.dataset.type}/` + elem.id + '-sm.avif';
    document.getElementById(elem.id + "-md-avif").srcset = `/media/${elem.dataset.type}/` + elem.id + '-md.avif';
}

function load() {
    // dark light toggle
    const toggleBtn = document.querySelector(".toggle-dark");
    const preferesDark = window.matchMedia("(prefers-color-scheme: dark)");

    function toggleDarkMode(state) {
        document.documentElement.classList.toggle("dark-mode", state);
        localStorage.setItem("dark-mode", state);
        currentModeState = state;
        toggleBtn.classList.toggle("active", state);

        // change chart colors
        if (typeof Chart !== "undefined") {
            Chart.helpers.each(Chart.instances, function (instance) {
                if (instance.config._config.type == "doughnut") {
                    instance.options = {
                        plugins: {
                            legend: {
                                labels: {
                                    color: cssvar('--main-color'),
                                }
                            }
                        }
                    }
                    instance.update();
                } else {
                    instance.options.scales = {
                        y: {
                            ticks: {
                                color: cssvar("--main-color"),
                            },
                            grid: {
                                color: cssvar("--stats-grid"),
                            }
                        },
                        x: {
                            ticks: {
                                color: cssvar("--main-color"),
                            },
                            title: {
                                color: cssvar("--main-color"),
                            },
                            grid: {
                                color: cssvar("--stats-grid"),
                            }
                        }
                    }
                    instance.update();
                }
            })
        }
    }

    if (localStorage.getItem("dark-mode") === null) {
        toggleDarkMode(preferesDark.matches);
    } else {
        toggleDarkMode(localStorage.getItem("dark-mode") == "true");
    }

    preferesDark.addEventListener("change", e => {
        toggleDarkMode(e.matches);
    })

    toggleBtn.addEventListener("click", () => {
        currentModeState = !currentModeState;
        toggleDarkMode(currentModeState);
    });

    // thumbnail hover
    document.querySelectorAll("picture, .has-preview").forEach(function (elem) {
        elem.addEventListener("touchstart", showPrev(elem), false);
        elem.addEventListener("touchend", hidePrev(elem), false);
    })

    // hotkey events
    document.onkeydown = function (event) {
        let searchActive = document.getElementById("search") == document.activeElement;
        if (event.key === "h" && !searchActive) {
            // toggle hotkey modal
            let hotkeyModal = document.getElementById("hotkeyModal");
            var modal = bootstrap.Modal.getInstance(hotkeyModal);
            if (!modal) {
                var modal = new bootstrap.Modal(hotkeyModal);
            }
            modal.toggle();
            return false;
        } else if (event.key === "/" && !searchActive) {
            // highlight search
            document.querySelector("#search").select();
            return false;
        } else if (event.key === "t" && !searchActive) {
            // toggle theme
            currentModeState = !currentModeState;
            toggleDarkMode(currentModeState);
            return false;
        } else if (typeof player !== "undefined" && event.key === "f" && !searchActive) {
            // toggle video fullscreen
            if (player.isFullscreen()) {
                player.exitFullscreen();
            } else {
                player.requestFullscreen();
            }
            return false;
        } else if (typeof player !== "undefined" && event.key === " " && !searchActive) {
            // toggle video play pause
            if (player.paused()) {
                player.play();
            } else {
                player.pause();
            }
            return false;
        } else if (typeof player !== "undefined" && event.key === "ArrowRight" && !searchActive) {
            // seek video forward 30s
            player.currentTime(player.currentTime() + 30)
            return false;
        } else if (typeof player !== "undefined" && event.key === "ArrowLeft" && !searchActive) {
            // seek video back 10s
            player.currentTime(player.currentTime() - 10)
            return false;
        } else if (typeof player !== "undefined" && event.key === "ArrowUp" && !searchActive) {
            // add 5% to volume
            player.volume(Math.round(player.volume() * 100) / 100 + 0.05)
            return false;
        } else if (typeof player !== "undefined" && event.key === "ArrowDown" && !searchActive) {
            // substract 5% from volume
            player.volume(Math.round(player.volume() * 100) / 100 - 0.05)
            return false;
        } else if (typeof player !== "undefined" && event.key === "m" && !searchActive) {
            // mute player
            if (player.muted()) {
                player.muted(false)
            } else {
                player.muted(true)
            }
            return false;
        }
    }
}

if (document.readyState !== "loading") {
    load();
} else {
    document.addEventListener("DOMContentLoaded", load);
}
