export function generateCountryCard(country) {
    return `
        <div class="country-item">
            <i class="country-flag-icon">
                ${country.flag}
            </i>
            <p class="country-name">
                ${country.name.common}
            </p>
            <p class="country-population">
                Population: ${country.population.toLocaleString()}
            </p>
            <p class="country-flag-description">
                ${country.flags.alt ? country.flags.alt : "No Description Available"}
            </p>
        </div>
    `
};

export function generateCountryCardPartial(country) {
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

export async function fetchCountryData(formData) {
    try {
        const response = await fetch("http://localhost:8000/search_continents", {
            method: "POST",
            body: formData,
        });
        const countryData = await response.json();
        return countryData;
    } catch(e) {
        console.error(e);
    }
};

export async function fetchSingleCountry() {
    try {
        const response = await fetch("http://localhost:8000/country");
        const countryData = await response.json();
        return countryData;
    } catch(e) {
        console.error(e);
    }
};


