import { fetchSingleCountry, generateCountryCard } from "./utils.js";
import { getGuess, isCorrectGuess, displayTargetCountry } from "./guessUtils.js";

const targetCountryContainer = document.getElementById("country-container");
const form = document.getElementById("guess-country-form");
const resultContainer = document.getElementById("result-container");
const next = document.getElementById("next-guess");

let singleCountryData;

document.addEventListener("DOMContentLoaded", async (e) => {
    e.preventDefault();
    singleCountryData = await fetchSingleCountry();
    displayTargetCountry(targetCountryContainer, singleCountryData);
});

form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const countryGuess = getGuess(formData);
    const result = isCorrectGuess(countryGuess, singleCountryData.name.common)
                    ? "Correct!"
                    : `Wrong! The correct answer is ${singleCountryData.name.common}`
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
