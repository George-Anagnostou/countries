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
    const capitalGuess = getGuess(formData);

    let result;
    if (!singleCountryData.capital) {
        singleCountryData.capital = ["None"];
    }
    if (singleCountryData.capital.length > 1) {
        const targetCapitals = singleCountryData.capital.map(capital => {
            return capital.toLowerCase()
        });
        if (targetCapitals.includes(capitalGuess.toLowerCase())) {
            result = "Correct!";
        } else {
            result = "Wrong! The correct answers are: ";
            for (let i = 0; i < singleCountryData.capital.length; i++) {
                result += singleCountryData.capital[i];
                if (singleCountryData.capital.length - i > 1) {
                    result += ", ";
                }
            }
        }
        result.trim();
    } else {
        result = isCorrectGuess(capitalGuess, singleCountryData.capital[0])
                    ? "Correct!"
                    : `Wrong! The correct answer is ${singleCountryData.capital[0]}`
    }

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
