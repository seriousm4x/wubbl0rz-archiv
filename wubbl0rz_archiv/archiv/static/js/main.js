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
                    }
                },
                x: {
                    ticks: {
                        color: cssvar('--main-color'),
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
                    title: {
                        display: true,
                        text: 'Anzahl Streams',
                        color: cssvar('--main-color'),
                        font: {
                            size: 16,
                        }
                    },
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
                }
            }
        }
    })
};

function showPrev(elem) {
    document.getElementById(elem.id + "-avif").srcset = "/media/" + elem.id + "-preview.webp";
    document.getElementById(elem.id + "-jpg").src = "/media/" + elem.id + "-preview.webp";
}

function hidePrev(elem) {
    document.getElementById(elem.id + "-avif").srcset = "/media/" + elem.id + '-sm.avif';
    document.getElementById(elem.id + "-jpg").src = "/media/" + elem.id + '-sm.jpg';
}

function load() {
    document.querySelectorAll("picture, .has-preview").forEach(function (elem) {
        elem.addEventListener("touchstart", showPrev(elem), false);
        elem.addEventListener("touchend", hidePrev(elem), false);
    })
}

window.addEventListener("DOMContentLoaded", load);
