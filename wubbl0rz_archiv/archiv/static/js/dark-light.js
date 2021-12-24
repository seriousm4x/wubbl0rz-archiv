function load() {
    const toggleBtn = document.querySelector(".toggle-dark");
    const preferesDark = window.matchMedia("(prefers-color-scheme: dark)");

    function toggleDarkMode(state) {
        document.documentElement.classList.toggle("dark-mode", state);
        localStorage.setItem("dark-mode", state);
        currentModeState = state;
        toggleBtn.classList.toggle("active", state);
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
}

if (document.readyState !== "loading") {
    load();
} else {
    document.addEventListener("DOMContentLoaded", load);
}
