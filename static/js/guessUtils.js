import { generateCountryCardPartial } from "./utils.js";

export function getGuess(formData) {
    const guess = formData.get("country-guess");
    return guess;
}

export function isCorrectGuess(guess, target) {
    return guess.toLowerCase() === target.toLowerCase()
}

export function displayTargetCountry(targetCountryContainer, targetCountry) {
    targetCountryContainer.innerHTML = generateCountryCardPartial(targetCountry);
}
