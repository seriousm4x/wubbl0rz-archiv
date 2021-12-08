function load() {
    let darkModeState = false;

    const toggle = document.querySelector(".toggle");

    // MediaQueryList object
    const useDark = window.matchMedia("(prefers-color-scheme: dark)");

    // Toggles the "dark-mode" class
    function toggleDarkMode(state) {
        document.documentElement.classList.toggle("dark-mode", state);
        darkModeState = state;
        toggle.checked = !state;
    }

    // Sets localStorage state
    function setDarkModeLocalStorage(state) {
        localStorage.setItem("dark-mode", state);
    }

    // Initial setting
    toggleDarkMode(localStorage.getItem("dark-mode") == "true");

    // Listen for changes in the OS settings.
    // Note: the arrow function shorthand works only in modern browsers, 
    // for older browsers define the function using the function keyword.
    useDark.addEventListener("change", e => {
        toggleDarkMode(e.matches);
    })

    // Toggles the "dark-mode" class on click and sets localStorage state
    toggle.addEventListener("click", () => {
        darkModeState = !darkModeState;
        toggleDarkMode(darkModeState);
        setDarkModeLocalStorage(darkModeState);
    });

}

window.addEventListener("DOMContentLoaded", load);