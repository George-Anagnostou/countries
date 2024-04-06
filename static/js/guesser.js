import { fetchSingleCountry, generateCountryCard, generateCountryCardPartial } from "./utils.js";

const targetCountryContainer = document.getElementById("country-container");
const form = document.getElementById("guess-country-form");
const resultContainer = document.getElementById("result-container");
const next = document.getElementById("next-guess");

let singleCountryData;

function getCountryGuess(formData) {
    const countryGuess = formData.get("country-guess");
    return countryGuess;
}

function isCorrectGuess(countryGuess, countryTarget) {
    return countryGuess.toLowerCase() === countryTarget.toLowerCase()
}

function displayTargetCountry(targetCountryContainer, targetCountry) {
    targetCountryContainer.innerHTML = generateCountryCardPartial(targetCountry);
}

document.addEventListener("DOMContentLoaded", async (e) => {
    e.preventDefault();
    singleCountryData = await fetchSingleCountry();
    displayTargetCountry(targetCountryContainer, singleCountryData);
});

form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const countryGuess = getCountryGuess(formData);
    const result = isCorrectGuess(countryGuess, singleCountryData.name.common)
                    ? "Correct!"
                    : "Wrong!"
    resultContainer.innerHTML = `<p>${result}</p>`
    targetCountryContainer.innerHTML = generateCountryCard(singleCountryData);
    next.focus();
});

next.addEventListener("click", async (e) => {
    e.preventDefault();
    resultContainer.innerHTML = "";
    singleCountryData = await fetchSingleCountry();
    displayTargetCountry(targetCountryContainer, singleCountryData);
    form.elements["country-guess"].value = "";
    form.elements["country-guess"].select();
});
