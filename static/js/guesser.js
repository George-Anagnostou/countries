import { fetchSingleCountry } from "./utils.js";

const targetCountryContainer = document.getElementById("country-container");
const form = document.getElementById("guess-country-form");
const resultContainer = document.getElementById("result-container");

let singleCountryData;

function getCountryGuess(formData) {
    const countryGuess = formData.get("country-guess");
    return countryGuess;
}

function isCorrectGuess(countryGuess, countryTarget) {
    if (countryGuess.toLowerCase() === countryTarget.toLowerCase()) {
        return true;
    } else {
        return false;
    }
}

function displayTargetCountry(targetCountryContainer, targetCountry) {
    targetCountryContainer.innerHTML = generateSingleFlagHTML(targetCountry);
}

function generateSingleFlagHTML(country) {
    return `
        <div class="country-item">
            <i class="country-flag-icon">
                ${country.flag}
            </i>
            <p class="country-population">
                Population: ${country.population.toLocaleString()}
            </p>
        </div>
    `
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
    console.log(countryGuess);
    console.log(singleCountryData.name.common);
    if (isCorrectGuess(countryGuess, singleCountryData.name.common)) {
        resultContainer.innerHTML = `
            <p>Right!</p>
        `;
    } else {
        resultContainer.innerHTML = `
            <p>Wrong!</p>
            <p>The right answer was ${singleCountryData.name.common}</p>
        `;
    }

    singleCountryData = await fetchSingleCountry();
    displayTargetCountry(targetCountryContainer, singleCountryData);
});
